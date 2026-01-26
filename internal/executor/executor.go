package executor

import (
	"context"
	"time"

	"unified-workflow/internal/common/model"
)

// Executor is the interface for workflow execution implementations
type Executor interface {
	// SubmitWorkflow submits a workflow for execution and returns a run ID
	SubmitWorkflow(ctx context.Context, workflow model.Workflow) (string, error)

	// SubmitWorkflowByID submits a workflow by ID for execution and returns a run ID
	SubmitWorkflowByID(ctx context.Context, workflowID string) (string, error)

	// GetExecutionStatus gets the status of a workflow execution
	GetExecutionStatus(ctx context.Context, runID string) (*ExecutionStatus, error)

	// GetExecutionData gets the data of a workflow execution
	GetExecutionData(ctx context.Context, runID string) (map[string]interface{}, error)

	// ListExecutions lists workflow executions with optional filters
	ListExecutions(ctx context.Context, filters ExecutionFilters) ([]*ExecutionInfo, error)

	// CancelExecution cancels a running workflow execution
	CancelExecution(ctx context.Context, runID string) error

	// PauseExecution pauses a running workflow execution
	PauseExecution(ctx context.Context, runID string) error

	// ResumeExecution resumes a paused workflow execution
	ResumeExecution(ctx context.Context, runID string) error

	// RetryExecution retries a failed workflow execution
	RetryExecution(ctx context.Context, runID string) error

	// GetMetrics gets execution metrics for a workflow run
	GetMetrics(ctx context.Context, runID string) (*ExecutionMetrics, error)

	// Start starts the executor
	Start(ctx context.Context) error

	// Stop stops the executor gracefully
	Stop(ctx context.Context) error

	// IsRunning checks if the executor is running
	IsRunning() bool
}

// ExecutionStatus represents the status of a workflow execution
type ExecutionStatus struct {
	RunID                 string                 `json:"run_id"`
	WorkflowID            string                 `json:"workflow_id"`
	Status                string                 `json:"status"`
	CurrentStep           string                 `json:"current_step,omitempty"`
	CurrentStepIndex      int                    `json:"current_step_index"`
	CurrentChildStepIndex int                    `json:"current_child_step_index"`
	Progress              float64                `json:"progress"` // 0.0 to 1.0
	StartTime             *time.Time             `json:"start_time,omitempty"`
	EndTime               *time.Time             `json:"end_time,omitempty"`
	ErrorMessage          string                 `json:"error_message,omitempty"`
	LastAttemptedStep     string                 `json:"last_attempted_step,omitempty"`
	IsTerminal            bool                   `json:"is_terminal"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
}

// ExecutionFilters represents filters for listing executions
type ExecutionFilters struct {
	WorkflowID    string     `json:"workflow_id,omitempty"`
	Status        string     `json:"status,omitempty"`
	IsTerminal    *bool      `json:"is_terminal,omitempty"`
	IsRunning     *bool      `json:"is_running,omitempty"`
	IsPending     *bool      `json:"is_pending,omitempty"`
	StartTimeFrom *time.Time `json:"start_time_from,omitempty"`
	StartTimeTo   *time.Time `json:"start_time_to,omitempty"`
	Limit         int        `json:"limit,omitempty"`
	Offset        int        `json:"offset,omitempty"`
	SortBy        string     `json:"sort_by,omitempty"`
	SortOrder     string     `json:"sort_order,omitempty"` // "asc" or "desc"
}

// ExecutionInfo represents workflow execution information for listing
type ExecutionInfo struct {
	RunID                 string     `json:"run_id"`
	WorkflowDefinitionID  string     `json:"workflow_definition_id"`
	WorkflowName          string     `json:"workflow_name,omitempty"`
	Status                string     `json:"status"`
	CurrentStepIndex      int        `json:"current_step_index"`
	CurrentChildStepIndex int        `json:"current_child_step_index"`
	StartTime             *time.Time `json:"start_time,omitempty"`
	EndTime               *time.Time `json:"end_time,omitempty"`
	ErrorMessage          string     `json:"error_message,omitempty"`
	LastAttemptedStep     string     `json:"last_attempted_step,omitempty"`
	IsTerminal            bool       `json:"is_terminal"`
	IsRunning             bool       `json:"is_running"`
	IsPending             bool       `json:"is_pending"`
	DurationMillis        *int64     `json:"duration_millis,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// ExecutionMetrics represents execution metrics
type ExecutionMetrics struct {
	RunID               string                 `json:"run_id"`
	WorkflowID          string                 `json:"workflow_id"`
	WorkflowMetrics     map[string]interface{} `json:"workflow_metrics"`
	StepMetrics         map[string]interface{} `json:"step_metrics"`
	ChildStepMetrics    map[string]interface{} `json:"child_step_metrics"`
	TotalSteps          int                    `json:"total_steps"`
	CompletedSteps      int                    `json:"completed_steps"`
	FailedSteps         int                    `json:"failed_steps"`
	TotalChildSteps     int                    `json:"total_child_steps"`
	CompletedChildSteps int                    `json:"completed_child_steps"`
	FailedChildSteps    int                    `json:"failed_child_steps"`
	TotalDurationMillis int64                  `json:"total_duration_millis,omitempty"`
	AverageStepDuration int64                  `json:"average_step_duration_millis,omitempty"`
	SuccessRate         float64                `json:"success_rate"` // 0.0 to 1.0
}

// Config represents executor configuration
type Config struct {
	WorkerCount            int           `json:"worker_count"`
	QueuePollInterval      time.Duration `json:"queue_poll_interval"`
	MaxRetries             int           `json:"max_retries"`
	RetryDelay             time.Duration `json:"retry_delay"`
	ExecutionTimeout       time.Duration `json:"execution_timeout"`
	StepTimeout            time.Duration `json:"step_timeout"`
	EnableMetrics          bool          `json:"enable_metrics"`
	EnableTracing          bool          `json:"enable_tracing"`
	MaxConcurrentWorkflows int           `json:"max_concurrent_workflows"`
}
