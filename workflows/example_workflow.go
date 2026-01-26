package workflows

import (
	"unified-workflow/internal/common/model"
	"unified-workflow/workflows/steps"
)

// GetExampleWorkflows returns a list of example workflows that should be available when the API starts
func GetExampleWorkflows() []model.Workflow {
	return []model.Workflow{
		createEchoWorkflow(),
		createSequentialWorkflow(),
		createPaymentProcessingWorkflow(),
		createMultiStepWorkflow(),
	}
}

// createEchoWorkflow creates a simple echo workflow
func createEchoWorkflow() model.Workflow {
	workflow := model.NewBaseWorkflow(
		"Echo Workflow",
		"A simple workflow that echoes input data",
	)

	// Add echo steps
	echoSteps := steps.GetEchoSteps()
	for _, step := range echoSteps {
		workflow.AddStep(step)
	}

	return workflow
}

// createSequentialWorkflow creates a sequential workflow example
func createSequentialWorkflow() model.Workflow {
	workflow := model.NewBaseWorkflow(
		"Sequential Workflow",
		"A workflow with sequential steps",
	)

	// Add sequential steps
	workflow.AddStep(model.NewSequentialStep("step-1"))
	workflow.AddStep(model.NewSequentialStep("step-2"))
	workflow.AddStep(model.NewSequentialStep("step-3"))

	return workflow
}

// createPaymentProcessingWorkflow creates a payment processing workflow example
func createPaymentProcessingWorkflow() model.Workflow {
	workflow := model.NewBaseWorkflow(
		"Payment Processing Workflow",
		"Processes payments with validation and fraud check",
	)

	// Add payment processing steps
	workflow.AddStep(model.NewSequentialStep("validate-payment"))
	workflow.AddStep(model.NewSequentialStep("check-fraud"))
	workflow.AddStep(model.NewSequentialStep("process-transaction"))
	workflow.AddStep(model.NewSequentialStep("send-receipt"))

	return workflow
}

// createMultiStepWorkflow creates a workflow with multiple step types
func createMultiStepWorkflow() model.Workflow {
	workflow := model.NewBaseWorkflow(
		"Multi-Step Workflow",
		"A workflow demonstrating different step types",
	)

	// Add various step types
	workflow.AddStep(model.NewSequentialStep("initial-step"))

	// Add echo steps
	echoSteps := steps.GetEchoSteps()
	for i, step := range echoSteps {
		if i < 2 { // Add first 2 echo steps
			workflow.AddStep(step)
		}
	}

	workflow.AddStep(model.NewSequentialStep("final-step"))

	return workflow
}
