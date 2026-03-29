# workflows/

Workflow definitions — each file defines a complete workflow with its step sequence.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `antifraud_workflow.go` | Antifraud workflow: FC validation → ML validation → AML → finalize transaction | Modifying antifraud workflow step sequence or configuration |
| `example_workflow.go` | Example workflow demonstrating basic step composition patterns | Understanding how to define a new workflow |

## Subdirectories

| Directory | What | When to read |
| --------- | ---- | ------------ |
| `steps/` | Individual step implementations used in workflows | Adding steps, modifying step logic |
| `child_steps/` | Child step implementations for nested workflow execution | Adding child steps, modifying sub-workflow behavior |
