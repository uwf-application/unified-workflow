# internal/primitive/

External service primitive layer — wraps third-party services behind workflow-friendly interfaces.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `interfaces.go` | Top-level primitive interfaces for use in workflows | Understanding primitive contracts, adding workflow steps |
| `service_interfaces.go` | Service-level interface definitions for primitive implementations | Adding new primitive service types |
| `global.go` | Global primitive registry and accessor functions | Understanding primitive registration, debugging global state |
| `proxies.go` | Proxy implementations wrapping service clients with resilience patterns | Modifying circuit breaking/retry around primitives |

## Subdirectories

| Directory | What | When to read |
| --------- | ---- | ------------ |
| `auth/` | Authentication primitives (token validation, JWT) | Adding auth to primitives, modifying auth behavior |
| `clients/` | Base HTTP client and S3 client used by primitive services | Modifying underlying HTTP behavior for primitives |
| `config/` | Primitive service configuration types | Adding config for new primitives |
| `interfaces/` | Detailed interface definitions for each primitive category | Adding new primitive operations |
| `model/` | Primitive domain models (workflow context, workflow data) | Modifying primitive data shapes |
| `services/` | Concrete primitive service implementations | Adding a new primitive service |
