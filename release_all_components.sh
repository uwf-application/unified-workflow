#!/bin/bash

# Unified Workflow - Complete Component Release Script
# This script builds and releases all Unified Workflow components

set -e

echo "=== Unified Workflow Component Release ==="
echo ""

# Configuration
VERSION=${1:-"v1.0.0"}
REGISTRY=${2:-"uwf-application"}  # Docker Hub organization or registry URL
COMPONENTS=("registry" "executor" "worker")
RELEASE_DIR="release-${VERSION}"
BUNDLE_NAME="unified-workflow-stack-${VERSION}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_step() {
    echo ""
    echo "=== $1 ==="
}

# Check prerequisites
check_prerequisites() {
    print_step "Checking Prerequisites"
    
    # Check Docker
    if command -v docker >/dev/null 2>&1; then
        print_success "Docker is installed"
    else
        print_error "Docker is not installed"
        exit 1
    fi
    
    # Check Docker daemon
    if docker info >/dev/null 2>&1; then
        print_success "Docker daemon is running"
    else
        print_error "Docker daemon is not running"
        exit 1
    fi
    
    # Check Go
    if command -v go >/dev/null 2>&1; then
        print_success "Go is installed"
    else
        print_warning "Go is not installed (needed for tests)"
    fi
    
    # Check git
    if command -v git >/dev/null 2>&1; then
        print_success "Git is installed"
    else
        print_warning "Git is not installed"
    fi
}

# Run tests
run_tests() {
    print_step "Running Tests"
    
    if command -v go >/dev/null 2>&1; then
        echo "Running Go tests with timeout..."
        # Run tests with timeout to avoid hanging tests
        # Use gtimeout on macOS, timeout on Linux
        if command -v gtimeout >/dev/null 2>&1; then
            TIMEOUT_CMD="gtimeout"
        elif command -v timeout >/dev/null 2>&1; then
            TIMEOUT_CMD="timeout"
        else
            TIMEOUT_CMD=""
        fi
        
        if [ -n "$TIMEOUT_CMD" ]; then
            $TIMEOUT_CMD 60s go test ./... -short || {
                if [ $? -eq 124 ]; then
                    print_warning "Tests timed out after 60 seconds (continuing anyway)"
                else
                    print_warning "Some tests failed (continuing anyway)"
                fi
            }
        else
            print_warning "No timeout command found, running tests without timeout..."
            go test ./... -short || print_warning "Some tests failed (continuing anyway)"
        fi
    else
        print_warning "Skipping tests (Go not installed)"
    fi
}

# Build binaries
build_binaries() {
    print_step "Building Binaries"
    
    # Create dist directory for binaries
    BIN_DIR="dist/binaries-${VERSION}"
    mkdir -p "$BIN_DIR"
    
    # Build registry binary
    echo "Building registry binary..."
    if go build -o "$BIN_DIR/registry-service" ./cmd/registry-api; then
        print_success "Built registry binary: $BIN_DIR/registry-service"
    else
        print_error "Failed to build registry binary"
        exit 1
    fi
    
    # Build executor binary
    echo "Building executor binary..."
    if go build -o "$BIN_DIR/executor-service" ./cmd/executor-api; then
        print_success "Built executor binary: $BIN_DIR/executor-service"
    else
        print_error "Failed to build executor binary"
        exit 1
    fi
    
    # Build worker binary
    echo "Building worker binary..."
    if go build -o "$BIN_DIR/workflow-worker" ./cmd/workflow-worker; then
        print_success "Built worker binary: $BIN_DIR/workflow-worker"
    else
        print_error "Failed to build worker binary"
        exit 1
    fi
    
    # Create binaries archive
    echo "Creating binaries archive..."
    tar -czf "dist/unified-workflow-binaries-${VERSION}.tar.gz" -C "$BIN_DIR" .
    
    print_success "Created binaries archive: dist/unified-workflow-binaries-${VERSION}.tar.gz"
    print_success "Binaries size: $(du -h "dist/unified-workflow-binaries-${VERSION}.tar.gz" | cut -f1)"
    
    # Create checksum for binaries
    shasum -a 256 "dist/unified-workflow-binaries-${VERSION}.tar.gz" > "dist/unified-workflow-binaries-${VERSION}.sha256"
    print_success "Created binaries checksum"
}

