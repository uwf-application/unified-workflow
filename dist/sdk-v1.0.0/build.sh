#!/bin/bash

# Build script for Unified Workflow SDK

set -e

echo "Building Unified Workflow SDK..."

# Run tests
echo "1. Running tests..."
go test ./...

# Build examples
echo "2. Building examples..."
go build -o bin/basic_example examples/basic_usage.go
go build -o bin/http_example examples/http_integration.go

echo "Build complete!"
echo "Binaries available in bin/"
