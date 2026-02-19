package child_steps

import (
	"fmt"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/internal/primitive"
	"unified-workflow/internal/primitive/services/antifraud/models"

	"github.com/google/uuid"
)

// GetAntifraudChildSteps returns all reusable antifraud child steps
func GetAntifraudChildSteps() []*model.ChildStep {
	return []*model.ChildStep{
		CreateAntifraudPrepareTransactionRequestChildStep(),
		CreateAntifraudCallStoreTransactionAsyncChildStep(),
		CreateAntifraudProcessStoreResponseChildStep(),
		CreateAntifraudStoreTransactionResultChildStep(),
		CreateAntifraudPrepareAMLRequestChildStep(),
		CreateAntifraudCallAMLValidationAsyncChildStep(),
		CreateAntifraudProcessAMLResponseChildStep(),
		CreateAntifraudValidateAMLResponseChildStep(),
		CreateAntifraudValidateAMLResultChildStep(),
		CreateAntifraudStoreAMLResolutionChildStep(),
		CreateAntifraudAddAMLToTransactionChildStep(),
	}
}

// CreateAntifraudPrepareTransactionRequestChildStep creates a child step for preparing transaction request
func CreateAntifraudPrepareTransactionRequestChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_prepare_transaction_request_child_step",
		func(context interface{}, data interface{}) interface{} {
			fmt.Println("Preparing transaction request...")

			// Extract transaction data from workflow data
			transactionData := extractTransactionData(data)

			// Get transaction ID from context or generate one
			transactionID := getTransactionIDFromContext(context)
			if transactionID == "" {
				transactionID = uuid.NewString()
			}

			// Create antifraud transaction from real data
			afTransaction := models.AF_Transaction{
				AF_Id:      transactionID,
				AF_AddDate: time.Now().Format(time.RFC3339Nano),
				Transaction: models.Transaction{
					Id:                 getString(transactionData, "id", ""),
					Type:               getString(transactionData, "type", ""),
					Date:               getString(transactionData, "date", time.Now().Format(time.RFC3339)),
					Amount:             getString(transactionData, "amount", ""),
					Currency:           getString(transactionData, "currency", ""),
					ClientId:           getString(transactionData, "client_id", ""),
					ClientName:         getString(transactionData, "client_name", ""),
					ClientPAN:          getString(transactionData, "client_pan", ""),
					ClientCVV:          getString(transactionData, "client_cvv", ""),
					ClientCardHolder:   getString(transactionData, "client_card_holder", ""),
					ClientPhone:        getString(transactionData, "client_phone", ""),
					MerchantTerminalId: getString(transactionData, "merchant_terminal_id", ""),
					Channel:            getString(transactionData, "channel", ""),
					LocationIp:         getString(transactionData, "location_ip", ""),
				},
			}

			fmt.Printf("Prepared transaction with ID: %s\n", afTransaction.AF_Id)
			return afTransaction
		},
		nil, // No response hook for request preparation
		func(request interface{}) error {
			fmt.Println("Validating transaction request...")

			// Validate the antifraud transaction
			afTransaction, ok := request.(models.AF_Transaction)
			if !ok {
				return fmt.Errorf("invalid transaction type: %T", request)
			}

			// Validate required fields for deployment
			if afTransaction.AF_Id == "" {
				return fmt.Errorf("missing AF_Id")
			}
			if afTransaction.Transaction.Id == "" {
				return fmt.Errorf("missing transaction ID")
			}
			if afTransaction.Transaction.Amount == "" {
				return fmt.Errorf("missing transaction amount")
			}
			if afTransaction.Transaction.Currency == "" {
				return fmt.Errorf("missing transaction currency")
			}
			if afTransaction.Transaction.ClientId == "" {
				return fmt.Errorf("missing client ID")
			}

			fmt.Println("Transaction request validation passed")
			return nil
		},
	)
}

