package serviceclients

import (
	"fmt"

	antifraudservice "unified-workflow/internal/primitive/services/antifraud"
	"unified-workflow/internal/primitive/services/antifraud/models"
	"unified-workflow/internal/serviceclients/antifraud"
)

// ClientProvider provides access to various service clients
type ClientProvider interface {
	// GetAntifraudClient returns an antifraud detection client
	GetAntifraudClient() (antifraudservice.AntifraudService, error)

	// Close closes all client connections
	Close() error
}

// Config holds configuration for all service clients
type Config struct {
	Antifraud models.ClientConfig `json:"antifraud"`
}

// DefaultClientProvider implements ClientProvider with default configurations
type DefaultClientProvider struct {
	config          *Config
	antifraudClient antifraudservice.AntifraudService
	initialized     bool
}

// NewClientProvider creates a new client provider with default configuration
func NewClientProvider() (*DefaultClientProvider, error) {
	config := &Config{
		Antifraud: models.ClientConfig{
			APIKey:                  "",
			Host:                    "https://api.antifraudservice.com/v1",
			Timeout:                 30,
			Enabled:                 false,
			MaxRetries:              3,
			CircuitBreakerEnabled:   true,
			CircuitBreakerThreshold: 5,
			CircuitBreakerTimeout:   30,
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

// GetAntifraudClient implements ClientProvider.GetAntifraudClient
func (p *DefaultClientProvider) GetAntifraudClient() (antifraudservice.AntifraudService, error) {
	if err := p.ensureInitialized(); err != nil {
		return nil, err
	}
	return p.antifraudClient, nil
}

// Close implements ClientProvider.Close
func (p *DefaultClientProvider) Close() error {
	// In a real implementation, this would close HTTP clients, database connections, etc.
	p.initialized = false
	p.antifraudClient = nil
	return nil
}

// ensureInitialized lazily initializes all service clients
func (p *DefaultClientProvider) ensureInitialized() error {
	if p.initialized {
		return nil
	}

	// Initialize antifraud client
	antifraudClient, err := antifraud.NewClient(p.config.Antifraud)
	if err != nil {
		return fmt.Errorf("failed to initialize antifraud client: %w", err)
	}
	p.antifraudClient = antifraudClient

	p.initialized = true
	return nil
}

// MockClientProvider implements ClientProvider with mock implementations for testing
type MockClientProvider struct {
	AntifraudClient antifraudservice.AntifraudService
}

// GetAntifraudClient implements ClientProvider.GetAntifraudClient
func (m *MockClientProvider) GetAntifraudClient() (antifraudservice.AntifraudService, error) {
	if m.AntifraudClient == nil {
		// Create a default mock antifraud client
		config := &models.ClientConfig{
			Enabled:               true,
			CircuitBreakerEnabled: false,
		}
		client, err := antifraud.NewClient(*config)
		if err != nil {
			return nil, err
		}
		m.AntifraudClient = client
	}
	return m.AntifraudClient, nil
}

// Close implements ClientProvider.Close
func (m *MockClientProvider) Close() error {
	// Nothing to close for mocks
	return nil
}
