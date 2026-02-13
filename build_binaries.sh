#!/bin/bash

# Simple script to build Unified Workflow binaries
# This script builds all component binaries and places them in dist/

set -e

echo "=== Building Unified Workflow Binaries ==="
echo ""

VERSION=${1:-"v1.0.0"}
BIN_DIR="dist/binaries-${VERSION}"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check prerequisites
echo "Checking prerequisites..."
if ! command -v go >/dev/null 2>&1; then
    print_error "Go is not installed"
    exit 1
fi
print_success "Go is installed"

# Create dist directory for binaries
mkdir -p "$BIN_DIR"
mkdir -p "dist"

echo ""
echo "Building binaries for version: $VERSION"
echo ""

# Build registry binary
echo "Building registry binary..."
if go build -o "$BIN_DIR/registry-service" ./cmd/registry-api; then
    print_success "Built registry binary: $BIN_DIR/registry-service"
    ls -lh "$BIN_DIR/registry-service"
else
    print_error "Failed to build registry binary"
    exit 1
fi

# Build executor binary
echo "Building executor binary..."
if go build -o "$BIN_DIR/executor-service" ./cmd/executor-api; then
    print_success "Built executor binary: $BIN_DIR/executor-service"
    ls -lh "$BIN_DIR/executor-service"
else
    print_error "Failed to build executor binary"
    exit 1
fi

# Build worker binary
echo "Building worker binary..."
if go build -o "$BIN_DIR/workflow-worker" ./cmd/workflow-worker; then
    print_success "Built worker binary: $BIN_DIR/workflow-worker"
    ls -lh "$BIN_DIR/workflow-worker"
else
    print_error "Failed to build worker binary"
    exit 1
fi

# Create binaries archive
echo ""
echo "Creating binaries archive..."
tar -czf "dist/unified-workflow-binaries-${VERSION}.tar.gz" -C "$BIN_DIR" .

print_success "Created binaries archive: dist/unified-workflow-binaries-${VERSION}.tar.gz"
print_success "Binaries size: $(du -h "dist/unified-workflow-binaries-${VERSION}.tar.gz" | cut -f1)"

# Create checksum for binaries
echo ""
echo "Creating checksum..."
shasum -a 256 "dist/unified-workflow-binaries-${VERSION}.tar.gz" > "dist/unified-workflow-binaries-${VERSION}.sha256"
print_success "Created binaries checksum: dist/unified-workflow-binaries-${VERSION}.sha256"

echo ""
echo "Checksum:"
cat "dist/unified-workflow-binaries-${VERSION}.sha256"

echo ""
echo "🎉 Binaries built successfully!"
echo ""
echo "Binaries location: $BIN_DIR/"
echo "  - registry-service"
echo "  - executor-service"
echo "  - workflow-worker"
echo ""
echo "Archive: dist/unified-workflow-binaries-${VERSION}.tar.gz"
echo "Checksum: dist/unified-workflow-binaries-${VERSION}.sha256"
echo ""
echo "To test a binary:"
echo "  ./$BIN_DIR/registry-service --help"
echo ""
echo "To extract and use:"
echo "  tar -xzf dist/unified-workflow-binaries-${VERSION}.tar.gz"
echo "  ./registry-service"