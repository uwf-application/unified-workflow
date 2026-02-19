package steps

import (
	"fmt"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/workflows/child_steps"

	"github.com/google/uuid"
)

// AMLValidationStep validates a transaction using the AML (Anti-Money Laundering) service
type AMLValidationStep struct {
	*AntifraudStep
}

// NewAMLValidationStep creates a new AMLValidationStep
func NewAMLValidationStep(endpoint string) *AMLValidationStep {
	step := &AMLValidationStep{
		AntifraudStep: NewAntifraudStep("aml-validation", endpoint),
	}

	// Add reusable child steps for complete AML validation flow
	step.AddChildSteps([]*model.ChildStep{
		// Child Step 1: Prepare AML request
		child_steps.CreateAntifraudPrepareAMLRequestChildStep(),
		// Child Step 2: Call AML validation async
		child_steps.CreateAntifraudCallAMLValidationAsyncChildStep(),
		// Child Step 3: Process AML response
		child_steps.CreateAntifraudProcessAMLResponseChildStep(),
		// Child Step 4: Validate AML response structure
		child_steps.CreateAntifraudValidateAMLResponseChildStep(),
		// Child Step 5: Validate AML result against business rules
		child_steps.CreateAntifraudValidateAMLResultChildStep(),
		// Child Step 6: Store AML resolution
		child_steps.CreateAntifraudStoreAMLResolutionChildStep(),
		// Child Step 7: Add AML to transaction
		child_steps.CreateAntifraudAddAMLToTransactionChildStep(),
	})

	return step
}

// prepareAMLRequest prepares the AML validation request
func (s *AMLValidationStep) prepareAMLRequest(context interface{}, data interface{}) interface{} {
	fmt.Println("Preparing AML validation request...")

	// In a real implementation, you would extract transaction from workflow data
	// For now, create a mock transaction
	transaction := map[string]interface{}{
		"af_id":       uuid.NewString(),
		"af_add_date": time.Now().Format(time.RFC3339Nano),
		"transaction": map[string]interface{}{
			"id":          uuid.NewString(),
			"type":        "deposit",
			"amount":      "100000",
			"currency":    "KZT",
			"client_id":   uuid.NewString(),
			"client_name": "John Smith",
		},
	}

	fmt.Printf("Prepared AML request for transaction: %s\n", transaction["af_id"])
	return transaction
}

// validateAMLRequest validates the AML request
func (s *AMLValidationStep) validateAMLRequest(request interface{}) error {
	fmt.Println("Validating AML request...")

	req, ok := request.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid AML request type: %T", request)
	}

	// Validate required fields
	requiredFields := []string{"af_id", "transaction"}
	for _, field := range requiredFields {
		if _, exists := req[field]; !exists {
			return fmt.Errorf("missing AML request field: %s", field)
		}
	}

	fmt.Println("AML request validation passed")
	return nil
}

// callAMLValidationAsync calls the AML validation service asynchronously
func (s *AMLValidationStep) callAMLValidationAsync(context interface{}, data interface{}) interface{} {
	fmt.Println("Calling AML validation async...")

	// Get antifraud service
	_, err := s.GetAntifraudService()
	if err != nil {
		return fmt.Errorf("failed to get antifraud service: %w", err)
	}

	// In production, this would initiate an actual async call
	// For high TPS, we don't simulate delays
	go func() {
		// Real async call would happen here
		fmt.Println("Async AML validation completed")
	}()

	return map[string]interface{}{
		"aml_validation_id": uuid.NewString(),
		"status":            "processing",
		"started_at":        time.Now().Format(time.RFC3339),
		"service":           "AML",
	}
}

// processAMLResponse processes the AML validation response
func (s *AMLValidationStep) processAMLResponse(context interface{}, data interface{}) interface{} {
	fmt.Println("Processing AML response...")

	// In production, this would process the actual SDK response
	// For high TPS, we process immediately without delays
	return map[string]interface{}{
		"service_name": "AML",
		"resolution":   "PASS", // PASS, FAIL, REVIEW
		"score":        85,     // Risk score 0-100
		"details":      "Transaction passed AML screening",
		"checked_at":   time.Now().Format(time.RFC3339),
		"risk_factors": []string{
			"Low risk country",
			"Amount within limits",
			"No PEP involvement",
		},
	}
}

