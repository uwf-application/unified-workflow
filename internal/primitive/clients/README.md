# Client Provider System

A comprehensive client management system for unified-workflow that provides a unified interface for various client types (S3, HTTP, etc.) with authentication, monitoring, and lifecycle management.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Provider                          │
│  ┌─────────────┐  ┌─────────────┐  ┌───────────────────┐  │
│  │ S3 Clients  │  │ HTTP Clients│  │ Other Clients     │  │
│  └─────────────┘  └─────────────┘  └───────────────────┘  │
│         │               │                       │          │
│  ┌─────────────────────────────────────────────────────┐  │
│  │              Authentication Strategies              │  │
│  │  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  │  │
│  │  │ JWT  │  │ API  │  │OAuth2│  │ AWS  │  │Basic │  │  │
│  │  │      │  │ Key  │  │      │  │      │  │      │  │  │
│  │  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘  │  │
│  └─────────────────────────────────────────────────────┘  │
│         │               │                       │          │
│  ┌─────────────────────────────────────────────────────┐  │
│  │              Configuration Management               │  │
│  │  ┌───────────────────────────────────────────────┐  │  │
│  │  │          Hashicorp Vault Integration          │  │  │
│  │  └───────────────────────────────────────────────┘  │  │
│  └─────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Base Client Interface (`clients/interfaces.go`)

The foundation of the system with interfaces for:
- `BaseClient`: Core client operations (connect, authenticate, health check, etc.)
- `S3Client`: S3-specific operations (put/get objects, list buckets, etc.)
- `HTTPClient`: HTTP/REST operations (GET, POST, PUT, DELETE, etc.)

### 2. Client Implementations

- **BaseClientImpl** (`clients/base_client.go`): Base implementation with metrics, error handling, and lifecycle management
- **S3ClientImpl** (`clients/s3_client.go`): S3 client implementation (mock/real AWS SDK)
- **HTTPClientImpl** (`clients/http_client.go`): HTTP client with authentication support

### 3. Authentication Strategies (`auth/interfaces.go`)

Support for multiple authentication types:
- JWT (JSON Web Tokens)
- API Key
- OAuth2
- AWS Credentials
- Basic Authentication
- Hashicorp Vault integration for secret management

### 4. Client Provider (`providers/`)

- **ClientProvider Interface**: Manages client lifecycle, pooling, and configuration
- **DefaultClientProvider**: Simple implementation for client management
- **Vault integration**: For secret resolution and secure credential management

## Usage Examples

### Creating and Using an S3 Client

```go
package main

import (
    "context"
    "time"
    
    "unified-workflow/internal/primitive/clients"
    "unified-workflow/internal/primitive/providers"
)

func main() {
    ctx := context.Background()
    provider := providers.NewDefaultClientProvider()
    
    // Configure S3 client
    s3Config := clients.S3ClientConfig{
        ClientConfig: clients.ClientConfig{
            Name:     "my-s3-client",
            Type:     clients.ClientTypeS3,
            Endpoint: "https://s3.amazonaws.com",
            Timeout:  30 * time.Second,
            AuthType: clients.AuthTypeAWS,
        },
        Region: "us-east-1",
        Bucket: "my-bucket",
    }
    
    // Create S3 client
    s3Client, err := provider.CreateS3Client(ctx, s3Config)
    if err != nil {
        panic(err)
    }
    defer s3Client.Shutdown(ctx)
    
    // Use the client
    buckets, err := s3Client.ListBuckets(ctx)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Buckets: %v\n", buckets)
}
```

### Creating and Using an HTTP Client

```go
// Configure HTTP client
httpConfig := clients.HTTPClientConfig{
    ClientConfig: clients.ClientConfig{
        Name:     "my-http-client",
        Type:     clients.ClientTypeHTTP,
        Timeout:  10 * time.Second,
        AuthType: clients.AuthTypeAPIKey,
        APIKey:   "my-api-key-12345",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    },
    BaseURL: "https://api.example.com",
}

// Create HTTP client
httpClient, err := provider.CreateHTTPClient(ctx, httpConfig)
if err != nil {
    panic(err)
}
defer httpClient.Shutdown(ctx)

// Make requests
resp, err := httpClient.Get(ctx, "/users", nil)
if err != nil {
    panic(err)
}

fmt.Printf("Response: %s\n", string(resp.Body))
```

### Authentication Examples

```go
// JWT Authentication
jwtConfig := clients.ClientConfig{
    AuthType: clients.AuthTypeJWT,
    JWTToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
}

// API Key Authentication
apiKeyConfig := clients.ClientConfig{
    AuthType: clients.AuthTypeAPIKey,
    APIKey:   "my-api-key",
    Headers: map[string]string{
        "X-API-Key": "my-api-key",
    },
}

// OAuth2 Authentication
oauth2Config := clients.ClientConfig{
    AuthType:    clients.AuthTypeOAuth2,
    OAuth2Token: "oauth2-token-here",
}
```