# Build Docker images
build_images() {
    print_step "Building Docker Images"
    
    for component in "${COMPONENTS[@]}"; do
        echo "Building ${component} service..."
        
        DOCKERFILE="Dockerfile.${component}"
        IMAGE_NAME="${REGISTRY}/uwf-${component}:${VERSION}"
        LATEST_IMAGE_NAME="${REGISTRY}/uwf-${component}:latest"
        
        if [ ! -f "$DOCKERFILE" ]; then
            print_error "Dockerfile not found: $DOCKERFILE"
            exit 1
        fi
        
        # Build the image
        if docker build -t "$IMAGE_NAME" -f "$DOCKERFILE" .; then
            print_success "Built ${component} image: $IMAGE_NAME"
            
            # Tag as latest
            docker tag "$IMAGE_NAME" "$LATEST_IMAGE_NAME"
            print_success "Tagged as latest: $LATEST_IMAGE_NAME"
        else
            print_error "Failed to build ${component} image"
            exit 1
        fi
    done
}

# Create release bundle
create_release_bundle() {
    print_step "Creating Release Bundle"
    
    # Create release directory
    rm -rf "$RELEASE_DIR"
    mkdir -p "$RELEASE_DIR"
    
    # Copy essential files
    echo "Copying configuration files..."
    cp docker-compose.yml "$RELEASE_DIR/"
    cp README.md "$RELEASE_DIR/"
    cp LICENSE "$RELEASE_DIR/" 2>/dev/null || echo "LICENSE not found, skipping"
    
    # Create version-specific docker-compose file
    echo "Creating versioned docker-compose file..."
    sed "s|build:|image: ${REGISTRY}/uwf-|g" docker-compose.yml | \
    sed "s|context: .||g" | \
    sed "s|dockerfile: Dockerfile.registry|registry:${VERSION}|g" | \
    sed "s|dockerfile: Dockerfile.executor|executor:${VERSION}|g" | \
    sed "s|dockerfile: Dockerfile.worker|worker:${VERSION}|g" > "$RELEASE_DIR/docker-compose.${VERSION}.yml"
    
    # Create deployment guide
    cat > "$RELEASE_DIR/DEPLOYMENT.md" << EOF
# Unified Workflow ${VERSION} - Deployment Guide

## Quick Start

### Using Docker Compose
\`\`\`bash
# Download and extract the release bundle
tar -xzf ${BUNDLE_NAME}.tar.gz
cd ${BUNDLE_NAME}

# Start all services
docker-compose -f docker-compose.${VERSION}.yml up -d
\`\`\`

### Using Docker Images Directly
\`\`\`bash
# Pull the images
docker pull ${REGISTRY}/uwf-registry:${VERSION}
docker pull ${REGISTRY}/uwf-executor:${VERSION}
docker pull ${REGISTRY}/uwf-worker:${VERSION}

# Run individually
docker run -d -p 8080:8080 --name uwf-registry ${REGISTRY}/uwf-registry:${VERSION}
docker run -d -p 8081:8081 --name uwf-executor ${REGISTRY}/uwf-executor:${VERSION}
docker run -d --name uwf-worker ${REGISTRY}/uwf-worker:${VERSION}
\`\`\`

## Service Ports

- **Registry Service**: 8080
- **Executor Service**: 8081
- **Worker Service**: No exposed ports (consumes from NATS)

## Health Checks

- Registry: http://localhost:8080/health
- Executor: http://localhost:8081/health

## Configuration

### Environment Variables

**Registry Service:**
- \`REGISTRY_PORT\`: Port to listen on (default: 8080)
- \`LOG_LEVEL\`: Log level (default: info)

**Executor Service:**
- \`EXECUTOR_PORT\`: Port to listen on (default: 8081)
- \`REGISTRY_SERVICE_URL\`: URL of registry service
- \`NATS_URL\`: NATS server URL
- \`LOG_LEVEL\`: Log level (default: info)

**Worker Service:**
- \`NATS_URL\`: NATS server URL
- \`REGISTRY_SERVICE_URL\`: URL of registry service
- \`LOG_LEVEL\`: Log level (default: info)
- \`QUEUE_TYPE\`: Queue type (default: nats)

## Troubleshooting

### Check Service Logs
\`\`\`bash
docker logs uwf-registry
docker logs uwf-executor
docker logs uwf-worker
\`\`\`

### Check Service Health
\`\`\`bash
curl http://localhost:8080/health
curl http://localhost:8081/health
\`\`\`

### Restart Services
\`\`\`bash
docker-compose -f docker-compose.${VERSION}.yml restart
\`\`\`

## Support

For issues and questions, please refer to:
- GitHub: https://github.com/uwf-application/unified-workflow
- Documentation: Included in release bundle
EOF
    
    # Create README for the bundle
    cat > "$RELEASE_DIR/README.md" << EOF
# Unified Workflow Stack ${VERSION}

This bundle contains everything needed to deploy the Unified Workflow platform.

## Contents

1. **Docker Images** (pre-built)
   - \`${REGISTRY}/uwf-registry:${VERSION}\`
   - \`${REGISTRY}/uwf-executor:${VERSION}\`
   - \`${REGISTRY}/uwf-worker:${VERSION}\`

2. **Configuration Files**
   - \`docker-compose.${VERSION}.yml\` - Complete stack configuration
   - \`DEPLOYMENT.md\` - Deployment guide

3. **Documentation**
   - Main README
   - License

## Quick Deployment

\`\`\`bash
# Extract the bundle
tar -xzf ${BUNDLE_NAME}.tar.gz
cd ${BUNDLE_NAME}

# Start all services
docker-compose -f docker-compose.${VERSION}.yml up -d

# Verify services are running
curl http://localhost:8080/health
curl http://localhost:8081/health
\`\`\`

## Component Details

### 1. Registry Service (Port 8080)
- Manages workflow definitions
- Provides REST API for workflow management
- Health endpoint: /health

### 2. Executor Service (Port 8081)
- Executes workflows
- Provides REST API for workflow execution
- Health endpoint: /health

### 3. Worker Service
- Processes workflow steps asynchronously
- Consumes messages from NATS queue
- No exposed ports

### 4. NATS JetStream
- Message queue for worker communication
- Included in docker-compose configuration

## Next Steps

1. **Test the deployment**: Use the health endpoints
2. **Create workflows**: Use the registry API
3. **Execute workflows**: Use the executor API
4. **Monitor**: Check service logs

## Support

- GitHub: https://github.com/uwf-application/unified-workflow
- Issues: https://github.com/uwf-application/unified-workflow/issues
EOF
    
    # Create bundle archive
    echo "Creating bundle archive..."
    tar -czf "${BUNDLE_NAME}.tar.gz" "$RELEASE_DIR"
    
    print_success "Created release bundle: ${BUNDLE_NAME}.tar.gz"
    print_success "Bundle size: $(du -h "${BUNDLE_NAME}.tar.gz" | cut -f1)"
}

# Generate component manifests
generate_manifests() {
    print_step "Generating Deployment Manifests"
    
    MANIFESTS_DIR="$RELEASE_DIR/manifests"
    mkdir -p "$MANIFESTS_DIR"
    
    # Generate Kubernetes deployment manifests
    for component in "${COMPONENTS[@]}"; do
        cat > "$MANIFESTS_DIR/${component}-deployment.yaml" << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: uwf-${component}
  labels:
    app: unified-workflow
    component: ${component}
    version: ${VERSION}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: unified-workflow
      component: ${component}
  template:
    metadata:
      labels:
        app: unified-workflow
        component: ${component}
        version: ${VERSION}
    spec:
      containers:
      - name: ${component}
        image: ${REGISTRY}/uwf-${component}:${VERSION}
        ports:
EOF
        
        # Add ports based on component
        if [ "$component" = "registry" ]; then
            cat >> "$MANIFESTS_DIR/${component}-deployment.yaml" << EOF
        - containerPort: 8080
          name: http
EOF
        elif [ "$component" = "executor" ]; then
            cat >> "$MANIFESTS_DIR/${component}-deployment.yaml" << EOF
        - containerPort: 8081
          name: http
EOF
        fi
        
        cat >> "$MANIFESTS_DIR/${component}-deployment.yaml" << EOF
        env:
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 5
          periodSeconds: 5
EOF
        
        print_success "Generated manifest: ${component}-deployment.yaml"
    done
    
    # Generate service manifests
    for component in "registry" "executor"; do
        PORT="8080"
        if [ "$component" = "executor" ]; then
            PORT="8081"
        fi
        
        cat > "$MANIFESTS_DIR/${component}-service.yaml" << EOF
apiVersion: v1
kind: Service
metadata:
  name: uwf-${component}
  labels:
    app: unified-workflow
    component: ${component}
spec:
  selector:
    app: unified-workflow
    component: ${component}
  ports:
  - port: ${PORT}
    targetPort: http
    name: http
  type: ClusterIP
EOF
        
        print_success "Generated service: ${component}-service.yaml"
    done
}

# Push images to registry (optional)
push_images() {
    print_step "Pushing Docker Images"
    
    read -p "Do you want to push images to registry? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Pushing images to registry: $REGISTRY"
        
        for component in "${COMPONENTS[@]}"; do
            IMAGE_NAME="${REGISTRY}/uwf-${component}:${VERSION}"
            LATEST_IMAGE_NAME="${REGISTRY}/uwf-${component}:latest"
            
            echo "Pushing ${component} image..."
            if docker push "$IMAGE_NAME"; then
                print_success "Pushed: $IMAGE_NAME"
            else
                print_error "Failed to push: $IMAGE_NAME"
                print_warning "Make sure you're logged into the registry:"
                print_warning "  docker login"
                break
            fi
            
            if docker push "$LATEST_IMAGE_NAME"; then
                print_success "Pushed: $LATEST_IMAGE_NAME"
            fi
        done
    else
        print_warning "Skipping image push"
    fi
}

# Create checksums
create_checksums() {
    print_step "Creating Checksums"
    
    echo "Creating SHA256 checksums..."
    shasum -a 256 "${BUNDLE_NAME}.tar.gz" > "${BUNDLE_NAME}.sha256"
    
    print_success "Created checksum: ${BUNDLE_NAME}.sha256"
    echo "Checksum: $(cat "${BUNDLE_NAME}.sha256")"
}

# Main execution
main() {
    echo "Release Version: $VERSION"
    echo "Registry: $REGISTRY"
    echo ""
    
    check_prerequisites
    run_tests
    build_binaries
    build_images
    create_release_bundle
    generate_manifests
    create_checksums
    push_images
    
    print_step "Release Complete"
    echo ""
    echo "🎉 Unified Workflow ${VERSION} release completed!"
    echo ""
    echo "Release artifacts created:"
    echo "  - ${BUNDLE_NAME}.tar.gz (Release bundle)"
    echo "  - ${BUNDLE_NAME}.sha256 (Checksum)"
    echo "  - dist/unified-workflow-binaries-${VERSION}.tar.gz (Binaries)"
    echo "  - dist/unified-workflow-binaries-${VERSION}.sha256 (Binaries checksum)"
    echo "  - ${RELEASE_DIR}/ (Release directory)"
    echo ""
    echo "Binaries created in dist/:"
    echo "  - registry-service"
    echo "  - executor-service"
    echo "  - workflow-worker"
    echo ""
    echo "Docker images built:"
    for component in "${COMPONENTS[@]}"; do
        echo "  - ${REGISTRY}/uwf-${component}:${VERSION}"
        echo "  - ${REGISTRY}/uwf-${component}:latest"
    done
    echo ""
    echo "Next steps:"
    echo "1. Test the binaries:"
    echo "   ./dist/binaries-${VERSION}/registry-service"
    echo ""
    echo "2. Test the images locally:"
    echo "   docker-compose -f ${RELEASE_DIR}/docker-compose.${VERSION}.yml up -d"
    echo ""
    echo "3. Push to GitHub release:"
    echo "   Upload ${BUNDLE_NAME}.tar.gz, ${BUNDLE_NAME}.sha256,"
    echo "   dist/unified-workflow-binaries-${VERSION}.tar.gz,"
    echo "   dist/unified-workflow-binaries-${VERSION}.sha256"
    echo ""
    echo "4. Update documentation"
    echo "5. Notify stakeholders"
    echo ""
}

# Run main function
main "$@"