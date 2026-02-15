package antifraud

import (
	"fmt"
	"time"

	primitiveantifraud "unified-workflow/internal/primitive/services/antifraud"
	"unified-workflow/internal/primitive/services/antifraud/models"

	af "github.com/baraic-io/antifraud-go"
)

// antifraudClientImpl is the implementation of AntifraudService
// This wraps the actual github.com/baraic-io/antifraud-go SDK
type antifraudClientImpl struct {
	config models.ClientConfig
	client af.Client
}

// NewClient creates a new antifraud client
func NewClient(config models.ClientConfig) (primitiveantifraud.AntifraudService, error) {
	if !config.Enabled {
		return &disabledClient{config: config}, nil
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required for antifraud service")
	}

	if config.Host == "" {
		return nil, fmt.Errorf("host is required for antifraud service")
	}

	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 30 // seconds
	}

	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	// Initialize the SDK client
	afConfig := af.ClientConfig{
		Host:   config.Host,
		APIKey: config.APIKey,
	}
	client, err := af.NewClient(afConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create antifraud client: %w", err)
	}

	return &antifraudClientImpl{
		config: config,
		client: client,
	}, nil
}

// convertToSDKTransaction converts our AF_Transaction to SDK's AF_Transaction
func convertToSDKTransaction(tx models.AF_Transaction) af.AF_Transaction {
	return af.AF_Transaction{
		AF_Id:      tx.AF_Id,
		AF_AddDate: tx.AF_AddDate,
		Transaction: af.Transaction{
			Id:                 tx.Transaction.Id,
			Type:               tx.Transaction.Type,
			Date:               tx.Transaction.Date,
			Amount:             tx.Transaction.Amount,
			Currency:           tx.Transaction.Currency,
			ClientId:           tx.Transaction.ClientId,
			ClientName:         tx.Transaction.ClientName,
			ClientPAN:          tx.Transaction.ClientPAN,
			ClientCVV:          tx.Transaction.ClientCVV,
			ClientCardHolder:   tx.Transaction.ClientCardHolder,
			ClientPhone:        tx.Transaction.ClientPhone,
			MerchantTerminalId: tx.Transaction.MerchantTerminalId,
			Channel:            tx.Transaction.Channel,
			LocationIp:         tx.Transaction.LocationIp,
		},
	}
}

// convertFromSDKServiceResolution converts SDK's ServiceResolution to our ServiceResolution
func convertFromSDKServiceResolution(res af.ServiceResolution) models.ServiceResolution {
	// Convert map[string]string details to string
	detailsStr := ""
	for k, v := range res.Details {
		if detailsStr != "" {
			detailsStr += "; "
		}
		detailsStr += k + ": " + v
	}

	// Determine resolution based on fields
	resolution := "APPROVED"
	if res.Fraud {
		resolution = "FRAUD"
	} else if res.Blocked {
		resolution = "BLOCKED"
	} else if res.Alert {
		resolution = "ALERT"
	}

	// Calculate score based on various factors
	score := 0
	if res.Validated {
		score += 50
	}
	if !res.Fraud && !res.Blocked {
		score += 30
	}
	if res.InWhiteList {
		score += 20
	}

	return models.ServiceResolution{
		ServiceName: res.Service,
		Resolution:  resolution,
		Score:       score,
		Details:     detailsStr,
	}
}

// convertToSDKServiceResolution converts our ServiceResolution to SDK's ServiceResolution
func convertToSDKServiceResolution(res models.ServiceResolution) af.ServiceResolution {
	// Parse details string back to map
	details := make(map[string]string)
	// Simple parsing - in real implementation would need more robust parsing
	if res.Details != "" {
		details["summary"] = res.Details
	}

	// Determine boolean fields from resolution
	fraud := res.Resolution == "FRAUD"
	blocked := res.Resolution == "BLOCKED"
	alert := res.Resolution == "ALERT"
	validated := res.Resolution == "APPROVED"

	return af.ServiceResolution{
		Service:     res.ServiceName,
		Details:     details,
		Fraud:       fraud,
		Blocked:     blocked,
		Alert:       alert,
		Validated:   validated,
		InWhiteList: res.Score > 70, // Assume high score means in whitelist
	}
}