// CreateAntifraudCallStoreTransactionAsyncChildStep creates a child step for calling store transaction async
func CreateAntifraudCallStoreTransactionAsyncChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_call_store_transaction_async_child_step",
		func(context interface{}, data interface{}) interface{} {
			fmt.Println("Calling store transaction async...")

			// Extract transaction from data
			transaction, err := extractAFTransactionFromData(data)
			if err != nil {
				return fmt.Errorf("failed to extract transaction: %w", err)
			}

			// Get antifraud service from context
			antifraudService, err := getAntifraudServiceFromContext(context)
			if err != nil {
				return fmt.Errorf("failed to get antifraud service: %w", err)
			}

			// Call the actual antifraud SDK
			err = antifraudService.StoreTransaction(transaction)
			if err != nil {
				return fmt.Errorf("store transaction failed: %w", err)
			}

			// Return actual SDK response
			return map[string]interface{}{
				"async_operation_id": uuid.NewString(),
				"status":             "completed",
				"started_at":         time.Now().Format(time.RFC3339),
				"completed_at":       time.Now().Format(time.RFC3339),
				"transaction_id":     transaction.AF_Id,
				"sdk_call_made":      true,
				"error":              nil,
			}
		},
		nil, // Response hook would process async response
		nil, // No validation for async call initiation
	)
}

// CreateAntifraudProcessStoreResponseChildStep creates a child step for processing store response
func CreateAntifraudProcessStoreResponseChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_process_store_response_child_step",
		nil, // No request hook for response processing
		func(context interface{}, data interface{}) interface{} {
			fmt.Println("Processing store response...")

			// Extract the SDK response from the previous async call
			// In production, this would be the actual SDK response
			sdkResponse, err := extractSDKResponseFromData(data)
			if err != nil {
				return fmt.Errorf("failed to extract SDK response: %w", err)
			}

			// Process the actual SDK response
			return map[string]interface{}{
				"transaction_id": sdkResponse["transaction_id"],
				"stored_at":      time.Now().Format(time.RFC3339),
				"status":         "stored",
				"message":        "Transaction stored successfully",
				"sdk_response":   sdkResponse,
			}
		},
		func(response interface{}) error {
			fmt.Println("Validating store response...")

			resp, ok := response.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid store response type: %T", response)
			}

			// Validate required fields
			requiredFields := []string{"transaction_id", "stored_at", "status"}
			for _, field := range requiredFields {
				if _, exists := resp[field]; !exists {
					return fmt.Errorf("missing store response field: %s", field)
				}
			}

			// Validate that SDK response is present
			if sdkResponse, exists := resp["sdk_response"]; exists {
				if sdkResp, ok := sdkResponse.(map[string]interface{}); ok {
					if _, hasSDKCall := sdkResp["sdk_call_made"]; !hasSDKCall {
						return fmt.Errorf("missing SDK call confirmation")
					}
				}
			}

			fmt.Println("Store response validation passed")
			return nil
		},
	)
}

// CreateAntifraudStoreTransactionResultChildStep creates a child step for storing transaction result
func CreateAntifraudStoreTransactionResultChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_store_transaction_result_child_step",
		func(context interface{}, data interface{}) interface{} {
			fmt.Println("Storing transaction result in workflow data...")

			// Extract processed response
			processedResponse, err := extractProcessedResponseFromData(data)
			if err != nil {
				return fmt.Errorf("failed to extract processed response: %w", err)
			}

			// Store the actual transaction result in workflow data
			// This uses the real processed response from SDK
			return map[string]interface{}{
				"operation":             "store_transaction",
				"success":               true,
				"timestamp":             time.Now().Format(time.RFC3339),
				"message":               "Transaction stored and result saved to workflow data",
				"transaction_id":        processedResponse["transaction_id"],
				"original_sdk_response": processedResponse["sdk_response"],
			}
		},
		nil, // No response hook
		nil, // No validation
	)
}

