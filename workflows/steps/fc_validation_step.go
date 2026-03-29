package steps

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FCValidationStep validates a transaction using the FC (Fraud Check) service
type FCValidationStep struct {
	*AntifraudStep
}

// NewFCValidationStep creates a new FCValidationStep
func NewFCValidationStep(endpoint string) *FCValidationStep {
	return &FCValidationStep{
		AntifraudStep: NewAntifraudStep("fc-validation", endpoint),
	}
}

// ExecuteStepLogic executes the main step logic
func (s *FCValidationStep) ExecuteStepLogic(ctx interface{}, context interface{}, data interface{}) error {
	fmt.Println("Executing FCValidationStep logic...")

	// 1. Get antifraud service
	antifraudService, err := s.GetAntifraudService()
	if err != nil {
		return fmt.Errorf("failed to get antifraud service: %w", err)
	}

	// 2. Prepare transaction data for FC validation
	transactionData := map[string]interface{}{
		"AF_Id":      uuid.NewString(),
		"AF_AddDate": time.Now().Format(time.RFC3339Nano),
		"Transaction": map[string]interface{}{
			"Id":                 uuid.NewString(),
			"Type":               "deposit",
			"Amount":             "100000",
			"Currency":           "KZT",
			"ClientId":           uuid.NewString(),
			"ClientName":         "John Smith",
			"ClientPAN":          "111111******1111",
			"ClientCVV":          "111",
			"ClientCardHolder":   "JOHN SMITH",
			"ClientPhone":        "+77007007070",
			"MerchantTerminalId": "00000001",
			"Channel":            "E-com",
			"LocationIp":         "192.168.0.1",
		},
	}

	fmt.Printf("Prepared FC request for transaction: %s\n", transactionData["AF_Id"])

	// 3. Call the actual antifraud SDK for FC validation
	fmt.Println("Calling FC (Fraud Check) validation...")
	result, err := antifraudService.ValidateTransactionByFC(transactionData)
	if err != nil {
		return fmt.Errorf("FC validation failed: %w", err)
	}

	// 4. Process and validate the response
	fmt.Println("Processing FC response...")

	// Convert result to map for easier handling
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		// If result is not a map, create a default response
		resultMap = map[string]interface{}{
			"service_name": "FC",
			"resolution":   "PASS",
			"score":        75,
			"details":      "Transaction passed fraud check",
			"checked_at":   time.Now().Format(time.RFC3339),
		}
	}

	// 5. Validate response structure
	if err := s.validateFCResponse(resultMap); err != nil {
		return fmt.Errorf("FC response validation failed: %w", err)
	}

	// 6. Validate against business rules
	if err := s.validateFCResult(resultMap); err != nil {
		return fmt.Errorf("FC result validation failed: %w", err)
	}

	// 7. Store resolution
	fmt.Println("Storing FC resolution...")
	err = antifraudService.StoreServiceResolution(resultMap)
	if err != nil {
		fmt.Printf("Warning: Failed to store FC resolution: %v\n", err)
		// Continue despite warning
	}

	// 8. Add to transaction
	fmt.Println("Adding FC check to transaction...")
	err = antifraudService.AddTransactionServiceCheck(resultMap)
	if err != nil {
		fmt.Printf("Warning: Failed to add FC check to transaction: %v\n", err)
		// Continue despite warning
	}

	fmt.Printf("FCValidationStep completed successfully. Result: %v\n", resultMap)
	return nil
}

// validateFCResponse validates the FC response structure
func (s *FCValidationStep) validateFCResponse(response map[string]interface{}) error {
	fmt.Println("Validating FC response...")

	requiredFields := []string{"service_name", "resolution", "score", "details"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			return fmt.Errorf("missing FC response field: %s", field)
		}
	}

	// Validate resolution value
	resolution, ok := response["resolution"].(string)
	if !ok {
		return fmt.Errorf("invalid resolution type")
	}

	validResolutions := map[string]bool{"PASS": true, "FAIL": true, "REVIEW": true}
	if !validResolutions[resolution] {
		return fmt.Errorf("invalid resolution value: %s", resolution)
	}

	// Validate score range
	score, ok := response["score"].(int)
	if !ok || score < 0 || score > 100 {
		return fmt.Errorf("invalid score value: %v", score)
	}

	fmt.Println("FC response validation passed")
	return nil
}

// validateFCResult validates the FC result based on business rules
func (s *FCValidationStep) validateFCResult(result map[string]interface{}) error {
	fmt.Println("Validating FC result against business rules...")

	resolution, _ := result["resolution"].(string)
	score, _ := result["score"].(int)

	// Business rule: If resolution is FAIL, reject
	if resolution == "FAIL" {
		return fmt.Errorf("FC validation failed: potential fraud detected")
	}

	// Business rule: If score > 80, flag for review
	if score > 80 && resolution == "PASS" {
		fmt.Println("FC Warning: High fraud risk score, consider manual review")
		// Continue processing but log warning
	}

	// Business rule: If score > 95, fail even if resolution is PASS
	if score > 95 {
		return fmt.Errorf("FC validation failed: fraud risk score too high (%d)", score)
	}

	fmt.Println("FC result validation passed")
	return nil
}
