# Workflow API Documentation

A RESTful API for managing and executing workflows built with Go and Gin.

## Overview

This workflow API provides a complete system for:
- Registering and managing workflows
- Executing workflows asynchronously
- Tracking workflow execution status
- Managing workflow state and data

## Architecture

The system follows a modular architecture with clear separation of concerns:

1. **Registry** - Manages workflow definitions (in-memory implementation)
2. **Queue** - Handles workflow execution requests (in-memory and NATS JetStream implementations)
3. **State Management** - Tracks workflow execution state (in-memory implementation)
4. **Executor** - Executes workflows (simplified implementation)
5. **API Layer** - RESTful endpoints for external interaction

## API Endpoints

### Workflow Management

#### List Workflows
```
GET /api/v1/workflows
```

**Response:**
```json
{
  "workflows": [
    {
      "id": "workflow-1234567890",
      "name": "Data Processing",
      "description": "Process data through multiple steps",
      "step_count": 3
    }
  ],
  "count": 1
}
```

#### Get Workflow Details
```
GET /api/v1/workflows/{id}
```

**Response:**
```json
{
  "id": "workflow-1234567890",
  "name": "Data Processing",
  "description": "Process data through multiple steps",
  "step_count": 3,
  "steps": [
    {
      "name": "Data Extraction",
      "child_step_count": 1,
      "is_parallel": false
    },
    {
      "name": "Data Transformation",
      "child_step_count": 1,
      "is_parallel": false
    },
    {
      "name": "Data Loading",
      "child_step_count": 1,
      "is_parallel": false
    }
  ]
}
```

#### Create Workflow
```
POST /api/v1/workflows
```

**Request Body:**
```json
{
  "name": "New Workflow",
  "description": "Description of the new workflow"
}
```

**Response:**
```json
{
  "id": "workflow-1234567891",
  "name": "New Workflow",
  "description": "Description of the new workflow",
  "message": "Workflow created successfully"
}
```

#### Delete Workflow
```
DELETE /api/v1/workflows/{id}
```

**Response:**
```json
{
  "message": "Workflow deleted successfully"
}
```

### Workflow Execution

#### Execute Workflow
```
POST /api/v1/workflows/{id}/execute
```

**Response:**
```json
{
  "run_id": "run-1234567890",
  "message": "Workflow execution started",
  "status_url": "/api/v1/executions/run-1234567890"
}
```

#### List Executions
```
GET /api/v1/executions
```

**Query Parameters:**
- `workflow_id` (optional) - Filter by workflow ID
- `status` (optional) - Filter by status
- `limit` (optional, default: 50) - Number of results
- `offset` (optional, default: 0) - Pagination offset

**Response:**
```json
{
  "executions": [
    {
      "run_id": "run-1234567890",
      "workflow_id": "workflow-1234567890",
      "status": "pending",
      "current_step_index": 0,
      "current_child_step_index": 0,
      "start_time": "2024-01-13T16:09:00Z",
      "end_time": null,
      "error_message": null,
      "last_attempted_step": null,
      "is_terminal": false,
      "is_running": false,
      "is_pending": true,
      "created_at": "2024-01-13T16:09:00Z",
      "updated_at": "2024-01-13T16:09:00Z"
    }
  ],
  "count": 1,
  "limit": 50,
  "offset": 0
}
```

#### Get Execution Status
```
GET /api/v1/executions/{runId}
```

**Response:**
```json
{
  "run_id": "run-1234567890",
  "workflow_id": "workflow-1234567890",
  "status": "running",
  "current_step": "Data Extraction",
  "current_step_index": 0,
  "current_child_step_index": 0,
  "progress": 0.33,
  "start_time": "2024-01-13T16:09:00Z",
  "end_time": null,
  "error_message": null,
  "last_attempted_step": null,
  "is_terminal": false,
  "metadata": {}
}
```

#### Execution Control

- `POST /api/v1/executions/{runId}/cancel` - Cancel execution
- `POST /api/v1/executions/{runId}/pause` - Pause execution
- `POST /api/v1/executions/{runId}/resume` - Resume execution
- `POST /api/v1/executions/{runId}/retry` - Retry failed execution

#### Get Execution Data
```
GET /api/v1/executions/{runId}/data
```

