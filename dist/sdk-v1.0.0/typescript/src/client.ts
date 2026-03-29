/**
 * Main client for Unified Workflow SDK
 */

import {
  WorkflowId,
  RunId,
  SDKExecuteWorkflowRequest,
  SDKExecuteWorkflowResponse,
  BatchExecuteWorkflowsRequest,
  BatchExecuteWorkflowsResponse,
  ValidationRule,
  WebhookConfiguration,
  ListWebhooksResponse,
  SDKConfiguration
} from './models';

import { SDKConfig } from './config';
import {
  WorkflowSDKError,
  RequestValidationError,
  BatchExecutionError,
  WebhookError,
  WorkflowNotFoundError,
  ExecutionNotFoundError,
  wrapError,
  ErrorCodes
} from './errors';

export interface WorkflowSDKClient {
  // Core SDK operations
  executeFromHTTPRequest(workflowId: WorkflowId, request: Request): Promise<SDKExecuteWorkflowResponse>;
  executeWorkflow(workflowId: WorkflowId, data: Record<string, any>): Promise<SDKExecuteWorkflowResponse>;
  executeWorkflowWithContext(workflowId: WorkflowId, request: SDKExecuteWorkflowRequest): Promise<SDKExecuteWorkflowResponse>;
  validateAndExecuteWorkflow(workflowId: WorkflowId, data: Record<string, any>, rules: ValidationRule[]): Promise<SDKExecuteWorkflowResponse>;
  
  // Batch operations
  batchExecuteWorkflows(request: BatchExecuteWorkflowsRequest): Promise<BatchExecuteWorkflowsResponse>;
  
  // Webhook operations
  registerWebhook(configuration: WebhookConfiguration): Promise<WebhookConfiguration>;
  unregisterWebhook(webhookId: string): Promise<void>;
  listWebhooks(): Promise<ListWebhooksResponse>;
  
  // SDK management
  getSDKConfiguration(): Promise<SDKConfiguration>;
  updateSDKConfiguration(configuration: SDKConfiguration): Promise<SDKConfiguration>;
  
  // Imported WorkflowApi operations
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
  
  // Health check
  ping(): Promise<void>;
  
  // Close client
  close(): Promise<void>;
}

export class UnifiedWorkflowSDK implements WorkflowSDKClient {
  private config: SDKConfig;
  private baseUrl: string;
  private headers: Record<string, string>;

  constructor(config: Partial<SDKConfiguration>) {
    this.config = new SDKConfig(config);
    this.baseUrl = this.config.getWorkflowApiEndpoint();
    this.headers = this.config.createHeaders();
  }

  // ============ CORE SDK OPERATIONS ============

  async executeFromHTTPRequest(workflowId: WorkflowId, request: Request): Promise<SDKExecuteWorkflowResponse> {
    try {
      // Extract HTTP context from request
      const httpRequest = await this.extractHTTPRequestContext(request);
      
      // Create SDK execution request
      const sdkRequest: SDKExecuteWorkflowRequest = {
        httpRequest,
        enableValidation: this.config.getValidationConfig().enableValidation,
        enableSanitization: this.config.getValidationConfig().enableSanitization,
        includeFullContext: this.config.shouldIncludeFullHttpContext(),
        validationRules: this.config.getValidationConfig().defaultValidationRules
      };

      return await this.executeWorkflowWithContext(workflowId, sdkRequest);
    } catch (error) {
      throw wrapError(error);
    }
  }

  async executeWorkflow(workflowId: WorkflowId, data: Record<string, any>): Promise<SDKExecuteWorkflowResponse> {
    try {
      const sdkRequest: SDKExecuteWorkflowRequest = {
        inputData: data,
        enableValidation: this.config.getValidationConfig().enableValidation,
        enableSanitization: this.config.getValidationConfig().enableSanitization,
        includeFullContext: this.config.shouldIncludeFullHttpContext(),
        validationRules: this.config.getValidationConfig().defaultValidationRules
      };

      return await this.executeWorkflowWithContext(workflowId, sdkRequest);
    } catch (error) {
      throw wrapError(error);
    }
  }

  async executeWorkflowWithContext(
    workflowId: WorkflowId,
    request: SDKExecuteWorkflowRequest
  ): Promise<SDKExecuteWorkflowResponse> {
    try {
      const url = `${this.baseUrl}/sdk/v1/workflows/${workflowId}/execute/context`;
      
      const response = await this.makeRequest<SDKExecuteWorkflowResponse>(url, {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(request)
      });

      return response;
    } catch (error) {
      throw wrapError(error);
    }
  }

  async validateAndExecuteWorkflow(
    workflowId: WorkflowId,
    data: Record<string, any>,
    rules: ValidationRule[]
  ): Promise<SDKExecuteWorkflowResponse> {
    try {
      const url = `${this.baseUrl}/sdk/v1/workflows/${workflowId}/execute/validate`;
      
      const request = {
        inputData: data,
        validationRules: rules,
        enableSanitization: this.config.getValidationConfig().enableSanitization
      };

      const response = await this.makeRequest<SDKExecuteWorkflowResponse>(url, {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(request)
      });

      return response;
    } catch (error) {
      throw wrapError(error);
    }
  }

