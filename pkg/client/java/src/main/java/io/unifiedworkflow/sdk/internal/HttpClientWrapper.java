package io.unifiedworkflow.sdk.internal;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import io.unifiedworkflow.sdk.SDKConfig;
import io.unifiedworkflow.sdk.errors.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.time.Instant;
import java.util.Base64;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;

/**
 * Internal HTTP client wrapper providing retry logic, exponential back-off,
 * circuit-breaker protection, and auth-header injection.
 *
 * <p>This class is <em>not</em> part of the public SDK API and may change without notice.
 * All methods are thread-safe.</p>
 */
public final class HttpClientWrapper implements AutoCloseable {

    private static final Logger log = LoggerFactory.getLogger(HttpClientWrapper.class);

    /** HTTP status codes that are considered transient and eligible for retry. */
    private static final Set<Integer> RETRYABLE_STATUS_CODES = Set.of(408, 429, 500, 502, 503, 504);

    // -------------------------------------------------------------------------
    // Circuit-breaker state machine
    // -------------------------------------------------------------------------

    private enum CircuitState { CLOSED, OPEN, HALF_OPEN }

    private final AtomicReference<CircuitState> circuitState =
            new AtomicReference<>(CircuitState.CLOSED);
    private final AtomicInteger failureCount = new AtomicInteger(0);
    private volatile Instant circuitOpenedAt = null;

    // -------------------------------------------------------------------------
    // Fields
    // -------------------------------------------------------------------------

    private final SDKConfig config;
    private final HttpClient httpClient;
    private final ObjectMapper mapper;

    /**
     * Constructs a wrapper around a new {@link HttpClient} backed by virtual threads.
     *
     * @param config SDK configuration; must not be null
     */
    public HttpClientWrapper(SDKConfig config) {
        this.config = config;
        this.httpClient = HttpClient.newBuilder()
                .executor(Executors.newVirtualThreadPerTaskExecutor())
                .connectTimeout(config.timeout())
                .build();
        this.mapper = new ObjectMapper()
                .registerModule(new JavaTimeModule())
                .disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
    }

    // -------------------------------------------------------------------------
    // Public API
    // -------------------------------------------------------------------------

    /**
     * Performs a GET request and deserialises the JSON response body into {@code responseType}.
     *
     * @param url          absolute URL to call
     * @param responseType Jackson TypeReference for the target type
     * @param <T>          response type
     * @return deserialised response body
     * @throws WorkflowSDKException on any error
     */
    public <T> T get(String url, TypeReference<T> responseType) {
        return execute("GET", url, null, responseType);
    }

    /**
     * Performs a POST request with a JSON body and deserialises the response.
     *
     * @param url          absolute URL to call
     * @param body         request body object; serialised to JSON; may be null for empty body
     * @param responseType Jackson TypeReference for the target type
     * @param <T>          response type
     * @return deserialised response body
     * @throws WorkflowSDKException on any error
     */
    public <T> T post(String url, Object body, TypeReference<T> responseType) {
        return execute("POST", url, body, responseType);
    }

    /**
     * Performs a DELETE request and returns {@code true} if the server responded with 2xx.
     *
     * @param url absolute URL to call
     * @return {@code true} on success
     * @throws WorkflowSDKException on any error
     */
    public boolean delete(String url) {
        execute("DELETE", url, null, new TypeReference<Void>() {});
        return true;
    }

    /**
     * Returns {@code true} if the circuit breaker is currently open (i.e. requests are blocked).
     *
     * @return {@code true} when the circuit is open
     */
    public boolean isCircuitOpen() {
        return circuitState.get() == CircuitState.OPEN;
    }

    // -------------------------------------------------------------------------
    // Core execution with retry and circuit-breaker
    // -------------------------------------------------------------------------

