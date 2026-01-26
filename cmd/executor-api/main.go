package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"unified-workflow/internal/config"
	"unified-workflow/internal/di"
	"unified-workflow/internal/executor"
	"unified-workflow/internal/primitive"
	"unified-workflow/internal/queue"
	"unified-workflow/internal/registry"
	"unified-workflow/internal/state"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("=== Executor Service with DI/IoC Framework ===")
	fmt.Println("Location: Bank Data Center")
	fmt.Println("Purpose: Workflow execution engine with dependency injection")
	fmt.Println("Traces: All transaction traces stay in Bank DC")
	fmt.Println("DI Framework: High-performance container for 5000+ TPS")
	fmt.Println("")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize DI container
	container, err := initializeContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}
	defer container.Stop()

	log.Println("DI container initialized successfully")

	// Resolve dependencies from container
	executorService, err := resolveExecutorService(container)
	if err != nil {
		log.Fatalf("Failed to resolve executor service: %v", err)
	}

	registryService, err := resolveRegistryService(container)
	if err != nil {
		log.Fatalf("Failed to resolve registry service: %v", err)
	}

	queueService, err := resolveQueueService(container, cfg)
	if err != nil {
		log.Fatalf("Failed to resolve queue service: %v", err)
	}

	// Start executor
	ctx := context.Background()
	if err := executorService.Start(ctx); err != nil {
		log.Printf("Failed to start executor: %v", err)
	}
	defer executorService.Stop(ctx)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API routes
	api := router.Group("/api/v1")
	{
		// Workflow execution
		api.POST("/execute", executeWorkflow(executorService, registryService))
		api.POST("/execute/async", asyncExecuteWorkflow(executorService, registryService, queueService))

		// Execution management
		api.GET("/executions", listExecutions(executorService))
		api.GET("/executions/:runId", getExecutionStatus(executorService))
		api.GET("/executions/:runId/details", getExecutionDetails(executorService))
		api.GET("/executions/:runId/data", getExecutionData(executorService))
		api.GET("/executions/:runId/metrics", getExecutionMetrics(executorService))
		api.POST("/executions/:runId/cancel", cancelExecution(executorService))
		api.POST("/executions/:runId/pause", pauseExecution(executorService))
		api.POST("/executions/:runId/resume", resumeExecution(executorService))
		api.POST("/executions/:runId/retry", retryExecution(executorService))

		// Step execution details
		api.GET("/executions/:runId/steps/:stepIndex", getStepExecution(executorService))
		api.GET("/executions/:runId/steps/:stepIndex/child-steps/:childStepIndex", getChildStepExecution(executorService))
	}

	// Health check with DI container health
	router.GET("/health", func(c *gin.Context) {
		health := container.HealthCheck()
		healthyComponents := 0
		totalComponents := len(health)

		for _, status := range health {
			if status.Healthy {
				healthyComponents++
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "executor-di",
			"timestamp": time.Now().Unix(),
			"location":  "bank-dc",
			"di_framework": gin.H{
				"container_type": "high-performance",
				"components": gin.H{
					"total":   totalComponents,
					"healthy": healthyComponents,
				},
			},
			"metrics": gin.H{
				"worker_count":             cfg.Executor.WorkerCount,
				"max_concurrent_workflows": cfg.Executor.MaxConcurrentWorkflows,
			},
		})
	})

	// DI container health endpoint
	router.GET("/health/di", func(c *gin.Context) {
		health := container.HealthCheck()
		c.JSON(http.StatusOK, gin.H{
			"container": "high-performance-di",
			"health":    health,
			"timestamp": time.Now().Unix(),
		})
	})

	// Start server
	port := getEnv("EXECUTOR_PORT", "8081")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Executor Service (DI) started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Executor Service (DI)...")

	// Give server time to finish existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Executor Service (DI) exited gracefully")
}

// initializeContainer creates and configures the DI container
func initializeContainer(cfg *config.Config) (di.Container, error) {
	// Create high-performance container for 5000+ TPS
	diConfig := di.DefaultConfig()
	diConfig.PoolSize = 1000
	diConfig.EnableMetrics = true
	container := di.NewWithConfig(diConfig)

	// Initialize global primitives
	primitiveConfig := &primitive.Config{
		EchoEnabled: true,
	}

	if err := primitive.Init(primitiveConfig); err != nil {
		log.Printf("Warning: Failed to initialize primitives: %v", err)
	} else {
		log.Printf("Initialized global primitives successfully")
	}

	// Register primitive services with DI
	if err := di.RegisterPrimitiveServices(container, primitiveConfig); err != nil {
		return nil, fmt.Errorf("failed to register primitive services: %w", err)
	}

	// Register core services
	if err := registerCoreServices(container, cfg); err != nil {
		return nil, fmt.Errorf("failed to register core services: %w", err)
	}

	// Register executor dependencies
	if err := executor.RegisterExecutorDependencies(container); err != nil {
		return nil, fmt.Errorf("failed to register executor dependencies: %w", err)
	}

	// Start the container
	if err := container.Start(); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Health check
	health := container.HealthCheck()
	log.Printf("Container health: %d components registered", len(health))
	for name, status := range health {
		if !status.Healthy {
			log.Printf("Warning: Component %s is unhealthy: %s", name, status.Error)
		}
	}

	return container, nil
}

