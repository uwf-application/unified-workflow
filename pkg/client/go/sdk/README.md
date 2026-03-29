# Workflow SDK

A comprehensive Go SDK for interacting with the Unified Workflow Execution Platform. This SDK provides a clean, type-safe interface for executing workflows, managing executions, and handling HTTP request parsing with validation.

## Features

- **HTTP Request Parsing**: Automatically parse HTTP requests into workflow execution requests
- **Data Validation**: Built-in validation with customizable rules
- **Request Sanitization**: Clean and sanitize input data
- **Session Extraction**: Extract user sessions from HTTP requests
- **Security Context**: Pass security information to workflows
- **Error Handling**: Comprehensive error handling with retry mechanisms
- **Async Execution**: Support for asynchronous workflow execution
- **Execution Management**: Monitor and manage workflow executions

## Installation

```bash
go get unified-workflow/pkg/client/go/sdk
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "unified-workflow/pkg/client/go/sdk"
)

func main() {
    // Create SDK configuration
    config := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://localhost:8080",
        Timeout:             30 * time.Second,
        MaxRetries:          3,
        AuthToken:          "your-auth-token",
        EnableValidation:   true,
        EnableSanitization: true,
    }

    // Create SDK client
    client, err := sdk.NewClient(config)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Execute a workflow
    ctx := context.Background()
    workflowID := "payment-processing-workflow"
    inputData := map[string]interface{}{
        "user_id": "user_12345",
        "amount":  99.99,
        "email":   "user@example.com",
    }

    resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Workflow execution started: %s\n", resp.RunID)
    fmt.Printf("Status URL: %s\n", resp.StatusURL)
}
```

### HTTP Request Integration

```go
func handlePayment(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    config := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://localhost:8080",
        Timeout:             30 * time.Second,
        AuthToken:          "your-auth-token",
    }

    client, err := sdk.NewClient(config)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer client.Close()

    // Execute workflow from HTTP request
    resp, err := client.ExecuteFromHTTPRequest(ctx, "payment-processing-workflow", r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Return execution information
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(resp)
}
```

## Configuration

The SDK supports multiple configuration sources in order of precedence:
1. **Programmatic configuration** - Directly in code
2. **Environment variables** - For runtime configuration
3. **Configuration files** - YAML or JSON files
4. **Default values** - Built-in sensible defaults

### Configuration Loading Methods

```go
// 1. Programmatic configuration
config := &sdk.SDKConfig{
    WorkflowAPIEndpoint: "http://localhost:8080",
    Timeout:             30 * time.Second,
    // ... other settings
}

// 2. Load from environment variables
config, err := sdk.LoadConfigFromEnvironment()

// 3. Load from configuration file
config, err := sdk.LoadConfigFromFile("sdk-config.yaml")

// 4. Load from default locations (env vars → config files → defaults)
config, err := sdk.LoadDefaultConfig()
```

### SDKConfig Fields

| Field | Type | Description | Default | Environment Variable |
|-------|------|-------------|---------|----------------------|
| `WorkflowAPIEndpoint` | `string` | Workflow API endpoint URL | Required | `SDK_WORKFLOW_API_ENDPOINT` |
| `Timeout` | `time.Duration` | Request timeout | `30s` | `SDK_TIMEOUT` |
| `MaxRetries` | `int` | Maximum retry attempts | `3` | `SDK_MAX_RETRIES` |
| `RetryDelay` | `time.Duration` | Delay between retries | `1s` | `SDK_RETRY_DELAY` |
| `AuthToken` | `string` | Authentication token | `""` | `SDK_AUTH_TOKEN` |
| `AuthType` | `AuthType` | Authentication type | `bearer_token` | `SDK_AUTH_TYPE` |
| `EnableValidation` | `bool` | Enable input validation | `true` | `SDK_ENABLE_VALIDATION` |
| `EnableSanitization` | `bool` | Enable data sanitization | `true` | `SDK_ENABLE_SANITIZATION` |
| `StrictValidation` | `bool` | Fail on validation errors | `false` | `SDK_STRICT_VALIDATION` |
| `EnableSessionExtraction` | `bool` | Extract session from HTTP requests | `true` | `SDK_ENABLE_SESSION_EXTRACTION` |
| `EnableSecurityContext` | `bool` | Extract security context | `true` | `SDK_ENABLE_SECURITY_CONTEXT` |
| `IncludeFullHTTPContext` | `bool` | Include full HTTP context | `true` | `SDK_INCLUDE_FULL_HTTP_CONTEXT` |
| `LogLevel` | `LogLevel` | Logging level | `info` | `SDK_LOG_LEVEL` |
| `EnableRequestLogging` | `bool` | Enable request logging | `true` | `SDK_ENABLE_REQUEST_LOGGING` |
| `EnableMetrics` | `bool` | Enable metrics collection | `true` | `SDK_ENABLE_METRICS` |
| `DefaultValidationRules` | `[]ValidationRule` | Default validation rules | `[]` | N/A |
| `CustomValidators` | `[]string` | Custom validator paths | `[]` | N/A |
| `AsyncExecution` | `bool` | Execute workflows asynchronously | `true` | `SDK_ASYNC_EXECUTION` |
| `DefaultPriority` | `int` | Default execution priority | `5` | `SDK_DEFAULT_PRIORITY` |
| `PollIntervalMs` | `int` | Poll interval for async results | `1000` | `SDK_POLL_INTERVAL_MS` |
| `EstimatedCompletionMs` | `int` | Estimated completion time | `5000` | `SDK_ESTIMATED_COMPLETION_MS` |
| `ExecutionExpiryDuration` | `time.Duration` | Execution expiry duration | `1h` | `SDK_EXECUTION_EXPIRY_DURATION` |
| `SDKVersion` | `string` | SDK version | `1.0.0` | `SDK_VERSION` |
| `EnableCircuitBreaker` | `bool` | Enable circuit breaker | `true` | `SDK_ENABLE_CIRCUIT_BREAKER` |
| `CircuitBreakerThreshold` | `int` | Circuit breaker failure threshold | `5` | `SDK_CIRCUIT_BREAKER_THRESHOLD` |
| `CircuitBreakerTimeout` | `time.Duration` | Circuit breaker timeout | `60s` | `SDK_CIRCUIT_BREAKER_TIMEOUT` |