    private <T> T execute(String method, String url, Object body, TypeReference<T> responseType) {
        checkCircuitBreaker();

        String requestBody = serializeBody(body);

        WorkflowSDKException lastException = null;
        int attempts = config.maxRetries() + 1;

        for (int attempt = 0; attempt < attempts; attempt++) {
            if (attempt > 0) {
                long backoffMs = config.retryDelay().toMillis() * (1L << Math.min(attempt - 1, 5));
                sleep(backoffMs);
            }

            try {
                T result = doRequest(method, url, requestBody, responseType);
                recordSuccess();
                return result;
            } catch (WorkflowSDKException ex) {
                lastException = ex;

                if (!shouldRetry(ex)) {
                    recordFailure(ex);
                    throw ex;
                }

                log.debug("Request to {} failed (attempt {}/{}): {}",
                        url, attempt + 1, attempts, ex.getMessage());
            }
        }

        recordFailure(lastException);
        throw new WorkflowSDKException(ErrorCode.RETRY_EXHAUSTED,
                "All " + config.maxRetries() + " retry attempts failed for " + url,
                lastException);
    }

    private <T> T doRequest(String method, String url,
                            String bodyStr, TypeReference<T> responseType) {
        HttpRequest.Builder reqBuilder = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .timeout(config.timeout())
                .header("Content-Type", "application/json")
                .header("Accept", "application/json");

        injectAuthHeader(reqBuilder);

        HttpRequest.BodyPublisher publisher = bodyStr != null
                ? HttpRequest.BodyPublishers.ofString(bodyStr)
                : HttpRequest.BodyPublishers.noBody();

        reqBuilder.method(method, publisher);
        HttpRequest request = reqBuilder.build();

        if (config.enableRequestLogging()) {
            log.debug("--> {} {}", method, url);
            if (bodyStr != null) {
                log.debug("    Body: {}", bodyStr);
            }
        }

        HttpResponse<String> response;
        try {
            response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        } catch (java.net.http.HttpTimeoutException e) {
            throw new TimeoutException("Request timed out: " + method + " " + url, e);
        } catch (IOException | InterruptedException e) {
            if (e instanceof InterruptedException) {
                Thread.currentThread().interrupt();
            }
            throw new NetworkException("Network error: " + method + " " + url, e);
        }

        if (config.enableRequestLogging()) {
            log.debug("<-- {} {} ({})", method, url, response.statusCode());
        }

        return handleResponse(response, responseType);
    }

    // -------------------------------------------------------------------------
    // Response handling
    // -------------------------------------------------------------------------

    private <T> T handleResponse(HttpResponse<String> response, TypeReference<T> responseType) {
        int status = response.statusCode();

        if (status >= 200 && status < 300) {
            if (status == 204 || response.body() == null || response.body().isBlank()) {
                return null;
            }
            try {
                return mapper.readValue(response.body(), responseType);
            } catch (IOException e) {
                throw new WorkflowSDKException(ErrorCode.UNKNOWN,
                        "Failed to deserialise server response", e);
            }
        }

        Map<String, Object> errorBody = parseErrorBody(response.body());

        throw switch (status) {
            case 400 -> new ValidationException(
                    extractMessage(errorBody, "Validation failed"),
                    java.util.Collections.emptyList());
            case 401 -> new AuthenticationException(
                    extractMessage(errorBody, "Authentication failed"));
            case 404 -> {
                String message = extractMessage(errorBody, "Resource not found");
                if (errorBody.containsKey("workflow_id")) {
                    yield new WorkflowNotFoundException(
                            String.valueOf(errorBody.get("workflow_id")), message, errorBody);
                }
                if (errorBody.containsKey("run_id")) {
                    yield new ExecutionNotFoundException(
                            String.valueOf(errorBody.get("run_id")), message, errorBody);
                }
                yield new WorkflowNotFoundException("unknown", message, errorBody);
            }
            case 429 -> {
                Duration retryAfter = parseRetryAfter(response);
                yield new RateLimitException(
                        extractMessage(errorBody, "Rate limit exceeded"), retryAfter);
            }
            default -> new WorkflowSDKException(ErrorCode.UNKNOWN,
                    extractMessage(errorBody, "HTTP " + status), status, errorBody);
        };
    }

    // -------------------------------------------------------------------------
    // Circuit breaker
    // -------------------------------------------------------------------------

