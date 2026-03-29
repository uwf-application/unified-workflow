# Registry Client

A Go client for interacting with the Unified Workflow Registry service. This client provides a comprehensive interface for managing workflow definitions, including creating, reading, updating, deleting, and listing workflows.

## Features

- **Workflow Management**: Create, read, update, and delete workflow definitions
- **Workflow Listing**: List workflows with filtering and pagination
- **Workflow Discovery**: Search and discover available workflows
- **Workflow Validation**: Validate workflow definitions before registration
- **Namespace Support**: Organize workflows by namespace
- **Metadata Management**: Attach and manage workflow metadata
- **Batch Operations**: Perform operations on multiple workflows
- **Error Handling**: Comprehensive error handling with retry mechanisms

## Installation

```bash
go get unified-workflow/pkg/client/registry
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "unified-workflow/pkg/client/registry"
)

func main() {
    // Create registry client configuration
    config := registry.DefaultConfig()
    config.Endpoint = "http://localhost:8080"
    config.Timeout = 30 * time.Second
    config.AuthToken = "your-auth-token"

    // Create registry client
    client := registry.NewClient(config)
    defer client.Close()

    // List all workflows
    ctx := context.Background()
    listReq := &registry.ListWorkflowsRequest{
        Request: client.Request{
            ID:        "req-123",
            Timestamp: time.Now(),
        },
        Limit: 10,
    }

    listResp, err := client.ListWorkflows(ctx, listReq)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Total workflows: %d\n", listResp.TotalCount)
    for _, workflow := range listResp.Workflows {
        fmt.Printf("Workflow: %s (%s) - %d steps\n", 
            workflow.Name, workflow.ID, workflow.StepCount)
    }
}
```

## API Reference

### Core Methods

#### `ListWorkflows(ctx context.Context, req *ListWorkflowsRequest) (*ListWorkflowsResponse, error)`
List workflows with optional filtering and pagination.

#### `GetWorkflow(ctx context.Context, req *GetWorkflowRequest) (*GetWorkflowResponse, error)`
Get detailed information about a specific workflow.

#### `CreateWorkflow(ctx context.Context, req *CreateWorkflowRequest) (*CreateWorkflowResponse, error)`
Create a new workflow definition.

#### `UpdateWorkflow(ctx context.Context, req *UpdateWorkflowRequest) (*UpdateWorkflowResponse, error)`
Update an existing workflow definition.

#### `DeleteWorkflow(ctx context.Context, req *DeleteWorkflowRequest) (*DeleteWorkflowResponse, error)`
Delete a workflow definition.

#### `ContainsWorkflow(ctx context.Context, req *ContainsWorkflowRequest) (*ContainsWorkflowResponse, error)`
Check if a workflow exists.

#### `GetWorkflowCount(ctx context.Context, req *GetWorkflowCountRequest) (*GetWorkflowCountResponse, error)`
Get the total number of workflows.

#### `Clear(ctx context.Context, req *ClearRequest) (*ClearResponse, error)`
Remove all workflows from the registry.

### Request Types

#### `ListWorkflowsRequest`
```go
type ListWorkflowsRequest struct {
    client.Request
    NameFilter        string // Filter by workflow name
    DescriptionFilter string // Filter by workflow description
    Limit             int    // Limit number of results
    Offset            int    // Pagination offset
}
```

#### `GetWorkflowRequest`
```go
type GetWorkflowRequest struct {
    client.Request
    ID string // Workflow ID
}
```

#### `CreateWorkflowRequest`
```go
type CreateWorkflowRequest struct {
    client.Request
    Name        string                 // Workflow name
    Description string                 // Workflow description
    Steps       []model.Step           // Workflow steps
    Metadata    map[string]interface{} // Additional metadata
}
```

#### `UpdateWorkflowRequest`
```go
type UpdateWorkflowRequest struct {
    client.Request
    ID          string                 // Workflow ID
    Name        *string                // Updated name (optional)
    Description *string                // Updated description (optional)
    Steps       []model.Step           // Updated steps (optional)
    Metadata    map[string]interface{} // Updated metadata (optional)
}
```

## Configuration

### Registry Client Configuration

```go
type Config struct {
    client.Config
    
    // Namespace is the registry namespace (optional)
    Namespace string
}
```

