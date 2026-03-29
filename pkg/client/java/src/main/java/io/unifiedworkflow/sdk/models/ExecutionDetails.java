package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * Full details for a workflow execution, including per-step results and resource usage.
 *
 * @param runId             the unique execution identifier
 * @param status            current status snapshot; may be null if not yet hydrated
 * @param steps             ordered list of step results; never null, may be empty
 * @param primitivesUsed    names of primitives (external services) invoked during execution;
 *                          never null, may be empty
 * @param totalDurationMs   total wall-clock duration in milliseconds; 0 while still running
 * @param executionContext  arbitrary key-value context attached to this execution;
 *                         never null, may be empty
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record ExecutionDetails(
        @JsonProperty("run_id") String runId,
        @JsonProperty("status") ExecutionStatus status,
        @JsonProperty("steps") List<ExecutionStep> steps,
        @JsonProperty("primitives_used") List<String> primitivesUsed,
        @JsonProperty("total_duration_ms") long totalDurationMs,
        @JsonProperty("execution_context") Map<String, Object> executionContext
) {

    /**
     * Compact constructor that guarantees collections are never null.
     */
    public ExecutionDetails {
        steps = steps != null ? Collections.unmodifiableList(steps) : Collections.emptyList();
        primitivesUsed = primitivesUsed != null
                ? Collections.unmodifiableList(primitivesUsed)
                : Collections.emptyList();
        executionContext = executionContext != null
                ? Collections.unmodifiableMap(executionContext)
                : Collections.emptyMap();
    }
}
