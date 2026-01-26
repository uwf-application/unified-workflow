package state

import (
	"context"
	"sync"
	"time"

	"unified-workflow/internal/primitive/model"
)

// InMemoryState implements the StateManagement interface using in-memory storage
type InMemoryState struct {
	mu               sync.RWMutex
	contexts         map[string]model.WorkflowContext
	data             map[string]model.WorkflowData
	locks            map[string]bool
	ttl              map[string]time.Time
	contextCreatedAt map[string]time.Time
	contextUpdatedAt map[string]time.Time
}

// NewInMemoryState creates a new in-memory state management
func NewInMemoryState() *InMemoryState {
	return &InMemoryState{
		contexts:         make(map[string]model.WorkflowContext),
		data:             make(map[string]model.WorkflowData),
		locks:            make(map[string]bool),
		ttl:              make(map[string]time.Time),
		contextCreatedAt: make(map[string]time.Time),
		contextUpdatedAt: make(map[string]time.Time),
	}
}

// SaveContext saves the workflow context to the store
func (s *InMemoryState) SaveContext(ctx context.Context, workflowContext model.WorkflowContext) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	runID := workflowContext.GetRunID()
	now := time.Now()

	if _, exists := s.contexts[runID]; !exists {
		s.contextCreatedAt[runID] = now
	}
	s.contexts[runID] = workflowContext
	s.contextUpdatedAt[runID] = now

	// Check TTL
	if expiry, ok := s.ttl[runID]; ok && now.After(expiry) {
		delete(s.contexts, runID)
		delete(s.data, runID)
		delete(s.ttl, runID)
		delete(s.contextCreatedAt, runID)
		delete(s.contextUpdatedAt, runID)
	}

	return nil
}

// GetContext retrieves the workflow context for the given run ID
func (s *InMemoryState) GetContext(ctx context.Context, runID string) (model.WorkflowContext, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check TTL
	if expiry, ok := s.ttl[runID]; ok && time.Now().After(expiry) {
		delete(s.contexts, runID)
		delete(s.data, runID)
		delete(s.ttl, runID)
		delete(s.contextCreatedAt, runID)
		delete(s.contextUpdatedAt, runID)
		return nil, ErrStateNotFound
	}

	context, exists := s.contexts[runID]
	if !exists {
		return nil, ErrStateNotFound
	}

	return context, nil
}

// SaveData saves the workflow data to the store
func (s *InMemoryState) SaveData(ctx context.Context, runID string, workflowData model.WorkflowData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check TTL
	if expiry, ok := s.ttl[runID]; ok && time.Now().After(expiry) {
		delete(s.contexts, runID)
		delete(s.data, runID)
		delete(s.ttl, runID)
		delete(s.contextCreatedAt, runID)
		delete(s.contextUpdatedAt, runID)
		return ErrStateExpired
	}

	s.data[runID] = workflowData
	return nil
}

// GetData retrieves the workflow data for the given run ID
func (s *InMemoryState) GetData(ctx context.Context, runID string) (model.WorkflowData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check TTL
	if expiry, ok := s.ttl[runID]; ok && time.Now().After(expiry) {
		delete(s.contexts, runID)
		delete(s.data, runID)
		delete(s.ttl, runID)
		delete(s.contextCreatedAt, runID)
		delete(s.contextUpdatedAt, runID)
		return nil, ErrStateNotFound
	}

	data, exists := s.data[runID]
	if !exists {
		return nil, ErrStateNotFound
	}

	return data, nil
}

// RemoveState removes all state (both context and data) for the given run ID
func (s *InMemoryState) RemoveState(ctx context.Context, runID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.contexts, runID)
	delete(s.data, runID)
	delete(s.locks, runID)
	delete(s.ttl, runID)
	delete(s.contextCreatedAt, runID)
	delete(s.contextUpdatedAt, runID)

	return nil
}

// ContainsContext checks if a workflow context exists for the given run ID
func (s *InMemoryState) ContainsContext(ctx context.Context, runID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check TTL
	if expiry, ok := s.ttl[runID]; ok && time.Now().After(expiry) {
		delete(s.contexts, runID)
		delete(s.data, runID)
		delete(s.ttl, runID)
		delete(s.contextCreatedAt, runID)
		delete(s.contextUpdatedAt, runID)
		return false, nil
	}

	_, exists := s.contexts[runID]
	return exists, nil
}

// ContainsData checks if workflow data exists for the given run ID
func (s *InMemoryState) ContainsData(ctx context.Context, runID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check TTL
	if expiry, ok := s.ttl[runID]; ok && time.Now().After(expiry) {
		delete(s.contexts, runID)
		delete(s.data, runID)
		delete(s.ttl, runID)
		delete(s.contextCreatedAt, runID)
		delete(s.contextUpdatedAt, runID)
		return false, nil
	}

	_, exists := s.data[runID]
	return exists, nil
}

// AcquireLock acquires a lock for a workflow run to prevent concurrent modifications
func (s *InMemoryState) AcquireLock(ctx context.Context, runID string, timeout time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if already locked
	if locked, ok := s.locks[runID]; ok && locked {
		return false, nil
	}

	s.locks[runID] = true
	return true, nil
}

// ReleaseLock releases a lock for a workflow run
func (s *InMemoryState) ReleaseLock(ctx context.Context, runID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.locks, runID)
	return nil
}

// SetTTL sets time-to-live for workflow state
func (s *InMemoryState) SetTTL(ctx context.Context, runID string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ttl <= 0 {
		delete(s.ttl, runID)
	} else {
		s.ttl[runID] = time.Now().Add(ttl)
	}

	return nil
}

