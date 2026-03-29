# Unified Workflow SDK - TypeScript/JavaScript

A comprehensive TypeScript/JavaScript SDK for interacting with the Unified Workflow Execution Platform. This SDK provides a clean, type-safe interface for executing workflows, managing executions, and handling HTTP request parsing with validation.

## Features

- **TypeScript First**: Full TypeScript support with comprehensive type definitions
- **HTTP Request Parsing**: Automatically parse HTTP requests into workflow execution requests
- **Data Validation**: Built-in validation with customizable rules
- **Request Sanitization**: Clean and sanitize input data
- **Session Extraction**: Extract user sessions from HTTP requests
- **Security Context**: Pass security information to workflows
- **Error Handling**: Comprehensive error handling with retry mechanisms
- **Async Execution**: Support for asynchronous workflow execution
- **Execution Management**: Monitor and manage workflow executions
- **Batch Operations**: Execute multiple workflows in batch
- **Webhook Support**: Register and manage webhooks for event notifications

## Installation

### npm
```bash
npm install @unified-workflow/sdk
```

### yarn
```bash
yarn add @unified-workflow/sdk
```

### pnpm
```bash
pnpm add @unified-workflow/sdk
```

## Quick Start

### Basic Usage

```typescript
import { createSDK } from '@unified-workflow/sdk';

// Create SDK instance
const sdk = createSDK({
  workflowApiEndpoint: 'http://localhost:8080',
  timeoutMs: 30000,
  maxRetries: 3,
  authToken: 'your-auth-token',
  enableValidation: true,
  enableSanitization: true,
});

// Execute a workflow
async function executeWorkflow() {
  try {
    const workflowId = 'payment-processing-workflow';
    const inputData = {
      user_id: 'user_12345',
      amount: 99.99,
      email: 'user@example.com',
    };

    const response = await sdk.executeWorkflow(workflowId, inputData);
    
    console.log('Workflow execution started:', response.runId);
    console.log('Status URL:', response.statusUrl);
    console.log('Result URL:', response.resultUrl);
    
  } catch (error) {
    console.error('Failed to execute workflow:', error);
  }
}

// Don't forget to close the SDK when done
sdk.close();
```

### HTTP Request Integration

```typescript
import { createSDK } from '@unified-workflow/sdk';
import express from 'express';

const app = express();
const sdk = createSDK({
  workflowApiEndpoint: 'http://localhost:8080',
  authToken: 'your-auth-token',
});

app.post('/api/workflows/:workflowId/execute', async (req, res) => {
  try {
    const { workflowId } = req.params;
    
    // Execute workflow from HTTP request
    const response = await sdk.executeFromHTTPRequest(workflowId, req);
    
    res.status(202).json({
      success: true,
      runId: response.runId,
      status: response.status,
      statusUrl: response.statusUrl,
      resultUrl: response.resultUrl,
    });
    
  } catch (error) {
    if (error.code === 'VALIDATION_FAILED') {
      res.status(400).json({
        success: false,
        error: 'Validation failed',
        details: error.validationResult,
      });
    } else {
      res.status(500).json({
        success: false,
        error: 'Internal server error',
      });
    }
  }
});

app.listen(3000, () => {
  console.log('Server running on port 3000');
});
```

## Configuration

### SDK Configuration Options

```typescript
import { AuthType, LogLevel } from '@unified-workflow/sdk';

const config = {
  // Required
  workflowApiEndpoint: 'http://localhost:8080',
  sdkVersion: '1.0.0',
  
  // Optional - Network
  timeoutMs: 30000,           // Request timeout in milliseconds
  maxRetries: 3,              // Maximum retry attempts
  retryDelayMs: 1000,         // Delay between retries
  
  // Optional - Authentication
  authType: AuthType.BEARER_TOKEN,
  authToken: 'your-auth-token',
  
  // Optional - Validation
  enableValidation: true,
  enableSanitization: true,
  strictValidation: false,
  defaultValidationRules: [],
  
  // Optional - Context extraction
  enableSessionExtraction: true,
  enableSecurityContext: true,
  includeFullHttpContext: true,
  
  // Optional - Logging
  logLevel: LogLevel.INFO,
  enableRequestLogging: true,
  enableMetrics: true,
  
  // Optional - Custom validators
  customValidators: [],
};
```

## API Reference

### Core Methods

#### `executeWorkflow(workflowId: string, data: Record<string, any>)`
Execute a workflow with raw data.

