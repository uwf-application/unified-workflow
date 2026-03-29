/**
 * Error handling for Unified Workflow SDK
 */
import { SDKError, ValidationResult, BatchExecutionResult, BatchError } from './models';
export declare class WorkflowSDKError extends Error implements SDKError {
    code: string;
    details?: Record<string, any>;
    field?: string;
    originalError?: Error;
    constructor(code: string, message: string, details?: Record<string, any>, field?: string, originalError?: Error);
    static fromError(error: Error, code?: string): WorkflowSDKError;
}
export declare class RequestValidationError extends WorkflowSDKError {
    validationResult: ValidationResult;
    constructor(message: string, validationResult: ValidationResult, details?: Record<string, any>);
}
export declare class BatchExecutionError extends WorkflowSDKError {
    batchId: string;
    failedExecutions?: BatchExecutionResult[];
    errors?: BatchError[];
    constructor(batchId: string, message: string, failedExecutions?: BatchExecutionResult[], errors?: BatchError[], details?: Record<string, any>);
}
export declare class WebhookError extends WorkflowSDKError {
    webhookId?: string;
    constructor(message: string, webhookId?: string, details?: Record<string, any>);
}
export declare class SDKConfigurationError extends WorkflowSDKError {
    constructor(message: string, field?: string, details?: Record<string, any>);
}
export declare class NetworkError extends WorkflowSDKError {
    constructor(message: string, originalError?: Error, details?: Record<string, any>);
}
export declare class AuthenticationError extends WorkflowSDKError {
    constructor(message: string, details?: Record<string, any>);
}
export declare class WorkflowNotFoundError extends WorkflowSDKError {
    workflowId: string;
    constructor(workflowId: string, message?: string, details?: Record<string, any>);
}
export declare class ExecutionNotFoundError extends WorkflowSDKError {
    runId: string;
    constructor(runId: string, message?: string, details?: Record<string, any>);
}
export declare class TimeoutError extends WorkflowSDKError {
    timeoutMs: number;
    constructor(timeoutMs: number, message?: string, details?: Record<string, any>);
}
export declare class RateLimitError extends WorkflowSDKError {
    retryAfter?: number;
    constructor(message: string, retryAfter?: number, details?: Record<string, any>);
}
export declare const ErrorCodes: {
    readonly VALIDATION_FAILED: "VALIDATION_FAILED";
    readonly REQUEST_PARSING_FAILED: "REQUEST_PARSING_FAILED";
    readonly WORKFLOW_EXECUTION_FAILED: "WORKFLOW_EXECUTION_FAILED";
    readonly BATCH_EXECUTION_FAILED: "BATCH_EXECUTION_FAILED";
    readonly WORKFLOW_NOT_FOUND: "WORKFLOW_NOT_FOUND";
    readonly EXECUTION_NOT_FOUND: "EXECUTION_NOT_FOUND";
    readonly WEBHOOK_NOT_FOUND: "WEBHOOK_NOT_FOUND";
    readonly INVALID_CONFIGURATION: "INVALID_CONFIGURATION";
    readonly AUTHENTICATION_FAILED: "AUTHENTICATION_FAILED";
    readonly NETWORK_ERROR: "NETWORK_ERROR";
    readonly TIMEOUT: "TIMEOUT";
    readonly RATE_LIMITED: "RATE_LIMITED";
    readonly WEBHOOK_ERROR: "WEBHOOK_ERROR";
    readonly UNKNOWN_ERROR: "UNKNOWN_ERROR";
    readonly INTERNAL_SERVER_ERROR: "INTERNAL_SERVER_ERROR";
};
export declare function createSDKError(code: keyof typeof ErrorCodes, message: string, details?: Record<string, any>, field?: string, originalError?: Error): WorkflowSDKError;
export declare function isSDKError(error: any): error is WorkflowSDKError;
export declare function wrapError(error: any): WorkflowSDKError;
//# sourceMappingURL=errors.d.ts.map