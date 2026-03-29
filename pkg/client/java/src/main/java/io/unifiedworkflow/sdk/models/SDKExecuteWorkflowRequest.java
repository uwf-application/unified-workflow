package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Objects;

/**
 * Request payload for executing a workflow via the SDK endpoint.
 *
 * <p>Construct instances using the nested {@link Builder}:</p>
 * <pre>{@code
 * SDKExecuteWorkflowRequest request = SDKExecuteWorkflowRequest.builder()
 *     .inputData(Map.of("transactionId", "txn-001", "amount", 1500.00))
 *     .callbackUrl("https://example.com/webhook")
 *     .priority(7)
 *     .build();
 * }</pre>
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public final class SDKExecuteWorkflowRequest {

    @JsonProperty("input_data")
    private final Map<String, Object> inputData;

    @JsonProperty("callback_url")
    private final String callbackUrl;

    @JsonProperty("timeout_ms")
    private final long timeoutMs;

    @JsonProperty("wait_for_completion")
    private final boolean waitForCompletion;

    @JsonProperty("metadata")
    private final Map<String, String> metadata;

    @JsonProperty("validation_rules")
    private final List<ValidationRule> validationRules;

    @JsonProperty("enable_validation")
    private final boolean enableValidation;

    @JsonProperty("enable_sanitization")
    private final boolean enableSanitization;

    @JsonProperty("priority")
    private final int priority;

    @JsonProperty("http_request_context")
    private final HTTPRequestContext httpRequestContext;

    @JsonProperty("session_context")
    private final SessionContext sessionContext;

    @JsonProperty("security_context")
    private final SecurityContext securityContext;

    private SDKExecuteWorkflowRequest(Builder builder) {
        this.inputData = builder.inputData != null
                ? Collections.unmodifiableMap(builder.inputData)
                : Collections.emptyMap();
        this.callbackUrl = builder.callbackUrl;
        this.timeoutMs = builder.timeoutMs;
        this.waitForCompletion = builder.waitForCompletion;
        this.metadata = builder.metadata != null
                ? Collections.unmodifiableMap(builder.metadata)
                : Collections.emptyMap();
        this.validationRules = builder.validationRules != null
                ? Collections.unmodifiableList(builder.validationRules)
                : Collections.emptyList();
        this.enableValidation = builder.enableValidation;
        this.enableSanitization = builder.enableSanitization;
        this.priority = builder.priority;
        this.httpRequestContext = builder.httpRequestContext;
        this.sessionContext = builder.sessionContext;
        this.securityContext = builder.securityContext;
    }

    /**
     * Returns a new {@link Builder} for constructing a request.
     *
     * @return a fresh builder instance
     */
    public static Builder builder() {
        return new Builder();
    }

    /** @return the workflow input data map; never null */
    public Map<String, Object> getInputData() { return inputData; }

    /** @return the callback URL to notify on completion; may be null */
    public String getCallbackUrl() { return callbackUrl; }

    /** @return per-request timeout in milliseconds; 0 means use SDK default */
    public long getTimeoutMs() { return timeoutMs; }

    /** @return {@code true} if the call should block until the execution finishes */
    public boolean isWaitForCompletion() { return waitForCompletion; }

    /** @return arbitrary metadata key-value pairs; never null */
    public Map<String, String> getMetadata() { return metadata; }

    /** @return validation rules to run before execution; never null */
    public List<ValidationRule> getValidationRules() { return validationRules; }

    /** @return {@code true} if server-side validation is enabled */
    public boolean isEnableValidation() { return enableValidation; }

    /** @return {@code true} if server-side sanitization is enabled */
    public boolean isEnableSanitization() { return enableSanitization; }

    /** @return execution priority in the range [1, 10] */
    public int getPriority() { return priority; }

    /** @return originating HTTP request context; may be null */
    public HTTPRequestContext getHttpRequestContext() { return httpRequestContext; }

    /** @return caller session context; may be null */
    public SessionContext getSessionContext() { return sessionContext; }

    /** @return caller security context; may be null */
    public SecurityContext getSecurityContext() { return securityContext; }

    // -------------------------------------------------------------------------
    // Builder
    // -------------------------------------------------------------------------

    /**
     * Fluent builder for {@link SDKExecuteWorkflowRequest}.
     */
    public static final class Builder {

        private Map<String, Object> inputData;
        private String callbackUrl;
        private long timeoutMs = 0;
        private boolean waitForCompletion = false;
        private Map<String, String> metadata;
        private List<ValidationRule> validationRules;
        private boolean enableValidation = true;
        private boolean enableSanitization = true;
        private int priority = 5;
        private HTTPRequestContext httpRequestContext;
        private SessionContext sessionContext;
        private SecurityContext securityContext;

        private Builder() {}

        /**
         * Sets the workflow input data.
         *
         * @param inputData input payload; must not be null
         * @return this builder
         */
        public Builder inputData(Map<String, Object> inputData) {
            this.inputData = Objects.requireNonNull(inputData, "inputData must not be null");
            return this;
        }

        /**
         * Sets the callback URL invoked when execution completes.
         *
         * @param callbackUrl webhook URL; may be null
         * @return this builder
         */
        public Builder callbackUrl(String callbackUrl) {
            this.callbackUrl = callbackUrl;
            return this;
        }

        /**
         * Sets a per-request timeout, overriding the SDK-level default.
         *
         * @param timeoutMs timeout in milliseconds; 0 defers to SDK config
         * @return this builder
         */
        public Builder timeoutMs(long timeoutMs) {
            this.timeoutMs = timeoutMs;
            return this;
        }

        /**
         * Requests that the SDK block until the execution reaches a terminal state.
         *
         * @param waitForCompletion {@code true} to block
         * @return this builder
         */
        public Builder waitForCompletion(boolean waitForCompletion) {
            this.waitForCompletion = waitForCompletion;
            return this;
        }

        /**
         * Attaches arbitrary metadata to the execution request.
         *
         * @param metadata key-value metadata map; may be null
         * @return this builder
         */
        public Builder metadata(Map<String, String> metadata) {
            this.metadata = metadata;
            return this;
        }

        /**
         * Sets the validation rules to evaluate against the input data.
         *
         * @param validationRules list of rules; may be null or empty
         * @return this builder
         */
        public Builder validationRules(List<ValidationRule> validationRules) {
            this.validationRules = validationRules;
            return this;
        }

        /**
         * Enables or disables server-side input validation.
         *
         * @param enableValidation {@code true} to enable (default)
         * @return this builder
         */
        public Builder enableValidation(boolean enableValidation) {
            this.enableValidation = enableValidation;
            return this;
        }

        /**
         * Enables or disables server-side input sanitization.
         *
         * @param enableSanitization {@code true} to enable (default)
         * @return this builder
         */
        public Builder enableSanitization(boolean enableSanitization) {
            this.enableSanitization = enableSanitization;
            return this;
        }

        /**
         * Sets the execution priority in the range [1, 10]; higher values are processed first.
         *
         * @param priority execution priority; defaults to 5
         * @return this builder
         */
        public Builder priority(int priority) {
            this.priority = priority;
            return this;
        }

        /**
         * Attaches originating HTTP request context for audit purposes.
         *
         * @param httpRequestContext context; may be null
         * @return this builder
         */
        public Builder httpRequestContext(HTTPRequestContext httpRequestContext) {
            this.httpRequestContext = httpRequestContext;
            return this;
        }

        /**
         * Attaches caller session context.
         *
         * @param sessionContext context; may be null
         * @return this builder
         */
        public Builder sessionContext(SessionContext sessionContext) {
            this.sessionContext = sessionContext;
            return this;
        }

        /**
         * Attaches caller security context.
         *
         * @param securityContext context; may be null
         * @return this builder
         */
        public Builder securityContext(SecurityContext securityContext) {
            this.securityContext = securityContext;
            return this;
        }

        /**
         * Builds and returns the request instance.
         *
         * @return a new immutable {@link SDKExecuteWorkflowRequest}
         */
        public SDKExecuteWorkflowRequest build() {
            return new SDKExecuteWorkflowRequest(this);
        }
    }
}
