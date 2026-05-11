package sdk

import (
	"time"
)

// HTTPRequestContext represents the complete HTTP request information
type HTTPRequestContext struct {
	Method      string              `json:"method"`
	Path        string              `json:"path"`
	Headers     map[string][]string `json:"headers"`
	QueryParams map[string][]string `json:"query_params,omitempty"`
	PathParams  map[string]string   `json:"path_params,omitempty"`
	Body        interface{}         `json:"body,omitempty"`
	RemoteAddr  string              `json:"remote_addr,omitempty"`
	UserAgent   string              `json:"user_agent,omitempty"`
	Timestamp   time.Time           `json:"timestamp"`
}

// SessionContext represents user session information
type SessionContext struct {
	UserID      string                 `json:"user_id,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	Roles       []string               `json:"roles,omitempty"`
	Permissions []string               `json:"permissions,omitempty"`
	AuthMethod  string                 `json:"auth_method,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// SecurityContext represents security-related information
type SecurityContext struct {
	Authenticated        bool                   `json:"authenticated"`
	AuthenticationMethod string                 `json:"authentication_method,omitempty"`
	Scopes               []string               `json:"scopes,omitempty"`
	Claims               map[string]interface{} `json:"claims,omitempty"`
	IPAddress            string                 `json:"ip_address,omitempty"`
	UserAgent            string                 `json:"user_agent,omitempty"`
	GeoLocation          *GeoLocation           `json:"geo_location,omitempty"`
}

// GeoLocation represents geographic location information
type GeoLocation struct {
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
	City      string  `json:"city,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// SDKExecuteWorkflowRequest represents the enhanced execution request
type SDKExecuteWorkflowRequest struct {
	// Workflow execution parameters
	InputData         map[string]interface{} `json:"input_data,omitempty"`
	CallbackURL       string                 `json:"callback_url,omitempty"`
	TimeoutMs         int64                  `json:"timeout_ms,omitempty"`
	WaitForCompletion bool                   `json:"wait_for_completion,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`

	// SDK-specific extensions
	HTTPRequest        *HTTPRequestContext `json:"http_request,omitempty"`
	Session            *SessionContext     `json:"session,omitempty"`
	Security           *SecurityContext    `json:"security,omitempty"`
	ValidationRules    []ValidationRule    `json:"validation_rules,omitempty"`
	EnableValidation   bool                `json:"enable_validation,omitempty"`
	EnableSanitization bool                `json:"enable_sanitization,omitempty"`
	IncludeFullContext bool                `json:"include_full_context,omitempty"`
}

// SDKExecuteWorkflowResponse represents the enhanced execution response
type SDKExecuteWorkflowResponse struct {
	// Base response fields
	RunID                 string    `json:"run_id"`
	Status                string    `json:"status"`
	Message               string    `json:"message"`
	StatusURL             string    `json:"status_url"`
	ResultURL             string    `json:"result_url"`
	PollAfterMs           int       `json:"poll_after_ms,omitempty"`
	EstimatedCompletionMs int       `json:"estimated_completion_ms,omitempty"`
	ExpiresAt             time.Time `json:"expires_at"`

	// SDK-specific extensions
	ValidationResult *ValidationResult `json:"validation_result,omitempty"`
	ContextIncluded  bool              `json:"context_included"`
	SDKVersion       string            `json:"sdk_version"`
	RequestID        string            `json:"request_id"`
}

// BatchExecutionItem represents a single workflow execution in a batch
type BatchExecutionItem struct {
	WorkflowID string                     `json:"workflow_id"`
	Request    *SDKExecuteWorkflowRequest `json:"request,omitempty"`
	Priority   int                        `json:"priority,omitempty"`
}

// BatchExecuteWorkflowsRequest represents a batch execution request
type BatchExecuteWorkflowsRequest struct {
	Executions         []BatchExecutionItem `json:"executions"`
	Parallel           bool                 `json:"parallel,omitempty"`
	MaxConcurrent      int                  `json:"max_concurrent,omitempty"`
	StopOnFirstFailure bool                 `json:"stop_on_first_failure,omitempty"`
}

// BatchExecutionResult represents the result of a single batch execution
type BatchExecutionResult struct {
	WorkflowID      string           `json:"workflow_id"`
	Success         bool             `json:"success"`
	RunID           string           `json:"run_id,omitempty"`
	Error           *ValidationError `json:"error,omitempty"`
	ExecutionTimeMs int64            `json:"execution_time_ms,omitempty"`
}

// BatchExecuteWorkflowsResponse represents the batch execution response
type BatchExecuteWorkflowsResponse struct {
	BatchID    string                 `json:"batch_id"`
	Total      int                    `json:"total"`
	Successful int                    `json:"successful"`
	Failed     int                    `json:"failed"`
	Pending    int                    `json:"pending"`
	Executions []BatchExecutionResult `json:"executions,omitempty"`
	Errors     []BatchError           `json:"errors,omitempty"`
}

// BatchError represents an error in batch execution
type BatchError struct {
	WorkflowID string           `json:"workflow_id"`
	Error      *ValidationError `json:"error"`
	Timestamp  time.Time        `json:"timestamp"`
}

// WebhookConfiguration represents webhook configuration
type WebhookConfiguration struct {
	WebhookID  string              `json:"webhook_id"`
	URL        string              `json:"url"`
	Events     []WebhookEvent      `json:"events"`
	Secret     string              `json:"secret,omitempty"`
	Enabled    bool                `json:"enabled,omitempty"`
	RetryCount int                 `json:"retry_count,omitempty"`
	TimeoutMs  int                 `json:"timeout_ms,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
}

// WebhookEvent represents a webhook event type
type WebhookEvent string

const (
	WebhookEventWorkflowStarted   WebhookEvent = "workflow_started"
	WebhookEventWorkflowCompleted WebhookEvent = "workflow_completed"
	WebhookEventWorkflowFailed    WebhookEvent = "workflow_failed"
	WebhookEventWorkflowCancelled WebhookEvent = "workflow_cancelled"
	WebhookEventStepCompleted     WebhookEvent = "step_completed"
	WebhookEventStepFailed        WebhookEvent = "step_failed"
)

// ListWebhooksResponse represents the response for listing webhooks
type ListWebhooksResponse struct {
	Webhooks []WebhookConfiguration `json:"webhooks"`
	Count    int                    `json:"count"`
}

// Helper functions

// NewHTTPRequestContext creates a new HTTP request context
func NewHTTPRequestContext(method, path string) *HTTPRequestContext {
	return &HTTPRequestContext{
		Method:      method,
		Path:        path,
		Headers:     make(map[string][]string),
		QueryParams: make(map[string][]string),
		PathParams:  make(map[string]string),
		Timestamp:   time.Now(),
	}
}

// NewSessionContext creates a new session context
func NewSessionContext(userID, sessionID string) *SessionContext {
	return &SessionContext{
		UserID:      userID,
		SessionID:   sessionID,
		Roles:       []string{},
		Permissions: []string{},
		Attributes:  make(map[string]interface{}),
	}
}

// NewSecurityContext creates a new security context
func NewSecurityContext(authenticated bool) *SecurityContext {
	return &SecurityContext{
		Authenticated: authenticated,
		Scopes:        []string{},
		Claims:        make(map[string]interface{}),
	}
}

// NewSDKExecuteWorkflowRequest creates a new SDK execution request
func NewSDKExecuteWorkflowRequest(inputData map[string]interface{}) *SDKExecuteWorkflowRequest {
	return &SDKExecuteWorkflowRequest{
		InputData:          inputData,
		Metadata:           make(map[string]interface{}),
		ValidationRules:    []ValidationRule{},
		EnableValidation:   true,
		EnableSanitization: true,
		IncludeFullContext: true,
	}
}

// AddHeader adds a header to the HTTP request context
func (h *HTTPRequestContext) AddHeader(key, value string) {
	if h.Headers == nil {
		h.Headers = make(map[string][]string)
	}
	h.Headers[key] = append(h.Headers[key], value)
}

// AddQueryParam adds a query parameter to the HTTP request context
func (h *HTTPRequestContext) AddQueryParam(key, value string) {
	if h.QueryParams == nil {
		h.QueryParams = make(map[string][]string)
	}
	h.QueryParams[key] = append(h.QueryParams[key], value)
}

// AddPathParam adds a path parameter to the HTTP request context
func (h *HTTPRequestContext) AddPathParam(key, value string) {
	if h.PathParams == nil {
		h.PathParams = make(map[string]string)
	}
	h.PathParams[key] = value
}

// AddRole adds a role to the session context
func (s *SessionContext) AddRole(role string) {
	s.Roles = append(s.Roles, role)
}

// AddPermission adds a permission to the session context
func (s *SessionContext) AddPermission(permission string) {
	s.Permissions = append(s.Permissions, permission)
}

// AddAttribute adds an attribute to the session context
func (s *SessionContext) AddAttribute(key string, value interface{}) {
	if s.Attributes == nil {
		s.Attributes = make(map[string]interface{})
	}
	s.Attributes[key] = value
}

// AddClaim adds a claim to the security context
func (s *SecurityContext) AddClaim(key string, value interface{}) {
	if s.Claims == nil {
		s.Claims = make(map[string]interface{})
	}
	s.Claims[key] = value
}

// AddScope adds a scope to the security context
func (s *SecurityContext) AddScope(scope string) {
	s.Scopes = append(s.Scopes, scope)
}
