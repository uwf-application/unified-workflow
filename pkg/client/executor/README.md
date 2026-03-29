# Executor Client

A Go client for interacting with the Unified Workflow Executor service. This client provides a comprehensive interface for executing workflows, managing executions, and monitoring execution status.

## Features

- **Workflow Execution**: Execute workflows synchronously or asynchronously
- **Execution Management**: Monitor, pause, resume, cancel, and retry executions
- **Status Tracking**: Get detailed execution status and progress
- **Execution Details**: Retrieve comprehensive execution information including steps and child steps
- **Metrics Collection**: Access execution metrics for monitoring and analysis
- **Batch Operations**: List and filter executions
- **Error Handling**: Comprehensive error handling with retry mechanisms

## Installation

```bash
go get unified-workflow/pkg/client/executor
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "unified-workflow/pkg/client/executor"
)

func main() {
    // Create executor client configuration
    config := executor.DefaultConfig()
    config.Endpoint = "http://localhost:8080"
    config.Timeout = 30 * time.Second
    config.AuthToken = "your-auth-token"

    // Create executor client
    client := executor.NewClient(config)
    defer client.Close()

    // Execute a workflow
    ctx := context.Background()
    req := &executor.ExecuteWorkflowRequest{
        Request: client.Request{
            ID:        "req-123",
            Timestamp: time.Now(),
        },
        WorkflowID: "payment-processing-workflow",
        InputData: map[string]interface{}{
            "user_id": "user_12345",
            "amount":  99.99,
            "email":   "user@example.com",
        },
        Async:    true,
        Priority: 5,
    }

    resp, err := client.ExecuteWorkflow(ctx, req)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Workflow execution started: %s\n", resp.RunID)
    fmt.Printf("Status URL: %s\n", resp.StatusURL)
    fmt.Printf("Estimated completion: %d ms\n", resp.EstimatedCompletionMs)
}
```

## API Reference

### Core Methods

#### `ExecuteWorkflow(ctx context.Context, req *ExecuteWorkflowRequest) (*ExecuteWorkflowResponse, error)`
Execute a workflow with the given input data.

#### `GetExecutionStatus(ctx context.Context, req *GetExecutionStatusRequest) (*GetExecutionStatusResponse, error)`
Get the current status of a workflow execution.

#### `GetExecutionDetails(ctx context.Context, req *GetExecutionDetailsRequest) (*GetExecutionDetailsResponse, error)`
Get detailed information about a workflow execution.

#### `CancelExecution(ctx context.Context, req *CancelExecutionRequest) (*CancelExecutionResponse, error)`
Cancel a running workflow execution.

#### `PauseExecution(ctx context.Context, req *PauseExecutionRequest) (*PauseExecutionResponse, error)`
Pause a running workflow execution.

#### `ResumeExecution(ctx context.Context, req *ResumeExecutionRequest) (*ResumeExecutionResponse, error)`
Resume a paused workflow execution.

#### `RetryExecution(ctx context.Context, req *RetryExecutionRequest) (*RetryExecutionResponse, error)`
Retry a failed workflow execution.

#### `GetExecutionData(ctx context.Context, req *GetExecutionDataRequest) (*GetExecutionDataResponse, error)`
Get the data associated with a workflow execution.

#### `ListExecutions(ctx context.Context, req *ListExecutionsRequest) (*ListExecutionsResponse, error)`
List workflow executions with optional filters.

#### `GetExecutionMetrics(ctx context.Context, req *GetExecutionMetricsRequest) (*GetExecutionMetricsResponse, error)`
Get metrics for a workflow execution.

#### `GetStepExecution(ctx context.Context, req *GetStepExecutionRequest) (*GetStepExecutionResponse, error)`
Get details of a specific step execution.

#### `GetChildStepExecution(ctx context.Context, req *GetChildStepExecutionRequest) (*GetChildStepExecutionResponse, error)`
Get details of a specific child step execution.

### Request Types

#### `ExecuteWorkflowRequest`
```go
type ExecuteWorkflowRequest struct {
    client.Request
    WorkflowID  string                 // Workflow ID to execute
    InputData   map[string]interface{} // Input data for the workflow
    Async       bool                   // Execute asynchronously
    CallbackURL string                 // Callback URL for async execution
    TimeoutMs   int64                  // Execution timeout in milliseconds
    Priority    int                    // Execution priority (1-10)
    Metadata    map[string]interface{} // Additional metadata
}
```

#### `GetExecutionStatusRequest`
```go
type GetExecutionStatusRequest struct {
    client.Request
    RunID         string // Execution run ID
    IncludeDetails bool  // Include detailed status
    WaitMs        int64  // Maximum time to wait for completion
}
```

#### `GetExecutionDetailsRequest`
```go
type GetExecutionDetailsRequest struct {
    client.Request
    RunID              string // Execution run ID
    IncludeSteps       bool   // Include step details
    IncludeChildSteps  bool   // Include child step details
    IncludePrimitives  bool   // Include primitive details
}
```

## Configuration

### Executor Client Configuration

