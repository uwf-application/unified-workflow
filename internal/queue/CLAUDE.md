# internal/queue/

Message queue interface and implementations (NATS JetStream, in-memory).

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `queue.go` | Queue interface: Publish, Subscribe, Ack contracts | Understanding queue contract, adding queue backends |
| `nats_queue.go` | NATS JetStream queue implementation | Modifying NATS consumer/producer behavior, debugging NATS issues |
| `nats_queue_enhanced.go` | Enhanced NATS queue with retry, backoff, and observability | Modifying retry behavior, tuning NATS performance |
| `in_memory_queue.go` | In-memory queue implementation for testing | Modifying test queue behavior |
