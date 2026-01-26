package di

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState int

const (
	// CircuitBreakerClosed - normal operation, requests pass through
	CircuitBreakerClosed CircuitBreakerState = iota
	// CircuitBreakerOpen - circuit is open, requests fail fast
	CircuitBreakerOpen
	// CircuitBreakerHalfOpen - testing if service has recovered
	CircuitBreakerHalfOpen
)

// String returns a string representation of the circuit breaker state
func (s CircuitBreakerState) String() string {
	switch s {
	case CircuitBreakerClosed:
		return "closed"
	case CircuitBreakerOpen:
		return "open"
	case CircuitBreakerHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig holds configuration for a circuit breaker
type CircuitBreakerConfig struct {
	// Name of the circuit breaker
	Name string

	// Failure threshold before opening circuit
	FailureThreshold int

	// Success threshold before closing circuit
	SuccessThreshold int

	// Timeout for open state before moving to half-open
	OpenTimeout time.Duration

	// Time window for counting failures
	FailureWindow time.Duration

	// Minimum number of requests before circuit can open
	MinimumRequests int

	// Enable metrics collection
	EnableMetrics bool
}

// DefaultCircuitBreakerConfig returns default circuit breaker configuration
func DefaultCircuitBreakerConfig(name string) CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:             name,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		OpenTimeout:      30 * time.Second,
		FailureWindow:    60 * time.Second,
		MinimumRequests:  10,
		EnableMetrics:    true,
	}
}

// CircuitBreakerMetrics holds metrics for a circuit breaker
type CircuitBreakerMetrics struct {
	// Total requests
	TotalRequests int64

	// Successful requests
	SuccessfulRequests int64

	// Failed requests
	FailedRequests int64

	// Rejected requests (circuit open)
	RejectedRequests int64

	// Timeouts
	TimeoutRequests int64

	// State transitions
	StateTransitions map[CircuitBreakerState]int64

	// Current state
	CurrentState CircuitBreakerState

	// Last state change time
	LastStateChange time.Time

	// Failure rate
	FailureRate float64

	// Latency percentiles (in nanoseconds)
	LatencyP50 int64
	LatencyP90 int64
	LatencyP95 int64
	LatencyP99 int64
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config CircuitBreakerConfig
	mu     sync.RWMutex

	// Current state
	state CircuitBreakerState

	// Failure tracking
	failures     []time.Time
	successes    []time.Time
	lastFailure  time.Time
	lastSuccess  time.Time
	stateChanged time.Time

	// Metrics
	metrics CircuitBreakerMetrics

	// Atomic counters for performance
	totalRequests    atomic.Int64
	successRequests  atomic.Int64
	failedRequests   atomic.Int64
	rejectedRequests atomic.Int64
	timeoutRequests  atomic.Int64
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	cb := &CircuitBreaker{
		config:       config,
		state:        CircuitBreakerClosed,
		failures:     make([]time.Time, 0, config.FailureThreshold*2),
		successes:    make([]time.Time, 0, config.SuccessThreshold*2),
		stateChanged: time.Now(),
		metrics: CircuitBreakerMetrics{
			StateTransitions: make(map[CircuitBreakerState]int64),
			CurrentState:     CircuitBreakerClosed,
			LastStateChange:  time.Now(),
		},
	}

	cb.metrics.StateTransitions[CircuitBreakerClosed] = 1

	return cb
}

// Execute runs a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// Check if circuit is open
	if !cb.allowRequest() {
		cb.rejectedRequests.Add(1)
		cb.metrics.RejectedRequests++
		return fmt.Errorf("circuit breaker '%s' is open", cb.config.Name)
	}

	cb.totalRequests.Add(1)
	cb.metrics.TotalRequests++

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	// Record result
	if err != nil {
		cb.recordFailure()
		cb.failedRequests.Add(1)
		cb.metrics.FailedRequests++
	} else {
		cb.recordSuccess()
		cb.successRequests.Add(1)
		cb.metrics.SuccessfulRequests++
	}

	// Update latency metrics
	cb.updateLatency(duration.Nanoseconds())

	return err
}

