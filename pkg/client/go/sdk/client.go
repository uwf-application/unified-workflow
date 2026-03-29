package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"unified-workflow/pkg/client"
	"unified-workflow/pkg/client/executor"
)

// WorkflowSDKClient is the main interface for the Workflow SDK
type WorkflowSDKClient interface {
	// Execute workflow from HTTP request
	ExecuteFromHTTPRequest(ctx context.Context, workflowID string, req *http.Request) (*SDKExecuteWorkflowResponse, error)

	// Execute workflow with raw data
	ExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}) (*SDKExecuteWorkflowResponse, error)

	// Execute workflow with full context
	ExecuteWorkflowWithContext(ctx context.Context, workflowID string, sdkReq *SDKExecuteWorkflowRequest) (*SDKExecuteWorkflowResponse, error)

	// Validate and execute workflow
	ValidateAndExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}, rules []ValidationRule) (*SDKExecuteWorkflowResponse, error)

	// Batch execute workflows
	BatchExecuteWorkflows(ctx context.Context, batchReq *BatchExecuteWorkflowsRequest) (*BatchExecuteWorkflowsResponse, error)

	// Get execution status
	GetExecutionStatus(ctx context.Context, runID string) (*executor.GetExecutionStatusResponse, error)

	// Get execution details
	GetExecutionDetails(ctx context.Context, runID string) (*executor.GetExecutionDetailsResponse, error)

	// Cancel execution
	CancelExecution(ctx context.Context, runID string) error

	// Health check
	Ping(ctx context.Context) error

	// Close the client
	Close() error
}

// workflowSDKClient implements the WorkflowSDKClient interface
type workflowSDKClient struct {
	config     *SDKConfig
	parser     *RequestParser
	validator  *Validator
	httpClient *client.HTTPClient
	executor   executor.Client
}

// NewClient creates a new Workflow SDK client
func NewClient(config *SDKConfig) (WorkflowSDKClient, error) {
	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Create HTTP client configuration
	httpConfig := client.Config{
		Endpoint:                config.WorkflowAPIEndpoint,
		Timeout:                 config.Timeout,
		MaxRetries:              config.MaxRetries,
		RetryDelay:              config.RetryDelay,
		AuthToken:               config.AuthToken,
		EnableCircuitBreaker:    config.EnableCircuitBreaker,
		CircuitBreakerThreshold: config.CircuitBreakerThreshold,
		CircuitBreakerTimeout:   config.CircuitBreakerTimeout,
	}

	// Create HTTP client
	httpClient := client.NewHTTPClient(httpConfig)

	// Create parser and validator
	parser := NewRequestParser(config)
	validator := NewValidator(config)

	// Create executor client (using existing executor client interface)
	// Note: We need to implement the actual executor client integration
	executor := newExecutorClient(httpClient, config)

	return &workflowSDKClient{
		config:     config,
		parser:     parser,
		validator:  validator,
		httpClient: httpClient,
		executor:   executor,
	}, nil
}

// ExecuteFromHTTPRequest executes a workflow from an HTTP request
func (c *workflowSDKClient) ExecuteFromHTTPRequest(ctx context.Context, workflowID string, req *http.Request) (*SDKExecuteWorkflowResponse, error) {
	// Parse HTTP request into SDK execution request
	sdkReq, err := c.parser.CreateSDKExecuteRequest(ctx, req, workflowID)
	if err != nil {
		return nil, WrapSDKError(err, ErrCodeRequestParsingFailed, "Failed to parse HTTP request")
	}

	// Execute with full context
	return c.ExecuteWorkflowWithContext(ctx, workflowID, sdkReq)
}

// ExecuteWorkflow executes a workflow with raw data
func (c *workflowSDKClient) ExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}) (*SDKExecuteWorkflowResponse, error) {
	// Create SDK execution request
	sdkReq := NewSDKExecuteWorkflowRequest(data)
	sdkReq.EnableValidation = c.config.EnableValidation
	sdkReq.EnableSanitization = c.config.EnableSanitization
	sdkReq.ValidationRules = c.config.DefaultValidationRules

	// Execute with context
	return c.ExecuteWorkflowWithContext(ctx, workflowID, sdkReq)
}

