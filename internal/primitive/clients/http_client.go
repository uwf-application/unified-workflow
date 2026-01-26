package clients

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPClientImpl implements the HTTPClient interface
type HTTPClientImpl struct {
	*BaseClientImpl
	httpClient *http.Client
	baseURL    string
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(config HTTPClientConfig) *HTTPClientImpl {
	baseConfig := config.ClientConfig
	baseConfig.Type = ClientTypeHTTP

	// Create HTTP client with configured timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	client := &HTTPClientImpl{
		BaseClientImpl: NewBaseClient(baseConfig),
		httpClient:     httpClient,
		baseURL:        config.BaseURL,
	}

	return client
}

// Connect implements BaseClient.Connect for HTTP
func (c *HTTPClientImpl) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return nil
	}

	// For HTTP client, connection is established on demand
	// We just mark it as connected
	c.isConnected = true
	c.metrics.ConnectionCount++
	c.lastUsed = time.Now()

	c.notifyCallbacks(EventConnected, c)
	return nil
}

// Authenticate implements BaseClient.Authenticate for HTTP
func (c *HTTPClientImpl) Authenticate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// HTTP authentication depends on the auth type
	switch c.config.AuthType {
	case AuthTypeAPIKey:
		// API key will be added to headers on each request
		c.authInfo.Type = AuthTypeAPIKey
		c.authInfo.Token = c.config.APIKey
	case AuthTypeJWT:
		// JWT token will be added to headers on each request
		c.authInfo.Type = AuthTypeJWT
		c.authInfo.Token = c.config.JWTToken
	case AuthTypeOAuth2:
		// OAuth2 token will be added to headers on each request
		c.authInfo.Type = AuthTypeOAuth2
		c.authInfo.Token = c.config.OAuth2Token
	case AuthTypeBasic:
		// Basic auth will be added to headers on each request
		c.authInfo.Type = AuthTypeBasic
	default:
		c.authInfo.Type = AuthTypeNone
	}

	c.authInfo.LastRefresh = time.Now()
	c.lastUsed = time.Now()

	c.notifyCallbacks(EventAuthenticated, c)
	return nil
}

// Get performs an HTTP GET request
func (c *HTTPClientImpl) Get(ctx context.Context, path string, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil, headers)
}

// Post performs an HTTP POST request
func (c *HTTPClientImpl) Post(ctx context.Context, path string, body []byte, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPost, path, body, headers)
}

// Put performs an HTTP PUT request
func (c *HTTPClientImpl) Put(ctx context.Context, path string, body []byte, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPut, path, body, headers)
}

// Delete performs an HTTP DELETE request
func (c *HTTPClientImpl) Delete(ctx context.Context, path string, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil, headers)
}

// Patch performs an HTTP PATCH request
func (c *HTTPClientImpl) Patch(ctx context.Context, path string, body []byte, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodPatch, path, body, headers)
}

// Head performs an HTTP HEAD request
func (c *HTTPClientImpl) Head(ctx context.Context, path string, headers map[string]string) (*HTTPResponse, error) {
	return c.doRequest(ctx, http.MethodHead, path, nil, headers)
}

// NewRequest creates a new HTTP request
func (c *HTTPClientImpl) NewRequest(method, path string) *HTTPRequest {
	url := c.buildURL(path)
	return &HTTPRequest{
		Method:  method,
		URL:     url,
		Headers: make(map[string]string),
		Body:    nil,
		Timeout: c.config.Timeout,
	}
}

// DoRequest executes an HTTP request
func (c *HTTPClientImpl) DoRequest(req *HTTPRequest) (*HTTPResponse, error) {
	return c.doRequest(context.Background(), req.Method, req.URL, req.Body, req.Headers)
}

