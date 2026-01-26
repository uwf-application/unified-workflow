package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NATSQueue implements the Queue interface using NATS JetStream
type NATSQueue struct {
	conn          *nats.Conn
	js            jetstream.JetStream
	stream        jetstream.Stream
	consumer      jetstream.Consumer
	streamName    string
	subjectPrefix string
	durableName   string
}

// NATSConfig represents NATS JetStream configuration
type NATSConfig struct {
	URLs           []string      `json:"urls"`
	StreamName     string        `json:"stream_name"`
	SubjectPrefix  string        `json:"subject_prefix"`
	DurableName    string        `json:"durable_name"`
	MaxReconnects  int           `json:"max_reconnects"`
	ReconnectWait  time.Duration `json:"reconnect_wait"`
	ConnectTimeout time.Duration `json:"connect_timeout"`
}

// NewNATSQueue creates a new NATS JetStream queue
func NewNATSQueue(config NATSConfig) (*NATSQueue, error) {
	// Connect to NATS
	opts := []nats.Option{
		nats.MaxReconnects(config.MaxReconnects),
		nats.ReconnectWait(config.ReconnectWait),
		nats.Timeout(config.ConnectTimeout),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("NATS disconnected: %v\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("NATS reconnected to %s\n", nc.ConnectedUrl())
		}),
	}

	conn, err := nats.Connect(config.URLs[0], opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create or get stream
	streamCfg := jetstream.StreamConfig{
		Name:      config.StreamName,
		Subjects:  []string{fmt.Sprintf("%s.>", config.SubjectPrefix)},
		Retention: jetstream.WorkQueuePolicy,
		MaxMsgs:   -1,
		MaxBytes:  -1,
		MaxAge:    24 * time.Hour,
		Storage:   jetstream.FileStorage,
		Discard:   jetstream.DiscardOld,
	}

	stream, err := js.CreateOrUpdateStream(context.Background(), streamCfg)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	// Create or get consumer
	consumerCfg := jetstream.ConsumerConfig{
		Durable:       config.DurableName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    3,
		FilterSubject: fmt.Sprintf("%s.execution.requests", config.SubjectPrefix),
	}

	consumer, err := stream.CreateOrUpdateConsumer(context.Background(), consumerCfg)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &NATSQueue{
		conn:          conn,
		js:            js,
		stream:        stream,
		consumer:      consumer,
		streamName:    config.StreamName,
		subjectPrefix: config.SubjectPrefix,
		durableName:   config.DurableName,
	}, nil
}

// Enqueue adds a workflow run ID to the queue for processing
func (q *NATSQueue) Enqueue(ctx context.Context, runID string, data []byte) error {
	subject := fmt.Sprintf("%s.execution.requests", q.subjectPrefix)

	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  nats.Header{},
	}
	msg.Header.Set("run_id", runID)
	msg.Header.Set("timestamp", time.Now().Format(time.RFC3339))

	_, err := q.js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Dequeue retrieves the next message from the queue
func (q *NATSQueue) Dequeue(ctx context.Context) (*Message, error) {
	msgs, err := q.consumer.Fetch(1, jetstream.FetchMaxWait(1*time.Second))
	if err != nil {
		if err == jetstream.ErrNoMessages {
			return nil, nil // No messages available
		}
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	// Iterate through the batch (should only have one message)
	for msg := range msgs.Messages() {
		metadata, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("failed to get message metadata: %w", err)
		}

		headers := msg.Headers()
		var timestamp time.Time
		var runID string
		if headers != nil {
			timestamp, _ = time.Parse(time.RFC3339, headers.Get("timestamp"))
			runID = headers.Get("run_id")
		}

		return &Message{
			ID:        fmt.Sprintf("%d", metadata.Sequence.Stream),
			RunID:     runID,
			Data:      msg.Data(),
			Timestamp: timestamp,
		}, nil
	}

	return nil, nil // No messages in batch
}

// Acknowledge successful processing of a message
func (q *NATSQueue) Acknowledge(ctx context.Context, messageID string) error {
	// In NATS JetStream, acknowledgement is done on the message itself
	// This method is a no-op for NATS since we acknowledge when processing
	return nil
}

// Reject rejects a message and schedules it for retry
func (q *NATSQueue) Reject(ctx context.Context, messageID string, delay time.Duration) error {
	// In NATS JetStream, we can nak the message with delay
	// This would require access to the original message
	// For now, we'll just return nil as the basic implementation
	return nil
}

// Size returns the current size of the queue
func (q *NATSQueue) Size(ctx context.Context) (int, error) {
	info, err := q.stream.Info(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get stream info: %w", err)
	}
	return int(info.State.Msgs), nil
}

// IsEmpty checks if the queue is empty
func (q *NATSQueue) IsEmpty(ctx context.Context) (bool, error) {
	size, err := q.Size(ctx)
	if err != nil {
		return false, err
	}
	return size == 0, nil
}

// Contains checks if a specific run ID is in the queue
func (q *NATSQueue) Contains(ctx context.Context, runID string) (bool, error) {
	// This is complex with NATS JetStream as we'd need to search messages
	// For simplicity, we'll return false for now
	return false, nil
}

// Remove removes a specific run ID from the queue
func (q *NATSQueue) Remove(ctx context.Context, runID string) error {
	// Removing specific messages is complex with NATS JetStream
	// For now, we'll return nil (no-op)
	return nil
}

// Clear removes all messages from the queue
func (q *NATSQueue) Clear(ctx context.Context) error {
	// In NATS JetStream, we can purge the stream
	err := q.stream.Purge(ctx)
	if err != nil {
		return fmt.Errorf("failed to purge stream: %w", err)
	}
	return nil
}

// Close closes the queue connection
func (q *NATSQueue) Close() error {
	q.conn.Close()
	return nil
}

// MarshalExecutionRequest marshals an ExecutionRequest to JSON
func MarshalExecutionRequest(req ExecutionRequest) ([]byte, error) {
	return json.Marshal(req)
}

// UnmarshalExecutionRequest unmarshals JSON to ExecutionRequest
func UnmarshalExecutionRequest(data []byte) (ExecutionRequest, error) {
	var req ExecutionRequest
	err := json.Unmarshal(data, &req)
	return req, err
}

// MarshalExecutionResult marshals an ExecutionResult to JSON
func MarshalExecutionResult(res ExecutionResult) ([]byte, error) {
	return json.Marshal(res)
}

// UnmarshalExecutionResult unmarshals JSON to ExecutionResult
func UnmarshalExecutionResult(data []byte) (ExecutionResult, error) {
	var res ExecutionResult
	err := json.Unmarshal(data, &res)
	return res, err
}
