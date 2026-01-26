package executor

import (
	"context"
	"time"

	"unified-workflow/pkg/client"
)

// Client is the interface for executor service client
type Client interface {
	client.Client

	// ExecuteWorkflow executes a workflow
	ExecuteWorkflow(ctx context.Context, req *ExecuteWorkflowRequest) (*ExecuteWorkflowResponse, error)

	// GetExecutionStatus gets the status of a workflow execution
	GetExecutionStatus(ctx context.Context, req *GetExecutionStatusRequest) (*GetExecutionStatusResponse, error)

	// GetExecutionDetails gets detailed execution information
	GetExecutionDetails(ctx context.Context, req *GetExecutionDetailsRequest) (*GetExecutionDetailsResponse, error)

	// CancelExecution cancels a running workflow execution
	CancelExecution(ctx context.Context, req *CancelExecutionRequest) (*CancelExecutionResponse, error)

	// PauseExecution pauses a running workflow execution
	PauseExecution(ctx context.Context, req *PauseExecutionRequest) (*PauseExecutionResponse, error)

	// ResumeExecution resumes a paused workflow execution
	ResumeExecution(ctx context.Context, req *ResumeExecutionRequest) (*ResumeExecutionResponse, error)

	// RetryExecution retries a failed workflow execution
	RetryExecution(ctx context.Context, req *RetryExecutionRequest) (*RetryExecutionResponse, error)

	// GetExecutionData gets the data of a workflow execution
	GetExecutionData(ctx context.Context, req *GetExecutionDataRequest) (*GetExecutionDataResponse, error)

	// ListExecutions lists workflow executions with optional filters
	ListExecutions(ctx context.Context, req *ListExecutionsRequest) (*ListExecutionsResponse, error)

	// GetExecutionMetrics gets execution metrics for a workflow run
	GetExecutionMetrics(ctx context.Context, req *GetExecutionMetricsRequest) (*GetExecutionMetricsResponse, error)

	// GetStepExecution gets step execution details
	GetStepExecution(ctx context.Context, req *GetStepExecutionRequest) (*GetStepExecutionResponse, error)

	// GetChildStepExecution gets child step execution details
	GetChildStepExecution(ctx context.Context, req *GetChildStepExecutionRequest) (*GetChildStepExecutionResponse, error)
}

// ExecuteWorkflowRequest is the request for executing a workflow
type ExecuteWorkflowRequest struct {
	client.Request

	// WorkflowID is the workflow ID to execute
	WorkflowID string `json:"workflow_id"`

	// InputData is the input data for the workflow
	InputData map[string]interface{} `json:"input_data,omitempty"`

	// Async indicates if execution should be asynchronous
	Async bool `json:"async,omitempty"`

	// CallbackURL is the callback URL for async execution
	CallbackURL string `json:"callback_url,omitempty"`

	// TimeoutMs is the execution timeout in milliseconds
	TimeoutMs int64 `json:"timeout_ms,omitempty"`

	// Priority is the execution priority (1-10, higher is more important)
	Priority int `json:"priority,omitempty"`

	// Metadata contains additional execution metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ExecuteWorkflowResponse is the response for executing a workflow
type ExecuteWorkflowResponse struct {
	client.Response

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// Status is the initial execution status
	Status string `json:"status"`

	// StatusURL is the URL to check execution status
	StatusURL string `json:"status_url,omitempty"`

	// ResultURL is the URL to get execution result
	ResultURL string `json:"result_url,omitempty"`

	// EstimatedCompletionMs is the estimated completion time in milliseconds
	EstimatedCompletionMs int64 `json:"estimated_completion_ms,omitempty"`

	// QueuePosition is the position in the execution queue (if async)
	QueuePosition int `json:"queue_position,omitempty"`
}

// GetExecutionStatusRequest is the request for getting execution status
type GetExecutionStatusRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// IncludeDetails indicates if detailed status should be included
	IncludeDetails bool `json:"include_details,omitempty"`

	// WaitMs is the maximum time to wait for completion (for long polling)
	WaitMs int64 `json:"wait_ms,omitempty"`
}

