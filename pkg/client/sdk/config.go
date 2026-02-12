package sdk

import (
	"time"
)

// SDKConfig represents the configuration for the Workflow SDK client
type SDKConfig struct {
	// Workflow API configuration
	WorkflowAPIEndpoint string        `json:"workflow_api_endpoint" yaml:"workflow_api_endpoint"`
	Timeout             time.Duration `json:"timeout" yaml:"timeout"`
	MaxRetries          int           `json:"max_retries" yaml:"max_retries"`
	RetryDelay          time.Duration `json:"retry_delay" yaml:"retry_delay"`

	// Authentication
	AuthToken string   `json:"auth_token" yaml:"auth_token"`
	AuthType  AuthType `json:"auth_type" yaml:"auth_type"`

	// Validation
	EnableValidation   bool `json:"enable_validation" yaml:"enable_validation"`
	EnableSanitization bool `json:"enable_sanitization" yaml:"enable_sanitization"`
	StrictValidation   bool `json:"strict_validation" yaml:"strict_validation"`

	// Context extraction
	EnableSessionExtraction bool `json:"enable_session_extraction" yaml:"enable_session_extraction"`
	EnableSecurityContext   bool `json:"enable_security_context" yaml:"enable_security_context"`
	IncludeFullHTTPContext  bool `json:"include_full_http_context" yaml:"include_full_http_context"`

	// Logging and monitoring
	LogLevel             LogLevel `json:"log_level" yaml:"log_level"`
	EnableRequestLogging bool     `json:"enable_request_logging" yaml:"enable_request_logging"`
	EnableMetrics        bool     `json:"enable_metrics" yaml:"enable_metrics"`

	// Default validation rules
	DefaultValidationRules []ValidationRule `json:"default_validation_rules" yaml:"default_validation_rules"`
	CustomValidators       []string         `json:"custom_validators" yaml:"custom_validators"`
}

// AuthType represents the authentication type
type AuthType string

const (
	AuthTypeBearerToken AuthType = "bearer_token"
	AuthTypeAPIKey      AuthType = "api_key"
	AuthTypeBasicAuth   AuthType = "basic_auth"
	AuthTypeOAuth2      AuthType = "oauth2"
	AuthTypeAWSSigV4    AuthType = "aws_sigv4"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// ValidationRule defines a validation rule for request data
type ValidationRule struct {
	Field           string             `json:"field"`
	RuleType        ValidationRuleType `json:"rule_type"`
	Required        bool               `json:"required,omitempty"`
	MinLength       *int               `json:"min_length,omitempty"`
	MaxLength       *int               `json:"max_length,omitempty"`
	Pattern         string             `json:"pattern,omitempty"`
	MinValue        *float64           `json:"min_value,omitempty"`
	MaxValue        *float64           `json:"max_value,omitempty"`
	AllowedValues   []string           `json:"allowed_values,omitempty"`
	CustomValidator string             `json:"custom_validator,omitempty"`
}

// ValidationRuleType represents the type of validation rule
type ValidationRuleType string

const (
	ValidationRuleTypeRequired ValidationRuleType = "required"
	ValidationRuleTypeString   ValidationRuleType = "string"
	ValidationRuleTypeNumber   ValidationRuleType = "number"
	ValidationRuleTypeBoolean  ValidationRuleType = "boolean"
	ValidationRuleTypeArray    ValidationRuleType = "array"
	ValidationRuleTypeObject   ValidationRuleType = "object"
	ValidationRuleTypeEmail    ValidationRuleType = "email"
	ValidationRuleTypeURL      ValidationRuleType = "url"
	ValidationRuleTypeUUID     ValidationRuleType = "uuid"
	ValidationRuleTypeCustom   ValidationRuleType = "custom"
)

// DefaultConfig returns the default SDK configuration
func DefaultConfig() SDKConfig {
	return SDKConfig{
		WorkflowAPIEndpoint:     "http://localhost:8080",
		Timeout:                 30 * time.Second,
		MaxRetries:              3,
		RetryDelay:              1 * time.Second,
		AuthType:                AuthTypeBearerToken,
		EnableValidation:        true,
		EnableSanitization:      true,
		StrictValidation:        false,
		EnableSessionExtraction: true,
		EnableSecurityContext:   true,
		IncludeFullHTTPContext:  true,
		LogLevel:                LogLevelInfo,
		EnableRequestLogging:    true,
		EnableMetrics:           true,
		DefaultValidationRules:  []ValidationRule{},
		CustomValidators:        []string{},
	}
}

// Validate validates the SDK configuration
func (c *SDKConfig) Validate() error {
	if c.WorkflowAPIEndpoint == "" {
		return &SDKError{
			Code:    "INVALID_CONFIG",
			Message: "WorkflowAPIEndpoint is required",
			Field:   "workflow_api_endpoint",
		}
	}

	if c.Timeout <= 0 {
		return &SDKError{
			Code:    "INVALID_CONFIG",
			Message: "Timeout must be positive",
			Field:   "timeout",
		}
	}

	if c.MaxRetries < 0 {
		return &SDKError{
			Code:    "INVALID_CONFIG",
			Message: "MaxRetries cannot be negative",
			Field:   "max_retries",
		}
	}

	if c.RetryDelay < 0 {
		return &SDKError{
			Code:    "INVALID_CONFIG",
			Message: "RetryDelay cannot be negative",
			Field:   "retry_delay",
		}
	}

	return nil
}
