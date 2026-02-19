# Unified Workflow Deployment CLI

A comprehensive CLI tool for building and deploying Unified Workflow services to the test environment via the jump server.

## 🚀 Overview

The deployment CLI automates the process of:
1. **Syncing code** from local machine to jump server (10.200.1.2)
2. **Building Docker images** on the jump server
3. **Pushing images** to Harbor registry (172.30.75.78:9080)
4. **Deploying services** to target servers (172.30.75.85)

## 📋 Prerequisites

### 1. Local Machine Requirements
- Go 1.20+ installed
- SSH access to jump server (10.200.1.2)
- SSH key configured for `khassangali@10.200.1.2`
- `rsync` command available

### 2. Jump Server Requirements
- Docker installed and running
- Access to Harbor registry (172.30.75.78:9080)
- `/tmp/uwf-deploy` directory writable

### 3. Network Access
- VPN connection to 10.200.0.0/16 network
- Access to jump server (10.200.1.2)
- Access to Harbor registry (172.30.75.78:9080)
- Access to target servers (172.30.75.85)

## 🔧 Installation

### Build from Source
```bash
cd /Users/khassangali/dev/unified-workflow-go
go build ./cmd/uwf-cli/
```

### Install Globally
```bash
go install ./cmd/uwf-cli/
```

## ⚙️ Configuration

### Initialize Configuration
```bash
uwf-cli deploy init \
  --jump-server 10.200.1.2 \
  --harbor 172.30.75.78:9080 \
  --environment test
```

### Configuration File (`deploy-config.yaml`)
The CLI uses a YAML configuration file with the following structure:

```yaml
# Jump server configuration
jump_server: "10.200.1.2"
jump_username: "khassangali"

# Harbor registry configuration
harbor_registry: "172.30.75.78:9080"
harbor_username: "zh.akhmetkarimov"

# Services configuration
services:
  workflow-worker:
    name: "workflow-worker"
    dockerfile: "Dockerfile.worker"
    registry: "172.30.75.78:9080/taf/workflow-worker"
    target_host: "172.30.75.85"
    compose_path: "/opt/workflow/docker-compose.yml"
```

## 🎯 Usage Examples

### 1. Check Deployment Status
```bash
uwf-cli deploy status
```

### 2. Sync Code to Jump Server
```bash
uwf-cli deploy sync
```

### 3. Build All Services
```bash
uwf-cli deploy build --all
```

### 4. Push All Services to Harbor
```bash
uwf-cli deploy push --all
```

### 5. Complete Deployment (Build + Push)
```bash
uwf-cli deploy all
```

### 6. Deploy Specific Service
```bash
uwf-cli deploy service workflow-worker
```

### 7. Verify Deployment
```bash
uwf-cli deploy verify
```

### 8. Initialize Configuration
```bash
uwf-cli deploy init
```

## 📊 Command Reference

### `uwf-cli deploy status`
Displays the current deployment configuration and status of all services.

### `uwf-cli deploy sync`
Syncs local code changes to the jump server using `rsync`.

**Options:**
- `--jump-server`: Override jump server address

### `uwf-cli deploy build`
Builds Docker images for specified services on the jump server.

**Options:**
- `--services`: Comma-separated list of services to build
- `--all`: Build all services

### `uwf-cli deploy push`
Pushes Docker images to Harbor registry.

**Options:**
- `--services`: Comma-separated list of services to push
- `--all`: Push all services

### `uwf-cli deploy all`
Complete deployment workflow (build + push).

**Options:**
- `--services`: Comma-separated list of services to deploy
- `--skip-build`: Skip build phase
- `--skip-push`: Skip push phase

### `uwf-cli deploy service [name]`
Deploys a specific service.

**Options:**
- `--skip-build`: Skip build phase
- `--skip-push`: Skip push phase

### `uwf-cli deploy verify`
Verifies deployment by checking:
- Jump server connectivity
- Harbor registry accessibility
- Service health (future)

### `uwf-cli deploy init`
Initializes deployment configuration.

**Options:**
- `--jump-server`: Jump server address (default: 10.200.1.2)
- `--harbor`: Harbor registry address (default: 172.30.75.78:9080)
- `--environment`: Environment name (default: test)

## 🔄 Deployment Workflow

### Phase 1: Code Sync
```
Local Machine → rsync → Jump Server (10.200.1.2:/tmp/uwf-deploy/)
```

### Phase 2: Build Images
```
Jump Server → docker build → Local Docker Images
```

### Phase 3: Push to Registry
```
Jump Server → docker push → Harbor Registry (172.30.75.78:9080)
```

