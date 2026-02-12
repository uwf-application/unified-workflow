# Unified Workflow SDK - Client Guide

## Overview
The Unified Workflow SDK provides a simple, type-safe interface for interacting with the Unified Workflow Execution Platform from your Go applications. This guide will help you get started with the SDK.

## Quick Start

### Installation

#### Option 1: Go Module (Recommended)
```bash
go get github.com/your-org/unified-workflow-sdk
```

#### Option 2: Direct Import
```go
import "unified-workflow/pkg/client/sdk"
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/your-org/unified-workflow-sdk"
)

func main() {
    // 1. Configure the SDK
    config := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://localhost:8081", // Your workflow API endpoint
        Timeout:             30 * time.Second,
        MaxRetries:          3,
        AuthToken:          "your-api-token",         // Optional authentication
        EnableValidation:    true,                     // Enable input validation
        EnableSanitization:  true,                     // Enable data sanitization
    }
    
    // 2. Create the SDK client
    client, err := sdk.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create SDK client: %v", err)
    }
    defer client.Close()
    
    // 3. Execute a workflow
    ctx := context.Background()
    workflowID := "payment-processing-workflow"
    
    inputData := map[string]interface{}{
        "user_id": "user_12345",
        "amount":  99.99,
        "email":   "user@example.com",
        "metadata": map[string]interface{}{
            "source": "web",
            "device": "mobile",
        },
    }
    
    // 4. Execute the workflow
    resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
    if err != nil {
        log.Fatalf("Failed to execute workflow: %v", err)
    }
    
    // 5. Handle the response
    fmt.Printf("✅ Workflow execution started!\n")
    fmt.Printf("   Run ID: %s\n", resp.RunID)
    fmt.Printf("   Status: %s\n", resp.Status)
    fmt.Printf("   Status URL: %s\n", resp.StatusURL)
    fmt.Printf("   Result URL: %s\n", resp.ResultURL)
    fmt.Printf("   Poll after: %d ms\n", resp.PollAfterMs)
    
    // 6. Check execution status (optional)
    time.Sleep(2 * time.Second)
    statusResp, err := client.GetExecutionStatus(ctx, resp.RunID)
    if err != nil {
        log.Printf("Failed to get execution status: %v", err)
    } else {
        fmt.Printf("✅ Execution Status: %s (Progress: %.2f)\n", 
            statusResp.Status.Status, statusResp.Status.Progress)
    }
}
```

## Core Concepts

### 1. SDK Configuration
The `SDKConfig` struct controls SDK behavior:

```go
type SDKConfig struct {
    WorkflowAPIEndpoint string        // Required: API endpoint URL
    Timeout             time.Duration // Request timeout (default: 30s)
    MaxRetries          int           // Max retry attempts (default: 3)
    RetryDelay          time.Duration // Delay between retries (default: 1s)
    AuthToken           string        // Authentication token
    EnableValidation    bool          // Enable input validation (default: true)
    EnableSanitization  bool          // Enable data sanitization (default: true)
    StrictValidation    bool          // Fail on validation errors (default: false)
    DefaultValidationRules []ValidationRule // Default validation rules
}
```

### 2. Workflow Execution
The SDK supports multiple ways to execute workflows:

#### Execute with Raw Data
```go
inputData := map[string]interface{}{
    "user_id": "user_123",
    "amount":  50.0,
    "items":   []string{"item1", "item2"},
}

resp, err := client.ExecuteWorkflow(ctx, "order-workflow", inputData)
```

#### Execute from HTTP Request
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    resp, err := client.ExecuteFromHTTPRequest(ctx, "payment-workflow", r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Handle response
}
```

#### Validate and Execute
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
}

resp, err := client.ValidateAndExecuteWorkflow(ctx, "payment-workflow", inputData, rules)
```

### 3. Execution Management
Monitor and manage workflow executions:

```go
// Get execution status
statusResp, err := client.GetExecutionStatus(ctx, runID)
if err != nil {
    // Handle error
}
fmt.Printf("Status: %s, Progress: %.2f\n", statusResp.Status.Status, statusResp.Status.Progress)

// Get execution details
detailsResp, err := client.GetExecutionDetails(ctx, runID)
if err != nil {
    // Handle error
}

// Cancel execution
err = client.CancelExecution(ctx, runID)
if err != nil {
    // Handle error
}
```

