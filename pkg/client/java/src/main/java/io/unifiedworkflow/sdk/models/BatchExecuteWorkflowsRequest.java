package io.unifiedworkflow.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Collections;
import java.util.List;
import java.util.Objects;

/**
 * Request to execute multiple workflows in a single API call.
 *
 * <p>Construct instances using the nested {@link Builder}:</p>
 * <pre>{@code
 * BatchExecuteWorkflowsRequest batch = BatchExecuteWorkflowsRequest.builder()
 *     .addExecution(new BatchExecutionItem("workflow-a", requestA, 5))
 *     .addExecution(new BatchExecutionItem("workflow-b", requestB, 3))
 *     .parallel(true)
 *     .maxConcurrent(4)
 *     .build();
 * }</pre>
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public final class BatchExecuteWorkflowsRequest {

    @JsonProperty("executions")
    private final List<BatchExecutionItem> executions;

    @JsonProperty("parallel")
    private final boolean parallel;

    @JsonProperty("max_concurrent")
    private final int maxConcurrent;

    @JsonProperty("stop_on_first_failure")
    private final boolean stopOnFirstFailure;

    private BatchExecuteWorkflowsRequest(Builder builder) {
        this.executions = builder.executions != null
                ? Collections.unmodifiableList(builder.executions)
                : Collections.emptyList();
        this.parallel = builder.parallel;
        this.maxConcurrent = builder.maxConcurrent;
        this.stopOnFirstFailure = builder.stopOnFirstFailure;
    }

    /**
     * Returns a new {@link Builder}.
     *
     * @return a fresh builder instance
     */
    public static Builder builder() {
        return new Builder();
    }

    /** @return list of execution items; never null */
    public List<BatchExecutionItem> getExecutions() { return executions; }

    /** @return {@code true} if executions should run in parallel */
    public boolean isParallel() { return parallel; }

    /** @return maximum number of concurrently running executions when parallel is enabled */
    public int getMaxConcurrent() { return maxConcurrent; }

    /** @return {@code true} if the batch should abort after the first execution failure */
    public boolean isStopOnFirstFailure() { return stopOnFirstFailure; }

    // -------------------------------------------------------------------------
    // Builder
    // -------------------------------------------------------------------------

    /**
     * Fluent builder for {@link BatchExecuteWorkflowsRequest}.
     */
    public static final class Builder {

        private List<BatchExecutionItem> executions;
        private boolean parallel = false;
        private int maxConcurrent = 10;
        private boolean stopOnFirstFailure = false;

        private Builder() {}

        /**
         * Sets the complete list of execution items.
         *
         * @param executions non-null list of items
         * @return this builder
         */
        public Builder executions(List<BatchExecutionItem> executions) {
            this.executions = Objects.requireNonNull(executions, "executions must not be null");
            return this;
        }

        /**
         * Adds a single execution item to the batch.
         *
         * @param item non-null item to add
         * @return this builder
         */
        public Builder addExecution(BatchExecutionItem item) {
            Objects.requireNonNull(item, "item must not be null");
            if (this.executions == null) {
                this.executions = new java.util.ArrayList<>();
            }
            this.executions.add(item);
            return this;
        }

        /**
         * Controls whether executions run in parallel.
         *
         * @param parallel {@code true} for parallel execution
         * @return this builder
         */
        public Builder parallel(boolean parallel) {
            this.parallel = parallel;
            return this;
        }

        /**
         * Sets the concurrency cap when running in parallel mode.
         *
         * @param maxConcurrent maximum simultaneous executions; defaults to 10
         * @return this builder
         */
        public Builder maxConcurrent(int maxConcurrent) {
            this.maxConcurrent = maxConcurrent;
            return this;
        }

        /**
         * Controls whether the batch stops immediately when one execution fails.
         *
         * @param stopOnFirstFailure {@code true} to stop on first failure
         * @return this builder
         */
        public Builder stopOnFirstFailure(boolean stopOnFirstFailure) {
            this.stopOnFirstFailure = stopOnFirstFailure;
            return this;
        }

        /**
         * Builds and returns the request.
         *
         * @return a new immutable {@link BatchExecuteWorkflowsRequest}
         */
        public BatchExecuteWorkflowsRequest build() {
            return new BatchExecuteWorkflowsRequest(this);
        }
    }
}
