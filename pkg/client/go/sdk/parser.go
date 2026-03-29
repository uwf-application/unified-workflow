package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequestParser parses HTTP requests into SDK data structures
type RequestParser struct {
	config *SDKConfig
}

// NewRequestParser creates a new request parser
func NewRequestParser(config *SDKConfig) *RequestParser {
	return &RequestParser{
		config: config,
	}
}

// ParseHTTPRequest parses an HTTP request into HTTPRequestContext
func (p *RequestParser) ParseHTTPRequest(req *http.Request) (*HTTPRequestContext, error) {
	if req == nil {
		return nil, NewSDKError(ErrCodeRequestParsingFailed, "HTTP request is nil")
	}

	// Parse request body
	var body interface{}
	if req.Body != nil && req.ContentLength > 0 {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, WrapSDKError(err, ErrCodeRequestParsingFailed, "Failed to read request body")
		}

		// Try to parse as JSON, otherwise keep as string
		var jsonBody interface{}
		if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
			body = jsonBody
		} else {
			body = string(bodyBytes)
		}

		// Restore the request body for potential reuse
		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}

	// Parse query parameters
	queryParams := make(map[string][]string)
	if req.URL != nil {
		for key, values := range req.URL.Query() {
			queryParams[key] = values
		}
	}

	// Parse headers
	headers := make(map[string][]string)
	for key, values := range req.Header {
		headers[key] = values
	}

	// Get remote address
	remoteAddr := req.RemoteAddr
	if forwardedFor := req.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		remoteAddr = forwardedFor
	}

	// Get user agent
	userAgent := req.UserAgent()

	return &HTTPRequestContext{
		Method:      req.Method,
		Path:        req.URL.Path,
		Headers:     headers,
		QueryParams: queryParams,
		Body:        body,
		RemoteAddr:  remoteAddr,
		UserAgent:   userAgent,
		Timestamp:   time.Now(),
	}, nil
}

// ExtractSessionInfo extracts session information from HTTP request
func (p *RequestParser) ExtractSessionInfo(req *http.Request) (*SessionContext, error) {
	if !p.config.EnableSessionExtraction {
		return nil, nil
	}

	if req == nil {
		return nil, NewSDKError(ErrCodeSessionExtraction, "HTTP request is nil")
	}

	session := &SessionContext{
		Roles:       []string{},
		Permissions: []string{},
		Attributes:  make(map[string]interface{}),
	}

	// Extract from Authorization header
	authHeader := req.Header.Get("Authorization")
	if authHeader != "" {
		// Parse JWT token if present
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			session.AuthMethod = "jwt"

			// In a real implementation, you would decode and validate the JWT
			// For now, we'll extract basic information
			session.SessionID = extractSessionIDFromToken(token)

			// Extract claims from token (simplified)
			if claims, err := extractClaimsFromToken(token); err == nil {
				if userID, ok := claims["sub"].(string); ok {
					session.UserID = userID
				}
				if roles, ok := claims["roles"].([]interface{}); ok {
					for _, role := range roles {
						if roleStr, ok := role.(string); ok {
							session.Roles = append(session.Roles, roleStr)
						}
					}
				}
			}
		} else if strings.HasPrefix(authHeader, "Basic ") {
			session.AuthMethod = "basic"
		} else if strings.HasPrefix(authHeader, "ApiKey ") {
			session.AuthMethod = "api_key"
		}
	}

	// Extract from cookies
	for _, cookie := range req.Cookies() {
		if cookie.Name == "session_id" || cookie.Name == "SESSIONID" {
			session.SessionID = cookie.Value
		}
		if cookie.Name == "user_id" {
			session.UserID = cookie.Value
		}
	}

	// Extract from custom headers
	if userID := req.Header.Get("X-User-ID"); userID != "" {
		session.UserID = userID
	}
	if sessionID := req.Header.Get("X-Session-ID"); sessionID != "" {
		session.SessionID = sessionID
	}
	if roles := req.Header.Get("X-User-Roles"); roles != "" {
		session.Roles = append(session.Roles, strings.Split(roles, ",")...)
	}

	return session, nil
}

// ExtractSecurityContext extracts security context from HTTP request
func (p *RequestParser) ExtractSecurityContext(req *http.Request) (*SecurityContext, error) {
	if !p.config.EnableSecurityContext {
		return nil, nil
	}

	if req == nil {
		return nil, NewSDKError(ErrCodeSecurityContext, "HTTP request is nil")
	}

	security := NewSecurityContext(false)

	// Check if request is authenticated
	authHeader := req.Header.Get("Authorization")
	if authHeader != "" {
		security.Authenticated = true

		if strings.HasPrefix(authHeader, "Bearer ") {
			security.AuthenticationMethod = "jwt"
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Extract claims from token
			if claims, err := extractClaimsFromToken(token); err == nil {
				security.Claims = claims

				// Extract scopes
				if scopes, ok := claims["scope"].(string); ok {
					security.Scopes = strings.Split(scopes, " ")
				} else if scopes, ok := claims["scopes"].([]interface{}); ok {
					for _, scope := range scopes {
						if scopeStr, ok := scope.(string); ok {
							security.Scopes = append(security.Scopes, scopeStr)
						}
					}
				}
			}
		} else if strings.HasPrefix(authHeader, "Basic ") {
			security.AuthenticationMethod = "basic"
		} else if strings.HasPrefix(authHeader, "ApiKey ") {
			security.AuthenticationMethod = "api_key"
		}
	}

	// Extract IP address
	security.IPAddress = req.RemoteAddr
	if forwardedFor := req.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		security.IPAddress = forwardedFor
	}

	// Extract user agent
	security.UserAgent = req.UserAgent()

	// Extract geo location from headers (if available)
	if country := req.Header.Get("X-Geo-Country"); country != "" {
		if security.GeoLocation == nil {
			security.GeoLocation = &GeoLocation{}
		}
		security.GeoLocation.Country = country
	}
	if region := req.Header.Get("X-Geo-Region"); region != "" {
		if security.GeoLocation == nil {
			security.GeoLocation = &GeoLocation{}
		}
		security.GeoLocation.Region = region
	}

	return security, nil
}

