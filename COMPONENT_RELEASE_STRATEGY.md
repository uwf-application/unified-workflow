# Unified Workflow - Component Release Strategy

## Overview
This document outlines the release strategy for all Unified Workflow components. After successfully releasing the SDK v1.0.0, we need to create a systematic approach for releasing all other components.

## Component Architecture

### 1. Core Services (Docker Containers)
```
1. Registry Service (port 8080)
   - Dockerfile: Dockerfile.registry
   - Binary: registry-service
   - Source: cmd/registry-api/

2. Executor Service (port 8081)
   - Dockerfile: Dockerfile.executor
   - Binary: executor-service
   - Source: cmd/executor-api/

3. Worker Service
   - Dockerfile: Dockerfile.worker
   - Binary: workflow-worker
   - Source: cmd/workflow-worker/
```

### 2. Client Libraries
```
1. Go SDK (Released: v1.0.0)
   - Location: pkg/client/sdk/
   - Purpose: Client integration with workflow API

2. Executor Client
   - Location: pkg/client/executor/
   - Purpose: Direct executor service communication

3. Registry Client
   - Location: pkg/client/registry/
   - Purpose: Registry service communication
```

### 3. Infrastructure
```
1. NATS JetStream
   - External dependency
   - Message queue for worker communication

2. Configuration Files
   - docker-compose.yml
   - Environment configurations
```

## Release Strategy

### Phase 1: Individual Component Releases

#### 1.1 Registry Service Release
```bash
# Build and tag Docker image
docker build -t uwf-registry:v1.0.0 -f Dockerfile.registry .
docker tag uwf-registry:v1.0.0 uwf-registry:latest

# Push to container registry
docker push your-registry/uwf-registry:v1.0.0
docker push your-registry/uwf-registry:latest

# Create GitHub release
# - Include Dockerfile
# - Include configuration examples
# - Include migration guides
```

#### 1.2 Executor Service Release
```bash
# Build and tag Docker image
docker build -t uwf-executor:v1.0.0 -f Dockerfile.executor .
docker tag uwf-executor:v1.0.0 uwf-executor:latest

# Push to container registry
docker push your-registry/uwf-executor:v1.0.0
docker push your-registry/uwf-executor:latest
```

#### 1.3 Worker Service Release
```bash
# Build and tag Docker image
docker build -t uwf-worker:v1.0.0 -f Dockerfile.worker .
docker tag uwf-worker:v1.0.0 uwf-worker:latest

# Push to container registry
docker push your-registry/uwf-worker:v1.0.0
docker push your-registry/uwf-worker:latest
```

### Phase 2: Bundled Release (Docker Compose)

#### 2.1 Complete Stack Release
```bash
# Create versioned docker-compose file
cp docker-compose.yml docker-compose.v1.0.0.yml

# Update image tags in docker-compose
# Change: image: your-registry/uwf-registry:v1.0.0
# Change: image: your-registry/uwf-executor:v1.0.0
# Change: image: your-registry/uwf-worker:v1.0.0

# Create release bundle
tar -czf unified-workflow-stack-v1.0.0.tar.gz \
  docker-compose.v1.0.0.yml \
  configs/ \
  README.md \
  LICENSE
```

#### 2.2 Helm Chart Release (Kubernetes)
```yaml
# Create Helm chart structure
unified-workflow/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── registry-deployment.yaml
│   ├── executor-deployment.yaml
│   ├── worker-deployment.yaml
│   ├── nats-deployment.yaml
│   └── service.yaml
└── README.md
```

### Phase 3: Client Library Releases

#### 3.1 SDK Updates
- Already released: v1.0.0
- Future: v1.1.0, v1.2.0 with new features

#### 3.2 Language-specific SDKs
```bash
# Python SDK
# Java SDK
# JavaScript/TypeScript SDK
# REST API clients
```

## Release Automation

### GitHub Actions Workflow

```yaml
name: Release Components
on:
  push:
    tags:
      - 'v*'

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Container Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ secrets.REGISTRY_URL }}
          username: ${{ secrets.REGISTRY_USERNAME }}
          password: ${{ secrets.REGISTRY_PASSWORD }}
      
      - name: Build and push Registry Service
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./Dockerfile.registry
          push: true
          tags: |
            ${{ secrets.REGISTRY_URL }}/uwf-registry:${{ github.ref_name }}
            ${{ secrets.REGISTRY_URL }}/uwf-registry:latest
      
      - name: Build and push Executor Service
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./Dockerfile.executor
          push: true
          tags: |
            ${{ secrets.REGISTRY_URL }}/uwf-executor:${{ github.ref_name }}
            ${{ secrets.REGISTRY_URL }}/uwf-executor:latest
      
      - name: Build and push Worker Service
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./Dockerfile.worker
          push: true
          tags: |
            ${{ secrets.REGISTRY_URL }}/uwf-worker:${{ github.ref_name }}
            ${{ secrets.REGISTRY_URL }}/uwf-worker:latest
      
      - name: Create Release Bundle
        run: |
          mkdir -p release-bundle
          cp docker-compose.yml release-bundle/
          cp README.md release-bundle/
          cp LICENSE release-bundle/
          tar -czf unified-workflow-${{ github.ref_name }}.tar.gz release-bundle/
      
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            unified-workflow-${{ github.ref_name }}.tar.gz
          generate_release_notes: true
```

