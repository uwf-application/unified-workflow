# internal/registry/

Workflow registry interface and implementations for storing and looking up workflow definitions.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `registry.go` | Registry interface: Register, Get, List workflow definitions | Understanding registry contract, adding registry backends |
| `http_registry.go` | HTTP-based registry client (calls external registry API) | Modifying remote registry communication |
| `in_memory_registry.go` | In-memory registry for testing and local development | Modifying test registry behavior |