### Phase 4: Deploy to Target
```
Harbor Registry → Target Server (172.30.75.85) → docker-compose up
```

## 🏗️ Service Architecture

### Supported Services
1. **workflow-worker** - Workflow execution engine
2. **workflow-registry** - Workflow definition registry
3. **workflow-executor** - Task execution service
4. **workflow-api** - REST API service

### Docker Images
- Built from: `Dockerfile.worker`, `Dockerfile.registry`, etc.
- Tagged as: `172.30.75.78:9080/taf/[service-name]:latest`
- Pushed to: Harbor registry at `172.30.75.78:9080`

### Target Deployment
- Host: `172.30.75.85` (Workflow SVC server)
- Compose file: `/opt/workflow/docker-compose.yml`
- Username: `zh.akhmetkarimov`

## 🔐 Security Considerations

### SSH Authentication
- Uses SSH key authentication for jump server
- Private key should be at `~/.ssh/id_rsa`
- User: `khassangali@10.200.1.2`

### Registry Authentication
- Harbor registry credentials configured on jump server
- User: `zh.akhmetkarimov`
- Docker login required on jump server

### Network Security
- All operations go through jump server
- No direct access from local to internal network
- VPN required for initial connection

## 🚨 Troubleshooting

### Common Issues

#### 1. SSH Connection Failed
```bash
# Test SSH connection
ssh khassangali@10.200.1.2 "echo Connected"
```

#### 2. rsync Not Available
```bash
# Install rsync on macOS
brew install rsync
```

#### 3. Docker Not Installed on Jump Server
```bash
# Check Docker on jump server
ssh khassangali@10.200.1.2 "docker --version"
```

#### 4. Harbor Registry Unreachable
```bash
# Test Harbor access from jump server
ssh khassangali@10.200.1.2 "curl -s http://172.30.75.78:9080/v2/_catalog"
```

### Debug Mode
```bash
# Enable verbose output
uwf-cli deploy --verbose [command]
```

## 📈 Advanced Usage

### Incremental Deployment
```bash
# Build and push only changed services (future feature)
uwf-cli deploy changed --since=HEAD~1
```

### Rollback Deployment
```bash
# Rollback to previous version (future feature)
uwf-cli deploy rollback v1.2.3
```

### CI/CD Integration
```yaml
# GitHub Actions example
name: Deploy to Test
on: [push]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: |
          make build-cli
          ./uwf-cli deploy all --auto-approve
```

## 🔮 Future Enhancements

### Planned Features
1. **Incremental builds** - Only build changed services
2. **Rollback support** - Automatic rollback on failure
3. **Health checks** - Verify service health after deployment
4. **Notifications** - Slack/email notifications
5. **Multi-environment** - Support for staging/production
6. **Dry-run mode** - Simulate deployment without changes
7. **Deployment history** - Track deployment versions
8. **Auto-scaling** - Scale services based on load

### Integration Points
- **GitHub Actions** - CI/CD pipeline integration
- **Slack** - Deployment notifications
- **Prometheus** - Deployment metrics
- **Grafana** - Deployment dashboards

## 📚 Related Documentation

- [VPN Setup Guide](ANTIFRAUD_PRIMITIVE_INTEGRATION.md)
- [Jump Server Access](CLI_INSTALLATION.md)
- [Harbor Registry](RELEASE_CLI.md)
- [Service Architecture](HLD_1000_TPS_DEPLOYMENT_PLAN.md)

## 🆘 Support

### Getting Help
1. Check `uwf-cli deploy --help` for command reference
2. Review configuration in `deploy-config.yaml`
3. Test SSH connectivity to jump server
4. Verify Docker installation on jump server

### Reporting Issues
1. Enable verbose mode: `uwf-cli deploy --verbose [command]`
2. Check deployment logs: `tail -f deploy.log`
3. Test individual components manually
4. Contact infrastructure team for network issues

## 🎉 Quick Start

### First-Time Setup
```bash
# 1. Build the CLI
go build ./cmd/uwf-cli/

# 2. Initialize configuration
./uwf-cli deploy init

# 3. Verify setup
./uwf-cli deploy verify

# 4. Sync code
./uwf-cli deploy sync

# 5. Deploy all services
./uwf-cli deploy all
```

### Daily Deployment
```bash
# Sync and deploy changes
./uwf-cli deploy sync
./uwf-cli deploy all

# Or deploy specific service
./uwf-cli deploy service workflow-worker
```

## 📄 License

This deployment CLI is part of the Unified Workflow system. See the main project README for license information.