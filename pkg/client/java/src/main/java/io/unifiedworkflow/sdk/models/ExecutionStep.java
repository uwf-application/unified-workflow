package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.Map;

/**
 * Represents the result of a single step within a workflow execution.
 *
 * @param stepId       unique identifier for this step within the workflow
 * @param stepName     human-readable name of the step
 * @param status       one of {@code "pending"}, {@code "running"}, {@code "completed"},
 *                     {@code "failed"}, {@code "skipped"}
 * @param durationMs   wall-clock execution duration in milliseconds; 0 while still running
 * @param errorMessage description of the failure when status is {@code "failed"}; null otherwise
 * @param output       key-value output data produced by this step; never null, may be empty
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record ExecutionStep(
        @JsonProperty("step_id") String stepId,
        @JsonProperty("step_name") String stepName,
        @JsonProperty("status") String status,
        @JsonProperty("duration_ms") long durationMs,
        @JsonProperty("error_message") String errorMessage,
        @JsonProperty("output") Map<String, Object> output
) {

    /**
     * Compact constructor that ensures output is never null.
     */
    public ExecutionStep {
        output = output != null ? Collections.unmodifiableMap(output) : Collections.emptyMap();
    }
}
