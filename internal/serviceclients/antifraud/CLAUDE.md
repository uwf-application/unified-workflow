# internal/serviceclients/antifraud/

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `client.go` | Antifraud service HTTP client: request construction, response parsing | Modifying antifraud API calls, adding endpoints |
| `client_test.go` | Unit tests for the antifraud HTTP client | Running antifraud client tests, understanding client behavior |
| `proxy.go` | Proxy wrapper for antifraud client (circuit breaker, retry) | Modifying resilience behavior around antifraud calls |
