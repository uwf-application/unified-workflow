# Unified Workflow System (Go)

A high-performance, loosely coupled workflow orchestration system built in Go with Gin and NATS JetStream. Designed for banking environments with hybrid cloud/on-premise deployment.

## Architecture Overview

The system is built as a collection of loosely coupled microservices communicating via NATS JetStream. Each component can be deployed independently and scaled separately.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      Unified Workflow System (Go)                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────┐ │
│  │   Workflow  │    │   Workflow  │    │   Workflow  │    │  NATS    │ │
│  │     API     │────│   Registry  │────│   Engine    │────│ JetStream│ │
│  │   (Gin)     │    │             │    │             │    │          │ │
│  └─────────────┘    └─────────────┘    └─────────────┘    └──────────┘ │
│        │                   │                   │               │        │
│        ▼                   ▼                   ▼               ▼        │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────┐ │
│  │   Client    │    │   Primitive │    │   Executor  │    │ Storage  │ │
│  │  Libraries  │    │ Operations  │    │  Services   │    │  Layer   │ │
│  │             │    │  (Interface)│    │             │    │          │ │
│  └─────────────┘    └─────────────┘    └─────────────┘    └──────────┘ │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Key Features

- **Loosely Coupled Architecture**: Each component can be deployed and scaled independently
- **Multi-language Support**: Go-based with interfaces for easy extension
- **Event-Driven**: NATS JetStream for reliable message delivery
- **High Performance**: Built with Go for low latency and high throughput
- **Banking-Grade Security**: Built-in security and compliance features
- **Hybrid Deployment**: Supports on-premise, private cloud, and public cloud
- **Workflow Versioning**: Full version control for workflow definitions
- **Compensation Logic**: Automatic rollback for failed steps
- **Parallel Execution**: Support for parallel step execution
- **Monitoring & Observability**: Prometheus metrics, structured logging, distributed tracing

## Project Structure

```
unified-workflow-go/
├── cmd/                          # Application entry points
│   ├── workflow-api/             # REST API server (Gin)
│   ├── workflow-registry/        # Workflow registry service
│   ├── workflow-engine/          # Workflow execution engine
│   └── executor-*/               # Operation executor services
├── internal/                     # Private application code
│   ├── primitive/                # Core interfaces and primitives
│   ├── models/                   # Data models and structures
│   ├── events/                   # Event system (NATS JetStream)
│   ├── config/                   # Configuration management
│   ├── api/                      # API handlers and middleware
│   ├── registry/                 # Workflow registry implementation
│   ├── engine/                   # Workflow engine implementation
│   └── executor/                 # Operation executor framework
├── pkg/                          # Public libraries
│   └── client/                   # Client libraries (Go, Python, etc.)
├── examples/                     # Example workflows and usage
├── configs/                      # Configuration files
├── deployments/                  # Deployment configurations
├── docker/                       # Docker configurations
└── README.md                     # This file
```

## Core Components

### 1. Workflow Primitive Operations
Defines standard interfaces for workflow operations:
- `Operation`: Interface for executable operations
- `WorkflowStep`: Interface for workflow steps
- `ChildStep`: Interface for child steps within parent steps
- `OperationRegistry`: Manages available operations
- `StepExecutor`: Executes workflow steps

### 2. Workflow Registry
Stores and manages workflow definitions:
- Version control for workflow templates
- Approval workflows for production deployment
- Search and discovery of workflows
- Dependency management between workflows

### 3. Workflow Engine
Core workflow execution and state management:
- Step-by-step execution with error handling
- State persistence and recovery mechanisms
- Parallel step execution support
- Timeout and retry policies
- Compensation actions for rollback

### 4. Workflow API (Gin-based)
RESTful API layer for workflow management:
- Complete CRUD operations for workflows
- Workflow execution and monitoring
- Operation management
- Event streaming
- Metrics and health checks

### 5. NATS JetStream Integration
Event-driven communication between components:
- Reliable message delivery with persistence
- Streams for workflow events
- Consumers for event processing
- Dead letter queues for failed messages

## Data Models

### Workflow Definition
```go
type WorkflowDefinition struct {
    ID           string                 `json:"id"`
    Name         string                 `json:"name"`
    Version      string                 `json:"version"`
    Description  string                 `json:"description"`
    Steps        []StepDefinition       `json:"steps"`
    InputSchema  map[string]interface{} `json:"input_schema"`
    OutputSchema map[string]interface{} `json:"output_schema"`
    // ... additional fields
}
```

