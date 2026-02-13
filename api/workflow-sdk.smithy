// api/workflow-sdk.smithy
$version: "2"

namespace unified.workflow.sdk

use unified.workflow.api#WorkflowApi
use unified.workflow.api#WorkflowId
use unified.workflow.api#RunId
use unified.workflow.api#WorkflowStatus
use unified.workflow.api#AsyncExecutionStatus
use unified.workflow.api#Document
use unified.workflow.api#Timestamp
use alloy#simpleRestJson
use smithy.framework#ValidationException

/// Workflow SDK Service
/// Provides enhanced operations for client applications to execute workflows
/// with HTTP request context, validation, and session information.
@simpleRestJson
service WorkflowSDK {
    version: "2024-01-13",
    resources: [SDKWorkflow, SDKExecution],
    operations: [
        // Core SDK operations
        ExecuteWorkflowFromHTTP,
        ExecuteWorkflowWithContext,
        ValidateAndExecuteWorkflow,
        
        // Batch operations
        BatchExecuteWorkflows,
        
        // Webhook operations  
        RegisterWebhook,
        UnregisterWebhook,
        ListWebhooks,
        
        // SDK management
        GetSDKConfiguration,
        UpdateSDKConfiguration,
        
        // Import existing WorkflowApi operations
        ListWorkflows,
        CreateWorkflow,
        GetWorkflow,
        DeleteWorkflow,
        ExecuteWorkflow,
        AsyncExecuteWorkflow,
        ListExecutions,
        GetExecutionStatus,
        GetExecutionResult,
        CancelExecution,
        PauseExecution,
        ResumeExecution,
        RetryExecution,
        GetExecutionData,
        GetExecutionMetrics,
        GetExecutionDetails,
        GetStepExecution,
        GetChildStepExecution
    ]
}

/// SDK Workflow resource (extends Workflow)
resource SDKWorkflow {
    identifiers: { workflowId: WorkflowId },
    read: GetWorkflow,
    delete: DeleteWorkflow,
    operations: [
        ExecuteWorkflow,
        AsyncExecuteWorkflow,
        ExecuteWorkflowFromHTTP,      // SDK-specific
        ExecuteWorkflowWithContext,   // SDK-specific
        ValidateAndExecuteWorkflow    // SDK-specific
    ],
    collectionOperations: [ListWorkflows, CreateWorkflow]
}

/// SDK Execution resource (extends Execution)
resource SDKExecution {
    identifiers: { runId: RunId },
    read: GetExecutionStatus,
    operations: [
        CancelExecution,
        PauseExecution,
        ResumeExecution,
        RetryExecution,
        GetExecutionData,
        GetExecutionMetrics,
        GetExecutionResult
    ],
    collectionOperations: [ListExecutions]
}

// ============ SDK-SPECIFIC TYPES ============

/// HTTP Request Context
/// Contains complete HTTP request information for workflow execution
structure HTTPRequestContext {
    @required
    method: String,  // GET, POST, PUT, DELETE, etc.
    
    @required
    path: String,
    
    @required
    headers: HTTPHeaders,
    
    queryParams: QueryParams,
    pathParams: PathParams,
    body: Document,
    remoteAddr: String,
    userAgent: String,
    
    @required
    timestamp: Timestamp
}

map HTTPHeaders {
    key: String,
    value: StringList
}

list StringList {
    member: String
}

map QueryParams {
    key: String,
    value: StringList
}

map PathParams {
    key: String,
    value: String
}

/// Session Context
/// Contains user session information for security and authorization
structure SessionContext {
    userId: String,
    sessionId: String,
    roles: StringList,
    permissions: StringList,
    authMethod: String,  // "jwt", "api_key", "oauth", "basic"
    expiresAt: Timestamp,
    attributes: Document
}

/// Security Context
/// Contains security-related information for the request
structure SecurityContext {
    @required
    authenticated: Boolean,
    
    authenticationMethod: String,
    scopes: StringList,
    claims: Document,
    ipAddress: String,
    userAgent: String,
    geoLocation: GeoLocation
}

structure GeoLocation {
    country: String,
    region: String,
    city: String,
    latitude: Float,
    longitude: Float
}

