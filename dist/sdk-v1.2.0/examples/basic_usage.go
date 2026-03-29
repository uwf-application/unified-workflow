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
