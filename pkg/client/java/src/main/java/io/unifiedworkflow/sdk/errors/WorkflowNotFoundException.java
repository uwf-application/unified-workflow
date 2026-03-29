package io.unifiedworkflow.sdk.errors;

import java.util.Map;
import java.util.Objects;

/**
 * Thrown when the requested workflow definition does not exist on the server.
 */
public class WorkflowNotFoundException extends WorkflowSDKException {

    private static final long serialVersionUID = 1L;

    private final String workflowId;

    /**
     * Constructs a new exception for the given workflow identifier.
     *
     * @param workflowId the identifier of the workflow that was not found; must not be null
     */
    public WorkflowNotFoundException(String workflowId) {
        super(ErrorCode.WORKFLOW_NOT_FOUND,
                "Workflow not found: " + workflowId, 404);
        this.workflowId = Objects.requireNonNull(workflowId, "workflowId must not be null");
    }

    /**
     * Constructs a new exception with a custom message and detail map.
     *
     * @param workflowId the identifier of the workflow that was not found; must not be null
     * @param message    human-readable description
     * @param details    additional structured detail; may be null
     */
    public WorkflowNotFoundException(String workflowId, String message, Map<String, Object> details) {
        super(ErrorCode.WORKFLOW_NOT_FOUND, message, 404, details);
        this.workflowId = Objects.requireNonNull(workflowId, "workflowId must not be null");
    }

    /**
     * Returns the workflow identifier that could not be found.
     *
     * @return the workflow ID, never null
     */
    public String getWorkflowId() {
        return workflowId;
    }
}
