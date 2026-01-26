package steps

import (
	"context"
	"fmt"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/workflows/child_steps"
)

// EchoStep is a simple step that echoes input data and processes child steps
type EchoStep struct {
	*model.BaseStep
	Message    string
	ChildSteps []*model.ChildStep
}

// NewEchoStep creates a new EchoStep with child steps
func NewEchoStep(name, message string) *EchoStep {
	baseStep := model.NewBaseStep(name, false)

	// Get example child steps
	exampleChildSteps := child_steps.GetExampleChildSteps()

	// Add child steps to the base step
	for _, childStep := range exampleChildSteps {
		baseStep.AddChildStep(childStep)
	}

	return &EchoStep{
		BaseStep:   baseStep,
		Message:    message,
		ChildSteps: exampleChildSteps,
	}
}

// GetChildSteps returns the child steps for this step
func (s *EchoStep) GetChildSteps() []*model.ChildStep {
	return s.BaseStep.GetChildSteps()
}

// Run executes the echo step and processes all child steps
func (s *EchoStep) Run(ctx context.Context, context interface{}, data interface{}) error {
	fmt.Printf("=== EchoStep [%s] Starting ===\n", s.GetName())
	fmt.Printf("Message: %s\n", s.Message)
	fmt.Printf("Child Steps Count: %d\n", len(s.ChildSteps))

	// Process each child step through full lifecycle
	for i, childStep := range s.ChildSteps {
		fmt.Printf("\n--- Processing Child Step %d: %s ---\n", i+1, childStep.GetName())

		// 1. Request hook
		fmt.Println("1. Executing Request Hook...")
		requestData := childStep.GetRequestHook()(context, data)
		fmt.Printf("   Request Data: %v\n", requestData)

		// 2. Response hook
		fmt.Println("2. Executing Response Hook...")
		responseData := childStep.GetResponseHook()(context, data)
		fmt.Printf("   Response Data: %v\n", responseData)

		// 3. Validate hook
		fmt.Println("3. Executing Validate Hook...")
		err := childStep.GetValidateHook()(responseData)
		if err != nil {
			fmt.Printf("   Validation Failed: %v\n", err)
			return fmt.Errorf("child step %s validation failed: %w", childStep.GetName(), err)
		}
		fmt.Println("   Validation Passed ✓")

		// Store child step metrics
		startTime := time.Now()
		endTime := time.Now()
		s.StoreChildStepMetrics(childStep, context, data, startTime, endTime, "")

		fmt.Printf("--- Child Step %s Completed Successfully ---\n", childStep.GetName())
	}

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("\n=== EchoStep [%s] Completed Successfully ===\n", s.GetName())
	return nil
}

// ExecuteWithTiming implements the Step interface with timing
func (s *EchoStep) ExecuteWithTiming(ctx context.Context, context interface{}, data interface{}, stepLogic func() error) error {
	startTime := time.Now()
	s.StartTime = &startTime

	err := stepLogic()

	endTime := time.Now()
	s.EndTime = &endTime

	return err
}

// ExecuteChildStep executes a specific child step
func (s *EchoStep) ExecuteChildStep(ctx context.Context, childStep *model.ChildStep, context interface{}, data interface{}) error {
	fmt.Printf("Executing child step: %s\n", childStep.GetName())

	// Execute request hook
	requestData := childStep.GetRequestHook()(context, data)
	fmt.Printf("Request data: %v\n", requestData)

	// Execute response hook
	responseData := childStep.GetResponseHook()(context, data)
	fmt.Printf("Response data: %v\n", responseData)

	// Execute validate hook
	return childStep.GetValidateHook()(responseData)
}

// ExecuteChildStepWithTiming executes a child step with timing
func (s *EchoStep) ExecuteChildStepWithTiming(ctx context.Context, childStep *model.ChildStep, context interface{}, data interface{}) error {
	startTime := time.Now()

	err := s.ExecuteChildStep(ctx, childStep, context, data)

	endTime := time.Now()
	s.StoreChildStepMetrics(childStep, context, data, startTime, endTime, "")

	return err
}

// StoreChildStepMetrics stores metrics for a child step execution
func (s *EchoStep) StoreChildStepMetrics(childStep *model.ChildStep, context interface{}, data interface{}, startTime, endTime time.Time, errorMessage string) {
	duration := endTime.Sub(startTime)
	fmt.Printf("Child Step Metrics - %s: Duration=%v, Error=%s\n",
		childStep.GetName(), duration, errorMessage)
}

// GetEchoSteps returns example echo steps with child steps
func GetEchoSteps() []model.Step {
	return []model.Step{
		NewEchoStep("echo-hello", "Hello, World! Processing child steps..."),
		NewEchoStep("echo-test", "This is a test message with child step execution"),
		NewEchoStep("echo-welcome", "Welcome to the workflow system with full lifecycle"),
	}
}
