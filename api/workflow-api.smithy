$version: "2"

namespace unified.workflow.api

use alloy#simpleRestJson
use alloy#uuidFormat
use smithy.framework#ValidationException

/// Workflow API Service
@simpleRestJson
service WorkflowApi {
    version: "2024-01-13",
    resources: [Workflow, Execution],
    operations: [
        ListWorkflows,
        CreateWorkflow,
        GetWorkflow,
        DeleteWorkflow,
        ExecuteWorkflow,
        AsyncExecuteWorkflow,  // NEW: Async execution
        ListExecutions,
        GetExecutionStatus,
        GetExecutionResult,    // NEW: Get async result
        CancelExecution,
        PauseExecution,
        ResumeExecution,
        RetryExecution,
        GetExecutionData,
        GetExecutionMetrics,
        GetExecutionDetails,   // NEW: Enhanced execution details
        GetStepExecution,      // NEW: Step execution details
        GetChildStepExecution  // NEW: Child step execution details
    ]
}

/// Workflow resource
resource Workflow {
    identifiers: { workflowId: WorkflowId },
    read: GetWorkflow,
    delete: DeleteWorkflow,
    operations: [ExecuteWorkflow, AsyncExecuteWorkflow],  // Added AsyncExecuteWorkflow
    collectionOperations: [ListWorkflows, CreateWorkflow]
}

/// Execution resource
resource Execution {
    identifiers: { runId: RunId },
    read: GetExecutionStatus,
    operations: [
        CancelExecution,
        PauseExecution,
        ResumeExecution,
        RetryExecution,
        GetExecutionData,
        GetExecutionMetrics,
        GetExecutionResult     // NEW: Added result retrieval
    ],
    collectionOperations: [ListExecutions]
}

// ============ COMMON TYPES ============

/// Workflow ID
@uuidFormat
string WorkflowId

/// Execution Run ID
@uuidFormat
string RunId

/// Workflow status
enum WorkflowStatus {
    PENDING = "pending"
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"
    PAUSED = "paused"
    QUEUED = "queued"      // NEW: For async execution
    PROCESSING = "processing" // NEW: For async execution
}

/// Async execution status (for polling responses)
enum AsyncExecutionStatus {
    QUEUED = "queued"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"
    TIMEOUT = "timeout"
}

/// Step information
structure StepInfo {
    @required
    name: String,

    @required
    childStepCount: Integer,

    @required
    isParallel: Boolean
}

/// Workflow information
structure WorkflowInfo {
    @required
    workflowId: WorkflowId,

    @required
    name: String,

    @required
    description: String,

    @required
    stepCount: Integer,

    steps: StepInfoList
}

list WorkflowInfoList {
    member: WorkflowInfo
}

list StepInfoList {
    member: StepInfo
}

/// Execution information
structure ExecutionInfo {
    @required
    runId: RunId,

    @required
    workflowId: WorkflowId,

    @required
    status: WorkflowStatus,

    @required
    currentStepIndex: Integer,

    @required
    currentChildStepIndex: Integer,

    startTime: Timestamp,
    endTime: Timestamp,
    errorMessage: String,
    lastAttemptedStep: String,

    @required
    isTerminal: Boolean,

    @required
    isRunning: Boolean,

    @required
    isPending: Boolean,

    @required
    createdAt: Timestamp,

    @required
    updatedAt: Timestamp
}

list ExecutionInfoList {
    member: ExecutionInfo
}

/// Execution status
structure ExecutionStatus {
    @required
    runId: RunId,

    @required
    workflowId: WorkflowId,

    @required
    status: WorkflowStatus,

    currentStep: String,
    
    @required
    currentStepIndex: Integer,

    @required
    currentChildStepIndex: Integer,

    @required
    progress: Float,

    startTime: Timestamp,
    endTime: Timestamp,
    errorMessage: String,
    lastAttemptedStep: String,

    @required
    isTerminal: Boolean,

    metadata: Document
}

/// Child step status (no RUNNING since primitives execute synchronously)
enum ChildStepStatus {
    PENDING = "pending"
    COMPLETED = "completed"
    FAILED = "failed"
    SKIPPED = "skipped"
}

/// Child step execution information
structure ChildStepExecution {
    @required
    stepIndex: Integer,
    
    @required
    childStepIndex: Integer,
    
    @required
    name: String,
    
    @required
    status: ChildStepStatus,
    
    @required
    primitiveName: String,
    
    startTime: Timestamp,
    endTime: Timestamp,
    durationMillis: Long,
    errorMessage: String,
    result: Document,
    parameters: Document
}

