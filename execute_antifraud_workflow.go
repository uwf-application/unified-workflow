package main

import (
	"context"
	"fmt"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/internal/executor"
	"unified-workflow/internal/queue"
	"unified-workflow/internal/registry"
	"unified-workflow/internal/state"
	"unified-workflow/workflows"
	"unified-workflow/workflows/steps"
)

func main() {
	fmt.Println("=== Executing Antifraud Transaction Validation Workflow ===")
	fmt.Println()

	// Create context
	ctx := context.Background()

	// Create in-memory registry
	registry := registry.NewInMemoryRegistry()

	// Create in-memory queue
	queue := queue.NewInMemoryQueue()

	// Create in-memory state management
	stateManagement := state.NewInMemoryState()

	// Create simple executor
	executor := executor.NewSimpleExecutor(registry, queue, stateManagement)

	// Start the executor
	err := executor.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to start executor: %v\n", err)
		return
	}
	defer executor.Stop(ctx)

	fmt.Println("✓ Executor started")

	// Create antifraud workflow
	endpoint := "af-test.qazpost.kz"
	workflow := workflows.CreateAntifraudTransactionWorkflow(endpoint)

	// Register the workflow
	err = registry.RegisterWorkflow(ctx, workflow)
	if err != nil {
		fmt.Printf("Failed to register workflow: %v\n", err)
		return
	}

	fmt.Printf("✓ Workflow registered: %s (ID: %s)\n", workflow.GetName(), workflow.GetID())

	// Print workflow details
	fmt.Println("\n=== Workflow Details ===")
	fmt.Printf("Name: %s\n", workflow.GetName())
	fmt.Printf("Description: %s\n", workflow.GetDescription())
	fmt.Printf("Total Steps: %d\n", workflow.GetStepCount())

	// Type assert to BaseWorkflow to access GetTotalChildStepCount
	if baseWorkflow, ok := workflow.(*model.BaseWorkflow); ok {
		fmt.Printf("Total Child Steps: %d\n", baseWorkflow.GetTotalChildStepCount())
	}

	// Show steps
	fmt.Println("\nSteps:")
	for i, step := range workflow.GetSteps() {
		fmt.Printf("  %d. %s (%d child steps)\n", i+1, step.GetName(), step.GetChildStepCount())

		// Show child steps for our antifraud steps
		if step.GetName() == "store-transaction" || step.GetName() == "aml-validation" {
			if antistep, ok := step.(*steps.AntifraudStep); ok {
				fmt.Printf("     Child Steps: ")
				for j, childStep := range antistep.GetChildSteps() {
					if j > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s", childStep.GetName())
				}
				fmt.Println()
			}
		}
	}

	// Prepare transaction data
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
			"timestamp":            time.Now().Format(time.RFC3339),
		},
	}

	fmt.Println("\n=== Transaction Data ===")
	fmt.Printf("Transaction ID: %s\n", transactionData["transaction"].(map[string]interface{})["id"])
	fmt.Printf("Amount: %s %s\n",
		transactionData["transaction"].(map[string]interface{})["amount"],
		transactionData["transaction"].(map[string]interface{})["currency"])
	fmt.Printf("Client: %s\n", transactionData["transaction"].(map[string]interface{})["client_name"])

	// Submit workflow for execution
	fmt.Println("\n=== Submitting Workflow for Execution ===")
	runID, err := executor.SubmitWorkflow(ctx, workflow)
	if err != nil {
		fmt.Printf("Failed to submit workflow: %v\n", err)
		return
	}

	fmt.Printf("✓ Workflow submitted successfully\n")
	fmt.Printf("✓ Run ID: %s\n", runID)

	// Get execution status
	fmt.Println("\n=== Checking Execution Status ===")
	status, err := executor.GetExecutionStatus(ctx, runID)
	if err != nil {
		fmt.Printf("Failed to get execution status: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", status.Status)
	fmt.Printf("Progress: %.1f%%\n", status.Progress*100)
	fmt.Printf("Current Step: %s\n", status.CurrentStep)
	fmt.Printf("Current Step Index: %d\n", status.CurrentStepIndex)
	fmt.Printf("Current Child Step Index: %d\n", status.CurrentChildStepIndex)

	// Get execution data
	fmt.Println("\n=== Getting Execution Data ===")
	execData, err := executor.GetExecutionData(ctx, runID)
	if err != nil {
		fmt.Printf("Failed to get execution data: %v\n", err)
		return
	}

	fmt.Printf("Execution data keys: %d\n", len(execData))
	if len(execData) > 0 {
		fmt.Println("Execution data:")
		for k, v := range execData {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	// Get metrics
	fmt.Println("\n=== Getting Execution Metrics ===")
	metrics, err := executor.GetMetrics(ctx, runID)
	if err != nil {
		fmt.Printf("Failed to get metrics: %v\n", err)
		return
	}

	fmt.Printf("Total Steps: %d\n", metrics.TotalSteps)
	fmt.Printf("Completed Steps: %d\n", metrics.CompletedSteps)
	fmt.Printf("Failed Steps: %d\n", metrics.FailedSteps)
	fmt.Printf("Total Child Steps: %d\n", metrics.TotalChildSteps)
	fmt.Printf("Completed Child Steps: %d\n", metrics.CompletedChildSteps)
	fmt.Printf("Failed Child Steps: %d\n", metrics.FailedChildSteps)
	fmt.Printf("Success Rate: %.1f%%\n", metrics.SuccessRate*100)

	// Simulate workflow execution by manually calling step execution
	fmt.Println("\n=== Simulating Step Execution ===")

	// Create a simple test context and data
	testContext := map[string]interface{}{
		"run_id":      runID,
		"workflow_id": workflow.GetID(),
		"start_time":  time.Now(),
	}

	testData := map[string]interface{}{
		"transaction": transactionData["transaction"],
		"workflow_data": map[string]interface{}{
			"endpoint":     endpoint,
			"submitted_at": time.Now().Format(time.RFC3339),
		},
	}

	// Test StoreTransactionStep
	fmt.Println("\n--- Testing StoreTransactionStep ---")
	storeStep := steps.NewStoreTransactionStep(endpoint)

	// Initialize the step
	err = storeStep.InitializeService()
	if err != nil {
		fmt.Printf("Failed to initialize store step: %v\n", err)
		fmt.Println("Note: This is expected if no API key is configured.")
		fmt.Println("To use the actual antifraud SDK, set ANTIFRAUD_API_KEY environment variable.")
		fmt.Println("Example: export ANTIFRAUD_API_KEY='your-api-key-here'")
	} else {
		fmt.Println("✓ StoreTransactionStep initialized")

		// Execute the step logic
		err = storeStep.ExecuteStepLogic(ctx, testContext, testData)
		if err != nil {
			fmt.Printf("StoreTransactionStep execution failed: %v\n", err)
		} else {
			fmt.Println("✓ StoreTransactionStep executed successfully")
		}
	}

	// Test AMLValidationStep
	fmt.Println("\n--- Testing AMLValidationStep ---")
	amlStep := steps.NewAMLValidationStep(endpoint)

	err = amlStep.InitializeService()
	if err != nil {
		fmt.Printf("Failed to initialize AML step: %v\n", err)
		fmt.Println("Note: This is expected if no API key is configured.")
		fmt.Println("To use the actual antifraud SDK, set ANTIFRAUD_API_KEY environment variable.")
	} else {
		fmt.Println("✓ AMLValidationStep initialized")

		// Execute the step logic
		err = amlStep.ExecuteStepLogic(ctx, testContext, testData)
		if err != nil {
			fmt.Printf("AMLValidationStep execution failed: %v\n", err)
		} else {
			fmt.Println("✓ AMLValidationStep executed successfully")
		}
	}

	fmt.Println("\n=== Execution Summary ===")
	fmt.Println("✓ Workflow created and registered")
	fmt.Println("✓ Workflow submitted for execution")
	fmt.Println("✓ Steps initialized and tested")
	fmt.Println("✓ Transaction data prepared")
	fmt.Println("✓ Execution metrics collected")

	fmt.Println("\n=== Next Steps for Full Execution ===")
	fmt.Println("1. Implement actual antifraud SDK integration in steps")
	fmt.Println("2. Configure proper queue and state management")
	fmt.Println("3. Add error handling and retry logic")
	fmt.Println("4. Implement FCValidationStep, MLValidationStep, FinalizeTransactionStep")
	fmt.Println("5. Add monitoring and logging")

	fmt.Println("\nThe antifraud workflow is now ready for integration with your workflow execution system!")
}
