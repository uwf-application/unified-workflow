# Antifraud Primitive Integration

This document describes the integration of the Antifraud Go SDK as a primitive operation within the Unified Workflow framework using Dependency Injection (DI) and the default namespace pattern.

## Overview

The Antifraud SDK (`github.com/baraic-io/antifraud-go`) has been integrated as a primitive service following the existing patterns in the codebase. This allows developers to use the antifraud service through:

1. **Global Primitive**: `primitive.Default.Antifraud`
2. **DI Container**: Through the `PrimitiveProvider`
3. **Direct Instantiation**: Using the `antifraud.NewClient()` function

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Primitive Layer                       │
├─────────────────────────────────────────────────────────┤
│  AntifraudService (Interface)                           │
│  ├── StoreTransaction()                                 │
│  ├── ValidateTransactionByAML()                         │
│  ├── ValidateTransactionByFC()                          │
│  ├── ValidateTransactionByML()                          │
│  ├── StoreServiceResolution()                           │
│  ├── AddTransactionServiceCheck()                       │
│  ├── FinalizeTransaction()                              │
│  └── StoreFinalResolution()                             │
├─────────────────────────────────────────────────────────┤
│  AntifraudClientImpl (Implementation)                   │
│  ├── Wraps github.com/baraic-io/antifraud-go           │
│  ├── Handles configuration & authentication            │
│  └── Implements error handling & logging               │
├─────────────────────────────────────────────────────────┤
│  AntifraudProxy (Decorator)                             │
│  ├── Adds logging, metrics, circuit breaking           │
│  ├── Implements retry logic                            │
│  └── Provides consistent error handling                │
└─────────────────────────────────────────────────────────┘
```

## Installation

### 1. Add SDK Dependency

Add the Antifraud SDK to your `go.mod`:

```bash
go get github.com/baraic-io/antifraud-go
```

### 2. Update Configuration

The antifraud service is configured through the `primitive.Config` struct:

```go
config := &primitive.Config{
    // Antifraud configuration
    AntifraudAPIKey:                  os.Getenv("ANTIFRAUD_API_KEY"),
    AntifraudAPIHost:                 os.Getenv("ANTIFRAUD_API_HOST"),
    AntifraudTimeout:                 30, // seconds
    AntifraudEnabled:                 true,
    AntifraudMaxRetries:              3,
    AntifraudCircuitBreakerEnabled:   true,
    AntifraudCircuitBreakerThreshold: 5,
    AntifraudCircuitBreakerTimeout:   60, // seconds
}
```

## Usage Examples

### Example 1: Using Global Primitive

```go
package main

import (
    "fmt"
    "os"
    "time"
    
    "unified-workflow/internal/primitive"
    "github.com/google/uuid"
)

func main() {
    // Initialize primitive
    config := &primitive.Config{
        AntifraudAPIKey:    os.Getenv("ANTIFRAUD_API_KEY"),
        AntifraudAPIHost:   os.Getenv("ANTIFRAUD_API_HOST"),
        AntifraudTimeout:   30,
        AntifraudEnabled:   true,
    }
    
    primitive.Init(config)
    
    // Create transaction
    transaction := map[string]interface{}{
        "af_id":       uuid.NewString(),
        "af_add_date": time.Now().Format(time.RFC3339Nano),
        "transaction": map[string]interface{}{
            "id":                   uuid.NewString(),
            "type":                 "deposit",
            "amount":               "100000",
            "currency":             "KZT",
            // ... other fields
        },
    }
    
    // Use antifraud service
    err := primitive.Default.Antifraud.StoreTransaction(transaction)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("Transaction stored successfully")
}
```

### Example 2: Using DI Container

```go
package main

import (
    "fmt"
    
    "unified-workflow/internal/di"
    "unified-workflow/internal/primitive"
)

func main() {
    // Create DI container
    container := di.New()
    
    // Initialize primitive with config
    config := &primitive.Config{
        AntifraudAPIKey:  "your-api-key",
        AntifraudAPIHost: "https://api.antifraud.example.com",
        AntifraudEnabled: true,
    }
    
    // Register primitive services
    err := di.RegisterPrimitiveServices(container, config)
    if err != nil {
        panic(err)
    }
    
    // Create primitive provider
    provider := di.NewPrimitiveProvider(container)
    
    // Get antifraud service
    antifraudService := provider.Antifraud()
    
    // Use the service
    // ... transaction validation logic
}
```

### Example 3: Direct Client Usage

```go
package main

import (
    "fmt"
    "time"
    
    "unified-workflow/internal/primitive/services/antifraud"
    "github.com/google/uuid"
)

