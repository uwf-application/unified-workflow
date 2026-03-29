package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.Collections;
import java.util.Map;

/**
 * Captures the full HTTP request context that triggered a workflow execution.
 *
 * <p>Populating this context allows the workflow engine to record audit trails,
 * apply geo-based policies, and correlate executions with originating requests.</p>
 *
 * @param method      HTTP method (e.g. {@code "POST"})
 * @param path        URL path (e.g. {@code "/api/payments"})
 * @param headers     HTTP headers as a name-to-value map; never null
 * @param queryParams URL query parameters as a name-to-value map; never null
 * @param pathParams  URL path parameters as a name-to-value map; never null
 * @param body        raw request body as a string; may be null
 * @param remoteAddr  client IP address; may be null
 * @param userAgent   value of the {@code User-Agent} header; may be null
 * @param timestamp   wall-clock time of the originating request
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record HTTPRequestContext(
        @JsonProperty("method") String method,
        @JsonProperty("path") String path,
        @JsonProperty("headers") Map<String, String> headers,
        @JsonProperty("query_params") Map<String, String> queryParams,
        @JsonProperty("path_params") Map<String, String> pathParams,
        @JsonProperty("body") String body,
        @JsonProperty("remote_addr") String remoteAddr,
        @JsonProperty("user_agent") String userAgent,
        @JsonProperty("timestamp") Instant timestamp
) {

    /**
     * Compact constructor that ensures map fields are never null.
     */
    public HTTPRequestContext {
        headers = headers != null ? Collections.unmodifiableMap(headers) : Collections.emptyMap();
        queryParams = queryParams != null
                ? Collections.unmodifiableMap(queryParams)
                : Collections.emptyMap();
        pathParams = pathParams != null
                ? Collections.unmodifiableMap(pathParams)
                : Collections.emptyMap();
    }
}