### Default Configuration

```go
func DefaultConfig() Config {
    return Config{
        Config: client.DefaultConfig(),
    }
}
```

## Examples

### Creating a Workflow

```go
func createWorkflow(client registry.Client) error {
    ctx := context.Background()
    
    // Define workflow steps
    steps := []model.Step{
        {
            Name: "validate-input",
            Type: "sequential",
            ChildSteps: []model.ChildStep{
                {
                    Name:         "validate-amount",
                    PrimitiveName: "validation",
                    Parameters: map[string]interface{}{
                        "field": "amount",
                        "min":   0.01,
                        "max":   10000,
                    },
                },
            },
        },
        {
            Name: "process-payment",
            Type: "sequential",
            ChildSteps: []model.ChildStep{
                {
                    Name:         "charge-card",
                    PrimitiveName: "payment",
                    Parameters: map[string]interface{}{
                        "provider": "stripe",
                    },
                },
            },
        },
    }
    
    // Create workflow request
    req := &registry.CreateWorkflowRequest{
        Name:        "payment-processing",
        Description: "Process payment transactions",
        Steps:       steps,
        Metadata: map[string]interface{}{
            "category": "finance",
            "version":  "1.0.0",
        },
    }
    
    resp, err := client.CreateWorkflow(ctx, req)
    if err != nil {
        return err
    }
    
    fmt.Printf("Workflow created: %s\n", resp.Workflow.ID)
    fmt.Printf("Workflow name: %s\n", resp.Workflow.Name)
    fmt.Printf("Step count: %d\n", resp.Workflow.StepCount)
    
    return nil
}
```

### Getting Workflow Details

```go
func getWorkflowDetails(client registry.Client, workflowID string) error {
    ctx := context.Background()
    
    req := &registry.GetWorkflowRequest{
        ID: workflowID,
    }
    
    resp, err := client.GetWorkflow(ctx, req)
    if err != nil {
        return err
    }
    
    workflow := resp.Workflow
    fmt.Printf("Workflow Details:\n")
    fmt.Printf("  ID: %s\n", workflow.ID)
    fmt.Printf("  Name: %s\n", workflow.Name)
    fmt.Printf("  Description: %s\n", workflow.Description)
    fmt.Printf("  Created: %v\n", workflow.CreatedAt.Format(time.RFC3339))
    fmt.Printf("  Updated: %v\n", workflow.UpdatedAt.Format(time.RFC3339))
    fmt.Printf("  Step count: %d\n", workflow.StepCount)
    
    if len(workflow.Steps) > 0 {
        fmt.Printf("  Steps:\n")
        for i, step := range workflow.Steps {
            fmt.Printf("    %d. %s (%s)\n", i+1, step.Name, step.Type)
            fmt.Printf("      Child steps: %d\n", step.ChildStepCount)
            fmt.Printf("      Parallel: %v\n", step.IsParallel)
        }
    }
    
    return nil
}
```

### Updating a Workflow

```go
func updateWorkflow(client registry.Client, workflowID string) error {
    ctx := context.Background()
    
    // Update workflow description and metadata
    description := "Updated payment processing workflow with fraud detection"
    req := &registry.UpdateWorkflowRequest{
        ID:          workflowID,
        Description: &description,
        Metadata: map[string]interface{}{
            "category":    "finance",
            "version":     "1.1.0",
            "features":    []string{"fraud-detection", "retry-logic"},
            "last_updated": time.Now().Format(time.RFC3339),
        },
    }
    
    resp, err := client.UpdateWorkflow(ctx, req)
    if err != nil {
        return err
    }
    
    fmt.Printf("Workflow updated: %s\n", resp.Workflow.ID)
    fmt.Printf("New description: %s\n", resp.Workflow.Description)
    fmt.Printf("Updated at: %v\n", resp.Workflow.UpdatedAt.Format(time.RFC3339))
    
    return nil
}
```

### Listing Workflows with Filtering

