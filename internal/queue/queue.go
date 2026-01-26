package queue

import (
	"context"
	"time"
)

// Message represents a message in the queue
type Message struct {
	ID        string
	RunID     string
	Data      []byte
	Timestamp time.Time
}

// Queue is the interface for workflow queue implementations
// Provides abstraction for different queue backends (in-memory, NATS JetStream, database, etc.)
type Queue interface {
	// Enqueue adds a workflow run ID to the queue for processing
	Enqueue(ctx context.Context, runID string, data []byte) error

	// Dequeue retrieves the next message from the queue
	Dequeue(ctx context.Context) (*Message, error)

	// Acknowledge successful processing of a message
	Acknowledge(ctx context.Context, messageID string) error

	// Reject rejects a message and schedules it for retry
	Reject(ctx context.Context, messageID string, delay time.Duration) error

	// Size returns the current size of the queue
	Size(ctx context.Context) (int, error)

	// IsEmpty checks if the queue is empty
	IsEmpty(ctx context.Context) (bool, error)

	// Contains checks if a specific run ID is in the queue
	Contains(ctx context.Context, runID string) (bool, error)

	// Remove removes a specific run ID from the queue
	Remove(ctx context.Context, runID string) error

	// Clear removes all messages from the queue
	Clear(ctx context.Context) error

	// Close closes the queue connection
	Close() error
}

// ExecutionRequest represents a workflow execution request
type ExecutionRequest struct {
	RunID       string                 `json:"run_id"`
	WorkflowID  string                 `json:"workflow_id"`
	InputData   map[string]interface{} `json:"input_data"`
	RequestedAt time.Time              `json:"requested_at"`
}

// ExecutionResult represents a workflow execution result
type ExecutionResult struct {
	RunID       string                 `json:"run_id"`
	WorkflowID  string                 `json:"workflow_id"`
	Status      string                 `json:"status"`
	OutputData  map[string]interface{} `json:"output_data,omitempty"`
	Error       string                 `json:"error,omitempty"`
	CompletedAt time.Time              `json:"completed_at"`
}
