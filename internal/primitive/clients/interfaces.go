package clients

import (
	"context"
	"crypto/tls"
	"time"
)

// ClientType represents the type of client
type ClientType string

const (
	ClientTypeS3       ClientType = "s3"
	ClientTypeHTTP     ClientType = "http"
	ClientTypeGRPC     ClientType = "grpc"
	ClientTypeKafka    ClientType = "kafka"
	ClientTypeDatabase ClientType = "database"
)

// AuthType represents the authentication type
type AuthType string

const (
	AuthTypeNone   AuthType = "none"
	AuthTypeAPIKey AuthType = "api_key"
	AuthTypeJWT    AuthType = "jwt"
	AuthTypeOAuth2 AuthType = "oauth2"
	AuthTypeAWS    AuthType = "aws"
	AuthTypeBasic  AuthType = "basic"
)

// HealthStatus represents client health status
type HealthStatus struct {
	Healthy   bool          `json:"healthy"`
	Message   string        `json:"message,omitempty"`
	CheckedAt time.Time     `json:"checked_at"`
	Latency   time.Duration `json:"latency,omitempty"`
}

// ClientMetrics represents client performance metrics
type ClientMetrics struct {
	RequestsTotal     int64         `json:"requests_total"`
	RequestsFailed    int64         `json:"requests_failed"`
	AvgLatency        time.Duration `json:"avg_latency"`
	LastRequestTime   time.Time     `json:"last_request_time"`
	ConnectionCount   int           `json:"connection_count"`
	ActiveConnections int           `json:"active_connections"`
}

// AuthInfo represents authentication information
type AuthInfo struct {
	Type        AuthType   `json:"type"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Token       string     `json:"token,omitempty"`
	Scopes      []string   `json:"scopes,omitempty"`
	LastRefresh time.Time  `json:"last_refresh"`
}

// RetryPolicy defines retry behavior
type RetryPolicy struct {
	MaxAttempts int           `json:"max_attempts"`
	BaseDelay   time.Duration `json:"base_delay"`
	MaxDelay    time.Duration `json:"max_delay"`
	Jitter      bool          `json:"jitter"`
}

// ConnectionPoolConfig defines connection pool settings
type ConnectionPoolConfig struct {
	MaxIdle     int           `json:"max_idle"`
	MaxActive   int           `json:"max_active"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	Wait        bool          `json:"wait"`
}

// ClientConfig is the base configuration for all clients
type ClientConfig struct {
	Name        string        `json:"name" yaml:"name"`
	Type        ClientType    `json:"type" yaml:"type"`
	Endpoint    string        `json:"endpoint" yaml:"endpoint"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`
	RetryPolicy RetryPolicy   `json:"retry_policy" yaml:"retry_policy"`

	// Authentication
	AuthType    AuthType `json:"auth_type" yaml:"auth_type"`
	APIKey      string   `json:"api_key,omitempty" yaml:"api_key,omitempty"`
	APIKeyPath  string   `json:"api_key_path,omitempty" yaml:"api_key_path,omitempty"`
	JWTToken    string   `json:"jwt_token,omitempty" yaml:"jwt_token,omitempty"`
	JWTPath     string   `json:"jwt_path,omitempty" yaml:"jwt_path,omitempty"`
	OAuth2Token string   `json:"oauth2_token,omitempty" yaml:"oauth2_token,omitempty"`
	OAuth2Path  string   `json:"oauth2_path,omitempty" yaml:"oauth2_path,omitempty"`

	// Headers and metadata
	Headers  map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// TLS/SSL configuration
	TLSConfig     *tls.Config `json:"-" yaml:"-"`
	TLSCACert     string      `json:"tls_ca_cert,omitempty" yaml:"tls_ca_cert,omitempty"`
	TLSClientCert string      `json:"tls_client_cert,omitempty" yaml:"tls_client_cert,omitempty"`
	TLSClientKey  string      `json:"tls_client_key,omitempty" yaml:"tls_client_key,omitempty"`

	// Connection pooling
	PoolConfig ConnectionPoolConfig `json:"pool_config" yaml:"pool_config"`

	// Monitoring
	EnableMetrics bool `json:"enable_metrics" yaml:"enable_metrics"`
	EnableTracing bool `json:"enable_tracing" yaml:"enable_tracing"`
}

