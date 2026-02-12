package sdk

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Validator validates and sanitizes request data
type Validator struct {
	config *SDKConfig
}

// NewValidator creates a new validator
func NewValidator(config *SDKConfig) *Validator {
	return &Validator{
		config: config,
	}
}

// Validate validates data against validation rules
func (v *Validator) Validate(data map[string]interface{}, rules []ValidationRule) *ValidationResult {
	result := NewValidationResult(true)

	if data == nil {
		result.AddError("", "DATA_REQUIRED", "Data cannot be nil", nil)
		return result
	}

	// Apply default validation rules if none provided
	if len(rules) == 0 {
		rules = v.config.DefaultValidationRules
	}

	// Validate each rule
	for _, rule := range rules {
		value, exists := data[rule.Field]

		// Check required field
		if rule.Required && !exists {
			result.AddError(rule.Field, "FIELD_REQUIRED",
				fmt.Sprintf("Field '%s' is required", rule.Field), nil)
			continue
		}

		// Skip validation if field doesn't exist and is not required
		if !exists {
			continue
		}

		// Apply validation based on rule type
		switch rule.RuleType {
		case ValidationRuleTypeString:
			v.validateString(result, rule, value)
		case ValidationRuleTypeNumber:
			v.validateNumber(result, rule, value)
		case ValidationRuleTypeBoolean:
			v.validateBoolean(result, rule, value)
		case ValidationRuleTypeArray:
			v.validateArray(result, rule, value)
		case ValidationRuleTypeObject:
			v.validateObject(result, rule, value)
		case ValidationRuleTypeEmail:
			v.validateEmail(result, rule, value)
		case ValidationRuleTypeURL:
			v.validateURL(result, rule, value)
		case ValidationRuleTypeUUID:
			v.validateUUID(result, rule, value)
		case ValidationRuleTypeCustom:
			v.validateCustom(result, rule, value)
		}
	}

	// If validation passed and sanitization is enabled, sanitize the data
	if result.Valid && v.config.EnableSanitization {
		result.SanitizedData = v.sanitizeData(data, rules)
	} else {
		result.SanitizedData = data
	}

	return result
}

// Sanitize sanitizes data according to rules
func (v *Validator) Sanitize(data map[string]interface{}, rules []ValidationRule) map[string]interface{} {
	if !v.config.EnableSanitization {
		return data
	}

	return v.sanitizeData(data, rules)
}

// AddDefaults adds default values to data
func (v *Validator) AddDefaults(data map[string]interface{}, rules []ValidationRule) map[string]interface{} {
	// This would be implemented based on default values in rules
	// For now, return the data as-is
	return data
}

// validateString validates a string value
func (v *Validator) validateString(result *ValidationResult, rule ValidationRule, value interface{}) {
	str, ok := value.(string)
	if !ok {
		result.AddError(rule.Field, "INVALID_TYPE",
			fmt.Sprintf("Field '%s' must be a string", rule.Field),
			map[string]interface{}{"actual_type": fmt.Sprintf("%T", value)})
		return
	}

	// Check min length
	if rule.MinLength != nil && len(str) < *rule.MinLength {
		result.AddError(rule.Field, "MIN_LENGTH",
			fmt.Sprintf("Field '%s' must be at least %d characters", rule.Field, *rule.MinLength),
			map[string]interface{}{"min_length": *rule.MinLength, "actual_length": len(str)})
	}

	// Check max length
	if rule.MaxLength != nil && len(str) > *rule.MaxLength {
		result.AddError(rule.Field, "MAX_LENGTH",
			fmt.Sprintf("Field '%s' must be at most %d characters", rule.Field, *rule.MaxLength),
			map[string]interface{}{"max_length": *rule.MaxLength, "actual_length": len(str)})
	}

	// Check pattern
	if rule.Pattern != "" {
		matched, err := regexp.MatchString(rule.Pattern, str)
		if err != nil {
			result.AddError(rule.Field, "INVALID_PATTERN",
				fmt.Sprintf("Invalid pattern for field '%s'", rule.Field),
				map[string]interface{}{"pattern": rule.Pattern, "error": err.Error()})
		} else if !matched {
			result.AddError(rule.Field, "PATTERN_MISMATCH",
				fmt.Sprintf("Field '%s' does not match pattern", rule.Field),
				map[string]interface{}{"pattern": rule.Pattern})
		}
	}

	// Check allowed values
	if len(rule.AllowedValues) > 0 {
		found := false
		for _, allowed := range rule.AllowedValues {
			if str == allowed {
				found = true
				break
			}
		}
		if !found {
			result.AddError(rule.Field, "INVALID_VALUE",
				fmt.Sprintf("Field '%s' must be one of: %s", rule.Field, strings.Join(rule.AllowedValues, ", ")),
				map[string]interface{}{"allowed_values": rule.AllowedValues})
		}
	}
}