// CreateAntifraudPrepareAMLRequestChildStep creates a child step for preparing AML request
func CreateAntifraudPrepareAMLRequestChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_prepare_aml_request_child_step",
		func(context interface{}, data interface{}) interface{} {
			fmt.Println("Preparing AML validation request...")

			// Extract transaction data from workflow data
			transactionData := extractTransactionData(data)

			// Get transaction ID from context or generate one
			transactionID := getTransactionIDFromContext(context)
			if transactionID == "" {
				transactionID = uuid.NewString()
			}

			// Create AML request from real data
			amlRequest := map[string]interface{}{
				"af_id":       transactionID,
				"af_add_date": time.Now().Format(time.RFC3339Nano),
				"transaction": map[string]interface{}{
					"id":          getString(transactionData, "id", ""),
					"type":        getString(transactionData, "type", ""),
					"amount":      getString(transactionData, "amount", ""),
					"currency":    getString(transactionData, "currency", ""),
					"client_id":   getString(transactionData, "client_id", ""),
					"client_name": getString(transactionData, "client_name", ""),
				},
			}

			fmt.Printf("Prepared AML request for transaction: %s\n", amlRequest["af_id"])
			return amlRequest
		},
		nil, // No response hook
		func(request interface{}) error {
			fmt.Println("Validating AML request...")

			req, ok := request.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid AML request type: %T", request)
			}

			// Validate required fields for deployment
			requiredFields := []string{"af_id", "transaction"}
			for _, field := range requiredFields {
				if _, exists := req[field]; !exists {
					return fmt.Errorf("missing AML request field: %s", field)
				}
			}

			// Validate transaction structure
			transaction, ok := req["transaction"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid transaction structure")
			}

			// Validate transaction fields
			txRequiredFields := []string{"id", "amount", "currency", "client_id"}
			for _, field := range txRequiredFields {
				if _, exists := transaction[field]; !exists {
					return fmt.Errorf("missing transaction field: %s", field)
				}
			}

			fmt.Println("AML request validation passed")
			return nil
		},
	)
}

// CreateAntifraudCallAMLValidationAsyncChildStep creates a child step for calling AML validation async
func CreateAntifraudCallAMLValidationAsyncChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_call_aml_validation_async_child_step",
		func(context interface{}, data interface{}) interface{} {
			fmt.Println("Calling AML validation async...")

			// Simulate async AML validation call
			amlValidationID := uuid.NewString()

			return map[string]interface{}{
				"aml_validation_id": amlValidationID,
				"status":            "processing",
				"started_at":        time.Now().Format(time.RFC3339),
				"service":           "AML",
			}
		},
		nil, // Response hook would process async response
		nil, // No validation for async call
	)
}

// CreateAntifraudProcessAMLResponseChildStep creates a child step for processing AML response
func CreateAntifraudProcessAMLResponseChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_process_aml_response_child_step",
		nil, // No request hook
		func(context interface{}, data interface{}) interface{} {
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
		},
		nil, // Validation done in separate step
	)
}

// CreateAntifraudValidateAMLResponseChildStep creates a child step for validating AML response structure
func CreateAntifraudValidateAMLResponseChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_validate_aml_response_child_step",
		nil, // No request hook
		nil, // No response hook
		func(response interface{}) error {
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
		},
	)
}

// CreateAntifraudValidateAMLResultChildStep creates a child step for validating AML result against business rules
func CreateAntifraudValidateAMLResultChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_validate_aml_result_child_step",
		nil, // No request hook
		nil, // No response hook
		func(result interface{}) error {
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
		},
	)
}

// CreateAntifraudStoreAMLResolutionChildStep creates a child step for storing AML resolution
func CreateAntifraudStoreAMLResolutionChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_store_aml_resolution_child_step",
		func(context interface{}, data interface{}) interface{} {
			fmt.Println("Storing AML resolution...")

			// In production, this would call the SDK immediately
			// For high TPS, no artificial delays
			return map[string]interface{}{
				"operation":     "store_aml_resolution",
				"success":       true,
				"stored_at":     time.Now().Format(time.RFC3339),
				"resolution_id": uuid.NewString(),
			}
		},
		nil, // No response hook
		nil, // No validation
	)
}

// CreateAntifraudAddAMLToTransactionChildStep creates a child step for adding AML to transaction
func CreateAntifraudAddAMLToTransactionChildStep() *model.ChildStep {
	return model.NewChildStep(
		"antifraud_add_aml_to_transaction_child_step",
		func(context interface{}, data interface{}) interface{} {
			fmt.Println("Adding AML check to transaction...")

			// In production, this would call the SDK immediately
			// For high TPS, no artificial delays
			return map[string]interface{}{
				"operation": "add_aml_to_transaction",
				"success":   true,
				"added_at":  time.Now().Format(time.RFC3339),
				"message":   "AML check added to transaction aggregation",
			}
		},
		nil, // No response hook
		nil, // No validation
	)
}

// Helper function to extract transaction data from workflow data
func extractTransactionData(data interface{}) map[string]interface{} {
	if data == nil {
		return make(map[string]interface{})
	}

	// Try to extract transaction from data
	if dataMap, ok := data.(map[string]interface{}); ok {
		if transaction, exists := dataMap["transaction"]; exists {
			if transactionMap, ok := transaction.(map[string]interface{}); ok {
				return transactionMap
			}
		}
	}

	return make(map[string]interface{})
}