## Advanced Features

### 1. Data Validation
The SDK includes a powerful validation system:

```go
// Define validation rules
rules := []sdk.ValidationRule{
    {
        Field:     "email",
        Required:  true,
        RuleType:  sdk.ValidationRuleTypeEmail,
    },
    {
        Field:    "phone",
        Required: false,
        RuleType: sdk.ValidationRuleTypeString,
        Pattern:  `^\+?[1-9]\d{1,14}$`, // E.164 format
    },
    {
        Field:         "status",
        Required:      true,
        RuleType:      sdk.ValidationRuleTypeString,
        AllowedValues: []string{"pending", "approved", "rejected"},
    },
}

// Validate and execute
resp, err := client.ValidateAndExecuteWorkflow(ctx, workflowID, inputData, rules)
```

### 2. Context Extraction
Extract context from HTTP requests:

```go
// The SDK automatically extracts:
// - HTTP headers and method
// - Query parameters
// - Session information
// - Security context
// - User agent and IP address

resp, err := client.ExecuteFromHTTPRequest(ctx, workflowID, httpRequest)
```

### 3. Error Handling
Comprehensive error handling with specific error types:

```go
resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
if err != nil {
    switch err := err.(type) {
    case *sdk.SDKError:
        switch err.Code {
        case sdk.ErrCodeValidationFailed:
            // Handle validation errors
            fmt.Printf("Validation failed: %v\n", err.Details)
        case sdk.ErrCodeRequestParsingFailed:
            // Handle parsing errors
            fmt.Printf("Failed to parse request: %v\n", err.OriginalError)
        case sdk.ErrCodeWorkflowExecution:
            // Handle execution errors
            fmt.Printf("Workflow execution failed: %v\n", err.Message)
        default:
            // Handle other SDK errors
            fmt.Printf("SDK error: %v\n", err)
        }
    default:
        // Handle non-SDK errors
        fmt.Printf("Unexpected error: %v\n", err)
    }
    return
}
```

### 4. Retry and Circuit Breaker
Built-in retry logic and circuit breaker:

```go
config := &sdk.SDKConfig{
    WorkflowAPIEndpoint: "http://localhost:8081",
    MaxRetries:          5,                    // Retry up to 5 times
    RetryDelay:          2 * time.Second,      // Wait 2 seconds between retries
    // Circuit breaker automatically opens after multiple failures
}
```

## Best Practices

### 1. Always Close the Client
```go
client, err := sdk.NewClient(config)
if err != nil {
    // Handle error
}
defer client.Close() // Always close the client
```

### 2. Use Context for Cancellation
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
```

### 3. Enable Validation for Production
```go
config := &sdk.SDKConfig{
    EnableValidation:   true,
    EnableSanitization: true,
    StrictValidation:   true, // Fail fast on validation errors
}
```

### 4. Handle Async Execution
```go
resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
if err != nil {
    // Handle error
}

// Poll for completion
for {
    statusResp, err := client.GetExecutionStatus(ctx, resp.RunID)
    if err != nil {
        // Handle error
        break
    }
    
    if statusResp.Status.IsTerminal {
        fmt.Printf("Workflow completed with status: %s\n", statusResp.Status.Status)
        break
    }
    
    time.Sleep(time.Duration(resp.PollAfterMs) * time.Millisecond)
}
```

### 5. Logging and Monitoring
```go
// Configure structured logging
config := &sdk.SDKConfig{
    WorkflowAPIEndpoint: "http://localhost:8081",
    // SDK automatically logs:
    // - Request/response timing
    // - Error details
    // - Validation results
    // - Retry attempts
}