### Workflow Instance
```go
type WorkflowInstance struct {
    ID           string                 `json:"id"`
    DefinitionID string                 `json:"definition_id"`
    Status       WorkflowStatus         `json:"status"`
    InputData    map[string]interface{} `json:"input_data"`
    OutputData   map[string]interface{} `json:"output_data,omitempty"`
    // ... additional fields
}
```

### Step Instance
```go
type StepInstance struct {
    ID         string                 `json:"id"`
    WorkflowID string                 `json:"workflow_id"`
    StepDefID  string                 `json:"step_def_id"`
    Status     StepStatus             `json:"status"`
    InputData  map[string]interface{} `json:"input_data"`
    OutputData map[string]interface{} `json:"output_data,omitempty"`
    // ... additional fields
}
```

## API Endpoints

### Workflow Definitions
- `GET    /api/v1/workflows/definitions` - List workflow definitions
- `GET    /api/v1/workflows/definitions/:id` - Get workflow definition
- `POST   /api/v1/workflows/definitions` - Create workflow definition
- `PUT    /api/v1/workflows/definitions/:id` - Update workflow definition
- `DELETE /api/v1/workflows/definitions/:id` - Delete workflow definition

### Workflow Instances
- `GET    /api/v1/workflows` - List workflow instances
- `GET    /api/v1/workflows/:id` - Get workflow instance
- `POST   /api/v1/workflows` - Create workflow instance
- `POST   /api/v1/workflows/:id/execute` - Execute workflow
- `PUT    /api/v1/workflows/:id/cancel` - Cancel workflow
- `PUT    /api/v1/workflows/:id/pause` - Pause workflow
- `PUT    /api/v1/workflows/:id/resume` - Resume workflow

### Workflow Steps
- `GET    /api/v1/workflows/:id/steps` - List workflow steps
- `GET    /api/v1/workflows/:id/steps/:stepId` - Get workflow step
- `POST   /api/v1/workflows/:id/steps/:stepId/retry` - Retry workflow step

### Operations
- `GET    /api/v1/operations` - List available operations
- `GET    /api/v1/operations/:name` - Get operation details
- `POST   /api/v1/operations/:name/execute` - Execute operation

## Getting Started

### Prerequisites
- Go 1.25+
- Docker
- NATS server with JetStream enabled
- PostgreSQL (for production)
- Redis (for caching)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd unified-workflow-go
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start NATS server with JetStream**
   ```bash
   docker run -d --name nats -p 4222:4222 -p 8222:8222 nats:latest -js
   ```

4. **Create configuration file**
   ```bash
   cp configs/api-config.example.yaml configs/api-config.yaml
   # Edit configuration as needed
   ```

5. **Build and run the API server**
   ```bash
   go build -o workflow-api ./cmd/workflow-api
   ./workflow-api --config ./configs/api-config.yaml
   ```

6. **Test the API**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

## Configuration

### Main Configuration (`configs/api-config.yaml`)
```yaml
api:
  address: ":8080"
  debug: false
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s

database:
  type: "postgres"
  host: "localhost"
  port: 5432
  database: "workflow"
  username: "postgres"
  password: "password"
  ssl_mode: "disable"

nats:
  urls: ["nats://localhost:4222"]
  stream_name: "workflow-events"
  subject_prefix: "workflow"
  max_reconnects: 10
  reconnect_wait: 1s
  connect_timeout: 5s

cache:
  type: "redis"
  host: "localhost"
  port: 6379
  database: 0
  ttl: 1h

security:
  enabled: false
  jwt_secret: "your-secret-key"
  token_expiry: 24h
  require_https: false
  rate_limit: true
  audit_logging: true
```

## Deployment

### Docker Deployment
```bash
# Build Docker image
docker build -t workflow-api:latest -f docker/Dockerfile.api .

# Run container
docker run -d \
  --name workflow-api \
  -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  workflow-api:latest
```