// GetAllContexts gets all workflow contexts stored in the state management
func (s *InMemoryState) GetAllContexts(ctx context.Context) ([]model.WorkflowContext, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	contexts := make([]model.WorkflowContext, 0, len(s.contexts))

	for runID, context := range s.contexts {
		// Check TTL
		if expiry, ok := s.ttl[runID]; ok && now.After(expiry) {
			continue // Skip expired contexts
		}
		contexts = append(contexts, context)
	}

	return contexts, nil
}

// Close closes the state management connection
func (s *InMemoryState) Close() error {
	// Nothing to close for in-memory state
	return nil
}

// GetExecutionInfo gets execution information for a run ID
func (s *InMemoryState) GetExecutionInfo(ctx context.Context, runID string) (*ExecutionInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check TTL
	if expiry, ok := s.ttl[runID]; ok && time.Now().After(expiry) {
		delete(s.contexts, runID)
		delete(s.data, runID)
		delete(s.ttl, runID)
		delete(s.contextCreatedAt, runID)
		delete(s.contextUpdatedAt, runID)
		return nil, ErrStateNotFound
	}

	context, exists := s.contexts[runID]
	if !exists {
		return nil, ErrStateNotFound
	}

	createdAt := time.Time{}
	if t, ok := s.contextCreatedAt[runID]; ok {
		createdAt = t
	}

	updatedAt := time.Time{}
	if t, ok := s.contextUpdatedAt[runID]; ok {
		updatedAt = t
	}

	status := context.GetStatus()
	statusStr := workflowStatusToString(status)
	isTerminal := isWorkflowStatusTerminal(status)

	var startTime, endTime *time.Time
	if st := context.GetStartTime(); st != nil {
		startTime = st
	}
	if et := context.GetEndTime(); et != nil {
		endTime = et
	}

	return &ExecutionInfo{
		RunID:                 runID,
		WorkflowDefinitionID:  context.GetWorkflowDefinitionID(),
		Status:                statusStr,
		CurrentStepIndex:      context.GetCurrentStepIndex(),
		CurrentChildStepIndex: context.GetCurrentChildStepIndex(),
		StartTime:             startTime,
		EndTime:               endTime,
		ErrorMessage:          context.GetErrorMessage(),
		LastAttemptedStep:     context.GetLastAttemptedStep(),
		IsTerminal:            isTerminal,
		IsRunning:             status == model.WorkflowStatusRunning,
		IsPending:             status == model.WorkflowStatusPending,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}, nil
}

// GetAllExecutionInfos gets execution information for all runs
func (s *InMemoryState) GetAllExecutionInfos(ctx context.Context) ([]*ExecutionInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	infos := make([]*ExecutionInfo, 0, len(s.contexts))

	for runID, context := range s.contexts {
		// Check TTL
		if expiry, ok := s.ttl[runID]; ok && now.After(expiry) {
			continue // Skip expired contexts
		}

		createdAt := time.Time{}
		if t, ok := s.contextCreatedAt[runID]; ok {
			createdAt = t
		}

		updatedAt := time.Time{}
		if t, ok := s.contextUpdatedAt[runID]; ok {
			updatedAt = t
		}

		status := context.GetStatus()
		statusStr := workflowStatusToString(status)
		isTerminal := isWorkflowStatusTerminal(status)

		var startTime, endTime *time.Time
		if st := context.GetStartTime(); st != nil {
			startTime = st
		}
		if et := context.GetEndTime(); et != nil {
			endTime = et
		}

		infos = append(infos, &ExecutionInfo{
			RunID:                 runID,
			WorkflowDefinitionID:  context.GetWorkflowDefinitionID(),
			Status:                statusStr,
			CurrentStepIndex:      context.GetCurrentStepIndex(),
			CurrentChildStepIndex: context.GetCurrentChildStepIndex(),
			StartTime:             startTime,
			EndTime:               endTime,
			ErrorMessage:          context.GetErrorMessage(),
			LastAttemptedStep:     context.GetLastAttemptedStep(),
			IsTerminal:            isTerminal,
			IsRunning:             status == model.WorkflowStatusRunning,
			IsPending:             status == model.WorkflowStatusPending,
			CreatedAt:             createdAt,
			UpdatedAt:             updatedAt,
		})
	}

	return infos, nil
}

// Helper functions for workflow status
func workflowStatusToString(status int) string {
	switch status {
	case model.WorkflowStatusPending:
		return "pending"
	case model.WorkflowStatusRunning:
		return "running"
	case model.WorkflowStatusCompleted:
		return "completed"
	case model.WorkflowStatusFailed:
		return "failed"
	case model.WorkflowStatusCancelled:
		return "cancelled"
	case model.WorkflowStatusPaused:
		return "paused"
	default:
		return "unknown"
	}
}

func isWorkflowStatusTerminal(status int) bool {
	return status == model.WorkflowStatusCompleted ||
		status == model.WorkflowStatusFailed ||
		status == model.WorkflowStatusCancelled
}

// Errors
var (
	ErrStateNotFound = &StateError{Message: "state not found", Code: "NOT_FOUND"}
	ErrStateExpired  = &StateError{Message: "state expired", Code: "EXPIRED"}
	ErrStateLocked   = &StateError{Message: "state locked", Code: "LOCKED"}
)

// StateError represents a state management error
type StateError struct {
	Message string
	Code    string
}

func (e *StateError) Error() string {
	return e.Message
}
