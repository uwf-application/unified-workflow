package model

import (
	"context"
	"time"
)

// Step is the interface that all steps must implement
// Similar to Java's Step abstract class in workflow-common
type Step interface {
	GetName() string
	GetChildSteps() []*ChildStep
	IsParallel() bool
	SetPrimitives(primitives interface{})
	GetPrimitives() interface{}
	GetStartTime() *time.Time
	GetEndTime() *time.Time
	GetDurationMillis() *int64
	GetContextBefore() interface{}
	GetContextAfter() interface{}
	GetDataBefore() interface{}
	GetDataAfter() interface{}
	AddChildStep(childStep *ChildStep) Step
	AddChildSteps(childSteps []*ChildStep) Step
	GetChildStepCount() int
	GetChildStep(index int) *ChildStep
	HasChildSteps() bool
	Run(ctx context.Context, context interface{}, data interface{}) error
	ExecuteWithTiming(ctx context.Context, context interface{}, data interface{}, stepLogic func() error) error
	ExecuteChildStep(ctx context.Context, childStep *ChildStep, context interface{}, data interface{}) error
	ExecuteChildStepWithTiming(ctx context.Context, childStep *ChildStep, context interface{}, data interface{}) error
	StoreStepMetrics(context interface{}, data interface{})
	StoreChildStepMetrics(childStep *ChildStep, context interface{}, data interface{}, startTime, endTime time.Time, errorMessage string)
}

// BaseStep is the base implementation of the Step interface
type BaseStep struct {
	Name          string
	ChildSteps    []*ChildStep
	Parallel      bool
	Primitives    interface{} // Will be typed as primitive.Primitives
	StartTime     *time.Time
	EndTime       *time.Time
	ContextBefore interface{} // Will be typed as primitive.WorkflowContext
	ContextAfter  interface{} // Will be typed as primitive.WorkflowContext
	DataBefore    interface{} // Will be typed as primitive.WorkflowData
	DataAfter     interface{} // Will be typed as primitive.WorkflowData
}

// NewBaseStep creates a new BaseStep
func NewBaseStep(name string, parallel bool) *BaseStep {
	return &BaseStep{
		Name:       name,
		ChildSteps: []*ChildStep{},
		Parallel:   parallel,
	}
}

// GetName returns the name of the step
func (s *BaseStep) GetName() string {
	return s.Name
}

// GetChildSteps returns the child steps
func (s *BaseStep) GetChildSteps() []*ChildStep {
	return s.ChildSteps
}

// IsParallel returns true if child steps should execute in parallel
func (s *BaseStep) IsParallel() bool {
	return s.Parallel
}

// SetPrimitives sets the primitives registry
func (s *BaseStep) SetPrimitives(primitives interface{}) {
	s.Primitives = primitives
}

// GetPrimitives returns the primitives registry
func (s *BaseStep) GetPrimitives() interface{} {
	return s.Primitives
}

// GetStartTime returns the start time
func (s *BaseStep) GetStartTime() *time.Time {
	return s.StartTime
}

// GetEndTime returns the end time
func (s *BaseStep) GetEndTime() *time.Time {
	return s.EndTime
}

// GetDurationMillis returns the duration in milliseconds
func (s *BaseStep) GetDurationMillis() *int64 {
	if s.StartTime != nil && s.EndTime != nil {
		duration := s.EndTime.Sub(*s.StartTime).Milliseconds()
		return &duration
	}
	return nil
}

// GetContextBefore returns the context before execution
func (s *BaseStep) GetContextBefore() interface{} {
	return s.ContextBefore
}

// GetContextAfter returns the context after execution
func (s *BaseStep) GetContextAfter() interface{} {
	return s.ContextAfter
}

// GetDataBefore returns the data before execution
func (s *BaseStep) GetDataBefore() interface{} {
	return s.DataBefore
}

// GetDataAfter returns the data after execution
func (s *BaseStep) GetDataAfter() interface{} {
	return s.DataAfter
}

// AddChildStep adds a child step
func (s *BaseStep) AddChildStep(childStep *ChildStep) Step {
	s.ChildSteps = append(s.ChildSteps, childStep)
	return s
}

// AddChildSteps adds multiple child steps
func (s *BaseStep) AddChildSteps(childSteps []*ChildStep) Step {
	s.ChildSteps = append(s.ChildSteps, childSteps...)
	return s
}

// GetChildStepCount returns the number of child steps
func (s *BaseStep) GetChildStepCount() int {
	return len(s.ChildSteps)
}

// GetChildStep returns a child step by index
func (s *BaseStep) GetChildStep(index int) *ChildStep {
	if index < 0 || index >= len(s.ChildSteps) {
		return nil
	}
	return s.ChildSteps[index]
}

// HasChildSteps returns true if the step has child steps
func (s *BaseStep) HasChildSteps() bool {
	return len(s.ChildSteps) > 0
}

// Run must be implemented by concrete step types
func (s *BaseStep) Run(ctx context.Context, context interface{}, data interface{}) error {
	// This must be implemented by concrete step types
	return nil
}

// ExecuteWithTiming is a template method for concrete Step implementations
func (s *BaseStep) ExecuteWithTiming(ctx context.Context, context interface{}, data interface{}, stepLogic func() error) error {
	// This will be implemented when we have the actual primitive types
	return nil
}

// ExecuteChildStep executes a child step
func (s *BaseStep) ExecuteChildStep(ctx context.Context, childStep *ChildStep, context interface{}, data interface{}) error {
	// This will be implemented when we have the actual primitive types
	return nil
}

// ExecuteChildStepWithTiming executes a child step with timing
func (s *BaseStep) ExecuteChildStepWithTiming(ctx context.Context, childStep *ChildStep, context interface{}, data interface{}) error {
	// This will be implemented when we have the actual primitive types
	return nil
}

// StoreStepMetrics stores step execution metrics
func (s *BaseStep) StoreStepMetrics(context interface{}, data interface{}) {
	// This will be implemented when we have the actual primitive types
}

// StoreChildStepMetrics stores child step execution metrics
func (s *BaseStep) StoreChildStepMetrics(childStep *ChildStep, context interface{}, data interface{}, startTime, endTime time.Time, errorMessage string) {
	// This will be implemented when we have the actual primitive types
}