func main() {
    // Create client config
    config := antifraud.ClientConfig{
        APIKey:    "your-api-key",
        Host:      "https://api.antifraud.example.com",
        Timeout:   30,
        Enabled:   true,
        MaxRetries: 3,
    }
    
    // Create client
    client, err := antifraud.NewClient(config)
    if err != nil {
        panic(err)
    }
    
    // Create transaction
    transaction := antifraud.AF_Transaction{
        AF_Id:      uuid.NewString(),
        AF_AddDate: time.Now().Format(time.RFC3339Nano),
        Transaction: antifraud.Transaction{
            Id:                 uuid.NewString(),
            Type:               "deposit",
            Amount:             "100000",
            Currency:           "KZT",
            ClientId:           uuid.NewString(),
            ClientName:         "John Smith",
            ClientPAN:          "111111******1111",
            ClientCVV:          "111",
            ClientCardHolder:   "JOHN SMITH",
            ClientPhone:        "+77007007070",
            MerchantTerminalId: "00000001",
            Channel:            "E-com",
            LocationIp:         "192.168.0.1",
        },
    }
    
    // Store transaction
    err = client.StoreTransaction(transaction)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("Transaction stored successfully")
}
```

## Full Transaction Validation Flow

The antifraud service supports a complete transaction validation flow:

```go
// 1. Store transaction
err := client.StoreTransaction(transaction)

// 2. Validate with AML service
amlResult, err := client.ValidateTransactionByAML(transaction)

// 3. Store AML resolution
err = client.StoreServiceResolution(amlResult)

// 4. Add AML check to transaction
err = client.AddTransactionServiceCheck(amlResult)

// 5. Validate with FC service
fcResult, err := client.ValidateTransactionByFC(transaction)

// 6. Store FC resolution
err = client.StoreServiceResolution(fcResult)

// 7. Add FC check to transaction
err = client.AddTransactionServiceCheck(fcResult)

// 8. Validate with ML service
mlResult, err := client.ValidateTransactionByML(transaction)

// 9. Store ML resolution
err = client.StoreServiceResolution(mlResult)

// 10. Add ML check to transaction
err = client.AddTransactionServiceCheck(mlResult)

// 11. Finalize transaction
finalResult, err := client.FinalizeTransaction(transaction)

// 12. Store final resolution
err = client.StoreFinalResolution(finalResult)
```

## Configuration Options

### ClientConfig

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `APIKey` | `string` | Required | API key for authentication |
| `Host` | `string` | Required | API host URL |
| `Timeout` | `int` | 30 | Request timeout in seconds |
| `Enabled` | `bool` | true | Enable/disable the service |
| `MaxRetries` | `int` | 3 | Maximum retry attempts for failed requests |
| `CircuitBreakerEnabled` | `bool` | false | Enable circuit breaker pattern |
| `CircuitBreakerThreshold` | `int` | 5 | Number of failures before opening circuit |
| `CircuitBreakerTimeout` | `int` | 60 | Circuit breaker reset timeout in seconds |

### Environment Variables

```bash
# Required
export ANTIFRAUD_API_KEY="your-api-key"
export ANTIFRAUD_API_HOST="https://api.antifraud.example.com"

