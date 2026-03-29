/**
 * Basic usage example for Unified Workflow SDK
 */

import { createSDK, ValidationRuleType } from '../src/index';

async function main() {
  console.log('=== Unified Workflow SDK Basic Example ===\n');

  // 1. Create SDK instance
  console.log('1. Creating SDK instance...');
  const sdk = createSDK({
    workflowApiEndpoint: 'http://localhost:8080',
    timeoutMs: 30000,
    maxRetries: 3,
    authToken: 'your-auth-token',
    enableValidation: true,
    enableSanitization: true,
  });

  try {
    // 2. Check service health
    console.log('2. Checking service health...');
    try {
      await sdk.ping();
      console.log('   ✅ Service is reachable\n');
    } catch (error) {
      console.log('   ⚠️  Service ping failed:', error.message);
      console.log('   Continuing with example...\n');
    }

    // 3. List available workflows
    console.log('3. Listing available workflows...');
    try {
      const workflows = await sdk.listWorkflows();
      console.log(`   Found ${workflows.length || 0} workflows\n`);
    } catch (error) {
      console.log('   ⚠️  Failed to list workflows:', error.message);
      console.log('   Continuing with example...\n');
    }

    // 4. Execute a workflow with validation
    console.log('4. Executing workflow with validation...');
    
    const workflowId = 'payment-processing-workflow';
    const inputData = {
      user_id: 'user_12345',
      amount: 99.99,
      email: 'user@example.com',
      currency: 'USD',
      payment_method: 'credit_card',
    };

    // Define validation rules
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
        field: 'currency',
        ruleType: ValidationRuleType.STRING,
        required: true,
        allowedValues: ['USD', 'EUR', 'GBP'],
      },
    ];

    try {
      const response = await sdk.validateAndExecuteWorkflow(
        workflowId,
        inputData,
        validationRules
      );

      console.log('   ✅ Workflow execution started!\n');
      console.log('   Execution Details:');
      console.log(`     Run ID: ${response.runId}`);
      console.log(`     Status: ${response.status}`);
      console.log(`     Status URL: ${response.statusUrl}`);
      console.log(`     Result URL: ${response.resultUrl}`);
      console.log(`     Poll after: ${response.pollAfterMs || 1000}ms`);
      console.log(`     Estimated completion: ${response.estimatedCompletionMs || 'N/A'}ms\n`);

      // 5. Check execution status
      console.log('5. Checking execution status...');
      
      // Wait a bit before checking status
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      try {
        const status = await sdk.getExecutionStatus(response.runId);
        console.log(`   Current Status: ${status.status}`);
        console.log(`   Progress: ${status.progress || 0}%`);
        console.log(`   Started at: ${status.startedAt || 'N/A'}`);
        console.log(`   Updated at: ${status.updatedAt || 'N/A'}\n`);
      } catch (statusError) {
        console.log('   ⚠️  Failed to get execution status:', statusError.message);
        console.log('   Continuing with example...\n');
      }

      // 6. Demonstrate batch execution
      console.log('6. Demonstrating batch execution...');
      
      const batchRequest = {
        executions: [
          {
            workflowId: 'workflow-1',
            request: {
              inputData: { test: 'batch-1', value: 100 },
            },
            priority: 1,
          },
          {
            workflowId: 'workflow-2',
            request: {
              inputData: { test: 'batch-2', value: 200 },
            },
            priority: 2,
          },
        ],
        parallel: true,
        maxConcurrent: 2,
        stopOnFirstFailure: false,
      };

      try {
        const batchResponse = await sdk.batchExecuteWorkflows(batchRequest);
        console.log(`   ✅ Batch execution started (ID: ${batchResponse.batchId})`);
        console.log(`     Total: ${batchResponse.total}`);
        console.log(`     Successful: ${batchResponse.successful}`);
        console.log(`     Failed: ${batchResponse.failed}`);
        console.log(`     Pending: ${batchResponse.pending}\n`);
      } catch (batchError) {
        console.log('   ⚠️  Batch execution failed:', batchError.message);
        console.log('   Continuing with example...\n');
      }

    } catch (executionError) {
      console.log('   ❌ Workflow execution failed:');
      
      if (executionError.code === 'VALIDATION_FAILED') {
        console.log('   Validation errors:');
        executionError.validationResult.errors?.forEach((error: any) => {
          console.log(`     - ${error.field}: ${error.message}`);
        });
      } else {
        console.log(`   Error: ${executionError.message}`);
        console.log(`   Code: ${executionError.code}`);
      }
      console.log();
    }

    // 7. Demonstrate error handling
    console.log('7. Demonstrating error handling...');
    
    try {
      // Try to execute a non-existent workflow
      await sdk.executeWorkflow('non-existent-workflow', { test: true });
    } catch (error) {
      if (error.code === 'WORKFLOW_NOT_FOUND') {
        console.log(`   ✅ Correctly handled workflow not found error`);
        console.log(`     Workflow ID: ${error.workflowId}`);
        console.log(`     Message: ${error.message}\n`);
      } else {
        console.log(`   ⚠️  Unexpected error: ${error.message}\n`);
      }
    }

  } finally {
    // 8. Always close the SDK
    console.log('8. Closing SDK...');
    await sdk.close();
    console.log('   ✅ SDK closed successfully\n');
  }

  console.log('=== Example Complete ===');
  console.log('\nSummary:');
  console.log('• Created and configured SDK instance');
  console.log('• Checked service health');
  console.log('• Executed workflow with validation');
  console.log('• Monitored execution status');
  console.log('• Demonstrated batch execution');
  console.log('• Handled errors appropriately');
  console.log('• Properly closed SDK connection');
}

// Run the example
main().catch(error => {
  console.error('Example failed:', error);
  process.exit(1);
});