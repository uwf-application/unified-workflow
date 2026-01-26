package registry

import (
	"context"
	"fmt"

	"unified-workflow/pkg/client"
)

// HTTPClient is an HTTP implementation of the Registry client
type HTTPClient struct {
	httpClient *client.HTTPClient
	config     Config
}

// NewHTTPClient creates a new HTTP registry client
func NewHTTPClient(config Config) *HTTPClient {
	baseConfig := client.Config{
		Endpoint:                config.Endpoint,
		Timeout:                 config.Timeout,
		MaxRetries:              config.MaxRetries,
		RetryDelay:              config.RetryDelay,
		EnableTLS:               config.EnableTLS,
		TLSCertPath:             config.TLSCertPath,
		TLSKeyPath:              config.TLSKeyPath,
		TLSCAPath:               config.TLSCAPath,
		AuthToken:               config.AuthToken,
		EnableCircuitBreaker:    config.EnableCircuitBreaker,
		CircuitBreakerThreshold: config.CircuitBreakerThreshold,
		CircuitBreakerTimeout:   config.CircuitBreakerTimeout,
	}

	return &HTTPClient{
		httpClient: client.NewHTTPClient(baseConfig),
		config:     config,
	}
}

// ListWorkflows lists all registered workflows
func (c *HTTPClient) ListWorkflows(ctx context.Context, req *ListWorkflowsRequest) (*ListWorkflowsResponse, error) {
	path := "/api/v1/workflows"
	if req.NameFilter != "" || req.DescriptionFilter != "" || req.Limit > 0 || req.Offset > 0 {
		path += "?"
		params := []string{}
		if req.NameFilter != "" {
			params = append(params, fmt.Sprintf("name=%s", req.NameFilter))
		}
		if req.DescriptionFilter != "" {
			params = append(params, fmt.Sprintf("description=%s", req.DescriptionFilter))
		}
		if req.Limit > 0 {
			params = append(params, fmt.Sprintf("limit=%d", req.Limit))
		}
		if req.Offset > 0 {
			params = append(params, fmt.Sprintf("offset=%d", req.Offset))
		}
		path += stringJoin(params, "&")
	}

	resp, err := c.httpClient.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response ListWorkflowsResponse
	if err := c.httpClient.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetWorkflow gets a specific workflow by ID
func (c *HTTPClient) GetWorkflow(ctx context.Context, req *GetWorkflowRequest) (*GetWorkflowResponse, error) {
	path := fmt.Sprintf("/api/v1/workflows/%s", req.ID)

	resp, err := c.httpClient.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response GetWorkflowResponse
	if err := c.httpClient.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateWorkflow creates a new workflow
func (c *HTTPClient) CreateWorkflow(ctx context.Context, req *CreateWorkflowRequest) (*CreateWorkflowResponse, error) {
	path := "/api/v1/workflows"

	resp, err := c.httpClient.DoRequest(ctx, "POST", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response CreateWorkflowResponse
	if err := c.httpClient.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateWorkflow updates an existing workflow
func (c *HTTPClient) UpdateWorkflow(ctx context.Context, req *UpdateWorkflowRequest) (*UpdateWorkflowResponse, error) {
	path := fmt.Sprintf("/api/v1/workflows/%s", req.ID)

	resp, err := c.httpClient.DoRequest(ctx, "PUT", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response UpdateWorkflowResponse
	if err := c.httpClient.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteWorkflow deletes a workflow
func (c *HTTPClient) DeleteWorkflow(ctx context.Context, req *DeleteWorkflowRequest) (*DeleteWorkflowResponse, error) {
	path := fmt.Sprintf("/api/v1/workflows/%s", req.ID)

	resp, err := c.httpClient.DoRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response DeleteWorkflowResponse
	if err := c.httpClient.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ContainsWorkflow checks if a workflow exists
func (c *HTTPClient) ContainsWorkflow(ctx context.Context, req *ContainsWorkflowRequest) (*ContainsWorkflowResponse, error) {
	path := fmt.Sprintf("/api/v1/workflows/%s/exists", req.ID)

	resp, err := c.httpClient.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response ContainsWorkflowResponse
	if err := c.httpClient.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetWorkflowCount gets the number of registered workflows
func (c *HTTPClient) GetWorkflowCount(ctx context.Context, req *GetWorkflowCountRequest) (*GetWorkflowCountResponse, error) {
	path := "/api/v1/workflows/count"

	resp, err := c.httpClient.DoRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response GetWorkflowCountResponse
	if err := c.httpClient.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Clear removes all workflows from the registry
func (c *HTTPClient) Clear(ctx context.Context, req *ClearRequest) (*ClearResponse, error) {
	path := "/api/v1/workflows"

	resp, err := c.httpClient.DoRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response ClearResponse
	if err := c.httpClient.ParseResponse(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Ping checks if the service is reachable
func (c *HTTPClient) Ping(ctx context.Context) error {
	return c.httpClient.Ping(ctx)
}

// Close closes the client connection
func (c *HTTPClient) Close() error {
	return c.httpClient.Close()
}

// GetEndpoint returns the service endpoint
func (c *HTTPClient) GetEndpoint() string {
	return c.httpClient.GetEndpoint()
}

// IsHealthy checks if the client is healthy
func (c *HTTPClient) IsHealthy() bool {
	return c.httpClient.IsHealthy()
}

// Helper function to join strings
func stringJoin(elems []string, sep string) string {
	result := ""
	for i, elem := range elems {
		if i > 0 {
			result += sep
		}
		result += elem
	}
	return result
}
