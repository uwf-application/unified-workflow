package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"unified-workflow/internal/common/model"
	registryClient "unified-workflow/pkg/client/registry"
)

// HTTPRegistry is an HTTP implementation of the Registry interface
type HTTPRegistry struct {
	client registryClient.Client
}

// NewHTTPRegistry creates a new HTTP registry
func NewHTTPRegistry(endpoint string) (*HTTPRegistry, error) {
	log.Printf("[HTTPRegistry] Creating HTTP registry for endpoint: %s", endpoint)
	config := registryClient.DefaultConfig()
	config.Endpoint = endpoint

	client := registryClient.NewHTTPClient(config)
	log.Printf("[HTTPRegistry] HTTP registry created successfully")
	return &HTTPRegistry{
		client: client,
	}, nil
}

// RegisterWorkflow registers a workflow with the registry
func (r *HTTPRegistry) RegisterWorkflow(ctx context.Context, workflow model.Workflow) error {
	req := &registryClient.CreateWorkflowRequest{
		Name:        workflow.GetName(),
		Description: workflow.GetDescription(),
		// Note: We need to convert model.Step to the appropriate format
		// For now, we'll create a simple workflow without steps
	}

	resp, err := r.client.CreateWorkflow(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register workflow: %w", err)
	}

	// Update the workflow ID if needed
	if workflow.GetID() != resp.Workflow.ID {
		// Note: We would need a way to update the workflow ID
		// This depends on how the model.Workflow interface is implemented
	}

	return nil
}

// GetWorkflow retrieves a workflow by its ID
func (r *HTTPRegistry) GetWorkflow(ctx context.Context, workflowID string) (model.Workflow, error) {
	log.Printf("[HTTPRegistry] GetWorkflow called for workflow ID: %s", workflowID)
	// Always use direct HTTP request since the HTTP client expects a different response format
	return r.getWorkflowDirect(ctx, workflowID)
}

// getWorkflowDirect makes a direct HTTP request to get workflow
func (r *HTTPRegistry) getWorkflowDirect(ctx context.Context, workflowID string) (model.Workflow, error) {
	endpoint := r.client.GetEndpoint()
	url := fmt.Sprintf("%s/api/v1/workflows/%s", endpoint, workflowID)

	// Log the URL for debugging
	log.Printf("[HTTPRegistry] Getting workflow from URL: %s", url)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("[HTTPRegistry] Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[HTTPRegistry] Failed to make request to %s: %v", url, err)
		return nil, fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read response body for error details
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[HTTPRegistry] Registry returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("registry returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[HTTPRegistry] Failed to parse response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract workflow info
	name, _ := result["name"].(string)
	description, _ := result["description"].(string)
	id, _ := result["id"].(string)

	log.Printf("[HTTPRegistry] Successfully retrieved workflow: %s (ID: %s)", name, id)

	workflow := model.NewBaseWorkflow(name, description)
	workflow.ID = id

	return workflow, nil
}

// ContainsWorkflow checks if a workflow exists in the registry
func (r *HTTPRegistry) ContainsWorkflow(ctx context.Context, workflowID string) (bool, error) {
	req := &registryClient.ContainsWorkflowRequest{
		ID: workflowID,
	}

	resp, err := r.client.ContainsWorkflow(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to check if workflow exists: %w", err)
	}

	return resp.Exists, nil
}

// RemoveWorkflow removes a workflow from the registry
func (r *HTTPRegistry) RemoveWorkflow(ctx context.Context, workflowID string) error {
	req := &registryClient.DeleteWorkflowRequest{
		ID: workflowID,
	}

	_, err := r.client.DeleteWorkflow(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to remove workflow: %w", err)
	}

	return nil
}

// GetAllWorkflowIDs gets all registered workflow IDs
func (r *HTTPRegistry) GetAllWorkflowIDs(ctx context.Context) ([]string, error) {
	req := &registryClient.ListWorkflowsRequest{}

	resp, err := r.client.ListWorkflows(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	ids := make([]string, 0, len(resp.Workflows))
	for _, workflow := range resp.Workflows {
		ids = append(ids, workflow.ID)
	}

	return ids, nil
}

// GetWorkflowCount gets the number of registered workflows
func (r *HTTPRegistry) GetWorkflowCount(ctx context.Context) (int, error) {
	req := &registryClient.GetWorkflowCountRequest{}

	resp, err := r.client.GetWorkflowCount(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to get workflow count: %w", err)
	}

	return resp.Count, nil
}

// Clear removes all workflows from the registry
func (r *HTTPRegistry) Clear(ctx context.Context) error {
	req := &registryClient.ClearRequest{}

	_, err := r.client.Clear(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to clear registry: %w", err)
	}

	return nil
}

// Initialize initializes the registry (e.g., creates tables, loads data)
func (r *HTTPRegistry) Initialize(ctx context.Context) error {
	// Nothing to initialize for HTTP registry
	return nil
}

// Shutdown shuts down the registry (e.g., closes connections)
func (r *HTTPRegistry) Shutdown(ctx context.Context) error {
	return r.client.Close()
}
