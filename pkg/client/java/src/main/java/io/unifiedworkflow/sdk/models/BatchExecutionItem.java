package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Objects;

/**
 * A single workflow execution entry within a {@link BatchExecuteWorkflowsRequest}.
 *
 * @param workflowId the identifier of the workflow to execute; must not be null
 * @param request    the execution request payload; may be null to use defaults
 * @param priority   execution priority for this item in the range [1, 10];
 *                   higher values are processed first
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record BatchExecutionItem(
        @JsonProperty("workflow_id") String workflowId,
        @JsonProperty("request") SDKExecuteWorkflowRequest request,
        @JsonProperty("priority") int priority
) {

    /**
     * Compact constructor that validates the required {@code workflowId}.
     */
    public BatchExecutionItem {
        Objects.requireNonNull(workflowId, "workflowId must not be null");
    }
}
