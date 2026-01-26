package model

import (
	"context"
	"fmt"
	"time"
)

// SequentialStep is a concrete Step implementation that executes child steps sequentially
type SequentialStep struct {
	*BaseStep
}

// NewSequentialStep creates a new SequentialStep
func NewSequentialStep(name string) *SequentialStep {
	return &SequentialStep{
		BaseStep: NewBaseStep(name, false),
	}
}

// Run executes the step by running all child steps sequentially
func (s *SequentialStep) Run(ctx context.Context, context interface{}, data interface{}) error {
	// Use ExecuteWithTiming template method
	return s.ExecuteWithTiming(ctx, context, data, func() error {
		return s.runChildSteps(ctx, context, data)
	})
}

// runChildSteps executes all child steps sequentially
func (s *SequentialStep) runChildSteps(ctx context.Context, context interface{}, data interface{}) error {
	for _, childStep := range s.ChildSteps {
		if err := s.ExecuteChildStep(ctx, childStep, context, data); err != nil {
			return fmt.Errorf("child step %s failed: %w", childStep.GetName(), err)
		}
	}
	return nil
}

// ExecuteWithTiming implements the template method for timing and metrics
func (s *SequentialStep) ExecuteWithTiming(ctx context.Context, context interface{}, data interface{}, stepLogic func() error) error {
	// Capture state before execution
	now := time.Now()
	s.StartTime = &now
	s.ContextBefore = context
	s.DataBefore = data

	var err error
	defer func() {
		now := time.Now()
		s.EndTime = &now
		s.ContextAfter = context
		s.DataAfter = data
		s.StoreStepMetrics(context, data)
	}()

	// Check for context cancellation before starting
	if ctx.Err() != nil {
		return ctx.Err()
	}

	err = stepLogic()
	return err
}

// ExecuteChildStep executes a child step with the standard request/response/validate pattern
func (s *SequentialStep) ExecuteChildStep(ctx context.Context, childStep *ChildStep, context interface{}, data interface{}) error {
	return s.ExecuteChildStepWithTiming(ctx, childStep, context, data)
}

// ExecuteChildStepWithTiming executes a child step with timing metrics
func (s *SequentialStep) ExecuteChildStepWithTiming(ctx context.Context, childStep *ChildStep, context interface{}, data interface{}) error {
	startTime := time.Now()
	var err error
	defer func() {
		endTime := time.Now()
		errorMessage := ""
		if err != nil {
			errorMessage = err.Error()
		}
		s.StoreChildStepMetrics(childStep, context, data, startTime, endTime, errorMessage)
	}()

	// Check for context cancellation before starting
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 1. Prepare request from context and data using request hook
	_ = childStep.GetRequestHook()(context, data)

	// 2. Execute response (calls Primitive) using response hook
	response := childStep.GetResponseHook()(context, data)

	// Check for context cancellation after response hook
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 3. If response is not nil, validate it using validate hook
	if response != nil && childStep.GetValidateHook() != nil {
		if validateErr := childStep.GetValidateHook()(response); validateErr != nil {
			err = fmt.Errorf("validation failed for child step %s: %w", childStep.GetName(), validateErr)
			return err
		}
	}

	// Store results in workflow data
	if response != nil {
		// Try to put response in data if data supports Put method
		if dataMap, ok := data.(map[string]interface{}); ok {
			dataMap[childStep.GetName()+"Result"] = response
			dataMap[childStep.GetName()+"Completed"] = true
		}
	}

	return nil
}

// StoreStepMetrics stores step execution metrics
func (s *SequentialStep) StoreStepMetrics(context interface{}, data interface{}) {
	// This is a simplified implementation
	// In a real implementation, you would store metrics in the data
	if dataMap, ok := data.(map[string]interface{}); ok {
		stepKey := "step_" + s.Name + "_metrics"
		metrics := map[string]interface{}{
			"stepName":       s.Name,
			"startTime":      s.StartTime,
			"endTime":        s.EndTime,
			"childStepCount": s.GetChildStepCount(),
			"parallel":       s.IsParallel(),
		}
		if s.StartTime != nil && s.EndTime != nil {
			metrics["durationMillis"] = s.EndTime.Sub(*s.StartTime).Milliseconds()
		}
		dataMap[stepKey] = metrics
	}
}

// StoreChildStepMetrics stores child step execution metrics
func (s *SequentialStep) StoreChildStepMetrics(childStep *ChildStep, context interface{}, data interface{}, startTime, endTime time.Time, errorMessage string) {
	// This is a simplified implementation
	if dataMap, ok := data.(map[string]interface{}); ok {
		childStepKey := "childStep_" + childStep.GetName() + "_metrics"
		metrics := map[string]interface{}{
			"childStepName":  childStep.GetName(),
			"parentStepName": s.Name,
			"startTime":      startTime,
			"endTime":        endTime,
			"errorMessage":   errorMessage,
		}
		if !startTime.IsZero() && !endTime.IsZero() {
			metrics["durationMillis"] = endTime.Sub(startTime).Milliseconds()
		}
		dataMap[childStepKey] = metrics
	}
}