list ChildStepExecutionList {
    member: ChildStepExecution
}

/// Step status (includes RUNNING for step progress tracking)
enum StepStatus {
    PENDING = "pending"
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"
}

/// Step execution information
structure StepExecution {
    @required
    stepIndex: Integer,
    
    @required
    name: String,
    
    @required
    status: StepStatus,
    
    @required
    isParallel: Boolean,
    
    @required
    childStepCount: Integer,
    
    @required
    completedChildSteps: Integer,
    
    @required
    failedChildSteps: Integer,
    
    startTime: Timestamp,
    endTime: Timestamp,
    durationMillis: Long,
    childSteps: ChildStepExecutionList,
    errorMessage: String
}

list StepExecutionList {
    member: StepExecution
}

/// Enhanced execution status with child-step details
structure EnhancedExecutionStatus {
    @required
    runId: RunId,
    
    @required
    workflowId: WorkflowId,
    
    @required
    status: WorkflowStatus,
    
    @required
    currentStepIndex: Integer,
    
    @required
    currentChildStepIndex: Integer,
    
    @required
    progress: Float,
    
    @required
    totalSteps: Integer,
    
    @required
    completedSteps: Integer,
    
    @required
    totalChildSteps: Integer,
    
    @required
    completedChildSteps: Integer,
    
    startTime: Timestamp,
    endTime: Timestamp,
    currentStep: StepExecution,
    completedStepsList: StepExecutionList,
    errorMessage: String,
    lastAttemptedStep: String,
    
    @required
    isTerminal: Boolean,
    
    metadata: Document
}

/// Execution result (for async execution)
structure ExecutionResult {
    @required
    runId: RunId,

    @required
    workflowId: WorkflowId,

    @required
    status: AsyncExecutionStatus,

    @required
    result: Document,

    @required
    completedAt: Timestamp,

    executionTimeMillis: Long,
    stepCount: Integer,
    errorDetails: String
}

/// Async execution response (202 Accepted)
structure AsyncExecuteWorkflowResponse {
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
    expiresAt: Timestamp
}

/// Async execution poll response
structure PollExecutionResultResponse {
    @required
    runId: RunId,

    @required
    status: AsyncExecutionStatus,

    result: ExecutionResult,
    pollAfterMs: Integer,
    estimatedCompletionMs: Integer,
    progress: Float
}

/// Execution metrics
structure ExecutionMetrics {
    @required
    runId: RunId,

    @required
    workflowId: WorkflowId,

    workflowMetrics: Document,
    stepMetrics: Document,
    childStepMetrics: Document,

    @required
    totalSteps: Integer,

    @required
    completedSteps: Integer,

    @required
    failedSteps: Integer,

    @required
    totalChildSteps: Integer,

    @required
    completedChildSteps: Integer,

    @required
    failedChildSteps: Integer,

    @required
    totalDurationMillis: Long,

    @required
    averageStepDuration: Float,

    @required
    successRate: Float
}

/// Execution filters
structure ExecutionFilters {
    workflowId: WorkflowId,
    status: WorkflowStatus,
    limit: Integer = 50,
    offset: Integer = 0
}

/// Create workflow request
structure CreateWorkflowRequest {
    @required
    name: String,

    description: String
}

/// Create workflow response
structure CreateWorkflowResponse {
    @required
    workflowId: WorkflowId,

    @required
    name: String,

    @required
    description: String,

    @required
    message: String
}

/// Execute workflow response (synchronous)
structure ExecuteWorkflowResponse {
    @required
    runId: RunId,

    @required
    message: String,

    @required
    statusUrl: String
}

/// Async execute workflow request
structure AsyncExecuteWorkflowRequest {
    @httpPayload
    inputData: Document,
    
    callbackUrl: String,
    timeoutMs: Integer = 30000,
    waitForCompletion: Boolean = false,
    metadata: Document
}

/// Execution data response
structure ExecutionDataResponse {
    @required
    runId: RunId,

    @required
    data: Document
}

/// List workflows response
structure ListWorkflowsResponse {
    @required
    workflows: WorkflowInfoList,

    @required
    count: Integer
}

/// List executions response
structure ListExecutionsResponse {
    @required
    executions: ExecutionInfoList,

    @required
    count: Integer,

    @required
    limit: Integer,

    @required
    offset: Integer
}