// ExecuteWorkflowWithContext executes a workflow with full SDK context
func (c *workflowSDKClient) ExecuteWorkflowWithContext(ctx context.Context, workflowID string, sdkReq *SDKExecuteWorkflowRequest) (*SDKExecuteWorkflowResponse, error) {
	// Validate input data if enabled
	if sdkReq.EnableValidation {
		validationResult := c.validator.Validate(sdkReq.InputData, sdkReq.ValidationRules)
		if !validationResult.Valid && c.config.StrictValidation {
			return nil, NewSDKErrorWithDetails(ErrCodeValidationFailed,
				"Validation failed",
				map[string]interface{}{"validation_result": validationResult})
		}

		// Use sanitized data if validation passed
		if validationResult.Valid && c.config.EnableSanitization {
			sdkReq.InputData = validationResult.SanitizedData
		}

		// Store validation result
		sdkReq.Metadata["validation_result"] = validationResult
	}

	// Validate HTTP request context if present
	if sdkReq.HTTPRequest != nil {
		httpValidationResult := c.validator.ValidateHTTPRequest(sdkReq.HTTPRequest)
		if httpValidationResult.HasErrors() {
			sdkReq.Metadata["http_validation_errors"] = httpValidationResult.Errors
		}
		if httpValidationResult.HasWarnings() {
			sdkReq.Metadata["http_validation_warnings"] = httpValidationResult.Warnings
		}
	}

	// Validate session context if present
	if sdkReq.Session != nil {
		sessionValidationResult := c.validator.ValidateSessionContext(sdkReq.Session)
		if sessionValidationResult.HasErrors() {
			sdkReq.Metadata["session_validation_errors"] = sessionValidationResult.Errors
		}
	}

	// Prepare execution request for workflow API
	executionReq := &executor.ExecuteWorkflowRequest{
		Request: client.Request{
			ID:        generateRequestID(),
			Timestamp: time.Now(),
			Metadata:  make(map[string]string),
		},
		WorkflowID: workflowID,
		InputData:  sdkReq.InputData,
		Async:      c.config.AsyncExecution,
		TimeoutMs:  sdkReq.TimeoutMs,
		Priority:   c.config.DefaultPriority,
		Metadata:   sdkReq.Metadata,
	}

	// Add SDK context to metadata
	if sdkReq.IncludeFullContext {
		executionReq.Metadata["sdk_context"] = map[string]interface{}{
			"http_request": sdkReq.HTTPRequest,
			"session":      sdkReq.Session,
			"security":     sdkReq.Security,
			"sdk_version":  c.config.SDKVersion,
		}
	}

	// Execute workflow
	executionResp, err := c.executor.ExecuteWorkflow(ctx, executionReq)
	if err != nil {
		return nil, WrapSDKError(err, ErrCodeWorkflowExecution, "Failed to execute workflow")
	}

	// Convert to SDK response
	sdkResp := &SDKExecuteWorkflowResponse{
		RunID:                 executionResp.RunID,
		Status:                executionResp.Status,
		Message:               "Workflow execution started",
		StatusURL:             executionResp.StatusURL,
		ResultURL:             executionResp.ResultURL,
		PollAfterMs:           c.config.PollIntervalMs,
		EstimatedCompletionMs: c.config.EstimatedCompletionMs,
		ExpiresAt:             time.Now().Add(c.config.ExecutionExpiryDuration),
		ContextIncluded:       sdkReq.IncludeFullContext,
		SDKVersion:            c.config.SDKVersion,
		RequestID:             executionReq.ID,
	}

	// Add validation result if available
	if validationResult, ok := sdkReq.Metadata["validation_result"].(*ValidationResult); ok {
		sdkResp.ValidationResult = validationResult
	}

	return sdkResp, nil
}

// ValidateAndExecuteWorkflow validates input data before execution
func (c *workflowSDKClient) ValidateAndExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}, rules []ValidationRule) (*SDKExecuteWorkflowResponse, error) {
	// Validate data
	validationResult := c.validator.Validate(data, rules)
	if !validationResult.Valid && c.config.StrictValidation {
		return nil, NewSDKErrorWithDetails(ErrCodeValidationFailed,
			"Validation failed",
			map[string]interface{}{"validation_result": validationResult})
	}

	// Create SDK execution request
	sdkReq := NewSDKExecuteWorkflowRequest(validationResult.SanitizedData)
	sdkReq.ValidationRules = rules
	sdkReq.Metadata["validation_result"] = validationResult

	// Execute with context
	return c.ExecuteWorkflowWithContext(ctx, workflowID, sdkReq)
}

