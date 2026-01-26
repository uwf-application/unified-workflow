package examples

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"unified-workflow/internal/common/model"
	"unified-workflow/internal/registry"
)

func RunSimpleWorkflowExample() {
	// Create a simple workflow
	workflow := model.NewBaseWorkflow("Simple Data Processing", "A workflow that processes data in multiple steps")

	// Create steps
	step1 := model.NewBaseStep("Data Extraction", false)
	step2 := model.NewBaseStep("Data Transformation", false)
	step3 := model.NewBaseStep("Data Loading", false)

	// Add steps to workflow
	workflow.AddStep(step1)
	workflow.AddStep(step2)
	workflow.AddStep(step3)

	// Create registry
	reg := registry.NewInMemoryRegistry()
	ctx := context.Background()

	// Register workflow
	err := reg.RegisterWorkflow(ctx, workflow)
	if err != nil {
		log.Fatalf("Failed to register workflow: %v", err)
	}

	fmt.Printf("Workflow created:\n")
	fmt.Printf("  ID: %s\n", workflow.GetID())
	fmt.Printf("  Name: %s\n", workflow.GetName())
	fmt.Printf("  Description: %s\n", workflow.GetDescription())
	fmt.Printf("  Steps: %d\n", workflow.GetStepCount())

	// Test API endpoints
	fmt.Println("\nTesting API endpoints:")
	fmt.Println("1. GET /api/v1/workflows - List all workflows")
	fmt.Println("2. GET /api/v1/workflows/{id} - Get workflow details")
	fmt.Println("3. POST /api/v1/workflows/{id}/execute - Execute workflow")
	fmt.Println("4. GET /api/v1/executions/{runId} - Get execution status")

	// Start a simple HTTP server for demonstration
	fmt.Println("\nStarting demo server on :8080...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Workflow API Demo\n")
		fmt.Fprintf(w, "Workflow ID: %s\n", workflow.GetID())
		fmt.Fprintf(w, "Workflow Name: %s\n", workflow.GetName())
		fmt.Fprintf(w, "Total Steps: %d\n", workflow.GetStepCount())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
