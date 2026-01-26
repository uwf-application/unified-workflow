package model

import (
	"fmt"
	"time"
)

// Workflow represents a complete workflow definition with steps
// This is a pure data structure - execution logic is handled by a separate executor
// Similar to Java's Workflow class in workflow-common
type Workflow interface {
	// GetID returns the workflow ID
	GetID() string

	// GetName returns the workflow name
	GetName() string

	// GetDescription returns the workflow description
	GetDescription() string

	// GetSteps returns all steps in the workflow
	GetSteps() []Step

	// GetStep returns a step by name
	GetStep(name string) (Step, error)

	// GetStepCount returns the number of steps
	GetStepCount() int

	// AddStep adds a step to the workflow
	AddStep(step Step) Workflow

	// AddSteps adds multiple steps to the workflow
	AddSteps(steps []Step) Workflow

	// SetPrimitives sets the primitives registry for all steps
	SetPrimitives(primitives interface{})

	// GetPrimitives returns the primitives registry
	GetPrimitives() interface{}

	// SetContext sets the workflow context
	SetContext(context interface{})

	// GetContext returns the workflow context
	GetContext() interface{}

	// SetData sets the workflow data
	SetData(data interface{})

	// GetData returns the workflow data
	GetData() interface{}
}

// BaseWorkflow is the base implementation of the Workflow interface
// This is a pure data structure - execution logic is handled by a separate executor
type BaseWorkflow struct {
	ID          string
	Name        string
	Description string
	Steps       []Step
	Primitives  interface{}
	Context     interface{}
	Data        interface{}
}

// NewBaseWorkflow creates a new BaseWorkflow
func NewBaseWorkflow(name, description string) *BaseWorkflow {
	return &BaseWorkflow{
		ID:          generateWorkflowID(),
		Name:        name,
		Description: description,
		Steps:       []Step{},
	}
}

// NewBaseWorkflowWithContext creates a new BaseWorkflow with context and data
func NewBaseWorkflowWithContext(name, description string, context, data interface{}) *BaseWorkflow {
	return &BaseWorkflow{
		ID:          generateWorkflowID(),
		Name:        name,
		Description: description,
		Steps:       []Step{},
		Context:     context,
		Data:        data,
	}
}

// GetID returns the workflow ID
func (w *BaseWorkflow) GetID() string {
	return w.ID
}

// GetName returns the workflow name
func (w *BaseWorkflow) GetName() string {
	return w.Name
}

// GetDescription returns the workflow description
func (w *BaseWorkflow) GetDescription() string {
	return w.Description
}

// GetSteps returns all steps in the workflow
func (w *BaseWorkflow) GetSteps() []Step {
	return w.Steps
}

// GetStep returns a step by name
func (w *BaseWorkflow) GetStep(name string) (Step, error) {
	for _, step := range w.Steps {
		if step.GetName() == name {
			return step, nil
		}
	}
	return nil, fmt.Errorf("step not found: %s", name)
}

// GetStepCount returns the number of steps
func (w *BaseWorkflow) GetStepCount() int {
	return len(w.Steps)
}

// AddStep adds a step to the workflow
func (w *BaseWorkflow) AddStep(step Step) Workflow {
	w.Steps = append(w.Steps, step)
	return w
}

// AddSteps adds multiple steps to the workflow
func (w *BaseWorkflow) AddSteps(steps []Step) Workflow {
	w.Steps = append(w.Steps, steps...)
	return w
}

// SetPrimitives sets the primitives registry for all steps
func (w *BaseWorkflow) SetPrimitives(primitives interface{}) {
	w.Primitives = primitives
	for _, step := range w.Steps {
		step.SetPrimitives(primitives)
	}
}

// GetPrimitives returns the primitives registry
func (w *BaseWorkflow) GetPrimitives() interface{} {
	return w.Primitives
}

// SetContext sets the workflow context
func (w *BaseWorkflow) SetContext(context interface{}) {
	w.Context = context
}

// GetContext returns the workflow context
func (w *BaseWorkflow) GetContext() interface{} {
	return w.Context
}

// SetData sets the workflow data
func (w *BaseWorkflow) SetData(data interface{}) {
	w.Data = data
}

// GetData returns the workflow data
func (w *BaseWorkflow) GetData() interface{} {
	return w.Data
}

// GetTotalChildStepCount returns the total number of child steps across all steps
func (w *BaseWorkflow) GetTotalChildStepCount() int {
	total := 0
	for _, step := range w.Steps {
		total += step.GetChildStepCount()
	}
	return total
}

// GetChildStepByGlobalIndex gets a child step by its global index across all steps
func (w *BaseWorkflow) GetChildStepByGlobalIndex(globalChildStepIndex int) (*ChildStep, error) {
	currentIndex := 0
	for _, step := range w.Steps {
		stepChildStepCount := step.GetChildStepCount()
		if globalChildStepIndex < currentIndex+stepChildStepCount {
			return step.GetChildStep(globalChildStepIndex - currentIndex), nil
		}
		currentIndex += stepChildStepCount
	}
	return nil, fmt.Errorf("global child step index out of range: %d", globalChildStepIndex)
}

// GetStepAndChildStepIndices gets the step and child step indices for a given global child step index
func (w *BaseWorkflow) GetStepAndChildStepIndices(globalChildStepIndex int) (int, int, error) {
	currentIndex := 0
	for stepIndex, step := range w.Steps {
		stepChildStepCount := step.GetChildStepCount()
		if globalChildStepIndex < currentIndex+stepChildStepCount {
			return stepIndex, globalChildStepIndex - currentIndex, nil
		}
		currentIndex += stepChildStepCount
	}
	return -1, -1, fmt.Errorf("global child step index out of range: %d", globalChildStepIndex)
}

// GetCurrentStep gets the current step for a given global child step index
func (w *BaseWorkflow) GetCurrentStep(globalChildStepIndex int) (Step, error) {
	if globalChildStepIndex >= w.GetTotalChildStepCount() {
		return nil, fmt.Errorf("global child step index out of range: %d", globalChildStepIndex)
	}

	stepIndex, _, err := w.GetStepAndChildStepIndices(globalChildStepIndex)
	if err != nil {
		return nil, err
	}
	return w.Steps[stepIndex], nil
}

// IsCurrentStepParallel checks if the current step (for the given global child step index) is parallel
func (w *BaseWorkflow) IsCurrentStepParallel(globalChildStepIndex int) (bool, error) {
	step, err := w.GetCurrentStep(globalChildStepIndex)
	if err != nil {
		return false, err
	}
	return step.IsParallel(), nil
}

// GetStepStartIndex gets the starting global index for a given step
func (w *BaseWorkflow) GetStepStartIndex(stepIndex int) (int, error) {
	if stepIndex < 0 || stepIndex >= len(w.Steps) {
		return -1, fmt.Errorf("step index out of range: %d", stepIndex)
	}

	startIndex := 0
	for i := 0; i < stepIndex; i++ {
		startIndex += w.Steps[i].GetChildStepCount()
	}
	return startIndex, nil
}

// HasSteps checks if this workflow contains any steps
func (w *BaseWorkflow) HasSteps() bool {
	return len(w.Steps) > 0
}

// Helper function to generate workflow ID
func generateWorkflowID() string {
	return fmt.Sprintf("workflow-%d", time.Now().UnixNano())
}
