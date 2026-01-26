package handlers

import (
	"net/http"
	"strconv"
	"time"

	"unified-workflow/internal/common/model"
	"unified-workflow/internal/executor"
	"unified-workflow/internal/registry"
	"unified-workflow/internal/state"

	"github.com/gin-gonic/gin"
)

// WorkflowHandler handles workflow API requests
type WorkflowHandler struct {
	executor        executor.Executor
	registry        registry.Registry
	stateManagement state.StateManagement
}

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler(
	executor executor.Executor,
	registry registry.Registry,
	stateManagement state.StateManagement,
) *WorkflowHandler {
	return &WorkflowHandler{
		executor:        executor,
		registry:        registry,
		stateManagement: stateManagement,
	}
}

// ListWorkflows lists all registered workflows
func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all workflow IDs
	workflowIDs, err := h.registry.GetAllWorkflowIDs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list workflows",
			"details": err.Error(),
		})
		return
	}

	// Get each workflow
	response := make([]gin.H, 0, len(workflowIDs))
	for _, workflowID := range workflowIDs {
		workflow, err := h.registry.GetWorkflow(ctx, workflowID)
		if err != nil {
			// Skip workflows that can't be retrieved
			continue
		}

		response = append(response, gin.H{
			"id":          workflow.GetID(),
			"name":        workflow.GetName(),
			"description": workflow.GetDescription(),
			"step_count":  workflow.GetStepCount(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"workflows": response,
		"count":     len(response),
	})
}

// GetWorkflow gets a specific workflow by ID
func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	ctx := c.Request.Context()
	workflowID := c.Param("id")

	workflow, err := h.registry.GetWorkflow(ctx, workflowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Workflow not found",
			"details": err.Error(),
		})
		return
	}

	// Get steps
	steps := workflow.GetSteps()
	stepDetails := make([]gin.H, 0, len(steps))
	for _, step := range steps {
		stepDetails = append(stepDetails, gin.H{
			"name":             step.GetName(),
			"child_step_count": step.GetChildStepCount(),
			"is_parallel":      step.IsParallel(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          workflow.GetID(),
		"name":        workflow.GetName(),
		"description": workflow.GetDescription(),
		"step_count":  workflow.GetStepCount(),
		"steps":       stepDetails,
	})
}

// CreateWorkflow creates a new workflow
func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	ctx := c.Request.Context()

	var request struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Create workflow using the common model
	workflow := model.NewBaseWorkflow(request.Name, request.Description)

	// Register workflow
	err := h.registry.RegisterWorkflow(ctx, workflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create workflow",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          workflow.GetID(),
		"name":        workflow.GetName(),
		"description": workflow.GetDescription(),
		"message":     "Workflow created successfully",
	})
}

// DeleteWorkflow deletes a workflow
func (h *WorkflowHandler) DeleteWorkflow(c *gin.Context) {
	ctx := c.Request.Context()
	workflowID := c.Param("id")

	err := h.registry.RemoveWorkflow(ctx, workflowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to delete workflow",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workflow deleted successfully",
	})
}

// ExecuteWorkflow executes a workflow
func (h *WorkflowHandler) ExecuteWorkflow(c *gin.Context) {
	ctx := c.Request.Context()
	workflowID := c.Param("id")

	// Get workflow
	workflow, err := h.registry.GetWorkflow(ctx, workflowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Workflow not found",
			"details": err.Error(),
		})
		return
	}

	// Submit workflow for execution
	runID, err := h.executor.SubmitWorkflow(ctx, workflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to execute workflow",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"run_id":     runID,
		"message":    "Workflow execution started",
		"status_url": "/api/v1/executions/" + runID,
	})
}