```go
func listWorkflowsWithFilters(client registry.Client) error {
    ctx := context.Background()
    
    // List workflows with filters
    req := &registry.ListWorkflowsRequest{
        NameFilter: "payment", // Filter by name containing "payment"
        Limit:      20,
        Offset:     0,
    }
    
    resp, err := client.ListWorkflows(ctx, req)
    if err != nil {
        return err
    }
    
    fmt.Printf("Total workflows: %d\n", resp.TotalCount)
    fmt.Printf("Filtered workflows: %d\n", resp.FilteredCount)
    fmt.Printf("Showing %d workflows:\n", len(resp.Workflows))
    
    for i, workflow := range resp.Workflows {
        fmt.Printf("%d. %s (%s)\n", i+1, workflow.Name, workflow.ID)
        fmt.Printf("   Description: %s\n", workflow.Description)
        fmt.Printf("   Steps: %d, Created: %v\n", 
            workflow.StepCount, workflow.CreatedAt.Format("2006-01-02"))
    }
    
    return nil
}
```

### Checking Workflow Existence

```go
func checkWorkflowExists(client registry.Client, workflowID string) error {
    ctx := context.Background()
    
    req := &registry.ContainsWorkflowRequest{
        ID: workflowID,
    }
    
    resp, err := client.ContainsWorkflow(ctx, req)
    if err != nil {
        return err
    }
    
    if resp.Exists {
        fmt.Printf("Workflow %s exists\n", workflowID)
    } else {
        fmt.Printf("Workflow %s does not exist\n", workflowID)
    }
    
    return nil
}
```

### Deleting a Workflow

```go
func deleteWorkflow(client registry.Client, workflowID string) error {
    ctx := context.Background()
    
    // First check if workflow exists
    containsReq := &registry.ContainsWorkflowRequest{
        ID: workflowID,
    }
    
    containsResp, err := client.ContainsWorkflow(ctx, containsReq)
    if err != nil {
        return err
    }
    
    if !containsResp.Exists {
        fmt.Printf("Workflow %s does not exist\n", workflowID)
        return nil
    }
    
    // Delete the workflow
    deleteReq := &registry.DeleteWorkflowRequest{
        ID: workflowID,
    }
    
    deleteResp, err := client.DeleteWorkflow(ctx, deleteReq)
    if err != nil {
        return err
    }
    
    if deleteResp.Deleted {
        fmt.Printf("Workflow %s deleted successfully\n", workflowID)
    } else {
        fmt.Printf("Failed to delete workflow %s\n", workflowID)
    }
    
    return nil
}
```

### Getting Workflow Count

```go
func getWorkflowCount(client registry.Client) error {
    ctx := context.Background()
    
    req := &registry.GetWorkflowCountRequest{}
    
    resp, err := client.GetWorkflowCount(ctx, req)
    if err != nil {
        return err
    }
    
    fmt.Printf("Total workflows in registry: %d\n", resp.Count)
    
    return nil
}
```

### Clearing All Workflows

```go
func clearAllWorkflows(client registry.Client) error {
    ctx := context.Background()
    
    // Get count before clearing
    countReq := &registry.GetWorkflowCountRequest{}
    countResp, err := client.GetWorkflowCount(ctx, countReq)
    if err != nil {
        return err
    }
    
    if countResp.Count == 0 {
        fmt.Println("Registry is already empty")
        return nil
    }
    
    // Ask for confirmation (in real application)
    fmt.Printf("About to delete %d workflows. Continue? (y/n): ", countResp.Count)
    // ... read user input ...
    
    // Clear all workflows
    clearReq := &registry.ClearRequest{}
    clearResp, err := client.Clear(ctx, clearReq)
    if err != nil {
        return err
    }
    
    fmt.Printf("Cleared %d workflows from registry\n", clearResp.ClearedCount)
    
    return nil
}
```

## Best Practices

### 1. Validate Workflow Before Creation
```go
func validateAndCreateWorkflow(client registry.Client, workflow *model.Workflow) error {
    // Validate workflow structure
    if workflow.Name == "" {
        return fmt.Errorf("workflow name is required")
    }
    
    if len(workflow.Steps) == 0 {
        return fmt.Errorf("workflow must have at least one step")
    }
    
    // Validate each step
    for i, step := range workflow.Steps {
        if step.Name == "" {
            return fmt.Errorf("step %d: name is required", i)
        }
        // ... additional validation ...
    }
    
    // Create workflow
    req := &registry.CreateWorkflowRequest{
        Name:        workflow.Name,
        Description: workflow.Description,
        Steps:       workflow.Steps,
        Metadata:    workflow.Metadata,
    }
    
    _, err := client.CreateWorkflow(ctx, req)
    return err
}
```

