/**
 * Unified Workflow SDK - TypeScript/JavaScript Client
 * 
 * This SDK provides a comprehensive interface for interacting with the
 * Unified Workflow Execution Platform from TypeScript/JavaScript applications.
 */

// Export all models
export * from './models';

// Export all errors
export * from './errors';

// Export configuration
export * from './config';

// Export main client
export * from './client';

// Re-export commonly used types for convenience
export {
  UnifiedWorkflowSDK,
  WorkflowSDKClient
} from './client';

export {
  SDKConfig,
  createConfig,
  validateConfig,
  DEFAULT_CONFIG
} from './config';

export {
  WorkflowSDKError,
  RequestValidationError,
  BatchExecutionError,
  WebhookError,
  SDKConfigurationError,
  NetworkError,
  AuthenticationError,
  WorkflowNotFoundError,
  ExecutionNotFoundError,
  TimeoutError,
  RateLimitError,
  ErrorCodes,
  createSDKError,
  isSDKError,
  wrapError
} from './errors';

// Default export for easier imports
import { UnifiedWorkflowSDK } from './client';
export default UnifiedWorkflowSDK;

// Helper function to create a new SDK instance
export function createSDK(config: Partial<import('./models').SDKConfiguration>) {
  return new UnifiedWorkflowSDK(config);
}

// Version information
export const VERSION = '1.0.0';
export const SDK_NAME = '@unified-workflow/sdk';