// validateNumber validates a number value
func (v *Validator) validateNumber(result *ValidationResult, rule ValidationRule, value interface{}) {
	var num float64
	switch val := value.(type) {
	case int:
		num = float64(val)
	case int64:
		num = float64(val)
	case float32:
		num = float64(val)
	case float64:
		num = val
	default:
		result.AddError(rule.Field, "INVALID_TYPE",
			fmt.Sprintf("Field '%s' must be a number", rule.Field),
			map[string]interface{}{"actual_type": fmt.Sprintf("%T", value)})
		return
	}

	// Check min value
	if rule.MinValue != nil && num < *rule.MinValue {
		result.AddError(rule.Field, "MIN_VALUE",
			fmt.Sprintf("Field '%s' must be at least %v", rule.Field, *rule.MinValue),
			map[string]interface{}{"min_value": *rule.MinValue, "actual_value": num})
	}

	// Check max value
	if rule.MaxValue != nil && num > *rule.MaxValue {
		result.AddError(rule.Field, "MAX_VALUE",
			fmt.Sprintf("Field '%s' must be at most %v", rule.Field, *rule.MaxValue),
			map[string]interface{}{"max_value": *rule.MaxValue, "actual_value": num})
	}
}

// validateBoolean validates a boolean value
func (v *Validator) validateBoolean(result *ValidationResult, rule ValidationRule, value interface{}) {
	_, ok := value.(bool)
	if !ok {
		result.AddError(rule.Field, "INVALID_TYPE",
			fmt.Sprintf("Field '%s' must be a boolean", rule.Field),
			map[string]interface{}{"actual_type": fmt.Sprintf("%T", value)})
	}
}

// validateArray validates an array value
func (v *Validator) validateArray(result *ValidationResult, rule ValidationRule, value interface{}) {
	_, ok := value.([]interface{})
	if !ok {
		result.AddError(rule.Field, "INVALID_TYPE",
			fmt.Sprintf("Field '%s' must be an array", rule.Field),
			map[string]interface{}{"actual_type": fmt.Sprintf("%T", value)})
	}
}

// validateObject validates an object value
func (v *Validator) validateObject(result *ValidationResult, rule ValidationRule, value interface{}) {
	_, ok := value.(map[string]interface{})
	if !ok {
		result.AddError(rule.Field, "INVALID_TYPE",
			fmt.Sprintf("Field '%s' must be an object", rule.Field),
			map[string]interface{}{"actual_type": fmt.Sprintf("%T", value)})
	}
}

// validateEmail validates an email value
func (v *Validator) validateEmail(result *ValidationResult, rule ValidationRule, value interface{}) {
	str, ok := value.(string)
	if !ok {
		result.AddError(rule.Field, "INVALID_TYPE",
			fmt.Sprintf("Field '%s' must be a string", rule.Field),
			map[string]interface{}{"actual_type": fmt.Sprintf("%T", value)})
		return
	}

	_, err := mail.ParseAddress(str)
	if err != nil {
		result.AddError(rule.Field, "INVALID_EMAIL",
			fmt.Sprintf("Field '%s' must be a valid email address", rule.Field),
			map[string]interface{}{"error": err.Error()})
	}
}

// validateURL validates a URL value
func (v *Validator) validateURL(result *ValidationResult, rule ValidationRule, value interface{}) {
	str, ok := value.(string)
	if !ok {
		result.AddError(rule.Field, "INVALID_TYPE",
			fmt.Sprintf("Field '%s' must be a string", rule.Field),
			map[string]interface{}{"actual_type": fmt.Sprintf("%T", value)})
		return
	}

	_, err := url.ParseRequestURI(str)
	if err != nil {
		result.AddError(rule.Field, "INVALID_URL",
			fmt.Sprintf("Field '%s' must be a valid URL", rule.Field),
			map[string]interface{}{"error": err.Error()})
	}
}

// validateUUID validates a UUID value
func (v *Validator) validateUUID(result *ValidationResult, rule ValidationRule, value interface{}) {
	str, ok := value.(string)
	if !ok {
		result.AddError(rule.Field, "INVALID_TYPE",
			fmt.Sprintf("Field '%s' must be a string", rule.Field),
			map[string]interface{}{"actual_type": fmt.Sprintf("%T", value)})
		return
	}

	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(strings.ToLower(str)) {
		result.AddError(rule.Field, "INVALID_UUID",
			fmt.Sprintf("Field '%s' must be a valid UUID", rule.Field), nil)
	}
}

// validateCustom validates using a custom validator
func (v *Validator) validateCustom(result *ValidationResult, rule ValidationRule, value interface{}) {
	// Custom validation would be implemented based on custom validator name
	// For now, just log that custom validation is not implemented
	if rule.CustomValidator != "" {
		result.AddWarning(rule.Field, "CUSTOM_VALIDATOR_NOT_IMPLEMENTED",
			fmt.Sprintf("Custom validator '%s' is not implemented", rule.CustomValidator), nil)
	}
}