/// Validation Rules
/// Defines validation rules for request data
structure ValidationRule {
    @required
    field: String,
    
    @required
    ruleType: ValidationRuleType,
    
    required: Boolean,
    minLength: Integer,
    maxLength: Integer,
    pattern: String,
    minValue: Double,
    maxValue: Double,
    allowedValues: StringList,
    customValidator: String
}

enum ValidationRuleType {
    REQUIRED = "required"
    STRING = "string"
    NUMBER = "number"
    BOOLEAN = "boolean"
    ARRAY = "array"
    OBJECT = "object"
    EMAIL = "email"
    URL = "url"
    UUID = "uuid"
    CUSTOM = "custom"
}

list ValidationRuleList {
    member: ValidationRule
}

/// Validation Result
/// Result of request validation
structure ValidationResult {
    @required
    valid: Boolean,
    
    errors: ValidationErrorList,
    warnings: ValidationWarningList,
    sanitizedData: Document
}

list ValidationErrorList {
    member: ValidationError
}

list ValidationWarningList {
    member: ValidationWarning
}

structure ValidationError {
    @required
    field: String,
    
    @required
    code: String,
    
    @required
    message: String,
    
    details: Document
}

structure ValidationWarning {
    @required
    field: String,
    
    @required
    code: String,
    
    @required
    message: String,
    
    details: Document
}

/// SDK Execution Request
/// Enhanced execution request with full HTTP and session context
structure SDKExecuteWorkflowRequest {
    /// Workflow execution parameters (from base)
    @httpPayload
    inputData: Document,
    
    callbackUrl: String,
    timeoutMs: Integer = 30000,
    waitForCompletion: Boolean = false,
    metadata: Document,
    
    /// SDK-specific extensions
    httpRequest: HTTPRequestContext,
    session: SessionContext,
    security: SecurityContext,
    validationRules: ValidationRuleList,
    enableValidation: Boolean = true,
    enableSanitization: Boolean = true,
    includeFullContext: Boolean = true
}

/// SDK Execution Response
/// Enhanced execution response with validation results
structure SDKExecuteWorkflowResponse {
    /// Base response fields
    @required
    runId: RunId,

    @required
    status: AsyncExecutionStatus,

    @required
    message: String,

    @required
    statusUrl: String,

    @required
    resultUrl: String,

    pollAfterMs: Integer = 1000,
    estimatedCompletionMs: Integer = 5000,
    expiresAt: Timestamp,
    
    /// SDK-specific extensions
    validationResult: ValidationResult,
    contextIncluded: Boolean,
    sdkVersion: String,
    requestId: String
}

/// Batch Execution Request
/// Request to execute multiple workflows in batch
structure BatchExecuteWorkflowsRequest {
    @required
    executions: BatchExecutionList,
    
    @required
    parallel: Boolean = false,
    
    maxConcurrent: Integer = 10,
    stopOnFirstFailure: Boolean = false
}

list BatchExecutionList {
    member: BatchExecutionItem
}

structure BatchExecutionItem {
    @required
    workflowId: WorkflowId,
    
    request: SDKExecuteWorkflowRequest,
    priority: Integer = 5
}

/// Batch Execution Response
/// Response for batch execution
structure BatchExecuteWorkflowsResponse {
    @required
    batchId: String,
    
    @required
    total: Integer,
    
    @required
    successful: Integer,
    
    @required
    failed: Integer,
    
    @required
    pending: Integer,
    
    executions: BatchExecutionResultList,
    errors: BatchErrorList
}

list BatchExecutionResultList {
    member: BatchExecutionResult
}

structure BatchExecutionResult {
    @required
    workflowId: WorkflowId,
    
    @required
    success: Boolean,
    
    runId: RunId,
    error: ValidationError,
    executionTimeMs: Long
}

list BatchErrorList {
    member: BatchError
}

structure BatchError {
    @required
    workflowId: WorkflowId,
    
    @required
    error: ValidationError,
    
    @required
    timestamp: Timestamp
}