// ExecuteWithTimeout runs a function with circuit breaker protection and timeout
func (cb *CircuitBreaker) ExecuteWithTimeout(fn func() error, timeout time.Duration) error {
	// Check if circuit is open
	if !cb.allowRequest() {
		cb.rejectedRequests.Add(1)
		cb.metrics.RejectedRequests++
		return fmt.Errorf("circuit breaker '%s' is open", cb.config.Name)
	}

	cb.totalRequests.Add(1)
	cb.metrics.TotalRequests++

	// Create channel for timeout
	resultChan := make(chan error, 1)
	start := time.Now()

	// Execute function in goroutine
	go func() {
		resultChan <- fn()
	}()

	// Wait for result or timeout
	select {
	case err := <-resultChan:
		duration := time.Since(start)

		// Record result
		if err != nil {
			cb.recordFailure()
			cb.failedRequests.Add(1)
			cb.metrics.FailedRequests++
		} else {
			cb.recordSuccess()
			cb.successRequests.Add(1)
			cb.metrics.SuccessfulRequests++
		}

		// Update latency metrics
		cb.updateLatency(duration.Nanoseconds())

		return err

	case <-time.After(timeout):
		cb.timeoutRequests.Add(1)
		cb.metrics.TimeoutRequests++
		cb.recordFailure()
		return fmt.Errorf("circuit breaker '%s': timeout after %v", cb.config.Name, timeout)
	}
}

// allowRequest checks if a request should be allowed
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case CircuitBreakerClosed:
		return true

	case CircuitBreakerOpen:
		// Check if open timeout has passed
		if time.Since(cb.stateChanged) > cb.config.OpenTimeout {
			// Move to half-open state
			cb.mu.RUnlock()
			cb.transitionToHalfOpen()
			cb.mu.RLock()
			return true
		}
		return false

	case CircuitBreakerHalfOpen:
		// Allow limited number of requests to test service
		return len(cb.successes) < cb.config.SuccessThreshold

	default:
		return false
	}
}

// recordFailure records a failure
func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	cb.lastFailure = now
	cb.failures = append(cb.failures, now)

	// Clean old failures
	cb.cleanOldFailures(now)

	// Check if we should open the circuit
	if cb.shouldOpenCircuit() {
		cb.transitionToOpen()
	}
}

// recordSuccess records a success
func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	cb.lastSuccess = now
	cb.successes = append(cb.successes, now)

	// Clean old successes
	cb.cleanOldSuccesses(now)

	// Check if we should close the circuit
	if cb.shouldCloseCircuit() {
		cb.transitionToClosed()
	}
}

// cleanOldFailures removes failures outside the failure window
func (cb *CircuitBreaker) cleanOldFailures(now time.Time) {
	cutoff := now.Add(-cb.config.FailureWindow)
	i := 0
	for i < len(cb.failures) && cb.failures[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		cb.failures = cb.failures[i:]
	}
}

// cleanOldSuccesses removes successes outside the failure window
func (cb *CircuitBreaker) cleanOldSuccesses(now time.Time) {
	cutoff := now.Add(-cb.config.FailureWindow)
	i := 0
	for i < len(cb.successes) && cb.successes[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		cb.successes = cb.successes[i:]
	}
}

// shouldOpenCircuit determines if circuit should open
func (cb *CircuitBreaker) shouldOpenCircuit() bool {
	if cb.state == CircuitBreakerOpen {
		return false
	}

	// Need minimum requests before opening
	totalRequests := len(cb.failures) + len(cb.successes)
	if totalRequests < cb.config.MinimumRequests {
		return false
	}

	// Check failure threshold
	if len(cb.failures) >= cb.config.FailureThreshold {
		// Calculate failure rate
		failureRate := float64(len(cb.failures)) / float64(totalRequests)
		cb.metrics.FailureRate = failureRate
		return true
	}

	return false
}

// shouldCloseCircuit determines if circuit should close
func (cb *CircuitBreaker) shouldCloseCircuit() bool {
	if cb.state != CircuitBreakerHalfOpen {
		return false
	}

	// Check success threshold
	return len(cb.successes) >= cb.config.SuccessThreshold
}

