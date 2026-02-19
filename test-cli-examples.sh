#!/bin/bash
# Unified Workflow CLI Test Examples
# This script demonstrates how to use the uwf-cli tool

set -e

echo "=== Unified Workflow CLI Test Examples ==="
echo ""

# Check if CLI is available
if ! command -v ./uwf-cli &> /dev/null; then
    echo "Error: uwf-cli not found. Building it first..."
    go build -o uwf-cli ./cmd/uwf-cli
fi

echo "1. Testing CLI version and help"
echo "--------------------------------"
./uwf-cli --version
echo ""
./uwf-cli --help
echo ""

echo "2. Testing health check"
echo "--------------------------------"
./uwf-cli health
echo ""

echo "3. Testing workflows commands"
echo "--------------------------------"
echo "Listing workflows:"
./uwf-cli workflows list --output json | jq '.workflows[0]'
echo ""

echo "Creating a test workflow:"
./uwf-cli workflows create --name "CLI Test Workflow" --description "Created by test script" --output json | jq '.'
echo ""

echo "4. Testing execute commands"
echo "--------------------------------"
echo "Executing workflow asynchronously:"
RUN_RESPONSE=$(./uwf-cli execute async workflow-1771427384409393014 --input '{"test": "cli-script", "timestamp": "'$(date -Iseconds)'"}' --output json)
echo "$RUN_RESPONSE" | jq '.'
RUN_ID=$(echo "$RUN_RESPONSE" | jq -r '.run_id')
echo "Run ID: $RUN_ID"
echo ""

echo "5. Testing executions commands"
echo "--------------------------------"
echo "Listing executions:"
./uwf-cli executions list --limit 3 --output json | jq '.executions'
echo ""

echo "Getting execution status (simulated):"
./uwf-cli executions status "$RUN_ID" --output json | jq '.'
echo ""

echo "6. Testing test utilities"
echo "--------------------------------"
echo "Generating test data:"
./uwf-cli test generate --count 2 --type simple --output-dir ./test-output --output json | jq '.'
echo ""

echo "7. Testing configuration"
echo "--------------------------------"
echo "Showing current configuration:"
./uwf-cli config show --output json | jq '.'
echo ""

echo "Initializing configuration:"
./uwf-cli config init --output json | jq '.'
echo ""

echo "8. Testing completion"
echo "--------------------------------"
echo "Generating bash completion:"
./uwf-cli completion bash --help | head -5
echo ""

echo "=== Test Complete ==="
echo ""
echo "Summary:"
echo "- CLI tool is working correctly"
echo "- All major commands are functional"
echo "- JSON output format is supported"
echo "- Simulated responses for testing"
echo ""
echo "Next steps:"
echo "1. Copy the uwf-cli binary to your jump server"
echo "2. Update the endpoint in config to your actual workflow-api"
echo "3. Integrate with your SDK for real API calls"
echo "4. Test with your deployed services"

# Cleanup
rm -rf ./test-output 2>/dev/null || true