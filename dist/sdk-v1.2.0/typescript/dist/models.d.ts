/**
 * TypeScript models for Unified Workflow SDK
 * Based on Smithy API definition: api/workflow-sdk.smithy
 */
export type WorkflowId = string;
export type RunId = string;
export type Timestamp = string | Date;
export declare enum WorkflowStatus {
    DRAFT = "draft",
    ACTIVE = "active",
    INACTIVE = "inactive",
    ARCHIVED = "archived"
}
export declare enum AsyncExecutionStatus {
    PENDING = "pending",
    QUEUED = "queued",
    RUNNING = "running",
    COMPLETED = "completed",
    FAILED = "failed",
    CANCELLED = "cancelled",
    PAUSED = "paused"
}
export declare enum AuthType {
    BEARER_TOKEN = "bearer_token",
    API_KEY = "api_key",
    BASIC_AUTH = "basic_auth",
    OAUTH2 = "oauth2",
    AWS_SIGV4 = "aws_sigv4"
}
export declare enum LogLevel {
    DEBUG = "debug",
    INFO = "info",
    WARN = "warn",
    ERROR = "error"
}
export declare enum ValidationRuleType {
    REQUIRED = "required",
    STRING = "string",
    NUMBER = "number",
    BOOLEAN = "boolean",
    ARRAY = "array",
    OBJECT = "object",
    EMAIL = "email",
    URL = "url",
    UUID = "uuid",
    CUSTOM = "custom"
}
export declare enum WebhookEvent {
    WORKFLOW_STARTED = "workflow_started",
    WORKFLOW_COMPLETED = "workflow_completed",
    WORKFLOW_FAILED = "workflow_failed",
    WORKFLOW_CANCELLED = "workflow_cancelled",
    STEP_COMPLETED = "step_completed",
    STEP_FAILED = "step_failed"
}
export interface HTTPHeaders {
    [key: string]: string[];
}
export interface QueryParams {
    [key: string]: string[];
}
export interface PathParams {
    [key: string]: string;
}
export interface HTTPRequestContext {
    method: string;
    path: string;
    headers: HTTPHeaders;
    queryParams?: QueryParams;
    pathParams?: PathParams;
    body?: any;
    remoteAddr?: string;
    userAgent?: string;
    timestamp: Timestamp;
}
export interface SessionContext {
    userId?: string;
    sessionId?: string;
    roles?: string[];
    permissions?: string[];
    authMethod?: string;
    expiresAt?: Timestamp;
    attributes?: Record<string, any>;
}
export interface GeoLocation {
    country?: string;
    region?: string;
    city?: string;
    latitude?: number;
    longitude?: number;
}
export interface SecurityContext {
    authenticated: boolean;
    authenticationMethod?: string;
    scopes?: string[];
    claims?: Record<string, any>;
    ipAddress?: string;
    userAgent?: string;
    geoLocation?: GeoLocation;
}
export interface ValidationRule {
    field: string;
    ruleType: ValidationRuleType;
    required?: boolean;
    minLength?: number;
    maxLength?: number;
    pattern?: string;
    minValue?: number;
    maxValue?: number;
    allowedValues?: string[];
    customValidator?: string;
}
export interface ValidationError {
    field: string;
    code: string;
    message: string;
    details?: Record<string, any>;
}
export interface ValidationWarning {
    field: string;
    code: string;
    message: string;
    details?: Record<string, any>;
}
export interface ValidationResult {
    valid: boolean;
    errors?: ValidationError[];
    warnings?: ValidationWarning[];
    sanitizedData?: Record<string, any>;
}
export interface SDKExecuteWorkflowRequest {
    inputData?: Record<string, any>;
    callbackUrl?: string;
    timeoutMs?: number;
    waitForCompletion?: boolean;
    metadata?: Record<string, any>;
    httpRequest?: HTTPRequestContext;
    session?: SessionContext;
    security?: SecurityContext;
    validationRules?: ValidationRule[];
    enableValidation?: boolean;
    enableSanitization?: boolean;
    includeFullContext?: boolean;
}
export interface SDKExecuteWorkflowResponse {
    runId: RunId;
    status: AsyncExecutionStatus;
    message: string;
    statusUrl: string;
    resultUrl: string;
    pollAfterMs?: number;
    estimatedCompletionMs?: number;
    expiresAt: Timestamp;
    validationResult?: ValidationResult;
    contextIncluded: boolean;
    sdkVersion: string;
    requestId: string;
}
export interface BatchExecutionItem {
    workflowId: WorkflowId;
    request?: SDKExecuteWorkflowRequest;
    priority?: number;
}
export interface BatchExecuteWorkflowsRequest {
    executions: BatchExecutionItem[];
    parallel?: boolean;
    maxConcurrent?: number;
    stopOnFirstFailure?: boolean;
}
export interface BatchExecutionResult {
    workflowId: WorkflowId;
    success: boolean;
    runId?: RunId;
    error?: ValidationError;
    executionTimeMs?: number;
}
export interface BatchError {
    workflowId: WorkflowId;
    error: ValidationError;
    timestamp: Timestamp;
}
export interface BatchExecuteWorkflowsResponse {
    batchId: string;
    total: number;
    successful: number;
    failed: number;
    pending: number;
    executions?: BatchExecutionResult[];
    errors?: BatchError[];
}
export interface WebhookConfiguration {
    webhookId: string;
    url: string;
    events: WebhookEvent[];
    secret?: string;
    enabled?: boolean;
    retryCount?: number;
    timeoutMs?: number;
    headers?: HTTPHeaders;
}
export interface ListWebhooksResponse {
    webhooks: WebhookConfiguration[];
    count: number;
}
export interface SDKConfiguration {
    sdkVersion: string;
    workflowApiEndpoint: string;
    timeoutMs?: number;
    maxRetries?: number;
    retryDelayMs?: number;
    authType?: AuthType;
    authToken?: string;
    enableValidation?: boolean;
    enableSanitization?: boolean;
    strictValidation?: boolean;
    enableSessionExtraction?: boolean;
    enableSecurityContext?: boolean;
    includeFullHttpContext?: boolean;
    logLevel?: LogLevel;
    enableRequestLogging?: boolean;
    enableMetrics?: boolean;
    defaultValidationRules?: ValidationRule[];
    customValidators?: string[];
    asyncExecution?: boolean;
    defaultPriority?: number;
    pollIntervalMs?: number;
    estimatedCompletionMs?: number;
    executionExpiryDurationMs?: number;
    enableCircuitBreaker?: boolean;
    circuitBreakerThreshold?: number;
    circuitBreakerTimeoutMs?: number;
    maxRedirects?: number;
}
export interface SDKError extends Error {
    code: string;
    message: string;
    details?: Record<string, any>;
    field?: string;
    originalError?: Error;
}
export interface RequestValidationError extends SDKError {
    validationResult: ValidationResult;
}
export interface BatchExecutionError extends SDKError {
    batchId: string;
    failedExecutions?: BatchExecutionResult[];
    errors?: BatchError[];
}
export interface WebhookError extends SDKError {
    webhookId?: string;
}
export interface SDKConfigurationError extends SDKError {
    field?: string;
}
export declare function createHTTPRequestContext(method: string, path: string): HTTPRequestContext;
export declare function createSessionContext(userId?: string, sessionId?: string): SessionContext;
export declare function createSecurityContext(authenticated: boolean): SecurityContext;
export declare function createSDKExecuteWorkflowRequest(inputData?: Record<string, any>): SDKExecuteWorkflowRequest;
//# sourceMappingURL=models.d.ts.map