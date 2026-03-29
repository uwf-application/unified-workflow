package io.unifiedworkflow.sdk.errors;

import java.util.Map;
import java.util.Objects;

/**
 * Thrown when the requested workflow execution run does not exist on the server.
 */
public class ExecutionNotFoundException extends WorkflowSDKException {

    private static final long serialVersionUID = 1L;

    private final String runId;

    /**
     * Constructs a new exception for the given run identifier.
     *
     * @param runId the run identifier that was not found; must not be null
     */
    public ExecutionNotFoundException(String runId) {
        super(ErrorCode.EXECUTION_NOT_FOUND,
                "Execution not found: " + runId, 404);
        this.runId = Objects.requireNonNull(runId, "runId must not be null");
    }

    /**
     * Constructs a new exception with a custom message and detail map.
     *
     * @param runId   the run identifier that was not found; must not be null
     * @param message human-readable description
     * @param details additional structured detail; may be null
     */
    public ExecutionNotFoundException(String runId, String message, Map<String, Object> details) {
        super(ErrorCode.EXECUTION_NOT_FOUND, message, 404, details);
        this.runId = Objects.requireNonNull(runId, "runId must not be null");
    }

    /**
     * Returns the run identifier that could not be found.
     *
     * @return the run ID, never null
     */
    public String getRunId() {
        return runId;
    }
}