// convertFromSDKFinalResolution converts SDK's FinalResolution to our FinalResolution
func convertFromSDKFinalResolution(res af.FinalResolution) models.FinalResolution {
	// Extract reasons from details
	var reasons []string
	if res.Details != nil {
		for k, v := range res.Details {
			reasons = append(reasons, fmt.Sprintf("%s: %v", k, v))
		}
	}

	// Add validated services as reasons
	for _, service := range res.ValidatedServices {
		reasons = append(reasons, fmt.Sprintf("Validated: %s", service))
	}

	// Calculate risk score
	riskScore := 0
	if res.Fraud {
		riskScore = 90
	} else if res.Blocked {
		riskScore = 80
	} else if res.Alert {
		riskScore = 60
	} else if res.Validated {
		riskScore = 20
	}

	return models.FinalResolution{
		TransactionId: res.AF_Id,
		FinalDecision: res.FinalizedAction,
		RiskScore:     riskScore,
		Reasons:       reasons,
	}
}

// convertToSDKFinalResolution converts our FinalResolution to SDK's FinalResolution
func convertToSDKFinalResolution(res models.FinalResolution) af.FinalResolution {
	// Convert reasons to details map
	details := make(map[string]interface{})
	for i, reason := range res.Reasons {
		details[fmt.Sprintf("reason_%d", i)] = reason
	}

	// Determine boolean fields from final decision
	fraud := res.FinalDecision == "FRAUD"
	blocked := res.FinalDecision == "BLOCKED"
	alert := res.FinalDecision == "ALERT"
	validated := res.FinalDecision == "APPROVED"

	return af.FinalResolution{
		AF_Id:           res.TransactionId,
		FinalizedAction: res.FinalDecision,
		Details:         details,
		Fraud:           fraud,
		Blocked:         blocked,
		Alert:           alert,
		Validated:       validated,
	}
}

// StoreTransaction stores a transaction in the antifraud system
func (c *antifraudClientImpl) StoreTransaction(afTransaction models.AF_Transaction) error {
	// Convert to SDK type
	sdkTransaction := convertToSDKTransaction(afTransaction)

	// Call the SDK
	return c.client.StoreTransaction(sdkTransaction)
}

// ValidateTransactionByAML validates a transaction using the AML service
func (c *antifraudClientImpl) ValidateTransactionByAML(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	// Convert to SDK type
	sdkTransaction := convertToSDKTransaction(afTransaction)

	// Call the SDK
	result, err := c.client.ValidateTransactionByAML(sdkTransaction)
	if err != nil {
		return models.ServiceResolution{}, fmt.Errorf("AML validation failed: %w", err)
	}

	// Convert back to our type
	return convertFromSDKServiceResolution(result), nil
}

// ValidateTransactionByFC validates a transaction using the FC service
func (c *antifraudClientImpl) ValidateTransactionByFC(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	// Convert to SDK type
	sdkTransaction := convertToSDKTransaction(afTransaction)

	// Call the SDK
	result, err := c.client.ValidateTransactionByFC(sdkTransaction)
	if err != nil {
		return models.ServiceResolution{}, fmt.Errorf("FC validation failed: %w", err)
	}

	// Convert back to our type
	return convertFromSDKServiceResolution(result), nil
}

// ValidateTransactionByML validates a transaction using the ML service
func (c *antifraudClientImpl) ValidateTransactionByML(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	// Convert to SDK type
	sdkTransaction := convertToSDKTransaction(afTransaction)

	// Call the SDK
	result, err := c.client.ValidateTransactionByML(sdkTransaction)
	if err != nil {
		return models.ServiceResolution{}, fmt.Errorf("ML validation failed: %w", err)
	}

	// Convert back to our type
	return convertFromSDKServiceResolution(result), nil
}