**Response:**
```json
{
  "run_id": "run-1234567890",
  "data": {
    "input": {},
    "output": {},
    "intermediate": {}
  }
}
```

#### Get Execution Metrics
```
GET /api/v1/executions/{runId}/metrics
```

**Response:**
```json
{
  "run_id": "run-1234567890",
  "workflow_id": "workflow-1234567890",
  "workflow_metrics": {},
  "step_metrics": {},
  "child_step_metrics": {},
  "total_steps": 3,
  "completed_steps": 1,
  "failed_steps": 0,
  "total_child_steps": 3,
  "completed_child_steps": 1,
  "failed_child_steps": 0,
  "total_duration_millis": 1500,
  "average_step_duration": 500,
  "success_rate": 0.33
}
```

### Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": 1736798940
}
```

## Getting Started

### Prerequisites
- Go 1.21 or later
- (Optional) NATS server for queue persistence

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd unified-workflow-go
```

2. Install dependencies:
```bash
go mod download
```

3. Build the API server:
```bash
go build ./cmd/workflow-api
```

4. Run the server:
```bash
./workflow-api
```

### Configuration

The API uses a configuration file (`config.yaml`) with the following structure:

```yaml
server:
  port: 8080
  host: "localhost"

queue:
  type: "in-memory"  # or "nats"
  nats:
    urls: ["nats://localhost:4222"]
    stream_name: "workflow-queue"
    subject_prefix: "workflow"
    durable_name: "workflow-consumer"

executor:
  worker_count: 5
  queue_poll_interval: "1s"
  max_retries: 3
  retry_delay: "5s"
  execution_timeout: "30m"
  step_timeout: "5m"
  enable_metrics: true
  enable_tracing: false
  max_concurrent_workflows: 100

logging:
  level: "info"
  format: "json"
```

### Example Usage

#### Using cURL

1. Create a workflow:
```bash
curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Data Processing Pipeline",
    "description": "Extract, transform, and load data"
  }'
```

2. List workflows:
```bash
curl http://localhost:8080/api/v1/workflows
```

3. Execute a workflow:
```bash
curl -X POST http://localhost:8080/api/v1/workflows/{workflow-id}/execute
```

4. Check execution status:
```bash
curl http://localhost:8080/api/v1/executions/{run-id}
```

## Architecture Details

### Workflow Model

Workflows consist of:
- **Steps** - Individual units of work
- **Child Steps** - Sub-steps within a step (for parallel execution)
- **Primitives** - Reusable business logic components
- **Context** - Shared data between steps
- **Data** - Workflow-specific data

### Queue System

The API supports multiple queue backends:
- **In-memory queue** - For development and testing
- **NATS JetStream** - For production with persistence and scalability

### State Management

Execution state is tracked with:
- Current step and child step indices
- Execution status (pending, running, completed, failed, cancelled, paused)
- Error messages and retry information
- Timing metrics

### Extensibility

The system is designed to be extensible:
- Add new queue implementations by implementing the `Queue` interface
- Add new registry implementations by implementing the `Registry` interface
- Add new state management implementations by implementing the `StateManagement` interface
- Add new executor implementations by implementing the `Executor` interface

## Development

### Running Tests
```bash
go test ./...
```

### Code Structure
```
unified-workflow-go/
├── cmd/
│   └── workflow-api/          # Main application entry point
├── internal/
│   ├── api/
│   │   └── handlers/         # HTTP request handlers
│   ├── common/
│   │   └── model/            # Core data models
│   ├── config/               # Configuration management
│   ├── executor/             # Workflow execution logic
│   ├── primitive/            # Business logic primitives
│   ├── queue/                # Queue implementations
│   ├── registry/             # Workflow registry
│   └── state/                # State management
├── examples/                 # Example workflows and usage
└── pkg/                      # Public packages (if any)
```

### Adding New Features

1. **New Queue Backend**: Implement the `Queue` interface in a new file in `internal/queue/`
2. **New Registry Backend**: Implement the `Registry` interface in a new file in `internal/registry/`
3. **New State Backend**: Implement the `StateManagement` interface in a new file in `internal/state/`
4. **New API Endpoint**: Add handler method in `internal/api/handlers/` and register route in `cmd/workflow-api/main.go`

## License

[Your License Here]

## Support

For issues and feature requests, please create an issue in the repository.
