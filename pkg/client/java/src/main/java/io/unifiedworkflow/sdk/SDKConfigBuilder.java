package io.unifiedworkflow.sdk;

import java.time.Duration;
import java.util.Objects;

/**
 * Fluent builder for {@link SDKConfig}.
 *
 * <p>All fields have defaults matching the recommended production settings. Override
 * only the values you need to change:</p>
 * <pre>{@code
 * SDKConfig config = SDKConfig.builder()
 *     .workflowApiEndpoint("https://workflow.example.com")
 *     .authToken(token)
 *     .authType(SDKConfig.AuthType.BEARER_TOKEN)
 *     .build();
 * }</pre>
 */
public final class SDKConfigBuilder {

    private String workflowApiEndpoint = "http://localhost:8080";
    private Duration timeout = Duration.ofSeconds(30);
    private int maxRetries = 3;
    private Duration retryDelay = Duration.ofSeconds(1);
    private String authToken = null;
    private SDKConfig.AuthType authType = SDKConfig.AuthType.NONE;
    private boolean enableValidation = true;
    private boolean enableSanitization = true;
    private boolean strictValidation = false;
    private boolean asyncExecution = true;
    private long pollIntervalMs = 2000L;
    private int defaultPriority = 5;
    private boolean enableCircuitBreaker = true;
    private int circuitBreakerThreshold = 5;
    private Duration circuitBreakerTimeout = Duration.ofSeconds(60);
    private boolean enableRequestLogging = false;

    /** Package-private constructor — use {@link SDKConfig#builder()} */
    SDKConfigBuilder() {}

    /**
     * Sets the base URL for the workflow API.
     *
     * @param workflowApiEndpoint endpoint URL; must not be null
     * @return this builder
     */
    public SDKConfigBuilder workflowApiEndpoint(String workflowApiEndpoint) {
        this.workflowApiEndpoint = Objects.requireNonNull(workflowApiEndpoint,
                "workflowApiEndpoint must not be null");
        return this;
    }

    /**
     * Sets the per-request HTTP timeout.
     *
     * @param timeout timeout duration; must not be null
     * @return this builder
     */
    public SDKConfigBuilder timeout(Duration timeout) {
        this.timeout = Objects.requireNonNull(timeout, "timeout must not be null");
        return this;
    }

    /**
     * Sets the maximum number of retry attempts for retryable errors.
     *
     * @param maxRetries retry count; must be &gt;= 0
     * @return this builder
     */
    public SDKConfigBuilder maxRetries(int maxRetries) {
        if (maxRetries < 0) throw new IllegalArgumentException("maxRetries must be >= 0");
        this.maxRetries = maxRetries;
        return this;
    }

    /**
     * Sets the base delay between retry attempts. Retries use exponential back-off starting
     * from this value.
     *
     * @param retryDelay base delay; must not be null
     * @return this builder
     */
    public SDKConfigBuilder retryDelay(Duration retryDelay) {
        this.retryDelay = Objects.requireNonNull(retryDelay, "retryDelay must not be null");
        return this;
    }

    /**
     * Sets the authentication token (bearer token or API key value).
     *
     * @param authToken token value; may be null to disable authentication
     * @return this builder
     */
    public SDKConfigBuilder authToken(String authToken) {
        this.authToken = authToken;
        return this;
    }

    /**
     * Sets the authentication scheme.
     *
     * @param authType auth type; must not be null
     * @return this builder
     */
    public SDKConfigBuilder authType(SDKConfig.AuthType authType) {
        this.authType = Objects.requireNonNull(authType, "authType must not be null");
        return this;
    }

    /**
     * Enables or disables server-side input validation.
     *
     * @param enableValidation {@code true} to enable (default)
     * @return this builder
     */
    public SDKConfigBuilder enableValidation(boolean enableValidation) {
        this.enableValidation = enableValidation;
        return this;
    }