#### `executeFromHTTPRequest(workflowId: string, request: Request)`
Execute a workflow from an HTTP request, automatically extracting context.

#### `executeWorkflowWithContext(workflowId: string, request: SDKExecuteWorkflowRequest)`
Execute a workflow with full context including HTTP request, session, and security information.

#### `validateAndExecuteWorkflow(workflowId: string, data: Record<string, any>, rules: ValidationRule[])`
Validate data against custom rules before executing the workflow.

### Batch Operations

#### `batchExecuteWorkflows(request: BatchExecuteWorkflowsRequest)`
Execute multiple workflows in a single batch request.

### Webhook Operations

#### `registerWebhook(configuration: WebhookConfiguration)`
Register a webhook for receiving workflow events.

#### `unregisterWebhook(webhookId: string)`
Remove a registered webhook.

#### `listWebhooks()`
List all registered webhooks.

### Workflow Management

#### `listWorkflows()`
List all available workflows.

#### `getWorkflow(workflowId: string)`
Get details of a specific workflow.

#### `createWorkflow(workflow: any)`
Create a new workflow.

#### `deleteWorkflow(workflowId: string)`
Delete a workflow.

### Execution Management

#### `getExecutionStatus(runId: string)`
Get the status of a workflow execution.

#### `getExecutionResult(runId: string)`
Get the result of a completed workflow execution.

#### `cancelExecution(runId: string)`
Cancel a running workflow execution.

#### `pauseExecution(runId: string)`
Pause a running workflow execution.

#### `resumeExecution(runId: string)`
Resume a paused workflow execution.

#### `retryExecution(runId: string)`
Retry a failed workflow execution.

## Error Handling

The SDK provides comprehensive error handling with specific error types:

```typescript
import {
  WorkflowSDKError,
  RequestValidationError,
  WorkflowNotFoundError,
  ExecutionNotFoundError,
  AuthenticationError,
  NetworkError,
  TimeoutError,
  RateLimitError,
  isSDKError,
  wrapError
} from '@unified-workflow/sdk';

try {
  const response = await sdk.executeWorkflow('workflow-id', data);
} catch (error) {
  if (isSDKError(error)) {
    switch (error.code) {
      case 'VALIDATION_FAILED':
        // Handle validation errors
        console.error('Validation failed:', error.validationResult);
        break;
        
      case 'WORKFLOW_NOT_FOUND':
        // Handle workflow not found
        console.error('Workflow not found:', error.workflowId);
        break;
        
      case 'AUTHENTICATION_FAILED':
        // Handle authentication errors
        console.error('Authentication failed');
        break;
        
      case 'NETWORK_ERROR':
        // Handle network errors
        console.error('Network error:', error.originalError);
        break;
        
      case 'TIMEOUT':
        // Handle timeout errors
        console.error('Request timed out');
        break;
        
      case 'RATE_LIMITED':
        // Handle rate limiting
        console.error('Rate limited, retry after:', error.retryAfter);
        break;
        
      default:
        // Handle other SDK errors
        console.error('SDK error:', error);
    }
  } else {
    // Handle non-SDK errors
    console.error('Unexpected error:', error);
  }
}
```

## Validation

The SDK includes a powerful validation system:

```typescript
import { ValidationRuleType } from '@unified-workflow/sdk';

const validationRules = [
  {
    field: 'user_id',
    ruleType: ValidationRuleType.STRING,
    required: true,
    minLength: 5,
    maxLength: 50,
  },
  {
    field: 'amount',
    ruleType: ValidationRuleType.NUMBER,
    required: true,
    minValue: 0.01,
    maxValue: 10000,
  },
  {
    field: 'email',
    ruleType: ValidationRuleType.EMAIL,
    required: true,
  },
  {
    field: 'status',
    ruleType: ValidationRuleType.STRING,
    required: true,
    allowedValues: ['pending', 'approved', 'rejected'],
  },
];

try {
  const response = await sdk.validateAndExecuteWorkflow(
    'payment-workflow',
    inputData,
    validationRules
  );
} catch (error) {
  if (error.code === 'VALIDATION_FAILED') {
    console.error('Validation errors:', error.validationResult.errors);
  }
}
```

## Batch Execution

Execute multiple workflows in a single request:

