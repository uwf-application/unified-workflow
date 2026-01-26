package model

import (
	"time"
)

// WorkflowContext represents workflow execution metadata (immutable)
// Similar to Java's WorkflowContext record in workflow-primitive
type WorkflowContext interface {
	// GetRunID returns the workflow run ID
	GetRunID() string

	// GetWorkflowDefinitionID returns the workflow definition ID
	GetWorkflowDefinitionID() string

	// GetStatus returns the workflow status
	GetStatus() int

	// GetCurrentStepIndex returns the current step index
	GetCurrentStepIndex() int

	// GetCurrentChildStepIndex returns the current child step index
	GetCurrentChildStepIndex() int

	// GetStartTime returns the start time
	GetStartTime() *time.Time

	// GetEndTime returns the end time
	GetEndTime() *time.Time

	// GetErrorMessage returns the error message if any
	GetErrorMessage() string

	// GetLastAttemptedStep returns the last attempted step name
	GetLastAttemptedStep() string

	// WithStatus creates a new context with updated status
	WithStatus(status int) WorkflowContext

	// WithIndices creates a new context with updated step indices
	WithIndices(stepIndex, childStepIndex int) WorkflowContext

	// WithErrorMessage creates a new context with updated error message
	WithErrorMessage(errorMessage string) WorkflowContext

	// WithStartTime creates a new context with updated start time
	WithStartTime(startTime time.Time) WorkflowContext

	// WithEndTime creates a new context with updated end time
	WithEndTime(endTime time.Time) WorkflowContext

	// WithCurrentStepIndex creates a new context with updated step index
	WithCurrentStepIndex(stepIndex int) WorkflowContext

	// WithCurrentChildStepIndex creates a new context with updated child step index
	WithCurrentChildStepIndex(childStepIndex int) WorkflowContext

	// WithLastAttemptedStep creates a new context with updated last attempted step
	WithLastAttemptedStep(stepName string) WorkflowContext
}

// WorkflowContextImpl implements the WorkflowContext interface
type WorkflowContextImpl struct {
	runID                 string
	workflowDefinitionID  string
	status                int // Using int for WorkflowStatus
	currentStepIndex      int
	currentChildStepIndex int
	startTime             *time.Time
	endTime               *time.Time
	errorMessage          string
	lastAttemptedStep     string
}

// NewWorkflowContext creates a new workflow context
func NewWorkflowContext(workflowDefinitionID string) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 generateUUID(),
		workflowDefinitionID:  workflowDefinitionID,
		status:                0, // Pending
		currentStepIndex:      -1,
		currentChildStepIndex: -1,
	}
}

// GetRunID returns the workflow run ID
func (wc *WorkflowContextImpl) GetRunID() string {
	return wc.runID
}

// GetWorkflowDefinitionID returns the workflow definition ID
func (wc *WorkflowContextImpl) GetWorkflowDefinitionID() string {
	return wc.workflowDefinitionID
}

// GetStatus returns the workflow status
func (wc *WorkflowContextImpl) GetStatus() int {
	return wc.status
}

// GetCurrentStepIndex returns the current step index
func (wc *WorkflowContextImpl) GetCurrentStepIndex() int {
	return wc.currentStepIndex
}

// GetCurrentChildStepIndex returns the current child step index
func (wc *WorkflowContextImpl) GetCurrentChildStepIndex() int {
	return wc.currentChildStepIndex
}

// GetStartTime returns the start time
func (wc *WorkflowContextImpl) GetStartTime() *time.Time {
	return wc.startTime
}

// GetEndTime returns the end time
func (wc *WorkflowContextImpl) GetEndTime() *time.Time {
	return wc.endTime
}

// GetErrorMessage returns the error message
func (wc *WorkflowContextImpl) GetErrorMessage() string {
	return wc.errorMessage
}

// GetLastAttemptedStep returns the last attempted step
func (wc *WorkflowContextImpl) GetLastAttemptedStep() string {
	return wc.lastAttemptedStep
}

// WithStatus creates a new context with updated status
func (wc *WorkflowContextImpl) WithStatus(status int) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 wc.runID,
		workflowDefinitionID:  wc.workflowDefinitionID,
		status:                status,
		currentStepIndex:      wc.currentStepIndex,
		currentChildStepIndex: wc.currentChildStepIndex,
		startTime:             wc.startTime,
		endTime:               wc.endTime,
		errorMessage:          wc.errorMessage,
		lastAttemptedStep:     wc.lastAttemptedStep,
	}
}