## Versioning Strategy

### Semantic Versioning
```
MAJOR.MINOR.PATCH

MAJOR: Breaking changes
MINOR: New features (backward compatible)
PATCH: Bug fixes (backward compatible)
```

### Component Version Alignment
```
All components share the same version number:
- Registry Service: v1.0.0
- Executor Service: v1.0.0
- Worker Service: v1.0.0
- SDK: v1.0.0

Benefits:
- Simplified dependency management
- Clear compatibility matrix
- Easier troubleshooting
```

## Release Checklist

### Pre-Release
- [ ] Run all tests
- [ ] Update version numbers
- [ ] Update CHANGELOG.md
- [ ] Update documentation
- [ ] Verify backward compatibility
- [ ] Security scan of Docker images

### Release Process
- [ ] Build Docker images
- [ ] Push to container registry
- [ ] Create GitHub release
- [ ] Update deployment guides
- [ ] Notify stakeholders

### Post-Release
- [ ] Monitor deployment
- [ ] Gather feedback
- [ ] Address critical issues
- [ ] Plan next release

## Container Registry Setup

### Docker Hub
```bash
# Organization: uwf-application
# Images:
# - uwf-application/registry
# - uwf-application/executor
# - uwf-application/worker
```

### AWS ECR
```bash
# Repository: unified-workflow
# Images:
# - registry.docker.aws.com/uwf/registry
# - registry.docker.aws.com/uwf/executor
# - registry.docker.aws.com/uwf/worker
```

### Google Container Registry
```bash
# Project: unified-workflow
# Images:
# - gcr.io/unified-workflow/registry
# - gcr.io/unified-workflow/executor
# - gcr.io/unified-workflow/worker
```

## Deployment Options

### Option 1: Docker Compose (Development)
```bash
# Using released images
docker-compose -f docker-compose.v1.0.0.yml up -d
```

### Option 2: Kubernetes (Production)
```bash
# Using Helm chart
helm install unified-workflow ./charts/unified-workflow \
  --version 1.0.0 \
  --set registry.image.tag=v1.0.0 \
  --set executor.image.tag=v1.0.0 \
  --set worker.image.tag=v1.0.0
```

### Option 3: Bare Metal
```bash
# Download binaries
# Configure and run services
./registry-service --config config.yaml
./executor-service --config config.yaml
./workflow-worker --config config.yaml
```

## Release Artifacts

### For Each Release
1. **Docker Images**
   - uwf-registry:v1.0.0
   - uwf-executor:v1.0.0
   - uwf-worker:v1.0.0

2. **Configuration Files**
   - docker-compose.v1.0.0.yml
   - Kubernetes manifests
   - Helm chart

3. **Documentation**
   - Release notes
   - Upgrade guide
   - Configuration reference
   - API documentation

4. **Client Libraries**
   - Go SDK package
   - REST API documentation
   - Example applications

## Quality Assurance

### Testing Strategy
1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Component interaction testing
3. **End-to-End Tests**: Complete workflow testing
4. **Performance Tests**: Load and stress testing
5. **Security Tests**: Vulnerability scanning

### Monitoring
1. **Health Checks**: Built into each service
2. **Metrics**: Prometheus metrics endpoint
3. **Logging**: Structured JSON logging
4. **Tracing**: Distributed tracing support

## Support and Maintenance

### Support Timeline
```
v1.0.0: Active support for 12 months
v1.1.0: Active support for 12 months
v1.x.x: Security patches for 24 months
```

### Upgrade Path
```
v1.0.0 → v1.1.0: Seamless upgrade
v1.x.x → v2.0.0: Migration guide required
```

## Next Steps

### Immediate Actions
1. [ ] Set up container registry
2. [ ] Create GitHub Actions workflow
3. [ ] Build and test Docker images
4. [ ] Create Helm chart
5. [ ] Prepare documentation

### Short-term Goals
1. [ ] Release Registry Service v1.0.0
2. [ ] Release Executor Service v1.0.0
3. [ ] Release Worker Service v1.0.0
4. [ ] Create bundled release
5. [ ] Set up automated testing

### Long-term Goals
1. [ ] Multi-architecture Docker images
2. [ ] Signed container images
3. [ ] Vulnerability scanning pipeline
4. [ ] Canary deployment strategy
5. [ ] Blue-green deployment support

## Conclusion

This release strategy provides a comprehensive approach for releasing all Unified Workflow components. By following this plan, we can ensure consistent, reliable releases that meet production requirements while maintaining backward compatibility and providing clear upgrade paths for users.