### Configuration File Example (YAML)

```yaml
# sdk-config.yaml
workflow_api_endpoint: "http://localhost:8080"
timeout: 30s
max_retries: 3
retry_delay: 1s
auth_token: "your-auth-token-here"
auth_type: "bearer_token"
enable_validation: true
enable_sanitization: true
strict_validation: false
enable_session_extraction: true
enable_security_context: true
include_full_http_context: true
log_level: "info"
enable_request_logging: true
enable_metrics: true
default_validation_rules:
  - field: "user_id"
    rule_type: "required"
    required: true
  - field: "email"
    rule_type: "email"
    required: true
  - field: "amount"
    rule_type: "number"
    min_value: 0.01
    max_value: 1000000
async_execution: true
default_priority: 5
poll_interval_ms: 1000
estimated_completion_ms: 5000
execution_expiry_duration: 1h
sdk_version: "1.0.0"
enable_circuit_breaker: true
circuit_breaker_threshold: 5
circuit_breaker_timeout: 60s
```

### Configuration File Example (JSON)

```json
{
  "workflow_api_endpoint": "http://localhost:8080",
  "timeout": "30s",
  "max_retries": 3,
  "retry_delay": "1s",
  "auth_token": "your-auth-token-here",
  "auth_type": "bearer_token",
  "enable_validation": true,
  "enable_sanitization": true,
  "strict_validation": false,
  "enable_session_extraction": true,
  "enable_security_context": true,
  "include_full_http_context": true,
  "log_level": "info",
  "enable_request_logging": true,
  "enable_metrics": true,
  "default_validation_rules": [
    {
      "field": "user_id",
      "rule_type": "required",
      "required": true
    },
    {
      "field": "email",
      "rule_type": "email",
      "required": true
    },
    {
      "field": "amount",
      "rule_type": "number",
      "min_value": 0.01,
      "max_value": 1000000
    }
  ],
  "async_execution": true,
  "default_priority": 5,
  "poll_interval_ms": 1000,
  "estimated_completion_ms": 5000,
  "execution_expiry_duration": "1h",
  "sdk_version": "1.0.0",
  "enable_circuit_breaker": true,
  "circuit_breaker_threshold": 5,
  "circuit_breaker_timeout": "60s"
}
```

### Validation Rules

```go
rules := []sdk.ValidationRule{
    {
        Field:    "user_id",
        Required: true,
        RuleType: sdk.ValidationRuleTypeString,
        MinLength: &[]int{5}[0],
        MaxLength: &[]int{50}[0],
    },
    {
        Field:    "amount",
        Required: true,
        RuleType: sdk.ValidationRuleTypeNumber,
        MinValue: &[]float64{0.01}[0],
        MaxValue: &[]float64{10000}[0],
    },
    {
        Field:    "email",
        Required: false,
        RuleType: sdk.ValidationRuleTypeEmail,
    },
    {
        Field:         "status",
        Required:      true,
        RuleType:      sdk.ValidationRuleTypeString,
        AllowedValues: []string{"pending", "approved", "rejected"},
    },
}
```

## API Reference

### WorkflowSDKClient Interface

```go
type WorkflowSDKClient interface {
    // Execute workflow from HTTP request
    ExecuteFromHTTPRequest(ctx context.Context, workflowID string, req *http.Request) (*SDKExecuteWorkflowResponse, error)
    
    // Execute workflow with raw data
    ExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}) (*SDKExecuteWorkflowResponse, error)
    
    // Execute workflow with full context
    ExecuteWorkflowWithContext(ctx context.Context, workflowID string, sdkReq *SDKExecuteWorkflowRequest) (*SDKExecuteWorkflowResponse, error)
    
    // Validate and execute workflow
    ValidateAndExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}, rules []ValidationRule) (*SDKExecuteWorkflowResponse, error)
    
    // Get execution status
    GetExecutionStatus(ctx context.Context, runID string) (*executor.GetExecutionStatusResponse, error)
    
    // Get execution details
    GetExecutionDetails(ctx context.Context, runID string) (*executor.GetExecutionDetailsResponse, error)
    
    // Cancel execution
    CancelExecution(ctx context.Context, runID string) error
    
    // Health check
    Ping(ctx context.Context) error
    
    // Close the client
    Close() error
}
```

