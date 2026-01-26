package executor

import (
	"context"
	"fmt"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/internal/queue"
	"unified-workflow/internal/registry"
	"unified-workflow/internal/state"
)

// SimpleExecutor is a simplified executor for the workflow API
type SimpleExecutor struct {
	registry        registry.Registry
	queue           queue.Queue
	stateManagement state.StateManagement
}

// NewSimpleExecutor creates a new simple executor
func NewSimpleExecutor(
	registry registry.Registry,
	queue queue.Queue,
	stateManagement state.StateManagement,
) *SimpleExecutor {
	return &SimpleExecutor{
		registry:        registry,
		queue:           queue,
		stateManagement: stateManagement,
	}
}

// SubmitWorkflow submits a workflow for execution and returns a run ID
func (e *SimpleExecutor) SubmitWorkflow(ctx context.Context, workflow model.Workflow) (string, error) {
	// Create a simple run ID
	runID := fmt.Sprintf("run-%d", time.Now().UnixNano())

	// Create execution request
	executionReq := queue.ExecutionRequest{
		RunID:       runID,
		WorkflowID:  workflow.GetID(),
		InputData:   make(map[string]interface{}),
		RequestedAt: time.Now(),
	}

	// Marshal request
	reqData, err := queue.MarshalExecutionRequest(executionReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal execution request: %w", err)
	}

	// Enqueue for execution
	err = e.queue.Enqueue(ctx, runID, reqData)
	if err != nil {
		return "", fmt.Errorf("failed to enqueue workflow: %w", err)
	}

	return runID, nil
}

// SubmitWorkflowByID submits a workflow by ID for execution and returns a run ID
func (e *SimpleExecutor) SubmitWorkflowByID(ctx context.Context, workflowID string) (string, error) {
	// Get workflow from registry
	workflow, err := e.registry.GetWorkflow(ctx, workflowID)
	if err != nil {
		return "", fmt.Errorf("failed to get workflow: %w", err)
	}

	return e.SubmitWorkflow(ctx, workflow)
}

// GetExecutionStatus gets the status of a workflow execution
func (e *SimpleExecutor) GetExecutionStatus(ctx context.Context, runID string) (*ExecutionStatus, error) {
	// For now, return a simple status
	return &ExecutionStatus{
		RunID:     runID,
		Status:    "pending",
		Progress:  0.0,
		StartTime: nil,
		EndTime:   nil,
	}, nil
}

// GetExecutionData gets the data of a workflow execution
func (e *SimpleExecutor) GetExecutionData(ctx context.Context, runID string) (map[string]interface{}, error) {
	// For now, return empty data
	return make(map[string]interface{}), nil
}

// ListExecutions lists workflow executions with optional filters
func (e *SimpleExecutor) ListExecutions(ctx context.Context, filters ExecutionFilters) ([]*ExecutionInfo, error) {
	// For now, return empty list
	return []*ExecutionInfo{}, nil
}

// CancelExecution cancels a running workflow execution
func (e *SimpleExecutor) CancelExecution(ctx context.Context, runID string) error {
	return fmt.Errorf("not implemented")
}

// PauseExecution pauses a running workflow execution
func (e *SimpleExecutor) PauseExecution(ctx context.Context, runID string) error {
	return fmt.Errorf("not implemented")
}

// ResumeExecution resumes a paused workflow execution
func (e *SimpleExecutor) ResumeExecution(ctx context.Context, runID string) error {
	return fmt.Errorf("not implemented")
}

// RetryExecution retries a failed workflow execution
func (e *SimpleExecutor) RetryExecution(ctx context.Context, runID string) error {
	return fmt.Errorf("not implemented")
}

// GetMetrics gets execution metrics for a workflow run
func (e *SimpleExecutor) GetMetrics(ctx context.Context, runID string) (*ExecutionMetrics, error) {
	return &ExecutionMetrics{
		RunID:      runID,
		TotalSteps: 0,
	}, nil
}

// Start starts the executor
func (e *SimpleExecutor) Start(ctx context.Context) error {
	return nil
}

// Stop stops the executor
func (e *SimpleExecutor) Stop(ctx context.Context) error {
	return nil
}

// IsRunning checks if the executor is running
func (e *SimpleExecutor) IsRunning() bool {
	return true
}
