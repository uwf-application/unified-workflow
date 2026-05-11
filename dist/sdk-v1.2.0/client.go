package sdk

import (
	"context"
	"net/http"
	"time"

	"github.com/uwf-application/unified-workflow-sdk/baseclient"
	"github.com/uwf-application/unified-workflow-sdk/executor"
)

// WorkflowSDKClient is the main interface for the Workflow SDK
type WorkflowSDKClient interface {
	ExecuteFromHTTPRequest(ctx context.Context, workflowID string, req *http.Request) (*SDKExecuteWorkflowResponse, error)
	ExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}) (*SDKExecuteWorkflowResponse, error)
	ExecuteWorkflowWithContext(ctx context.Context, workflowID string, sdkReq *SDKExecuteWorkflowRequest) (*SDKExecuteWorkflowResponse, error)
	ValidateAndExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}, rules []ValidationRule) (*SDKExecuteWorkflowResponse, error)
	BatchExecuteWorkflows(ctx context.Context, batchReq *BatchExecuteWorkflowsRequest) (*BatchExecuteWorkflowsResponse, error)
	GetExecutionStatus(ctx context.Context, runID string) (*executor.GetExecutionStatusResponse, error)
	GetExecutionDetails(ctx context.Context, runID string) (*executor.GetExecutionDetailsResponse, error)
	CancelExecution(ctx context.Context, runID string) error
	Ping(ctx context.Context) error
	Close() error
}

type workflowSDKClient struct {
	config     *SDKConfig
	parser     *RequestParser
	validator  *Validator
	httpClient *baseclient.HTTPClient
	executor   executor.Client
}

// NewClient creates a new Workflow SDK client
func NewClient(config *SDKConfig) (WorkflowSDKClient, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	httpConfig := baseclient.Config{
		Endpoint:                config.WorkflowAPIEndpoint,
		Timeout:                 config.Timeout,
		MaxRetries:              config.MaxRetries,
		RetryDelay:              config.RetryDelay,
		AuthToken:               config.AuthToken,
		EnableCircuitBreaker:    config.EnableCircuitBreaker,
		CircuitBreakerThreshold: config.CircuitBreakerThreshold,
		CircuitBreakerTimeout:   config.CircuitBreakerTimeout,
	}

	httpClient := baseclient.NewHTTPClient(httpConfig)
	parser := NewRequestParser(config)
	validator := NewValidator(config)
	exec := newExecutorClient(httpClient, config)

	return &workflowSDKClient{
		config:     config,
		parser:     parser,
		validator:  validator,
		httpClient: httpClient,
		executor:   exec,
	}, nil
}

func (c *workflowSDKClient) ExecuteFromHTTPRequest(ctx context.Context, workflowID string, req *http.Request) (*SDKExecuteWorkflowResponse, error) {
	sdkReq, err := c.parser.CreateSDKExecuteRequest(ctx, req, workflowID)
	if err != nil {
		return nil, WrapSDKError(err, ErrCodeRequestParsingFailed, "Failed to parse HTTP request")
	}
	return c.ExecuteWorkflowWithContext(ctx, workflowID, sdkReq)
}

func (c *workflowSDKClient) ExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}) (*SDKExecuteWorkflowResponse, error) {
	sdkReq := NewSDKExecuteWorkflowRequest(data)
	sdkReq.EnableValidation = c.config.EnableValidation
	sdkReq.EnableSanitization = c.config.EnableSanitization
	sdkReq.ValidationRules = c.config.DefaultValidationRules
	return c.ExecuteWorkflowWithContext(ctx, workflowID, sdkReq)
}

