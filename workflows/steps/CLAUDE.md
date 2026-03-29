# workflows/steps/

Individual workflow step implementations.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `antifraud_step.go` | Main antifraud step: calls antifraud primitive, handles fraud decision | Modifying antifraud check logic |
| `fc_validation_step.go` | Financial compliance (FC) validation step | Modifying compliance validation behavior |
| `ml_validation_step.go` | ML model validation step (risk scoring) | Modifying ML risk scoring integration |
| `aml_validation_step.go` | Anti-money laundering (AML) validation step | Modifying AML check behavior |
| `finalize_transaction_step.go` | Transaction finalization step (commit or reject) | Modifying transaction outcome handling |
| `store_transaction_step.go` | Transaction storage step (persist to database) | Modifying transaction persistence logic |
| `echo_step.go` | Echo/passthrough step for testing and debugging | Adding test steps, debugging workflows |
