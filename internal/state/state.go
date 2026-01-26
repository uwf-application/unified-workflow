package state

import (
	"context"
	"time"

	"unified-workflow/internal/primitive/model"
)

// StateManagement is the interface for workflow state management implementations
// Provides abstraction for different state storage backends (in-memory, database, distributed cache)
type StateManagement interface {
	// SaveContext saves the workflow context to the store
	SaveContext(ctx context.Context, workflowContext model.WorkflowContext) error

	// GetContext retrieves the workflow context for the given run ID
	GetContext(ctx context.Context, runID string) (model.WorkflowContext, error)

	// SaveData saves the workflow data to the store
	SaveData(ctx context.Context, runID string, workflowData model.WorkflowData) error

	// GetData retrieves the workflow data for the given run ID
	GetData(ctx context.Context, runID string) (model.WorkflowData, error)

	// RemoveState removes all state (both context and data) for the given run ID
	RemoveState(ctx context.Context, runID string) error

	// ContainsContext checks if a workflow context exists for the given run ID
	ContainsContext(ctx context.Context, runID string) (bool, error)

	// ContainsData checks if workflow data exists for the given run ID
	ContainsData(ctx context.Context, runID string) (bool, error)

	// AcquireLock acquires a lock for a workflow run to prevent concurrent modifications
	AcquireLock(ctx context.Context, runID string, timeout time.Duration) (bool, error)

	// ReleaseLock releases a lock for a workflow run
	ReleaseLock(ctx context.Context, runID string) error

	// SetTTL sets time-to-live for workflow state
	SetTTL(ctx context.Context, runID string, ttl time.Duration) error

	// GetAllContexts gets all workflow contexts stored in the state management
	GetAllContexts(ctx context.Context) ([]model.WorkflowContext, error)

	// Close closes the state management connection
	Close() error
}

// ExecutionInfo represents workflow execution information for listing
type ExecutionInfo struct {
	RunID                 string     `json:"run_id"`
	WorkflowDefinitionID  string     `json:"workflow_definition_id"`
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
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