// BatchExecuteWorkflows executes multiple workflows in batch
func (c *workflowSDKClient) BatchExecuteWorkflows(ctx context.Context, batchReq *BatchExecuteWorkflowsRequest) (*BatchExecuteWorkflowsResponse, error) {
	// TODO: Implement batch execution
	// This would involve executing multiple workflows in parallel or sequentially
	// based on the batch configuration

	return nil, NewSDKError(ErrCodeWorkflowExecution, "Batch execution not yet implemented")
}

// GetExecutionStatus gets the status of a workflow execution
func (c *workflowSDKClient) GetExecutionStatus(ctx context.Context, runID string) (*executor.GetExecutionStatusResponse, error) {
	req := &executor.GetExecutionStatusRequest{
		Request: client.Request{
			ID:        generateRequestID(),
			Timestamp: time.Now(),
		},
		RunID: runID,
	}

	return c.executor.GetExecutionStatus(ctx, req)
}

// GetExecutionDetails gets detailed execution information
func (c *workflowSDKClient) GetExecutionDetails(ctx context.Context, runID string) (*executor.GetExecutionDetailsResponse, error) {
	req := &executor.GetExecutionDetailsRequest{
		Request: client.Request{
			ID:        generateRequestID(),
			Timestamp: time.Now(),
		},
		RunID: runID,
	}

	return c.executor.GetExecutionDetails(ctx, req)
}

// CancelExecution cancels a workflow execution
func (c *workflowSDKClient) CancelExecution(ctx context.Context, runID string) error {
	req := &executor.CancelExecutionRequest{
		Request: client.Request{
			ID:        generateRequestID(),
			Timestamp: time.Now(),
		},
		RunID: runID,
	}

	_, err := c.executor.CancelExecution(ctx, req)
	return err
}

// Ping performs a health check
func (c *workflowSDKClient) Ping(ctx context.Context) error {
	return c.httpClient.Ping(ctx)
}

// Close closes the client
func (c *workflowSDKClient) Close() error {
	return c.httpClient.Close()
}

// Helper functions

func generateRequestID() string {
	return "req_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().Nanosecond()%len(charset)]
	}
	return string(b)
}

// executorClient is a wrapper around the existing executor client
type executorClient struct {
	httpClient *client.HTTPClient
	config     *SDKConfig
	endpoint   string
}

func newExecutorClient(httpClient *client.HTTPClient, config *SDKConfig) *executorClient {
	return &executorClient{
		httpClient: httpClient,
		config:     config,
		endpoint:   config.WorkflowAPIEndpoint,
	}
}

// Ping checks if the service is reachable
func (ec *executorClient) Ping(ctx context.Context) error {
	return ec.httpClient.Ping(ctx)
}

// GetEndpoint returns the service endpoint
func (ec *executorClient) GetEndpoint() string {
	return ec.endpoint
}

// IsHealthy checks if the client is healthy
func (ec *executorClient) IsHealthy() bool {
	return ec.httpClient.IsHealthy()
}

// ExecuteWorkflow executes a workflow
func (ec *executorClient) ExecuteWorkflow(ctx context.Context, req *executor.ExecuteWorkflowRequest) (*executor.ExecuteWorkflowResponse, error) {
	// Make actual HTTP call to workflow API
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/workflows/"+req.WorkflowID+"/execute", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var executionResp executor.ExecuteWorkflowResponse
	if err := ec.httpClient.ParseResponse(resp, &executionResp); err != nil {
		return nil, err
	}

	return &executionResp, nil
}

// GetExecutionStatus gets the status of a workflow execution
func (ec *executorClient) GetExecutionStatus(ctx context.Context, req *executor.GetExecutionStatusRequest) (*executor.GetExecutionStatusResponse, error) {
	// Make actual HTTP call to workflow API
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/status", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var statusResp executor.GetExecutionStatusResponse
	if err := ec.httpClient.ParseResponse(resp, &statusResp); err != nil {
		return nil, err
	}

	return &statusResp, nil
}

