package antifraud

import (
	"fmt"
	"time"

	primitiveantifraud "unified-workflow/internal/primitive/services/antifraud"
	"unified-workflow/internal/primitive/services/antifraud/models"
)

// antifraudProxy wraps an AntifraudService to add logging, metrics, and circuit breaking
type antifraudProxy struct {
	service       primitiveantifraud.AntifraudService
	config        models.ClientConfig
	failureCount  int
	lastFailure   time.Time
	circuitOpen   bool
	circuitOpenAt time.Time
}

// NewProxy creates a new antifraud proxy
func NewProxy(service primitiveantifraud.AntifraudService, config models.ClientConfig) primitiveantifraud.AntifraudService {
	return &antifraudProxy{
		service: service,
		config:  config,
	}
}

// checkCircuitBreaker checks if the circuit breaker should be open
func (p *antifraudProxy) checkCircuitBreaker() bool {
	if !p.config.CircuitBreakerEnabled {
		return false
	}

	if p.circuitOpen {
		// Check if we should try to close the circuit
		timeSinceOpen := time.Since(p.circuitOpenAt)
		if timeSinceOpen > time.Duration(p.config.CircuitBreakerTimeout)*time.Second {
			// Try to close the circuit
			p.circuitOpen = false
			p.failureCount = 0
			return false
		}
		return true
	}

	// Check if we should open the circuit
	if p.failureCount >= p.config.CircuitBreakerThreshold {
		p.circuitOpen = true
		p.circuitOpenAt = time.Now()
		return true
	}

	return false
}

// recordFailure records a failure for circuit breaker
func (p *antifraudProxy) recordFailure() {
	p.failureCount++
	p.lastFailure = time.Now()
}

// recordSuccess resets the failure count
func (p *antifraudProxy) recordSuccess() {
	p.failureCount = 0
}

// StoreTransaction stores a transaction in the antifraud system
func (p *antifraudProxy) StoreTransaction(afTransaction models.AF_Transaction) error {
	// Check circuit breaker
	if p.checkCircuitBreaker() {
		return fmt.Errorf("circuit breaker is open for antifraud service")
	}

	// Log start
	startTime := time.Now()

	// Call the underlying service
	err := p.service.StoreTransaction(afTransaction)

	// Record metrics
	duration := time.Since(startTime)

	if err != nil {
		p.recordFailure()
		// Log error
		return fmt.Errorf("antifraud.StoreTransaction failed after %v: %w", duration, err)
	}

	p.recordSuccess()
	// Log success
	return nil
}

// ValidateTransactionByAML validates a transaction using the AML service
func (p *antifraudProxy) ValidateTransactionByAML(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	// Check circuit breaker
	if p.checkCircuitBreaker() {
		return models.ServiceResolution{}, fmt.Errorf("circuit breaker is open for antifraud service")
	}

	// Log start
	startTime := time.Now()

	// Call the underlying service
	result, err := p.service.ValidateTransactionByAML(afTransaction)

	// Record metrics
	duration := time.Since(startTime)

	if err != nil {
		p.recordFailure()
		// Log error
		return models.ServiceResolution{}, fmt.Errorf("antifraud.ValidateTransactionByAML failed after %v: %w", duration, err)
	}

	p.recordSuccess()
	// Log success
	return result, nil
}

// ValidateTransactionByFC validates a transaction using the FC service
func (p *antifraudProxy) ValidateTransactionByFC(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	// Check circuit breaker
	if p.checkCircuitBreaker() {
		return models.ServiceResolution{}, fmt.Errorf("circuit breaker is open for antifraud service")
	}

	// Log start
	startTime := time.Now()

	// Call the underlying service
	result, err := p.service.ValidateTransactionByFC(afTransaction)

	// Record metrics
	duration := time.Since(startTime)

	if err != nil {
		p.recordFailure()
		// Log error
		return models.ServiceResolution{}, fmt.Errorf("antifraud.ValidateTransactionByFC failed after %v: %w", duration, err)
	}

	p.recordSuccess()
	// Log success
	return result, nil
}