// transitionToOpen transitions to open state
func (cb *CircuitBreaker) transitionToOpen() {
	if cb.state == CircuitBreakerOpen {
		return
	}

	cb.state = CircuitBreakerOpen
	cb.stateChanged = time.Now()
	cb.metrics.CurrentState = CircuitBreakerOpen
	cb.metrics.LastStateChange = cb.stateChanged
	cb.metrics.StateTransitions[CircuitBreakerOpen]++

	// Clear successes when opening
	cb.successes = cb.successes[:0]
}

// transitionToHalfOpen transitions to half-open state
func (cb *CircuitBreaker) transitionToHalfOpen() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == CircuitBreakerHalfOpen {
		return
	}

	cb.state = CircuitBreakerHalfOpen
	cb.stateChanged = time.Now()
	cb.metrics.CurrentState = CircuitBreakerHalfOpen
	cb.metrics.LastStateChange = cb.stateChanged
	cb.metrics.StateTransitions[CircuitBreakerHalfOpen]++

	// Clear successes and failures when moving to half-open
	cb.successes = cb.successes[:0]
	cb.failures = cb.failures[:0]
}

// transitionToClosed transitions to closed state
func (cb *CircuitBreaker) transitionToClosed() {
	if cb.state == CircuitBreakerClosed {
		return
	}

	cb.state = CircuitBreakerClosed
	cb.stateChanged = time.Now()
	cb.metrics.CurrentState = CircuitBreakerClosed
	cb.metrics.LastStateChange = cb.stateChanged
	cb.metrics.StateTransitions[CircuitBreakerClosed]++

	// Clear failures and successes when closing
	cb.failures = cb.failures[:0]
	cb.successes = cb.successes[:0]
}

// updateLatency updates latency metrics
func (cb *CircuitBreaker) updateLatency(latencyNs int64) {
	if !cb.config.EnableMetrics {
		return
	}

	// Simple percentile estimation (for demo purposes)
	// In production, use a proper histogram library
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Update percentiles (simplified)
	if latencyNs > cb.metrics.LatencyP99 {
		cb.metrics.LatencyP99 = latencyNs
	}
	if latencyNs > cb.metrics.LatencyP95 {
		cb.metrics.LatencyP95 = latencyNs
	}
	if latencyNs > cb.metrics.LatencyP90 {
		cb.metrics.LatencyP90 = latencyNs
	}
	if latencyNs > cb.metrics.LatencyP50 {
		cb.metrics.LatencyP50 = latencyNs
	}
}

// GetState returns the current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetMetrics returns circuit breaker metrics
func (cb *CircuitBreaker) GetMetrics() CircuitBreakerMetrics {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	// Update atomic counters
	cb.metrics.TotalRequests = cb.totalRequests.Load()
	cb.metrics.SuccessfulRequests = cb.successRequests.Load()
	cb.metrics.FailedRequests = cb.failedRequests.Load()
	cb.metrics.RejectedRequests = cb.rejectedRequests.Load()
	cb.metrics.TimeoutRequests = cb.timeoutRequests.Load()

	// Calculate failure rate
	total := cb.metrics.TotalRequests
	if total > 0 {
		cb.metrics.FailureRate = float64(cb.metrics.FailedRequests) / float64(total)
	}

	return cb.metrics
}

// Reset resets the circuit breaker
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = CircuitBreakerClosed
	cb.stateChanged = time.Now()
	cb.failures = cb.failures[:0]
	cb.successes = cb.successes[:0]
	cb.lastFailure = time.Time{}
	cb.lastSuccess = time.Time{}

	// Reset metrics
	cb.totalRequests.Store(0)
	cb.successRequests.Store(0)
	cb.failedRequests.Store(0)
	cb.rejectedRequests.Store(0)
	cb.timeoutRequests.Store(0)

	cb.metrics = CircuitBreakerMetrics{
		StateTransitions: make(map[CircuitBreakerState]int64),
		CurrentState:     CircuitBreakerClosed,
		LastStateChange:  time.Now(),
	}
	cb.metrics.StateTransitions[CircuitBreakerClosed] = 1
}

