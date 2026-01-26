package queue

import (
	"context"
	"sync"
	"time"
)

// InMemoryQueue implements the Queue interface using in-memory storage
type InMemoryQueue struct {
	mu       sync.RWMutex
	messages []*Message
	runIDMap map[string]*Message
}

// NewInMemoryQueue creates a new in-memory queue
func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		messages: make([]*Message, 0),
		runIDMap: make(map[string]*Message),
	}
}

// Enqueue adds a workflow run ID to the queue for processing
func (q *InMemoryQueue) Enqueue(ctx context.Context, runID string, data []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	msg := &Message{
		ID:        generateMessageID(),
		RunID:     runID,
		Data:      data,
		Timestamp: time.Now(),
	}

	q.messages = append(q.messages, msg)
	q.runIDMap[runID] = msg

	return nil
}

// Dequeue retrieves the next message from the queue
func (q *InMemoryQueue) Dequeue(ctx context.Context) (*Message, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.messages) == 0 {
		return nil, nil
	}

	msg := q.messages[0]
	q.messages = q.messages[1:]
	delete(q.runIDMap, msg.RunID)

	return msg, nil
}

// Acknowledge successful processing of a message
func (q *InMemoryQueue) Acknowledge(ctx context.Context, messageID string) error {
	// For in-memory queue, messages are removed on dequeue
	return nil
}

// Reject rejects a message and schedules it for retry
func (q *InMemoryQueue) Reject(ctx context.Context, messageID string, delay time.Duration) error {
	// For in-memory queue, we don't implement retry logic
	return nil
}

// Size returns the current size of the queue
func (q *InMemoryQueue) Size(ctx context.Context) (int, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.messages), nil
}

// IsEmpty checks if the queue is empty
func (q *InMemoryQueue) IsEmpty(ctx context.Context) (bool, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.messages) == 0, nil
}

// Contains checks if a specific run ID is in the queue
func (q *InMemoryQueue) Contains(ctx context.Context, runID string) (bool, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	_, exists := q.runIDMap[runID]
	return exists, nil
}

// Remove removes a specific run ID from the queue
func (q *InMemoryQueue) Remove(ctx context.Context, runID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Find and remove the message
	for i, msg := range q.messages {
		if msg.RunID == runID {
			q.messages = append(q.messages[:i], q.messages[i+1:]...)
			delete(q.runIDMap, runID)
			return nil
		}
	}

	return nil
}

// Clear removes all messages from the queue
func (q *InMemoryQueue) Clear(ctx context.Context) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.messages = make([]*Message, 0)
	q.runIDMap = make(map[string]*Message)

	return nil
}

// Close closes the queue connection
func (q *InMemoryQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.messages = make([]*Message, 0)
	q.runIDMap = make(map[string]*Message)

	return nil
}

// Helper function to generate message ID
func generateMessageID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
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
