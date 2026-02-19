package workflows

import (
	"fmt"

	"unified-workflow/internal/common/model"
	"unified-workflow/workflows/steps"
)

// CreateAntifraudTransactionWorkflow creates a complete transaction validation workflow
// using the antifraud SDK at the specified endpoint
func CreateAntifraudTransactionWorkflow(endpoint string) model.Workflow {
	workflow := model.NewBaseWorkflow(
		"antifraud-transaction-validation",
		fmt.Sprintf("Complete transaction validation using antifraud SDK at %s", endpoint),
	)

	fmt.Printf("Creating antifraud transaction workflow for endpoint: %s\n", endpoint)

	// Step 1: Store Transaction
	storeStep := steps.NewStoreTransactionStep(endpoint)
	workflow.AddStep(storeStep)
	fmt.Println("✓ Added StoreTransactionStep")

	// Step 2: AML Validation
	amlStep := steps.NewAMLValidationStep(endpoint)
	workflow.AddStep(amlStep)
	fmt.Println("✓ Added AMLValidationStep")

	// Note: FCValidationStep, MLValidationStep, and FinalizeTransactionStep
	// would be implemented similarly to AMLValidationStep
	// For now, we'll use sequential steps as placeholders

	// Step 3: FC Validation (Fraud Check) - placeholder
	fcStep := model.NewSequentialStep("fc-validation-placeholder")
	workflow.AddStep(fcStep)
	fmt.Println("✓ Added FCValidationStep (placeholder)")

	// Step 4: ML Validation (Machine Learning) - placeholder
	mlStep := model.NewSequentialStep("ml-validation-placeholder")
	workflow.AddStep(mlStep)
	fmt.Println("✓ Added MLValidationStep (placeholder)")

	// Step 5: Finalize Transaction - placeholder
	finalizeStep := model.NewSequentialStep("finalize-transaction-placeholder")
	workflow.AddStep(finalizeStep)
	fmt.Println("✓ Added FinalizeTransactionStep (placeholder)")

	fmt.Printf("✓ Antifraud workflow created with %d steps\n", workflow.GetStepCount())
	fmt.Printf("✓ Total child steps: %d\n", workflow.GetTotalChildStepCount())

	return workflow
}

// GetAntifraudWorkflows returns a list of antifraud workflows
func GetAntifraudWorkflows() []model.Workflow {
	endpoints := []string{
		"af-test.qazpost.kz",         // Primary endpoint
		"https://af-test.qazpost.kz", // HTTPS endpoint
	}

	var workflows []model.Workflow
	for _, endpoint := range endpoints {
		workflow := CreateAntifraudTransactionWorkflow(endpoint)
		workflows = append(workflows, workflow)
	}

	return workflows
}

// GetAntifraudExampleWorkflows returns all example workflows including antifraud workflows
func GetAntifraudExampleWorkflows() []model.Workflow {
	// Get existing example workflows
	existingWorkflows := []model.Workflow{
		createEchoWorkflow(),
		createSequentialWorkflow(),
		createPaymentProcessingWorkflow(),
		createMultiStepWorkflow(),
	}

	// Add antifraud workflows
	antifraudWorkflows := GetAntifraudWorkflows()

	// Combine all workflows
	allWorkflows := append(existingWorkflows, antifraudWorkflows...)

	fmt.Printf("Total workflows available: %d (including %d antifraud workflows)\n",
		len(allWorkflows), len(antifraudWorkflows))

	return allWorkflows
}

// Helper function to print workflow summary
func PrintWorkflowSummary(workflow model.Workflow) {
	fmt.Printf("\n=== Workflow Summary ===\n")
	fmt.Printf("Name: %s\n", workflow.GetName())
	fmt.Printf("Description: %s\n", workflow.GetDescription())
	fmt.Printf("ID: %s\n", workflow.GetID())
	fmt.Printf("Total Steps: %d\n", workflow.GetStepCount())

	// Type assert to BaseWorkflow to access GetTotalChildStepCount
	if baseWorkflow, ok := workflow.(*model.BaseWorkflow); ok {
		fmt.Printf("Total Child Steps: %d\n", baseWorkflow.GetTotalChildStepCount())
	}

	fmt.Printf("\nSteps:\n")
	for i, step := range workflow.GetSteps() {
		fmt.Printf("  %d. %s (%d child steps)\n", i+1, step.GetName(), step.GetChildStepCount())
	}
	fmt.Println()
}

// TestAntifraudWorkflow demonstrates how to use the antifraud workflow
func TestAntifraudWorkflow() {
	fmt.Println("=== Testing Antifraud Transaction Workflow ===")

	// Create workflow for test endpoint
	endpoint := "af-test.qazpost.kz"
	workflow := CreateAntifraudTransactionWorkflow(endpoint)

	// Print summary
	PrintWorkflowSummary(workflow)

	// Simulate workflow execution
	fmt.Println("Simulating workflow execution...")

	// In a real implementation, you would execute the workflow through the executor
	// For now, just show the structure
	fmt.Println("Workflow ready for execution!")
	fmt.Println("Clients can use this workflow to validate transactions with antifraud SDK")
	fmt.Println()

	// Show example transaction data
	fmt.Println("Example transaction data format:")
	exampleTransaction := map[string]interface{}{
		"transaction": map[string]interface{}{
			"id":                   "txn-12345",
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
		},
	}

	fmt.Printf("%+v\n", exampleTransaction)
}
