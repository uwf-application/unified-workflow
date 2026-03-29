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
