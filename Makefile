# Unified Workflow System Makefile

.PHONY: help build test clean run lint deps generate

# Variables
BINARY_NAME=workflow-api
CONFIG_PATH=configs/api-config.yaml
GO=go
GOFMT=gofmt
GOLINT=golangci-lint

# Default target
help:
	@echo "Unified Workflow System Build Tools"
	@echo ""
	@echo "Available targets:"
	@echo "  help     - Show this help message"
	@echo "  deps     - Download dependencies"
	@echo "  build    - Build all binaries"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter"
	@echo "  format   - Format code"
	@echo "  clean    - Clean build artifacts"
	@echo "  run      - Run the API server"
	@echo "  generate - Generate code"
	@echo "  docker   - Build Docker images"
	@echo "  all      - Run deps, format, lint, test, build"

# Download dependencies
deps:
	$(GO) mod download
	$(GO) mod verify

# Build all binaries
build: deps
	@echo "Building workflow-api..."
	$(GO) build -o bin/$(BINARY_NAME) ./cmd/workflow-api
	@echo "Build complete: bin/$(BINARY_NAME)"

# Run tests
test: deps
	@echo "Running tests..."
	$(GO) test ./... -v -cover

# Run linter
lint: deps
	@echo "Running linter..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOLINT) run ./...

# Format code
format:
	@echo "Formatting code..."
	$(GOFMT) -w ./cmd/ ./internal/ ./pkg/

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf coverage.out
	rm -rf logs/
	$(GO) clean

# Run the API server
run: build
	@echo "Starting workflow API server..."
	./bin/$(BINARY_NAME) --config $(CONFIG_PATH)

# Generate code
generate: deps
	@echo "Generating code..."
	$(GO) generate ./...

# Build Docker images
docker:
	@echo "Building Docker images..."
	docker build -t workflow-api:latest -f docker/Dockerfile.api .
	docker build -t workflow-registry:latest -f docker/Dockerfile.registry .
	docker build -t workflow-engine:latest -f docker/Dockerfile.engine .

# Run all checks and build with coverage
all: deps format lint test build coverage
	@echo "All tasks completed successfully!"

# CI pipeline target
ci: deps lint test ci-coverage build
	@echo "CI pipeline completed successfully!"

# Development environment setup
dev-setup:
	@echo "Setting up development environment..."
	
	# Create necessary directories
	mkdir -p bin logs configs deployments docker
	
	# Copy example config if it doesn't exist
	if [ ! -f configs/api-config.yaml ]; then \
		cp configs/api-config.example.yaml configs/api-config.yaml; \
		echo "Created configs/api-config.yaml from example"; \
	fi
	
	# Install development tools
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/go-delve/delve/cmd/dlv@latest
	$(GO) install github.com/vektra/mockery/v2@latest
	
	@echo "Development environment setup complete!"

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	# Add migration commands here

migrate-down:
	@echo "Rolling back database migrations..."
	# Add migration rollback commands here

# Start all services with Docker Compose
compose-up:
	@echo "Starting all services with Docker Compose..."
	docker-compose -f deployments/docker-compose.yml up -d

compose-down:
	@echo "Stopping all services..."
	docker-compose -f deployments/docker-compose.yml down

# Performance testing
bench:
	@echo "Running benchmarks..."
	$(GO) test ./... -bench=. -benchmem

# Coverage report with minimum threshold (80%)
coverage:
	@echo "Generating coverage report with minimum 80% threshold..."
	chmod +x scripts/coverage.sh
	./scripts/coverage.sh
	@echo "Coverage report generated: coverage/coverage.html"

# Quick test coverage (no threshold check)
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test ./... -coverprofile=coverage.out -covermode=atomic
	$(GO) tool cover -func=coverage.out
	@echo "Coverage report generated: coverage.out"

# Coverage for CI/CD pipeline
ci-coverage:
	@echo "Running CI coverage check..."
	mkdir -p coverage
	$(GO) test ./... -coverprofile=coverage/coverage.out -covermode=atomic
	COVERAGE=$(shell go tool cover -func=coverage/coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	MIN_COVERAGE=80; \
	if [ $$(echo "$$COVERAGE < $$MIN_COVERAGE" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $$COVERAGE% is below minimum $$MIN_COVERAGE%"; \
		exit 1; \
	else \
		echo "✅ Coverage $$COVERAGE% meets minimum $$MIN_COVERAGE%"; \
	fi

# Security scanning
security:
	@echo "Running security scan..."
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec ./...

# Documentation
docs:
	@echo "Generating documentation..."
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	# Add swagger generation commands here
	@echo "See README.md for detailed documentation"

# Release build
release: clean test
	@echo "Building release binaries..."
	
	# Linux
	GOOS=linux GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/workflow-api
	
	# macOS
	GOOS=darwin GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/workflow-api
	GOOS=darwin GOARCH=arm64 $(GO) build -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/workflow-api
	
	# Windows
	GOOS=windows GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/workflow-api
	
	@echo "Release binaries created in bin/"