## Configuration

### Client Configuration Structure

```go
type ClientConfig struct {
    Name        string            `json:"name" yaml:"name"`
    Type        ClientType        `json:"type" yaml:"type"`
    Endpoint    string            `json:"endpoint" yaml:"endpoint"`
    Timeout     time.Duration     `json:"timeout" yaml:"timeout"`
    AuthType    AuthType          `json:"auth_type" yaml:"auth_type"`
    
    // Authentication fields
    APIKey      string            `json:"api_key,omitempty" yaml:"api_key,omitempty"`
    APIKeyPath  string            `json:"api_key_path,omitempty" yaml:"api_key_path,omitempty"`
    JWTToken    string            `json:"jwt_token,omitempty" yaml:"jwt_token,omitempty"`
    JWTPath     string            `json:"jwt_path,omitempty" yaml:"jwt_path,omitempty"`
    
    // Headers and metadata
    Headers     map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
    
    // Monitoring
    EnableMetrics bool            `json:"enable_metrics" yaml:"enable_metrics"`
}
```

### S3-Specific Configuration

```go
type S3ClientConfig struct {
    ClientConfig
    
    Region       string `json:"region" yaml:"region"`
    Bucket       string `json:"bucket,omitempty" yaml:"bucket,omitempty"`
    AccessKey    string `json:"access_key,omitempty" yaml:"access_key,omitempty"`
    SecretKey    string `json:"secret_key,omitempty" yaml:"secret_key,omitempty"`
    
    // Vault paths for secret resolution
    AccessKeyPath string `json:"access_key_path,omitempty" yaml:"access_key_path,omitempty"`
    SecretKeyPath string `json:"secret_key_path,omitempty" yaml:"secret_key_path,omitempty"`
}
```

## Monitoring and Metrics

The system provides built-in metrics collection:

```go
// Get client metrics
metrics := client.GetMetrics()
fmt.Printf("Total Requests: %d\n", metrics.RequestsTotal)
fmt.Printf("Failed Requests: %d\n", metrics.RequestsFailed)
fmt.Printf("Average Latency: %v\n", metrics.AvgLatency)

// Health checking
health, err := client.HealthCheck(ctx)
if err != nil {
    panic(err)
}
fmt.Printf("Healthy: %v, Message: %s\n", health.Healthy, health.Message)

// Provider-level metrics
allMetrics := provider.GetAllMetrics()
allHealth := provider.HealthCheckAll(ctx)
```

## Integration with Primitive Wrappers

The client provider system can be integrated with existing primitive wrappers:

```go
type ClientAwarePrimitiveWrapper struct {
    service        interface{}
    clientName     string
    clientType     string
    clientProvider providers.ClientProvider
}

func (w *ClientAwarePrimitiveWrapper) Execute(ctx model.WorkflowContext, data model.WorkflowData) (interface{}, error) {
    // Get client from provider
    client, err := w.clientProvider.GetClient(context.Background(), w.clientName)
    if err != nil {
        return nil, fmt.Errorf("failed to get client %s: %w", w.clientName, err)
    }
    
    // Execute operation with client
    switch w.clientType {
    case "s3":
        return w.executeS3Operation(client.(clients.S3Client), ctx, data)
    case "http":
        return w.executeHTTPOperation(client.(clients.HTTPClient), ctx, data)
    default:
        return nil, fmt.Errorf("unsupported client type: %s", w.clientType)
    }
}
```

## Security Features

1. **Secret Management**: Integration with Hashicorp Vault for secure credential storage
2. **Token Rotation**: Automatic refresh of JWT and OAuth2 tokens
3. **TLS/SSL**: Support for secure connections with certificate validation
4. **Audit Logging**: Comprehensive logging of client operations
5. **Access Control**: Vault policies for fine-grained secret access control

## Extending the System

### Adding a New Client Type

1. Implement the `BaseClient` interface
2. Add client type to `ClientType` enum
3. Create configuration structure
4. Register with client factory

### Adding a New Authentication Strategy

1. Implement the `AuthStrategy` interface
2. Add authentication type to `AuthType` enum
3. Create configuration structure
4. Integrate with client implementations

## Running the Example

See `examples/client_provider_example.go` for a complete working example.

```bash
cd unified-workflow-go
go run examples/client_provider_example.go
```

## Future Enhancements

1. **Real AWS SDK Integration**: Replace mock S3 implementation with real AWS SDK
2. **Connection Pooling**: Advanced connection pooling with load balancing
3. **Circuit Breaker**: Implement circuit breaker pattern for fault tolerance
4. **Rate Limiting**: Client-side rate limiting
5. **Distributed Tracing**: Integration with OpenTelemetry for distributed tracing
6. **Configuration Hot Reload**: Dynamic configuration updates without restart
7. **Multi-cloud Support**: Support for other cloud providers (GCP, Azure, etc.)