# Optional
export ANTIFRAUD_TIMEOUT="30"
export ANTIFRAUD_MAX_RETRIES="3"
export ANTIFRAUD_CIRCUIT_BREAKER_ENABLED="true"
export ANTIFRAUD_CIRCUIT_BREAKER_THRESHOLD="5"
export ANTIFRAUD_CIRCUIT_BREAKER_TIMEOUT="60"
```

## Error Handling

The antifraud service provides comprehensive error handling:

1. **Validation Errors**: Invalid configuration or transaction data
2. **Network Errors**: Connection timeouts, HTTP errors
3. **Circuit Breaker**: Automatic fail-fast when service is unhealthy
4. **Disabled Service**: Graceful handling when service is disabled

```go
// Example error handling
err := client.StoreTransaction(transaction)
if err != nil {
    switch {
    case errors.Is(err, antifraud.ErrInvalidConfig):
        // Handle configuration error
    case errors.Is(err, antifraud.ErrServiceDisabled):
        // Handle disabled service
    case errors.Is(err, antifraud.ErrCircuitBreakerOpen):
        // Handle circuit breaker open
    default:
        // Handle other errors
    }
}
```

## Circuit Breaker Pattern

The antifraud proxy implements a circuit breaker pattern to prevent cascading failures:

- **Closed State**: Normal operation, requests pass through
- **Open State**: Circuit is open, requests fail fast
- **Half-Open State**: Testing if service has recovered

```go
// Check circuit breaker status
proxy, ok := client.(*antifraud.antifraudProxy)
if ok {
    isOpen, failureCount, openedAt := proxy.GetCircuitBreakerStatus()
    fmt.Printf("Circuit breaker: open=%v, failures=%d, opened at=%v\n",
        isOpen, failureCount, openedAt)
}
```

## Health Checking

```go
// Check service health
healthy, err := client.HealthCheck()
if err != nil {
    fmt.Printf("Health check failed: %v\n", err)
} else if healthy {
    fmt.Println("Service is healthy")
} else {
    fmt.Println("Service is unhealthy")
}
```

## Testing

### Unit Tests

Run the antifraud service tests:

```bash
go test ./internal/primitive/services/antifraud/... -v
```

### Example Tests

See the example file for comprehensive usage examples:

```bash
go run examples/antifraud_demo.go
```

## Integration with Workflows

The antifraud service has been fully integrated into the workflow system with complete transaction validation workflow:

### Complete Antifraud Workflow

The `antifraud-transaction-validation` workflow includes 5 sequential steps:

1. **StoreTransactionStep** - Stores transaction in antifraud system
2. **AMLValidationStep** - Anti-Money Laundering validation  
3. **FCValidationStep** - Fraud Check validation
4. **MLValidationStep** - Machine Learning validation
5. **FinalizeTransactionStep** - Final decision and resolution

### Workflow Implementation

```go
// workflows/antifraud_workflow.go
func CreateAntifraudTransactionWorkflow(endpoint string) model.Workflow {
    workflow := model.NewBaseWorkflow(
        "antifraud-transaction-validation",
        fmt.Sprintf("Complete transaction validation using antifraud SDK at %s", endpoint),
    )
    
    // Step 1: Store Transaction
    storeStep := steps.NewStoreTransactionStep(endpoint)
    workflow.AddStep(storeStep)
    
    // Step 2: AML Validation
    amlStep := steps.NewAMLValidationStep(endpoint)
    workflow.AddStep(amlStep)
    
    // Step 3: FC Validation (Fraud Check)
    fcStep := steps.NewFCValidationStep(endpoint)
    workflow.AddStep(fcStep)
    
    // Step 4: ML Validation (Machine Learning)
    mlStep := steps.NewMLValidationStep(endpoint)
    workflow.AddStep(mlStep)
    
    // Step 5: Finalize Transaction
    finalizeStep := steps.NewFinalizeTransactionStep(endpoint)
    workflow.AddStep(finalizeStep)
    
    return workflow
}
```

### Step Implementation Example

```go
// workflows/steps/fc_validation_step.go
type FCValidationStep struct {
    *AntifraudStep
}

func (s *FCValidationStep) ExecuteStepLogic(ctx interface{}, context interface{}, data interface{}) error {
    // Get antifraud service
    antifraudService, err := s.GetAntifraudService()
    if err != nil {
        return fmt.Errorf("failed to get antifraud service: %w", err)
    }
    
    // Prepare transaction data
    transactionData := map[string]interface{}{
        "AF_Id":      uuid.NewString(),
        "AF_AddDate": time.Now().Format(time.RFC3339Nano),
        "Transaction": map[string]interface{}{
            "Id":                 uuid.NewString(),
            "Type":               "deposit",
            "Amount":             "100000",
            "Currency":           "KZT",
            // ... transaction fields
        },
    }
    
    // Call antifraud SDK
    result, err := antifraudService.ValidateTransactionByFC(transactionData)
    if err != nil {
        return fmt.Errorf("FC validation failed: %w", err)
    }
    
    // Validate response
    if err := s.validateFCResponse(result); err != nil {
        return fmt.Errorf("FC response validation failed: %w", err)
    }
    
    // Store resolution
    err = antifraudService.StoreServiceResolution(result)
    if err != nil {
        fmt.Printf("Warning: Failed to store FC resolution: %v\n", err)
    }
    
    return nil
}
```

### Demo Workflow Execution

A complete demo is available at `examples/antifraud_workflow_demo.go`:

```go
// examples/antifraud_workflow_demo.go
func main() {
    // Configure SDK for TAF environment
    sdkConfig := sdk.DefaultConfig()
    sdkConfig.WorkflowAPIEndpoint = "https://af-test.qazpost.kz"
    sdkConfig.Timeout = 30
    
    // Create SDK client
    client, err := sdk.NewClient(&sdkConfig)
    
    // Prepare transaction data
    transactionData := map[string]interface{}{
        "transaction": map[string]interface{}{
            "id":                   fmt.Sprintf("txn-%d", time.Now().Unix()),
            "type":                 "deposit",
            "amount":               "100000",
            "currency":             "KZT",
            // ... transaction fields
        },
    }
    
    // Execute antifraud workflow
    resp, err := client.ExecuteWorkflow(ctx, "antifraud-transaction-validation", transactionData)
    
    // Poll for execution status
    statusResp, err := client.GetExecutionStatus(ctx, resp.RunID)
    
    // Get execution details
    detailsResp, err := client.GetExecutionDetails(ctx, resp.RunID)
}
```

## Testing on Jump Server

The antifraud workflow can be tested on the Qazpost jump server (10.200.1.2) using the following steps:

### 1. Sync Code to Jump Server

```bash
./uwf-cli deploy sync --verbose
```

### 2. Build Docker Images

```bash
./uwf-cli deploy build --verbose
```

### 3. Push Images to Harbor Registry

```bash
./uwf-cli deploy push --all --verbose
```

### 4. Test Connectivity

```bash
# From jump server
ssh khassangali@10.200.1.2

