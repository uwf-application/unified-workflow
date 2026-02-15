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

The antifraud service can be integrated into workflow steps:

```go
package workflows

import (
    "unified-workflow/internal/primitive"
)

type AntifraudStep struct {
    // Step implementation
}

func (s *AntifraudStep) Execute(ctx WorkflowContext, data WorkflowData) error {
    // Get transaction from workflow data
    transaction, ok := data.Get("transaction")
    if !ok {
        return fmt.Errorf("transaction not found in workflow data")
    }
    
    // Use antifraud service
    err := primitive.Default.Antifraud.StoreTransaction(transaction)
    if err != nil {
        return fmt.Errorf("failed to store transaction: %w", err)
    }
    
    // Validate transaction
    result, err := primitive.Default.Antifraud.ValidateTransactionByAML(transaction)
    if err != nil {
        return fmt.Errorf("AML validation failed: %w", err)
    }
    
    // Store result in workflow data
    data.Put("aml_result", result)
    
    return nil
}
```

## Migration from Direct SDK Usage

### Before (Direct SDK)

```go
import (
    af "github.com/baraic-io/antifraud-go"
)

func main() {
    client, err := af.NewClient(af.ClientConfig{
        Host:   os.Getenv("API_HOST"),
        APIKey: os.Getenv("API_KEY"),
    })
    
    err = client.StoreTransaction(transaction)
}
```

### After (Primitive Integration)

```go
import (
    "unified-workflow/internal/primitive"
)

func main() {
    primitive.Init(&primitive.Config{
        AntifraudAPIKey:  os.Getenv("API_KEY"),
        AntifraudAPIHost: os.Getenv("API_HOST"),
    })
    
    err := primitive.Default.Antifraud.StoreTransaction(transaction)
}
```

## Best Practices

1. **Always check if service is enabled** before using it
2. **Use circuit breaker** for production deployments
3. **Set appropriate timeouts** based on your SLA requirements
4. **Monitor health checks** in production environments
5. **Use environment variables** for configuration
6. **Implement retry logic** for transient failures
7. **Log all antifraud operations** for audit trails

## Troubleshooting

### Common Issues

1. **"API key is required"**: Set `ANTIFRAUD_API_KEY` environment variable
2. **"Host is required"**: Set `ANTIFRAUD_API_HOST` environment variable
3. **Circuit breaker constantly open**: Check service health and adjust thresholds
4. **Timeouts**: Increase `AntifraudTimeout` configuration
5. **Disabled service**: Set `AntifraudEnabled: true` in configuration

### Debugging

Enable debug logging:

```go
config := &primitive.Config{
    AntifraudAPIKey:  "your-key",
    AntifraudAPIHost: "your-host",
    // Add debug logging configuration
}
```

## Performance Considerations

1. **Connection Pooling**: The SDK manages HTTP connection pooling
2. **Circuit Breaker**: Prevents cascading failures during outages
3. **Timeout Management**: Configurable timeouts prevent hung requests
4. **Retry Logic**: Automatic retries for transient failures
5. **Caching**: Consider implementing caching for repeated validations

## Security

1. **API Key Security**: Store API keys in secure vaults (not in code)
2. **TLS/SSL**: Always use HTTPS for API communication
3. **Input Validation**: Validate all transaction data before sending
4. **Audit Logging**: Log all antifraud operations for compliance
5. **Access Control**: Implement proper access controls for antifraud operations

## Support

For issues with the antifraud primitive integration:

1. Check the [examples/antifraud_demo.go](examples/antifraud_demo.go) file
2. Review the test cases in [internal/primitive/services/antifraud/client_test.go](internal/primitive/services/antifraud/client_test.go)
3. Consult the Antifraud SDK documentation at `github.com/baraic-io/antifraud-go`

## License

This integration is part of the Unified Workflow project. See the project LICENSE for details.

---

**Last Updated**: February 15, 2026  
**Version**: 1.0.0  
**Author**: Unified Workflow Team