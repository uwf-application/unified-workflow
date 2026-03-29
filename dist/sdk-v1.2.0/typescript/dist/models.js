/**
 * TypeScript models for Unified Workflow SDK
 * Based on Smithy API definition: api/workflow-sdk.smithy
 */
export var WorkflowStatus;
(function (WorkflowStatus) {
    WorkflowStatus["DRAFT"] = "draft";
    WorkflowStatus["ACTIVE"] = "active";
    WorkflowStatus["INACTIVE"] = "inactive";
    WorkflowStatus["ARCHIVED"] = "archived";
})(WorkflowStatus || (WorkflowStatus = {}));
export var AsyncExecutionStatus;
(function (AsyncExecutionStatus) {
    AsyncExecutionStatus["PENDING"] = "pending";
    AsyncExecutionStatus["QUEUED"] = "queued";
    AsyncExecutionStatus["RUNNING"] = "running";
    AsyncExecutionStatus["COMPLETED"] = "completed";
    AsyncExecutionStatus["FAILED"] = "failed";
    AsyncExecutionStatus["CANCELLED"] = "cancelled";
    AsyncExecutionStatus["PAUSED"] = "paused";
})(AsyncExecutionStatus || (AsyncExecutionStatus = {}));
export var AuthType;
(function (AuthType) {
    AuthType["BEARER_TOKEN"] = "bearer_token";
    AuthType["API_KEY"] = "api_key";
    AuthType["BASIC_AUTH"] = "basic_auth";
    AuthType["OAUTH2"] = "oauth2";
    AuthType["AWS_SIGV4"] = "aws_sigv4";
})(AuthType || (AuthType = {}));
export var LogLevel;
(function (LogLevel) {
    LogLevel["DEBUG"] = "debug";
    LogLevel["INFO"] = "info";
    LogLevel["WARN"] = "warn";
    LogLevel["ERROR"] = "error";
})(LogLevel || (LogLevel = {}));
export var ValidationRuleType;
(function (ValidationRuleType) {
    ValidationRuleType["REQUIRED"] = "required";
    ValidationRuleType["STRING"] = "string";
    ValidationRuleType["NUMBER"] = "number";
    ValidationRuleType["BOOLEAN"] = "boolean";
    ValidationRuleType["ARRAY"] = "array";
    ValidationRuleType["OBJECT"] = "object";
    ValidationRuleType["EMAIL"] = "email";
    ValidationRuleType["URL"] = "url";
    ValidationRuleType["UUID"] = "uuid";
    ValidationRuleType["CUSTOM"] = "custom";
})(ValidationRuleType || (ValidationRuleType = {}));
export var WebhookEvent;
(function (WebhookEvent) {
    WebhookEvent["WORKFLOW_STARTED"] = "workflow_started";
    WebhookEvent["WORKFLOW_COMPLETED"] = "workflow_completed";
    WebhookEvent["WORKFLOW_FAILED"] = "workflow_failed";
    WebhookEvent["WORKFLOW_CANCELLED"] = "workflow_cancelled";
    WebhookEvent["STEP_COMPLETED"] = "step_completed";
    WebhookEvent["STEP_FAILED"] = "step_failed";
})(WebhookEvent || (WebhookEvent = {}));
// ============ HELPER FUNCTIONS ============
export function createHTTPRequestContext(method, path) {
    return {
        method,
        path,
        headers: {},
        timestamp: new Date().toISOString()
    };
}
export function createSessionContext(userId, sessionId) {
    return {
        userId,
        sessionId,
        roles: [],
        permissions: [],
        attributes: {}
    };
}
export function createSecurityContext(authenticated) {
    return {
        authenticated,
        scopes: [],
        claims: {}
    };
}
export function createSDKExecuteWorkflowRequest(inputData) {
    return {
        inputData: inputData || {},
        metadata: {},
        validationRules: [],
        enableValidation: true,
        enableSanitization: true,
        includeFullContext: true
    };
}
//# sourceMappingURL=models.js.map