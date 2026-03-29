package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;

/**
 * Lightweight status snapshot for a running or completed workflow execution.
 *
 * <p>Use {@link #isTerminal()} to determine whether polling should stop.</p>
 *
 * @param runId         the unique execution identifier
 * @param status        one of {@code "pending"}, {@code "running"}, {@code "completed"},
 *                      {@code "failed"}, {@code "cancelled"}
 * @param progress      completion percentage in the range [0, 100]
 * @param currentStep   name of the step currently executing; may be null when not running
 * @param isTerminal    {@code true} when the execution has reached a terminal state and
 *                      will not change further
 * @param errorMessage  human-readable error description when {@code status} is {@code "failed"};
 *                      null otherwise
 * @param startTime     wall-clock time at which the execution began; may be null for pending runs
 * @param endTime       wall-clock time at which the execution finished; null while still running
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record ExecutionStatus(
        @JsonProperty("run_id") String runId,
        @JsonProperty("status") String status,
        @JsonProperty("progress") int progress,
        @JsonProperty("current_step") String currentStep,
        @JsonProperty("is_terminal") boolean isTerminal,
        @JsonProperty("error_message") String errorMessage,
        @JsonProperty("start_time") Instant startTime,
        @JsonProperty("end_time") Instant endTime
) {

    /** Status constant for an execution that has not yet started. */
    public static final String STATUS_PENDING = "pending";

    /** Status constant for a currently executing run. */
    public static final String STATUS_RUNNING = "running";

    /** Status constant for a successfully completed execution. */
    public static final String STATUS_COMPLETED = "completed";

    /** Status constant for a failed execution. */
    public static final String STATUS_FAILED = "failed";

    /** Status constant for a cancelled execution. */
    public static final String STATUS_CANCELLED = "cancelled";

    /**
     * Returns {@code true} when this execution completed successfully.
     *
     * @return {@code true} if status is {@code "completed"}
     */
    public boolean isCompleted() {
        return STATUS_COMPLETED.equals(status);
    }

    /**
     * Returns {@code true} when this execution failed.
     *
     * @return {@code true} if status is {@code "failed"}
     */
    public boolean isFailed() {
        return STATUS_FAILED.equals(status);
    }
}
