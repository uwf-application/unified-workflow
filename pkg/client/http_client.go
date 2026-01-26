package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient is a base HTTP client implementation
type HTTPClient struct {
	config     Config
	httpClient *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(config Config) *HTTPClient {
	// Create HTTP client with TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !config.EnableTLS,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return &HTTPClient{
		config:     config,
		httpClient: httpClient,
	}
}

// DoRequest performs an HTTP request with retry logic
func (c *HTTPClient) DoRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			time.Sleep(c.config.RetryDelay * time.Duration(attempt))
		}

		// Execute request
		resp, err := c.doRequestOnce(ctx, method, path, body)
		if err != nil {
			lastErr = err
			continue
		}

		// Check response status
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		// Handle error response
		resp.Body.Close()

		// Determine if retryable
		retryable := isRetryableStatusCode(resp.StatusCode)
		if !retryable || attempt == c.config.MaxRetries {
			return nil, &Error{
				Code:          getErrorCode(resp.StatusCode),
				Message:       fmt.Sprintf("HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode)),
				Retryable:     retryable,
				OriginalError: fmt.Errorf("request failed with status %d", resp.StatusCode),
			}
		}

		lastErr = fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	return nil, &Error{
		Code:          ErrCodeConnectionFailed,
		Message:       "Max retries exceeded",
		Details:       map[string]interface{}{"attempts": c.config.MaxRetries + 1},
		Retryable:     false,
		OriginalError: lastErr,
	}
}

// doRequestOnce performs a single HTTP request
func (c *HTTPClient) doRequestOnce(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, &Error{
				Code:          ErrCodeValidation,
				Message:       "Failed to marshal request body",
				Retryable:     false,
				OriginalError: err,
			}
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	url := c.config.Endpoint + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, &Error{
			Code:          ErrCodeConnectionFailed,
			Message:       "Failed to create request",
			Retryable:     false,
			OriginalError: err,
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.config.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.AuthToken)
	}

	// Add tracing headers from context if available
	if traceID := getTraceIDFromContext(ctx); traceID != "" {
		req.Header.Set("X-Trace-ID", traceID)
	}
	if spanID := getSpanIDFromContext(ctx); spanID != "" {
		req.Header.Set("X-Span-ID", spanID)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &Error{
			Code:          ErrCodeConnectionFailed,
			Message:       "Failed to execute request",
			Retryable:     true,
			OriginalError: err,
		}
	}

	return resp, nil
}

// ParseResponse parses the HTTP response into the target struct
func (c *HTTPClient) ParseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Error{
			Code:          ErrCodeInternal,
			Message:       "Failed to read response body",
			Retryable:     false,
			OriginalError: err,
		}
	}

	if err := json.Unmarshal(body, target); err != nil {
		return &Error{
			Code:          ErrCodeInternal,
			Message:       "Failed to parse response body",
			Retryable:     false,
			OriginalError: err,
		}
	}

	return nil
}

// Ping performs a health check
func (c *HTTPClient) Ping(ctx context.Context) error {
	resp, err := c.DoRequest(ctx, "GET", "/health", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &Error{
			Code:          getErrorCode(resp.StatusCode),
			Message:       "Health check failed",
			Retryable:     true,
			OriginalError: fmt.Errorf("health check returned status %d", resp.StatusCode),
		}
	}

	return nil
}

// Close closes the HTTP client
func (c *HTTPClient) Close() error {
	// HTTP client doesn't need explicit closing in Go
	return nil
}

// GetEndpoint returns the service endpoint
func (c *HTTPClient) GetEndpoint() string {
	return c.config.Endpoint
}

// IsHealthy checks if the client is healthy
func (c *HTTPClient) IsHealthy() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.Ping(ctx) == nil
}

// Helper functions

func isRetryableStatusCode(statusCode int) bool {
	// 5xx errors are retryable, 429 (Too Many Requests) is retryable
	return statusCode == 429 || (statusCode >= 500 && statusCode < 600)
}

func getErrorCode(statusCode int) string {
	switch statusCode {
	case 400:
		return ErrCodeValidation
	case 401:
		return ErrCodeUnauthorized
	case 403:
		return ErrCodeForbidden
	case 404:
		return ErrCodeNotFound
	case 408, 429:
		return ErrCodeTimeout
	case 500, 502, 503, 504:
		return ErrCodeInternal
	default:
		if statusCode >= 400 && statusCode < 500 {
			return ErrCodeValidation
		}
		return ErrCodeInternal
	}
}

func getTraceIDFromContext(ctx context.Context) string {
	// Extract trace ID from context
	// This is a simplified implementation
	// In real implementation, you would use OpenTelemetry or similar
	return ""
}

func getSpanIDFromContext(ctx context.Context) string {
	// Extract span ID from context
	// This is a simplified implementation
	// In real implementation, you would use OpenTelemetry or similar
	return ""
}