    private void checkCircuitBreaker() {
        if (!config.enableCircuitBreaker()) {
            return;
        }

        CircuitState state = circuitState.get();
        switch (state) {
            case CLOSED -> { /* normal operation */ }
            case OPEN -> {
                Duration elapsed = Duration.between(circuitOpenedAt, Instant.now());
                if (elapsed.compareTo(config.circuitBreakerTimeout()) >= 0) {
                    log.info("Circuit breaker transitioning OPEN -> HALF_OPEN");
                    circuitState.compareAndSet(CircuitState.OPEN, CircuitState.HALF_OPEN);
                } else {
                    throw new WorkflowSDKException(ErrorCode.CIRCUIT_BREAKER_OPEN,
                            "Circuit breaker is open; requests are temporarily blocked");
                }
            }
            case HALF_OPEN -> { /* allow one probe request */ }
        }
    }

    private void recordSuccess() {
        if (!config.enableCircuitBreaker()) return;

        CircuitState previous = circuitState.get();
        if (previous == CircuitState.HALF_OPEN) {
            log.info("Circuit breaker transitioning HALF_OPEN -> CLOSED");
        }
        circuitState.set(CircuitState.CLOSED);
        failureCount.set(0);
    }

    private void recordFailure(WorkflowSDKException ex) {
        if (!config.enableCircuitBreaker()) return;
        if (ex != null && !ex.isRetryable()) return;

        int failures = failureCount.incrementAndGet();
        if (failures >= config.circuitBreakerThreshold()
                && circuitState.get() != CircuitState.OPEN) {
            log.warn("Circuit breaker threshold reached ({} failures); transitioning to OPEN",
                    failures);
            circuitState.set(CircuitState.OPEN);
            circuitOpenedAt = Instant.now();
        }
    }

    // -------------------------------------------------------------------------
    // Helpers
    // -------------------------------------------------------------------------

    private void injectAuthHeader(HttpRequest.Builder builder) {
        if (config.authToken() == null || config.authType() == SDKConfig.AuthType.NONE) {
            return;
        }
        switch (config.authType()) {
            case BEARER_TOKEN ->
                    builder.header("Authorization", "Bearer " + config.authToken());
            case API_KEY ->
                    builder.header("X-API-Key", config.authToken());
            case BASIC_AUTH -> {
                String encoded = Base64.getEncoder()
                        .encodeToString(config.authToken().getBytes());
                builder.header("Authorization", "Basic " + encoded);
            }
            default -> { /* NONE */ }
        }
    }

    private String serializeBody(Object body) {
        if (body == null) return null;
        try {
            return mapper.writeValueAsString(body);
        } catch (IOException e) {
            throw new WorkflowSDKException(ErrorCode.UNKNOWN,
                    "Failed to serialise request body", e);
        }
    }

    private Map<String, Object> parseErrorBody(String body) {
        if (body == null || body.isBlank()) {
            return Map.of();
        }
        try {
            return mapper.readValue(body, new TypeReference<Map<String, Object>>() {});
        } catch (IOException e) {
            return Map.of("message", body);
        }
    }

    private String extractMessage(Map<String, Object> errorBody, String defaultMessage) {
        Object msg = errorBody.get("error");
        if (msg == null) msg = errorBody.get("message");
        return msg != null ? String.valueOf(msg) : defaultMessage;
    }

    private Duration parseRetryAfter(HttpResponse<?> response) {
        return response.headers().firstValue("Retry-After")
                .map(v -> {
                    try {
                        return Duration.ofSeconds(Long.parseLong(v));
                    } catch (NumberFormatException e) {
                        return Duration.ofSeconds(60);
                    }
                })
                .orElse(Duration.ZERO);
    }

    private boolean shouldRetry(WorkflowSDKException ex) {
        if (!ex.isRetryable()) return false;
        if (ex.getHttpStatus() != 0 && !RETRYABLE_STATUS_CODES.contains(ex.getHttpStatus())) {
            return false;
        }
        return true;
    }

    private void sleep(long ms) {
        try {
            Thread.sleep(ms);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new WorkflowSDKException(ErrorCode.NETWORK_ERROR,
                    "Interrupted during retry back-off", e);
        }
    }

    @Override
    public void close() {
        // HttpClient does not require explicit closure, but we reset circuit state
        circuitState.set(CircuitState.CLOSED);
        failureCount.set(0);
    }
}