// GetExecutionStatusResponse is the response for getting execution status
type GetExecutionStatusResponse struct {
	client.Response

	// Status is the execution status
	Status *ExecutionStatus `json:"status"`
}

// GetExecutionDetailsRequest is the request for getting execution details
type GetExecutionDetailsRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// IncludeSteps indicates if step details should be included
	IncludeSteps bool `json:"include_steps,omitempty"`

	// IncludeChildSteps indicates if child step details should be included
	IncludeChildSteps bool `json:"include_child_steps,omitempty"`

	// IncludePrimitives indicates if primitive details should be included
	IncludePrimitives bool `json:"include_primitives,omitempty"`
}

// GetExecutionDetailsResponse is the response for getting execution details
type GetExecutionDetailsResponse struct {
	client.Response

	// Details is the execution details
	Details *ExecutionDetails `json:"details"`
}

// CancelExecutionRequest is the request for canceling execution
type CancelExecutionRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// Reason is the cancellation reason
	Reason string `json:"reason,omitempty"`

	// Force indicates if cancellation should be forced
	Force bool `json:"force,omitempty"`
}

// CancelExecutionResponse is the response for canceling execution
type CancelExecutionResponse struct {
	client.Response

	// Cancelled indicates if the execution was cancelled
	Cancelled bool `json:"cancelled"`

	// Message is the cancellation message
	Message string `json:"message,omitempty"`
}

// PauseExecutionRequest is the request for pausing execution
type PauseExecutionRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// Reason is the pause reason
	Reason string `json:"reason,omitempty"`
}

// PauseExecutionResponse is the response for pausing execution
type PauseExecutionResponse struct {
	client.Response

	// Paused indicates if the execution was paused
	Paused bool `json:"paused"`

	// Message is the pause message
	Message string `json:"message,omitempty"`
}

// ResumeExecutionRequest is the request for resuming execution
type ResumeExecutionRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// Reason is the resume reason
	Reason string `json:"reason,omitempty"`
}

// ResumeExecutionResponse is the response for resuming execution
type ResumeExecutionResponse struct {
	client.Response

	// Resumed indicates if the execution was resumed
	Resumed bool `json:"resumed"`

	// Message is the resume message
	Message string `json:"message,omitempty"`
}

// RetryExecutionRequest is the request for retrying execution
type RetryExecutionRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// FromStepIndex is the step index to retry from (optional)
	FromStepIndex int `json:"from_step_index,omitempty"`

	// FromChildStepIndex is the child step index to retry from (optional)
	FromChildStepIndex int `json:"from_child_step_index,omitempty"`

	// ResetData indicates if execution data should be reset
	ResetData bool `json:"reset_data,omitempty"`
}

// RetryExecutionResponse is the response for retrying execution
type RetryExecutionResponse struct {
	client.Response

	// Retried indicates if the execution was retried
	Retried bool `json:"retried"`

	// NewRunID is the new run ID for the retry
	NewRunID string `json:"new_run_id,omitempty"`

	// Message is the retry message
	Message string `json:"message,omitempty"`
}

// GetExecutionDataRequest is the request for getting execution data
type GetExecutionDataRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// IncludeContext indicates if context data should be included
	IncludeContext bool `json:"include_context,omitempty"`

	// IncludeIntermediate indicates if intermediate data should be included
	IncludeIntermediate bool `json:"include_intermediate,omitempty"`
}

// GetExecutionDataResponse is the response for getting execution data
type GetExecutionDataResponse struct {
	client.Response

	// Data is the execution data
	Data map[string]interface{} `json:"data"`
}

// ListExecutionsRequest is the request for listing executions
type ListExecutionsRequest struct {
	client.Request

	// WorkflowID filters executions by workflow ID
	WorkflowID string `json:"workflow_id,omitempty"`

	// Status filters executions by status
	Status string `json:"status,omitempty"`

	// StartTime filters executions started after this time
	StartTime *time.Time `json:"start_time,omitempty"`

	// EndTime filters executions started before this time
	EndTime *time.Time `json:"end_time,omitempty"`

	// Limit limits the number of results
	Limit int `json:"limit,omitempty"`

	// Offset is the pagination offset
	Offset int `json:"offset,omitempty"`

	// SortBy is the field to sort by
	SortBy string `json:"sort_by,omitempty"`

	// SortOrder is the sort order (asc, desc)
	SortOrder string `json:"sort_order,omitempty"`
}