  // ============ BATCH OPERATIONS ============

  async batchExecuteWorkflows(request: BatchExecuteWorkflowsRequest): Promise<BatchExecuteWorkflowsResponse> {
    try {
      const url = `${this.baseUrl}/sdk/v1/batch/execute`;
      
      const response = await this.makeRequest<BatchExecuteWorkflowsResponse>(url, {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify(request)
      });

      return response;
    } catch (error) {
      throw wrapError(error);
    }
  }

  // ============ WEBHOOK OPERATIONS ============

  async registerWebhook(configuration: WebhookConfiguration): Promise<WebhookConfiguration> {
    try {
      const url = `${this.baseUrl}/sdk/v1/webhooks`;
      
      const response = await this.makeRequest<WebhookConfiguration>(url, {
        method: 'POST',
        headers: this.headers,
        body: JSON.stringify({ configuration })
      });

      return response;
    } catch (error) {
      throw wrapError(error);
    }
  }

  async unregisterWebhook(webhookId: string): Promise<void> {
    try {
      const url = `${this.baseUrl}/sdk/v1/webhooks/${webhookId}`;
      
      await this.makeRequest<void>(url, {
        method: 'DELETE',
        headers: this.headers
      });
    } catch (error) {
      throw wrapError(error);
    }
  }

  async listWebhooks(): Promise<ListWebhooksResponse> {
    try {
      const url = `${this.baseUrl}/sdk/v1/webhooks`;
      
      const response = await this.makeRequest<ListWebhooksResponse>(url, {
        method: 'GET',
        headers: this.headers
      });

      return response;
    } catch (error) {
      throw wrapError(error);
    }
  }

  // ============ SDK MANAGEMENT ============

  async getSDKConfiguration(): Promise<SDKConfiguration> {
    try {
      const url = `${this.baseUrl}/sdk/v1/configuration`;
      
      const response = await this.makeRequest<SDKConfiguration>(url, {
        method: 'GET',
        headers: this.headers
      });

      return response;
    } catch (error) {
      throw wrapError(error);
    }
  }

  async updateSDKConfiguration(configuration: SDKConfiguration): Promise<SDKConfiguration> {
    try {
      const url = `${this.baseUrl}/sdk/v1/configuration`;
      
      const response = await this.makeRequest<SDKConfiguration>(url, {
        method: 'PUT',
        headers: this.headers,
        body: JSON.stringify({ configuration })
      });

      // Update local configuration
      this.config.updateConfig(configuration);
      this.headers = this.config.createHeaders();

      return response;
    } catch (error) {
      throw wrapError(error);
    }
  }

  // ============ IMPORTED WORKFLOW API OPERATIONS ============

