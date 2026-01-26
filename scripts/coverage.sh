#!/bin/bash

echo "=== Running Test Coverage Analysis ==="
echo ""

# Create coverage directory
mkdir -p coverage

# Run tests with coverage for all packages
echo "Running tests with coverage..."
go test ./... -coverprofile=coverage/coverage.out -covermode=atomic

if [ $? -ne 0 ]; then
    echo "❌ Tests failed"
    exit 1
fi

echo ""
echo "✅ Tests passed"

# Generate HTML coverage report
echo "Generating HTML coverage report..."
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Generate coverage summary
echo "Generating coverage summary..."
go tool cover -func=coverage/coverage.out > coverage/coverage.txt

# Display coverage summary
echo ""
echo "=== Coverage Summary ==="
cat coverage/coverage.txt | tail -5

# Check if coverage meets minimum threshold (80%)
COVERAGE=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}' | sed 's/%//')
MIN_COVERAGE=80

echo ""
echo "Current coverage: ${COVERAGE}%"
echo "Minimum required: ${MIN_COVERAGE}%"

if (( $(echo "$COVERAGE < $MIN_COVERAGE" | bc -l) )); then
    echo "❌ Coverage below minimum threshold"
    echo "Please add more tests to improve coverage"
    exit 1
else
    echo "✅ Coverage meets minimum threshold"
fi

echo ""
echo "Coverage report available at: coverage/coverage.html"
echo "Coverage summary available at: coverage/coverage.txt"
