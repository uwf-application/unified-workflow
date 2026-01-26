package serviceclients

import (
	"fmt"

	fraudservice "unified-workflow/internal/primitive/services/fraud"
	paymentservice "unified-workflow/internal/primitive/services/payment"
	"unified-workflow/internal/serviceclients/fraud"
	"unified-workflow/internal/serviceclients/payment"
)

// ClientProvider provides access to various service clients
type ClientProvider interface {
	// GetFraudClient returns a fraud detection client
	GetFraudClient() (fraudservice.FraudService, error)

	// GetPaymentClient returns a payment processing client
	GetPaymentClient() (paymentservice.PaymentService, error)

	// Close closes all client connections
	Close() error
}

// Config holds configuration for all service clients
type Config struct {
	Fraud   fraud.Config   `json:"fraud"`
	Payment payment.Config `json:"payment"`
}

// DefaultClientProvider implements ClientProvider with default configurations
type DefaultClientProvider struct {
	config        *Config
	fraudClient   fraudservice.FraudService
	paymentClient paymentservice.PaymentService
	initialized   bool
}

// NewClientProvider creates a new client provider with default configuration
func NewClientProvider() (*DefaultClientProvider, error) {
	config := &Config{
		Fraud: fraud.Config{
			APIKey:        "demo_fraud_api_key",
			APIURL:        "https://api.fraudservice.com/v1",
			Timeout:       30,
			MaxRetries:    3,
			RiskThreshold: 75,
			EnableMock:    true,
			MockResponses: fraud.MockConfig{
				DefaultRiskScore: 65,
				DefaultAction:    "review",
			},
		},
		Payment: payment.Config{
			APIKey:     "demo_payment_api_key",
			APIURL:     "https://api.paymentservice.com/v1",
			Timeout:    30,
			MaxRetries: 3,
			EnableMock: true,
			MockResponses: payment.MockConfig{
				DefaultStatus: "completed",
				SuccessRate:   0.95,
			},
		},
	}

	return NewClientProviderWithConfig(config)
}

// NewClientProviderWithConfig creates a new client provider with custom configuration
func NewClientProviderWithConfig(config *Config) (*DefaultClientProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &DefaultClientProvider{
		config:      config,
		initialized: false,
	}, nil
}

// GetFraudClient implements ClientProvider.GetFraudClient
func (p *DefaultClientProvider) GetFraudClient() (fraudservice.FraudService, error) {
	if err := p.ensureInitialized(); err != nil {
		return nil, err
	}
	return p.fraudClient, nil
}

// GetPaymentClient implements ClientProvider.GetPaymentClient
func (p *DefaultClientProvider) GetPaymentClient() (paymentservice.PaymentService, error) {
	if err := p.ensureInitialized(); err != nil {
		return nil, err
	}
	return p.paymentClient, nil
}

// Close implements ClientProvider.Close
func (p *DefaultClientProvider) Close() error {
	// In a real implementation, this would close HTTP clients, database connections, etc.
	p.initialized = false
	p.fraudClient = nil
	p.paymentClient = nil
	return nil
}

// ensureInitialized lazily initializes all service clients
func (p *DefaultClientProvider) ensureInitialized() error {
	if p.initialized {
		return nil
	}

	// Initialize fraud client
	fraudClient, err := fraud.NewFraudClientWithConfig(&p.config.Fraud)
	if err != nil {
		return fmt.Errorf("failed to initialize fraud client: %w", err)
	}
	p.fraudClient = fraudClient

	// Initialize payment client
	paymentClient, err := payment.NewPaymentClientWithConfig(&p.config.Payment)
	if err != nil {
		return fmt.Errorf("failed to initialize payment client: %w", err)
	}
	p.paymentClient = paymentClient

	p.initialized = true
	return nil
}

// MockClientProvider implements ClientProvider with mock implementations for testing
type MockClientProvider struct {
	FraudClient   fraudservice.FraudService
	PaymentClient paymentservice.PaymentService
}

// GetFraudClient implements ClientProvider.GetFraudClient
func (m *MockClientProvider) GetFraudClient() (fraudservice.FraudService, error) {
	if m.FraudClient == nil {
		// Create a default mock fraud client
		config := &fraud.Config{
			EnableMock: true,
			MockResponses: fraud.MockConfig{
				DefaultRiskScore: 30,
				DefaultAction:    "allow",
			},
		}
		client, err := fraud.NewFraudClientWithConfig(config)
		if err != nil {
			return nil, err
		}
		m.FraudClient = client
	}
	return m.FraudClient, nil
}

// GetPaymentClient implements ClientProvider.GetPaymentClient
func (m *MockClientProvider) GetPaymentClient() (paymentservice.PaymentService, error) {
	if m.PaymentClient == nil {
		// Create a default mock payment client
		config := &payment.Config{
			EnableMock: true,
			MockResponses: payment.MockConfig{
				DefaultStatus: "completed",
				SuccessRate:   1.0,
			},
		}
		client, err := payment.NewPaymentClientWithConfig(config)
		if err != nil {
			return nil, err
		}
		m.PaymentClient = client
	}
	return m.PaymentClient, nil
}

// Close implements ClientProvider.Close
func (m *MockClientProvider) Close() error {
	// Nothing to close for mocks
	return nil
}