// ParseRequestBody parses the HTTP request body into a map
func (p *RequestParser) ParseRequestBody(req *http.Request) (map[string]interface{}, error) {
	if req == nil {
		return nil, NewSDKError(ErrCodeRequestParsingFailed, "HTTP request is nil")
	}

	if req.Body == nil || req.ContentLength == 0 {
		return make(map[string]interface{}), nil
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, WrapSDKError(err, ErrCodeRequestParsingFailed, "Failed to read request body")
	}

	// Restore the request body for potential reuse
	req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		// If not JSON, return as raw string
		data = map[string]interface{}{
			"raw": string(bodyBytes),
		}
	}

	return data, nil
}

// ParseQueryParams parses query parameters from URL
func (p *RequestParser) ParseQueryParams(rawURL string) (map[string][]string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, WrapSDKError(err, ErrCodeRequestParsingFailed, "Failed to parse URL")
	}

	queryParams := make(map[string][]string)
	for key, values := range parsedURL.Query() {
		queryParams[key] = values
	}

	return queryParams, nil
}

// ParsePathParams parses path parameters from URL pattern
func (p *RequestParser) ParsePathParams(path, pattern string) (map[string]string, error) {
	pathSegments := strings.Split(strings.Trim(path, "/"), "/")
	patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")

	if len(pathSegments) != len(patternSegments) {
		return nil, NewSDKError(ErrCodeRequestParsingFailed,
			fmt.Sprintf("Path segments mismatch: path=%s, pattern=%s", path, pattern))
	}

	pathParams := make(map[string]string)
	for i, patternSegment := range patternSegments {
		if strings.HasPrefix(patternSegment, "{") && strings.HasSuffix(patternSegment, "}") {
			paramName := strings.TrimSuffix(strings.TrimPrefix(patternSegment, "{"), "}")
			pathParams[paramName] = pathSegments[i]
		}
	}

	return pathParams, nil
}

// CreateSDKExecuteRequest creates a complete SDK execution request from HTTP request
func (p *RequestParser) CreateSDKExecuteRequest(ctx context.Context, req *http.Request, workflowID string) (*SDKExecuteWorkflowRequest, error) {
	// Parse HTTP request context
	httpContext, err := p.ParseHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	// Parse request body
	inputData, err := p.ParseRequestBody(req)
	if err != nil {
		return nil, err
	}

	// Extract session information
	var session *SessionContext
	if p.config.EnableSessionExtraction {
		session, err = p.ExtractSessionInfo(req)
		if err != nil {
			return nil, err
		}
	}

	// Extract security context
	var security *SecurityContext
	if p.config.EnableSecurityContext {
		security, err = p.ExtractSecurityContext(req)
		if err != nil {
			return nil, err
		}
	}

	// Create SDK execution request
	sdkRequest := NewSDKExecuteWorkflowRequest(inputData)
	sdkRequest.HTTPRequest = httpContext
	sdkRequest.Session = session
	sdkRequest.Security = security
	sdkRequest.EnableValidation = p.config.EnableValidation
	sdkRequest.EnableSanitization = p.config.EnableSanitization
	sdkRequest.IncludeFullContext = p.config.IncludeFullHTTPContext
	sdkRequest.ValidationRules = p.config.DefaultValidationRules

	// Add metadata
	if sdkRequest.Metadata == nil {
		sdkRequest.Metadata = make(map[string]interface{})
	}
	sdkRequest.Metadata["sdk_version"] = "1.0.0"
	sdkRequest.Metadata["parsed_at"] = time.Now().Format(time.RFC3339)
	sdkRequest.Metadata["request_id"] = getRequestIDFromContext(ctx)

	return sdkRequest, nil
}

// Helper functions

func extractSessionIDFromToken(token string) string {
	// In a real implementation, you would decode the JWT and extract the session ID
	// For now, return a hash of the token
	return fmt.Sprintf("session_%x", token)
}

func extractClaimsFromToken(token string) (map[string]interface{}, error) {
	// In a real implementation, you would decode and validate the JWT
	// For now, return empty claims
	return make(map[string]interface{}), nil
}

func getRequestIDFromContext(ctx context.Context) string {
	// Extract request ID from context
	// In a real implementation, you would use OpenTelemetry or similar
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