// WithIndices creates a new context with updated indices
func (wc *WorkflowContextImpl) WithIndices(stepIndex, childStepIndex int) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 wc.runID,
		workflowDefinitionID:  wc.workflowDefinitionID,
		status:                wc.status,
		currentStepIndex:      stepIndex,
		currentChildStepIndex: childStepIndex,
		startTime:             wc.startTime,
		endTime:               wc.endTime,
		errorMessage:          wc.errorMessage,
		lastAttemptedStep:     wc.lastAttemptedStep,
	}
}

// WithErrorMessage creates a new context with updated error message
func (wc *WorkflowContextImpl) WithErrorMessage(errorMessage string) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 wc.runID,
		workflowDefinitionID:  wc.workflowDefinitionID,
		status:                wc.status,
		currentStepIndex:      wc.currentStepIndex,
		currentChildStepIndex: wc.currentChildStepIndex,
		startTime:             wc.startTime,
		endTime:               wc.endTime,
		errorMessage:          errorMessage,
		lastAttemptedStep:     wc.lastAttemptedStep,
	}
}

// WithStartTime creates a new context with updated start time
func (wc *WorkflowContextImpl) WithStartTime(startTime time.Time) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 wc.runID,
		workflowDefinitionID:  wc.workflowDefinitionID,
		status:                wc.status,
		currentStepIndex:      wc.currentStepIndex,
		currentChildStepIndex: wc.currentChildStepIndex,
		startTime:             &startTime,
		endTime:               wc.endTime,
		errorMessage:          wc.errorMessage,
		lastAttemptedStep:     wc.lastAttemptedStep,
	}
}

// WithEndTime creates a new context with updated end time
func (wc *WorkflowContextImpl) WithEndTime(endTime time.Time) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 wc.runID,
		workflowDefinitionID:  wc.workflowDefinitionID,
		status:                wc.status,
		currentStepIndex:      wc.currentStepIndex,
		currentChildStepIndex: wc.currentChildStepIndex,
		startTime:             wc.startTime,
		endTime:               &endTime,
		errorMessage:          wc.errorMessage,
		lastAttemptedStep:     wc.lastAttemptedStep,
	}
}

// WithCurrentStepIndex creates a new context with updated step index
func (wc *WorkflowContextImpl) WithCurrentStepIndex(stepIndex int) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 wc.runID,
		workflowDefinitionID:  wc.workflowDefinitionID,
		status:                wc.status,
		currentStepIndex:      stepIndex,
		currentChildStepIndex: wc.currentChildStepIndex,
		startTime:             wc.startTime,
		endTime:               wc.endTime,
		errorMessage:          wc.errorMessage,
		lastAttemptedStep:     wc.lastAttemptedStep,
	}
}

// WithCurrentChildStepIndex creates a new context with updated child step index
func (wc *WorkflowContextImpl) WithCurrentChildStepIndex(childStepIndex int) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 wc.runID,
		workflowDefinitionID:  wc.workflowDefinitionID,
		status:                wc.status,
		currentStepIndex:      wc.currentStepIndex,
		currentChildStepIndex: childStepIndex,
		startTime:             wc.startTime,
		endTime:               wc.endTime,
		errorMessage:          wc.errorMessage,
		lastAttemptedStep:     wc.lastAttemptedStep,
	}
}

// WithLastAttemptedStep creates a new context with updated last attempted step
func (wc *WorkflowContextImpl) WithLastAttemptedStep(stepName string) *WorkflowContextImpl {
	return &WorkflowContextImpl{
		runID:                 wc.runID,
		workflowDefinitionID:  wc.workflowDefinitionID,
		status:                wc.status,
		currentStepIndex:      wc.currentStepIndex,
		currentChildStepIndex: wc.currentChildStepIndex,
		startTime:             wc.startTime,
		endTime:               wc.endTime,
		errorMessage:          wc.errorMessage,
		lastAttemptedStep:     stepName,
	}
}

// Helper function to generate UUID
func generateUUID() string {
	return "uuid-" + time.Now().Format("20060102150405") + "-" + randomString(8)
}

// Helper function to generate random string
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
