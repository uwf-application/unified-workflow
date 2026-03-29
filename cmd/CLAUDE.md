# cmd/

Application entry points — one `main` package per deployable service.

## Subdirectories

| Directory | What | When to read |
| --------- | ---- | ------------ |
| `workflow-api/` | Workflow REST API server (Gin) entry point | Modifying API startup, middleware, or port configuration |
| `workflow-worker/` | Workflow execution worker entry point | Modifying worker startup or consumer configuration |
| `uwf-cli/` | CLI tool entry point with Cobra commands | Adding CLI commands, modifying CLI behavior |
| `executor-api/` | Executor service entry point | Modifying executor startup |
| `primitive-api/` | Primitive operation API entry point | Modifying primitive API startup |
| `registry-api/` | Workflow registry API entry point | Modifying registry startup |
