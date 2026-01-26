package clients

import (
	"context"
	"sync"
	"time"
)

// BaseClientImpl provides a base implementation for clients
type BaseClientImpl struct {
	config      ClientConfig
	authInfo    AuthInfo
	metrics     ClientMetrics
	lastError   error
	isConnected bool
	mu          sync.RWMutex
	createdAt   time.Time
	lastUsed    time.Time
	callbacks   []ClientCallback
}

// ClientCallback is a function that gets called on client events
type ClientCallback func(event ClientEvent, client BaseClient)

// ClientEvent represents an event that occurs in a client's lifecycle
type ClientEvent string

const (
	EventConnected     ClientEvent = "connected"
	EventDisconnected  ClientEvent = "disconnected"
	EventAuthenticated ClientEvent = "authenticated"
	EventError         ClientEvent = "error"
	EventShutdown      ClientEvent = "shutdown"
	EventHealthCheck   ClientEvent = "health_check"
)

// NewBaseClient creates a new base client
func NewBaseClient(config ClientConfig) *BaseClientImpl {
	return &BaseClientImpl{
		config:    config,
		createdAt: time.Now(),
		metrics: ClientMetrics{
			RequestsTotal:     0,
			RequestsFailed:    0,
			AvgLatency:        0,
			LastRequestTime:   time.Time{},
			ConnectionCount:   0,
			ActiveConnections: 0,
		},
		authInfo: AuthInfo{
			Type:        config.AuthType,
			ExpiresAt:   nil,
			Token:       "",
			Scopes:      nil,
			LastRefresh: time.Time{},
		},
	}
}

// Connect implements BaseClient.Connect
func (c *BaseClientImpl) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return nil
	}

	// Default implementation - should be overridden by specific clients
	c.isConnected = true
	c.metrics.ConnectionCount++
	c.lastUsed = time.Now()

	c.notifyCallbacks(EventConnected, c)
	return nil
}

// Disconnect implements BaseClient.Disconnect
func (c *BaseClientImpl) Disconnect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return nil
	}

	c.isConnected = false
	c.metrics.ActiveConnections--
	c.lastUsed = time.Now()

	c.notifyCallbacks(EventDisconnected, c)
	return nil
}

// IsConnected implements BaseClient.IsConnected
func (c *BaseClientImpl) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// Ping implements BaseClient.Ping
func (c *BaseClientImpl) Ping(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return ErrClientNotConnected
	}

	// Default implementation - should be overridden by specific clients
	c.lastUsed = time.Now()
	return nil
}

// Authenticate implements BaseClient.Authenticate
func (c *BaseClientImpl) Authenticate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Default implementation - should be overridden by specific clients
	c.authInfo.LastRefresh = time.Now()
	c.lastUsed = time.Now()

	c.notifyCallbacks(EventAuthenticated, c)
	return nil
}

// GetAuthInfo implements BaseClient.GetAuthInfo
func (c *BaseClientImpl) GetAuthInfo() AuthInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.authInfo
}

// RefreshAuth implements BaseClient.RefreshAuth
func (c *BaseClientImpl) RefreshAuth(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Default implementation - should be overridden by specific clients
	c.authInfo.LastRefresh = time.Now()
	c.lastUsed = time.Now()
	return nil
}

// Configure implements BaseClient.Configure
func (c *BaseClientImpl) Configure(config ClientConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config = config
	c.authInfo.Type = config.AuthType
	return nil
}

// GetConfig implements BaseClient.GetConfig
func (c *BaseClientImpl) GetConfig() ClientConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// UpdateConfig implements BaseClient.UpdateConfig
func (c *BaseClientImpl) UpdateConfig(config ClientConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Store old config for comparison
	oldConfig := c.config
	c.config = config
	c.authInfo.Type = config.AuthType

	// If endpoint changed, we should disconnect and reconnect
	if oldConfig.Endpoint != config.Endpoint && c.isConnected {
		c.isConnected = false
		c.metrics.ActiveConnections--
	}

	return nil
}

// HealthCheck implements BaseClient.HealthCheck
func (c *BaseClientImpl) HealthCheck(ctx context.Context) (HealthStatus, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	start := time.Now()

	// Default health check - just check connection status
	status := HealthStatus{
		Healthy:   c.isConnected,
		Message:   "",
		CheckedAt: time.Now(),
		Latency:   0,
	}

	if !c.isConnected {
		status.Message = "Client is not connected"
	}

	status.Latency = time.Since(start)
	c.lastUsed = time.Now()

	c.notifyCallbacks(EventHealthCheck, c)
	return status, nil
}

// GetMetrics implements BaseClient.GetMetrics
func (c *BaseClientImpl) GetMetrics() ClientMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.metrics
}