// ValidateTransactionByML validates a transaction using the ML service
func (p *antifraudProxy) ValidateTransactionByML(afTransaction models.AF_Transaction) (models.ServiceResolution, error) {
	// Check circuit breaker
	if p.checkCircuitBreaker() {
		return models.ServiceResolution{}, fmt.Errorf("circuit breaker is open for antifraud service")
	}

	// Log start
	startTime := time.Now()

	// Call the underlying service
	result, err := p.service.ValidateTransactionByML(afTransaction)

	// Record metrics
	duration := time.Since(startTime)

	if err != nil {
		p.recordFailure()
		// Log error
		return models.ServiceResolution{}, fmt.Errorf("antifraud.ValidateTransactionByML failed after %v: %w", duration, err)
	}

	p.recordSuccess()
	// Log success
	return result, nil
}

// StoreServiceResolution stores the resolution from a service check
func (p *antifraudProxy) StoreServiceResolution(resolution models.ServiceResolution) error {
	// Check circuit breaker
	if p.checkCircuitBreaker() {
		return fmt.Errorf("circuit breaker is open for antifraud service")
	}

	// Log start
	startTime := time.Now()

	// Call the underlying service
	err := p.service.StoreServiceResolution(resolution)

	// Record metrics
	duration := time.Since(startTime)

	if err != nil {
		p.recordFailure()
		// Log error
		return fmt.Errorf("antifraud.StoreServiceResolution failed after %v: %w", duration, err)
	}

	p.recordSuccess()
	// Log success
	return nil
}

// AddTransactionServiceCheck adds a completed service check resolution
func (p *antifraudProxy) AddTransactionServiceCheck(resolution models.ServiceResolution) error {
	// Check circuit breaker
	if p.checkCircuitBreaker() {
		return fmt.Errorf("circuit breaker is open for antifraud service")
	}

	// Log start
	startTime := time.Now()

	// Call the underlying service
	err := p.service.AddTransactionServiceCheck(resolution)

	// Record metrics
	duration := time.Since(startTime)

	if err != nil {
		p.recordFailure()
		// Log error
		return fmt.Errorf("antifraud.AddTransactionServiceCheck failed after %v: %w", duration, err)
	}

	p.recordSuccess()
	// Log success
	return nil
}

// FinalizeTransaction finalizes the transaction validation process
func (p *antifraudProxy) FinalizeTransaction(afTransaction models.AF_Transaction) (models.FinalResolution, error) {
	// Check circuit breaker
	if p.checkCircuitBreaker() {
		return models.FinalResolution{}, fmt.Errorf("circuit breaker is open for antifraud service")
	}

	// Log start
	startTime := time.Now()

	// Call the underlying service
	result, err := p.service.FinalizeTransaction(afTransaction)

	// Record metrics
	duration := time.Since(startTime)

	if err != nil {
		p.recordFailure()
		// Log error
		return models.FinalResolution{}, fmt.Errorf("antifraud.FinalizeTransaction failed after %v: %w", duration, err)
	}

	p.recordSuccess()
	// Log success
	return result, nil
}

// StoreFinalResolution stores the final resolution of the transaction
func (p *antifraudProxy) StoreFinalResolution(resolution models.FinalResolution) error {
	// Check circuit breaker
	if p.checkCircuitBreaker() {
		return fmt.Errorf("circuit breaker is open for antifraud service")
	}

	// Log start
	startTime := time.Now()

	// Call the underlying service
	err := p.service.StoreFinalResolution(resolution)

	// Record metrics
	duration := time.Since(startTime)

	if err != nil {
		p.recordFailure()
		// Log error
		return fmt.Errorf("antifraud.StoreFinalResolution failed after %v: %w", duration, err)
	}

	p.recordSuccess()
	// Log success
	return nil
}

// HealthCheck checks the health of the antifraud service
func (p *antifraudProxy) HealthCheck() (bool, error) {
	// Circuit breaker doesn't apply to health checks
	return p.service.HealthCheck()
}

// GetConfig returns the current configuration
func (p *antifraudProxy) GetConfig() models.ClientConfig {
	return p.config
}

// GetCircuitBreakerStatus returns the current circuit breaker status
func (p *antifraudProxy) GetCircuitBreakerStatus() (bool, int, time.Time) {
	return p.circuitOpen, p.failureCount, p.circuitOpenAt
}
