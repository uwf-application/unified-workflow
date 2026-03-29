# internal/api/

HTTP API layer for the workflow service.

## Subdirectories

| Directory | What | When to read |
| --------- | ---- | ------------ |
| `handlers/` | Gin HTTP route handlers | Adding endpoints, modifying request/response logic |
| `middleware/` | Gin middleware (auth, logging, etc.) | Adding middleware, modifying request pipeline |
