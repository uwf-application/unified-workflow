/**
 * Error handling for Unified Workflow SDK
 */

import { SDKError, ValidationResult, BatchExecutionResult, BatchError } from './models';

export class WorkflowSDKError extends Error implements SDKError {
  code: string;
  details?: Record<string, any>;
  field?: string;
  originalError?: Error;

  constructor(
    code: string,
    message: string,
    details?: Record<string, any>,
    field?: string,
    originalError?: Error
  ) {
    super(message);
    this.name = 'WorkflowSDKError';
    this.code = code;
    this.details = details;
    this.field = field;
    this.originalError = originalError;
    
    // Ensure proper prototype chain
    Object.setPrototypeOf(this, WorkflowSDKError.prototype);
  }

  static fromError(error: Error, code: string = 'UNKNOWN_ERROR'): WorkflowSDKError {
    return new WorkflowSDKError(
      code,
      error.message,
      { originalMessage: error.message },
      undefined,
      error
    );
  }
}

export class RequestValidationError extends WorkflowSDKError {
  validationResult: ValidationResult;

  constructor(
    message: string,
    validationResult: ValidationResult,
    details?: Record<string, any>
  ) {
    super('VALIDATION_FAILED', message, details);
    this.name = 'RequestValidationError';
    this.validationResult = validationResult;
    
    Object.setPrototypeOf(this, RequestValidationError.prototype);
  }
}

export class BatchExecutionError extends WorkflowSDKError {
  batchId: string;
  failedExecutions?: BatchExecutionResult[];
  errors?: BatchError[];

  constructor(
    batchId: string,
    message: string,
    failedExecutions?: BatchExecutionResult[],
    errors?: BatchError[],
    details?: Record<string, any>
  ) {
    super('BATCH_EXECUTION_FAILED', message, details);
    this.name = 'BatchExecutionError';
    this.batchId = batchId;
    this.failedExecutions = failedExecutions;
    this.errors = errors;
    
    Object.setPrototypeOf(this, BatchExecutionError.prototype);
  }
}

export class WebhookError extends WorkflowSDKError {
  webhookId?: string;

  constructor(
    message: string,
    webhookId?: string,
    details?: Record<string, any>
  ) {
    super('WEBHOOK_ERROR', message, details);
    this.name = 'WebhookError';
    this.webhookId = webhookId;
    
    Object.setPrototypeOf(this, WebhookError.prototype);
  }
}

export class SDKConfigurationError extends WorkflowSDKError {
  constructor(
    message: string,
    field?: string,
    details?: Record<string, any>
  ) {
    super('INVALID_CONFIGURATION', message, details, field);
    this.name = 'SDKConfigurationError';
    
    Object.setPrototypeOf(this, SDKConfigurationError.prototype);
  }
}

export class NetworkError extends WorkflowSDKError {
  constructor(
    message: string,
    originalError?: Error,
    details?: Record<string, any>
  ) {
    super('NETWORK_ERROR', message, details, undefined, originalError);
    this.name = 'NetworkError';
    
    Object.setPrototypeOf(this, NetworkError.prototype);
  }
}

export class AuthenticationError extends WorkflowSDKError {
  constructor(
    message: string,
    details?: Record<string, any>
  ) {
    super('AUTHENTICATION_FAILED', message, details);
    this.name = 'AuthenticationError';
    
    Object.setPrototypeOf(this, AuthenticationError.prototype);
  }
}

export class WorkflowNotFoundError extends WorkflowSDKError {
  workflowId: string;

  constructor(
    workflowId: string,
    message: string = `Workflow not found: ${workflowId}`,
    details?: Record<string, any>
  ) {
    super('WORKFLOW_NOT_FOUND', message, details);
    this.name = 'WorkflowNotFoundError';
    this.workflowId = workflowId;
    
    Object.setPrototypeOf(this, WorkflowNotFoundError.prototype);
  }
}

export class ExecutionNotFoundError extends WorkflowSDKError {
  runId: string;

  constructor(
    runId: string,
    message: string = `Execution not found: ${runId}`,
    details?: Record<string, any>
  ) {
    super('EXECUTION_NOT_FOUND', message, details);
    this.name = 'ExecutionNotFoundError';
    this.runId = runId;
    
    Object.setPrototypeOf(this, ExecutionNotFoundError.prototype);
  }
}

export class TimeoutError extends WorkflowSDKError {
  timeoutMs: number;

  constructor(
    timeoutMs: number,
    message: string = `Operation timed out after ${timeoutMs}ms`,
    details?: Record<string, any>
  ) {
    super('TIMEOUT', message, details);
    this.name = 'TimeoutError';
    this.timeoutMs = timeoutMs;
    
    Object.setPrototypeOf(this, TimeoutError.prototype);
  }
}

export class RateLimitError extends WorkflowSDKError {
  retryAfter?: number;

  constructor(
    message: string,
    retryAfter?: number,
    details?: Record<string, any>
  ) {
    super('RATE_LIMITED', message, details);
    this.name = 'RateLimitError';
    this.retryAfter = retryAfter;
    
    Object.setPrototypeOf(this, RateLimitError.prototype);
  }
}

// Error codes
export const ErrorCodes = {
  // Validation errors
  VALIDATION_FAILED: 'VALIDATION_FAILED',
  REQUEST_PARSING_FAILED: 'REQUEST_PARSING_FAILED',
  
  // Execution errors
  WORKFLOW_EXECUTION_FAILED: 'WORKFLOW_EXECUTION_FAILED',
  BATCH_EXECUTION_FAILED: 'BATCH_EXECUTION_FAILED',
  
  // Resource errors
  WORKFLOW_NOT_FOUND: 'WORKFLOW_NOT_FOUND',
  EXECUTION_NOT_FOUND: 'EXECUTION_NOT_FOUND',
  WEBHOOK_NOT_FOUND: 'WEBHOOK_NOT_FOUND',
  
  // Configuration errors
  INVALID_CONFIGURATION: 'INVALID_CONFIGURATION',
  AUTHENTICATION_FAILED: 'AUTHENTICATION_FAILED',
  
  // Network errors
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT: 'TIMEOUT',
  RATE_LIMITED: 'RATE_LIMITED',
  
  // Webhook errors
  WEBHOOK_ERROR: 'WEBHOOK_ERROR',
  
  // Unknown errors
  UNKNOWN_ERROR: 'UNKNOWN_ERROR',
  INTERNAL_SERVER_ERROR: 'INTERNAL_SERVER_ERROR'
} as const;

// Helper function to create SDK errors
export function createSDKError(
  code: keyof typeof ErrorCodes,
  message: string,
  details?: Record<string, any>,
  field?: string,
  originalError?: Error
): WorkflowSDKError {
  return new WorkflowSDKError(code, message, details, field, originalError);
}

// Helper function to check if an error is an SDK error
export function isSDKError(error: any): error is WorkflowSDKError {
  return error instanceof WorkflowSDKError || 
         (error && typeof error === 'object' && 'code' in error && 'message' in error);
}

// Helper function to wrap unknown errors
export function wrapError(error: any): WorkflowSDKError {
  if (isSDKError(error)) {
    return error;
  }
  
  if (error instanceof Error) {
    return WorkflowSDKError.fromError(error);
  }
  
  return new WorkflowSDKError(
    ErrorCodes.UNKNOWN_ERROR,
    typeof error === 'string' ? error : 'An unknown error occurred',
    typeof error === 'object' ? error : { original: error }
  );
}