func (c *workflowSDKClient) ExecuteWorkflowWithContext(ctx context.Context, workflowID string, sdkReq *SDKExecuteWorkflowRequest) (*SDKExecuteWorkflowResponse, error) {
	if sdkReq.EnableValidation {
		validationResult := c.validator.Validate(sdkReq.InputData, sdkReq.ValidationRules)
		if !validationResult.Valid && c.config.StrictValidation {
			return nil, NewSDKErrorWithDetails(ErrCodeValidationFailed, "Validation failed",
				map[string]interface{}{"validation_result": validationResult})
		}
		if validationResult.Valid && c.config.EnableSanitization {
			sdkReq.InputData = validationResult.SanitizedData
		}
		sdkReq.Metadata["validation_result"] = validationResult
	}

	if sdkReq.HTTPRequest != nil {
		httpVal := c.validator.ValidateHTTPRequest(sdkReq.HTTPRequest)
		if httpVal.HasErrors() {
			sdkReq.Metadata["http_validation_errors"] = httpVal.Errors
		}
	}

	if sdkReq.Session != nil {
		sessionVal := c.validator.ValidateSessionContext(sdkReq.Session)
		if sessionVal.HasErrors() {
			sdkReq.Metadata["session_validation_errors"] = sessionVal.Errors
		}
	}

	executionReq := &executor.ExecuteWorkflowRequest{
		Request: baseclient.Request{
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

	executionResp, err := c.executor.ExecuteWorkflow(ctx, executionReq)
	if err != nil {
		return nil, WrapSDKError(err, ErrCodeWorkflowExecution, "Failed to execute workflow")
	}

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

	if validationResult, ok := sdkReq.Metadata["validation_result"].(*ValidationResult); ok {
		sdkResp.ValidationResult = validationResult
	}

	return sdkResp, nil
}

func (c *workflowSDKClient) ValidateAndExecuteWorkflow(ctx context.Context, workflowID string, data map[string]interface{}, rules []ValidationRule) (*SDKExecuteWorkflowResponse, error) {
	validationResult := c.validator.Validate(data, rules)
	if !validationResult.Valid && c.config.StrictValidation {
		return nil, NewSDKErrorWithDetails(ErrCodeValidationFailed, "Validation failed",
			map[string]interface{}{"validation_result": validationResult})
	}
	sdkReq := NewSDKExecuteWorkflowRequest(validationResult.SanitizedData)
	sdkReq.ValidationRules = rules
	sdkReq.Metadata["validation_result"] = validationResult
	return c.ExecuteWorkflowWithContext(ctx, workflowID, sdkReq)
}

func (c *workflowSDKClient) BatchExecuteWorkflows(_ context.Context, _ *BatchExecuteWorkflowsRequest) (*BatchExecuteWorkflowsResponse, error) {
	return nil, NewSDKError(ErrCodeWorkflowExecution, "Batch execution not yet implemented")
}

func (c *workflowSDKClient) GetExecutionStatus(ctx context.Context, runID string) (*executor.GetExecutionStatusResponse, error) {
	req := &executor.GetExecutionStatusRequest{
		Request: baseclient.Request{ID: generateRequestID(), Timestamp: time.Now()},
		RunID:   runID,
	}
	return c.executor.GetExecutionStatus(ctx, req)
}

func (c *workflowSDKClient) GetExecutionDetails(ctx context.Context, runID string) (*executor.GetExecutionDetailsResponse, error) {
	req := &executor.GetExecutionDetailsRequest{
		Request: baseclient.Request{ID: generateRequestID(), Timestamp: time.Now()},
		RunID:   runID,
	}
	return c.executor.GetExecutionDetails(ctx, req)
}

func (c *workflowSDKClient) CancelExecution(ctx context.Context, runID string) error {
	req := &executor.CancelExecutionRequest{
		Request: baseclient.Request{ID: generateRequestID(), Timestamp: time.Now()},
		RunID:   runID,
	}
	_, err := c.executor.CancelExecution(ctx, req)
	return err
}

func (c *workflowSDKClient) Ping(ctx context.Context) error { return c.httpClient.Ping(ctx) }
func (c *workflowSDKClient) Close() error                   { return c.httpClient.Close() }

// ---- internal executor adapter ----

type executorClient struct {
	httpClient *baseclient.HTTPClient
	config     *SDKConfig
}

func newExecutorClient(httpClient *baseclient.HTTPClient, config *SDKConfig) *executorClient {
	return &executorClient{httpClient: httpClient, config: config}
}

func (ec *executorClient) Ping(ctx context.Context) error { return ec.httpClient.Ping(ctx) }
func (ec *executorClient) GetEndpoint() string            { return ec.config.WorkflowAPIEndpoint }
func (ec *executorClient) IsHealthy() bool                { return ec.httpClient.IsHealthy() }
func (ec *executorClient) Close() error                   { return ec.httpClient.Close() }

func (ec *executorClient) ExecuteWorkflow(ctx context.Context, req *executor.ExecuteWorkflowRequest) (*executor.ExecuteWorkflowResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/workflows/"+req.WorkflowID+"/execute", req)
	if err != nil {
		return nil, err
	}
	var result executor.ExecuteWorkflowResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) GetExecutionStatus(ctx context.Context, req *executor.GetExecutionStatusRequest) (*executor.GetExecutionStatusResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/status", nil)
	if err != nil {
		return nil, err
	}
	var result executor.GetExecutionStatusResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) GetExecutionDetails(ctx context.Context, req *executor.GetExecutionDetailsRequest) (*executor.GetExecutionDetailsResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID, nil)
	if err != nil {
		return nil, err
	}
	var result executor.GetExecutionDetailsResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) CancelExecution(ctx context.Context, req *executor.CancelExecutionRequest) (*executor.CancelExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/executions/"+req.RunID+"/cancel", req)
	if err != nil {
		return nil, err
	}
	var result executor.CancelExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) PauseExecution(ctx context.Context, req *executor.PauseExecutionRequest) (*executor.PauseExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/executions/"+req.RunID+"/pause", req)
	if err != nil {
		return nil, err
	}
	var result executor.PauseExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) ResumeExecution(ctx context.Context, req *executor.ResumeExecutionRequest) (*executor.ResumeExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/executions/"+req.RunID+"/resume", req)
	if err != nil {
		return nil, err
	}
	var result executor.ResumeExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) RetryExecution(ctx context.Context, req *executor.RetryExecutionRequest) (*executor.RetryExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "POST", "/api/v1/executions/"+req.RunID+"/retry", req)
	if err != nil {
		return nil, err
	}
	var result executor.RetryExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) GetExecutionData(ctx context.Context, req *executor.GetExecutionDataRequest) (*executor.GetExecutionDataResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/data", nil)
	if err != nil {
		return nil, err
	}
	var result executor.GetExecutionDataResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) ListExecutions(ctx context.Context, req *executor.ListExecutionsRequest) (*executor.ListExecutionsResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions", req)
	if err != nil {
		return nil, err
	}
	var result executor.ListExecutionsResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) GetExecutionMetrics(ctx context.Context, req *executor.GetExecutionMetricsRequest) (*executor.GetExecutionMetricsResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/metrics", nil)
	if err != nil {
		return nil, err
	}
	var result executor.GetExecutionMetricsResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) GetStepExecution(ctx context.Context, req *executor.GetStepExecutionRequest) (*executor.GetStepExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/steps", nil)
	if err != nil {
		return nil, err
	}
	var result executor.GetStepExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (ec *executorClient) GetChildStepExecution(ctx context.Context, req *executor.GetChildStepExecutionRequest) (*executor.GetChildStepExecutionResponse, error) {
	resp, err := ec.httpClient.DoRequest(ctx, "GET", "/api/v1/executions/"+req.RunID+"/steps/child", nil)
	if err != nil {
		return nil, err
	}
	var result executor.GetChildStepExecutionResponse
	if err := ec.httpClient.ParseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

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