// ListExecutionsResponse is the response for listing executions
type ListExecutionsResponse struct {
	client.Response

	// Executions is the list of executions
	Executions []*ExecutionInfo `json:"executions"`

	// TotalCount is the total number of executions
	TotalCount int `json:"total_count"`

	// FilteredCount is the number of executions after filtering
	FilteredCount int `json:"filtered_count"`
}

// GetExecutionMetricsRequest is the request for getting execution metrics
type GetExecutionMetricsRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// IncludeStepMetrics indicates if step metrics should be included
	IncludeStepMetrics bool `json:"include_step_metrics,omitempty"`

	// IncludeChildStepMetrics indicates if child step metrics should be included
	IncludeChildStepMetrics bool `json:"include_child_step_metrics,omitempty"`
}

// GetExecutionMetricsResponse is the response for getting execution metrics
type GetExecutionMetricsResponse struct {
	client.Response

	// Metrics is the execution metrics
	Metrics *ExecutionMetrics `json:"metrics"`
}

// GetStepExecutionRequest is the request for getting step execution details
type GetStepExecutionRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// StepIndex is the step index
	StepIndex int `json:"step_index"`
}

// GetStepExecutionResponse is the response for getting step execution details
type GetStepExecutionResponse struct {
	client.Response

	// StepExecution is the step execution details
	StepExecution *StepExecution `json:"step_execution"`
}

// GetChildStepExecutionRequest is the request for getting child step execution details
type GetChildStepExecutionRequest struct {
	client.Request

	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// StepIndex is the step index
	StepIndex int `json:"step_index"`

	// ChildStepIndex is the child step index
	ChildStepIndex int `json:"child_step_index"`
}

// GetChildStepExecutionResponse is the response for getting child step execution details
type GetChildStepExecutionResponse struct {
	client.Response

	// ChildStepExecution is the child step execution details
	ChildStepExecution *ChildStepExecution `json:"child_step_execution"`
}

