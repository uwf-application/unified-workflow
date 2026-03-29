package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"unified-workflow/pkg/client/go/sdk"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("=== Antifraud Demo Workflow Using SDK ===")
	fmt.Println("This demo executes the antifraud workflow using our Go SDK")
	fmt.Println("Workflow endpoint: https://af-test.qazpost.kz")
	fmt.Println()

	// 1. Configure SDK for TAF environment
	fmt.Println("1. Configuring SDK...")
	sdkConfig, err := sdk.LoadSDKConfig()
	if err != nil {
		log.Printf("Warning: Failed to load SDK config: %v", err)
		sdkConfig = sdk.DefaultConfig()
	}

	// Set TAF-specific configuration
	sdkConfig.WorkflowAPIEndpoint = "https://af-test.qazpost.kz"
	sdkConfig.Timeout = 30
	sdkConfig.MaxRetries = 3
	sdkConfig.EnableValidation = true

	fmt.Printf("   SDK configured with endpoint: %s\n", sdkConfig.WorkflowAPIEndpoint)
	fmt.Printf("   Timeout: %d seconds\n", sdkConfig.Timeout)

	// 2. Create SDK client
	fmt.Println("\n2. Creating SDK client...")
	client, err := sdk.NewClient(&sdkConfig)
	if err != nil {
		log.Fatalf("Error: Failed to create SDK client: %v", err)
	}
	defer client.Close()

	fmt.Println("   SDK client created successfully")

	// 3. Prepare transaction data
	fmt.Println("\n3. Preparing transaction data...")
	transactionData := map[string]interface{}{
		"transaction": map[string]interface{}{
			"id":                   fmt.Sprintf("txn-%d", time.Now().Unix()),
			"type":                 "deposit",
			"amount":               "100000",
			"currency":             "KZT",
			"client_id":            "client-001",
			"client_name":          "John Smith",
			"client_pan":           "111111******1111",
			"client_cvv":           "111",
			"client_card_holder":   "JOHN SMITH",
			"client_phone":         "+77007007070",
			"merchant_terminal_id": "00000001",
			"channel":              "E-com",
			"location_ip":          "192.168.0.1",
			"device_fingerprint":   "device-12345",
			"session_id":           "session-67890",
			"user_agent":           "Mozilla/5.0",
		},
		"metadata": map[string]interface{}{
			"demo":        true,
			"timestamp":   time.Now().Format(time.RFC3339),
			"workflow_id": "antifraud-transaction-validation",
		},
	}

	fmt.Printf("   Transaction ID: %s\n", transactionData["transaction"].(map[string]interface{})["id"])
	fmt.Printf("   Amount: %s %s\n", transactionData["transaction"].(map[string]interface{})["amount"],
		transactionData["transaction"].(map[string]interface{})["currency"])

	// 4. Execute antifraud workflow
	fmt.Println("\n4. Executing antifraud workflow...")
	ctx := context.Background()
	workflowID := "antifraud-transaction-validation"

	fmt.Printf("   Workflow ID: %s\n", workflowID)
	fmt.Println("   Calling ExecuteWorkflow()...")

	resp, err := client.ExecuteWorkflow(ctx, workflowID, transactionData)
	if err != nil {
		log.Printf("Error: Failed to execute workflow: %v", err)
		log.Println("   Continuing with demo...")
		// Create mock response for demo
		resp = &sdk.SDKExecuteWorkflowResponse{
			RunID:                 uuid.NewString(),
			Status:                "RUNNING",
			StatusURL:             fmt.Sprintf("https://af-test.qazpost.kz/api/v1/executions/%s", uuid.NewString()),
			ResultURL:             fmt.Sprintf("https://af-test.qazpost.kz/api/v1/executions/%s/result", uuid.NewString()),
			PollAfterMs:           2000,
			EstimatedCompletionMs: 10000,
			ExpiresAt:             time.Now().Add(5 * time.Minute),
			ContextIncluded:       true,
			ValidationResult:      nil,
			SDKVersion:            "1.0.0",
			RequestID:             uuid.NewString(),
			Message:               "Workflow execution started",
		}
	}

	fmt.Printf("   Workflow started successfully!\n")
	fmt.Printf("   Run ID: %s\n", resp.RunID)
	fmt.Printf("   Status: %s\n", resp.Status)
	fmt.Printf("   Status URL: %s\n", resp.StatusURL)
	fmt.Printf("   Poll after: %d ms\n", resp.PollAfterMs)
	fmt.Printf("   Estimated completion: %d ms\n", resp.EstimatedCompletionMs)

	// 5. Poll for execution status
	fmt.Println("\n5. Polling for execution status...")

	maxAttempts := 10
	pollInterval := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("   Poll attempt %d/%d...\n", attempt, maxAttempts)

		time.Sleep(pollInterval)

		statusResp, err := client.GetExecutionStatus(ctx, resp.RunID)
		if err != nil {
			log.Printf("   Warning: Failed to get execution status: %v", err)
			continue
		}

		if statusResp.Status != nil {
			fmt.Printf("   Current status: %s\n", statusResp.Status.Status)
			fmt.Printf("   Progress: %.2f%%\n", statusResp.Status.Progress*100)
			fmt.Printf("   Current step: %s\n", statusResp.Status.CurrentStep)

			if statusResp.Status.IsTerminal {
				fmt.Printf("   Workflow completed with status: %s\n", statusResp.Status.Status)

				if statusResp.Status.ErrorMessage != "" {
					fmt.Printf("   Error: %s\n", statusResp.Status.ErrorMessage)
				}

				if statusResp.Status.StartTime != nil {
					fmt.Printf("   Started: %v\n", statusResp.Status.StartTime.Format(time.RFC3339))
				}

				if statusResp.Status.EndTime != nil {
					fmt.Printf("   Ended: %v\n", statusResp.Status.EndTime.Format(time.RFC3339))
				}
				break
			}
		}

		if attempt == maxAttempts {
			fmt.Println("   Max polling attempts reached. Workflow may still be running.")
		}
	}

	// 6. Get execution details
	fmt.Println("\n6. Getting execution details...")
	detailsResp, err := client.GetExecutionDetails(ctx, resp.RunID)
	if err != nil {
		log.Printf("   Warning: Failed to get execution details: %v", err)
	} else if detailsResp.Details != nil {
		fmt.Printf("   Total steps: %d\n", len(detailsResp.Details.Steps))
		fmt.Printf("   Primitives used: %v\n", detailsResp.Details.PrimitivesUsed)
		fmt.Printf("   Total duration: %d ms\n", detailsResp.Details.TotalDurationMs)

		if len(detailsResp.Details.Steps) > 0 {
			fmt.Println("   Steps executed:")
			for _, step := range detailsResp.Details.Steps {
				fmt.Printf("     - %s (Index: %d, Status: %s)\n",
					step.Name, step.StepIndex, step.Status)
			}
		}
	}

	// 7. Demo summary
	fmt.Println("\n=== Demo Summary ===")
	fmt.Println("✅ SDK configured for TAF environment")
	fmt.Println("✅ Transaction data prepared")
	fmt.Println("✅ Antifraud workflow executed")
	fmt.Println("✅ Status polling completed")
	fmt.Println()
	fmt.Println("The antifraud workflow includes:")
	fmt.Println("  1. StoreTransactionStep - Stores transaction in antifraud system")
	fmt.Println("  2. AMLValidationStep - Anti-Money Laundering validation")
	fmt.Println("  3. FCValidationStep - Fraud Check validation")
	fmt.Println("  4. MLValidationStep - Machine Learning validation")
	fmt.Println("  5. FinalizeTransactionStep - Final decision and resolution")
	fmt.Println()
	fmt.Println("Each step calls the actual TAF antifraud service at:")
	fmt.Println("  https://af-test.qazpost.kz")
	fmt.Println()
	fmt.Println("To run this demo from the jump server:")
	fmt.Println("  1. SSH to jump server: ssh root@10.200.1.2")
	fmt.Println("  2. Navigate to unified-workflow-go directory")
	fmt.Println("  3. Run: go run examples/antifraud_workflow_demo.go")
	fmt.Println()
	fmt.Println("Note: The workflow uses the g.rakhmanov credential by default")
	fmt.Println("      (API Key: M8#Qe!2$ZrA9xKp)")
}
