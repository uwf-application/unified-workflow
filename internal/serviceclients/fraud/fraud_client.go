package fraud

import (
	"fmt"
	"time"

	"unified-workflow/internal/primitive/services/fraud"
)

// FraudClient implements the FraudService interface with actual business logic
type FraudClient struct {
	apiKey     string
	apiURL     string
	httpClient interface{} // Would be *http.Client in real implementation
	config     *Config
}

// Config holds configuration for the fraud client
type Config struct {
	APIKey        string        `json:"api_key"`
	APIURL        string        `json:"api_url"`
	Timeout       time.Duration `json:"timeout"`
	MaxRetries    int           `json:"max_retries"`
	RiskThreshold int           `json:"risk_threshold"` // Default: 75
	EnableMock    bool          `json:"enable_mock"`
	MockResponses MockConfig    `json:"mock_responses,omitempty"`
}

// MockConfig holds mock response configurations for testing
type MockConfig struct {
	DefaultRiskScore int    `json:"default_risk_score"`
	DefaultAction    string `json:"default_action"` // "allow", "review", "block"
}

// NewFraudClient creates a new fraud client
func NewFraudClient(apiKey, apiURL string) (*FraudClient, error) {
	return &FraudClient{
		apiKey: apiKey,
		apiURL: apiURL,
		config: &Config{
			APIKey:        apiKey,
			APIURL:        apiURL,
			Timeout:       30 * time.Second,
			MaxRetries:    3,
			RiskThreshold: 75,
			EnableMock:    false,
		},
	}, nil
}

