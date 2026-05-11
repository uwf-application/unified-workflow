# Unified Workflow SDK v1.2.0

A Go SDK for integrating with the Unified Workflow Execution Platform.

## Installation

### Download the SDK tarball

Download `unified-workflow-sdk-v1.2.0.tar.gz` from the [release page](https://github.com/uwf-application/unified-workflow/releases/tag/v1.2.0):

```bash
curl -L -o unified-workflow-sdk-v1.2.0.tar.gz \
  https://github.com/uwf-application/unified-workflow/releases/download/v1.2.0/unified-workflow-sdk-v1.2.0.tar.gz

tar -xzf unified-workflow-sdk-v1.2.0.tar.gz
```

Then in your project's `go.mod`:

```go
module your-project

go 1.21.0

require github.com/uwf-application/unified-workflow-sdk v0.0.0

replace github.com/uwf-application/unified-workflow-sdk => ./unified-workflow-sdk-v1.2.0
```

Run:
```bash
go mod tidy
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    sdk "github.com/uwf-application/unified-workflow-sdk"
)

func main() {
    config := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://<workflow-host>:8081",
        Timeout:             30 * time.Second,
        MaxRetries:          3,
        AuthToken:           "your-api-token",
        EnableValidation:    true,
    }

    client, err := sdk.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()

    // Health check
    if err := client.Ping(ctx); err != nil {
        log.Fatalf("Service unreachable: %v", err)
    }

    // Execute a workflow
    resp, err := client.ExecuteWorkflow(ctx, "your-workflow-id", map[string]interface{}{
        "user_id": "12345",
        "amount":  99.99,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Run ID: %s  Status: %s\n", resp.RunID, resp.Status)

    // Poll for completion
    time.Sleep(2 * time.Second)
    status, err := client.GetExecutionStatus(ctx, resp.RunID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %s  Progress: %.0f%%\n", status.Status.Status, status.Status.Progress*100)
}
```

## Configuration via environment variables

```bash
export SDK_WORKFLOW_API_ENDPOINT=http://<workflow-host>:8081
export SDK_AUTH_TOKEN=your-api-token
export SDK_TIMEOUT=30           # seconds
export SDK_MAX_RETRIES=3
export SDK_ASYNC_EXECUTION=true
export SDK_ENABLE_VALIDATION=true
```

Then use:
```go
config, err := sdk.LoadDefaultConfig()
```

## API Reference

| Method | Description |
|--------|-------------|
| `NewClient(config)` | Create a new SDK client |
| `client.Ping(ctx)` | Health check the workflow service |
| `client.ExecuteWorkflow(ctx, workflowID, data)` | Trigger a workflow with input data |
| `client.ExecuteFromHTTPRequest(ctx, workflowID, req)` | Trigger from an incoming HTTP request |
| `client.ValidateAndExecuteWorkflow(ctx, id, data, rules)` | Validate input then execute |
| `client.BatchExecuteWorkflows(ctx, batchReq)` | Execute multiple workflows in one call |
| `client.GetExecutionStatus(ctx, runID)` | Poll execution status |
| `client.GetExecutionDetails(ctx, runID)` | Full execution details and output |
| `client.CancelExecution(ctx, runID)` | Cancel a running workflow |

## Examples

- [`examples/basic/`](examples/basic/) — Basic workflow execution
- [`examples/http/`](examples/http/) — Wrapping an HTTP handler

## Support

- Issues: https://github.com/uwf-application/unified-workflow/issues

## License

MIT License - See LICENSE file for details.
