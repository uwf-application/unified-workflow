# pkg/client/

Client libraries for interacting with the workflow system from external applications.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `client.go` | Top-level unified client entry point | Using the Go client, understanding the client API |
| `http_client.go` | Shared HTTP client utilities for all client libraries | Modifying shared HTTP behavior across clients |

## Subdirectories

| Directory | What | When to read |
| --------- | ---- | ------------ |
| `go/` | Go SDK client library | Using the Go SDK, modifying Go client behavior |
| `typescript/` | TypeScript/Node.js SDK client library | Using the TypeScript SDK, modifying TS client behavior |
| `python/` | Python SDK client library | Using the Python SDK, modifying Python client |
| `executor/` | Executor service client | Calling executor operations from external code |
| `queue/` | Queue client for publishing workflow events | Publishing messages to the workflow queue |
| `registry/` | Registry client for workflow registration/lookup | Registering or looking up workflows externally |
| `state/` | State client for reading/writing workflow state | Reading or writing workflow execution state |
