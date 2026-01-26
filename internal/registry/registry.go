package registry

import (
	"context"

	"unified-workflow/internal/common/model"
)

// Registry is the interface for workflow registry implementations
// Provides abstraction for different storage backends (in-memory, database, etc.)
type Registry interface {
	// RegisterWorkflow registers a workflow with the registry
	RegisterWorkflow(ctx context.Context, workflow model.Workflow) error

	// GetWorkflow retrieves a workflow by its ID
	GetWorkflow(ctx context.Context, workflowID string) (model.Workflow, error)

	// ContainsWorkflow checks if a workflow exists in the registry
	ContainsWorkflow(ctx context.Context, workflowID string) (bool, error)

	// RemoveWorkflow removes a workflow from the registry
	RemoveWorkflow(ctx context.Context, workflowID string) error

	// GetAllWorkflowIDs gets all registered workflow IDs
	GetAllWorkflowIDs(ctx context.Context) ([]string, error)

	// GetWorkflowCount gets the number of registered workflows
	GetWorkflowCount(ctx context.Context) (int, error)

	// Clear removes all workflows from the registry
	Clear(ctx context.Context) error

	// Initialize initializes the registry (e.g., creates tables, loads data)
	Initialize(ctx context.Context) error

	// Shutdown shuts down the registry (e.g., closes connections)
	Shutdown(ctx context.Context) error
}

// WorkflowInfo represents simplified workflow information for listing
type WorkflowInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StepCount   int    `json:"step_count"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
