package sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unified-workflow/internal/config"

	"github.com/goccy/go-yaml"
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

	// Execution configuration
	AsyncExecution          bool          `json:"async_execution" yaml:"async_execution"`
	DefaultPriority         int           `json:"default_priority" yaml:"default_priority"`
	PollIntervalMs          int           `json:"poll_interval_ms" yaml:"poll_interval_ms"`
	EstimatedCompletionMs   int           `json:"estimated_completion_ms" yaml:"estimated_completion_ms"`
	ExecutionExpiryDuration time.Duration `json:"execution_expiry_duration" yaml:"execution_expiry_duration"`
	SDKVersion              string        `json:"sdk_version" yaml:"sdk_version"`

	// Circuit breaker configuration
	EnableCircuitBreaker    bool          `json:"enable_circuit_breaker" yaml:"enable_circuit_breaker"`
	CircuitBreakerThreshold int           `json:"circuit_breaker_threshold" yaml:"circuit_breaker_threshold"`
	CircuitBreakerTimeout   time.Duration `json:"circuit_breaker_timeout" yaml:"circuit_breaker_timeout"`
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
		AsyncExecution:          true,
		DefaultPriority:         5,
		PollIntervalMs:          1000,
		EstimatedCompletionMs:   5000,
		ExecutionExpiryDuration: 1 * time.Hour,
		SDKVersion:              "1.0.0",
		EnableCircuitBreaker:    true,
		CircuitBreakerThreshold: 5,
		CircuitBreakerTimeout:   60 * time.Second,
	}
}

// NewConfigFromAppConfig creates SDK configuration from the main application config
func NewConfigFromAppConfig(appConfig interface{}) SDKConfig {
	// Try to cast to the main config type
	if cfg, ok := appConfig.(*config.Config); ok {
		return SDKConfig{
			WorkflowAPIEndpoint:     cfg.Clients.SDK.WorkflowAPIEndpoint,
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

	// Fall back to default config
	return DefaultConfig()
}

// LoadSDKConfig loads SDK configuration from the main config file
func LoadSDKConfig() (SDKConfig, error) {
	// Load the main application config
	appConfig, err := config.LoadConfig()
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to load application config: %w", err)
	}

	// Create SDK config from app config
	sdkConfig := NewConfigFromAppConfig(appConfig)

	// Apply environment variable overrides for SDK-specific settings
	applySDKEnvOverrides(&sdkConfig)

	return sdkConfig, nil
}

// LoadConfigFromFile loads SDK configuration from a YAML or JSON file
func LoadConfigFromFile(filePath string) (SDKConfig, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to read config file: %w", err)
	}

	// Determine file type by extension
	ext := strings.ToLower(filepath.Ext(filePath))
	var sdkConfig SDKConfig

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &sdkConfig); err != nil {
			return DefaultConfig(), fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &sdkConfig); err != nil {
			return DefaultConfig(), fmt.Errorf("failed to parse JSON config: %w", err)
		}
	default:
		return DefaultConfig(), fmt.Errorf("unsupported config file format: %s. Use .yaml, .yml, or .json", ext)
	}

	// Apply environment variable overrides
	applySDKEnvOverrides(&sdkConfig)

	// Validate the configuration
	if err := sdkConfig.Validate(); err != nil {
		return DefaultConfig(), fmt.Errorf("invalid configuration: %w", err)
	}

	return sdkConfig, nil
}

// LoadDefaultConfig loads SDK configuration from default locations
// Looks for config files in the following order:
// 1. Environment variables
// 2. sdk-config.yaml in current directory
// 3. sdk-config.yml in current directory
// 4. sdk-config.json in current directory
// 5. Default configuration
func LoadDefaultConfig() (SDKConfig, error) {
	// Try to load from environment variables first
	if envConfig, err := LoadConfigFromEnvironment(); err == nil {
		return envConfig, nil
	}

	// Try config files in order
	configFiles := []string{
		"sdk-config.yaml",
		"sdk-config.yml",
		"sdk-config.json",
	}

	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			return LoadConfigFromFile(configFile)
		}
	}

	// Fall back to default configuration
	return DefaultConfig(), nil
}

// LoadConfigFromEnvironment loads SDK configuration from environment variables only
func LoadConfigFromEnvironment() (SDKConfig, error) {
	sdkConfig := DefaultConfig()
	applySDKEnvOverrides(&sdkConfig)

	// Validate the configuration
	if err := sdkConfig.Validate(); err != nil {
		return DefaultConfig(), fmt.Errorf("invalid environment configuration: %w", err)
	}

	return sdkConfig, nil
}

