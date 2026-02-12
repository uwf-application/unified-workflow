package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Example demonstrates how to use the Workflow SDK
func Example() {
	// Create SDK configuration
	config := &SDKConfig{
		WorkflowAPIEndpoint: "http://localhost:8080",
		Timeout:             30 * time.Second,
		MaxRetries:          3,
		RetryDelay:          1 * time.Second,
		AuthToken:           "your-auth-token",
		EnableValidation:    true,
		EnableSanitization:  true,
		StrictValidation:    false,
		DefaultValidationRules: []ValidationRule{
			{
				Field:     "user_id",
				Required:  true,
				RuleType:  ValidationRuleTypeString,
				MinLength: &[]int{1}[0],
				MaxLength: &[]int{255}[0],
			},
			{
				Field:    "amount",
				Required: true,
				RuleType: ValidationRuleTypeNumber,
				MinValue: &[]float64{0.01}[0],
			},
			{
				Field:    "email",
				Required: false,
				RuleType: ValidationRuleTypeEmail,
			},
		},
	}

	// Create SDK client
	client, err := NewClient(config)
	if err != nil {
		fmt.Printf("Failed to create SDK client: %v\n", err)
		return
	}
	defer client.Close()

	// Example 1: Execute workflow with raw data
	fmt.Println("=== Example 1: Execute workflow with raw data ===")
	executeWorkflowExample(client)

	// Example 2: Execute workflow from HTTP request
	fmt.Println("\n=== Example 2: Execute workflow from HTTP request ===")
	executeFromHTTPRequestExample(client)

	// Example 3: Validate and execute workflow
	fmt.Println("\n=== Example 3: Validate and execute workflow ===")
	validateAndExecuteExample(client)

	// Example 4: Get execution status
	fmt.Println("\n=== Example 4: Get execution status ===")
	getExecutionStatusExample(client)
}

func executeWorkflowExample(client WorkflowSDKClient) {
	ctx := context.Background()
	workflowID := "payment-processing-workflow"

	// Prepare input data
	inputData := map[string]interface{}{
		"user_id": "user_12345",
		"amount":  99.99,
		"email":   "user@example.com",
		"metadata": map[string]interface{}{
			"source": "web",
			"device": "mobile",
		},
	}

	// Execute workflow
	resp, err := client.ExecuteWorkflow(ctx, workflowID, inputData)
	if err != nil {
		fmt.Printf("Failed to execute workflow: %v\n", err)
		return
	}

	fmt.Printf("Workflow execution started!\n")
	fmt.Printf("Run ID: %s\n", resp.RunID)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Status URL: %s\n", resp.StatusURL)
	fmt.Printf("Poll after: %d ms\n", resp.PollAfterMs)
	fmt.Printf("Estimated completion: %d ms\n", resp.EstimatedCompletionMs)

	// Check validation result if available
	if resp.ValidationResult != nil {
		fmt.Printf("Validation passed: %v\n", resp.ValidationResult.Valid)
		if len(resp.ValidationResult.Warnings) > 0 {
			fmt.Printf("Validation warnings: %d\n", len(resp.ValidationResult.Warnings))
		}
	}
}