// NewFraudClientWithConfig creates a new fraud client with custom configuration
func NewFraudClientWithConfig(config *Config) (*FraudClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if config.APIURL == "" {
		config.APIURL = "https://api.fraudservice.com/v1"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RiskThreshold == 0 {
		config.RiskThreshold = 75
	}

	return &FraudClient{
		apiKey: config.APIKey,
		apiURL: config.APIURL,
		config: config,
	}, nil
}

// CheckSuspiciousTransaction implements FraudService.CheckSuspiciousTransaction
func (c *FraudClient) CheckSuspiciousTransaction(input fraud.CheckSuspiciousTransactionInput) (fraud.CheckSuspiciousTransactionOutput, error) {
	if c.config.EnableMock {
		return c.mockCheckSuspiciousTransaction(input)
	}

	// In a real implementation, this would call an external fraud detection API
	// For now, we'll implement a simple rule-based check
	riskScore := c.calculateRiskScore(input)
	isSuspicious := riskScore >= c.config.RiskThreshold
	action := "allow"
	if isSuspicious {
		if riskScore >= 90 {
			action = "block"
		} else {
			action = "review"
		}
	}

	return fraud.CheckSuspiciousTransactionOutput{
		IsSuspicious: isSuspicious,
		RiskScore:    riskScore,
		Reason:       c.generateReason(input, riskScore),
		Action:       action,
		Confidence:   c.calculateConfidence(riskScore),
	}, nil
}

// FlagTransaction implements FraudService.FlagTransaction
func (c *FraudClient) FlagTransaction(input fraud.FlagTransactionInput) (fraud.FlagTransactionOutput, error) {
	if c.config.EnableMock {
		return c.mockFlagTransaction(input)
	}

	// In a real implementation, this would call an external API to flag the transaction
	// For now, we'll simulate a successful flagging
	flagID := fmt.Sprintf("flag_%d", time.Now().UnixNano())

	return fraud.FlagTransactionOutput{
		FlagID:         flagID,
		Success:        true,
		Message:        fmt.Sprintf("Transaction %s flagged as %s: %s", input.TransactionID, input.Severity, input.Reason),
		FlaggedAt:      time.Now().UTC().Format(time.RFC3339),
		ReviewRequired: input.Severity == "high" || input.Severity == "critical",
	}, nil
}

// GetRiskScore implements FraudService.GetRiskScore
func (c *FraudClient) GetRiskScore(input fraud.GetRiskScoreInput) (fraud.GetRiskScoreOutput, error) {
	if c.config.EnableMock {
		return c.mockGetRiskScore(input)
	}

	// In a real implementation, this would fetch risk score from external service
	// For now, we'll return a simulated risk score
	riskScore := 45 // Default medium-low risk
	riskLevel := "low"
	if riskScore >= 70 {
		riskLevel = "high"
	} else if riskScore >= 40 {
		riskLevel = "medium"
	}

	factors := []fraud.RiskFactor{
		{
			Factor:     "transaction_history",
			Score:      30,
			Confidence: 0.85,
			Details:    "Customer has 5 transactions in last 30 days",
		},
		{
			Factor:     "device_trust",
			Score:      15,
			Confidence: 0.75,
			Details:    "Device recognized from previous logins",
		},
	}

	return fraud.GetRiskScoreOutput{
		CustomerID:  input.CustomerID,
		RiskScore:   riskScore,
		RiskLevel:   riskLevel,
		Factors:     factors,
		LastUpdated: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// BatchCheck implements FraudService.BatchCheck
func (c *FraudClient) BatchCheck(inputs []fraud.CheckSuspiciousTransactionInput) ([]fraud.CheckSuspiciousTransactionOutput, error) {
	outputs := make([]fraud.CheckSuspiciousTransactionOutput, len(inputs))
	for i, input := range inputs {
		output, err := c.CheckSuspiciousTransaction(input)
		if err != nil {
			return nil, fmt.Errorf("error checking transaction %s: %w", input.TransactionID, err)
		}
		outputs[i] = output
	}
	return outputs, nil
}

// Helper methods for rule-based fraud detection
func (c *FraudClient) calculateRiskScore(input fraud.CheckSuspiciousTransactionInput) int {
	score := 0

	// Amount-based risk
	if input.Amount > 10000 {
		score += 40
	} else if input.Amount > 5000 {
		score += 25
	} else if input.Amount > 1000 {
		score += 15
	}

	// Currency risk (simplified)
	if input.Currency != "USD" && input.Currency != "EUR" {
		score += 10
	}

	// Location risk (simplified)
	if input.Location == "high_risk_country" {
		score += 30
	}

	// IP risk (simplified)
	if input.IPAddress == "suspicious_ip" {
		score += 25
	}

	// Ensure score is between 0-100
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

func (c *FraudClient) generateReason(input fraud.CheckSuspiciousTransactionInput, riskScore int) string {
	if riskScore >= 90 {
		return "High risk transaction based on amount, location, and IP address"
	} else if riskScore >= 75 {
		return "Medium risk transaction requiring review"
	} else if riskScore >= 50 {
		return "Low risk transaction, routine check passed"
	}
	return "Transaction appears normal"
}

func (c *FraudClient) calculateConfidence(riskScore int) float64 {
	// Higher confidence for extreme scores, lower for middle scores
	if riskScore >= 90 || riskScore <= 10 {
		return 0.95
	} else if riskScore >= 75 || riskScore <= 25 {
		return 0.85
	}
	return 0.75
}

// Mock implementations for testing
func (c *FraudClient) mockCheckSuspiciousTransaction(input fraud.CheckSuspiciousTransactionInput) (fraud.CheckSuspiciousTransactionOutput, error) {
	return fraud.CheckSuspiciousTransactionOutput{
		IsSuspicious: c.config.MockResponses.DefaultRiskScore >= c.config.RiskThreshold,
		RiskScore:    c.config.MockResponses.DefaultRiskScore,
		Reason:       "Mock response for testing",
		Action:       c.config.MockResponses.DefaultAction,
		Confidence:   0.95,
	}, nil
}

func (c *FraudClient) mockFlagTransaction(input fraud.FlagTransactionInput) (fraud.FlagTransactionOutput, error) {
	return fraud.FlagTransactionOutput{
		FlagID:         "mock_flag_123",
		Success:        true,
		Message:        "Mock flagging successful",
		FlaggedAt:      time.Now().UTC().Format(time.RFC3339),
		ReviewRequired: true,
	}, nil
}

func (c *FraudClient) mockGetRiskScore(input fraud.GetRiskScoreInput) (fraud.GetRiskScoreOutput, error) {
	return fraud.GetRiskScoreOutput{
		CustomerID: input.CustomerID,
		RiskScore:  65,
		RiskLevel:  "medium",
		Factors: []fraud.RiskFactor{
			{
				Factor:     "mock_factor",
				Score:      65,
				Confidence: 0.9,
				Details:    "Mock data for testing",
			},
		},
		LastUpdated: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