// registerCoreServices registers core application services
func registerCoreServices(container di.Container, cfg *config.Config) error {
	// Register registry service - use HTTP registry to connect to the actual registry service
	err := container.RegisterFactory((*registry.Registry)(nil), func(c di.Container) (interface{}, error) {
		// Get registry service URL from environment variable
		registryURL := getEnv("REGISTRY_SERVICE_URL", "http://registry-service:8080")
		log.Printf("Creating HTTP registry client for endpoint: %s", registryURL)
		return registry.NewHTTPRegistry(registryURL)
	}, di.Singleton)
	if err != nil {
		return err
	}

	// Register state management
	err = container.RegisterFactory((*state.StateManagement)(nil), func(c di.Container) (interface{}, error) {
		return state.NewInMemoryState(), nil
	}, di.Singleton)
	if err != nil {
		return err
	}

	// Register queue service factory
	err = container.RegisterFactory((*queue.Queue)(nil), func(c di.Container) (interface{}, error) {
		return createQueueService(cfg), nil
	}, di.Singleton)
	if err != nil {
		return err
	}

	// Register configuration
	err = container.RegisterInstance((*config.Config)(nil), cfg)
	if err != nil {
		return err
	}

	return nil
}

// createQueueService creates the appropriate queue service based on config
func createQueueService(cfg *config.Config) queue.Queue {
	if cfg.Queue.Type == "nats" {
		enhancedConfig := queue.EnhancedNATSConfig{
			URLs:           cfg.Queue.NATS.URLs,
			StreamName:     cfg.Queue.NATS.StreamName,
			SubjectPrefix:  cfg.Queue.NATS.SubjectPrefix,
			DurableName:    cfg.Queue.NATS.DurableName,
			MaxReconnects:  cfg.Queue.NATS.MaxReconnects,
			ReconnectWait:  cfg.Queue.NATS.ReconnectWait,
			ConnectTimeout: cfg.Queue.NATS.ConnectTimeout,
		}
		enhancedQueue, err := queue.NewEnhancedNATSQueue(enhancedConfig)
		if err != nil {
			log.Printf("Failed to create enhanced NATS queue, falling back to in-memory: %v", err)
			return queue.NewInMemoryQueue()
		}
		return enhancedQueue
	}
	return queue.NewInMemoryQueue()
}

// resolveExecutorService resolves the executor service from container
func resolveExecutorService(container di.Container) (executor.Executor, error) {
	// Resolve executor factory
	factoryInstance, err := container.Resolve((*executor.DIFactory)(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve executor factory: %w", err)
	}

	factory := factoryInstance.(*executor.DIFactory)

	// Resolve executor config
	configInstance, err := container.Resolve((*executor.Config)(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve executor config: %w", err)
	}

	execConfig := configInstance.(executor.Config)

	// Create executor with DI
	exec, err := factory.CreateWorkflowExecutor(execConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow executor: %w", err)
	}

	return exec, nil
}

// resolveRegistryService resolves the registry service from container
func resolveRegistryService(container di.Container) (registry.Registry, error) {
	instance, err := container.Resolve((*registry.Registry)(nil))
	if err != nil {
		log.Printf("Failed to resolve registry from DI container: %v", err)
		log.Printf("Falling back to in-memory registry")
		// Fall back to creating directly
		return registry.NewInMemoryRegistry(), nil
	}
	log.Printf("Successfully resolved HTTP registry from DI container")
	return instance.(registry.Registry), nil
}

// resolveQueueService resolves the queue service from container
func resolveQueueService(container di.Container, cfg *config.Config) (queue.Queue, error) {
	instance, err := container.Resolve((*queue.Queue)(nil))
	if err != nil {
		// Fall back to creating directly
		return createQueueService(cfg), nil
	}
	return instance.(queue.Queue), nil
}

// Handler functions

func executeWorkflow(exec executor.Executor, reg registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var request struct {
			WorkflowID string                 `json:"workflow_id" binding:"required"`
			InputData  map[string]interface{} `json:"input_data,omitempty"`
			TimeoutMs  int64                  `json:"timeout_ms,omitempty"`
			Priority   int                    `json:"priority,omitempty"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		log.Printf("[Executor] Getting workflow with ID: %s", request.WorkflowID)
		// Get workflow
		workflow, err := reg.GetWorkflow(ctx, request.WorkflowID)
		if err != nil {
			log.Printf("[Executor] Failed to get workflow: %v", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Workflow not found",
				"details": err.Error(),
			})
			return
		}
		log.Printf("[Executor] Successfully retrieved workflow: %s", workflow.GetName())

		// Submit workflow for execution
		runID, err := exec.SubmitWorkflowByID(ctx, workflow.GetID())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to execute workflow",
				"details": err.Error(),
			})
			return
		}

		// Get execution status
		status, err := exec.GetExecutionStatus(ctx, runID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get execution status",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"run_id":      runID,
			"workflow_id": workflow.GetID(),
			"status":      status.Status,
			"progress":    status.Progress,
			"start_time":  status.StartTime,
			"end_time":    status.EndTime,
			"message":     "Workflow execution initiated",
		})
	}
}

func asyncExecuteWorkflow(exec executor.Executor, reg registry.Registry, q queue.Queue) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var request struct {
			WorkflowID        string                 `json:"workflow_id" binding:"required"`
			InputData         map[string]interface{} `json:"input_data,omitempty"`
			CallbackURL       string                 `json:"callback_url,omitempty"`
			TimeoutMs         int64                  `json:"timeout_ms,omitempty"`
			WaitForCompletion bool                   `json:"wait_for_completion,omitempty"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Get workflow
		workflow, err := reg.GetWorkflow(ctx, request.WorkflowID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Workflow not found",
				"details": err.Error(),
			})
			return
		}

		// Submit workflow for execution
		runID, err := exec.SubmitWorkflow(ctx, workflow)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to submit workflow",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"run_id":                  runID,
			"status":                  "queued",
			"message":                 "Workflow execution queued",
			"status_url":              fmt.Sprintf("/api/v1/executions/%s", runID),
			"result_url":              fmt.Sprintf("/api/v1/executions/%s/data", runID),
			"poll_after_ms":           1000,
			"estimated_completion_ms": 5000,
			"expires_at":              time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		})
	}
}

