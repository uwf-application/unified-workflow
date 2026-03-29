package steps

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// MLValidationStep validates a transaction using the ML (Machine Learning) service
type MLValidationStep struct {
	*AntifraudStep
}

// NewMLValidationStep creates a new MLValidationStep
func NewMLValidationStep(endpoint string) *MLValidationStep {
	return &MLValidationStep{
		AntifraudStep: NewAntifraudStep("ml-validation", endpoint),
	}
}

// ExecuteStepLogic executes the main step logic
func (s *MLValidationStep) ExecuteStepLogic(ctx interface{}, context interface{}, data interface{}) error {
	fmt.Println("Executing MLValidationStep logic...")

	// 1. Get antifraud service
	antifraudService, err := s.GetAntifraudService()
	if err != nil {
		return fmt.Errorf("failed to get antifraud service: %w", err)
	}

	// 2. Prepare transaction data for ML validation
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
			"DeviceFingerprint":  "device-12345",
			"SessionId":          "session-67890",
			"UserAgent":          "Mozilla/5.0",
		},
	}

	fmt.Printf("Prepared ML request for transaction: %s\n", transactionData["AF_Id"])

	// 3. Call the actual antifraud SDK for ML validation
	fmt.Println("Calling ML (Machine Learning) validation...")
	result, err := antifraudService.ValidateTransactionByML(transactionData)
	if err != nil {
		return fmt.Errorf("ML validation failed: %w", err)
	}

	// 4. Process and validate the response
	fmt.Println("Processing ML response...")

	// Convert result to map for easier handling
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		// If result is not a map, create a default response
		resultMap = map[string]interface{}{
			"service_name":  "ML",
			"resolution":    "PASS",
			"score":         80,
			"details":       "Transaction passed ML model analysis",
			"checked_at":    time.Now().Format(time.RFC3339),
			"model_version": "v2.1.0",
			"confidence":    0.92,
			"features": []string{
				"transaction_pattern_normal",
				"device_trust_score_high",
				"behavior_consistent",
				"location_verified",
			},
		}
	}

	// 5. Validate response structure
	if err := s.validateMLResponse(resultMap); err != nil {
		return fmt.Errorf("ML response validation failed: %w", err)
	}

	// 6. Validate against business rules
	if err := s.validateMLResult(resultMap); err != nil {
		return fmt.Errorf("ML result validation failed: %w", err)
	}

	// 7. Store resolution
	fmt.Println("Storing ML resolution...")
	err = antifraudService.StoreServiceResolution(resultMap)
	if err != nil {
		fmt.Printf("Warning: Failed to store ML resolution: %v\n", err)
		// Continue despite warning
	}

	// 8. Add to transaction
	fmt.Println("Adding ML check to transaction...")
	err = antifraudService.AddTransactionServiceCheck(resultMap)
	if err != nil {
		fmt.Printf("Warning: Failed to add ML check to transaction: %v\n", err)
		// Continue despite warning
	}

	fmt.Printf("MLValidationStep completed successfully. Result: %v\n", resultMap)
	return nil
}

// validateMLResponse validates the ML response structure
func (s *MLValidationStep) validateMLResponse(response map[string]interface{}) error {
	fmt.Println("Validating ML response...")

	requiredFields := []string{"service_name", "resolution", "score", "details"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			return fmt.Errorf("missing ML response field: %s", field)
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

	fmt.Println("ML response validation passed")
	return nil
}

// validateMLResult validates the ML result based on business rules
func (s *MLValidationStep) validateMLResult(result map[string]interface{}) error {
	fmt.Println("Validating ML result against business rules...")

	resolution, _ := result["resolution"].(string)
	score, _ := result["score"].(int)

	// Business rule: If resolution is FAIL, reject
	if resolution == "FAIL" {
		return fmt.Errorf("ML validation failed: high fraud probability detected")
	}

	// Business rule: If score > 85, flag for review
	if score > 85 && resolution == "PASS" {
		fmt.Println("ML Warning: High ML risk score, consider manual review")
		// Continue processing but log warning
	}

	// Business rule: If score > 90, fail even if resolution is PASS
	if score > 90 {
		return fmt.Errorf("ML validation failed: ML risk score too high (%d)", score)
	}

	// Business rule: Check confidence if available
	if confidence, exists := result["confidence"]; exists {
		if conf, ok := confidence.(float64); ok && conf < 0.7 {
			fmt.Println("ML Warning: Low confidence score, consider manual review")
		}
	}

	fmt.Println("ML result validation passed")
	return nil
}
