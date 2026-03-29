/**
 * Main client for Unified Workflow SDK
 */
import { WorkflowId, RunId, SDKExecuteWorkflowRequest, SDKExecuteWorkflowResponse, BatchExecuteWorkflowsRequest, BatchExecuteWorkflowsResponse, ValidationRule, WebhookConfiguration, ListWebhooksResponse, SDKConfiguration } from './models';
export interface WorkflowSDKClient {
    executeFromHTTPRequest(workflowId: WorkflowId, request: Request): Promise<SDKExecuteWorkflowResponse>;
    executeWorkflow(workflowId: WorkflowId, data: Record<string, any>): Promise<SDKExecuteWorkflowResponse>;
    executeWorkflowWithContext(workflowId: WorkflowId, request: SDKExecuteWorkflowRequest): Promise<SDKExecuteWorkflowResponse>;
    validateAndExecuteWorkflow(workflowId: WorkflowId, data: Record<string, any>, rules: ValidationRule[]): Promise<SDKExecuteWorkflowResponse>;
    batchExecuteWorkflows(request: BatchExecuteWorkflowsRequest): Promise<BatchExecuteWorkflowsResponse>;
    registerWebhook(configuration: WebhookConfiguration): Promise<WebhookConfiguration>;
    unregisterWebhook(webhookId: string): Promise<void>;
    listWebhooks(): Promise<ListWebhooksResponse>;
    getSDKConfiguration(): Promise<SDKConfiguration>;
    updateSDKConfiguration(configuration: SDKConfiguration): Promise<SDKConfiguration>;
    listWorkflows(): Promise<any>;
    createWorkflow(workflow: any): Promise<any>;
    getWorkflow(workflowId: WorkflowId): Promise<any>;
    deleteWorkflow(workflowId: WorkflowId): Promise<void>;
    executeWorkflowAsync(workflowId: WorkflowId, request: any): Promise<any>;
    listExecutions(): Promise<any>;
    getExecutionStatus(runId: RunId): Promise<any>;
    getExecutionResult(runId: RunId): Promise<any>;
    cancelExecution(runId: RunId): Promise<void>;
    pauseExecution(runId: RunId): Promise<void>;
    resumeExecution(runId: RunId): Promise<void>;
    retryExecution(runId: RunId): Promise<any>;
    getExecutionData(runId: RunId): Promise<any>;
    getExecutionMetrics(runId: RunId): Promise<any>;
    getExecutionDetails(runId: RunId): Promise<any>;
    getStepExecution(runId: RunId, stepId: string): Promise<any>;
    getChildStepExecution(runId: RunId, childStepId: string): Promise<any>;
    ping(): Promise<void>;
    close(): Promise<void>;
}
export declare class UnifiedWorkflowSDK implements WorkflowSDKClient {
    private config;
    private baseUrl;
    private headers;
    constructor(config: Partial<SDKConfiguration>);
    executeFromHTTPRequest(workflowId: WorkflowId, request: Request): Promise<SDKExecuteWorkflowResponse>;
    executeWorkflow(workflowId: WorkflowId, data: Record<string, any>): Promise<SDKExecuteWorkflowResponse>;
    executeWorkflowWithContext(workflowId: WorkflowId, request: SDKExecuteWorkflowRequest): Promise<SDKExecuteWorkflowResponse>;
    validateAndExecuteWorkflow(workflowId: WorkflowId, data: Record<string, any>, rules: ValidationRule[]): Promise<SDKExecuteWorkflowResponse>;
    batchExecuteWorkflows(request: BatchExecuteWorkflowsRequest): Promise<BatchExecuteWorkflowsResponse>;
    registerWebhook(configuration: WebhookConfiguration): Promise<WebhookConfiguration>;
    unregisterWebhook(webhookId: string): Promise<void>;
    listWebhooks(): Promise<ListWebhooksResponse>;
    getSDKConfiguration(): Promise<SDKConfiguration>;
    updateSDKConfiguration(configuration: SDKConfiguration): Promise<SDKConfiguration>;
    listWorkflows(): Promise<any>;
    createWorkflow(workflow: any): Promise<any>;
    getWorkflow(workflowId: WorkflowId): Promise<any>;
    deleteWorkflow(workflowId: WorkflowId): Promise<void>;
    executeWorkflowAsync(workflowId: WorkflowId, request: any): Promise<any>;
    listExecutions(): Promise<any>;
    getExecutionStatus(runId: RunId): Promise<any>;
    getExecutionResult(runId: RunId): Promise<any>;
    cancelExecution(runId: RunId): Promise<void>;
    pauseExecution(runId: RunId): Promise<void>;
    resumeExecution(runId: RunId): Promise<void>;
    retryExecution(runId: RunId): Promise<any>;
    getExecutionData(runId: RunId): Promise<any>;
    getExecutionMetrics(runId: RunId): Promise<any>;
    getExecutionDetails(runId: RunId): Promise<any>;
    getStepExecution(runId: RunId, stepId: string): Promise<any>;
    getChildStepExecution(runId: RunId, childStepId: string): Promise<any>;
    ping(): Promise<void>;
    close(): Promise<void>;
    private makeRequest;
    private handleHTTPError;
    private shouldNotRetry;
    private extractHTTPRequestContext;
}
//# sourceMappingURL=client.d.ts.map