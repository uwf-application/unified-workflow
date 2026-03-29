# pkg/client/go/sdk/

Go SDK for the unified workflow system — public API for Go consumers.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `client.go` | SDK client: workflow execution, registration, status polling | Integrating the Go SDK, adding SDK methods |
| `config.go` | SDK configuration (server URL, timeouts, auth) | Configuring the Go SDK |
| `models.go` | SDK-facing data types (WorkflowRequest, ExecutionResult, etc.) | Understanding SDK data shapes, adding SDK types |
| `errors.go` | SDK error types and error handling helpers | Handling SDK errors, adding error types |
| `parser.go` | Response parser for SDK API responses | Modifying response deserialization |
| `validator.go` | Input validation for SDK requests | Adding input validation rules |
| `example.go` | Usage examples for the Go SDK | Understanding SDK usage patterns |
| `README.md` | Go SDK documentation | Onboarding to the Go SDK |
| `config.example.yaml` | Example configuration file for the Go SDK | Setting up SDK configuration |
