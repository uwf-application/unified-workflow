package model

// WorkflowStatus constants (similar to Java's WorkflowStatus enum)
const (
	WorkflowStatusPending   = 0
	WorkflowStatusRunning   = 1
	WorkflowStatusCompleted = 2
	WorkflowStatusFailed    = 3
	WorkflowStatusCancelled = 4
	WorkflowStatusPaused    = 5
)

// StepStatus constants (similar to Java's StepStatus enum)
const (
	StepStatusPending   = "pending"
	StepStatusRunning   = "running"
	StepStatusCompleted = "completed"
	StepStatusFailed    = "failed"
	StepStatusSkipped   = "skipped"
	StepStatusCancelled = "cancelled"
)
