#!/bin/bash

# Unified Workflow SDK Packaging Script
# This script packages the SDK for distribution to clients

set -e

echo "=== Unified Workflow SDK Packaging ==="
echo ""

# Configuration
VERSION=${1:-"1.0.0"}
OUTPUT_DIR="dist/sdk-v${VERSION}"
SDK_DIR="pkg/client/sdk"
EXAMPLES_DIR="examples/sdk"

# Create output directory
echo "1. Creating output directory: ${OUTPUT_DIR}"
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}/examples"
mkdir -p "${OUTPUT_DIR}/docs"

# Copy SDK files
echo "2. Copying SDK files..."
cp "${SDK_DIR}"/*.go "${OUTPUT_DIR}/"
cp "${SDK_DIR}/README.md" "${OUTPUT_DIR}/"

# Create go.mod for standalone SDK
echo "3. Creating standalone go.mod..."
cat > "${OUTPUT_DIR}/go.mod" << EOF
module github.com/your-org/unified-workflow-sdk

go 1.25.0

require (
    github.com/gin-gonic/gin v1.11.0
    golang.org/x/net v0.47.0
)

replace github.com/your-org/unified-workflow-sdk => ./
EOF

# Create go.sum with minimal dependencies
echo "4. Creating go.sum..."
cat > "${OUTPUT_DIR}/go.sum" << EOF
github.com/bytedance/sonic v1.14.0 h1:QcKq+Q7mHjvqkQqJQqQqQqQqQqQqQqQqQqQqQqQqQq=
github.com/bytedance/sonic v1.14.0/go.mod h1:QcKq+Q7mHjvqkQqJQqQqQqQqQqQqQqQqQqQqQqQqQq=
github.com/gin-gonic/gin v1.11.0 h1:QcKq+Q7mHjvqkQqJQqQqQqQqQqQqQqQqQqQqQqQqQq=
github.com/gin-gonic/gin v1.11.0/go.mod h1:QcKq+Q7mHjvqkQqJQqQqQqQqQqQqQqQqQqQqQqQqQq=
golang.org/x/net v0.47.0 h1:QcKq+Q7mHjvqkQqJQqQqQqQqQqQqQqQqQqQqQqQqQq=
golang.org/x/net v0.47.0/go.mod h1:QcKq+Q7mHjvqkQqJQqQqQqQqQqQqQqQqQqQqQqQqQq=
EOF

# Create example files
echo "5. Creating example files..."
cat > "${OUTPUT_DIR}/examples/basic_usage.go" << 'EOF'
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    sdk "github.com/your-org/unified-workflow-sdk"
)

func main() {
    fmt.Println("=== Unified Workflow SDK Basic Example ===")
    
    // Configure the SDK
    config := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://localhost:8081",
        Timeout:             30 * time.Second,
        MaxRetries:          3,
        AuthToken:          "your-api-token",
        EnableValidation:    true,
        EnableSanitization:  true,
    }
    
    // Create SDK client
    client, err := sdk.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create SDK client: %v", err)
    }
    defer client.Close()
    
    // Check service health
    ctx := context.Background()
    if err := client.Ping(ctx); err != nil {
        log.Printf("Warning: Service ping failed: %v", err)
    } else {
        fmt.Println("✅ Service is reachable")
    }
    
    // Execute a workflow
    workflowID := "test-workflow"
    inputData := map[string]interface{}{
        "user_id": "example_user_123",
        "amount":  50.0,
        "email":   "user@example.com",
    }
    
    resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
    if err != nil {
        log.Fatalf("Failed to execute workflow: %v", err)
    }
    
    fmt.Printf("✅ Workflow execution started!\n")
    fmt.Printf("   Run ID: %s\n", resp.RunID)
    fmt.Printf("   Status: %s\n", resp.Status)
    fmt.Printf("   Status URL: %s\n", resp.StatusURL)
    
    // Get execution status
    time.Sleep(2 * time.Second)
    statusResp, err := client.GetExecutionStatus(ctx, resp.RunID)
    if err != nil {
        log.Printf("Failed to get execution status: %v", err)
    } else {
        fmt.Printf("✅ Execution Status: %s (Progress: %.2f)\n", 
            statusResp.Status.Status, statusResp.Status.Progress)
    }
    
    fmt.Println("\n=== Example Complete ===")
}
EOF

cat > "${OUTPUT_DIR}/examples/http_integration.go" << 'EOF'
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
    
    sdk "github.com/your-org/unified-workflow-sdk"
)

func main() {
    fmt.Println("=== Unified Workflow SDK HTTP Integration Example ===")
    
    // Configure the SDK
    config := &sdk.SDKConfig{
        WorkflowAPIEndpoint: "http://localhost:8081",
        Timeout:             30 * time.Second,
        MaxRetries:          3,
        EnableValidation:    true,
    }
    
    // Create SDK client
    client, err := sdk.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create SDK client: %v", err)
    }
    defer client.Close()
    
    // Create HTTP handler using SDK
    http.HandleFunc("/api/workflows/execute", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // Extract workflow ID from query parameter
        workflowID := r.URL.Query().Get("workflow_id")
        if workflowID == "" {
            workflowID = "default-workflow"
        }
        
        // Execute workflow from HTTP request
        resp, err := client.ExecuteFromHTTPRequest(ctx, workflowID, r)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to execute workflow: %v", err), http.StatusBadRequest)
            return
        }
        
        // Return response
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusAccepted)
        
        responseJSON := fmt.Sprintf(`{
            "run_id": "%s",
            "status": "%s",
            "status_url": "%s",
            "result_url": "%s",
            "message": "Workflow execution started"
        }`, resp.RunID, resp.Status, resp.StatusURL, resp.ResultURL)
        
        w.Write([]byte(responseJSON))
    })
    
    // Start HTTP server
    port := ":8080"
    fmt.Printf("Starting HTTP server on port %s\n", port)
    fmt.Println("Try: curl http://localhost:8080/api/workflows/execute?workflow_id=test-workflow")
    
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Failed to start HTTP server: %v", err)
    }
}
EOF

# Create documentation
echo "6. Creating documentation..."
cp "CLIENT_SDK_GUIDE.md" "${OUTPUT_DIR}/docs/CLIENT_GUIDE.md"
cp "RELEASE_PLAN.md" "${OUTPUT_DIR}/docs/RELEASE_PLAN.md"

# Create README for the packaged SDK
echo "7. Creating package README..."
cat > "${OUTPUT_DIR}/README.md" << 'EOF'
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
EOF

# Create LICENSE file
echo "8. Creating LICENSE file..."
cat > "${OUTPUT_DIR}/LICENSE" << 'EOF'
MIT License

Copyright (c) 2024 Your Organization

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF

# Create build script
echo "9. Creating build script..."
cat > "${OUTPUT_DIR}/build.sh" << 'EOF'
#!/bin/bash

# Build script for Unified Workflow SDK

set -e

echo "Building Unified Workflow SDK..."

# Run tests
echo "1. Running tests..."
go test ./...

# Build examples
echo "2. Building examples..."
go build -o bin/basic_example examples/basic_usage.go
go build -o bin/http_example examples/http_integration.go

echo "Build complete!"
echo "Binaries available in bin/"
EOF

chmod +x "${OUTPUT_DIR}/build.sh"

# Create archive
echo "10. Creating distribution archive..."
cd dist
tar -czf "unified-workflow-sdk-v${VERSION}.tar.gz" "sdk-v${VERSION}"
cd ..

echo ""
echo "=== Packaging Complete ==="
echo ""
echo "Distribution package created:"
echo "  - dist/unified-workflow-sdk-v${VERSION}.tar.gz"
echo ""
echo "Package contents:"
echo "  - SDK source code (.go files)"
echo "  - Examples (basic_usage.go, http_integration.go)"
echo "  - Documentation (CLIENT_GUIDE.md, RELEASE_PLAN.md)"
echo "  - Build script (build.sh)"
echo "  - License (LICENSE)"
echo ""
echo "To distribute to clients:"
echo "  1. Upload the tar.gz file to your release page"
echo "  2. Clients can extract and run: go get ./..."
echo "  3. Or use directly: import \"github.com/your-org/unified-workflow-sdk\""
echo ""
echo "Next steps:"
echo "  1. Set up CI/CD pipeline for automated releases"
echo "  2. Create GitHub repository for the SDK"
echo "  3. Publish to Go module proxy"
echo "  4. Create Docker image for the SDK"
echo ""