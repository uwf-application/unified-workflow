package client

import (
	"context"
	"time"
)

// Client is the base interface for all service clients
type Client interface {
	// Ping checks if the service is reachable
	Ping(ctx context.Context) error

	// Close closes the client connection
	Close() error

	// GetEndpoint returns the service endpoint
	GetEndpoint() string

	// IsHealthy checks if the client is healthy
	IsHealthy() bool
}

// Config is the base configuration for all clients
type Config struct {
	// Endpoint is the service endpoint URL
	Endpoint string `json:"endpoint" yaml:"endpoint"`

	// Timeout is the request timeout
	Timeout time.Duration `json:"timeout" yaml:"timeout"`

	// MaxRetries is the maximum number of retries
	MaxRetries int `json:"max_retries" yaml:"max_retries"`

	// RetryDelay is the delay between retries
	RetryDelay time.Duration `json:"retry_delay" yaml:"retry_delay"`

	// EnableTLS enables TLS for the connection
	EnableTLS bool `json:"enable_tls" yaml:"enable_tls"`

	// TLSCertPath is the path to the TLS certificate
	TLSCertPath string `json:"tls_cert_path" yaml:"tls_cert_path"`

	// TLSKeyPath is the path to the TLS key
	TLSKeyPath string `json:"tls_key_path" yaml:"tls_key_path"`

	// TLSCAPath is the path to the TLS CA certificate
	TLSCAPath string `json:"tls_ca_path" yaml:"tls_ca_path"`

	// AuthToken is the authentication token
	AuthToken string `json:"auth_token" yaml:"auth_token"`

	// EnableCircuitBreaker enables circuit breaker pattern
	EnableCircuitBreaker bool `json:"enable_circuit_breaker" yaml:"enable_circuit_breaker"`

	// CircuitBreakerThreshold is the failure threshold for circuit breaker
	CircuitBreakerThreshold int `json:"circuit_breaker_threshold" yaml:"circuit_breaker_threshold"`

	// CircuitBreakerTimeout is the timeout for circuit breaker
	CircuitBreakerTimeout time.Duration `json:"circuit_breaker_timeout" yaml:"circuit_breaker_timeout"`
}

// DefaultConfig returns the default client configuration
func DefaultConfig() Config {
	return Config{
		Endpoint:                "http://localhost:8080",
		Timeout:                 30 * time.Second,
		MaxRetries:              3,
		RetryDelay:              1 * time.Second,
		EnableTLS:               false,
		EnableCircuitBreaker:    true,
		CircuitBreakerThreshold: 5,
		CircuitBreakerTimeout:   60 * time.Second,
	}
}

// Error represents a client error
type Error struct {
	// Code is the error code
	Code string `json:"code"`

	// Message is the error message
	Message string `json:"message"`

	// Details contains additional error details
	Details map[string]interface{} `json:"details,omitempty"`

	// Retryable indicates if the error is retryable
	Retryable bool `json:"retryable"`

	// OriginalError is the original error
	OriginalError error `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// IsRetryable checks if the error is retryable
func (e *Error) IsRetryable() bool {
	return e.Retryable
}

// Common error codes
const (
	ErrCodeConnectionFailed = "CONNECTION_FAILED"
	ErrCodeTimeout          = "TIMEOUT"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeCircuitBreaker   = "CIRCUIT_BREAKER_OPEN"
)

// Request is the base request structure
type Request struct {
	// ID is the request ID for tracing
	ID string `json:"id,omitempty"`

	// TraceID is the distributed trace ID
	TraceID string `json:"trace_id,omitempty"`

	// SpanID is the span ID for tracing
	SpanID string `json:"span_id,omitempty"`

	// Timestamp is the request timestamp
	Timestamp time.Time `json:"timestamp,omitempty"`

	// Metadata contains additional request metadata
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Response is the base response structure
type Response struct {
	// ID is the response ID
	ID string `json:"id,omitempty"`

	// TraceID is the distributed trace ID
	TraceID string `json:"trace_id,omitempty"`

	// RequestID is the original request ID
	RequestID string `json:"request_id,omitempty"`

	// Timestamp is the response timestamp
	Timestamp time.Time `json:"timestamp,omitempty"`

	// ProcessingTime is the processing time in milliseconds
	ProcessingTime int64 `json:"processing_time,omitempty"`

	// Success indicates if the request was successful
	Success bool `json:"success"`

	// Error contains error details if Success is false
	Error *Error `json:"error,omitempty"`
}

// NewRequest creates a new request with tracing information
func NewRequest(ctx context.Context) Request {
	return Request{
		ID:        generateID(),
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}
}

// Helper function to generate ID
func generateID() string {
	return "req-" + time.Now().Format("20060102150405") + "-" + randomString(8)
}

// Helper function to generate random string
func randomString(length int) string {
	// Simplified implementation
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		// Using time-based seed for simplicity
		b[i] = charset[time.Now().Nanosecond()%len(charset)]
	}
	return string(b)
}