  async listWorkflows(): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/workflows`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async createWorkflow(workflow: any): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/workflows`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(workflow)
    });
  }

  async getWorkflow(workflowId: WorkflowId): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/workflows/${workflowId}`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async deleteWorkflow(workflowId: WorkflowId): Promise<void> {
    await this.makeRequest(`${this.baseUrl}/api/v1/workflows/${workflowId}`, {
      method: 'DELETE',
      headers: this.headers
    });
  }

  async executeWorkflowAsync(workflowId: WorkflowId, request: any): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/workflows/${workflowId}/execute`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(request)
    });
  }

  async listExecutions(): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async getExecutionStatus(runId: RunId): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/status`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async getExecutionResult(runId: RunId): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/result`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async cancelExecution(runId: RunId): Promise<void> {
    await this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/cancel`, {
      method: 'POST',
      headers: this.headers
    });
  }

  async pauseExecution(runId: RunId): Promise<void> {
    await this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/pause`, {
      method: 'POST',
      headers: this.headers
    });
  }

  async resumeExecution(runId: RunId): Promise<void> {
    await this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/resume`, {
      method: 'POST',
      headers: this.headers
    });
  }

  async retryExecution(runId: RunId): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/retry`, {
      method: 'POST',
      headers: this.headers
    });
  }

  async getExecutionData(runId: RunId): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/data`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async getExecutionMetrics(runId: RunId): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/metrics`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async getExecutionDetails(runId: RunId): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/details`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async getStepExecution(runId: RunId, stepId: string): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/steps/${stepId}`, {
      method: 'GET',
      headers: this.headers
    });
  }

  async getChildStepExecution(runId: RunId, childStepId: string): Promise<any> {
    return this.makeRequest(`${this.baseUrl}/api/v1/executions/${runId}/child-steps/${childStepId}`, {
      method: 'GET',
      headers: this.headers
    });
  }

  // ============ HEALTH CHECK ============

  async ping(): Promise<void> {
    try {
      await this.makeRequest(`${this.baseUrl}/health`, {
        method: 'GET',
        headers: this.headers
      });
    } catch (error) {
      throw wrapError(error);
    }
  }

  // ============ CLIENT MANAGEMENT ============

  async close(): Promise<void> {
    // Clean up any resources if needed
    // For now, just clear references
    this.headers = {};
  }

  // ============ PRIVATE HELPER METHODS ============

  private async makeRequest<T>(url: string, options: RequestInit): Promise<T> {
    const timeoutMs = this.config.getTimeoutConfig().timeoutMs;
    const maxRetries = this.config.getTimeoutConfig().maxRetries;
    const retryDelayMs = this.config.getTimeoutConfig().retryDelayMs;

    let lastError: Error | null = null;

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

        const response = await fetch(url, {
          ...options,
          signal: controller.signal
        });

        clearTimeout(timeoutId);

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          throw this.handleHTTPError(response.status, errorData);
        }

        if (response.status === 204) {
          return undefined as T;
        }

        return await response.json() as T;
      } catch (error) {
        lastError = error as Error;

        // Don't retry on certain errors
        if (this.shouldNotRetry(error as Error)) {
          throw wrapError(error);
        }

        // If this was the last attempt, throw the error
        if (attempt === maxRetries) {
          throw wrapError(error);
        }

        // Wait before retrying
        if (retryDelayMs > 0) {
          await new Promise(resolve => setTimeout(resolve, retryDelayMs));
        }
      }
    }

    throw wrapError(lastError);
  }

  private handleHTTPError(status: number, errorData: any): Error {
    switch (status) {
      case 400:
        if (errorData.validationResult) {
          return new RequestValidationError(
            errorData.error || 'Validation failed',
            errorData.validationResult,
            errorData
          );
        }
        return new WorkflowSDKError(
          ErrorCodes.VALIDATION_FAILED,
          errorData.error || 'Bad request',
          errorData
        );
      
      case 401:
        return new WorkflowSDKError(
          ErrorCodes.AUTHENTICATION_FAILED,
          errorData.error || 'Authentication failed',
          errorData
        );
      
      case 403:
        return new WorkflowSDKError(
          'FORBIDDEN',
          errorData.error || 'Forbidden',
          errorData
        );
      
      case 404:
        if (errorData.workflowId) {
          return new WorkflowNotFoundError(
            errorData.workflowId,
            errorData.error || 'Workflow not found',
            errorData
          );
        }
        if (errorData.runId) {
          return new ExecutionNotFoundError(
            errorData.runId,
            errorData.error || 'Execution not found',
            errorData
          );
        }
        return new WorkflowSDKError(
          ErrorCodes.WORKFLOW_NOT_FOUND,
          errorData.error || 'Resource not found',
          errorData
        );
      
      case 429:
        return new WorkflowSDKError(
          ErrorCodes.RATE_LIMITED,
          errorData.error || 'Rate limited',
          { ...errorData, retryAfter: errorData.retryAfter }
        );
      
      case 500:
        return new WorkflowSDKError(
          ErrorCodes.INTERNAL_SERVER_ERROR,
          errorData.error || 'Internal server error',
          errorData
        );
      
      default:
        return new WorkflowSDKError(
          ErrorCodes.UNKNOWN_ERROR,
          errorData.error || `HTTP error ${status}`,
          errorData
        );
    }
  }

  private shouldNotRetry(error: Error): boolean {
    // Don't retry on 4xx errors (except 429)
    if (error instanceof WorkflowSDKError) {
      return [
        ErrorCodes.VALIDATION_FAILED,
        ErrorCodes.AUTHENTICATION_FAILED,
        ErrorCodes.WORKFLOW_NOT_FOUND,
        ErrorCodes.EXECUTION_NOT_FOUND
      ].includes(error.code as any);
    }
    return false;
  }

  private async extractHTTPRequestContext(request: Request): Promise<any> {
    // Extract HTTP request context from Fetch API Request object
    const headers: Record<string, string[]> = {};
    
    // Convert headers to the format expected by our API
    request.headers.forEach((value, key) => {
      if (!headers[key]) {
        headers[key] = [];
      }
      headers[key].push(value);
    });

    // Extract URL components
    const url = new URL(request.url);
    
    // Extract query parameters
    const queryParams: Record<string, string[]> = {};
    url.searchParams.forEach((value, key) => {
      if (!queryParams[key]) {
        queryParams[key] = [];
      }
      queryParams[key].push(value);
    });

    // Extract body if present
    let body: any = undefined;
    if (request.method !== 'GET' && request.method !== 'HEAD') {
      try {
        body = await request.json();
      } catch (error) {
        // If we can't parse as JSON, try as text
        try {
          body = await request.text();
        } catch (textError) {
          // Body might be empty or in unsupported format
          body = undefined;
        }
      }
    }

    return {
      method: request.method,
      path: url.pathname,
      headers,
      queryParams: Object.keys(queryParams).length > 0 ? queryParams : undefined,
      body,
      remoteAddr: request.headers.get('x-forwarded-for') || 
                  request.headers.get('x-real-ip') || 
                  'unknown',
      userAgent: request.headers.get('user-agent') || 'unknown',
      timestamp: new Date().toISOString()
    };
  }
}