// BaseClient is the fundamental interface that all clients must implement
type BaseClient interface {
	// Connection Management
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	IsConnected() bool
	Ping(ctx context.Context) error

	// Authentication
	Authenticate(ctx context.Context) error
	GetAuthInfo() AuthInfo
	RefreshAuth(ctx context.Context) error

	// Configuration
	Configure(config ClientConfig) error
	GetConfig() ClientConfig
	UpdateConfig(config ClientConfig) error

	// Health and Status
	HealthCheck(ctx context.Context) (HealthStatus, error)
	GetMetrics() ClientMetrics
	GetLastError() error
	ResetMetrics()

	// Client Information
	GetClientType() ClientType
	GetClientName() string
	GetEndpoint() string

	// Lifecycle
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// S3Client defines operations for S3/Object Storage clients
type S3Client interface {
	BaseClient

	// Bucket operations
	CreateBucket(ctx context.Context, bucket string) error
	DeleteBucket(ctx context.Context, bucket string) error
	ListBuckets(ctx context.Context) ([]string, error)
	BucketExists(ctx context.Context, bucket string) (bool, error)

	// Object operations
	PutObject(ctx context.Context, bucket, key string, data []byte) error
	GetObject(ctx context.Context, bucket, key string) ([]byte, error)
	DeleteObject(ctx context.Context, bucket, key string) error
	ListObjects(ctx context.Context, bucket, prefix string) ([]ObjectInfo, error)
	ObjectExists(ctx context.Context, bucket, key string) (bool, error)

	// Multipart uploads
	CreateMultipartUpload(ctx context.Context, bucket, key string) (string, error)
	UploadPart(ctx context.Context, bucket, key, uploadID string, partNumber int, data []byte) (string, error)
	CompleteMultipartUpload(ctx context.Context, bucket, key, uploadID string, parts []PartInfo) error
	AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error

	// Presigned URLs
	GeneratePresignedURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error)
}

// ObjectInfo contains information about an S3 object
type ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	StorageClass string    `json:"storage_class,omitempty"`
}

// PartInfo contains information about a multipart upload part
type PartInfo struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
	Size       int64  `json:"size"`
}

// HTTPClient defines operations for HTTP/REST clients
type HTTPClient interface {
	BaseClient

	// HTTP methods
	Get(ctx context.Context, path string, headers map[string]string) (*HTTPResponse, error)
	Post(ctx context.Context, path string, body []byte, headers map[string]string) (*HTTPResponse, error)
	Put(ctx context.Context, path string, body []byte, headers map[string]string) (*HTTPResponse, error)
	Delete(ctx context.Context, path string, headers map[string]string) (*HTTPResponse, error)
	Patch(ctx context.Context, path string, body []byte, headers map[string]string) (*HTTPResponse, error)
	Head(ctx context.Context, path string, headers map[string]string) (*HTTPResponse, error)

	// Request building
	NewRequest(method, path string) *HTTPRequest
	DoRequest(req *HTTPRequest) (*HTTPResponse, error)

	// Batch operations
	BatchRequest(ctx context.Context, requests []*HTTPRequest) ([]*HTTPResponse, error)
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	Latency    time.Duration
	Request    *HTTPRequest
}

// ClientFactory creates clients based on configuration
type ClientFactory interface {
	CreateClient(config ClientConfig) (BaseClient, error)
	CreateS3Client(config S3ClientConfig) (S3Client, error)
	CreateHTTPClient(config HTTPClientConfig) (HTTPClient, error)
}

// S3ClientConfig extends ClientConfig with S3-specific settings
type S3ClientConfig struct {
	ClientConfig

	Region       string `json:"region" yaml:"region"`
	Bucket       string `json:"bucket,omitempty" yaml:"bucket,omitempty"`
	AccessKey    string `json:"access_key,omitempty" yaml:"access_key,omitempty"`
	SecretKey    string `json:"secret_key,omitempty" yaml:"secret_key,omitempty"`
	SessionToken string `json:"session_token,omitempty" yaml:"session_token,omitempty"`
	UsePathStyle bool   `json:"use_path_style" yaml:"use_path_style"`
	DisableSSL   bool   `json:"disable_ssl" yaml:"disable_ssl"`

	// Vault paths for secret resolution
	AccessKeyPath    string `json:"access_key_path,omitempty" yaml:"access_key_path,omitempty"`
	SecretKeyPath    string `json:"secret_key_path,omitempty" yaml:"secret_key_path,omitempty"`
	SessionTokenPath string `json:"session_token_path,omitempty" yaml:"session_token_path,omitempty"`
}

// HTTPClientConfig extends ClientConfig with HTTP-specific settings
type HTTPClientConfig struct {
	ClientConfig

	BaseURL         string            `json:"base_url" yaml:"base_url"`
	DefaultHeaders  map[string]string `json:"default_headers" yaml:"default_headers"`
	FollowRedirects bool              `json:"follow_redirects" yaml:"follow_redirects"`
	MaxRedirects    int               `json:"max_redirects" yaml:"max_redirects"`
	Compression     bool              `json:"compression" yaml:"compression"`

	// Rate limiting
	RateLimit  int `json:"rate_limit" yaml:"rate_limit"` // requests per second
	BurstLimit int `json:"burst_limit" yaml:"burst_limit"`

	// Circuit breaker
	CircuitBreakerEnabled bool          `json:"circuit_breaker_enabled" yaml:"circuit_breaker_enabled"`
	FailureThreshold      int           `json:"failure_threshold" yaml:"failure_threshold"`
	SuccessThreshold      int           `json:"success_threshold" yaml:"success_threshold"`
	TimeoutWindow         time.Duration `json:"timeout_window" yaml:"timeout_window"`
}