// GetLastError implements BaseClient.GetLastError
func (c *BaseClientImpl) GetLastError() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastError
}

// ResetMetrics implements BaseClient.ResetMetrics
func (c *BaseClientImpl) ResetMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = ClientMetrics{
		RequestsTotal:     0,
		RequestsFailed:    0,
		AvgLatency:        0,
		LastRequestTime:   time.Time{},
		ConnectionCount:   0,
		ActiveConnections: 0,
	}
}

// GetClientType implements BaseClient.GetClientType
func (c *BaseClientImpl) GetClientType() ClientType {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.Type
}

// GetClientName implements BaseClient.GetClientName
func (c *BaseClientImpl) GetClientName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.Name
}

// GetEndpoint implements BaseClient.GetEndpoint
func (c *BaseClientImpl) GetEndpoint() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.Endpoint
}

// Initialize implements BaseClient.Initialize
func (c *BaseClientImpl) Initialize(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Connect directly without calling Connect() to avoid deadlock
	if !c.isConnected {
		c.isConnected = true
		c.metrics.ConnectionCount++
		c.lastUsed = time.Now()
		c.notifyCallbacks(EventConnected, c)
	}

	// Authenticate directly without calling Authenticate() to avoid deadlock
	if c.config.AuthType != AuthTypeNone {
		c.authInfo.LastRefresh = time.Now()
		c.lastUsed = time.Now()
		c.notifyCallbacks(EventAuthenticated, c)
	}

	return nil
}

// Shutdown implements BaseClient.Shutdown
func (c *BaseClientImpl) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		// Disconnect directly without calling Disconnect() to avoid deadlock
		c.isConnected = false
		c.metrics.ActiveConnections--
		c.lastUsed = time.Now()
		c.notifyCallbacks(EventDisconnected, c)
	}

	c.notifyCallbacks(EventShutdown, c)
	return nil
}

// RegisterCallback registers a callback for client events
func (c *BaseClientImpl) RegisterCallback(callback ClientCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.callbacks = append(c.callbacks, callback)
}

// UnregisterCallback unregisters a callback
func (c *BaseClientImpl) UnregisterCallback(callback ClientCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, cb := range c.callbacks {
		if &cb == &callback {
			c.callbacks = append(c.callbacks[:i], c.callbacks[i+1:]...)
			break
		}
	}
}

// notifyCallbacks notifies all registered callbacks of an event
func (c *BaseClientImpl) notifyCallbacks(event ClientEvent, client BaseClient) {
	for _, callback := range c.callbacks {
		callback(event, client)
	}
}

// recordRequest records a request for metrics
func (c *BaseClientImpl) recordRequest(latency time.Duration, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics.RequestsTotal++
	if !success {
		c.metrics.RequestsFailed++
	}

	// Update average latency
	if c.metrics.RequestsTotal == 1 {
		c.metrics.AvgLatency = latency
	} else {
		// Exponential moving average
		alpha := 0.1
		c.metrics.AvgLatency = time.Duration(float64(c.metrics.AvgLatency)*(1-alpha) + float64(latency)*alpha)
	}

	c.metrics.LastRequestTime = time.Now()
	c.lastUsed = time.Now()
}

// recordError records an error
func (c *BaseClientImpl) recordError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastError = err
	c.notifyCallbacks(EventError, c)
}

// GetCreatedAt returns when the client was created
func (c *BaseClientImpl) GetCreatedAt() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.createdAt
}

// GetLastUsed returns when the client was last used
func (c *BaseClientImpl) GetLastUsed() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastUsed
}

// SetConnected sets the connection status (for use by subclasses)
func (c *BaseClientImpl) SetConnected(connected bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isConnected = connected
	if connected {
		c.metrics.ActiveConnections++
	} else {
		c.metrics.ActiveConnections--
	}
}

// SetAuthInfo sets the authentication info (for use by subclasses)
func (c *BaseClientImpl) SetAuthInfo(authInfo AuthInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.authInfo = authInfo
}

// Errors
var (
	ErrClientNotConnected   = &ClientError{Message: "client is not connected", Code: "NOT_CONNECTED"}
	ErrAuthenticationFailed = &ClientError{Message: "authentication failed", Code: "AUTH_FAILED"}
	ErrInvalidConfiguration = &ClientError{Message: "invalid configuration", Code: "INVALID_CONFIG"}
	ErrOperationFailed      = &ClientError{Message: "operation failed", Code: "OPERATION_FAILED"}
)

// ClientError represents a client error
type ClientError struct {
	Message string
	Code    string
	Details map[string]interface{}
}

func (e *ClientError) Error() string {
	return e.Message
}

func (e *ClientError) WithDetails(details map[string]interface{}) *ClientError {
	e.Details = details
	return e
}