# Test TAF service connectivity
ping 172.30.75.91
curl -k https://af-test.qazpost.kz

# Run test script
cd /tmp/uwf-deploy
./test_antifraud_workflow.sh
```

### 5. Execute Antifraud Workflow

From the jump server, you can execute the antifraud workflow using:

```bash
# Option 1: Using uwf-cli
./uwf-cli execute sync antifraud-transaction-validation \
  --input '{
    "transaction": {
      "id": "txn-12345",
      "type": "deposit",
      "amount": "100000",
      "currency": "KZT",
      "client_id": "client-001",
      "client_name": "John Smith",
      "client_pan": "111111******1111",
      "client_cvv": "111",
      "client_card_holder": "JOHN SMITH",
      "client_phone": "+77007007070",
      "merchant_terminal_id": "00000001",
      "channel": "E-com",
      "location_ip": "192.168.0.1"
    }
  }'

# Option 2: Using curl directly
curl -X POST https://af-test.qazpost.kz/api/v1/workflows/antifraud-transaction-validation/execute \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer M8#Qe!2$ZrA9xKp' \
  -d '{
    "input_data": {
      "transaction": {
        "id": "txn-12345",
        "type": "deposit",
        "amount": "100000",
        "currency": "KZT",
        "client_id": "client-001",
        "client_name": "John Smith",
        "client_pan": "111111******1111",
        "client_cvv": "111",
        "client_card_holder": "JOHN SMITH",
        "client_phone": "+77007007070",
        "merchant_terminal_id": "00000001",
        "channel": "E-com",
        "location_ip": "192.168.0.1"
      }
    }
  }'
```

### 6. Verify Deployment

```bash
./uwf-cli deploy verify
```

## Workflow Lifecycle Verification

The antifraud workflow follows this complete lifecycle:

```
┌─────────────────────────────────────────────────────────┐
│         Workflow Initiation (Using SDK)                 │
├─────────────────────────────────────────────────────────┤
│  SDK Client → ExecuteWorkflow() → Workflow API          │
├─────────────────────────────────────────────────────────┤
│         Workflow Execution                              │
├─────────────────────────────────────────────────────────┤
│  Workflow → Steps → ChildSteps → PrimitiveOperations    │
├─────────────────────────────────────────────────────────┤
│         Primitive Operations                            │
├─────────────────────────────────────────────────────────┤
│  Primitive → ServiceClients → SDK → Actual Client       │
├─────────────────────────────────────────────────────────┤
│         TAF Antifraud Service                           │
├─────────────────────────────────────────────────────────┤
│  Actual Client → TAF API (172.30.75.91)                 │
└─────────────────────────────────────────────────────────┘
```

### Verification Points

1. **SDK Level**: `ExecuteWorkflow()` returns execution ID
2. **Workflow Level**: Steps execute in sequence (Store → AML → FC → ML → Finalize)
3. **Primitive Level**: Antifraud service calls are made with proper authentication
4. **Service Client Level**: HTTP requests to TAF service with correct payload
5. **TAF Service Level**: Transaction validation results returned

## Support

For issues with the antifraud primitive integration:

1. Check the [examples/antifraud_demo.go](examples/antifraud_demo.go) file
2. Check the [examples/antifraud_workflow_demo.go](examples/antifraud_workflow_demo.go) file
3. Review the test cases in [internal/primitive/services/antifraud/client_test.go](internal/primitive/services/antifraud/client_test.go)
4. Run the test script: `./test_antifraud_workflow.sh` on jump server
5. Consult the Antifraud SDK documentation at `github.com/baraic-io/antifraud-go`

This integration is part of the Unified Workflow project. See the project LICENSE for details.

---

**Last Updated**: March 2, 2026  
**Version**: 1.1.0  
**Author**: Unified Workflow Team  
**Status**: ✅ Complete - Ready for deployment and testing on Qazpost TAF environment