// GetExecutionDetails gets detailed execution information
func (ec *executorClient) GetExecutionDetails(ctx context.Context, req *executor.GetExecutionDetailsRequest) (*executor.GetExecutionDetailsResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/details", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var detailsResp executor.GetExecutionDetailsResponse
	if err := ec.httpClient.ParseResponse(resp, &detailsResp); err != nil {
		return nil, err
	}

	return &detailsResp, nil
}

// CancelExecution cancels a running workflow execution
func (ec *executorClient) CancelExecution(ctx context.Context, req *executor.CancelExecutionRequest) (*executor.CancelExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/executions/"+req.RunID+"/cancel", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cancelResp executor.CancelExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &cancelResp); err != nil {
		return nil, err
	}

	return &cancelResp, nil
}

// PauseExecution pauses a running workflow execution
func (ec *executorClient) PauseExecution(ctx context.Context, req *executor.PauseExecutionRequest) (*executor.PauseExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/executions/"+req.RunID+"/pause", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pauseResp executor.PauseExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &pauseResp); err != nil {
		return nil, err
	}

	return &pauseResp, nil
}

// ResumeExecution resumes a paused workflow execution
func (ec *executorClient) ResumeExecution(ctx context.Context, req *executor.ResumeExecutionRequest) (*executor.ResumeExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/executions/"+req.RunID+"/resume", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var resumeResp executor.ResumeExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &resumeResp); err != nil {
		return nil, err
	}

	return &resumeResp, nil
}

// RetryExecution retries a failed workflow execution
func (ec *executorClient) RetryExecution(ctx context.Context, req *executor.RetryExecutionRequest) (*executor.RetryExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/executions/"+req.RunID+"/retry", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var retryResp executor.RetryExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &retryResp); err != nil {
		return nil, err
	}

	return &retryResp, nil
}

// GetExecutionData gets the data of a workflow execution
func (ec *executorClient) GetExecutionData(ctx context.Context, req *executor.GetExecutionDataRequest) (*executor.GetExecutionDataResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/data", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dataResp executor.GetExecutionDataResponse
	if err := ec.httpClient.ParseResponse(resp, &dataResp); err != nil {
		return nil, err
	}

	return &dataResp, nil
}

// ListExecutions lists workflow executions with optional filters
func (ec *executorClient) ListExecutions(ctx context.Context, req *executor.ListExecutionsRequest) (*executor.ListExecutionsResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var listResp executor.ListExecutionsResponse
	if err := ec.httpClient.ParseResponse(resp, &listResp); err != nil {
		return nil, err
	}

	return &listResp, nil
}

// GetExecutionMetrics gets execution metrics for a workflow run
func (ec *executorClient) GetExecutionMetrics(ctx context.Context, req *executor.GetExecutionMetricsRequest) (*executor.GetExecutionMetricsResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/metrics", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metricsResp executor.GetExecutionMetricsResponse
	if err := ec.httpClient.ParseResponse(resp, &metricsResp); err != nil {
		return nil, err
	}

	return &metricsResp, nil
}

// GetStepExecution gets step execution details
func (ec *executorClient) GetStepExecution(ctx context.Context, req *executor.GetStepExecutionRequest) (*executor.GetStepExecutionResponse, error) {
	// TODO: Implement actual HTTP call
	return &executor.GetStepExecutionResponse{
		Response: client.Response{
			ID:        generateRequestID(),
			Timestamp: time.Now(),
			Success:   true,
		},
		StepExecution: &executor.StepExecution{},
	}, nil
}

// GetChildStepExecution gets child step execution details
func (ec *executorClient) GetChildStepExecution(ctx context.Context, req *executor.GetChildStepExecutionRequest) (*executor.GetChildStepExecutionResponse, error) {
	// TODO: Implement actual HTTP call
	return &executor.GetChildStepExecutionResponse{
		Response: client.Response{
			ID:        generateRequestID(),
			Timestamp: time.Now(),
			Success:   true,
		},
		ChildStepExecution: &executor.ChildStepExecution{},
	}, nil
}

// Close closes the client
func (ec *executorClient) Close() error {
	return ec.httpClient.Close()
}