/// Webhook Configuration
/// Configuration for webhook notifications
structure WebhookConfiguration {
    @required
    webhookId: String,
    
    @required
    url: String,
    
    @required
    events: WebhookEventList,
    
    secret: String,
    enabled: Boolean = true,
    retryCount: Integer = 3,
    timeoutMs: Integer = 5000,
    headers: HTTPHeaders
}

list WebhookEventList {
    member: WebhookEvent
}

enum WebhookEvent {
    WORKFLOW_STARTED = "workflow_started"
    WORKFLOW_COMPLETED = "workflow_completed"
    WORKFLOW_FAILED = "workflow_failed"
    WORKFLOW_CANCELLED = "workflow_cancelled"
    STEP_COMPLETED = "step_completed"
    STEP_FAILED = "step_failed"
}

/// SDK Configuration
/// Configuration for the SDK client
structure SDKConfiguration {
    @required
    sdkVersion: String,
    
    @required
    workflowApiEndpoint: String,
    
    timeoutMs: Integer = 30000,
    maxRetries: Integer = 3,
    retryDelayMs: Integer = 1000,
    
    authType: AuthType = "BEARER_TOKEN",
    authToken: String,
    
    enableValidation: Boolean = true,
    enableSanitization: Boolean = true,
    strictValidation: Boolean = false,
    
    enableSessionExtraction: Boolean = true,
    enableSecurityContext: Boolean = true,
    includeFullHttpContext: Boolean = true,
    
    logLevel: LogLevel = "INFO",
    enableRequestLogging: Boolean = true,
    enableMetrics: Boolean = true,
    
    defaultValidationRules: ValidationRuleList,
    customValidators: StringList
}

enum AuthType {
    BEARER_TOKEN = "bearer_token"
    API_KEY = "api_key"
    BASIC_AUTH = "basic_auth"
    OAUTH2 = "oauth2"
    AWS_SIGV4 = "aws_sigv4"
}

enum LogLevel {
    DEBUG = "debug"
    INFO = "info"
    WARN = "warn"
    ERROR = "error"
}

// ============ SDK-SPECIFIC OPERATIONS ============

/// Execute workflow from HTTP request
/// Accepts HTTP request context and executes workflow with full context
@http(method: "POST", uri: "/sdk/v1/workflows/{workflowId}/execute/http", code: 202)
operation ExecuteWorkflowFromHTTP {
    input: ExecuteWorkflowFromHTTPInput
    output: SDKExecuteWorkflowResponse
    errors: [ValidationException, WorkflowNotFoundError, RequestValidationError]
}

structure ExecuteWorkflowFromHTTPInput {
    @httpLabel
    @required
    workflowId: WorkflowId,
    
    @httpPayload
    @required
    httpRequest: HTTPRequestContext,
    
    validationRules: ValidationRuleList,
    enableValidation: Boolean = true,
    enableSanitization: Boolean = true,
    includeSession: Boolean = true,
    includeSecurity: Boolean = true
}

/// Execute workflow with full context
/// Accepts complete SDK execution request with all context
@http(method: "POST", uri: "/sdk/v1/workflows/{workflowId}/execute/context", code: 202)
operation ExecuteWorkflowWithContext {
    input: ExecuteWorkflowWithContextInput
    output: SDKExecuteWorkflowResponse
    errors: [ValidationException, WorkflowNotFoundError, RequestValidationError]
}

structure ExecuteWorkflowWithContextInput {
    @httpLabel
    @required
    workflowId: WorkflowId,
    
    @httpPayload
    @required
    request: SDKExecuteWorkflowRequest
}

/// Validate and execute workflow
/// Validates input data before execution
@http(method: "POST", uri: "/sdk/v1/workflows/{workflowId}/execute/validate", code: 202)
operation ValidateAndExecuteWorkflow {
    input: ValidateAndExecuteWorkflowInput
    output: SDKExecuteWorkflowResponse
    errors: [ValidationException, WorkflowNotFoundError, RequestValidationError]
}

structure ValidateAndExecuteWorkflowInput {
    @httpLabel
    @required
    workflowId: WorkflowId,
    
    @httpPayload
    @required
    inputData: Document,
    
    @required
    validationRules: ValidationRuleList,
    
    httpRequest: HTTPRequestContext,
    session: SessionContext,
    security: SecurityContext,
    enableSanitization: Boolean = true
}

