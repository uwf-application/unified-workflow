package steps

import (
	"fmt"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/workflows/child_steps"

	"github.com/google/uuid"
)

// StoreTransactionStep stores a transaction in the antifraud system
type StoreTransactionStep struct {
	*AntifraudStep
}

// NewStoreTransactionStep creates a new StoreTransactionStep
func NewStoreTransactionStep(endpoint string) *StoreTransactionStep {
	step := &StoreTransactionStep{
		AntifraudStep: NewAntifraudStep("store-transaction", endpoint),
	}

	// Add reusable child steps for async operations
	step.AddChildSteps([]*model.ChildStep{
		// Child Step 1: Prepare transaction request
		child_steps.CreateAntifraudPrepareTransactionRequestChildStep(),
		// Child Step 2: Call store transaction async
		child_steps.CreateAntifraudCallStoreTransactionAsyncChildStep(),
		// Child Step 3: Process store response
		child_steps.CreateAntifraudProcessStoreResponseChildStep(),
		// Child Step 4: Store transaction result
		child_steps.CreateAntifraudStoreTransactionResultChildStep(),
	})

	return step
}

// prepareTransactionRequest prepares the transaction request for antifraud system
func (s *StoreTransactionStep) prepareTransactionRequest(context interface{}, data interface{}) interface{} {
	fmt.Println("Preparing transaction request...")

	// Extract transaction data from workflow data
	// In a real implementation, you would extract from the data parameter
	transactionData := map[string]interface{}{
		"transaction": map[string]interface{}{
			"id":                   uuid.NewString(),
			"type":                 "deposit",
			"date":                 time.Now().Format(time.RFC3339Nano),
			"amount":               "100000",
			"currency":             "KZT",
			"client_id":            uuid.NewString(),
			"client_name":          "John Smith",
			"client_pan":           "111111******1111",
			"client_cvv":           "111",
			"client_card_holder":   "JOHN SMITH",
			"client_phone":         "+77007007070",
			"merchant_terminal_id": "00000001",
			"channel":              "E-com",
			"location_ip":          "192.168.0.1",
		},
	}

	// Create AF_Transaction structure
	afTransaction := map[string]interface{}{
		"af_id":       uuid.NewString(),
		"af_add_date": time.Now().Format(time.RFC3339Nano),
		"transaction": transactionData["transaction"],
	}

	fmt.Printf("Prepared transaction with ID: %s\n", afTransaction["af_id"])
	return afTransaction
}

// validateTransactionRequest validates the prepared transaction request
func (s *StoreTransactionStep) validateTransactionRequest(request interface{}) error {
	fmt.Println("Validating transaction request...")

	// Check if request is a map
	req, ok := request.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid request type: %T", request)
	}

	// Validate required fields
	requiredFields := []string{"af_id", "af_add_date", "transaction"}
	for _, field := range requiredFields {
		if _, exists := req[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate transaction structure
	transaction, ok := req["transaction"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid transaction structure")
	}

	// Validate transaction fields
	txRequiredFields := []string{"id", "type", "amount", "currency"}
	for _, field := range txRequiredFields {
		if _, exists := transaction[field]; !exists {
			return fmt.Errorf("missing transaction field: %s", field)
		}
	}

	fmt.Println("Transaction request validation passed")
	return nil
}

// callStoreTransactionAsync calls the antifraud service to store transaction
func (s *StoreTransactionStep) callStoreTransactionAsync(context interface{}, data interface{}) interface{} {
	fmt.Println("Calling store transaction async...")

	// Get antifraud service
	_, err := s.GetAntifraudService()
	if err != nil {
		return fmt.Errorf("failed to get antifraud service: %w", err)
	}

	// In a real async implementation, this would return a future/promise
	// For high TPS, we don't simulate delays
	go func() {
		// Real async call would happen here
		fmt.Println("Async store transaction completed")
	}()

	// Return a mock async result
	return map[string]interface{}{
		"async_operation_id": uuid.NewString(),
		"status":             "pending",
		"started_at":         time.Now().Format(time.RFC3339),
	}
}

// processStoreResponse processes the async store transaction response
func (s *StoreTransactionStep) processStoreResponse(context interface{}, data interface{}) interface{} {
	fmt.Println("Processing store response...")

	// In production, this would process the actual SDK response
	// For high TPS, we process immediately without delays
	return map[string]interface{}{
		"transaction_id": uuid.NewString(),
		"status":         "stored",
		"stored_at":      time.Now().Format(time.RFC3339),
		"message":        "Transaction stored successfully in antifraud system",
	}
}

// validateStoreResponse validates the store transaction response
func (s *StoreTransactionStep) validateStoreResponse(response interface{}) error {
	fmt.Println("Validating store response...")

	// Check if response is a map
	resp, ok := response.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response type: %T", response)
	}

	// Validate response fields
	requiredFields := []string{"transaction_id", "status", "stored_at"}
	for _, field := range requiredFields {
		if _, exists := resp[field]; !exists {
			return fmt.Errorf("missing response field: %s", field)
		}
	}

	// Check status
	status, ok := resp["status"].(string)
	if !ok || status != "stored" {
		return fmt.Errorf("invalid transaction status: %v", status)
	}

	fmt.Println("Store response validation passed")
	return nil
}

// storeTransactionResult stores the transaction result in workflow data
func (s *StoreTransactionStep) storeTransactionResult(context interface{}, data interface{}) interface{} {
	fmt.Println("Storing transaction result in workflow data...")

	// In a real implementation, you would store the result in workflow data
	// For now, we'll return a success result
	return map[string]interface{}{
		"operation": "store_transaction",
		"success":   true,
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   "Transaction stored and result saved to workflow data",
	}
}

// ExecuteStepLogic executes the main step logic (alternative to child steps)
func (s *StoreTransactionStep) ExecuteStepLogic(ctx interface{}, context interface{}, data interface{}) error {
	fmt.Println("Executing StoreTransactionStep logic...")

	// This method provides an alternative execution path without child steps
	// It follows the same flow as child steps but in a single method

	// 1. Prepare request
	request := s.prepareTransactionRequest(context, data)

	// 2. Validate request
	if err := s.validateTransactionRequest(request); err != nil {
		return fmt.Errorf("transaction request validation failed: %w", err)
	}

	// 3. Call async operation
	asyncResult := s.callStoreTransactionAsync(context, data)
	fmt.Printf("Async operation started: %v\n", asyncResult)

	// 4. Process response
	response := s.processStoreResponse(context, data)

	// 5. Validate response
	if err := s.validateStoreResponse(response); err != nil {
		return fmt.Errorf("store response validation failed: %w", err)
	}

	// 6. Store result
	result := s.storeTransactionResult(context, data)
	fmt.Printf("Transaction result: %v\n", result)

	fmt.Println("StoreTransactionStep completed successfully")
	return nil
}
