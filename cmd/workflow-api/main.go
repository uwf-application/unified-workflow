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

	"unified-workflow/internal/api/handlers"
	"unified-workflow/internal/config"
	"unified-workflow/internal/executor"
	"unified-workflow/internal/queue"
	"unified-workflow/internal/registry"
	"unified-workflow/internal/state"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize components
	reg := registry.NewInMemoryRegistry()
	stateMgmt := state.NewInMemoryState()

	// Initialize queue based on configuration
	var q queue.Queue
	if cfg.Queue.Type == "nats" {
		natsConfig := queue.NATSConfig{
			URLs:           cfg.Queue.NATS.URLs,
			StreamName:     cfg.Queue.NATS.StreamName,
			SubjectPrefix:  cfg.Queue.NATS.SubjectPrefix,
			DurableName:    cfg.Queue.NATS.DurableName,
			MaxReconnects:  cfg.Queue.NATS.MaxReconnects,
			ReconnectWait:  cfg.Queue.NATS.ReconnectWait,
			ConnectTimeout: cfg.Queue.NATS.ConnectTimeout,
		}
		q, err = queue.NewNATSQueue(natsConfig)
		if err != nil {
			log.Printf("Failed to create NATS queue, falling back to in-memory: %v", err)
			q = queue.NewInMemoryQueue()
		}
	} else {
		q = queue.NewInMemoryQueue()
	}

	// Initialize executor
	exec := executor.NewSimpleExecutor(reg, q, stateMgmt)

	// Start executor
	ctx := context.Background()
	if err := exec.Start(ctx); err != nil {
		log.Printf("Failed to start executor: %v", err)
	}
	defer exec.Stop(ctx)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize handlers
	handler := handlers.NewWorkflowHandler(exec, reg, stateMgmt)

	// API routes
	api := router.Group("/api/v1")
	{
		// Workflow management
		api.GET("/workflows", handler.ListWorkflows)
		api.POST("/workflows", handler.CreateWorkflow)
		api.GET("/workflows/:id", handler.GetWorkflow)
		api.DELETE("/workflows/:id", handler.DeleteWorkflow)

		// Workflow execution
		api.POST("/workflows/:id/execute", handler.ExecuteWorkflow)
		api.POST("/workflows/:id/async-execute", handler.AsyncExecuteWorkflow) // NEW: Async execution
		api.GET("/executions", handler.ListExecutions)
		api.GET("/executions/:runId", handler.GetExecutionStatus)
		api.GET("/executions/:runId/result", handler.GetExecutionResult) // NEW: Get async result
		api.POST("/executions/:runId/cancel", handler.CancelExecution)
		api.POST("/executions/:runId/pause", handler.PauseExecution)
		api.POST("/executions/:runId/resume", handler.ResumeExecution)
		api.POST("/executions/:runId/retry", handler.RetryExecution)
		api.GET("/executions/:runId/data", handler.GetExecutionData)
		api.GET("/executions/:runId/metrics", handler.GetExecutionMetrics)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %d", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give server time to finish existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