// ExecutionStatus represents workflow execution status
type ExecutionStatus struct {
	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// WorkflowID is the workflow ID
	WorkflowID string `json:"workflow_id"`

	// Status is the execution status
	Status string `json:"status"`

	// CurrentStep is the current step name
	CurrentStep string `json:"current_step,omitempty"`

	// CurrentStepIndex is the current step index
	CurrentStepIndex int `json:"current_step_index,omitempty"`

	// CurrentChildStepIndex is the current child step index
	CurrentChildStepIndex int `json:"current_child_step_index,omitempty"`

	// Progress is the execution progress (0.0 to 1.0)
	Progress float64 `json:"progress"`

	// StartTime is the execution start time
	StartTime *time.Time `json:"start_time,omitempty"`

	// EndTime is the execution end time
	EndTime *time.Time `json:"end_time,omitempty"`

	// ErrorMessage is the error message if execution failed
	ErrorMessage string `json:"error_message,omitempty"`

	// LastAttemptedStep is the last attempted step
	LastAttemptedStep string `json:"last_attempted_step,omitempty"`

	// IsTerminal indicates if the execution is in a terminal state
	IsTerminal bool `json:"is_terminal"`

	// Metadata contains additional status metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ExecutionDetails represents detailed execution information
type ExecutionDetails struct {
	ExecutionStatus

	// Steps is the list of step executions
	Steps []*StepExecution `json:"steps,omitempty"`

	// InputData is the execution input data
	InputData map[string]interface{} `json:"input_data,omitempty"`

	// OutputData is the execution output data
	OutputData map[string]interface{} `json:"output_data,omitempty"`

	// ContextData is the execution context data
	ContextData map[string]interface{} `json:"context_data,omitempty"`

	// PrimitivesUsed is the list of primitives used
	PrimitivesUsed []string `json:"primitives_used,omitempty"`

	// TotalDurationMs is the total execution duration in milliseconds
	TotalDurationMs int64 `json:"total_duration_ms,omitempty"`
}

// ExecutionInfo represents execution information for listing
type ExecutionInfo struct {
	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// WorkflowDefinitionID is the workflow definition ID
	WorkflowDefinitionID string `json:"workflow_definition_id"`

	// Status is the execution status
	Status string `json:"status"`

	// CurrentStepIndex is the current step index
	CurrentStepIndex int `json:"current_step_index"`

	// CurrentChildStepIndex is the current child step index
	CurrentChildStepIndex int `json:"current_child_step_index"`

	// StartTime is the execution start time
	StartTime *time.Time `json:"start_time,omitempty"`

	// EndTime is the execution end time
	EndTime *time.Time `json:"end_time,omitempty"`

	// ErrorMessage is the error message
	ErrorMessage string `json:"error_message,omitempty"`

	// LastAttemptedStep is the last attempted step
	LastAttemptedStep string `json:"last_attempted_step,omitempty"`

	// IsTerminal indicates if the execution is terminal
	IsTerminal bool `json:"is_terminal"`

	// IsRunning indicates if the execution is running
	IsRunning bool `json:"is_running"`

	// IsPending indicates if the execution is pending
	IsPending bool `json:"is_pending"`

	// CreatedAt is the creation timestamp
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
}

// ExecutionMetrics represents execution metrics
type ExecutionMetrics struct {
	// RunID is the execution run ID
	RunID string `json:"run_id"`

	// WorkflowID is the workflow ID
	WorkflowID string `json:"workflow_id"`

	// WorkflowMetrics contains workflow-level metrics
	WorkflowMetrics map[string]interface{} `json:"workflow_metrics,omitempty"`

	// StepMetrics contains step-level metrics
	StepMetrics []*StepMetrics `json:"step_metrics,omitempty"`

	// ChildStepMetrics contains child step-level metrics
	ChildStepMetrics []*ChildStepMetrics `json:"child_step_metrics,omitempty"`

	// TotalSteps is the total number of steps
	TotalSteps int `json:"total_steps"`

	// CompletedSteps is the number of completed steps
	CompletedSteps int `json:"completed_steps"`

	// FailedSteps is the number of failed steps
	FailedSteps int `json:"failed_steps"`

	// TotalChildSteps is the total number of child steps
	TotalChildSteps int `json:"total_child_steps"`

	// CompletedChildSteps is the number of completed child steps
	CompletedChildSteps int `json:"completed_child_steps"`

	// FailedChildSteps is the number of failed child steps
	FailedChildSteps int `json:"failed_child_steps"`

	// TotalDurationMillis is the total execution duration in milliseconds
	TotalDurationMillis int64 `json:"total_duration_millis"`

	// AverageStepDuration is the average step duration in milliseconds
	AverageStepDuration int64 `json:"average_step_duration"`

	// SuccessRate is the execution success rate (0.0 to 1.0)
	SuccessRate float64 `json:"success_rate"`
}

// StepExecution represents step execution details
type StepExecution struct {
	// StepIndex is the step index
	StepIndex int `json:"step_index"`

	// Name is the step name
	Name string `json:"name"`

	// Status is the step status
	Status string `json:"status"`

	// IsParallel indicates if the step is parallel
	IsParallel bool `json:"is_parallel"`

	// ChildStepCount is the number of child steps
	ChildStepCount int `json:"child_step_count"`

	// CompletedChildSteps is the number of completed child steps
	CompletedChildSteps int `json:"completed_child_steps"`

	// FailedChildSteps is the number of failed child steps
	FailedChildSteps int `json:"failed_child_steps"`

	// StartTime is the step start time
	StartTime *time.Time `json:"start_time,omitempty"`

	// EndTime is the step end time
	EndTime *time.Time `json:"end_time,omitempty"`

	// DurationMillis is the step duration in milliseconds
	DurationMillis int64 `json:"duration_millis,omitempty"`

	// ChildSteps is the list of child step executions
	ChildSteps []*ChildStepExecution `json:"child_steps,omitempty"`

	// ErrorMessage is the error message if step failed
	ErrorMessage string `json:"error_message,omitempty"`
}

// ChildStepExecution represents child step execution details
type ChildStepExecution struct {
	// StepIndex is the parent step index
	StepIndex int `json:"step_index"`

	// ChildStepIndex is the child step index
	ChildStepIndex int `json:"child_step_index"`

	// Name is the child step name
	Name string `json:"name"`

	// Status is the child step status
	Status string `json:"status"`

	// PrimitiveName is the primitive name used
	PrimitiveName string `json:"primitive_name,omitempty"`

	// StartTime is the child step start time
	StartTime *time.Time `json:"start_time,omitempty"`

	// EndTime is the child step end time
	EndTime *time.Time `json:"end_time,omitempty"`

	// DurationMillis is the child step duration in milliseconds
	DurationMillis int64 `json:"duration_millis,omitempty"`

	// ErrorMessage is the error message if child step failed
	ErrorMessage string `json:"error_message,omitempty"`

	// Result is the child step result
	Result interface{} `json:"result,omitempty"`

	// Parameters is the child step parameters
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// StepMetrics represents step execution metrics
type StepMetrics struct {
	// StepIndex is the step index
	StepIndex int `json:"step_index"`

	// Name is the step name
	Name string `json:"name"`

	// SuccessCount is the number of successful executions
	SuccessCount int `json:"success_count"`

	// FailureCount is the number of failed executions
	FailureCount int `json:"failure_count"`

	// TotalDurationMillis is the total duration in milliseconds
	TotalDurationMillis int64 `json:"total_duration_millis"`

	// AverageDurationMillis is the average duration in milliseconds
	AverageDurationMillis int64 `json:"average_duration_millis"`

	// MinDurationMillis is the minimum duration in milliseconds
	MinDurationMillis int64 `json:"min_duration_millis"`

	// MaxDurationMillis is the maximum duration in milliseconds
	MaxDurationMillis int64 `json:"max_duration_millis"`

	// SuccessRate is the success rate (0.0 to 1.0)
	SuccessRate float64 `json:"success_rate"`
}

// ChildStepMetrics represents child step execution metrics
type ChildStepMetrics struct {
	// StepIndex is the parent step index
	StepIndex int `json:"step_index"`

	// ChildStepIndex is the child step index
	ChildStepIndex int `json:"child_step_index"`

	// Name is the child step name
	Name string `json:"name"`

	// PrimitiveName is the primitive name
	PrimitiveName string `json:"primitive_name"`

	// SuccessCount is the number of successful executions
	SuccessCount int `json:"success_count"`

	// FailureCount is the number of failed executions
	FailureCount int `json:"failure_count"`

	// TotalDurationMillis is the total duration in milliseconds
	TotalDurationMillis int64 `json:"total_duration_millis"`

	// AverageDurationMillis is the average duration in milliseconds
	AverageDurationMillis int64 `json:"average_duration_millis"`

	// MinDurationMillis is the minimum duration in milliseconds
	MinDurationMillis int64 `json:"min_duration_millis"`

	// MaxDurationMillis is the maximum duration in milliseconds
	MaxDurationMillis int64 `json:"max_duration_millis"`

	// SuccessRate is the success rate (0.0 to 1.0)
	SuccessRate float64 `json:"success_rate"`
}

// Config is the executor client configuration
type Config struct {
	client.Config

	// DefaultTimeoutMs is the default execution timeout in milliseconds
	DefaultTimeoutMs int64 `json:"default_timeout_ms" yaml:"default_timeout_ms"`

	// DefaultPriority is the default execution priority
	DefaultPriority int `json:"default_priority" yaml:"default_priority"`

	// EnableAsync indicates if async execution is enabled
	EnableAsync bool `json:"enable_async" yaml:"enable_async"`
}

// DefaultConfig returns the default executor client configuration
func DefaultConfig() Config {
	return Config{
		Config:           client.DefaultConfig(),
		DefaultTimeoutMs: 30000, // 30 seconds
		DefaultPriority:  5,
		EnableAsync:      true,
	}
}