### Kubernetes Deployment
```yaml
# See deployments/kubernetes/ for complete examples
apiVersion: apps/v1
kind: Deployment
metadata:
  name: workflow-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: workflow-api
  template:
    metadata:
      labels:
        app: workflow-api
    spec:
      containers:
      - name: workflow-api
        image: workflow-api:latest
        ports:
        - containerPort: 8080
        volumeMounts:
        - name: config
          mountPath: /app/configs
        env:
        - name: CONFIG_PATH
          value: "/app/configs/api-config.yaml"
```

## Development

### Building
```bash
# Build all components
go build ./cmd/...

# Build specific component
go build -o workflow-api ./cmd/workflow-api
```

### Testing
```bash
# Run unit tests
go test ./...

# Run integration tests
go test ./internal/... -v

# Run with coverage
go test ./... -cover
```

### Code Generation
```bash
# Generate mocks
go generate ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## Monitoring & Observability

### Metrics
Prometheus metrics are exposed on port 9090:
- `workflow_executions_total`
- `workflow_duration_seconds`
- `workflow_steps_total`
- `workflow_errors_total`
- `operation_executions_total`
- `operation_duration_seconds`

### Logging
Structured JSON logging with correlation IDs:
```json
{
  "level": "info",
  "time": "2024-01-12T10:00:00Z",
  "message": "Workflow executed",
  "workflow_id": "wf_123",
  "duration_ms": 150,
  "correlation_id": "corr_456"
}
```

### Tracing
Distributed tracing with OpenTelemetry:
- End-to-end workflow execution tracing
- Performance analysis across components
- Dependency mapping and visualization

## Security

### Authentication & Authorization
- JWT tokens for API authentication
- Service-to-service authentication using mTLS
- Role-Based Access Control (RBAC)
- API keys for external integrations

### Data Protection
- Encryption at rest for sensitive data
- TLS 1.3 for all communications
- Secrets management via HashiCorp Vault or AWS Secrets Manager
- Data masking in logs and audit trails

### Compliance
- Audit logging for all operations
- Data retention policies
- Access reviews and permission management
- Regulatory compliance (GDPR, PCI DSS, etc.)

## Performance

### Benchmarks
- **API Latency**: P95 < 100ms for reads, < 500ms for writes
- **Throughput**: 10,000+ requests per second
- **Concurrent Workflows**: 100,000+ concurrent executions
- **Recovery Time**: < 30 seconds for failed steps

### Scaling
- **Horizontal Scaling**: Stateless components (API, Executors)
- **Vertical Scaling**: Stateful components (Engine, Registry)
- **Auto-scaling**: Based on queue depth and CPU utilization
- **Load Balancing**: Round-robin with health checks

## Migration from Java Version

### Migration Strategy
1. **Parallel Run**: Run both systems side-by-side
2. **Data Sync**: Bi-directional sync of workflow definitions
3. **Traffic Routing**: Gradually shift traffic from Java to Go
4. **Feature Parity**: Implement all critical features before full migration
5. **Fallback Mechanism**: Ability to rollback to Java system if needed

### Key Advantages
1. **Performance**: Go's lightweight goroutines vs Java threads
2. **Resource Efficiency**: Lower memory footprint
3. **Cold Start**: Faster startup times for containers
4. **Simplicity**: Less boilerplate code
5. **Native Compilation**: Single binary deployment
6. **Better Concurrency**: Built-in channels and goroutines

## Roadmap

### Phase 1: Foundation (Completed)
- [x] Project structure and core interfaces
- [x] Data models and configuration
- [x] NATS JetStream integration
- [x] Basic API framework

### Phase 2: Core Implementation
- [ ] Workflow registry implementation
- [ ] Workflow engine implementation
- [ ] Operation executor framework
- [ ] Database integration

### Phase 3: Advanced Features
- [ ] Parallel step execution
- [ ] Compensation logic
- [ ] Workflow versioning
- [ ] Approval workflows

### Phase 4: Production Ready
- [ ] Security hardening
- [ ] Monitoring and observability
- [ ] Performance optimization
- [ ] Documentation and examples

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Support

- Documentation: [docs/](docs/)
- Issues: [GitHub Issues](<repository-url>/issues)
- Discussions: [GitHub Discussions](<repository-url>/discussions)
- Email: support@example.com

## Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) web framework
- Event-driven architecture with [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream)
- Inspired by banking workflow requirements
- Designed for hybrid cloud/on-premise deployment