// validateAMLResponse validates the AML response structure
func (s *AMLValidationStep) validateAMLResponse(response interface{}) error {
	fmt.Println("Validating AML response...")

	resp, ok := response.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid AML response type: %T", response)
	}

	requiredFields := []string{"service_name", "resolution", "score", "details"}
	for _, field := range requiredFields {
		if _, exists := resp[field]; !exists {
			return fmt.Errorf("missing AML response field: %s", field)
		}
	}

	// Validate resolution value
	resolution, ok := resp["resolution"].(string)
	if !ok {
		return fmt.Errorf("invalid resolution type")
	}

	validResolutions := map[string]bool{"PASS": true, "FAIL": true, "REVIEW": true}
	if !validResolutions[resolution] {
		return fmt.Errorf("invalid resolution value: %s", resolution)
	}

	// Validate score range
	score, ok := resp["score"].(int)
	if !ok || score < 0 || score > 100 {
		return fmt.Errorf("invalid score value: %v", score)
	}

	fmt.Println("AML response validation passed")
	return nil
}

// validateAMLResult validates the AML result based on business rules
func (s *AMLValidationStep) validateAMLResult(result interface{}) error {
	fmt.Println("Validating AML result against business rules...")

	res, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid AML result type: %T", result)
	}

	resolution, _ := res["resolution"].(string)
	score, _ := res["score"].(int)

	// Business rule: If resolution is FAIL, reject
	if resolution == "FAIL" {
		return fmt.Errorf("AML validation failed: transaction rejected")
	}

	// Business rule: If score > 70, flag for review
	if score > 70 && resolution == "PASS" {
		fmt.Println("AML Warning: High risk score, consider manual review")
		// Continue processing but log warning
	}

	// Business rule: If score > 90, fail even if resolution is PASS
	if score > 90 {
		return fmt.Errorf("AML validation failed: risk score too high (%d)", score)
	}

	fmt.Println("AML result validation passed")
	return nil
}

// storeAMLResolution stores the AML resolution in antifraud system
func (s *AMLValidationStep) storeAMLResolution(context interface{}, data interface{}) interface{} {
	fmt.Println("Storing AML resolution...")

	// In a real implementation, this would call antifraudService.StoreServiceResolution()
	// For high TPS, no artificial delays
	return map[string]interface{}{
		"operation":     "store_aml_resolution",
		"success":       true,
		"stored_at":     time.Now().Format(time.RFC3339),
		"resolution_id": uuid.NewString(),
	}
}

// addAMLToTransaction adds the AML check to transaction aggregation
func (s *AMLValidationStep) addAMLToTransaction(context interface{}, data interface{}) interface{} {
	fmt.Println("Adding AML check to transaction...")

	// In a real implementation, this would call antifraudService.AddTransactionServiceCheck()
	// For high TPS, no artificial delays
	return map[string]interface{}{
		"operation": "add_aml_to_transaction",
		"success":   true,
		"added_at":  time.Now().Format(time.RFC3339),
		"message":   "AML check added to transaction aggregation",
	}
}

// ExecuteStepLogic executes the main step logic (alternative to child steps)
func (s *AMLValidationStep) ExecuteStepLogic(ctx interface{}, context interface{}, data interface{}) error {
	fmt.Println("Executing AMLValidationStep logic...")

	// 1. Prepare AML request
	request := s.prepareAMLRequest(context, data)

	// 2. Validate request
	if err := s.validateAMLRequest(request); err != nil {
		return fmt.Errorf("AML request validation failed: %w", err)
	}

	// 3. Call async AML validation
	asyncResult := s.callAMLValidationAsync(context, data)
	fmt.Printf("AML validation started: %v\n", asyncResult)

	// 4. Process response
	response := s.processAMLResponse(context, data)

	// 5. Validate response
	if err := s.validateAMLResponse(response); err != nil {
		return fmt.Errorf("AML response validation failed: %w", err)
	}

	// 6. Validate result against business rules
	if err := s.validateAMLResult(response); err != nil {
		return fmt.Errorf("AML result validation failed: %w", err)
	}

	// 7. Store resolution
	resolutionResult := s.storeAMLResolution(context, data)
	fmt.Printf("AML resolution stored: %v\n", resolutionResult)

	// 8. Add to transaction
	transactionResult := s.addAMLToTransaction(context, data)
	fmt.Printf("AML added to transaction: %v\n", transactionResult)

	fmt.Println("AMLValidationStep completed successfully")
	return nil
}
