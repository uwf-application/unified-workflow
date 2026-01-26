package primitive

import (
	"time"
)

// WorkflowContext represents workflow execution metadata (immutable)
// Similar to Java's WorkflowContext record
type WorkflowContext interface {
	// GetRunID returns the workflow run ID
	GetRunID() string

	// GetWorkflowDefinitionID returns the workflow definition ID
	GetWorkflowDefinitionID() string

	// GetStatus returns the workflow status
	GetStatus() WorkflowStatus

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
	WithStatus(status WorkflowStatus) WorkflowContext

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

// WorkflowData represents shared workflow data (mutable)
// Similar to Java's WorkflowData class
type WorkflowData interface {
	// Get returns a value by key
	Get(key string) interface{}

	// Put sets a value by key
	Put(key string, value interface{})

	// Remove removes a value by key
	Remove(key string)

	// Contains checks if a key exists
	Contains(key string) bool

	// Size returns the number of entries
	Size() int

	// Clear removes all entries
	Clear()

	// ToMap returns a copy of the data as a map
	ToMap() map[string]interface{}

	// DeepCopy creates a deep copy of the workflow data
	DeepCopy() WorkflowData

	// Merge merges another WorkflowData into this one
	Merge(other WorkflowData)

	// GetString returns a string value by key
	GetString(key string) (string, bool)

	// GetInt returns an int value by key
	GetInt(key string) (int, bool)

	// GetBool returns a bool value by key
	GetBool(key string) (bool, bool)

	// GetFloat returns a float64 value by key
	GetFloat(key string) (float64, bool)

	// GetMap returns a map value by key
	GetMap(key string) (map[string]interface{}, bool)

	// GetSlice returns a slice value by key
	GetSlice(key string) ([]interface{}, bool)
}

// WorkflowStatus represents the status of a workflow
type WorkflowStatus int

const (
	WorkflowStatusPending WorkflowStatus = iota
	WorkflowStatusRunning
	WorkflowStatusCompleted
	WorkflowStatusFailed
	WorkflowStatusCancelled
	WorkflowStatusPaused
)

// StepLogic is a functional interface for step execution logic
type StepLogic func() error

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	WithFields(fields map[string]interface{}) Logger
}

// Metrics defines the metrics collection interface
type Metrics interface {
	IncCounter(name string, labels map[string]string)
	ObserveHistogram(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
	RecordTiming(name string, duration time.Duration, labels map[string]string)
}
