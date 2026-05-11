package sdk

import (
	"fmt"
)

// SDKError represents an SDK-specific error
type SDKError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Field   string                 `json:"field,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
	Err     error                  `json:"-"`
}

// Error implements the error interface
func (e *SDKError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *SDKError) Unwrap() error {
	return e.Err
}

// Common SDK error codes
const (
	ErrCodeInvalidConfig        = "INVALID_CONFIG"
	ErrCodeValidationFailed     = "VALIDATION_FAILED"
	ErrCodeRequestParsingFailed = "REQUEST_PARSING_FAILED"
	ErrCodeSessionExtraction    = "SESSION_EXTRACTION_FAILED"
	ErrCodeSecurityContext      = "SECURITY_CONTEXT_FAILED"
	ErrCodeWorkflowExecution    = "WORKFLOW_EXECUTION_FAILED"
	ErrCodeHTTPRequest          = "HTTP_REQUEST_FAILED"
	ErrCodeTimeout              = "TIMEOUT"
	ErrCodeCircuitBreaker       = "CIRCUIT_BREAKER_OPEN"
	ErrCodeRetryExhausted       = "RETRY_EXHAUSTED"
)

// NewSDKError creates a new SDK error
func NewSDKError(code, message string) *SDKError {
	return &SDKError{
		Code:    code,
		Message: message,
	}
}

// NewSDKErrorWithField creates a new SDK error with field information
func NewSDKErrorWithField(code, message, field string) *SDKError {
	return &SDKError{
		Code:    code,
		Message: message,
		Field:   field,
	}
}

// NewSDKErrorWithDetails creates a new SDK error with details
func NewSDKErrorWithDetails(code, message string, details map[string]interface{}) *SDKError {
	return &SDKError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// WrapSDKError wraps an existing error as an SDK error
func WrapSDKError(err error, code, message string) *SDKError {
	return &SDKError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string                 `json:"field"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s (field: %s)", e.Code, e.Message, e.Field)
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid         bool                   `json:"valid"`
	Errors        []ValidationError      `json:"errors,omitempty"`
	Warnings      []ValidationError      `json:"warnings,omitempty"`
	SanitizedData map[string]interface{} `json:"sanitized_data,omitempty"`
}

// NewValidationResult creates a new validation result
func NewValidationResult(valid bool) *ValidationResult {
	return &ValidationResult{
		Valid:         valid,
		Errors:        []ValidationError{},
		Warnings:      []ValidationError{},
		SanitizedData: make(map[string]interface{}),
	}
}

// AddError adds a validation error
func (vr *ValidationResult) AddError(field, code, message string, details map[string]interface{}) {
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Code:    code,
		Message: message,
		Details: details,
	})
	vr.Valid = false
}

// AddWarning adds a validation warning
func (vr *ValidationResult) AddWarning(field, code, message string, details map[string]interface{}) {
	vr.Warnings = append(vr.Warnings, ValidationError{
		Field:   field,
		Code:    code,
		Message: message,
		Details: details,
	})
}

// HasErrors returns true if there are validation errors
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are validation warnings
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// ErrorMessages returns all error messages as a slice
func (vr *ValidationResult) ErrorMessages() []string {
	messages := make([]string, len(vr.Errors))
	for i, err := range vr.Errors {
		messages[i] = err.Error()
	}
	return messages
}

// WarningMessages returns all warning messages as a slice
func (vr *ValidationResult) WarningMessages() []string {
	messages := make([]string, len(vr.Warnings))
	for i, warning := range vr.Warnings {
		messages[i] = warning.Error()
	}
	return messages
}
