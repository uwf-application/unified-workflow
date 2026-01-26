package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
)

func main() {
	log.Println("Starting workflow worker with DI/IoC framework...")

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
	queueService, err := resolveQueueService(container, cfg)
	if err != nil {
		log.Fatalf("Failed to resolve queue service: %v", err)
	}

	executorService, err := resolveExecutorService(container)
	if err != nil {
		log.Fatalf("Failed to resolve executor service: %v", err)
	}

	// Start executor
	ctx := context.Background()
	if err := executorService.Start(ctx); err != nil {
		log.Printf("Failed to start executor: %v", err)
	}
	defer executorService.Stop(ctx)

	log.Println("Workflow worker started with DI-enabled executor")

	// Process messages
	processMessages(ctx, queueService, executorService)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
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
	// Register registry service
	err := container.RegisterFactory((*registry.Registry)(nil), func(c di.Container) (interface{}, error) {
		return registry.NewInMemoryRegistry(), nil
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

// resolveQueueService resolves the queue service from container
func resolveQueueService(container di.Container, cfg *config.Config) (queue.Queue, error) {
	instance, err := container.Resolve((*queue.Queue)(nil))
	if err != nil {
		// Fall back to creating directly
		return createQueueService(cfg), nil
	}
	return instance.(queue.Queue), nil
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

// processMessages continuously processes messages from the queue
func processMessages(ctx context.Context, q queue.Queue, exec executor.Executor) {
	// Try to cast to EnhancedNATSQueue for enhanced features
	enhancedQueue, isEnhanced := q.(*queue.EnhancedNATSQueue)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if isEnhanced {
				// Use enhanced dequeue for better message handling
				processEnhancedMessage(ctx, enhancedQueue, exec)
			} else {
				// Use standard dequeue
				processStandardMessage(ctx, q, exec)
			}
		}
	}
}

// processEnhancedMessage processes a message using enhanced queue features
func processEnhancedMessage(ctx context.Context, q *queue.EnhancedNATSQueue, exec executor.Executor) {
	// Dequeue message with enhanced info
	enhancedMsg, err := q.DequeueEnhanced(ctx)
	if err != nil {
		log.Printf("Failed to dequeue message: %v", err)
		time.Sleep(1 * time.Second)
		return
	}

	if enhancedMsg == nil {
		// No messages available, wait before trying again
		time.Sleep(100 * time.Millisecond)
		return
	}

	// Process the message
	if err := processWorkflowExecution(ctx, enhancedMsg.Message, exec, q, enhancedMsg); err != nil {
		log.Printf("Failed to process workflow execution %s: %v", enhancedMsg.RunID, err)

		// Reject message for retry
		if err := q.RejectEnhanced(ctx, enhancedMsg, 5*time.Second); err != nil {
			log.Printf("Failed to reject message %s: %v", enhancedMsg.RunID, err)
		}
	} else {
		// Acknowledge successful processing
		if err := q.AcknowledgeEnhanced(ctx, enhancedMsg); err != nil {
			log.Printf("Failed to acknowledge message %s: %v", enhancedMsg.RunID, err)
		}
	}
}

// processStandardMessage processes a message using standard queue interface
func processStandardMessage(ctx context.Context, q queue.Queue, exec executor.Executor) {
	// Dequeue message
	msg, err := q.Dequeue(ctx)
	if err != nil {
		log.Printf("Failed to dequeue message: %v", err)
		time.Sleep(1 * time.Second)
		return
	}

	if msg == nil {
		// No messages available, wait before trying again
		time.Sleep(100 * time.Millisecond)
		return
	}

	// Process the message
	if err := processWorkflowExecution(ctx, msg, exec, q, nil); err != nil {
		log.Printf("Failed to process workflow execution %s: %v", msg.RunID, err)

		// Reject message for retry
		if err := q.Reject(ctx, msg.ID, 5*time.Second); err != nil {
			log.Printf("Failed to reject message %s: %v", msg.RunID, err)
		}
	} else {
		// Acknowledge successful processing
		if err := q.Acknowledge(ctx, msg.ID); err != nil {
			log.Printf("Failed to acknowledge message %s: %v", msg.RunID, err)
		}
	}
}

// processWorkflowExecution processes a workflow execution request
func processWorkflowExecution(ctx context.Context, msg *queue.Message, exec executor.Executor, q queue.Queue, enhancedMsg *queue.EnhancedMessage) error {
	// Unmarshal execution request
	var execReq queue.ExecutionRequest
	if err := json.Unmarshal(msg.Data, &execReq); err != nil {
		return fmt.Errorf("failed to unmarshal execution request: %w", err)
	}

	log.Printf("Processing workflow execution %s for workflow %s", execReq.RunID, execReq.WorkflowID)

	// Try to cast to WorkflowExecutor to use real execution
	if workflowExecutor, ok := exec.(*executor.WorkflowExecutor); ok {
		// Use real workflow execution with child-step tracking
		result, err := workflowExecutor.ExecuteWorkflow(ctx, execReq.WorkflowID, execReq.InputData)
		if err != nil {
			// Publish error if using enhanced queue
			if enhancedQueue, ok := q.(*queue.EnhancedNATSQueue); ok && enhancedMsg != nil {
				errorResult := queue.ExecutionResult{
					RunID:       execReq.RunID,
					WorkflowID:  execReq.WorkflowID,
					Status:      "failed",
					Error:       err.Error(),
					CompletedAt: time.Now(),
				}
				errorData, _ := json.Marshal(errorResult)
				if err := enhancedQueue.PublishError(ctx, execReq.RunID, errorData); err != nil {
					log.Printf("Failed to publish error for %s: %v", execReq.RunID, err)
				}
			}
			return fmt.Errorf("workflow execution failed: %w", err)
		}

		// Publish result if using enhanced queue
		if enhancedQueue, ok := q.(*queue.EnhancedNATSQueue); ok && enhancedMsg != nil {
			execResult := queue.ExecutionResult{
				RunID:       result.RunID,
				WorkflowID:  result.WorkflowID,
				Status:      result.Status,
				OutputData:  result.Result,
				CompletedAt: result.EndTime,
			}
			resultData, _ := json.Marshal(execResult)
			if err := enhancedQueue.PublishResult(ctx, execReq.RunID, resultData); err != nil {
				log.Printf("Failed to publish result for %s: %v", execReq.RunID, err)
			}
		}

		log.Printf("Successfully processed workflow execution %s with real executor, status: %s", result.RunID, result.Status)
		return nil
	}

	// Fall back to simulation if not a WorkflowExecutor
	log.Printf("Using simulated execution (executor is not WorkflowExecutor)")
	return simulateWorkflowExecution(ctx, execReq)
}

// simulateWorkflowExecution simulates workflow execution (fallback)
func simulateWorkflowExecution(ctx context.Context, execReq queue.ExecutionRequest) error {
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	log.Printf("Simulated workflow execution %s completed", execReq.RunID)
	return nil
}
