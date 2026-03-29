package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;

/**
 * Response returned by the batch-execute endpoint.
 *
 * @param batchId    server-assigned unique identifier for this batch
 * @param total      total number of executions submitted
 * @param successful number of executions that were accepted or completed successfully
 * @param failed     number of executions that failed
 * @param pending    number of executions still queued or running
 * @param executions per-item result list; never null, may be empty
 * @param errors     list of batch-level error messages; never null, may be empty
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record BatchExecuteWorkflowsResponse(
        @JsonProperty("batch_id") String batchId,
        @JsonProperty("total") int total,
        @JsonProperty("successful") int successful,
        @JsonProperty("failed") int failed,
        @JsonProperty("pending") int pending,
        @JsonProperty("executions") List<BatchExecutionResult> executions,
        @JsonProperty("errors") List<String> errors
) {

    /**
     * Compact constructor that ensures collections are never null.
     */
    public BatchExecuteWorkflowsResponse {
        executions = executions != null
                ? Collections.unmodifiableList(executions)
                : Collections.emptyList();
        errors = errors != null ? Collections.unmodifiableList(errors) : Collections.emptyList();
    }
}