### 2. Use Pagination for Large Workflow Lists
```go
func listAllWorkflows(client registry.Client) ([]*registry.WorkflowInfo, error) {
    ctx := context.Background()
    var allWorkflows []*registry.WorkflowInfo
    limit := 50
    offset := 0
    
    for {
        req := &registry.ListWorkflowsRequest{
            Limit:  limit,
            Offset: offset,
        }
        
        resp, err := client.ListWorkflows(ctx, req)
        if err != nil {
            return nil, err
        }
        
        allWorkflows = append(allWorkflows, resp.Workflows...)
        
        if offset+len(resp.Workflows) >= resp.TotalCount {
            break
        }
        
        offset += limit
    }
    
    return allWorkflows, nil
}
```

### 3. Cache Workflow Details
```go
type WorkflowCache struct {
    client    registry.Client
    cache     map[string]*registry.WorkflowDetail
    cacheTTL  time.Duration
    lastFetch map[string]time.Time
}

func (c *WorkflowCache) GetWorkflow(workflowID string) (*registry.WorkflowDetail, error) {
    // Check cache
    if detail, ok := c.cache[workflowID]; ok {
        if time.Since(c.lastFetch[workflowID]) < c.cacheTTL {
            return detail, nil
        }
    }
    
    // Fetch from registry
    req := &registry.GetWorkflowRequest{
        ID: workflowID,
    }
    
    resp, err := c.client.GetWorkflow(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Update cache
    c.cache[workflowID] = resp.Workflow
    c.lastFetch[workflowID] = time.Now()
    
    return resp.Workflow, nil
}
```

### 4. Handle Concurrent Updates
```go
func updateWorkflowSafely(client registry.Client, workflowID string, updates map[string]interface{}) error {
    ctx := context.Background()
    
    // Get current workflow
    getReq := &registry.GetWorkflowRequest{
        ID: workflowID,
    }
    
    getResp, err := client.GetWorkflow(ctx, getReq)
    if err != nil {
        return err
    }
    
    currentWorkflow := getResp.Workflow
    
    // Apply updates
    updateReq := &registry.UpdateWorkflowRequest{
        ID: workflowID,
    }
    
    // Only update if workflow hasn't changed since we read it
    if updates["name"] != nil {
        name := updates["name"].(string)
        updateReq.Name = &name
    }
    
    if updates["description"] != nil {
        description := updates["description"].(string)
        updateReq.Description = &description
    }
    
    // Include original metadata and merge updates
    metadata := make(map[string]interface{})
    for k, v := range currentWorkflow.Metadata {
        metadata[k] = v
    }
    for k, v := range updates["metadata"].(map[string]interface{}) {
        metadata[k] = v
    }
    updateReq.Metadata = metadata
    
    // Perform update
    _, err = client.UpdateWorkflow(ctx, updateReq)
    return err
}
```

### 5. Monitor Registry Health
```go
func monitorRegistryHealth(client registry.Client) error {
    ctx := context.Background()
    
    // Ping registry
    if err := client.Ping(ctx); err != nil {
        return fmt.Errorf("registry ping failed: %v", err)
    }
    
    // Get workflow count
    countReq := &registry.GetWorkflowCountRequest{}
    countResp, err := client.GetWorkflowCount(ctx, countReq)
    if err != nil {
        return fmt.Errorf("failed to get workflow count: %v", err)
    }
    
    fmt.Printf("Registry health check passed\n")
    fmt.Printf("  Status: healthy\n")
    fmt.Printf("  Workflow count: %d\n", countResp.Count)
    
    return nil
}
```

## Integration with Executor Client

The registry client is often used together with the executor client:

```go
import (
    "context"
    "fmt"
    
    "unified-workflow/pkg/client/registry"
    "unified-workflow/pkg/client/executor"
)

func executeWorkflowByName(registryClient registry.Client, executorClient executor.Client, workflowName string, inputData map[string]interface{}) error {
    ctx := context.Background()
    
    // First, find workflow by name