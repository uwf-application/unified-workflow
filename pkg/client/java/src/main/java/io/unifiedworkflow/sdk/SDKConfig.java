package io.unifiedworkflow.sdk;

import io.unifiedworkflow.sdk.errors.WorkflowSDKException;
import io.unifiedworkflow.sdk.errors.ErrorCode;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;
import java.util.Objects;

/**
 * Immutable configuration record for the Unified Workflow SDK.
 *
 * <p>Use {@link #builder()} to create instances, or the static factory methods
 * {@link #fromYaml(Path)} and {@link #fromEnvironment()} for environment-specific configuration.</p>
 *
 * <pre>{@code
 * SDKConfig config = SDKConfig.builder()
 *     .workflowApiEndpoint("https://workflow.internal:8080")
 *     .authToken(System.getenv("WORKFLOW_TOKEN"))
 *     .authType(SDKConfig.AuthType.BEARER_TOKEN)
 *     .timeout(Duration.ofSeconds(60))
 *     .build();
 * }</pre>
 *
 * @param workflowApiEndpoint    base URL of the workflow API; defaults to {@code http://localhost:8080}
 * @param timeout                per-request timeout; defaults to 30 seconds
 * @param maxRetries             maximum number of retry attempts for retryable errors; defaults to 3
 * @param retryDelay             base delay between retries (doubles with each attempt); defaults to 1 second
 * @param authToken              bearer token or API key value; may be null when auth is disabled
 * @param authType               authentication scheme; defaults to {@link AuthType#NONE}
 * @param enableValidation       send validation rules to the server; defaults to {@code true}
 * @param enableSanitization     request server-side input sanitization; defaults to {@code true}
 * @param strictValidation       treat validation warnings as errors; defaults to {@code false}
 * @param asyncExecution         submit executions asynchronously; defaults to {@code true}
 * @param pollIntervalMs         polling interval for {@code waitForCompletion}; defaults to 2000 ms
 * @param defaultPriority        default execution priority [1–10]; defaults to 5
 * @param enableCircuitBreaker   enable the client-side circuit breaker; defaults to {@code true}
 * @param circuitBreakerThreshold consecutive failure count that trips the breaker; defaults to 5
 * @param circuitBreakerTimeout  how long the breaker stays open before attempting recovery;
 *                               defaults to 60 seconds
 * @param enableRequestLogging   log outgoing requests and responses at DEBUG level; defaults to {@code false}
 */
