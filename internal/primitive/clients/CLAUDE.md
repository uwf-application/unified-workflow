# internal/primitive/clients/

Base HTTP and storage clients shared across primitive service implementations.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `interfaces.go` | Client interface contracts (HTTP, S3) | Adding new client types, understanding client abstractions |
| `base_client.go` | Base HTTP client with retries, timeouts, and auth headers | Modifying shared HTTP behavior for all primitives |
| `http_client.go` | Concrete HTTP client implementation | Debugging HTTP calls, modifying connection pooling |
| `s3_client.go` | S3 storage client for primitives that need object storage | Adding S3-backed primitives |
