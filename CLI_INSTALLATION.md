# Unified Workflow CLI (uwf-cli) Installation Guide

A comprehensive command-line interface for testing and managing your workflow-api from your jump server.

## Quick Installation

### Option 1: Use the Linux Binary (Recommended for Jump Server)

For AlmaLinux 9.4 (x86_64):

```bash
# Download the Linux binary from your local machine or build it
# On your local machine, build for Linux:
# GOOS=linux GOARCH=amd64 go build -o uwf-cli-linux ./cmd/uwf-cli

# Copy to jump server
scp uwf-cli-linux user@jump-server:/tmp/uwf-cli

# On jump server:
sudo mv /tmp/uwf-cli /usr/local/bin/uwf-cli
sudo chmod +x /usr/local/bin/uwf-cli

# Verify installation
uwf-cli --version
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/uwf-application/unified-workflow.git
cd unified-workflow

# Build the CLI
go build -o uwf-cli ./cmd/uwf-cli

# Install to system path
sudo mv uwf-cli /usr/local/bin/
```

### Option 3: Docker (Containerized)

```bash
# Pull the image
docker pull ghcr.io/uwf-application/uwf-cli:latest

# Run as a container
docker run --rm ghcr.io/uwf-application/uwf-cli:latest --help

# Create an alias for convenience
alias uwf-cli='docker run --rm -v $(pwd):/data ghcr.io/uwf-application/uwf-cli:latest'
```

## Configuration

### Quick Setup

```bash
# Initialize configuration
uwf-cli config init

# Set your endpoint (default: https://af-test.qazpost.kz)
uwf-cli config set endpoint https://af-test.qazpost.kz

# Set authentication token (if needed)
uwf-cli config set auth_token "your-token-here"
```

### Configuration File

The CLI uses `~/.uwf-cli/config.yaml`:

```yaml
# Unified Workflow CLI Configuration
endpoint: "https://af-test.qazpost.kz"
auth_token: ""
timeout: 30
max_retries: 3
output_format: "json"
default_workflow_name: "Test Workflow"
test_data_path: "./test-data"
```

### Environment Variables

```bash
# Override config with environment variables
export WORKFLOW_API_ENDPOINT="https://af-test.qazpost.kz"
export WORKFLOW_AUTH_TOKEN="your-token"
export UWF_CLI_OUTPUT="json"
```

## Quick Start Examples

### 1. Check Service Health

```bash
uwf-cli health
uwf-cli health --verbose
```

### 2. List Workflows

```bash
# List all workflows
uwf-cli workflows list

# List with JSON output
uwf-cli workflows list --output json

# List with pagination
uwf-cli workflows list --limit 10 --offset 0
```

### 3. Create and Execute Workflows

```bash
# Create a workflow
uwf-cli workflows create --name "Payment Processing" --description "Process payments"

# Execute workflow synchronously
uwf-cli execute sync workflow-123 --input '{"amount": 100, "currency": "USD"}'

# Execute workflow asynchronously
uwf-cli execute async workflow-123 --input '{"amount": 100, "currency": "USD"}'

# Execute with input file
echo '{"amount": 100, "currency": "USD"}' > input.json
uwf-cli execute async workflow-123 --input-file input.json
```

### 4. Monitor Executions

```bash
# List all executions
uwf-cli executions list

# Get execution status
uwf-cli executions status run-1771433588043084557

# Get execution result
uwf-cli executions result run-1771433588043084557

# Watch execution in real-time
uwf-cli executions watch run-1771433588043084557 --interval 2s
```

### 5. Test Utilities

```bash
# Generate test data
uwf-cli test generate --count 5 --type payment

# Run test suite
uwf-cli test run --suite basic

# Cleanup test data
uwf-cli test cleanup --confirm
```

## Complete Command Reference

### Global Flags

```
-t, --auth-token string   Authentication token
-c, --config string       Config file path
-e, --endpoint string     Workflow API endpoint (default "https://af-test.qazpost.kz")
-h, --help                help for uwf-cli
-o, --output string       Output format (json, yaml, table) (default "json")
-v, --verbose             Verbose output
    --version             version for uwf-cli
```

### Available Commands

