package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// EnhancedMessage extends the standard Message with response routing info
type EnhancedMessage struct {
	*Message
	CorrelationID   string
	ResponseSubject string
	JetStreamMsg    jetstream.Msg // Internal handle for ACK/NACK
}

// EnhancedNATSQueue implements the Queue interface with response routing support
type EnhancedNATSQueue struct {
	conn          *nats.Conn
	js            jetstream.JetStream
	stream        jetstream.Stream
	consumer      jetstream.Consumer
	streamName    string
	subjectPrefix string
	durableName   string
}

// EnhancedNATSConfig represents enhanced NATS JetStream configuration
type EnhancedNATSConfig struct {
	URLs           []string      `json:"urls"`
	StreamName     string        `json:"stream_name"`
	SubjectPrefix  string        `json:"subject_prefix"`
	DurableName    string        `json:"durable_name"`
	MaxReconnects  int           `json:"max_reconnects"`
	ReconnectWait  time.Duration `json:"reconnect_wait"`
	ConnectTimeout time.Duration `json:"connect_timeout"`
}

// NewEnhancedNATSQueue creates a new enhanced NATS JetStream queue with response routing
func NewEnhancedNATSQueue(config EnhancedNATSConfig) (*EnhancedNATSQueue, error) {
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

	// Create or get stream with multiple subjects for request/response
	streamCfg := jetstream.StreamConfig{
		Name: config.StreamName,
		Subjects: []string{
			fmt.Sprintf("%s.execution.requests", config.SubjectPrefix),
			fmt.Sprintf("%s.execution.results.>", config.SubjectPrefix),
			fmt.Sprintf("%s.execution.errors.>", config.SubjectPrefix),
		},
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

	// Create or get consumer for execution requests
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

	return &EnhancedNATSQueue{
		conn:          conn,
		js:            js,
		stream:        stream,
		consumer:      consumer,
		streamName:    config.StreamName,
		subjectPrefix: config.SubjectPrefix,
		durableName:   config.DurableName,
	}, nil
}

// EnqueueWithResponse publishes a message and returns a response channel
func (q *EnhancedNATSQueue) EnqueueWithResponse(ctx context.Context, runID string, data []byte, responseTimeout time.Duration) (chan []byte, error) {
	subject := fmt.Sprintf("%s.execution.requests", q.subjectPrefix)
	responseSubject := fmt.Sprintf("%s.execution.results.%s", q.subjectPrefix, runID)

	// Create a channel for the response
	responseCh := make(chan []byte, 1)

	// Subscribe to response subject
	sub, err := q.conn.Subscribe(responseSubject, func(msg *nats.Msg) {
		responseCh <- msg.Data
		close(responseCh)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to response subject: %w", err)
	}

	// Set up auto-unsubscribe after 1 message
	sub.AutoUnsubscribe(1)

	// Create message with headers
	natsMsg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  nats.Header{},
	}
	natsMsg.Header.Set("run_id", runID)
	natsMsg.Header.Set("timestamp", time.Now().Format(time.RFC3339))
	natsMsg.Header.Set("response_subject", responseSubject)
	natsMsg.Header.Set("correlation_id", runID)

	// Publish message
	_, err = q.js.PublishMsg(ctx, natsMsg)
	if err != nil {
		sub.Unsubscribe()
		close(responseCh)
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	// Set up timeout for response
	if responseTimeout > 0 {
		go func() {
			time.Sleep(responseTimeout)
			select {
			case <-responseCh:
				// Response already received
			default:
				sub.Unsubscribe()
				close(responseCh)
			}
		}()
	}

	return responseCh, nil
}

// Enqueue adds a workflow run ID to the queue for processing
func (q *EnhancedNATSQueue) Enqueue(ctx context.Context, runID string, data []byte) error {
	subject := fmt.Sprintf("%s.execution.requests", q.subjectPrefix)

	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
		Header:  nats.Header{},
	}
	msg.Header.Set("run_id", runID)
	msg.Header.Set("timestamp", time.Now().Format(time.RFC3339))
	msg.Header.Set("correlation_id", runID)

	_, err := q.js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// PublishResult publishes execution result to response subject
func (q *EnhancedNATSQueue) PublishResult(ctx context.Context, runID string, resultData []byte) error {
	subject := fmt.Sprintf("%s.execution.results.%s", q.subjectPrefix, runID)

	msg := &nats.Msg{
		Subject: subject,
		Data:    resultData,
		Header:  nats.Header{},
	}
	msg.Header.Set("run_id", runID)
	msg.Header.Set("timestamp", time.Now().Format(time.RFC3339))
	msg.Header.Set("result_type", "success")

	_, err := q.js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to publish result: %w", err)
	}

	return nil
}

// PublishError publishes execution error to error subject
func (q *EnhancedNATSQueue) PublishError(ctx context.Context, runID string, errorData []byte) error {
	subject := fmt.Sprintf("%s.execution.errors.%s", q.subjectPrefix, runID)

	msg := &nats.Msg{
		Subject: subject,
		Data:    errorData,
		Header:  nats.Header{},
	}
	msg.Header.Set("run_id", runID)
	msg.Header.Set("timestamp", time.Now().Format(time.RFC3339))
	msg.Header.Set("result_type", "error")

	_, err := q.js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to publish error: %w", err)
	}

	return nil
}