    /**
     * Enables or disables server-side input sanitization.
     *
     * @param enableSanitization {@code true} to enable (default)
     * @return this builder
     */
    public SDKConfigBuilder enableSanitization(boolean enableSanitization) {
        this.enableSanitization = enableSanitization;
        return this;
    }

    /**
     * When {@code true}, validation warnings are treated as errors.
     *
     * @param strictValidation strict mode flag; defaults to {@code false}
     * @return this builder
     */
    public SDKConfigBuilder strictValidation(boolean strictValidation) {
        this.strictValidation = strictValidation;
        return this;
    }

    /**
     * When {@code true}, executions are submitted without waiting for completion.
     *
     * @param asyncExecution async flag; defaults to {@code true}
     * @return this builder
     */
    public SDKConfigBuilder asyncExecution(boolean asyncExecution) {
        this.asyncExecution = asyncExecution;
        return this;
    }

    /**
     * Sets the polling interval used by {@code waitForCompletion}.
     *
     * @param pollIntervalMs interval in milliseconds; defaults to 2000
     * @return this builder
     */
    public SDKConfigBuilder pollIntervalMs(long pollIntervalMs) {
        if (pollIntervalMs <= 0) throw new IllegalArgumentException("pollIntervalMs must be > 0");
        this.pollIntervalMs = pollIntervalMs;
        return this;
    }

    /**
     * Sets the default execution priority applied when the request does not specify one.
     *
     * @param defaultPriority priority in [1, 10]; defaults to 5
     * @return this builder
     */
    public SDKConfigBuilder defaultPriority(int defaultPriority) {
        if (defaultPriority < 1 || defaultPriority > 10) {
            throw new IllegalArgumentException("defaultPriority must be in [1, 10]");
        }
        this.defaultPriority = defaultPriority;
        return this;
    }

    /**
     * Enables or disables the client-side circuit breaker.
     *
     * @param enableCircuitBreaker {@code true} to enable (default)
     * @return this builder
     */
    public SDKConfigBuilder enableCircuitBreaker(boolean enableCircuitBreaker) {
        this.enableCircuitBreaker = enableCircuitBreaker;
        return this;
    }

    /**
     * Sets the consecutive failure count that trips the circuit breaker.
     *
     * @param circuitBreakerThreshold failure threshold; must be &gt; 0
     * @return this builder
     */
    public SDKConfigBuilder circuitBreakerThreshold(int circuitBreakerThreshold) {
        if (circuitBreakerThreshold <= 0) {
            throw new IllegalArgumentException("circuitBreakerThreshold must be > 0");
        }
        this.circuitBreakerThreshold = circuitBreakerThreshold;
        return this;
    }

    /**
     * Sets how long the circuit breaker stays open before attempting recovery.
     *
     * @param circuitBreakerTimeout open duration; must not be null
     * @return this builder
     */
    public SDKConfigBuilder circuitBreakerTimeout(Duration circuitBreakerTimeout) {
        this.circuitBreakerTimeout = Objects.requireNonNull(circuitBreakerTimeout,
                "circuitBreakerTimeout must not be null");
        return this;
    }

    /**
     * Enables or disables DEBUG-level request/response logging.
     *
     * @param enableRequestLogging {@code true} to enable; defaults to {@code false}
     * @return this builder
     */
    public SDKConfigBuilder enableRequestLogging(boolean enableRequestLogging) {
        this.enableRequestLogging = enableRequestLogging;
        return this;
    }

    /**
     * Builds and returns an immutable {@link SDKConfig}.
     *
     * @return a new {@code SDKConfig} with the configured values
     */
    public SDKConfig build() {
        return new SDKConfig(
                workflowApiEndpoint,
                timeout,
                maxRetries,
                retryDelay,
                authToken,
                authType,
                enableValidation,
                enableSanitization,
                strictValidation,
                asyncExecution,
                pollIntervalMs,
                defaultPriority,
                enableCircuitBreaker,
                circuitBreakerThreshold,
                circuitBreakerTimeout,
                enableRequestLogging
        );
    }
}
