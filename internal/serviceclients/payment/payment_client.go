package payment

import (
	"fmt"
	"time"

	"unified-workflow/internal/primitive/services/payment"
)

// PaymentClient implements the PaymentService interface with actual business logic
type PaymentClient struct {
	apiKey     string
	apiURL     string
	httpClient interface{} // Would be *http.Client in real implementation
	config     *Config
}

// Config holds configuration for the payment client
type Config struct {
	APIKey        string        `json:"api_key"`
	APIURL        string        `json:"api_url"`
	Timeout       time.Duration `json:"timeout"`
	MaxRetries    int           `json:"max_retries"`
	EnableMock    bool          `json:"enable_mock"`
	MockResponses MockConfig    `json:"mock_responses,omitempty"`
}

// MockConfig holds mock response configurations for testing
type MockConfig struct {
	DefaultStatus string  `json:"default_status"` // "completed", "failed", "pending"
	SuccessRate   float64 `json:"success_rate"`   // 0.0 to 1.0
}

// NewPaymentClient creates a new payment client
func NewPaymentClient(apiKey, apiURL string) (*PaymentClient, error) {
	return &PaymentClient{
		apiKey: apiKey,
		apiURL: apiURL,
		config: &Config{
			APIKey:     apiKey,
			APIURL:     apiURL,
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			EnableMock: false,
		},
	}, nil
}