// AsyncExecuteWorkflow executes a workflow asynchronously with immediate 202 response
func (h *WorkflowHandler) AsyncExecuteWorkflow(c *gin.Context) {
	ctx := c.Request.Context()
	workflowID := c.Param("id")

	// Parse request body
	var request struct {
		InputData         map[string]interface{} `json:"input_data"`
		CallbackURL       string                 `json:"callback_url"`
		TimeoutMs         int                    `json:"timeout_ms"`
		WaitForCompletion bool                   `json:"wait_for_completion"`
		Metadata          map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get workflow
	workflow, err := h.registry.GetWorkflow(ctx, workflowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Workflow not found",
			"details": err.Error(),
		})
		return
	}

	// For now, use the same executor but we'll enhance this later
	// to publish directly to NATS with response routing
	runID, err := h.executor.SubmitWorkflow(ctx, workflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to execute workflow",
			"details": err.Error(),
		})
		return
	}

	// Return 202 Accepted with polling information
	c.JSON(http.StatusAccepted, gin.H{
		"run_id":                  runID,
		"status":                  "queued",
		"message":                 "Workflow execution queued",
		"status_url":              "/api/v1/executions/" + runID,
		"result_url":              "/api/v1/executions/" + runID + "/result",
		"poll_after_ms":           1000,
		"estimated_completion_ms": 5000,
		"expires_at":              time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	})
}

// GetExecutionResult gets the result of an async workflow execution
func (h *WorkflowHandler) GetExecutionResult(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("runId")

	// Parse query parameters
	waitMsStr := c.Query("wait_ms")
	longPollStr := c.Query("long_poll")

	waitMs := 0
	if waitMsStr != "" {
		if val, err := strconv.Atoi(waitMsStr); err == nil && val >= 0 {
			waitMs = val
		}
	}

	longPoll := false
	if longPollStr == "true" {
		longPoll = true
	}

	// For now, check execution status
	// TODO: Implement proper result storage and retrieval from Redis
	status, err := h.executor.GetExecutionStatus(ctx, runID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Execution not found",
			"details": err.Error(),
		})
		return
	}

	// Check if execution is complete
	if status.Status == "completed" || status.Status == "failed" || status.Status == "cancelled" {
		// Get execution data for result
		data, err := h.executor.GetExecutionData(ctx, runID)
		if err != nil {
			data = make(map[string]interface{})
		}

		// Return completed result
		c.JSON(http.StatusOK, gin.H{
			"run_id": runID,
			"status": status.Status,
			"result": gin.H{
				"run_id":                runID,
				"workflow_id":           status.WorkflowID,
				"status":                status.Status,
				"result":                data,
				"completed_at":          time.Now().Format(time.RFC3339),
				"execution_time_millis": 0, // TODO: Calculate actual execution time
				"step_count":            0, // TODO: Get actual step count
			},
		})
		return
	}

	// Execution still in progress
	// If long polling is enabled and waitMs > 0, wait for completion
	if longPoll && waitMs > 0 {
		// Simple implementation: poll with timeout
		timeout := time.Duration(waitMs) * time.Millisecond
		startTime := time.Now()

		for time.Since(startTime) < timeout {
			status, err := h.executor.GetExecutionStatus(ctx, runID)
			if err == nil && (status.Status == "completed" || status.Status == "failed" || status.Status == "cancelled") {
				// Get result and return
				data, _ := h.executor.GetExecutionData(ctx, runID)
				c.JSON(http.StatusOK, gin.H{
					"run_id": runID,
					"status": status.Status,
					"result": gin.H{
						"run_id":                runID,
						"workflow_id":           status.WorkflowID,
						"status":                status.Status,
						"result":                data,
						"completed_at":          time.Now().Format(time.RFC3339),
						"execution_time_millis": 0,
						"step_count":            0,
					},
				})
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Return "not ready" response
	c.JSON(http.StatusAccepted, gin.H{
		"run_id":                  runID,
		"status":                  status.Status,
		"poll_after_ms":           1000,
		"estimated_completion_ms": 5000,
		"progress":                status.Progress,
	})
}

// ListExecutions lists workflow executions with optional filters
func (h *WorkflowHandler) ListExecutions(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	workflowID := c.Query("workflow_id")
	status := c.Query("status")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 50
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	offset := 0
	if offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	// Build filters
	filters := executor.ExecutionFilters{
		WorkflowID: workflowID,
		Status:     status,
		Limit:      limit,
		Offset:     offset,
	}

	// Get executions
	executions, err := h.executor.ListExecutions(ctx, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list executions",
			"details": err.Error(),
		})
		return
	}

	// Convert to response format
	response := make([]gin.H, 0, len(executions))
	for _, exec := range executions {
		response = append(response, gin.H{
			"run_id":                   exec.RunID,
			"workflow_id":              exec.WorkflowDefinitionID,
			"status":                   exec.Status,
			"current_step_index":       exec.CurrentStepIndex,
			"current_child_step_index": exec.CurrentChildStepIndex,
			"start_time":               exec.StartTime,
			"end_time":                 exec.EndTime,
			"error_message":            exec.ErrorMessage,
			"last_attempted_step":      exec.LastAttemptedStep,
			"is_terminal":              exec.IsTerminal,
			"is_running":               exec.IsRunning,
			"is_pending":               exec.IsPending,
			"created_at":               exec.CreatedAt,
			"updated_at":               exec.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"executions": response,
		"count":      len(response),
		"limit":      limit,
		"offset":     offset,
	})
}

// GetExecutionStatus gets the status of a workflow execution
func (h *WorkflowHandler) GetExecutionStatus(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("runId")

	status, err := h.executor.GetExecutionStatus(ctx, runID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Execution not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"run_id":                   status.RunID,
		"workflow_id":              status.WorkflowID,
		"status":                   status.Status,
		"current_step":             status.CurrentStep,
		"current_step_index":       status.CurrentStepIndex,
		"current_child_step_index": status.CurrentChildStepIndex,
		"progress":                 status.Progress,
		"start_time":               status.StartTime,
		"end_time":                 status.EndTime,
		"error_message":            status.ErrorMessage,
		"last_attempted_step":      status.LastAttemptedStep,
		"is_terminal":              status.IsTerminal,
		"metadata":                 status.Metadata,
	})
}

