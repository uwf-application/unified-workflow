# Unified Workflow System (Go)

High-performance workflow orchestration system for banking environments using Gin and NATS JetStream.

## Files

| File | What | When to read |
| ---- | ---- | ------------ |
| `README.md` | Project overview, architecture diagram, features | Starting development, understanding the system |
| `go.mod` | Module name (`unified-workflow`), Go version, all dependencies | Adding dependencies, resolving version conflicts |
| `go.sum` | Dependency checksums | Never edit directly |
| `Makefile` | Build, test, run, and Docker targets | Running builds, tests, or services locally |
| `config.yaml` | Default runtime configuration (NATS, Redis, ports) | Configuring local dev environment, understanding defaults |
| `deploy-config.yaml` | Production deployment configuration | Deploying services, configuring production |
| `docker-compose.yml` | Local development stack (NATS, Redis, all services) | Starting local dev environment |
| `docker-compose.yml.backup` | Previous docker-compose snapshot | Reverting docker-compose changes |
| `VERSION` | Current release version string | Versioning builds, release automation |
| `Dockerfile.workflow-api` | Workflow API service container image | Building/deploying the API service |
| `Dockerfile.worker` | Workflow worker container image | Building/deploying worker service |
| `Dockerfile.executor` | Executor service container image | Building/deploying executor service |
| `Dockerfile.registry` | Registry service container image | Building/deploying registry service |
| `build_binaries.sh` | Cross-platform binary build script | Building release binaries |
| `create_release.sh` | Release creation automation script | Cutting a new release |
| `package_sdk.sh` | SDK packaging script | Publishing client SDKs |
| `register_antifraud_workflow.sh` | Script to register antifraud workflow via API | Setting up antifraud workflow in an environment |
| `test_antifraud_workflow.sh` | End-to-end antifraud workflow test script | Testing antifraud workflow integration |
| `postman_test_commands.sh` | Postman/curl API test commands | Manual API testing |
| `execute_antifraud_workflow.go` | Standalone antifraud workflow execution example | Understanding workflow triggering |
| `API_README.md` | REST API endpoint reference | Integrating with the workflow API |
| `CLIENT_SDK_GUIDE.md` | Client SDK usage guide (Go, TypeScript, Python) | Integrating client SDKs |
| `ANTIFRAUD_PRIMITIVE_INTEGRATION.md` | Antifraud primitive integration guide | Integrating antifraud services as primitives |
| `DI_FRAMEWORK_SUMMARY.md` | Dependency injection framework summary | Understanding the DI container |
| `CLI_INSTALLATION.md` | CLI tool installation instructions | Installing the uwf-cli tool |
| `DEPLOYMENT_CLI.md` | CLI-based deployment guide | Deploying via CLI |
| `HARDWARE_REQUIREMENTS.md` | Infrastructure/hardware sizing requirements | Planning deployment capacity |
| `HLD_1000_TPS_DEPLOYMENT_PLAN.md` | High-level design for 1000 TPS deployment | Scaling to production load |
| `POC_DEPLOYMENT_PLAN_200_TPS.md` | POC deployment plan for 200 TPS | Setting up a proof-of-concept environment |
| `COMPONENT_RELEASE_STRATEGY.md` | Per-component release strategy | Planning component-level releases |
| `RELEASE_CLI.md` | CLI release notes | Reviewing CLI changes |
| `RELEASE_PLAN.md` | Overall release plan | Planning releases |
| `RELEASE_v1.2.0.md` | v1.2.0 release notes | Reviewing v1.2.0 changes |

## Subdirectories

| Directory | What | When to read |
| --------- | ---- | ------------ |
| `api/` | Smithy API definition files | Modifying API contracts, understanding API shape |
| `cmd/` | Application entry points (main packages) | Adding a new service, modifying startup behavior |
| `internal/` | Private application packages (not importable externally) | Implementing features, debugging internal logic |
| `workflows/` | Workflow definitions and step implementations | Adding workflows, modifying workflow logic |
| `pkg/client/` | Public client libraries (Go, TypeScript, Python) | Using or modifying client SDKs |
| `examples/` | Runnable usage examples | Understanding how to use the system |
| `scripts/` | Build and deployment automation scripts | Modifying CI/CD, build automation |

## Build

```bash
make build          # Build all services
make build-cli      # Build uwf-cli binary
./build_binaries.sh # Cross-compile release binaries
```

## Test

```bash
make test           # Run all tests
make test-coverage  # Run tests with coverage report
```

## Development

```bash
docker-compose up -d    # Start NATS + Redis dependencies
make run-api            # Run workflow API locally
make run-worker         # Run workflow worker locally
```