// IsClosed returns true if circuit is closed
func (cb *CircuitBreaker) IsClosed() bool {
	return cb.GetState() == CircuitBreakerClosed
}

// IsOpen returns true if circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.GetState() == CircuitBreakerOpen
}

// IsHalfOpen returns true if circuit is half-open
func (cb *CircuitBreaker) IsHalfOpen() bool {
	return cb.GetState() == CircuitBreakerHalfOpen
}

// String returns a string representation of the circuit breaker
func (cb *CircuitBreaker) String() string {
	state := cb.GetState()
	metrics := cb.GetMetrics()

	return fmt.Sprintf(
		"CircuitBreaker '%s': state=%s, total=%d, success=%d, failed=%d, rejected=%d, failure_rate=%.2f%%",
		cb.config.Name,
		state.String(),
		metrics.TotalRequests,
		metrics.SuccessfulRequests,
		metrics.FailedRequests,
		metrics.RejectedRequests,
		metrics.FailureRate*100,
	)
}

// CircuitBreakerManager manages multiple circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate gets or creates a circuit breaker
func (m *CircuitBreakerManager) GetOrCreate(name string) *CircuitBreaker {
	return m.GetOrCreateWithConfig(DefaultCircuitBreakerConfig(name))
}

// GetOrCreateWithConfig gets or creates a circuit breaker with custom config
func (m *CircuitBreakerManager) GetOrCreateWithConfig(config CircuitBreakerConfig) *CircuitBreaker {
	m.mu.RLock()
	breaker, exists := m.breakers[config.Name]
	m.mu.RUnlock()

	if exists {
		return breaker
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	breaker, exists = m.breakers[config.Name]
	if exists {
		return breaker
	}

	breaker = NewCircuitBreaker(config)
	m.breakers[config.Name] = breaker
	return breaker
}

// Get returns a circuit breaker by name
func (m *CircuitBreakerManager) Get(name string) (*CircuitBreaker, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	breaker, exists := m.breakers[name]
	return breaker, exists
}

// GetAll returns all circuit breakers
func (m *CircuitBreakerManager) GetAll() map[string]*CircuitBreaker {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*CircuitBreaker)
	for name, breaker := range m.breakers {
		result[name] = breaker
	}
	return result
}

// Remove removes a circuit breaker
func (m *CircuitBreakerManager) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.breakers, name)
}

// ResetAll resets all circuit breakers
func (m *CircuitBreakerManager) ResetAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, breaker := range m.breakers {
		breaker.Reset()
	}
}

// GetMetrics returns metrics for all circuit breakers
func (m *CircuitBreakerManager) GetMetrics() map[string]CircuitBreakerMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]CircuitBreakerMetrics)
	for name, breaker := range m.breakers {
		metrics[name] = breaker.GetMetrics()
	}
	return metrics
}

// CircuitBreakerMiddleware provides circuit breaker middleware for HTTP handlers
type CircuitBreakerMiddleware struct {
	manager *CircuitBreakerManager
}

// NewCircuitBreakerMiddleware creates a new circuit breaker middleware
func NewCircuitBreakerMiddleware() *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{
		manager: NewCircuitBreakerManager(),
	}
}

// Wrap wraps an HTTP handler with circuit breaker protection
func (m *CircuitBreakerMiddleware) Wrap(handler http.Handler, breakerName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		breaker := m.manager.GetOrCreate(breakerName)

		err := breaker.Execute(func() error {
			// Create a response writer that captures the status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			handler.ServeHTTP(rw, r)

			// Consider 5xx errors as failures
			if rw.statusCode >= 500 {
				return fmt.Errorf("HTTP %d: %s", rw.statusCode, http.StatusText(rw.statusCode))
			}
			return nil
		})

		if err != nil {
			// Circuit breaker rejected the request or handler failed
			if breaker.IsOpen() {
				http.Error(w, "Service unavailable (circuit breaker open)", http.StatusServiceUnavailable)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})
}

// responseWriter captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// GetManager returns the circuit breaker manager
func (m *CircuitBreakerMiddleware) GetManager() *CircuitBreakerManager {
	return m.manager
}