// DequeueEnhanced retrieves the next message from the queue with enhanced info
func (q *EnhancedNATSQueue) DequeueEnhanced(ctx context.Context) (*EnhancedMessage, error) {
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
		var correlationID string
		var responseSubject string

		if headers != nil {
			timestamp, _ = time.Parse(time.RFC3339, headers.Get("timestamp"))
			runID = headers.Get("run_id")
			correlationID = headers.Get("correlation_id")
			responseSubject = headers.Get("response_subject")
		}

		// If correlation ID not set, use run ID
		if correlationID == "" {
			correlationID = runID
		}

		baseMessage := &Message{
			ID:        fmt.Sprintf("%d", metadata.Sequence.Stream),
			RunID:     runID,
			Data:      msg.Data(),
			Timestamp: timestamp,
		}

		enhancedMsg := &EnhancedMessage{
			Message:         baseMessage,
			CorrelationID:   correlationID,
			ResponseSubject: responseSubject,
			JetStreamMsg:    msg,
		}

		return enhancedMsg, nil
	}

	return nil, nil // No messages in batch
}

// Dequeue retrieves the next message from the queue (implements Queue interface)
func (q *EnhancedNATSQueue) Dequeue(ctx context.Context) (*Message, error) {
	enhancedMsg, err := q.DequeueEnhanced(ctx)
	if err != nil {
		return nil, err
	}
	if enhancedMsg == nil {
		return nil, nil
	}
	return enhancedMsg.Message, nil
}

// Acknowledge successful processing of a message
func (q *EnhancedNATSQueue) Acknowledge(ctx context.Context, messageID string) error {
	// This is a simplified implementation
	// In production, we would need to track message handles
	return nil
}

// AcknowledgeEnhanced acknowledges an enhanced message
func (q *EnhancedNATSQueue) AcknowledgeEnhanced(ctx context.Context, msg *EnhancedMessage) error {
	if msg == nil || msg.JetStreamMsg == nil {
		return fmt.Errorf("message or jetstream message is nil")
	}

	// Use DoubleAck for reliability (waits for ack from server)
	if err := msg.JetStreamMsg.DoubleAck(ctx); err != nil {
		return fmt.Errorf("failed to acknowledge message: %w", err)
	}

	return nil
}

// Reject rejects a message and schedules it for retry
func (q *EnhancedNATSQueue) Reject(ctx context.Context, messageID string, delay time.Duration) error {
	// This is a simplified implementation
	return nil
}

// RejectEnhanced rejects an enhanced message
func (q *EnhancedNATSQueue) RejectEnhanced(ctx context.Context, msg *EnhancedMessage, delay time.Duration) error {
	if msg == nil || msg.JetStreamMsg == nil {
		return fmt.Errorf("message or jetstream message is nil")
	}

	if delay > 0 {
		if err := msg.JetStreamMsg.NakWithDelay(delay); err != nil {
			return fmt.Errorf("failed to reject message with delay: %w", err)
		}
	} else {
		if err := msg.JetStreamMsg.Nak(); err != nil {
			return fmt.Errorf("failed to reject message: %w", err)
		}
	}

	return nil
}

// Size returns the current size of the queue
func (q *EnhancedNATSQueue) Size(ctx context.Context) (int, error) {
	info, err := q.stream.Info(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get stream info: %w", err)
	}
	return int(info.State.Msgs), nil
}

// IsEmpty checks if the queue is empty
func (q *EnhancedNATSQueue) IsEmpty(ctx context.Context) (bool, error) {
	size, err := q.Size(ctx)
	if err != nil {
		return false, err
	}
	return size == 0, nil
}

// Contains checks if a specific run ID is in the queue
func (q *EnhancedNATSQueue) Contains(ctx context.Context, runID string) (bool, error) {
	// This is complex with NATS JetStream as we'd need to search messages
	// For simplicity, we'll return false for now
	return false, nil
}

// Remove removes a specific run ID from the queue
func (q *EnhancedNATSQueue) Remove(ctx context.Context, runID string) error {
	// Removing specific messages is complex with NATS JetStream
	// For now, we'll return nil (no-op)
	return nil
}

// Clear removes all messages from the queue
func (q *EnhancedNATSQueue) Clear(ctx context.Context) error {
	// In NATS JetStream, we can purge the stream
	err := q.stream.Purge(ctx)
	if err != nil {
		return fmt.Errorf("failed to purge stream: %w", err)
	}
	return nil
}

// Close closes the queue connection
func (q *EnhancedNATSQueue) Close() error {
	q.conn.Close()
	return nil
}

// SubscribeToResults subscribes to results for a specific run ID
func (q *EnhancedNATSQueue) SubscribeToResults(runID string, handler func([]byte)) (*nats.Subscription, error) {
	subject := fmt.Sprintf("%s.execution.results.%s", q.subjectPrefix, runID)
	return q.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
}

// SubscribeToErrors subscribes to errors for a specific run ID
func (q *EnhancedNATSQueue) SubscribeToErrors(runID string, handler func([]byte)) (*nats.Subscription, error) {
	subject := fmt.Sprintf("%s.execution.errors.%s", q.subjectPrefix, runID)
	return q.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
}