### SDKExecuteWorkflowRequest

```go
type SDKExecuteWorkflowRequest struct {
    InputData          map[string]interface{}
    HTTPRequest        *HTTPRequestContext
    Session            *SessionContext
    Security           *SecurityContext
    EnableValidation   bool
    EnableSanitization bool
    ValidationRules    []ValidationRule
    IncludeFullContext bool
    TimeoutMs          int64
    Metadata           map[string]interface{}
}
```

### SDKExecuteWorkflowResponse

```go
type SDKExecuteWorkflowResponse struct {
    RunID                 string
    Status                string
    Message               string
    StatusURL             string
    ResultURL             string
    PollAfterMs           int64
    EstimatedCompletionMs int64
    ExpiresAt             time.Time
    ContextIncluded       bool
    SDKVersion            string
    RequestID             string
    ValidationResult      *ValidationResult
}
```

## Advanced Usage

### Custom Validation

```go
// Define custom validation rules
customRules := []sdk.ValidationRule{
    {
        Field:    "transaction_id",
        Required: true,
        RuleType: sdk.ValidationRuleTypeUUID,
    },
    {
        Field:    "callback_url",
        Required: false,
        RuleType: sdk.ValidationRuleTypeURL,
    },
    {
        Field:    "metadata.tags",
        Required: false,
        RuleType: sdk.ValidationRuleTypeArray,
    },
}

// Validate and execute
resp, err := client.ValidateAndExecuteWorkflow(ctx, workflowID, inputData, customRules)
if err != nil {
    // Handle validation errors
    if sdkErr, ok := err.(*sdk.SDKError); ok && sdkErr.Code == sdk.ErrCodeValidationFailed {
        fmt.Println("Validation failed:", sdkErr.Details)
    }
}
```

### Full Context Execution

```go
// Create SDK execution request with full context
sdkReq := sdk.NewSDKExecuteWorkflowRequest(inputData)
sdkReq.HTTPRequest = &sdk.HTTPRequestContext{
    Method:    r.Method,
    Path:      r.URL.Path,
    Headers:   headers,
    Timestamp: time.Now(),
}
sdkReq.Session = &sdk.SessionContext{
    UserID:    "user_12345",
    SessionID: "session_67890",
    AuthMethod: "jwt",
}
sdkReq.Security = &sdk.SecurityContext{
    IPAddress:  r.RemoteAddr,
    UserAgent:  r.UserAgent(),
    IsSecure:   r.TLS != nil,
}
sdkReq.IncludeFullContext = true

// Execute with full context
resp, err := client.ExecuteWorkflowWithContext(ctx, workflowID, sdkReq)
```

### Error Handling

```go
resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
if err != nil {
    switch err := err.(type) {
    case *sdk.SDKError:
        switch err.Code {
        case sdk.ErrCodeValidationFailed:
            // Handle validation errors
            fmt.Println("Validation failed:", err.Details)
        case sdk.ErrCodeRequestParsingFailed:
            // Handle parsing errors
            fmt.Println("Failed to parse request:", err.OriginalError)
        case sdk.ErrCodeWorkflowExecution:
            // Handle execution errors
            fmt.Println("Workflow execution failed:", err.Message)
        default:
            // Handle other SDK errors
            fmt.Println("SDK error:", err)
        }
    default:
        // Handle non-SDK errors
        fmt.Println("Unexpected error:", err)
    }
    return
}
```

## Examples

See the [example.go](./example.go) file for complete usage examples including:

1. Basic workflow execution
2. HTTP request integration
3. Validation and error handling
4. Execution status monitoring
5. HTTP handler implementation

## Best Practices

1. **Always close the client**: Use `defer client.Close()` to ensure proper cleanup
2. **Use context for cancellation**: Pass context with timeout for long-running operations
3. **Enable validation**: Always enable validation for production use
4. **Handle errors properly**: Check for specific error types and handle accordingly
5. **Monitor execution status**: Use the provided status URLs to track execution progress
6. **Use async execution**: For long-running workflows, use async execution and poll for results

## Testing

```go
func TestWorkflowExecution(t *testing.T) {
    config := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://localhost:8080",
        Timeout:             5 * time.Second,
        EnableValidation:    true,
    }

    client, err := sdk.NewClient(config)
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()
    inputData := map[string]interface{}{
        "test": true,
        "data": "test-data",
    }

    resp, err := client.ExecuteWorkflow(ctx, "test-workflow", inputData)
    require.NoError(t, err)
    assert.NotEmpty(t, resp.RunID)
    assert.NotEmpty(t, resp.StatusURL)
}
```

## License

This SDK is part of the Unified Workflow Execution Platform. See the main project LICENSE for details.

## Support

For issues, questions, or contributions, please refer to the main project repository.