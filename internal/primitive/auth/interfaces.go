package auth

import (
	"context"
	"time"

	"unified-workflow/internal/primitive/clients"
)

// AuthStrategy defines the interface for authentication strategies
type AuthStrategy interface {
	// Authenticate performs authentication
	Authenticate(ctx context.Context) error

	// Refresh refreshes authentication tokens
	Refresh(ctx context.Context) error

	// GetAuthHeaders returns headers to include in requests
	GetAuthHeaders() map[string]string

	// GetAuthInfo returns authentication information
	GetAuthInfo() clients.AuthInfo

	// GetType returns the authentication type
	GetType() clients.AuthType

	// IsValid checks if authentication is still valid
	IsValid() bool

	// GetExpiry returns when authentication expires
	GetExpiry() *time.Time
}

// CredentialProvider provides credentials for authentication
type CredentialProvider interface {
	// GetCredentials retrieves credentials
	GetCredentials(ctx context.Context) (interface{}, error)

	// RefreshCredentials refreshes credentials
	RefreshCredentials(ctx context.Context) error

	// GetCredentialType returns the type of credentials
	GetCredentialType() string
}

// JWTConfig configuration for JWT authentication
type JWTConfig struct {
	Token         string                 `json:"token" yaml:"token"`
	PrivateKey    []byte                 `json:"private_key,omitempty" yaml:"private_key,omitempty"`
	PublicKey     []byte                 `json:"public_key,omitempty" yaml:"public_key,omitempty"`
	KeyPath       string                 `json:"key_path,omitempty" yaml:"key_path,omitempty"`
	Claims        map[string]interface{} `json:"claims,omitempty" yaml:"claims,omitempty"`
	Expiration    time.Duration          `json:"expiration" yaml:"expiration"`
	RefreshBefore time.Duration          `json:"refresh_before" yaml:"refresh_before"`
	Issuer        string                 `json:"issuer,omitempty" yaml:"issuer,omitempty"`
	Audience      []string               `json:"audience,omitempty" yaml:"audience,omitempty"`
	Subject       string                 `json:"subject,omitempty" yaml:"subject,omitempty"`
}

// APIKeyConfig configuration for API key authentication
type APIKeyConfig struct {
	APIKey     string `json:"api_key" yaml:"api_key"`
	HeaderName string `json:"header_name" yaml:"header_name"`                       // e.g., "X-API-Key"
	QueryParam string `json:"query_param,omitempty" yaml:"query_param,omitempty"`   // e.g., "api_key"
	AuthScheme string `json:"auth_scheme,omitempty" yaml:"auth_scheme,omitempty"`   // e.g., "Bearer", "Token"
	APIKeyPath string `json:"api_key_path,omitempty" yaml:"api_key_path,omitempty"` // Vault path
}

// OAuth2Config configuration for OAuth2 authentication
type OAuth2Config struct {
	TokenURL       string            `json:"token_url" yaml:"token_url"`
	ClientID       string            `json:"client_id" yaml:"client_id"`
	ClientSecret   string            `json:"client_secret" yaml:"client_secret"`
	Scopes         []string          `json:"scopes" yaml:"scopes"`
	AuthStyle      int               `json:"auth_style" yaml:"auth_style"` // 0 = auto, 1 = params, 2 = header
	EndpointParams map[string]string `json:"endpoint_params,omitempty" yaml:"endpoint_params,omitempty"`
	RedirectURL    string            `json:"redirect_url,omitempty" yaml:"redirect_url,omitempty"`

	// Vault paths for secret resolution
	ClientIDPath     string `json:"client_id_path,omitempty" yaml:"client_id_path,omitempty"`
	ClientSecretPath string `json:"client_secret_path,omitempty" yaml:"client_secret_path,omitempty"`
}

// AWSConfig configuration for AWS authentication
type AWSConfig struct {
	AccessKey    string `json:"access_key,omitempty" yaml:"access_key,omitempty"`
	SecretKey    string `json:"secret_key,omitempty" yaml:"secret_key,omitempty"`
	SessionToken string `json:"session_token,omitempty" yaml:"session_token,omitempty"`
	Region       string `json:"region" yaml:"region"`
	Profile      string `json:"profile,omitempty" yaml:"profile,omitempty"`

	// Vault paths for secret resolution
	AccessKeyPath    string `json:"access_key_path,omitempty" yaml:"access_key_path,omitempty"`
	SecretKeyPath    string `json:"secret_key_path,omitempty" yaml:"secret_key_path,omitempty"`
	SessionTokenPath string `json:"session_token_path,omitempty" yaml:"session_token_path,omitempty"`
}

// BasicAuthConfig configuration for basic authentication
type BasicAuthConfig struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	// Vault paths for secret resolution
	UsernamePath string `json:"username_path,omitempty" yaml:"username_path,omitempty"`
	PasswordPath string `json:"password_path,omitempty" yaml:"password_path,omitempty"`
}

// VaultConfig configuration for Hashicorp Vault
type VaultConfig struct {
	Address     string `json:"address" yaml:"address"`
	Token       string `json:"token,omitempty" yaml:"token,omitempty"`
	TokenPath   string `json:"token_path,omitempty" yaml:"token_path,omitempty"` // File path for token
	RoleID      string `json:"role_id,omitempty" yaml:"role_id,omitempty"`
	SecretID    string `json:"secret_id,omitempty" yaml:"secret_id,omitempty"`
	AppRolePath string `json:"approle_path" yaml:"approle_path"` // e.g., "auth/approle/login"
	SecretsPath string `json:"secrets_path" yaml:"secrets_path"` // e.g., "secret/data"
	Namespace   string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	MountPath   string `json:"mount_path,omitempty" yaml:"mount_path,omitempty"`
	Engine      string `json:"engine" yaml:"engine"` // "kv", "transit", "aws", etc.
}

// TokenCache caches authentication tokens
type TokenCache interface {
	// Get retrieves a token from cache
	Get(key string) (interface{}, bool)

	// Set stores a token in cache
	Set(key string, value interface{}, ttl time.Duration)

	// Delete removes a token from cache
	Delete(key string)

	// Clear clears all tokens from cache
	Clear()

	// Size returns the number of items in cache
	Size() int
}

// AuthManager manages multiple authentication strategies
type AuthManager interface {
	// RegisterStrategy registers an authentication strategy
	RegisterStrategy(name string, strategy AuthStrategy) error

	// GetStrategy retrieves an authentication strategy
	GetStrategy(name string) (AuthStrategy, error)

	// Authenticate performs authentication for a strategy
	Authenticate(ctx context.Context, name string) error

	// Refresh refreshes authentication for a strategy
	Refresh(ctx context.Context, name string) error

	// GetAuthHeaders returns headers for a strategy
	GetAuthHeaders(name string) (map[string]string, error)

	// RemoveStrategy removes an authentication strategy
	RemoveStrategy(name string) error

	// ListStrategies lists all registered strategies
	ListStrategies() []string
}