func listExecutions(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Parse query parameters
		workflowID := c.Query("workflow_id")
		status := c.Query("status")
		limitStr := c.Query("limit")
		offsetStr := c.Query("offset")

		limit := 50
		if limitStr != "" {
			if val, err := parseInt(limitStr); err == nil && val > 0 {
				limit = val
			}
		}

		offset := 0
		if offsetStr != "" {
			if val, err := parseInt(offsetStr); err == nil && val >= 0 {
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
		executions, err := exec.ListExecutions(ctx, filters)
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
}

func getExecutionStatus(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		runID := c.Param("runId")

		status, err := exec.GetExecutionStatus(ctx, runID)
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
}

func getExecutionDetails(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		runID := c.Param("runId")

		// TODO: Implement detailed execution information
		// For now, return basic status
		status, err := exec.GetExecutionStatus(ctx, runID)
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
			"message":                  "Detailed execution information (implementation pending)",
		})
	}
}

func getExecutionData(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		runID := c.Param("runId")

		data, err := exec.GetExecutionData(ctx, runID)
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
}

func getExecutionMetrics(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		runID := c.Param("runId")

		metrics, err := exec.GetMetrics(ctx, runID)
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
}

func cancelExecution(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		runID := c.Param("runId")

		err := exec.CancelExecution(ctx, runID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to cancel execution",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Execution cancelled successfully",
			"run_id":  runID,
		})
	}
}

func pauseExecution(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		runID := c.Param("runId")

		err := exec.PauseExecution(ctx, runID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to pause execution",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Execution paused successfully",
			"run_id":  runID,
		})
	}
}

func resumeExecution(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		runID := c.Param("runId")

		err := exec.ResumeExecution(ctx, runID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to resume execution",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Execution resumed successfully",
			"run_id":  runID,
		})
	}
}

func retryExecution(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		runID := c.Param("runId")

		err := exec.RetryExecution(ctx, runID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retry execution",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Execution retry initiated successfully",
			"run_id":  runID,
		})
	}
}

func getStepExecution(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		runID := c.Param("runId")
		stepIndexStr := c.Param("stepIndex")

		stepIndex, err := parseInt(stepIndexStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid step index",
				"details": err.Error(),
			})
			return
		}

		// TODO: Implement step execution details
		c.JSON(http.StatusOK, gin.H{
			"run_id":     runID,
			"step_index": stepIndex,
			"message":    "Step execution details (implementation pending)",
		})
	}
}

func getChildStepExecution(exec executor.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		runID := c.Param("runId")
		stepIndexStr := c.Param("stepIndex")
		childStepIndexStr := c.Param("childStepIndex")

		stepIndex, err := parseInt(stepIndexStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid step index",
				"details": err.Error(),
			})
			return
		}

		childStepIndex, err := parseInt(childStepIndexStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid child step index",
				"details": err.Error(),
			})
			return
		}

		// TODO: Implement child step execution details
		c.JSON(http.StatusOK, gin.H{
			"run_id":           runID,
			"step_index":       stepIndex,
			"child_step_index": childStepIndex,
			"message":          "Child step execution details (implementation pending)",
		})
	}
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