```typescript
const batchRequest = {
  executions: [
    {
      workflowId: 'workflow-1',
      request: {
        inputData: { user_id: 'user1', amount: 50 },
      },
      priority: 1,
    },
    {
      workflowId: 'workflow-2',
      request: {
        inputData: { user_id: 'user2', amount: 100 },
      },
      priority: 2,
    },
  ],
  parallel: true,
  maxConcurrent: 5,
  stopOnFirstFailure: false,
};

const batchResponse = await sdk.batchExecuteWorkflows(batchRequest);

console.log('Batch ID:', batchResponse.batchId);
console.log('Total executions:', batchResponse.total);
console.log('Successful:', batchResponse.successful);
console.log('Failed:', batchResponse.failed);
```

## Webhooks

Register webhooks to receive workflow events:

```typescript
import { WebhookEvent } from '@unified-workflow/sdk';

const webhookConfig = {
  url: 'https://your-server.com/webhooks/workflow-events',
  events: [
    WebhookEvent.WORKFLOW_STARTED,
    WebhookEvent.WORKFLOW_COMPLETED,
    WebhookEvent.WORKFLOW_FAILED,
  ],
  secret: 'your-webhook-secret',
  enabled: true,
  retryCount: 3,
  timeoutMs: 5000,
};

const webhook = await sdk.registerWebhook(webhookConfig);
console.log('Webhook registered:', webhook.webhookId);

// List all webhooks
const webhooks = await sdk.listWebhooks();
console.log('Registered webhooks:', webhooks.webhooks);

// Unregister webhook
await sdk.unregisterWebhook(webhook.webhookId);
```

## Best Practices

### 1. Always Close the SDK
```typescript
const sdk = createSDK(config);

try {
  // Use the SDK
  await sdk.ping();
  // ... other operations
} finally {
  // Always close the SDK
  await sdk.close();
}
```

### 2. Use Async/Await with Error Handling
```typescript
async function processWorkflow(workflowId: string, data: any) {
  try {
    const response = await sdk.executeWorkflow(workflowId, data);
    return response;
  } catch (error) {
    if (isSDKError(error)) {
      // Handle SDK errors
      throw new Error(`Workflow execution failed: ${error.message}`);
    }
    throw error;
  }
}
```

### 3. Configure Appropriate Timeouts
```typescript
const sdk = createSDK({
  workflowApiEndpoint: 'http://localhost:8080',
  timeoutMs: 60000, // 60 seconds for long-running workflows
  maxRetries: 5,
  retryDelayMs: 2000,
});
```

### 4. Enable Validation for Production
```typescript
const sdk = createSDK({
  workflowApiEndpoint: 'http://localhost:8080',
  enableValidation: true,
  enableSanitization: true,
  strictValidation: true, // Fail fast on validation errors
});
```

### 5. Monitor Execution Status
```typescript
async function monitorExecution(runId: string) {
  let status = 'running';
  
  while (status === 'running' || status === 'pending') {
    const executionStatus = await sdk.getExecutionStatus(runId);
    status = executionStatus.status;
    
    if (status === 'completed') {
      const result = await sdk.getExecutionResult(runId);
      return result;
    } else if (status === 'failed') {
      throw new Error('Workflow execution failed');
    }
    
    // Wait before polling again
    await new Promise(resolve => setTimeout(resolve, 1000));
  }
}
```

## Examples

See the [examples directory](./examples) for complete examples including:

1. Basic workflow execution
2. HTTP server integration
3. Validation and error handling
4. Batch execution
5. Webhook registration
6. Execution monitoring

## Testing

```typescript
import { createSDK } from '@unified-workflow/sdk';

describe('Workflow SDK', () => {
  let sdk: ReturnType<typeof createSDK>;
  
  beforeEach(() => {
    sdk = createSDK({
      workflowApiEndpoint: 'http://localhost:8080',
      timeoutMs: 5000,
    });
  });
  
  afterEach(async () => {
    await sdk.close();
  });
  
  test('should execute workflow', async () => {
    const response = await sdk.executeWorkflow('test-workflow', {
      test: true,
      data: 'test-data',
    });
    
    expect(response.runId).toBeDefined();
    expect(response.statusUrl).toBeDefined();
  });
});
```

## Support

- **Documentation**: Complete API reference and examples
- **GitHub Issues**: Report bugs and request features
- **Email Support**: enterprise-support@your-org.com

## License

MIT License - See LICENSE file for details.

## Contributing

Contributions are welcome! Please see the CONTRIBUTING.md file for details.

## Changelog

See CHANGELOG.md for version history and changes.