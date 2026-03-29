# cmd/uwf-cli/

CLI tool (`uwf-cli`) built with Cobra for managing workflows from the command line.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `main.go` | CLI entry point, root Cobra command setup | Modifying CLI startup or root flags |
| `workflows.go` | Workflow management commands (list, register, get) | Adding workflow commands |
| `execute.go` | Workflow execution command | Modifying workflow execution via CLI |
| `executions.go` | Execution history/status commands | Adding execution query commands |
| `deploy.go` | Deployment commands | Modifying deploy workflow |
| `config.go` | CLI configuration management commands | Adding config commands |
| `health.go` | Health check command | Modifying health check behavior |
| `test.go` | Workflow test command | Modifying workflow testing via CLI |
| `completion.go` | Shell completion setup | Modifying shell completion |
