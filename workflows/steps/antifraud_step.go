package steps

import (
	"context"
	"fmt"
	"os"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/internal/primitive"

	"github.com/baraic-io/antifraud-go"
)

// AntifraudStep is the base class for all antifraud steps
type AntifraudStep struct {
	*model.BaseStep
	endpoint         string
	antifraudService primitive.AntifraudService
	initialized      bool
}

// NewAntifraudStep creates a new AntifraudStep
func NewAntifraudStep(name, endpoint string) *AntifraudStep {
	return &AntifraudStep{
		BaseStep:    model.NewBaseStep(name, false), // sequential execution
		endpoint:    endpoint,
		initialized: false,
	}
}

// InitializeService initializes the antifraud service
func (s *AntifraudStep) InitializeService() error {
	if s.initialized {
		return nil
	}

	// Try to get API key from environment variable
	apiKey := getAntifraudAPIKey()
	if apiKey == "" {
		// If no API key is available, use mock service for testing
		fmt.Printf("Note: No API key found. Using mock antifraud service for testing (endpoint: %s)\n", s.endpoint)
		fmt.Println("To use the actual antifraud SDK, set ANTIFRAUD_API_KEY environment variable.")
		fmt.Println("Example: export ANTIFRAUD_API_KEY='your-api-key-here'")
		s.antifraudService = createMockAntifraudService()
		s.initialized = true
		return nil
	}

	// Initialize antifraud service with endpoint and API key
	config := &primitive.Config{
		AntifraudAPIKey:                  apiKey,
		AntifraudAPIHost:                 s.endpoint,
		AntifraudEnabled:                 true,
		AntifraudTimeout:                 30,
		AntifraudMaxRetries:              3,
		AntifraudCircuitBreakerEnabled:   true,
		AntifraudCircuitBreakerThreshold: 5,
		AntifraudCircuitBreakerTimeout:   60,
	}

	// Try to initialize with the actual SDK
	err := primitive.Init(config)
	if err != nil {
		return fmt.Errorf("failed to initialize primitive: %w", err)
	}

	if primitive.Default == nil {
		return fmt.Errorf("primitive not initialized")
	}

	s.antifraudService = primitive.Default.Antifraud
	s.initialized = true
	fmt.Printf("✓ Antifraud service initialized with endpoint: %s\n", s.endpoint)
	return nil
}

// GetAntifraudService returns the antifraud service
func (s *AntifraudStep) GetAntifraudService() (primitive.AntifraudService, error) {
	if !s.initialized {
		err := s.InitializeService()
		if err != nil {
			return nil, err
		}
	}
	return s.antifraudService, nil
}

// Run implements the Step interface
func (s *AntifraudStep) Run(ctx context.Context, context interface{}, data interface{}) error {
	// Initialize service
	err := s.InitializeService()
	if err != nil {
		return fmt.Errorf("failed to initialize antifraud service: %w", err)
	}

	// Execute with timing
	startTime := time.Now()
	s.StartTime = &startTime

	// Execute child steps if any
	if s.HasChildSteps() {
		for _, childStep := range s.GetChildSteps() {
			err := s.ExecuteChildStepWithTiming(ctx, childStep, context, data)
			if err != nil {
				return fmt.Errorf("child step %s failed: %w", childStep.GetName(), err)
			}
		}
	} else {
		// Execute step logic directly if no child steps
		err := s.ExecuteStepLogic(ctx, context, data)
		if err != nil {
			return fmt.Errorf("step logic failed: %w", err)
		}
	}

	endTime := time.Now()
	s.EndTime = &endTime

	return nil
}

// ExecuteStepLogic is the main logic for the step (to be overridden by concrete steps)
func (s *AntifraudStep) ExecuteStepLogic(ctx context.Context, context interface{}, data interface{}) error {
	// Base implementation does nothing
	// Concrete steps should override this method
	return nil
}

// ExecuteChildStepWithTiming executes a child step with timing and error handling
func (s *AntifraudStep) ExecuteChildStepWithTiming(ctx context.Context, childStep *model.ChildStep, context interface{}, data interface{}) error {
	startTime := time.Now()

	// Execute request hook if present
	var requestResult interface{}
	if requestHook := childStep.GetRequestHook(); requestHook != nil {
		requestResult = requestHook(context, data)
	}

	// Execute response hook if present (simulating async response processing)
	var responseResult interface{}
	if responseHook := childStep.GetResponseHook(); responseHook != nil {
		responseResult = responseHook(context, data)
	}

	// Execute validate hook if present
	if validateHook := childStep.GetValidateHook(); validateHook != nil {
		// Use response result if available, otherwise use request result
		validateTarget := responseResult
		if validateTarget == nil {
			validateTarget = requestResult
		}
		if validateTarget != nil {
			err := validateHook(validateTarget)
			if err != nil {
				endTime := time.Now()
				s.StoreChildStepMetrics(childStep, context, data, startTime, endTime, err.Error())
				return fmt.Errorf("validation failed for child step %s: %w", childStep.GetName(), err)
			}
		}
	}

	endTime := time.Now()
	s.StoreChildStepMetrics(childStep, context, data, startTime, endTime, "")

	return nil
}

