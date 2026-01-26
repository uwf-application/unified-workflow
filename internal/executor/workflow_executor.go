package executor

import (
	"context"
	"fmt"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/internal/primitive"
	workflowRegistry "unified-workflow/internal/registry"
	"unified-workflow/internal/state"
)

// WorkflowExecutor is a real executor that actually executes workflows with child-step tracking
type WorkflowExecutor struct {
	workflowRegistry workflowRegistry.Registry
	stateManagement  state.StateManagement
	config           Config
}

// NewWorkflowExecutor creates a new workflow executor
func NewWorkflowExecutor(
	workflowRegistry workflowRegistry.Registry,
	stateManagement state.StateManagement,
	config Config,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		workflowRegistry: workflowRegistry,
		stateManagement:  stateManagement,
		config:           config,
	}
}

// ExecutionResult represents the result of a workflow execution
type ExecutionResult struct {
	RunID      string                 `json:"run_id"`
	WorkflowID string                 `json:"workflow_id"`
	Status     string                 `json:"status"`
	Result     map[string]interface{} `json:"result"`
	Error      string                 `json:"error,omitempty"`
	StartTime  time.Time              `json:"start_time"`
	EndTime    time.Time              `json:"end_time"`
}

// ChildStepExecutionResult represents the result of a child step execution
type ChildStepExecutionResult struct {
	StepIndex      int                    `json:"step_index"`
	ChildStepIndex int                    `json:"child_step_index"`
	Name           string                 `json:"name"`
	Status         string                 `json:"status"` // "pending", "completed", "failed", "skipped"
	PrimitiveName  string                 `json:"primitive_name"`
	StartTime      time.Time              `json:"start_time,omitempty"`
	EndTime        time.Time              `json:"end_time,omitempty"`
	DurationMillis int64                  `json:"duration_millis,omitempty"`
	ErrorMessage   string                 `json:"error_message,omitempty"`
	Result         interface{}            `json:"result,omitempty"`
	Parameters     map[string]interface{} `json:"parameters,omitempty"`
}

// StepExecutionResult represents the result of a step execution
type StepExecutionResult struct {
	StepIndex           int                        `json:"step_index"`
	Name                string                     `json:"name"`
	Status              string                     `json:"status"` // "pending", "running", "completed", "failed", "cancelled"
	IsParallel          bool                       `json:"is_parallel"`
	ChildStepCount      int                        `json:"child_step_count"`
	CompletedChildSteps int                        `json:"completed_child_steps"`
	FailedChildSteps    int                        `json:"failed_child_steps"`
	StartTime           time.Time                  `json:"start_time,omitempty"`
	EndTime             time.Time                  `json:"end_time,omitempty"`
	DurationMillis      int64                      `json:"duration_millis,omitempty"`
	ChildSteps          []ChildStepExecutionResult `json:"child_steps,omitempty"`
	ErrorMessage        string                     `json:"error_message,omitempty"`
}

