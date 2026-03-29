# internal/executor/

Workflow and step execution engines.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `executor.go` | Executor interface and core types | Understanding executor contract, adding executor implementations |
| `workflow_executor.go` | Full workflow execution engine: step sequencing, compensation, state management | Modifying workflow execution logic, adding parallel execution |
| `simple_executor.go` | Simplified executor for single-step or test use | Understanding basic execution, testing executor logic |
| `workflow_executor.go.bak` | Backup of previous executor implementation | Reviewing previous executor design — do not use |