```go
type Config struct {
    client.Config
    
    // DefaultTimeoutMs is the default execution timeout in milliseconds
    DefaultTimeoutMs int64
    
    // DefaultPriority is the default execution priority
    DefaultPriority int
    
    // EnableAsync indicates if async execution is enabled
    EnableAsync bool
}
```

### Default Configuration

```go
func DefaultConfig() Config {
    return Config{
        Config:           client.DefaultConfig(),
        DefaultTimeoutMs: 30000, // 30 seconds
        DefaultPriority:  5,
        EnableAsync:      true,
    }
}
```

## Examples

### Synchronous Execution

```go
func executeWorkflowSync(client executor.Client) error {
    ctx := context.Background()
    req := &executor.ExecuteWorkflowRequest{
        WorkflowID: "sync-workflow",
        InputData: map[string]interface{}{
            "data": "test-data",
        },
        Async: false, // Synchronous execution
    }

    resp, err := client.ExecuteWorkflow(ctx, req)
    if err != nil {
        return err
    }

    fmt.Printf("Execution completed: %s\n", resp.RunID)
    fmt.Printf("Status: %s\n", resp.Status)
    return nil
}
```

### Asynchronous Execution with Status Polling

```go
func executeWorkflowAsync(client executor.Client) error {
    ctx := context.Background()
    
    // Start async execution
    req := &executor.ExecuteWorkflowRequest{
        WorkflowID: "async-workflow",
        InputData: map[string]interface{}{
            "data": "test-data",
        },
        Async:    true,
        Priority: 5,
    }

    resp, err := client.ExecuteWorkflow(ctx, req)
    if err != nil {
        return err
    }

    fmt.Printf("Async execution started: %s\n", resp.RunID)
    
    // Poll for completion
    for {
        statusReq := &executor.GetExecutionStatusRequest{
            RunID: resp.RunID,
        }
        
        statusResp, err := client.GetExecutionStatus(ctx, statusReq)
        if err != nil {
            return err
        }

        if statusResp.Status.IsTerminal {
            if statusResp.Status.Status == "completed" {
                fmt.Println("Execution completed successfully")
            } else {
                fmt.Printf("Execution failed: %s\n", statusResp.Status.ErrorMessage)
            }
            break
        }

        fmt.Printf("Progress: %.2f%%\n", statusResp.Status.Progress*100)
        time.Sleep(1 * time.Second)
    }
    
    return nil
}
```

### Execution Management

```go
func manageExecution(client executor.Client, runID string) error {
    ctx := context.Background()
    
    // Get execution status
    statusReq := &executor.GetExecutionStatusRequest{
        RunID: runID,
    }
    statusResp, err := client.GetExecutionStatus(ctx, statusReq)
    if err != nil {
        return err
    }
    
    fmt.Printf("Current status: %s\n", statusResp.Status.Status)
    
    // Pause execution if running
    if statusResp.Status.Status == "running" {
        pauseReq := &executor.PauseExecutionRequest{
            RunID:  runID,
            Reason: "Manual pause",
        }
        pauseResp, err := client.PauseExecution(ctx, pauseReq)
        if err != nil {
            return err
        }
        fmt.Printf("Execution paused: %v\n", pauseResp.Paused)
    }
    
    // Get execution details
    detailsReq := &executor.GetExecutionDetailsRequest{
        RunID:             runID,
        IncludeSteps:      true,
        IncludeChildSteps: true,
    }
    detailsResp, err := client.GetExecutionDetails(ctx, detailsReq)
    if err != nil {
        return err
    }
    
    fmt.Printf("Total steps: %d\n", len(detailsResp.Details.Steps))
    fmt.Printf("Total duration: %d ms\n", detailsResp.Details.TotalDurationMs)
    
    return nil
}
```

### Listing Executions

```go
func listExecutions(client executor.Client) error {
    ctx := context.Background()
    
    // List recent executions
    listReq := &executor.ListExecutionsRequest{
        Limit:    10,
        SortBy:   "created_at",
        SortOrder: "desc",
    }
    
    listResp, err := client.ListExecutions(ctx, listReq)
    if err != nil {
        return err
    }
    
    fmt.Printf("Total executions: %d\n", listResp.TotalCount)
    fmt.Printf("Showing: %d executions\n", len(listResp.Executions))
    
    for _, exec := range listResp.Executions {
        fmt.Printf("Run ID: %s, Status: %s, Created: %v\n",
            exec.RunID, exec.Status, exec.CreatedAt.Format(time.RFC3339))
    }
    
    return nil
}
```

### Error Handling

```go
func executeWithRetry(client executor.Client) error {
    ctx := context.Background()
    maxRetries := 3
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        req := &executor.ExecuteWorkflowRequest{
            WorkflowID: "retry-workflow",
            InputData: map[string]interface{}{
                "attempt": attempt,
            },
        }
        
        resp, err := client.ExecuteWorkflow(ctx, req)
        if err != nil {
            fmt.Printf("Attempt %d failed: %v\n", attempt, err)
            
            // Check if error is retryable
            if clientErr, ok := err.(*client.Error); ok && clientErr.Retryable {
                if attempt < maxRetries {
                    fmt.Printf("Retrying in 1 second...\n")
                    time.Sleep(1 * time.Second)
                    continue
                }
            }
            return err
        }
        
        fmt.Printf("Execution successful on attempt %d: %s\n", attempt, resp.RunID)
        return nil
    }
    
    return fmt.Errorf("max retries exceeded")
}
```

