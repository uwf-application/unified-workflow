package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Response returned by the SDK execute-workflow endpoint.
 *
 * @param runId                   unique identifier for this execution run
 * @param status                  initial status of the run (typically {@code "pending"} or
 *                                {@code "running"})
 * @param message                 human-readable status message from the server
 * @param statusUrl               URL to poll for the current execution status
 * @param resultUrl               URL to retrieve the final execution result once complete
 * @param pollAfterMs             suggested delay in milliseconds before the first status poll
 * @param estimatedCompletionMs   server's estimated execution duration in milliseconds; 0 if unknown
 * @param expiresAt               ISO-8601 timestamp after which the run record may be garbage-collected
 * @param validationResult        result of input validation performed before execution; may be null
 * @param sdkVersion              SDK server-side version string
 * @param requestId               correlation identifier echoed from the request
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record SDKExecuteWorkflowResponse(
        @JsonProperty("run_id") String runId,
        @JsonProperty("status") String status,
        @JsonProperty("message") String message,
        @JsonProperty("status_url") String statusUrl,
        @JsonProperty("result_url") String resultUrl,
        @JsonProperty("poll_after_ms") long pollAfterMs,
        @JsonProperty("estimated_completion_ms") long estimatedCompletionMs,
        @JsonProperty("expires_at") String expiresAt,
        @JsonProperty("validation_result") ValidationResult validationResult,
        @JsonProperty("sdk_version") String sdkVersion,
        @JsonProperty("request_id") String requestId
) {}
