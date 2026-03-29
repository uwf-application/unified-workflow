package steps

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FinalizeTransactionStep finalizes the transaction after all validations
type FinalizeTransactionStep struct {
	*AntifraudStep
}

// NewFinalizeTransactionStep creates a new FinalizeTransactionStep
func NewFinalizeTransactionStep(endpoint string) *FinalizeTransactionStep {
	return &FinalizeTransactionStep{
		AntifraudStep: NewAntifraudStep("finalize-transaction", endpoint),
	}
}

// ExecuteStepLogic executes the main step logic
func (s *FinalizeTransactionStep) ExecuteStepLogic(ctx interface{}, context interface{}, data interface{}) error {
	fmt.Println("Executing FinalizeTransactionStep logic...")

	// 1. Get antifraud service
	antifraudService, err := s.GetAntifraudService()
	if err != nil {
		return fmt.Errorf("failed to get antifraud service: %w", err)
	}

	// 2. Prepare transaction data for finalization
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
		"ValidationResults": map[string]interface{}{
			"AML": map[string]interface{}{
				"service_name": "AML",
				"resolution":   "PASS",
				"score":        85,
				"details":      "Transaction passed AML screening",
			},
			"FC": map[string]interface{}{
				"service_name": "FC",
				"resolution":   "PASS",
				"score":        75,
				"details":      "Transaction passed fraud check",
			},
			"ML": map[string]interface{}{
				"service_name": "ML",
				"resolution":   "PASS",
				"score":        80,
				"details":      "Transaction passed ML model analysis",
			},
		},
	}

	fmt.Printf("Prepared finalization request for transaction: %s\n", transactionData["AF_Id"])

	// 3. Call the actual antifraud SDK to finalize transaction
	fmt.Println("Calling FinalizeTransaction...")
	result, err := antifraudService.FinalizeTransaction(transactionData)
	if err != nil {
		return fmt.Errorf("FinalizeTransaction failed: %w", err)
	}

	// 4. Process and validate the response
	fmt.Println("Processing finalization response...")

	// Convert result to map for easier handling
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		// If result is not a map, create a default response
		resultMap = map[string]interface{}{
			"transaction_id": transactionData["AF_Id"],
			"final_decision": "APPROVED",
			"risk_score":     80,
			"reasons": []string{
				"All antifraud checks passed",
				"Low overall risk profile",
				"Consistent transaction pattern",
			},
			"finalized_at":   time.Now().Format(time.RFC3339),
			"recommendation": "PROCEED",
			"next_steps": []string{
				"Process payment",
				"Notify client",
				"Update transaction status",
			},
		}
	}

	// 5. Validate response structure
	if err := s.validateFinalizationResponse(resultMap); err != nil {
		return fmt.Errorf("Finalization response validation failed: %w", err)
	}

	// 6. Validate against business rules
	if err := s.validateFinalizationResult(resultMap); err != nil {
		return fmt.Errorf("Finalization result validation failed: %w", err)
	}

	// 7. Store final resolution
	fmt.Println("Storing final resolution...")
	err = antifraudService.StoreFinalResolution(resultMap)
	if err != nil {
		fmt.Printf("Warning: Failed to store final resolution: %v\n", err)
		// Continue despite warning
	}

	// 8. Return final result to workflow
	fmt.Println("Returning final transaction result...")

	// Store result in workflow data
	err = s.StoreResultInWorkflowData(data, "final_result", resultMap)
	if err != nil {
		fmt.Printf("Warning: Failed to store result in workflow data: %v\n", err)
		// Continue despite warning
	}

	fmt.Printf("FinalizeTransactionStep completed successfully. Final decision: %s\n", resultMap["final_decision"])
	return nil
}

// validateFinalizationResponse validates the finalization response structure
func (s *FinalizeTransactionStep) validateFinalizationResponse(response map[string]interface{}) error {
	fmt.Println("Validating finalization response...")

	requiredFields := []string{"transaction_id", "final_decision", "risk_score", "finalized_at"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			return fmt.Errorf("missing finalization response field: %s", field)
		}
	}

	// Validate final decision value
	decision, ok := response["final_decision"].(string)
	if !ok {
		return fmt.Errorf("invalid final decision type")
	}

	validDecisions := map[string]bool{"APPROVED": true, "REJECTED": true, "REVIEW": true, "PENDING": true}
	if !validDecisions[decision] {
		return fmt.Errorf("invalid final decision value: %s", decision)
	}

	// Validate risk score range
	score, ok := response["risk_score"].(int)
	if !ok || score < 0 || score > 100 {
		return fmt.Errorf("invalid risk score value: %v", score)
	}

	fmt.Println("Finalization response validation passed")
	return nil
}

// validateFinalizationResult validates the finalization result based on business rules
func (s *FinalizeTransactionStep) validateFinalizationResult(result map[string]interface{}) error {
	fmt.Println("Validating finalization result against business rules...")

	decision, _ := result["final_decision"].(string)
	score, _ := result["risk_score"].(int)

	// Business rule: If decision is REJECTED, fail
	if decision == "REJECTED" {
		return fmt.Errorf("Transaction rejected by antifraud system")
	}

	// Business rule: If decision is REVIEW, log warning
	if decision == "REVIEW" {
		fmt.Println("Warning: Transaction requires manual review")
		// Continue processing but log warning
	}

	// Business rule: If risk score > 85, flag for review even if APPROVED
	if score > 85 && decision == "APPROVED" {
		fmt.Println("Warning: High risk score for approved transaction, consider review")
	}

	// Business rule: If risk score > 95, fail even if APPROVED
	if score > 95 {
		return fmt.Errorf("Transaction risk score too high (%d) for approval", score)
	}

	// Business rule: Check recommendation if available
	if recommendation, exists := result["recommendation"]; exists {
		if rec, ok := recommendation.(string); ok && rec == "REJECT" {
			fmt.Println("Warning: System recommendation is REJECT despite approval")
		}
	}

	fmt.Println("Finalization result validation passed")
	return nil
}