// BatchRequest executes multiple HTTP requests
func (c *HTTPClientImpl) BatchRequest(ctx context.Context, requests []*HTTPRequest) ([]*HTTPResponse, error) {
	// Simple sequential implementation
	// Could be enhanced with goroutines for parallel execution
	responses := make([]*HTTPResponse, 0, len(requests))

	for _, req := range requests {
		resp, err := c.DoRequest(req)
		if err != nil {
			// Create error response
			resp = &HTTPResponse{
				StatusCode: 0,
				Headers:    nil,
				Body:       []byte(err.Error()),
				Latency:    0,
				Request:    req,
			}
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// doRequest performs the actual HTTP request
func (c *HTTPClientImpl) doRequest(ctx context.Context, method, path string, body []byte, headers map[string]string) (*HTTPResponse, error) {
	start := time.Now()

	if !c.isConnected {
		c.recordError(ErrClientNotConnected)
		return nil, ErrClientNotConnected
	}

	url := c.buildURL(path)

	// Create request
	var bodyReader io.Reader
	if body != nil {
		bodyReader = strings.NewReader(string(body))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		c.recordError(err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add authentication headers
	c.addAuthHeaders(req)

	// Add default headers from config
	for key, value := range c.config.Headers {
		if _, exists := req.Header[key]; !exists {
			req.Header.Set(key, value)
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	latency := time.Since(start)

	if err != nil {
		c.recordRequest(latency, false)
		c.recordError(err)
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.recordRequest(latency, false)
		c.recordError(err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert headers
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	// Create HTTP response
	httpResp := &HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       respBody,
		Latency:    latency,
		Request: &HTTPRequest{
			Method:  method,
			URL:     url,
			Headers: headers,
			Body:    body,
			Timeout: c.config.Timeout,
		},
	}

	// Record successful request
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	c.recordRequest(latency, success)

	if !success {
		c.recordError(fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody)))
	}

	return httpResp, nil
}

// buildURL builds the full URL from base URL and path
func (c *HTTPClientImpl) buildURL(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	if c.baseURL == "" {
		return path
	}

	// Ensure baseURL ends with slash if path doesn't start with slash
	if !strings.HasSuffix(c.baseURL, "/") && !strings.HasPrefix(path, "/") {
		return c.baseURL + "/" + path
	}

	if strings.HasSuffix(c.baseURL, "/") && strings.HasPrefix(path, "/") {
		return c.baseURL + path[1:]
	}

	return c.baseURL + path
}

// addAuthHeaders adds authentication headers to the request
func (c *HTTPClientImpl) addAuthHeaders(req *http.Request) {
	switch c.config.AuthType {
	case AuthTypeAPIKey:
		if c.config.APIKey != "" {
			// Try to determine where to put the API key
			if c.config.Headers["X-API-Key"] != "" {
				req.Header.Set("X-API-Key", c.config.APIKey)
			} else if c.config.Headers["Authorization"] != "" {
				req.Header.Set("Authorization", c.config.APIKey)
			} else {
				// Default to Authorization header
				req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
			}
		}
	case AuthTypeJWT:
		if c.config.JWTToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.config.JWTToken)
		}
	case AuthTypeOAuth2:
		if c.config.OAuth2Token != "" {
			req.Header.Set("Authorization", "Bearer "+c.config.OAuth2Token)
		}
	case AuthTypeBasic:
		// Basic auth would need username/password
		// This is a simplified implementation
		if c.config.Headers["Authorization"] != "" {
			req.Header.Set("Authorization", c.config.Headers["Authorization"])
		}
	}
}

// HealthCheck implements BaseClient.HealthCheck for HTTP
func (c *HTTPClientImpl) HealthCheck(ctx context.Context) (HealthStatus, error) {
	start := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	// For HTTP client, health check typically involves making a request to a health endpoint
	// For now, we'll just check if we can make a request to the base URL
	status := HealthStatus{
		Healthy:   c.isConnected,
		Message:   "",
		CheckedAt: time.Now(),
		Latency:   0,
	}

	if !c.isConnected {
		status.Message = "HTTP client is not connected"
	} else if c.baseURL != "" {
		// Try to make a HEAD request to the base URL
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.baseURL, nil)
		if err == nil {
			resp, err := c.httpClient.Do(req)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode < 500 {
					status.Message = "HTTP client is healthy"
				} else {
					status.Message = fmt.Sprintf("HTTP client returned status %d", resp.StatusCode)
					status.Healthy = false
				}
			} else {
				status.Message = fmt.Sprintf("HTTP health check failed: %v", err)
				status.Healthy = false
			}
		} else {
			status.Message = fmt.Sprintf("Failed to create health check request: %v", err)
			status.Healthy = false
		}
	} else {
		status.Message = "HTTP client has no base URL configured"
	}

	status.Latency = time.Since(start)
	c.lastUsed = time.Now()

	c.notifyCallbacks(EventHealthCheck, c)
	return status, nil
}

// GetBaseURL returns the configured base URL
func (c *HTTPClientImpl) GetBaseURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.baseURL
}

// GetHTTPClient returns the underlying http.Client
func (c *HTTPClientImpl) GetHTTPClient() *http.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.httpClient
}