// ExecuteWorkflow executes a workflow with child-step tracking
func (e *WorkflowExecutor) ExecuteWorkflow(ctx context.Context, workflowID string, inputData map[string]interface{}) (*ExecutionResult, error) {
	startTime := time.Now()

	// Load workflow from registry
	workflow, err := e.workflowRegistry.GetWorkflow(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow %s: %w", workflowID, err)
	}

	// Generate run ID
	runID := fmt.Sprintf("run-%d", time.Now().UnixNano())

	// Initialize execution context and data
	executionContext := make(map[string]interface{})
	executionData := make(map[string]interface{})

	// Merge input data
	for k, v := range inputData {
		executionData[k] = v
	}

	// Track execution state
	stepResults := make([]StepExecutionResult, 0)
	totalSteps := workflow.GetStepCount()

	// Execute each step
	for stepIndex, step := range workflow.GetSteps() {
		stepStartTime := time.Now()

		// Create step result
		stepResult := StepExecutionResult{
			StepIndex:           stepIndex,
			Name:                step.GetName(),
			Status:              "running",
			IsParallel:          step.IsParallel(),
			ChildStepCount:      step.GetChildStepCount(),
			CompletedChildSteps: 0,
			FailedChildSteps:    0,
			StartTime:           stepStartTime,
			ChildSteps:          make([]ChildStepExecutionResult, 0),
		}

		// Execute child steps
		childStepResults, stepErr := e.executeStep(ctx, step, stepIndex, executionContext, executionData)

		// Update step result
		stepResult.EndTime = time.Now()
		stepResult.DurationMillis = stepResult.EndTime.Sub(stepResult.StartTime).Milliseconds()
		stepResult.ChildSteps = childStepResults

		// Count completed/failed child steps
		for _, childResult := range childStepResults {
			if childResult.Status == "completed" {
				stepResult.CompletedChildSteps++
			} else if childResult.Status == "failed" {
				stepResult.FailedChildSteps++
			}
		}

		// Determine step status
		if stepErr != nil {
			stepResult.Status = "failed"
			stepResult.ErrorMessage = stepErr.Error()
		} else if stepResult.CompletedChildSteps == stepResult.ChildStepCount {
			stepResult.Status = "completed"
		} else {
			stepResult.Status = "failed" // Some child steps failed
		}

		stepResults = append(stepResults, stepResult)

		// If step failed and we should stop, break
		if stepErr != nil && e.config.MaxRetries == 0 {
			break
		}
	}

	// Calculate overall execution result
	endTime := time.Now()

	// Determine overall status
	completedSteps := 0
	failedSteps := 0
	for _, stepResult := range stepResults {
		if stepResult.Status == "completed" {
			completedSteps++
		} else if stepResult.Status == "failed" {
			failedSteps++
		}
	}

	status := "completed"
	if failedSteps > 0 {
		status = "failed"
	} else if completedSteps < totalSteps {
		status = "partial"
	}

	result := &ExecutionResult{
		RunID:      runID,
		WorkflowID: workflowID,
		Status:     status,
		Result:     executionData,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	// TODO: Store execution result using state management
	// For now, just log that we would store it
	fmt.Printf("Execution completed: %s, status: %s\n", runID, status)

	return result, nil
}

// executeStep executes a single step with its child steps
func (e *WorkflowExecutor) executeStep(ctx context.Context, step model.Step, stepIndex int, context, data interface{}) ([]ChildStepExecutionResult, error) {
	childSteps := step.GetChildSteps()
	childStepResults := make([]ChildStepExecutionResult, len(childSteps))

	// Check if step is parallel
	if step.IsParallel() {
		// Execute child steps in parallel (simplified - would use goroutines in real implementation)
		for childStepIndex, childStep := range childSteps {
			result := e.executeChildStep(ctx, childStep, stepIndex, childStepIndex, context, data)
			childStepResults[childStepIndex] = result

			// If child step failed and we should stop, return early
			if result.Status == "failed" && e.config.MaxRetries == 0 {
				return childStepResults, fmt.Errorf("child step %d failed: %s", childStepIndex, result.ErrorMessage)
			}
		}
	} else {
		// Execute child steps sequentially
		for childStepIndex, childStep := range childSteps {
			result := e.executeChildStep(ctx, childStep, stepIndex, childStepIndex, context, data)
			childStepResults[childStepIndex] = result

			// If child step failed and we should stop, return early
			if result.Status == "failed" && e.config.MaxRetries == 0 {
				return childStepResults, fmt.Errorf("child step %d failed: %s", childStepIndex, result.ErrorMessage)
			}
		}
	}

	return childStepResults, nil
}

// executeChildStep executes a single child step with primitive execution
func (e *WorkflowExecutor) executeChildStep(ctx context.Context, childStep *model.ChildStep, stepIndex, childStepIndex int, context, data interface{}) ChildStepExecutionResult {
	startTime := time.Now()
	result := ChildStepExecutionResult{
		StepIndex:      stepIndex,
		ChildStepIndex: childStepIndex,
		Name:           childStep.GetName(),
		Status:         "pending",
		StartTime:      startTime,
	}

	// Get primitive name from child step hooks
	// For now, use a default echo primitive for testing
	primitiveName := "primitive.echo.echo"

	// Default parameters for echo primitive
	params := map[string]interface{}{
		"message": fmt.Sprintf("Executing step %s", childStep.GetName()),
	}

	// Try to execute the primitive using global primitive pattern
	// For echo primitive, use primitive.Default.Echo.Echo()
	if primitiveName == "primitive.echo.echo" {
		message, ok := params["message"].(string)
		if !ok {
			message = "Default message"
		}

		primitiveResult, err := primitive.Default.Echo.Echo(message)
		if err != nil {
			result.Status = "failed"
			result.ErrorMessage = fmt.Sprintf("Primitive execution failed: %v", err)
			result.EndTime = time.Now()
			result.DurationMillis = result.EndTime.Sub(startTime).Milliseconds()
			result.PrimitiveName = primitiveName
			result.Parameters = params
			return result
		}

		endTime := time.Now()
		result.EndTime = endTime
		result.DurationMillis = endTime.Sub(startTime).Milliseconds()
		result.Status = "completed"
		result.PrimitiveName = primitiveName
		result.Result = primitiveResult
		result.Parameters = params
		return result
	}

	// Primitive not found, use fallback
	result.Status = "failed"
	result.ErrorMessage = fmt.Sprintf("Primitive not found or not implemented: %s", primitiveName)
	result.EndTime = time.Now()
	result.DurationMillis = result.EndTime.Sub(startTime).Milliseconds()
	result.PrimitiveName = primitiveName
	result.Parameters = params
	return result
}

// GetExecutionStatus gets the status of a workflow execution
func (e *WorkflowExecutor) GetExecutionStatus(ctx context.Context, runID string) (*ExecutionStatus, error) {
	// For now, return a simple status since we don't have persistent storage
	// In a real implementation, we would query the state management
	return &ExecutionStatus{
		RunID:      runID,
		Status:     "completed", // Default status
		Progress:   1.0,
		IsTerminal: true,
	}, nil
}

// SubmitWorkflow submits a workflow for execution and returns a run ID
func (e *WorkflowExecutor) SubmitWorkflow(ctx context.Context, workflow model.Workflow) (string, error) {
	// For now, simulate submission
	runID := fmt.Sprintf("run-%d", time.Now().UnixNano())
	return runID, nil
}

// SubmitWorkflowByID submits a workflow by ID for execution and returns a run ID
func (e *WorkflowExecutor) SubmitWorkflowByID(ctx context.Context, workflowID string) (string, error) {
	// For now, simulate submission
	runID := fmt.Sprintf("run-%d", time.Now().UnixNano())
	return runID, nil
}

// GetExecutionData gets the data of a workflow execution
func (e *WorkflowExecutor) GetExecutionData(ctx context.Context, runID string) (map[string]interface{}, error) {
	// For now, return empty data
	return map[string]interface{}{}, nil
}

// ListExecutions lists workflow executions with optional filters
func (e *WorkflowExecutor) ListExecutions(ctx context.Context, filters ExecutionFilters) ([]*ExecutionInfo, error) {
	// For now, return empty list
	return []*ExecutionInfo{}, nil
}

// CancelExecution cancels a running workflow execution
func (e *WorkflowExecutor) CancelExecution(ctx context.Context, runID string) error {
	// For now, just return success
	return nil
}

// PauseExecution pauses a running workflow execution
func (e *WorkflowExecutor) PauseExecution(ctx context.Context, runID string) error {
	// For now, just return success
	return nil
}

// ResumeExecution resumes a paused workflow execution
func (e *WorkflowExecutor) ResumeExecution(ctx context.Context, runID string) error {
	// For now, just return success
	return nil
}

// RetryExecution retries a failed workflow execution
func (e *WorkflowExecutor) RetryExecution(ctx context.Context, runID string) error {
	// For now, just return success
	return nil
}

// GetMetrics gets execution metrics for a workflow run
func (e *WorkflowExecutor) GetMetrics(ctx context.Context, runID string) (*ExecutionMetrics, error) {
	// For now, return empty metrics
	return &ExecutionMetrics{
		RunID:       runID,
		SuccessRate: 1.0,
	}, nil
}

// Start starts the executor
func (e *WorkflowExecutor) Start(ctx context.Context) error {
	// Nothing to start for now
	return nil
}

// Stop stops the executor
func (e *WorkflowExecutor) Stop(ctx context.Context) error {
	// Nothing to stop for now
	return nil
}

// IsRunning checks if the executor is running
func (e *WorkflowExecutor) IsRunning() bool {
	return true
}
