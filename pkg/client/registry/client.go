package registry

import (
	"context"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/pkg/client"
)

// Client is the interface for registry service client
type Client interface {
	client.Client

	// ListWorkflows lists all registered workflows
	ListWorkflows(ctx context.Context, req *ListWorkflowsRequest) (*ListWorkflowsResponse, error)

	// GetWorkflow gets a specific workflow by ID
	GetWorkflow(ctx context.Context, req *GetWorkflowRequest) (*GetWorkflowResponse, error)

	// CreateWorkflow creates a new workflow
	CreateWorkflow(ctx context.Context, req *CreateWorkflowRequest) (*CreateWorkflowResponse, error)

	// UpdateWorkflow updates an existing workflow
	UpdateWorkflow(ctx context.Context, req *UpdateWorkflowRequest) (*UpdateWorkflowResponse, error)

	// DeleteWorkflow deletes a workflow
	DeleteWorkflow(ctx context.Context, req *DeleteWorkflowRequest) (*DeleteWorkflowResponse, error)

	// ContainsWorkflow checks if a workflow exists
	ContainsWorkflow(ctx context.Context, req *ContainsWorkflowRequest) (*ContainsWorkflowResponse, error)

	// GetWorkflowCount gets the number of registered workflows
	GetWorkflowCount(ctx context.Context, req *GetWorkflowCountRequest) (*GetWorkflowCountResponse, error)

	// Clear removes all workflows from the registry
	Clear(ctx context.Context, req *ClearRequest) (*ClearResponse, error)
}

// ListWorkflowsRequest is the request for listing workflows
type ListWorkflowsRequest struct {
	client.Request

	// Filter by workflow name (optional)
	NameFilter string `json:"name_filter,omitempty"`

	// Filter by workflow description (optional)
	DescriptionFilter string `json:"description_filter,omitempty"`

	// Limit the number of results (optional)
	Limit int `json:"limit,omitempty"`

	// Offset for pagination (optional)
	Offset int `json:"offset,omitempty"`
}

// ListWorkflowsResponse is the response for listing workflows
type ListWorkflowsResponse struct {
	client.Response

	// Workflows is the list of workflows
	Workflows []*WorkflowInfo `json:"workflows"`

	// TotalCount is the total number of workflows
	TotalCount int `json:"total_count"`

	// FilteredCount is the number of workflows after filtering
	FilteredCount int `json:"filtered_count"`
}

// GetWorkflowRequest is the request for getting a workflow
type GetWorkflowRequest struct {
	client.Request

	// ID is the workflow ID
	ID string `json:"id"`
}

// GetWorkflowResponse is the response for getting a workflow
type GetWorkflowResponse struct {
	client.Response

	// Workflow is the workflow information
	Workflow *WorkflowDetail `json:"workflow"`
}

// CreateWorkflowRequest is the request for creating a workflow
type CreateWorkflowRequest struct {
	client.Request

	// Name is the workflow name
	Name string `json:"name"`

	// Description is the workflow description
	Description string `json:"description,omitempty"`

	// Steps are the workflow steps (optional, can be added later)
	Steps []model.Step `json:"steps,omitempty"`

	// Metadata contains additional workflow metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CreateWorkflowResponse is the response for creating a workflow
type CreateWorkflowResponse struct {
	client.Response

	// Workflow is the created workflow information
	Workflow *WorkflowInfo `json:"workflow"`
}

// UpdateWorkflowRequest is the request for updating a workflow
type UpdateWorkflowRequest struct {
	client.Request

	// ID is the workflow ID
	ID string `json:"id"`

	// Name is the updated workflow name (optional)
	Name *string `json:"name,omitempty"`

	// Description is the updated workflow description (optional)
	Description *string `json:"description,omitempty"`

	// Steps are the updated workflow steps (optional)
	Steps []model.Step `json:"steps,omitempty"`

	// Metadata contains updated workflow metadata (optional)
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateWorkflowResponse is the response for updating a workflow
type UpdateWorkflowResponse struct {
	client.Response

	// Workflow is the updated workflow information
	Workflow *WorkflowInfo `json:"workflow"`
}

// DeleteWorkflowRequest is the request for deleting a workflow
type DeleteWorkflowRequest struct {
	client.Request

	// ID is the workflow ID
	ID string `json:"id"`
}

// DeleteWorkflowResponse is the response for deleting a workflow
type DeleteWorkflowResponse struct {
	client.Response

	// Deleted indicates if the workflow was deleted
	Deleted bool `json:"deleted"`
}

// ContainsWorkflowRequest is the request for checking if a workflow exists
type ContainsWorkflowRequest struct {
	client.Request

	// ID is the workflow ID
	ID string `json:"id"`
}

// ContainsWorkflowResponse is the response for checking if a workflow exists
type ContainsWorkflowResponse struct {
	client.Response

	// Exists indicates if the workflow exists
	Exists bool `json:"exists"`
}

// GetWorkflowCountRequest is the request for getting workflow count
type GetWorkflowCountRequest struct {
	client.Request
}

// GetWorkflowCountResponse is the response for getting workflow count
type GetWorkflowCountResponse struct {
	client.Response

	// Count is the number of workflows
	Count int `json:"count"`
}

// ClearRequest is the request for clearing all workflows
type ClearRequest struct {
	client.Request
}

// ClearResponse is the response for clearing all workflows
type ClearResponse struct {
	client.Response

	// ClearedCount is the number of workflows cleared
	ClearedCount int `json:"cleared_count"`
}

// WorkflowInfo represents simplified workflow information
type WorkflowInfo struct {
	// ID is the workflow ID
	ID string `json:"id"`

	// Name is the workflow name
	Name string `json:"name"`

	// Description is the workflow description
	Description string `json:"description,omitempty"`

	// StepCount is the number of steps in the workflow
	StepCount int `json:"step_count"`

	// CreatedAt is the creation timestamp
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the last update timestamp
	UpdatedAt time.Time `json:"updated_at"`

	// Metadata contains additional workflow metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// WorkflowDetail represents detailed workflow information
type WorkflowDetail struct {
	WorkflowInfo

	// Steps are the workflow steps
	Steps []*StepInfo `json:"steps,omitempty"`

	// Primitives contains primitive information
	Primitives interface{} `json:"primitives,omitempty"`

	// Context contains workflow context
	Context interface{} `json:"context,omitempty"`

	// Data contains workflow data
	Data interface{} `json:"data,omitempty"`
}

// StepInfo represents step information
type StepInfo struct {
	// Name is the step name
	Name string `json:"name"`

	// Type is the step type (sequential, parallel, etc.)
	Type string `json:"type"`

	// ChildStepCount is the number of child steps
	ChildStepCount int `json:"child_step_count"`

	// IsParallel indicates if the step is parallel
	IsParallel bool `json:"is_parallel"`

	// Description is the step description
	Description string `json:"description,omitempty"`

	// Parameters contains step parameters
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// Config is the registry client configuration
type Config struct {
	client.Config

	// Namespace is the registry namespace (optional)
	Namespace string `json:"namespace" yaml:"namespace"`
}

// DefaultConfig returns the default registry client configuration
func DefaultConfig() Config {
	return Config{
		Config: client.DefaultConfig(),
	}
}
