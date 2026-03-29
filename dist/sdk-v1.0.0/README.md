# Unified Workflow SDK

A comprehensive Go SDK for interacting with the Unified Workflow Execution Platform.

## Installation

```bash
go get github.com/your-org/unified-workflow-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/your-org/unified-workflow-sdk"
)

func main() {
    config := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://localhost:8081",
        Timeout:             30 * time.Second,
        MaxRetries:          3,
        EnableValidation:    true,
    }
    
    client, err := sdk.NewClient(config)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    ctx := context.Background()
    inputData := map[string]interface{}{
        "user_id": "test_user",
        "amount":  99.99,
    }
    
    resp, err := client.ExecuteWorkflow(ctx, "payment-workflow", inputData)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Workflow execution started: %s\n", resp.RunID)
}
```

## Features

- **Simple API**: Clean, type-safe interface for workflow execution
- **Validation**: Built-in data validation with customizable rules
- **Error Handling**: Comprehensive error handling with retry logic
- **HTTP Integration**: Easy integration with HTTP handlers
- **Async Support**: Support for asynchronous workflow execution
- **Monitoring**: Execution status and progress tracking

## Documentation

- [Client Guide](docs/CLIENT_GUIDE.md) - Complete usage guide
- [Release Plan](docs/RELEASE_PLAN.md) - Release strategy and packaging
- [Examples](examples/) - Complete example applications

## Examples

See the [examples directory](examples/) for complete examples:

1. `basic_usage.go` - Basic SDK usage
2. `http_integration.go` - HTTP server integration

## Support

- **Issues**: [GitHub Issues](https://github.com/your-org/unified-workflow-sdk/issues)
- **Documentation**: [Client Guide](docs/CLIENT_GUIDE.md)
- **Email**: support@your-org.com

## License

MIT License - See LICENSE file for details.