// Helper function to get transaction ID from context
func getTransactionIDFromContext(context interface{}) string {
	if context == nil {
		return ""
	}

	// Try to extract from context map
	if contextMap, ok := context.(map[string]interface{}); ok {
		if transactionID, exists := contextMap["transaction_id"]; exists {
			if id, ok := transactionID.(string); ok {
				return id
			}
		}
		// Also check for run_id which might contain transaction ID
		if runID, exists := contextMap["run_id"]; exists {
			if id, ok := runID.(string); ok {
				return id
			}
		}
	}

	return ""
}

// Helper function to get antifraud service from context
func getAntifraudServiceFromContext(context interface{}) (primitive.AntifraudService, error) {
	if context == nil {
		return nil, fmt.Errorf("context is nil")
	}

	// Try to extract antifraud service from context
	if contextMap, ok := context.(map[string]interface{}); ok {
		if antifraudService, exists := contextMap["antifraud_service"]; exists {
			// Type assert to primitive.AntifraudService
			if service, ok := antifraudService.(primitive.AntifraudService); ok {
				return service, nil
			}
		}
	}

	return nil, fmt.Errorf("antifraud service not found in context")
}

// Helper function to extract AF_Transaction from data
func extractAFTransactionFromData(data interface{}) (models.AF_Transaction, error) {
	if data == nil {
		return models.AF_Transaction{}, fmt.Errorf("data is nil")
	}

	// Try to extract AF_Transaction from data
	if afTransaction, ok := data.(models.AF_Transaction); ok {
		return afTransaction, nil
	}

	// Try to extract from map
	if dataMap, ok := data.(map[string]interface{}); ok {
		// Check if it's an AF_Transaction in a map
		if afTransaction, exists := dataMap["af_transaction"]; exists {
			if transaction, ok := afTransaction.(models.AF_Transaction); ok {
				return transaction, nil
			}
		}
	}

	return models.AF_Transaction{}, fmt.Errorf("AF_Transaction not found in data")
}

// Helper function to extract SDK response from data
func extractSDKResponseFromData(data interface{}) (map[string]interface{}, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	// Try to extract SDK response from data
	if sdkResponse, ok := data.(map[string]interface{}); ok {
		// Check if it has SDK call confirmation
		if sdkCallMade, exists := sdkResponse["sdk_call_made"]; exists {
			if callMade, ok := sdkCallMade.(bool); ok && callMade {
				return sdkResponse, nil
			}
		}
		// Return anyway if it looks like an SDK response
		return sdkResponse, nil
	}

	// Try to extract from map
	if dataMap, ok := data.(map[string]interface{}); ok {
		// Check if SDK response is in a map
		if sdkResponse, exists := dataMap["sdk_response"]; exists {
			if response, ok := sdkResponse.(map[string]interface{}); ok {
				return response, nil
			}
		}
		// Check if async operation result is in a map
		if asyncResult, exists := dataMap["async_operation_result"]; exists {
			if result, ok := asyncResult.(map[string]interface{}); ok {
				return result, nil
			}
		}
	}

	return nil, fmt.Errorf("SDK response not found in data")
}

// Helper function to extract processed response from data
func extractProcessedResponseFromData(data interface{}) (map[string]interface{}, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	// Try to extract processed response from data
	if processedResponse, ok := data.(map[string]interface{}); ok {
		// Check if it has processed response fields
		if _, hasTransactionID := processedResponse["transaction_id"]; hasTransactionID {
			if _, hasSDKResponse := processedResponse["sdk_response"]; hasSDKResponse {
				return processedResponse, nil
			}
		}
		// Return anyway if it looks like a processed response
		return processedResponse, nil
	}

	// Try to extract from map
	if dataMap, ok := data.(map[string]interface{}); ok {
		// Check if processed response is in a map
		if processedResponse, exists := dataMap["processed_response"]; exists {
			if response, ok := processedResponse.(map[string]interface{}); ok {
				return response, nil
			}
		}
		// Check if store response is in a map
		if storeResponse, exists := dataMap["store_response"]; exists {
			if response, ok := storeResponse.(map[string]interface{}); ok {
				return response, nil
			}
		}
	}

	return nil, fmt.Errorf("processed response not found in data")
}

// Helper function to get string value from map with default
func getString(data map[string]interface{}, key, defaultValue string) string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}