/// Batch execute multiple workflows
/// Execute multiple workflows in a single request
@http(method: "POST", uri: "/sdk/v1/batch/execute", code: 202)
operation BatchExecuteWorkflows {
    input: BatchExecuteWorkflowsRequest
    output: BatchExecuteWorkflowsResponse
    errors: [ValidationException, BatchExecutionError]
}

/// Register webhook for events
/// Register a webhook to receive workflow events
@http(method: "POST", uri: "/sdk/v1/webhooks")
operation RegisterWebhook {
    input: RegisterWebhookInput
    output: WebhookConfiguration
    errors: [ValidationException, WebhookError]
}

structure RegisterWebhookInput {
    @httpPayload
    @required
    configuration: WebhookConfiguration
}

/// Unregister webhook
/// Remove a registered webhook
@http(method: "DELETE", uri: "/sdk/v1/webhooks/{webhookId}")
operation UnregisterWebhook {
    input: UnregisterWebhookInput
    output: SuccessResponse
    errors: [ValidationException, WebhookNotFoundError]
}

structure UnregisterWebhookInput {
    @httpLabel
    @required
    webhookId: String
}

/// List registered webhooks
/// Get all registered webhooks
@http(method: "GET", uri: "/sdk/v1/webhooks")
@readonly
operation ListWebhooks {
    input := {}
    output: ListWebhooksResponse
    errors: [ValidationException]
}

structure ListWebhooksResponse {
    @required
    webhooks: WebhookList,
    
    @required
    count: Integer
}

list WebhookList {
    member: WebhookConfiguration
}

/// Get SDK configuration
/// Get current SDK configuration
@http(method: "GET", uri: "/sdk/v1/configuration")
@readonly
operation GetSDKConfiguration {
    input := {}
    output: SDKConfiguration
    errors: [ValidationException]
}

/// Update SDK configuration
/// Update SDK configuration
@http(method: "PUT", uri: "/sdk/v1/configuration")
operation UpdateSDKConfiguration {
    input: UpdateSDKConfigurationInput
    output: SDKConfiguration
    errors: [ValidationException]
}

structure UpdateSDKConfigurationInput {
    @httpPayload
    @required
    configuration: SDKConfiguration
}

// ============ SDK-SPECIFIC ERRORS ============

/// Request validation error
@httpError(400)
@error("client")
structure RequestValidationError {
    @required
    error: String,
    
    @required
    details: String,
    
    @required
    validationResult: ValidationResult
}

/// Batch execution error
@httpError(400)
@error("client")
structure BatchExecutionError {
    @required
    error: String,
    
    @required
    details: String,
    
    @required
    batchId: String,
    
    failedExecutions: BatchExecutionResultList,
    errors: BatchErrorList
}

/// Webhook error
@httpError(400)
@error("client")
structure WebhookError {
    @required
    error: String,
    
    @required
    details: String,
    
    webhookId: String
}

/// Webhook not found error
@httpError(404)
@error("client")
structure WebhookNotFoundError {
    @required
    error: String,
    
    @required
    details: String,
    
    @required
    webhookId: String
}

/// SDK configuration error
@httpError(400)
@error("client")
structure SDKConfigurationError {
    @required
    error: String,
    
    @required
    details: String,
    
    field: String
}

// ============ RE-EXPORT EXISTING TYPES ============

/// Re-exported from unified.workflow.api
@documentation("Re-exported from unified.workflow.api")
string WorkflowId

/// Re-exported from unified.workflow.api
@documentation("Re-exported from unified.workflow.api")
string RunId

/// Re-exported from unified.workflow.api
@documentation("Re-exported from unified.workflow.api")
enum WorkflowStatus

/// Re-exported from unified.workflow.api
@documentation("Re-exported from unified.workflow.api")
enum AsyncExecutionStatus

/// Re-exported from unified.workflow.api
@documentation("Re-exported from unified.workflow.api")
structure SuccessResponse {
    @required
    message: String
}