package baseclient

import (
	"context"
	"time"
)

// Client is the base interface for all service clients
type Client interface {
	Ping(ctx context.Context) error
	Close() error
	GetEndpoint() string
	IsHealthy() bool
}

// Config is the base configuration for all clients
type Config struct {
	Endpoint                string        `json:"endpoint" yaml:"endpoint"`
	Timeout                 time.Duration `json:"timeout" yaml:"timeout"`
	MaxRetries              int           `json:"max_retries" yaml:"max_retries"`
	RetryDelay              time.Duration `json:"retry_delay" yaml:"retry_delay"`
	EnableTLS               bool          `json:"enable_tls" yaml:"enable_tls"`
	TLSCertPath             string        `json:"tls_cert_path" yaml:"tls_cert_path"`
	TLSKeyPath              string        `json:"tls_key_path" yaml:"tls_key_path"`
	TLSCAPath               string        `json:"tls_ca_path" yaml:"tls_ca_path"`
	AuthToken               string        `json:"auth_token" yaml:"auth_token"`
	EnableCircuitBreaker    bool          `json:"enable_circuit_breaker" yaml:"enable_circuit_breaker"`
	CircuitBreakerThreshold int           `json:"circuit_breaker_threshold" yaml:"circuit_breaker_threshold"`
	CircuitBreakerTimeout   time.Duration `json:"circuit_breaker_timeout" yaml:"circuit_breaker_timeout"`
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
	Code          string                 `json:"code"`
	Message       string                 `json:"message"`
	Details       map[string]interface{} `json:"details,omitempty"`
	Retryable     bool                   `json:"retryable"`
	OriginalError error                  `json:"-"`
}

func (e *Error) Error() string  { return e.Message }
func (e *Error) IsRetryable() bool { return e.Retryable }

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
	ID        string            `json:"id,omitempty"`
	TraceID   string            `json:"trace_id,omitempty"`
	SpanID    string            `json:"span_id,omitempty"`
	Timestamp time.Time         `json:"timestamp,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Response is the base response structure
type Response struct {
	ID             string    `json:"id,omitempty"`
	TraceID        string    `json:"trace_id,omitempty"`
	RequestID      string    `json:"request_id,omitempty"`
	Timestamp      time.Time `json:"timestamp,omitempty"`
	ProcessingTime int64     `json:"processing_time,omitempty"`
	Success        bool      `json:"success"`
	Error          *Error    `json:"error,omitempty"`
}

// NewRequest creates a new request with tracing information
func NewRequest(_ context.Context) Request {
	return Request{
		ID:        generateID(),
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}
}

func generateID() string {
	return "req-" + time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().Nanosecond()%len(charset)]
	}
	return string(b)
}