## Best Practices

### 1. Use Context for Cancellation
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

req := &executor.ExecuteWorkflowRequest{
    WorkflowID: "workflow-id",
    InputData:  data,
}

resp, err := client.ExecuteWorkflow(ctx, req)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("Request timed out")
    }
    return err
}
```

### 2. Configure Appropriate Timeouts
```go
config := executor.DefaultConfig()
config.Timeout = 60 * time.Second // Longer timeout for complex workflows
config.DefaultTimeoutMs = 120000   // 2 minutes for execution timeout
```

### 3. Handle Async Execution Properly
```go
// For long-running workflows, use async execution
req := &executor.ExecuteWorkflowRequest{
    WorkflowID: "long-running-workflow",
    InputData:  data,
    Async:      true,
    CallbackURL: "https://your-server.com/callbacks",
}

resp, err := client.ExecuteWorkflow(ctx, req)
if err != nil {
    return err
}

// Store run ID for later status checking
storeRunID(resp.RunID)
```

### 4. Monitor Execution Progress
```go
func monitorExecution(client executor.Client, runID string) error {
    for {
        statusReq := &executor.GetExecutionStatusRequest{
            RunID: runID,
        }
        
        statusResp, err := client.GetExecutionStatus(ctx, statusReq)
        if err != nil {
            return err
        }
        
        if statusResp.Status.IsTerminal {
            return handleTerminalState(statusResp.Status)
        }
        
        // Update progress
        updateProgress(statusResp.Status.Progress)
        
        // Wait before polling again
        time.Sleep(2 * time.Second)
    }
}
```

### 5. Use Execution Details for Debugging
```go
func debugExecution(client executor.Client, runID string) error {
    detailsReq := &executor.GetExecutionDetailsRequest{
        RunID:             runID,
        IncludeSteps:      true,
        IncludeChildSteps: true,
        IncludePrimitives: true,
    }
    
    detailsResp, err := client.GetExecutionDetails(ctx, detailsReq)
    if err != nil {
        return err
    }
    
    // Analyze execution details
    for _, step := range detailsResp.Details.Steps {
        if step.Status == "failed" {
            fmt.Printf("Step %d failed: %s\n", step.StepIndex, step.ErrorMessage)
            for _, childStep := range step.ChildSteps {
                if childStep.Status == "failed" {
                    fmt.Printf("  Child step %d failed: %s\n", 
                        childStep.ChildStepIndex, childStep.ErrorMessage)
                }
            }
        }
    }
    
    return nil
}
```

## Testing

```go
func TestExecutorClient(t *testing.T) {
    // Create test configuration
    config := executor.DefaultConfig()
    config.Endpoint = "http://localhost:8080"
    
    // Create client
    client := executor.NewClient(config)
    defer client.Close()
    
    // Test ping
    ctx := context.Background()
    err := client.Ping(ctx)
    require.NoError(t, err)
    
    // Test workflow execution
    req := &executor.ExecuteWorkflowRequest{
        WorkflowID: "test-workflow",
        InputData: map[string]interface{}{
            "test": true,
        },
        Async: true,
    }
    
    resp, err := client.ExecuteWorkflow(ctx, req)
    require.NoError(t, err)
    assert.NotEmpty(t, resp.RunID)
    assert.NotEmpty(t, resp.StatusURL)
}
```

## Integration with Workflow SDK

The executor client can be used directly or as part of the higher-level Workflow SDK:

```go
import (
    "unified-workflow/pkg/client/executor"
    "unified-workflow/pkg/client/go/sdk"
)

func main() {
    // Create executor client
    execConfig := executor.DefaultConfig()
    execConfig.Endpoint = "http://localhost:8080"
    executorClient := executor.NewClient(execConfig)
    
    // Create SDK client with executor client
    sdkConfig := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://localhost:8080",
    }
    
    sdkClient, err := sdk.NewClientWithExecutor(sdkConfig, executorClient)
    if err != nil {
        panic(err)
    }
    defer sdkClient.Close()
    
    // Use SDK for higher-level operations
    // ...
}
```

## Error Codes

Common error codes returned by the executor client:

- `EXECUTION_NOT_FOUND`: Execution with the given run ID was not found
- `WORKFLOW_NOT_FOUND`: Workflow with the given ID was not found
- `EXECUTION_ALREADY_COMPLETED`: Cannot modify a completed execution
- `EXECUTION_NOT_RUNNING`: Cannot pause/resume a non-running execution
- `INVALID_PRIORITY`: Invalid execution priority value
- `TIMEOUT`: Execution timeout exceeded
- `CIRCUIT_BREAKER_OPEN`: Circuit breaker is open, requests are blocked

## Support

For issues, questions, or contributions, please refer to the main project repository.

## License

This client is part of the Unified Workflow Execution Platform. See the main project LICENSE for details.