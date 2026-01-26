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

	"unified-workflow/internal/common/model"
	"unified-workflow/internal/registry"
	"unified-workflow/workflows"
	"unified-workflow/workflows/steps"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("=== Registry Service ===")
	fmt.Println("Location: Bank Data Center")
	fmt.Println("Purpose: Workflow definition storage and management")
	fmt.Println("")

	// Initialize registry
	reg := registry.NewInMemoryRegistry()

	// Load example workflows on startup
	loadExampleWorkflows(reg)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API routes
	api := router.Group("/api/v1")
	{
		// Workflow management
		api.GET("/workflows", listWorkflows(reg))
		api.POST("/workflows", createWorkflow(reg))
		api.GET("/workflows/:id", getWorkflow(reg))
		api.PUT("/workflows/:id", updateWorkflow(reg))
		api.DELETE("/workflows/:id", deleteWorkflow(reg))
		api.GET("/workflows/:id/exists", containsWorkflow(reg))
		api.GET("/workflows/count", getWorkflowCount(reg))
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "registry",
			"timestamp": time.Now().Unix(),
			"location":  "bank-dc",
		})
	})

	// Start server
	port := getEnv("REGISTRY_PORT", "8080")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Registry Service started on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Registry Service...")

	// Give server time to finish existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Registry Service exited gracefully")
}

// Handler functions

func listWorkflows(reg registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get query parameters
		nameFilter := c.Query("name")
		descriptionFilter := c.Query("description")
		limit := c.Query("limit")
		offset := c.Query("offset")

		// Get all workflow IDs
		workflowIDs, err := reg.GetAllWorkflowIDs(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list workflows",
				"details": err.Error(),
			})
			return
		}

		// Get each workflow
		workflows := make([]gin.H, 0, len(workflowIDs))
		for _, workflowID := range workflowIDs {
			workflow, err := reg.GetWorkflow(ctx, workflowID)
			if err != nil {
				// Skip workflows that can't be retrieved
				continue
			}

			// Apply filters
			if nameFilter != "" && workflow.GetName() != nameFilter {
				continue
			}
			if descriptionFilter != "" && workflow.GetDescription() != descriptionFilter {
				continue
			}

			workflows = append(workflows, gin.H{
				"id":          workflow.GetID(),
				"name":        workflow.GetName(),
				"description": workflow.GetDescription(),
				"step_count":  workflow.GetStepCount(),
				"created_at":  time.Now().Format(time.RFC3339),
				"updated_at":  time.Now().Format(time.RFC3339),
			})
		}

		// Apply pagination
		totalCount := len(workflows)
		filteredCount := totalCount

		// TODO: Implement proper pagination
		_ = limit
		_ = offset

		c.JSON(http.StatusOK, gin.H{
			"workflows":      workflows,
			"total_count":    totalCount,
			"filtered_count": filteredCount,
		})
	}
}

func getWorkflow(reg registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		workflowID := c.Param("id")

		workflow, err := reg.GetWorkflow(ctx, workflowID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Workflow not found",
				"details": err.Error(),
			})
			return
		}

		// Get steps
		workflowSteps := workflow.GetSteps()
		stepDetails := make([]gin.H, 0, len(workflowSteps))
		for _, step := range workflowSteps {
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
			"created_at":  time.Now().Format(time.RFC3339),
			"updated_at":  time.Now().Format(time.RFC3339),
		})
	}
}

func createWorkflow(reg registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var request struct {
			Name        string                   `json:"name" binding:"required"`
			Description string                   `json:"description"`
			Steps       []map[string]interface{} `json:"steps"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Create workflow
		workflow := model.NewBaseWorkflow(request.Name, request.Description)

		// Add steps if provided
		for _, stepData := range request.Steps {
			if stepType, ok := stepData["type"].(string); ok {
				stepName, _ := stepData["name"].(string)
				if stepName == "" {
					stepName = stepType + "-step"
				}

				// Create appropriate step based on type
				var step model.Step
				switch stepType {
				case "sequential":
					step = model.NewSequentialStep(stepName)
				case "echo":
					// Create echo step with child steps
					step = steps.NewEchoStep(stepName, "Echo step created via API")
				default:
					// Default to sequential step
					step = model.NewSequentialStep(stepName)
				}

				// Add step to workflow
				workflow.AddStep(step)
			}
		}

		// Register workflow
		err := reg.RegisterWorkflow(ctx, workflow)
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
			"created_at":  time.Now().Format(time.RFC3339),
		})
	}
}

func updateWorkflow(reg registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		workflowID := c.Param("id")

		var request struct {
			Name        *string `json:"name,omitempty"`
			Description *string `json:"description,omitempty"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Get existing workflow
		workflow, err := reg.GetWorkflow(ctx, workflowID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Workflow not found",
				"details": err.Error(),
			})
			return
		}

		// TODO: Implement update logic
		// For now, just return the workflow
		_ = request

		c.JSON(http.StatusOK, gin.H{
			"id":          workflow.GetID(),
			"name":        workflow.GetName(),
			"description": workflow.GetDescription(),
			"message":     "Workflow update endpoint (implementation pending)",
			"updated_at":  time.Now().Format(time.RFC3339),
		})
	}
}

func deleteWorkflow(reg registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		workflowID := c.Param("id")

		err := reg.RemoveWorkflow(ctx, workflowID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Failed to delete workflow",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Workflow deleted successfully",
			"deleted_id": workflowID,
			"deleted_at": time.Now().Format(time.RFC3339),
		})
	}
}

func containsWorkflow(reg registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		workflowID := c.Param("id")

		_, err := reg.GetWorkflow(ctx, workflowID)
		exists := err == nil

		c.JSON(http.StatusOK, gin.H{
			"exists": exists,
			"id":     workflowID,
		})
	}
}

func getWorkflowCount(reg registry.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		workflowIDs, err := reg.GetAllWorkflowIDs(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get workflow count",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count": len(workflowIDs),
		})
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// loadExampleWorkflows loads example workflows from the workflows package
func loadExampleWorkflows(reg registry.Registry) {
	ctx := context.Background()

	// Get example workflows
	exampleWorkflows := workflows.GetExampleWorkflows()

	// Register each workflow
	for _, workflow := range exampleWorkflows {
		err := reg.RegisterWorkflow(ctx, workflow)
		if err != nil {
			log.Printf("Failed to register example workflow '%s': %v", workflow.GetName(), err)
		} else {
			log.Printf("Registered example workflow: %s", workflow.GetName())
		}
	}

	log.Printf("Loaded %d example workflows", len(exampleWorkflows))
}