/// Success response
structure SuccessResponse {
    @required
    message: String
}

// ============ OPERATIONS ============

/// List all workflows
@http(method: "GET", uri: "/api/v1/workflows")
@readonly
operation ListWorkflows {
    input := {}
    output: ListWorkflowsResponse
    errors: [ValidationException]
}

/// Get workflow details
@http(method: "GET", uri: "/api/v1/workflows/{workflowId}")
@readonly
operation GetWorkflow {
    input: GetWorkflowInput
    output: WorkflowInfo
    errors: [ValidationException, WorkflowNotFoundError]
}

structure GetWorkflowInput {
    @httpLabel
    @required
    workflowId: WorkflowId
}

/// Create a new workflow
@http(method: "POST", uri: "/api/v1/workflows")
operation CreateWorkflow {
    input: CreateWorkflowRequest
    output: CreateWorkflowResponse
    errors: [ValidationException]
}

/// Delete a workflow
@http(method: "DELETE", uri: "/api/v1/workflows/{workflowId}")
operation DeleteWorkflow {
    input: DeleteWorkflowInput
    output: SuccessResponse
    errors: [ValidationException, WorkflowNotFoundError]
}

structure DeleteWorkflowInput {
    @httpLabel
    @required
    workflowId: WorkflowId
}

/// Execute a workflow (synchronous - existing)
@http(method: "POST", uri: "/api/v1/workflows/{workflowId}/execute")
operation ExecuteWorkflow {
    input: ExecuteWorkflowInput
    output: ExecuteWorkflowResponse
    errors: [ValidationException, WorkflowNotFoundError]
}

structure ExecuteWorkflowInput {
    @httpLabel
    @required
    workflowId: WorkflowId
}

/// Execute a workflow asynchronously (NEW)
@http(method: "POST", uri: "/api/v1/workflows/{workflowId}/async-execute", code: 202)
operation AsyncExecuteWorkflow {
    input: AsyncExecuteWorkflowInput
    output: AsyncExecuteWorkflowResponse
    errors: [ValidationException, WorkflowNotFoundError, RequestTimeoutError]
}

structure AsyncExecuteWorkflowInput {
    @httpLabel
    @required
    workflowId: WorkflowId,
    
    @httpPayload
    request: AsyncExecuteWorkflowRequest
}

/// Get execution result (for async execution) - NEW
@http(method: "GET", uri: "/api/v1/executions/{runId}/result")
@readonly
operation GetExecutionResult {
    input: GetExecutionResultInput
    output: PollExecutionResultResponse
    errors: [ValidationException, ExecutionNotFoundError, ResultNotReadyError]
}

structure GetExecutionResultInput {
    @httpLabel
    @required
    runId: RunId,
    
    @httpQuery("wait_ms")
    waitMs: Integer = 0,
    
    @httpQuery("long_poll")
    longPoll: Boolean = false
}

/// List executions
@http(method: "GET", uri: "/api/v1/executions")
@readonly
operation ListExecutions {
    input: ListExecutionsInput
    output: ListExecutionsResponse
    errors: [ValidationException]
}

structure ListExecutionsInput {
    @httpQuery("workflow_id")
    workflowId: WorkflowId,

    @httpQuery("status")
    status: WorkflowStatus,

    @httpQuery("limit")
    limit: Integer,

    @httpQuery("offset")
    offset: Integer
}

/// Get execution status
@http(method: "GET", uri: "/api/v1/executions/{runId}")
@readonly
operation GetExecutionStatus {
    input: GetExecutionStatusInput
    output: ExecutionStatus
    errors: [ValidationException, ExecutionNotFoundError]
}

structure GetExecutionStatusInput {
    @httpLabel
    @required
    runId: RunId
}

/// Cancel execution
@http(method: "POST", uri: "/api/v1/executions/{runId}/cancel")
operation CancelExecution {
    input: CancelExecutionInput
    output: SuccessResponse
    errors: [ValidationException, ExecutionNotFoundError]
}

structure CancelExecutionInput {
    @httpLabel
    @required
    runId: RunId
}

/// Pause execution
@http(method: "POST", uri: "/api/v1/executions/{runId}/pause")
operation PauseExecution {
    input: PauseExecutionInput
    output: SuccessResponse
    errors: [ValidationException, ExecutionNotFoundError]
}

