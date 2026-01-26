package registry

import (
	"context"
	"sync"
	"time"

	"unified-workflow/internal/common/model"
)

// InMemoryRegistry implements the Registry interface using in-memory storage
type InMemoryRegistry struct {
	mu        sync.RWMutex
	workflows map[string]model.Workflow
	createdAt map[string]time.Time
	updatedAt map[string]time.Time
}

// NewInMemoryRegistry creates a new in-memory registry
func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		workflows: make(map[string]model.Workflow),
		createdAt: make(map[string]time.Time),
		updatedAt: make(map[string]time.Time),
	}
}

// RegisterWorkflow registers a workflow with the registry
func (r *InMemoryRegistry) RegisterWorkflow(ctx context.Context, workflow model.Workflow) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	workflowID := workflow.GetID()
	now := time.Now()

	if _, exists := r.workflows[workflowID]; !exists {
		r.createdAt[workflowID] = now
	}
	r.workflows[workflowID] = workflow
	r.updatedAt[workflowID] = now

	return nil
}

// GetWorkflow retrieves a workflow by its ID
func (r *InMemoryRegistry) GetWorkflow(ctx context.Context, workflowID string) (model.Workflow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workflow, exists := r.workflows[workflowID]
	if !exists {
		return nil, ErrWorkflowNotFound
	}

	return workflow, nil
}

// ContainsWorkflow checks if a workflow exists in the registry
func (r *InMemoryRegistry) ContainsWorkflow(ctx context.Context, workflowID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.workflows[workflowID]
	return exists, nil
}

// RemoveWorkflow removes a workflow from the registry
func (r *InMemoryRegistry) RemoveWorkflow(ctx context.Context, workflowID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workflows[workflowID]; !exists {
		return ErrWorkflowNotFound
	}

	delete(r.workflows, workflowID)
	delete(r.createdAt, workflowID)
	delete(r.updatedAt, workflowID)

	return nil
}

// GetAllWorkflowIDs gets all registered workflow IDs
func (r *InMemoryRegistry) GetAllWorkflowIDs(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.workflows))
	for id := range r.workflows {
		ids = append(ids, id)
	}

	return ids, nil
}

// GetWorkflowCount gets the number of registered workflows
func (r *InMemoryRegistry) GetWorkflowCount(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.workflows), nil
}

// Clear removes all workflows from the registry
func (r *InMemoryRegistry) Clear(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.workflows = make(map[string]model.Workflow)
	r.createdAt = make(map[string]time.Time)
	r.updatedAt = make(map[string]time.Time)

	return nil
}

// Initialize initializes the registry
func (r *InMemoryRegistry) Initialize(ctx context.Context) error {
	// Nothing to initialize for in-memory registry
	return nil
}

// Shutdown shuts down the registry
func (r *InMemoryRegistry) Shutdown(ctx context.Context) error {
	// Nothing to shutdown for in-memory registry
	return nil
}

// GetWorkflowInfo gets workflow information for listing
func (r *InMemoryRegistry) GetWorkflowInfo(ctx context.Context, workflowID string) (*WorkflowInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workflow, exists := r.workflows[workflowID]
	if !exists {
		return nil, ErrWorkflowNotFound
	}

	createdAt := ""
	if t, ok := r.createdAt[workflowID]; ok {
		createdAt = t.Format(time.RFC3339)
	}

	updatedAt := ""
	if t, ok := r.updatedAt[workflowID]; ok {
		updatedAt = t.Format(time.RFC3339)
	}

	return &WorkflowInfo{
		ID:          workflowID,
		Name:        workflow.GetName(),
		Description: workflow.GetDescription(),
		StepCount:   workflow.GetStepCount(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

// GetAllWorkflowInfos gets information for all workflows
func (r *InMemoryRegistry) GetAllWorkflowInfos(ctx context.Context) ([]*WorkflowInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]*WorkflowInfo, 0, len(r.workflows))
	for workflowID, workflow := range r.workflows {
		createdAt := ""
		if t, ok := r.createdAt[workflowID]; ok {
			createdAt = t.Format(time.RFC3339)
		}

		updatedAt := ""
		if t, ok := r.updatedAt[workflowID]; ok {
			updatedAt = t.Format(time.RFC3339)
		}

		infos = append(infos, &WorkflowInfo{
			ID:          workflowID,
			Name:        workflow.GetName(),
			Description: workflow.GetDescription(),
			StepCount:   workflow.GetStepCount(),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	return infos, nil
}

// Errors
var (
	ErrWorkflowNotFound = &RegistryError{Message: "workflow not found", Code: "NOT_FOUND"}
)

// RegistryError represents a registry error
type RegistryError struct {
	Message string
	Code    string
}

func (e *RegistryError) Error() string {
	return e.Message
}