// CancelExecution cancels a running workflow execution
func (h *WorkflowHandler) CancelExecution(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("runId")

	err := h.executor.CancelExecution(ctx, runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cancel execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Execution cancelled successfully",
	})
}

// PauseExecution pauses a running workflow execution
func (h *WorkflowHandler) PauseExecution(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("runId")

	err := h.executor.PauseExecution(ctx, runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to pause execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Execution paused successfully",
	})
}

// ResumeExecution resumes a paused workflow execution
func (h *WorkflowHandler) ResumeExecution(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("runId")

	err := h.executor.ResumeExecution(ctx, runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to resume execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Execution resumed successfully",
	})
}

// RetryExecution retries a failed workflow execution
func (h *WorkflowHandler) RetryExecution(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("runId")

	err := h.executor.RetryExecution(ctx, runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retry execution",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Execution retry initiated successfully",
	})
}

// GetExecutionData gets the data of a workflow execution
func (h *WorkflowHandler) GetExecutionData(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("runId")

	data, err := h.executor.GetExecutionData(ctx, runID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to get execution data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"run_id": runID,
		"data":   data,
	})
}

// GetExecutionMetrics gets execution metrics for a workflow run
func (h *WorkflowHandler) GetExecutionMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	runID := c.Param("runId")

	metrics, err := h.executor.GetMetrics(ctx, runID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Failed to get execution metrics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"run_id":                metrics.RunID,
		"workflow_id":           metrics.WorkflowID,
		"workflow_metrics":      metrics.WorkflowMetrics,
		"step_metrics":          metrics.StepMetrics,
		"child_step_metrics":    metrics.ChildStepMetrics,
		"total_steps":           metrics.TotalSteps,
		"completed_steps":       metrics.CompletedSteps,
		"failed_steps":          metrics.FailedSteps,
		"total_child_steps":     metrics.TotalChildSteps,
		"completed_child_steps": metrics.CompletedChildSteps,
		"failed_child_steps":    metrics.FailedChildSteps,
		"total_duration_millis": metrics.TotalDurationMillis,
		"average_step_duration": metrics.AverageStepDuration,
		"success_rate":          metrics.SuccessRate,
	})
}