structure PauseExecutionInput {
    @httpLabel
    @required
    runId: RunId
}

/// Resume execution
@http(method: "POST", uri: "/api/v1/executions/{runId}/resume")
operation ResumeExecution {
    input: ResumeExecutionInput
    output: SuccessResponse
    errors: [ValidationException, ExecutionNotFoundError]
}

structure ResumeExecutionInput {
    @httpLabel
    @required
    runId: RunId
}

/// Retry execution
@http(method: "POST", uri: "/api/v1/executions/{runId}/retry")
operation RetryExecution {
    input: RetryExecutionInput
    output: SuccessResponse
    errors: [ValidationException, ExecutionNotFoundError]
}

structure RetryExecutionInput {
    @httpLabel
    @required
    runId: RunId
}

/// Get execution data
@http(method: "GET", uri: "/api/v1/executions/{runId}/data")
@readonly
operation GetExecutionData {
    input: GetExecutionDataInput
    output: ExecutionDataResponse
    errors: [ValidationException, ExecutionNotFoundError]
}

structure GetExecutionDataInput {
    @httpLabel
    @required
    runId: RunId
}

/// Get execution metrics
@http(method: "GET", uri: "/api/v1/executions/{runId}/metrics")
@readonly
operation GetExecutionMetrics {
    input: GetExecutionMetricsInput
    output: ExecutionMetrics
    errors: [ValidationException, ExecutionNotFoundError]
}

structure GetExecutionMetricsInput {
    @httpLabel
    @required
    runId: RunId
}

/// Get detailed execution status with child-step information
@http(method: "GET", uri: "/api/v1/executions/{runId}/details")
@readonly
operation GetExecutionDetails {
    input: GetExecutionDetailsInput
    output: EnhancedExecutionStatus
    errors: [ValidationException, ExecutionNotFoundError]
}

structure GetExecutionDetailsInput {
    @httpLabel
    @required
    runId: RunId,
    
    @httpQuery("include_child_steps")
    includeChildSteps: Boolean = true
}

/// Get step execution details
@http(method: "GET", uri: "/api/v1/executions/{runId}/steps/{stepIndex}")
@readonly
operation GetStepExecution {
    input: GetStepExecutionInput
    output: StepExecution
    errors: [ValidationException, ExecutionNotFoundError, StepNotFoundError]
}

structure GetStepExecutionInput {
    @httpLabel
    @required
    runId: RunId,
    
    @httpLabel
    @required
    stepIndex: Integer
}

/// Get child step execution details
@http(method: "GET", uri: "/api/v1/executions/{runId}/steps/{stepIndex}/child-steps/{childStepIndex}")
@readonly
operation GetChildStepExecution {
    input: GetChildStepExecutionInput
    output: ChildStepExecution
    errors: [ValidationException, ExecutionNotFoundError, ChildStepNotFoundError]
}

structure GetChildStepExecutionInput {
    @httpLabel
    @required
    runId: RunId,
    
    @httpLabel
    @required
    stepIndex: Integer,
    
    @httpLabel
    @required
    childStepIndex: Integer
}

// ============ ERRORS ============

/// Workflow not found error
@httpError(404)
@error("client")
structure WorkflowNotFoundError {
    @required
    error: String,

    @required
    details: String
}

/// Execution not found error
@httpError(404)
@error("client")
structure ExecutionNotFoundError {
    @required
    error: String,

    @required
    details: String
}

/// Step not found error
@httpError(404)
@error("client")
structure StepNotFoundError {
    @required
    error: String,
    
    @required
    details: String,
    
    @required
    stepIndex: Integer
}

/// Child step not found error
@httpError(404)
@error("client")
structure ChildStepNotFoundError {
    @required
    error: String,
    
    @required
    details: String,
    
    @required
    stepIndex: Integer,
    
    @required
    childStepIndex: Integer
}

/// Result not ready error (for async polling)
@httpError(202)
@error("client")
structure ResultNotReadyError {
    @required
    runId: RunId,

    @required
    status: AsyncExecutionStatus,

    @required
    message: String,

    pollAfterMs: Integer,
    estimatedCompletionMs: Integer
}

/// Request timeout error
@httpError(408)
@error("client")
structure RequestTimeoutError {
    @required
    error: String,

    @required
    details: String
}

/// Internal server error
@httpError(500)
@error("server")
structure InternalServerError {
    @required
    error: String,

    @required
    details: String
}
