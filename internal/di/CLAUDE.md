# internal/di/

Dependency injection container for wiring services and their dependencies.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `container.go` | DI container: service registration, resolution, lifecycle management | Adding services, modifying dependency wiring |
| `highperf_container.go` | High-performance DI container variant for production use | Optimizing DI performance, understanding production wiring |
| `di_factory.go` | Factory functions for constructing wired service graphs | Adding factory functions, modifying service construction |
| `circuit_breaker.go` | Circuit breaker integration within the DI container | Configuring circuit breakers for external services |
| `cluster.go` | Cluster-aware DI configuration | Modifying cluster topology, multi-node setup |
| `metrics.go` | Metrics/Prometheus integration within DI | Adding metrics to services |
| `primitive_integration.go` | DI wiring for primitive (external service) integrations | Adding new primitive services to the DI graph |
| `tracing.go` | Distributed tracing setup within DI | Configuring tracing providers |
| `example_test.go` | DI container usage examples | Understanding how to use the DI container |