```
completion  Generate shell completion
config      Manage CLI configuration
execute     Execute workflows
executions  Manage workflow executions
health      Check API health
test        Test utilities
workflows   Manage workflows
```

### Workflows Commands

```
workflows list              List all workflows
workflows create            Create a new workflow
workflows get <id>          Get workflow details
workflows delete <id>       Delete a workflow
workflows export <id>       Export workflow to file
workflows import <file>     Import workflow from file
workflows bulk-create       Create multiple workflows from file
workflows bulk-delete       Delete multiple workflows from file
```

### Execute Commands

```
execute sync <id>           Execute workflow synchronously
execute async <id>          Execute workflow asynchronously
execute bulk <file>         Execute multiple workflows from file
```

### Executions Commands

```
executions list             List all executions
executions status <id>      Get execution status
executions result <id>      Get execution result
executions data <id>        Get execution data
executions metrics <id>     Get execution metrics
executions cancel <id>      Cancel execution
executions pause <id>       Pause execution
executions resume <id>      Resume execution
executions retry <id>       Retry execution
executions watch <id>       Watch execution status in real-time
```

## Advanced Usage

### Bulk Operations

```bash
# Create multiple workflows from JSON file
uwf-cli workflows bulk-create workflows.json

# Execute multiple workflows
uwf-cli execute bulk executions.json --parallel --concurrency 5

# Delete multiple workflows
uwf-cli workflows bulk-delete workflows-to-delete.json
```

### Shell Completion

```bash
# Bash
uwf-cli completion bash > /etc/bash_completion.d/uwf-cli

# Zsh
uwf-cli completion zsh > "${fpath[1]}/_uwf-cli"

# Fish
uwf-cli completion fish > ~/.config/fish/completions/uwf-cli.fish
```

### Scripting Examples

```bash
#!/bin/bash
# Example script to test workflow execution

# Set endpoint
export WORKFLOW_API_ENDPOINT="https://af-test.qazpost.kz"

# Create test workflow
WORKFLOW_ID=$(uwf-cli workflows create --name "Test Script" --description "Created by script" --output json | jq -r '.id')

echo "Created workflow: $WORKFLOW_ID"

# Execute workflow
RUN_ID=$(uwf-cli execute async "$WORKFLOW_ID" --input '{"test": "script"}' --output json | jq -r '.run_id')

echo "Started execution: $RUN_ID"

# Wait for completion
uwf-cli executions watch "$RUN_ID" --interval 5s --timeout 5m

# Get result
uwf-cli executions result "$RUN_ID" --output json | jq '.result'
```

## Troubleshooting

### Common Issues

1. **Connection refused**: Check if the endpoint is correct and the service is running
   ```bash
   curl -v https://af-test.qazpost.kz/health
   ```

2. **Authentication errors**: Verify your auth token
   ```bash
   uwf-cli config show
   ```

3. **JSON parsing errors**: Ensure input JSON is valid
   ```bash
   echo '{"test": "data"}' | jq .
   ```

### Debug Mode

```bash
# Enable verbose output
uwf-cli --verbose health

# See raw HTTP requests (when integrated with real API)
uwf-cli --verbose workflows list
```

## Integration with Your SDK

The CLI currently uses simulated responses. To integrate with your actual workflow-api:

1. Update the CLI to use your `pkg/client/sdk` package
2. Replace simulated responses with actual API calls
3. Add proper error handling and retry logic

Example integration point in `health.go`:
```go
// Replace simulated check with actual SDK call
config := &sdk.SDKConfig{
    WorkflowAPIEndpoint: endpoint,
    Timeout: 10 * time.Second,
}
client, err := sdk.NewClient(config)
if err != nil {
    return err
}
defer client.Close()

err = client.Ping(ctx)
if err != nil {
    return fmt.Errorf("health check failed: %v", err)
}
```

## Next Steps

1. **Test with your deployed services**: Use the CLI to test your workflow-api
2. **Integrate with SDK**: Replace simulated responses with actual API calls
3. **Add advanced features**: Implement missing commands and options
4. **Create CI/CD pipeline**: Automate CLI builds and releases

## Support

- **Issues**: GitHub repository issues
- **Documentation**: `uwf-cli --help` for command reference
- **Examples**: See `examples/` directory for usage patterns