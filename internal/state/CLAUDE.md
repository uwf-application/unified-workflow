# internal/state/

Workflow execution state storage interface and implementations.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `state.go` | State interface: Get, Set, Delete execution state | Understanding state contract, adding state backends |
| `redis_state.go` | Redis-backed state implementation | Modifying state persistence, debugging Redis state issues |
| `in_memory_state.go` | In-memory state implementation for testing | Modifying test state behavior |