// sanitizeData sanitizes data according to rules
func (v *Validator) sanitizeData(data map[string]interface{}, rules []ValidationRule) map[string]interface{} {
	sanitized := make(map[string]interface{})

	for key, value := range data {
		sanitized[key] = v.sanitizeValue(key, value, rules)
	}

	return sanitized
}

// sanitizeValue sanitizes a single value
func (v *Validator) sanitizeValue(key string, value interface{}, rules []ValidationRule) interface{} {
	// Find rule for this field
	var rule *ValidationRule
	for _, r := range rules {
		if r.Field == key {
			rule = &r
			break
		}
	}

	// If no rule found, return value as-is
	if rule == nil {
		return value
	}

	// Apply sanitization based on rule type
	switch rule.RuleType {
	case ValidationRuleTypeString:
		return v.sanitizeString(value)
	case ValidationRuleTypeNumber:
		return v.sanitizeNumber(value)
	case ValidationRuleTypeEmail:
		return v.sanitizeEmail(value)
	case ValidationRuleTypeURL:
		return v.sanitizeURL(value)
	default:
		return value
	}
}

// sanitizeString sanitizes a string value
func (v *Validator) sanitizeString(value interface{}) interface{} {
	str, ok := value.(string)
	if !ok {
		return value
	}

	// Trim whitespace
	str = strings.TrimSpace(str)

	// Remove control characters
	str = regexp.MustCompile(`[\x00-\x1F\x7F]`).ReplaceAllString(str, "")

	return str
}

// sanitizeNumber sanitizes a number value
func (v *Validator) sanitizeNumber(value interface{}) interface{} {
	// For numbers, just return as-is
	// In a real implementation, you might round or format numbers
	return value
}

// sanitizeEmail sanitizes an email value
func (v *Validator) sanitizeEmail(value interface{}) interface{} {
	str, ok := value.(string)
	if !ok {
		return value
	}

	// Trim and lowercase email
	str = strings.TrimSpace(strings.ToLower(str))

	return str
}

// sanitizeURL sanitizes a URL value
func (v *Validator) sanitizeURL(value interface{}) interface{} {
	str, ok := value.(string)
	if !ok {
		return value
	}

	// Trim URL
	str = strings.TrimSpace(str)

	return str
}

// ValidateHTTPRequest validates an HTTP request context
func (v *Validator) ValidateHTTPRequest(httpContext *HTTPRequestContext) *ValidationResult {
	result := NewValidationResult(true)

	if httpContext == nil {
		result.AddError("http_request", "HTTP_CONTEXT_REQUIRED", "HTTP request context is required", nil)
		return result
	}

	// Validate method
	if httpContext.Method == "" {
		result.AddError("method", "METHOD_REQUIRED", "HTTP method is required", nil)
	} else {
		allowedMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
		validMethod := false
		for _, method := range allowedMethods {
			if httpContext.Method == method {
				validMethod = true
				break
			}
		}
		if !validMethod {
			result.AddError("method", "INVALID_METHOD",
				fmt.Sprintf("Invalid HTTP method: %s", httpContext.Method),
				map[string]interface{}{"allowed_methods": allowedMethods})
		}
	}

	// Validate path
	if httpContext.Path == "" {
		result.AddError("path", "PATH_REQUIRED", "HTTP path is required", nil)
	}

	// Validate timestamp
	if httpContext.Timestamp.IsZero() {
		result.AddError("timestamp", "TIMESTAMP_REQUIRED", "Timestamp is required", nil)
	} else if httpContext.Timestamp.After(time.Now().Add(5 * time.Minute)) {
		result.AddWarning("timestamp", "FUTURE_TIMESTAMP",
			"Timestamp is in the future",
			map[string]interface{}{"timestamp": httpContext.Timestamp})
	}

	return result
}

// ValidateSessionContext validates a session context
func (v *Validator) ValidateSessionContext(session *SessionContext) *ValidationResult {
	result := NewValidationResult(true)

	if session == nil {
		return result
	}

	// Validate user ID if present
	if session.UserID != "" && len(session.UserID) > 255 {
		result.AddError("user_id", "USER_ID_TOO_LONG",
			"User ID must be at most 255 characters",
			map[string]interface{}{"length": len(session.UserID)})
	}

	// Validate session ID if present
	if session.SessionID != "" && len(session.SessionID) > 255 {
		result.AddError("session_id", "SESSION_ID_TOO_LONG",
			"Session ID must be at most 255 characters",
			map[string]interface{}{"length": len(session.SessionID)})
	}

	// Validate auth method if present
	if session.AuthMethod != "" {
		allowedMethods := []string{"jwt", "api_key", "basic", "oauth"}
		validMethod := false
		for _, method := range allowedMethods {
			if session.AuthMethod == method {
				validMethod = true
				break
			}
		}
		if !validMethod {
			result.AddError("auth_method", "INVALID_AUTH_METHOD",
				fmt.Sprintf("Invalid auth method: %s", session.AuthMethod),
				map[string]interface{}{"allowed_methods": allowedMethods})
		}
	}

	return result
}
