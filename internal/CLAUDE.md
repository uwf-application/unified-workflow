# internal/

Private application packages — not importable by external modules.

## Subdirectories

| Directory | What | When to read |
| --------- | ---- | ------------ |
| `api/` | HTTP handlers and middleware for the workflow API | Adding endpoints, modifying request handling |
| `common/` | Shared models used across packages | Adding shared data types |
| `config/` | Configuration loading and validation | Adding config fields, changing defaults |
| `di/` | Dependency injection container and wiring | Adding services, modifying DI setup |
| `executor/` | Workflow and step execution engines | Modifying execution logic, adding executor types |
| `primitive/` | External service integrations (antifraud, etc.) | Adding primitive services, modifying service clients |
| `queue/` | Message queue implementations (NATS, in-memory) | Modifying queue behavior, adding queue backends |
| `registry/` | Workflow registry implementations (HTTP, in-memory) | Modifying workflow registration/lookup |
| `serviceclients/` | HTTP clients for external services | Adding external service integrations |
| `state/` | Workflow state storage implementations (Redis, in-memory) | Modifying state persistence, adding state backends |