public record SDKConfig(
        String workflowApiEndpoint,
        Duration timeout,
        int maxRetries,
        Duration retryDelay,
        String authToken,
        AuthType authType,
        boolean enableValidation,
        boolean enableSanitization,
        boolean strictValidation,
        boolean asyncExecution,
        long pollIntervalMs,
        int defaultPriority,
        boolean enableCircuitBreaker,
        int circuitBreakerThreshold,
        Duration circuitBreakerTimeout,
        boolean enableRequestLogging
) {

    // -------------------------------------------------------------------------
    // AuthType enum
    // -------------------------------------------------------------------------

    /**
     * Authentication scheme used when making requests to the workflow API.
     */
    public enum AuthType {
        /** Attach a {@code Authorization: Bearer <token>} header. */
        BEARER_TOKEN,
        /** Attach an {@code X-API-Key: <key>} header. */
        API_KEY,
        /** Attach an {@code Authorization: Basic <base64>} header. */
        BASIC_AUTH,
        /** No authentication header is attached. */
        NONE
    }

    // -------------------------------------------------------------------------
    // Compact constructor – defaults validation
    // -------------------------------------------------------------------------

    /**
     * Compact constructor that validates mandatory fields and normalises the endpoint.
     */
    public SDKConfig {
        Objects.requireNonNull(workflowApiEndpoint, "workflowApiEndpoint must not be null");
        Objects.requireNonNull(timeout, "timeout must not be null");
        Objects.requireNonNull(retryDelay, "retryDelay must not be null");
        Objects.requireNonNull(authType, "authType must not be null");
        Objects.requireNonNull(circuitBreakerTimeout, "circuitBreakerTimeout must not be null");

        // Strip trailing slash for consistent URL construction
        if (workflowApiEndpoint.endsWith("/")) {
            workflowApiEndpoint = workflowApiEndpoint.substring(0, workflowApiEndpoint.length() - 1);
        }
    }

    // -------------------------------------------------------------------------
    // Static factory: builder
    // -------------------------------------------------------------------------

    /**
     * Returns a new {@link SDKConfigBuilder} pre-filled with the default values.
     *
     * @return a fresh builder instance
     */
    public static SDKConfigBuilder builder() {
        return new SDKConfigBuilder();
    }

    // -------------------------------------------------------------------------
    // Static factory: fromEnvironment
    // -------------------------------------------------------------------------

    /**
     * Creates a configuration populated from well-known environment variables.
     *
     * <table>
     *   <caption>Supported environment variables</caption>
     *   <tr><th>Variable</th><th>Maps to</th><th>Default</th></tr>
     *   <tr><td>{@code WORKFLOW_API_ENDPOINT}</td><td>workflowApiEndpoint</td><td>http://localhost:8080</td></tr>
     *   <tr><td>{@code WORKFLOW_AUTH_TOKEN}</td><td>authToken</td><td>(none)</td></tr>
     *   <tr><td>{@code WORKFLOW_AUTH_TYPE}</td><td>authType</td><td>NONE</td></tr>
     *   <tr><td>{@code WORKFLOW_TIMEOUT_MS}</td><td>timeout (milliseconds)</td><td>30000</td></tr>
     *   <tr><td>{@code WORKFLOW_MAX_RETRIES}</td><td>maxRetries</td><td>3</td></tr>
     *   <tr><td>{@code WORKFLOW_ENABLE_REQUEST_LOGGING}</td><td>enableRequestLogging</td><td>false</td></tr>
     * </table>
     *
     * @return a configuration instance populated from the environment
     */
    public static SDKConfig fromEnvironment() {
        SDKConfigBuilder builder = new SDKConfigBuilder();

        String endpoint = System.getenv("WORKFLOW_API_ENDPOINT");
        if (endpoint != null && !endpoint.isBlank()) {
            builder.workflowApiEndpoint(endpoint);
        }

        String token = System.getenv("WORKFLOW_AUTH_TOKEN");
        if (token != null && !token.isBlank()) {
            builder.authToken(token);
        }

        String authTypeStr = System.getenv("WORKFLOW_AUTH_TYPE");
        if (authTypeStr != null && !authTypeStr.isBlank()) {
            try {
                builder.authType(AuthType.valueOf(authTypeStr.toUpperCase()));
            } catch (IllegalArgumentException ignored) {
                // fall back to default NONE
            }
        }

        String timeoutMs = System.getenv("WORKFLOW_TIMEOUT_MS");
        if (timeoutMs != null && !timeoutMs.isBlank()) {
            try {
                builder.timeout(Duration.ofMillis(Long.parseLong(timeoutMs)));
            } catch (NumberFormatException ignored) {
                // fall back to default
            }
        }

        String maxRetries = System.getenv("WORKFLOW_MAX_RETRIES");
        if (maxRetries != null && !maxRetries.isBlank()) {
            try {
                builder.maxRetries(Integer.parseInt(maxRetries));
            } catch (NumberFormatException ignored) {
                // fall back to default
            }
        }

        String logging = System.getenv("WORKFLOW_ENABLE_REQUEST_LOGGING");
        if ("true".equalsIgnoreCase(logging)) {
            builder.enableRequestLogging(true);
        }

        return builder.build();
    }

    // -------------------------------------------------------------------------
    // Static factory: fromYaml
    // -------------------------------------------------------------------------

    /**
     * Creates a configuration by reading a YAML file.
     *
     * <p>The YAML file must follow this structure:</p>
     * <pre>
     * workflowApiEndpoint: "http://workflow-api:8080"
     * authToken: "my-token"
     * authType: "BEARER_TOKEN"
     * timeoutMs: 30000
     * maxRetries: 3
     * enableRequestLogging: false
     * </pre>
     *
     * <p>This implementation uses a simple line-by-line key-value parser so that
     * the SDK has no dependency on a YAML library. Jackson's YAML module is
     * intentionally not required.</p>
     *
     * @param path path to the YAML configuration file; must not be null
     * @return a configuration instance populated from the file
     * @throws WorkflowSDKException if the file cannot be read or parsed
     */
    public static SDKConfig fromYaml(Path path) {
        Objects.requireNonNull(path, "path must not be null");

        SDKConfigBuilder builder = new SDKConfigBuilder();

        try {
            for (String line : Files.readAllLines(path)) {
                String trimmed = line.trim();
                if (trimmed.isEmpty() || trimmed.startsWith("#")) {
                    continue;
                }
                int colon = trimmed.indexOf(':');
                if (colon < 0) {
                    continue;
                }
                String key = trimmed.substring(0, colon).trim();
                String value = trimmed.substring(colon + 1).trim()
                        .replaceAll("^[\"']|[\"']$", ""); // strip quotes

                switch (key) {
                    case "workflowApiEndpoint" -> builder.workflowApiEndpoint(value);
                    case "authToken" -> builder.authToken(value);
                    case "authType" -> {
                        try {
                            builder.authType(AuthType.valueOf(value.toUpperCase()));
                        } catch (IllegalArgumentException ignored) {}
                    }
                    case "timeoutMs" -> {
                        try {
                            builder.timeout(Duration.ofMillis(Long.parseLong(value)));
                        } catch (NumberFormatException ignored) {}
                    }
                    case "maxRetries" -> {
                        try {
                            builder.maxRetries(Integer.parseInt(value));
                        } catch (NumberFormatException ignored) {}
                    }
                    case "retryDelayMs" -> {
                        try {
                            builder.retryDelay(Duration.ofMillis(Long.parseLong(value)));
                        } catch (NumberFormatException ignored) {}
                    }
                    case "enableValidation" ->
                            builder.enableValidation(Boolean.parseBoolean(value));
                    case "enableSanitization" ->
                            builder.enableSanitization(Boolean.parseBoolean(value));
                    case "strictValidation" ->
                            builder.strictValidation(Boolean.parseBoolean(value));
                    case "asyncExecution" ->
                            builder.asyncExecution(Boolean.parseBoolean(value));
                    case "pollIntervalMs" -> {
                        try {
                            builder.pollIntervalMs(Long.parseLong(value));
                        } catch (NumberFormatException ignored) {}
                    }
                    case "defaultPriority" -> {
                        try {
                            builder.defaultPriority(Integer.parseInt(value));
                        } catch (NumberFormatException ignored) {}
                    }
                    case "enableCircuitBreaker" ->
                            builder.enableCircuitBreaker(Boolean.parseBoolean(value));
                    case "circuitBreakerThreshold" -> {
                        try {
                            builder.circuitBreakerThreshold(Integer.parseInt(value));
                        } catch (NumberFormatException ignored) {}
                    }
                    case "circuitBreakerTimeoutMs" -> {
                        try {
                            builder.circuitBreakerTimeout(Duration.ofMillis(Long.parseLong(value)));
                        } catch (NumberFormatException ignored) {}
                    }
                    case "enableRequestLogging" ->
                            builder.enableRequestLogging(Boolean.parseBoolean(value));
                    default -> { /* unknown keys are silently ignored */ }
                }
            }
        } catch (IOException e) {
            throw new WorkflowSDKException(ErrorCode.INVALID_CONFIG,
                    "Failed to read SDK configuration file: " + path, e);
        }

        return builder.build();
    }
}