// Check service health
err := client.Ping(ctx)
if err != nil {
    log.Printf("Service is unreachable: %v", err)
}
```

## Common Use Cases

### 1. Payment Processing
```go
func processPayment(ctx context.Context, paymentData map[string]interface{}) error {
    client, err := sdk.NewClient(config)
    if err != nil {
        return err
    }
    defer client.Close()
    
    // Define payment validation rules
    rules := []sdk.ValidationRule{
        {
            Field:    "amount",
            Required: true,
            RuleType: sdk.ValidationRuleTypeNumber,
            MinValue: &[]float64{0.01}[0],
        },
        {
            Field:    "currency",
            Required: true,
            RuleType: sdk.ValidationRuleTypeString,
            AllowedValues: []string{"USD", "EUR", "GBP"},
        },
        {
            Field:    "payment_method",
            Required: true,
            RuleType: sdk.ValidationRuleTypeString,
        },
    }
    
    // Execute payment workflow
    resp, err := client.ValidateAndExecuteWorkflow(ctx, "payment-processing", paymentData, rules)
    if err != nil {
        return fmt.Errorf("payment processing failed: %w", err)
    }
    
    // Monitor payment status
    return monitorPayment(ctx, client, resp.RunID)
}
```

### 2. Order Fulfillment
```go
func fulfillOrder(ctx context.Context, orderData map[string]interface{}) (string, error) {
    client, err := sdk.NewClient(config)
    if err != nil {
        return "", err
    }
    defer client.Close()
    
    resp, err := client.ExecuteWorkflow(ctx, "order-fulfillment", orderData)
    if err != nil {
        return "", fmt.Errorf("order fulfillment failed: %w", err)
    }
    
    return resp.RunID, nil
}
```

### 3. User Registration
```go
func registerUser(ctx context.Context, userData map[string]interface{}) error {
    client, err := sdk.NewClient(config)
    if err != nil {
        return err
    }
    defer client.Close()
    
    // Extract session from HTTP context
    resp, err := client.ExecuteFromHTTPRequest(ctx, "user-registration", httpRequest)
    if err != nil {
        return fmt.Errorf("user registration failed: %w", err)
    }
    
    log.Printf("User registration started: %s", resp.RunID)
    return nil
}
```

## Troubleshooting

### Common Issues

#### 1. Connection Errors
```go
// Check if the service is reachable
err := client.Ping(ctx)
if err != nil {
    // Service is unreachable
    // Check:
    // - Network connectivity
    // - Firewall rules
    // - Service status
}
```

#### 2. Validation Errors
```go
resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
if err != nil {
    if sdkErr, ok := err.(*sdk.SDKError); ok && sdkErr.Code == sdk.ErrCodeValidationFailed {
        // Print validation details
        if details, ok := sdkErr.Details["validation_result"].(*sdk.ValidationResult); ok {
            for _, validationErr := range details.Errors {
                fmt.Printf("Field: %s, Error: %s\n", validationErr.Field, validationErr.Message)
            }
        }
    }
}
```

#### 3. Authentication Errors
```go
// Ensure auth token is valid
config := &sdk.SDKConfig{
    WorkflowAPIEndpoint: "http://localhost:8081",
    AuthToken:          "your-valid-token", // Check token expiration
}
```

#### 4. Timeout Errors
```go
// Increase timeout for long-running workflows
config := &sdk.SDKConfig{
    WorkflowAPIEndpoint: "http://localhost:8081",
    Timeout:             5 * time.Minute, // Increase timeout
}

// Or use context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
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

### Response Types
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

## Support

### Getting Help
- **Documentation**: Complete API reference and examples
- **GitHub Issues**: Report bugs and request features
- **Stack Overflow**: Tag questions with `unified-workflow-sdk`
- **Email Support**: enterprise-support@your-org.com

### Community Resources
- **Sample Applications**: GitHub repository examples
- **Blog Posts**: Integration tutorials and best practices
- **Video Tutorials**: Step-by-step guides
- **Office Hours**: Weekly Q&A sessions

## Migration Guide

### From Direct API Calls
If you're migrating from direct HTTP API calls to the SDK:

```go
// Before: Direct HTTP calls
resp, err := http.Post("http://localhost:8081/api/v1/execute", "application/json", bytes)

// After: SDK usage
client, err := sdk.NewClient(config)
resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
```

Benefits:
- Type safety
- Built-in validation
- Automatic retry logic
- Better error handling
- Connection pooling

## Conclusion

The Unified Workflow SDK simplifies integration with the workflow execution platform while providing robust features for production use. Follow the best practices outlined in this guide to build reliable, maintainable applications.

For more detailed information, refer to the complete API documentation and example applications in the SDK repository.