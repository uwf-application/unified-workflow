/**
 * Unified Workflow SDK - TypeScript/JavaScript Client
 *
 * This SDK provides a comprehensive interface for interacting with the
 * Unified Workflow Execution Platform from TypeScript/JavaScript applications.
 */
export * from './models';
export * from './errors';
export * from './config';
export * from './client';
export { UnifiedWorkflowSDK, WorkflowSDKClient } from './client';
export { SDKConfig, createConfig, validateConfig, DEFAULT_CONFIG } from './config';
export { WorkflowSDKError, RequestValidationError, BatchExecutionError, WebhookError, SDKConfigurationError, NetworkError, AuthenticationError, WorkflowNotFoundError, ExecutionNotFoundError, TimeoutError, RateLimitError, ErrorCodes, createSDKError, isSDKError, wrapError } from './errors';
import { UnifiedWorkflowSDK } from './client';
export default UnifiedWorkflowSDK;
export declare function createSDK(config: Partial<import('./models').SDKConfiguration>): UnifiedWorkflowSDK;
export declare const VERSION = "1.0.0";
export declare const SDK_NAME = "@unified-workflow/sdk";
//# sourceMappingURL=index.d.ts.map