// applySDKEnvOverrides applies environment variable overrides to SDK configuration
func applySDKEnvOverrides(sdkConfig *SDKConfig) {
	if val := os.Getenv("SDK_WORKFLOW_API_ENDPOINT"); val != "" {
		sdkConfig.WorkflowAPIEndpoint = val
	}
	if val := os.Getenv("SDK_TIMEOUT"); val != "" {
		if timeout, err := strconv.Atoi(val); err == nil {
			sdkConfig.Timeout = time.Duration(timeout) * time.Second
		}
	}
	if val := os.Getenv("SDK_MAX_RETRIES"); val != "" {
		if retries, err := strconv.Atoi(val); err == nil {
			sdkConfig.MaxRetries = retries
		}
	}
	if val := os.Getenv("SDK_RETRY_DELAY"); val != "" {
		if retryDelay, err := strconv.Atoi(val); err == nil {
			sdkConfig.RetryDelay = time.Duration(retryDelay) * time.Second
		}
	}
	if val := os.Getenv("SDK_AUTH_TOKEN"); val != "" {
		sdkConfig.AuthToken = val
	}
	if val := os.Getenv("SDK_AUTH_TYPE"); val != "" {
		sdkConfig.AuthType = AuthType(val)
	}
	if val := os.Getenv("SDK_ENABLE_VALIDATION"); val != "" {
		sdkConfig.EnableValidation = strings.ToLower(val) == "true"
	}
	if val := os.Getenv("SDK_ENABLE_SANITIZATION"); val != "" {
		sdkConfig.EnableSanitization = strings.ToLower(val) == "true"
	}
	if val := os.Getenv("SDK_STRICT_VALIDATION"); val != "" {
		sdkConfig.StrictValidation = strings.ToLower(val) == "true"
	}
	if val := os.Getenv("SDK_LOG_LEVEL"); val != "" {
		sdkConfig.LogLevel = LogLevel(val)
	}
	if val := os.Getenv("SDK_ENABLE_REQUEST_LOGGING"); val != "" {
		sdkConfig.EnableRequestLogging = strings.ToLower(val) == "true"
	}
	if val := os.Getenv("SDK_ENABLE_METRICS"); val != "" {
		sdkConfig.EnableMetrics = strings.ToLower(val) == "true"
	}
	if val := os.Getenv("SDK_ASYNC_EXECUTION"); val != "" {
		sdkConfig.AsyncExecution = strings.ToLower(val) == "true"
	}
	if val := os.Getenv("SDK_DEFAULT_PRIORITY"); val != "" {
		if priority, err := strconv.Atoi(val); err == nil {
			sdkConfig.DefaultPriority = priority
		}
	}
	if val := os.Getenv("SDK_POLL_INTERVAL_MS"); val != "" {
		if pollInterval, err := strconv.Atoi(val); err == nil {
			sdkConfig.PollIntervalMs = pollInterval
		}
	}
	if val := os.Getenv("SDK_ESTIMATED_COMPLETION_MS"); val != "" {
		if estimatedCompletion, err := strconv.Atoi(val); err == nil {
			sdkConfig.EstimatedCompletionMs = estimatedCompletion
		}
	}
	if val := os.Getenv("SDK_EXECUTION_EXPIRY_DURATION"); val != "" {
		if expiryDuration, err := strconv.Atoi(val); err == nil {
			sdkConfig.ExecutionExpiryDuration = time.Duration(expiryDuration) * time.Second
		}
	}
	if val := os.Getenv("SDK_VERSION"); val != "" {
		sdkConfig.SDKVersion = val
	}
	if val := os.Getenv("SDK_ENABLE_CIRCUIT_BREAKER"); val != "" {
		sdkConfig.EnableCircuitBreaker = strings.ToLower(val) == "true"
	}
	if val := os.Getenv("SDK_CIRCUIT_BREAKER_THRESHOLD"); val != "" {
		if threshold, err := strconv.Atoi(val); err == nil {
			sdkConfig.CircuitBreakerThreshold = threshold
		}
	}
	if val := os.Getenv("SDK_CIRCUIT_BREAKER_TIMEOUT"); val != "" {
		if timeout, err := strconv.Atoi(val); err == nil {
			sdkConfig.CircuitBreakerTimeout = time.Duration(timeout) * time.Second
		}
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
