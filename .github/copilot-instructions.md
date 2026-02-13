# Copilot Instructions for Unified Workflow System (Go)

## Overview
This document provides guidance for AI coding agents working with the Unified Workflow System codebase. The system is a high-performance, loosely coupled workflow orchestration system built in Go, utilizing Gin and NATS JetStream.

## Architecture
- **Microservices**: The system is composed of multiple microservices that communicate via NATS JetStream, allowing for independent deployment and scaling.
- **Key Components**:
  - **Workflow API**: RESTful API for managing workflows.
  - **Workflow Registry**: Manages workflow definitions and versions.
  - **Workflow Engine**: Executes workflows and manages state.
  - **NATS JetStream**: Handles event-driven communication between components.

## Developer Workflows
- **Building the Project**: Use `go build -o workflow-api ./cmd/workflow-api` to build the API server.
- **Running Tests**: Execute `go test ./...` to run all tests in the codebase.
- **Debugging**: Use `dlv debug ./cmd/workflow-api` for debugging the API server.

## Project Conventions
- **File Structure**: Follow the established directory structure for organizing code, with `cmd/` for entry points and `internal/` for private application code.
- **Naming Conventions**: Use camelCase for function names and PascalCase for types.

## Integration Points
- **NATS JetStream**: Ensure the NATS server is running with JetStream enabled for message delivery.
- **Database Connections**: Configure PostgreSQL for production environments and Redis for caching.

## Examples
- **Creating a Workflow**: Use the `POST /api/v1/workflows` endpoint to create a new workflow instance.
- **Executing a Workflow**: Trigger execution with `POST /api/v1/workflows/:id/execute`.

## Conclusion
This document should serve as a foundational guide for AI agents to navigate and contribute effectively to the Unified Workflow System codebase. For further details, refer to the README.md and other documentation within the repository.