func executeFromHTTPRequestExample(client WorkflowSDKClient) {
	ctx := context.Background()
	workflowID := "payment-processing-workflow"

	// Create a mock HTTP request
	req, _ := http.NewRequest("POST", "/api/payments", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("X-User-ID", "user_12345")
	req.Header.Set("X-Session-ID", "session_67890")

	// Add query parameters
	q := req.URL.Query()
	q.Add("amount", "99.99")
	q.Add("currency", "USD")
	req.URL.RawQuery = q.Encode()

	// Execute workflow from HTTP request
	resp, err := client.ExecuteFromHTTPRequest(ctx, workflowID, req)
	if err != nil {
		fmt.Printf("Failed to execute workflow from HTTP request: %v\n", err)
		return
	}

	fmt.Printf("Workflow execution started from HTTP request!\n")
	fmt.Printf("Run ID: %s\n", resp.RunID)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Context included: %v\n", resp.ContextIncluded)
}

func validateAndExecuteExample(client WorkflowSDKClient) {
	ctx := context.Background()
	workflowID := "payment-processing-workflow"

	// Prepare input data with potential issues
	inputData := map[string]interface{}{
		"user_id": "user_12345",
		"amount":  -10.0,           // Invalid: negative amount
		"email":   "invalid-email", // Invalid email format
		"metadata": map[string]interface{}{
			"source": "web",
		},
	}

	// Define custom validation rules
	customRules := []ValidationRule{
		{
			Field:     "user_id",
			Required:  true,
			RuleType:  ValidationRuleTypeString,
			MinLength: &[]int{5}[0],
			MaxLength: &[]int{50}[0],
		},
		{
			Field:    "amount",
			Required: true,
			RuleType: ValidationRuleTypeNumber,
			MinValue: &[]float64{0.01}[0],
			MaxValue: &[]float64{10000}[0],
		},
		{
			Field:    "email",
			Required: false,
			RuleType: ValidationRuleTypeEmail,
		},
		{
			Field:         "metadata.source",
			Required:      true,
			RuleType:      ValidationRuleTypeString,
			AllowedValues: []string{"web", "mobile", "api"},
		},
	}

	// Validate and execute workflow
	resp, err := client.ValidateAndExecuteWorkflow(ctx, workflowID, inputData, customRules)
	if err != nil {
		fmt.Printf("Failed to validate and execute workflow: %v\n", err)

		// Check if it's a validation error
		if sdkErr, ok := err.(*SDKError); ok && sdkErr.Code == ErrCodeValidationFailed {
			fmt.Println("Validation failed with the following errors:")
			if details, ok := sdkErr.Details["validation_result"].(*ValidationResult); ok {
				for _, err := range details.Errors {
					fmt.Printf("  - Field: %s, Code: %s, Message: %s\n",
						err.Field, err.Code, err.Message)
				}
			}
		}
		return
	}

	fmt.Printf("Workflow execution started after validation!\n")
	fmt.Printf("Run ID: %s\n", resp.RunID)
	fmt.Printf("Status: %s\n", resp.Status)

	if resp.ValidationResult != nil {
		fmt.Printf("Validation passed: %v\n", resp.ValidationResult.Valid)
		if len(resp.ValidationResult.Errors) > 0 {
			fmt.Printf("Validation errors: %d\n", len(resp.ValidationResult.Errors))
		}
		if len(resp.ValidationResult.Warnings) > 0 {
			fmt.Printf("Validation warnings: %d\n", len(resp.ValidationResult.Warnings))
		}
	}
}

func getExecutionStatusExample(client WorkflowSDKClient) {
	ctx := context.Background()
	runID := "run_1234567890"

	// Get execution status
	statusResp, err := client.GetExecutionStatus(ctx, runID)
	if err != nil {
		fmt.Printf("Failed to get execution status: %v\n", err)
		return
	}

	if statusResp.Status != nil {
		fmt.Printf("Execution Status:\n")
		fmt.Printf("  Run ID: %s\n", statusResp.Status.RunID)
		fmt.Printf("  Status: %s\n", statusResp.Status.Status)
		fmt.Printf("  Progress: %.2f\n", statusResp.Status.Progress)
		fmt.Printf("  Current Step: %s\n", statusResp.Status.CurrentStep)
		fmt.Printf("  Is Terminal: %v\n", statusResp.Status.IsTerminal)

		if statusResp.Status.ErrorMessage != "" {
			fmt.Printf("  Error: %s\n", statusResp.Status.ErrorMessage)
		}

		if statusResp.Status.StartTime != nil {
			fmt.Printf("  Started: %v\n", statusResp.Status.StartTime.Format(time.RFC3339))
		}

		if statusResp.Status.EndTime != nil {
			fmt.Printf("  Ended: %v\n", statusResp.Status.EndTime.Format(time.RFC3339))
		}
	}

	// Get execution details
	detailsResp, err := client.GetExecutionDetails(ctx, runID)
	if err != nil {
		fmt.Printf("Failed to get execution details: %v\n", err)
		return
	}

	if detailsResp.Details != nil {
		fmt.Printf("\nExecution Details:\n")
		fmt.Printf("  Total Steps: %d\n", len(detailsResp.Details.Steps))
		fmt.Printf("  Primitives Used: %v\n", detailsResp.Details.PrimitivesUsed)
		fmt.Printf("  Total Duration: %d ms\n", detailsResp.Details.TotalDurationMs)

		if len(detailsResp.Details.Steps) > 0 {
			fmt.Printf("  Steps:\n")
			for _, step := range detailsResp.Details.Steps {
				fmt.Printf("    - %s (Index: %d, Status: %s)\n",
					step.Name, step.StepIndex, step.Status)
			}
		}
	}
}

// ExampleHTTPHandler demonstrates how to use the SDK in an HTTP handler
func ExampleHTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Create SDK configuration
		config := &SDKConfig{
			WorkflowAPIEndpoint: "http://localhost:8080",
			Timeout:             30 * time.Second,
			MaxRetries:          3,
			AuthToken:           "your-auth-token",
			EnableValidation:    true,
			EnableSanitization:  true,
		}

		// Create SDK client
		client, err := NewClient(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create SDK client: %v", err), http.StatusInternalServerError)
			return
		}
		defer client.Close()

		// Extract workflow ID from URL path or query parameter
		workflowID := r.URL.Query().Get("workflow_id")
		if workflowID == "" {
			workflowID = "default-workflow"
		}

		// Execute workflow from HTTP request
		resp, err := client.ExecuteFromHTTPRequest(ctx, workflowID, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute workflow: %v", err), http.StatusBadRequest)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		responseJSON := fmt.Sprintf(`{
			"run_id": "%s",
			"status": "%s",
			"status_url": "%s",
			"result_url": "%s",
			"poll_after_ms": %d,
			"estimated_completion_ms": %d,
			"expires_at": "%s"
		}`,
			resp.RunID, resp.Status, resp.StatusURL, resp.ResultURL,
			resp.PollAfterMs, resp.EstimatedCompletionMs, resp.ExpiresAt.Format(time.RFC3339))

		w.Write([]byte(responseJSON))
	}
}
