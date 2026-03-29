package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Objects;

/**
 * Configuration for a registered webhook endpoint.
 *
 * <p>Construct instances with the nested {@link Builder}:</p>
 * <pre>{@code
 * WebhookConfiguration webhook = WebhookConfiguration.builder()
 *     .url("https://example.com/hooks/workflow")
 *     .events(List.of("workflow_completed", "workflow_failed"))
 *     .secret("s3cr3t")
 *     .build();
 * }</pre>
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public final class WebhookConfiguration {

    @JsonProperty("webhook_id")
    private final String webhookId;

    @JsonProperty("url")
    private final String url;

    @JsonProperty("events")
    private final List<String> events;

    @JsonProperty("secret")
    private final String secret;

    @JsonProperty("enabled")
    private final boolean enabled;

    @JsonProperty("retry_count")
    private final int retryCount;

    @JsonProperty("timeout_ms")
    private final long timeoutMs;

    @JsonProperty("headers")
    private final Map<String, String> headers;

    private WebhookConfiguration(Builder builder) {
        this.webhookId = builder.webhookId;
        this.url = builder.url;
        this.events = builder.events != null
                ? Collections.unmodifiableList(builder.events)
                : Collections.emptyList();
        this.secret = builder.secret;
        this.enabled = builder.enabled;
        this.retryCount = builder.retryCount;
        this.timeoutMs = builder.timeoutMs;
        this.headers = builder.headers != null
                ? Collections.unmodifiableMap(builder.headers)
                : Collections.emptyMap();
    }

    /**
     * Returns a new {@link Builder}.
     *
     * @return a fresh builder instance
     */
    public static Builder builder() {
        return new Builder();
    }

    /** @return server-assigned webhook identifier; null before registration */
    public String getWebhookId() { return webhookId; }

    /** @return the URL that receives webhook HTTP POST requests; must not be null */
    public String getUrl() { return url; }

    /** @return list of event types that trigger this webhook; never null */
    public List<String> getEvents() { return events; }

    /** @return HMAC signing secret used for payload verification; may be null */
    public String getSecret() { return secret; }

    /** @return {@code true} if the webhook is active */
    public boolean isEnabled() { return enabled; }

    /** @return number of delivery retry attempts before giving up */
    public int getRetryCount() { return retryCount; }

    /** @return per-delivery timeout in milliseconds */
    public long getTimeoutMs() { return timeoutMs; }

    /** @return additional HTTP headers sent with each webhook delivery; never null */
    public Map<String, String> getHeaders() { return headers; }

    // -------------------------------------------------------------------------
    // Builder
    // -------------------------------------------------------------------------

    /**
     * Fluent builder for {@link WebhookConfiguration}.
     */
    public static final class Builder {

        private String webhookId;
        private String url;
        private List<String> events;
        private String secret;
        private boolean enabled = true;
        private int retryCount = 3;
        private long timeoutMs = 5000;
        private Map<String, String> headers;

        private Builder() {}

        /**
         * Sets the server-assigned webhook ID (required when updating an existing webhook).
         *
         * @param webhookId identifier; may be null for new registrations
         * @return this builder
         */
        public Builder webhookId(String webhookId) {
            this.webhookId = webhookId;
            return this;
        }

        /**
         * Sets the target URL.
         *
         * @param url delivery URL; must not be null
         * @return this builder
         */
        public Builder url(String url) {
            this.url = Objects.requireNonNull(url, "url must not be null");
            return this;
        }

        /**
         * Sets the event types that trigger this webhook.
         *
         * @param events list of event type strings; must not be null
         * @return this builder
         */
        public Builder events(List<String> events) {
            this.events = Objects.requireNonNull(events, "events must not be null");
            return this;
        }

        /**
         * Sets the HMAC signing secret.
         *
         * @param secret signing secret; may be null
         * @return this builder
         */
        public Builder secret(String secret) {
            this.secret = secret;
            return this;
        }

        /**
         * Controls whether the webhook is active.
         *
         * @param enabled {@code true} to enable (default)
         * @return this builder
         */
        public Builder enabled(boolean enabled) {
            this.enabled = enabled;
            return this;
        }

        /**
         * Sets the number of delivery retry attempts.
         *
         * @param retryCount retry count; defaults to 3
         * @return this builder
         */
        public Builder retryCount(int retryCount) {
            this.retryCount = retryCount;
            return this;
        }

        /**
         * Sets the per-delivery timeout in milliseconds.
         *
         * @param timeoutMs timeout; defaults to 5000
         * @return this builder
         */
        public Builder timeoutMs(long timeoutMs) {
            this.timeoutMs = timeoutMs;
            return this;
        }

        /**
         * Sets additional HTTP headers sent with each delivery.
         *
         * @param headers header map; may be null
         * @return this builder
         */
        public Builder headers(Map<String, String> headers) {
            this.headers = headers;
            return this;
        }

        /**
         * Builds and returns the configuration.
         *
         * @return a new immutable {@link WebhookConfiguration}
         */
        public WebhookConfiguration build() {
            Objects.requireNonNull(url, "url is required");
            return new WebhookConfiguration(this);
        }
    }
}