// NewPaymentClientWithConfig creates a new payment client with custom configuration
func NewPaymentClientWithConfig(config *Config) (*PaymentClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if config.APIURL == "" {
		config.APIURL = "https://api.paymentservice.com/v1"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	return &PaymentClient{
		apiKey: config.APIKey,
		apiURL: config.APIURL,
		config: config,
	}, nil
}

// ProcessPayment implements PaymentService.ProcessPayment
func (c *PaymentClient) ProcessPayment(input payment.ProcessPaymentInput) (payment.ProcessPaymentOutput, error) {
	if c.config.EnableMock {
		return c.mockProcessPayment(input)
	}

	// In a real implementation, this would call an external payment gateway API
	// For now, we'll simulate payment processing with basic validation

	// Validate input
	if input.Amount <= 0 {
		return payment.ProcessPaymentOutput{
			TransactionID: input.TransactionID,
			Status:        "failed",
			FailureReason: "Invalid amount",
			ProcessedAt:   time.Now().UTC().Format(time.RFC3339),
		}, fmt.Errorf("invalid amount: %f", input.Amount)
	}

	if input.Currency == "" {
		return payment.ProcessPaymentOutput{
			TransactionID: input.TransactionID,
			Status:        "failed",
			FailureReason: "Currency is required",
			ProcessedAt:   time.Now().UTC().Format(time.RFC3339),
		}, fmt.Errorf("currency is required")
	}

	// Simulate payment processing
	paymentID := fmt.Sprintf("pay_%d", time.Now().UnixNano())
	status := "completed"

	// Simulate random failures (10% chance in demo)
	if time.Now().UnixNano()%10 == 0 {
		status = "failed"
	}

	output := payment.ProcessPaymentOutput{
		TransactionID:     input.TransactionID,
		PaymentID:         paymentID,
		Status:            status,
		AuthorizationCode: fmt.Sprintf("AUTH%d", time.Now().UnixNano()%1000000),
		ProcessedAt:       time.Now().UTC().Format(time.RFC3339),
		GatewayResponse: map[string]interface{}{
			"gateway":        "demo_payment_gateway",
			"transaction_id": input.TransactionID,
			"amount":         input.Amount,
			"currency":       input.Currency,
			"timestamp":      time.Now().UTC().Format(time.RFC3339),
		},
	}

	if status == "failed" {
		output.FailureReason = "Payment gateway declined transaction"
	}

	return output, nil
}

// RefundPayment implements PaymentService.RefundPayment
func (c *PaymentClient) RefundPayment(input payment.RefundPaymentInput) (payment.RefundPaymentOutput, error) {
	if c.config.EnableMock {
		return c.mockRefundPayment(input)
	}

	// Validate input
	if input.Amount <= 0 {
		return payment.RefundPaymentOutput{
			PaymentID:     input.PaymentID,
			Status:        "failed",
			FailureReason: "Invalid refund amount",
			RefundedAt:    time.Now().UTC().Format(time.RFC3339),
		}, fmt.Errorf("invalid refund amount: %f", input.Amount)
	}

	// Simulate refund processing
	refundID := fmt.Sprintf("ref_%d", time.Now().UnixNano())
	status := "completed"

	output := payment.RefundPaymentOutput{
		RefundID:       refundID,
		PaymentID:      input.PaymentID,
		Status:         status,
		RefundedAmount: input.Amount,
		RefundedAt:     time.Now().UTC().Format(time.RFC3339),
	}

	return output, nil
}

// GetPaymentStatus implements PaymentService.GetPaymentStatus
func (c *PaymentClient) GetPaymentStatus(input payment.GetPaymentStatusInput) (payment.GetPaymentStatusOutput, error) {
	if c.config.EnableMock {
		return c.mockGetPaymentStatus(input)
	}

	// In a real implementation, this would query the payment gateway
	// For now, return mock status
	createdAt := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	updatedAt := time.Now().Add(-30 * time.Minute).UTC().Format(time.RFC3339)

	return payment.GetPaymentStatusOutput{
		PaymentID:     input.PaymentID,
		TransactionID: input.TransactionID,
		Status:        "completed",
		Amount:        99.99,
		Currency:      "USD",
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Details: map[string]interface{}{
			"method":          "credit_card",
			"last_four":       "4242",
			"card_type":       "visa",
			"expiry_month":    12,
			"expiry_year":     2026,
			"billing_address": "123 Main St, Anytown, USA",
		},
	}, nil
}

// CapturePayment implements PaymentService.CapturePayment
func (c *PaymentClient) CapturePayment(input payment.CapturePaymentInput) (payment.CapturePaymentOutput, error) {
	if c.config.EnableMock {
		return c.mockCapturePayment(input)
	}

	// Simulate capture
	captureID := fmt.Sprintf("cap_%d", time.Now().UnixNano())

	return payment.CapturePaymentOutput{
		CaptureID:      captureID,
		PaymentID:      input.PaymentID,
		Status:         "completed",
		CapturedAmount: input.Amount,
		CapturedAt:     time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// VoidPayment implements PaymentService.VoidPayment
func (c *PaymentClient) VoidPayment(input payment.VoidPaymentInput) (payment.VoidPaymentOutput, error) {
	if c.config.EnableMock {
		return c.mockVoidPayment(input)
	}

	// Simulate void
	voidID := fmt.Sprintf("void_%d", time.Now().UnixNano())

	return payment.VoidPaymentOutput{
		VoidID:    voidID,
		PaymentID: input.PaymentID,
		Status:    "completed",
		VoidedAt:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// Mock implementations for testing
func (c *PaymentClient) mockProcessPayment(input payment.ProcessPaymentInput) (payment.ProcessPaymentOutput, error) {
	status := c.config.MockResponses.DefaultStatus
	if status == "" {
		status = "completed"
	}

	// Apply success rate
	if c.config.MockResponses.SuccessRate > 0 {
		// Simple random success based on success rate
		if float64(time.Now().UnixNano()%100)/100.0 > c.config.MockResponses.SuccessRate {
			status = "failed"
		}
	}

	return payment.ProcessPaymentOutput{
		TransactionID:     input.TransactionID,
		PaymentID:         fmt.Sprintf("mock_pay_%d", time.Now().UnixNano()),
		Status:            status,
		AuthorizationCode: "MOCK_AUTH_123456",
		ProcessedAt:       time.Now().UTC().Format(time.RFC3339),
		FailureReason:     "",
		GatewayResponse: map[string]interface{}{
			"gateway":   "mock_gateway",
			"mock":      true,
			"success":   status == "completed",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	}, nil
}

func (c *PaymentClient) mockRefundPayment(input payment.RefundPaymentInput) (payment.RefundPaymentOutput, error) {
	return payment.RefundPaymentOutput{
		RefundID:       fmt.Sprintf("mock_ref_%d", time.Now().UnixNano()),
		PaymentID:      input.PaymentID,
		Status:         "completed",
		RefundedAmount: input.Amount,
		RefundedAt:     time.Now().UTC().Format(time.RFC3339),
		FailureReason:  "",
	}, nil
}

func (c *PaymentClient) mockGetPaymentStatus(input payment.GetPaymentStatusInput) (payment.GetPaymentStatusOutput, error) {
	return payment.GetPaymentStatusOutput{
		PaymentID:     input.PaymentID,
		TransactionID: input.TransactionID,
		Status:        "completed",
		Amount:        99.99,
		Currency:      "USD",
		CreatedAt:     time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339),
		UpdatedAt:     time.Now().Add(-30 * time.Minute).UTC().Format(time.RFC3339),
		Details: map[string]interface{}{
			"mock":      true,
			"test_data": "This is mock payment status data",
		},
	}, nil
}

func (c *PaymentClient) mockCapturePayment(input payment.CapturePaymentInput) (payment.CapturePaymentOutput, error) {
	return payment.CapturePaymentOutput{
		CaptureID:      fmt.Sprintf("mock_cap_%d", time.Now().UnixNano()),
		PaymentID:      input.PaymentID,
		Status:         "completed",
		CapturedAmount: input.Amount,
		CapturedAt:     time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (c *PaymentClient) mockVoidPayment(input payment.VoidPaymentInput) (payment.VoidPaymentOutput, error) {
	return payment.VoidPaymentOutput{
		VoidID:    fmt.Sprintf("mock_void_%d", time.Now().UnixNano()),
		PaymentID: input.PaymentID,
		Status:    "completed",
		VoidedAt:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}