// StoreServiceResolution stores the resolution from a service check
func (c *antifraudClientImpl) StoreServiceResolution(resolution models.ServiceResolution) error {
	// Convert to SDK type
	sdkResolution := convertToSDKServiceResolution(resolution)

	// Call the SDK
	return c.client.StoreServiceResolution(sdkResolution)
}

// AddTransactionServiceCheck adds a completed service check resolution
func (c *antifraudClientImpl) AddTransactionServiceCheck(resolution models.ServiceResolution) error {
	// Convert to SDK type
	sdkResolution := convertToSDKServiceResolution(resolution)

	// Call the SDK
	return c.client.AddTransactionServiceCheck(sdkResolution)
}

// FinalizeTransaction finalizes the transaction validation process
func (c *antifraudClientImpl) FinalizeTransaction(afTransaction models.AF_Transaction) (models.FinalResolution, error) {
	// Convert to SDK type
	sdkTransaction := convertToSDKTransaction(afTransaction)

	// Call the SDK
	result, err := c.client.FinalizeTransaction(sdkTransaction)
	if err != nil {
		return models.FinalResolution{}, fmt.Errorf("failed to finalize transaction: %w", err)
	}

	// Convert back to our type
	return convertFromSDKFinalResolution(result), nil
}

// StoreFinalResolution stores the final resolution of the transaction
func (c *antifraudClientImpl) StoreFinalResolution(resolution models.FinalResolution) error {
	// Convert to SDK type
	sdkResolution := convertToSDKFinalResolution(resolution)

	// Call the SDK
	return c.client.StoreFinalResolution(sdkResolution)
}

// HealthCheck checks the health of the antifraud service
func (c *antifraudClientImpl) HealthCheck() (bool, error) {
	// Note: The SDK doesn't have a HealthCheck method
	// We'll implement a simple health check by trying to make a lightweight API call
	// For now, we'll assume the service is healthy if we can create a simple transaction
	testTransaction := models.AF_Transaction{
		AF_Id:      "health-check",
		AF_AddDate: time.Now().Format(time.RFC3339Nano),
		Transaction: models.Transaction{
			Id:   "health-check-txn",
			Type: "health-check",
		},
	}

	// Try to store a test transaction (this should fail gracefully if service is down)
	err := c.StoreTransaction(testTransaction)
	if err != nil {
		// Check if it's an authentication error vs connection error
		if err.Error() == "antifraud service is disabled" {
			return false, fmt.Errorf("service disabled")
		}
		// For health check purposes, we consider any error as unhealthy
		return false, fmt.Errorf("health check failed: %w", err)
	}

	return true, nil
}

// GetConfig returns the current configuration
func (c *antifraudClientImpl) GetConfig() models.ClientConfig {
	return c.config
}

// disabledClient is used when the antifraud service is disabled
type disabledClient struct {
	config models.ClientConfig
}

func (c *disabledClient) StoreTransaction(afTransaction models.AF_Transaction) error {
	return fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) ValidateTransactionByAML(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	return models.ServiceResolution{}, fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) ValidateTransactionByFC(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	return models.ServiceResolution{}, fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) ValidateTransactionByML(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	return models.ServiceResolution{}, fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) StoreServiceResolution(resolution models.ServiceResolution) error {
	return fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) AddTransactionServiceCheck(resolution models.ServiceResolution) error {
	return fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) FinalizeTransaction(afTransaction models.AF_Transaction) (models.FinalResolution, error) {
	return models.FinalResolution{}, fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) StoreFinalResolution(resolution models.FinalResolution) error {
	return fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) HealthCheck() (bool, error) {
	return false, fmt.Errorf("antifraud service is disabled")
}

func (c *disabledClient) GetConfig() models.ClientConfig {
	return c.config
}
