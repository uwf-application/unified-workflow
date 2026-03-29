package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Result for a single workflow execution within a batch response.
 *
 * @param workflowId      the identifier of the workflow that was executed
 * @param success         {@code true} if the execution was accepted or completed successfully
 * @param runId           the execution run identifier; may be null if the execution was rejected
 * @param error           human-readable error description when {@code success} is {@code false};
 *                        null otherwise
 * @param executionTimeMs wall-clock time in milliseconds that this execution took; 0 if not yet complete
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record BatchExecutionResult(
        @JsonProperty("workflow_id") String workflowId,
        @JsonProperty("success") boolean success,
        @JsonProperty("run_id") String runId,
        @JsonProperty("error") String error,
        @JsonProperty("execution_time_ms") long executionTimeMs
) {}
