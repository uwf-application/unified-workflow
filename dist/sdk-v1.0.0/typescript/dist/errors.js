/**
 * Error handling for Unified Workflow SDK
 */
export class WorkflowSDKError extends Error {
    constructor(code, message, details, field, originalError) {
        super(message);
        this.name = 'WorkflowSDKError';
        this.code = code;
        this.details = details;
        this.field = field;
        this.originalError = originalError;
        // Ensure proper prototype chain
        Object.setPrototypeOf(this, WorkflowSDKError.prototype);
    }
    static fromError(error, code = 'UNKNOWN_ERROR') {
        return new WorkflowSDKError(code, error.message, { originalMessage: error.message }, undefined, error);
    }
}
export class RequestValidationError extends WorkflowSDKError {
    constructor(message, validationResult, details) {
        super('VALIDATION_FAILED', message, details);
        this.name = 'RequestValidationError';
        this.validationResult = validationResult;
        Object.setPrototypeOf(this, RequestValidationError.prototype);
    }
}
export class BatchExecutionError extends WorkflowSDKError {
    constructor(batchId, message, failedExecutions, errors, details) {
        super('BATCH_EXECUTION_FAILED', message, details);
        this.name = 'BatchExecutionError';
        this.batchId = batchId;
        this.failedExecutions = failedExecutions;
        this.errors = errors;
        Object.setPrototypeOf(this, BatchExecutionError.prototype);
    }
}
export class WebhookError extends WorkflowSDKError {
    constructor(message, webhookId, details) {
        super('WEBHOOK_ERROR', message, details);
        this.name = 'WebhookError';
        this.webhookId = webhookId;
        Object.setPrototypeOf(this, WebhookError.prototype);
    }
}
export class SDKConfigurationError extends WorkflowSDKError {
    constructor(message, field, details) {
        super('INVALID_CONFIGURATION', message, details, field);
        this.name = 'SDKConfigurationError';
        Object.setPrototypeOf(this, SDKConfigurationError.prototype);
    }
}
export class NetworkError extends WorkflowSDKError {
    constructor(message, originalError, details) {
        super('NETWORK_ERROR', message, details, undefined, originalError);
        this.name = 'NetworkError';
        Object.setPrototypeOf(this, NetworkError.prototype);
    }
}
export class AuthenticationError extends WorkflowSDKError {
    constructor(message, details) {
        super('AUTHENTICATION_FAILED', message, details);
        this.name = 'AuthenticationError';
        Object.setPrototypeOf(this, AuthenticationError.prototype);
    }
}
export class WorkflowNotFoundError extends WorkflowSDKError {
    constructor(workflowId, message = `Workflow not found: ${workflowId}`, details) {
        super('WORKFLOW_NOT_FOUND', message, details);
        this.name = 'WorkflowNotFoundError';
        this.workflowId = workflowId;
        Object.setPrototypeOf(this, WorkflowNotFoundError.prototype);
    }
}
export class ExecutionNotFoundError extends WorkflowSDKError {
    constructor(runId, message = `Execution not found: ${runId}`, details) {
        super('EXECUTION_NOT_FOUND', message, details);
        this.name = 'ExecutionNotFoundError';
        this.runId = runId;
        Object.setPrototypeOf(this, ExecutionNotFoundError.prototype);
    }
}
export class TimeoutError extends WorkflowSDKError {
    constructor(timeoutMs, message = `Operation timed out after ${timeoutMs}ms`, details) {
        super('TIMEOUT', message, details);
        this.name = 'TimeoutError';
        this.timeoutMs = timeoutMs;
        Object.setPrototypeOf(this, TimeoutError.prototype);
    }
}
export class RateLimitError extends WorkflowSDKError {
    constructor(message, retryAfter, details) {
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
};
// Helper function to create SDK errors
export function createSDKError(code, message, details, field, originalError) {
    return new WorkflowSDKError(code, message, details, field, originalError);
}
// Helper function to check if an error is an SDK error
export function isSDKError(error) {
    return error instanceof WorkflowSDKError ||
        (error && typeof error === 'object' && 'code' in error && 'message' in error);
}
// Helper function to wrap unknown errors
export function wrapError(error) {
    if (isSDKError(error)) {
        return error;
    }
    if (error instanceof Error) {
        return WorkflowSDKError.fromError(error);
    }
    return new WorkflowSDKError(ErrorCodes.UNKNOWN_ERROR, typeof error === 'string' ? error : 'An unknown error occurred', typeof error === 'object' ? error : { original: error });
}
//# sourceMappingURL=errors.js.map