// StoreChildStepMetrics stores metrics for child step execution
func (s *AntifraudStep) StoreChildStepMetrics(childStep *model.ChildStep, context interface{}, data interface{}, startTime, endTime time.Time, errorMessage string) {
	// Store timing information
	duration := endTime.Sub(startTime)

	// In a real implementation, this would store metrics to a monitoring system
	// For now, we'll just log the information
	fmt.Printf("Child step %s executed in %v", childStep.GetName(), duration)
	if errorMessage != "" {
		fmt.Printf(" with error: %s", errorMessage)
	}
	fmt.Println()
}

// Helper function to extract transaction from workflow data
func (s *AntifraudStep) ExtractTransaction(data interface{}) (antifraud.AF_Transaction, error) {
	// This is a simplified implementation
	// In a real implementation, you would properly extract the transaction from workflow data
	// For now, we'll return a mock transaction
	return antifraud.AF_Transaction{}, nil
}

// Helper function to store result in workflow data
func (s *AntifraudStep) StoreResultInWorkflowData(data interface{}, key string, value interface{}) error {
	// This is a simplified implementation
	// In a real implementation, you would properly store data in the workflow data structure
	return nil
}

// Helper function to get result from workflow data
func (s *AntifraudStep) GetResultFromWorkflowData(data interface{}, key string) (interface{}, bool) {
	// This is a simplified implementation
	// In a real implementation, you would properly retrieve data from the workflow data structure
	return nil, false
}

// getAntifraudAPIKey retrieves the antifraud API key from environment variables
func getAntifraudAPIKey() string {
	// Try environment variable first
	apiKey := os.Getenv("ANTIFRAUD_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	// Try alternative environment variable names
	apiKey = os.Getenv("AF_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	// Check config file or other sources
	// For now, return empty string if not found
	return ""
}

// createMockAntifraudService creates a mock antifraud service for testing
func createMockAntifraudService() primitive.AntifraudService {
	return &mockAntifraudService{}
}

// mockAntifraudService is a mock implementation for testing
type mockAntifraudService struct{}

func (m *mockAntifraudService) StoreTransaction(afTransaction interface{}) error {
	fmt.Println("Mock: StoreTransaction called")
	return nil
}

func (m *mockAntifraudService) ValidateTransactionByAML(afTransaction interface{}) (interface{}, error) {
	fmt.Println("Mock: ValidateTransactionByAML called")
	return map[string]interface{}{
		"service_name": "AML",
		"resolution":   "PASS",
		"score":        85,
		"details":      "Mock AML validation passed",
	}, nil
}

func (m *mockAntifraudService) ValidateTransactionByFC(afTransaction interface{}) (interface{}, error) {
	fmt.Println("Mock: ValidateTransactionByFC called")
	return map[string]interface{}{
		"service_name": "FC",
		"resolution":   "PASS",
		"score":        75,
		"details":      "Mock FC validation passed",
	}, nil
}

func (m *mockAntifraudService) ValidateTransactionByML(afTransaction interface{}) (interface{}, error) {
	fmt.Println("Mock: ValidateTransactionByML called")
	return map[string]interface{}{
		"service_name": "ML",
		"resolution":   "PASS",
		"score":        80,
		"details":      "Mock ML validation passed",
	}, nil
}

func (m *mockAntifraudService) StoreServiceResolution(resolution interface{}) error {
	fmt.Println("Mock: StoreServiceResolution called")
	return nil
}

func (m *mockAntifraudService) AddTransactionServiceCheck(resolution interface{}) error {
	fmt.Println("Mock: AddTransactionServiceCheck called")
	return nil
}

func (m *mockAntifraudService) FinalizeTransaction(afTransaction interface{}) (interface{}, error) {
	fmt.Println("Mock: FinalizeTransaction called")
	return map[string]interface{}{
		"transaction_id": "mock-txn-123",
		"final_decision": "APPROVED",
		"risk_score":     75,
		"reasons":        []string{"All checks passed", "Low risk profile"},
	}, nil
}

func (m *mockAntifraudService) StoreFinalResolution(resolution interface{}) error {
	fmt.Println("Mock: StoreFinalResolution called")
	return nil
}

func (m *mockAntifraudService) HealthCheck() (bool, error) {
	return true, nil
}

func (m *mockAntifraudService) GetConfig() interface{} {
	return map[string]interface{}{
		"endpoint": "mock-endpoint",
		"enabled":  true,
	}
}
