package executor

import (
	"unified-workflow/internal/di"
	"unified-workflow/internal/primitive"
	"unified-workflow/internal/registry"
	"unified-workflow/internal/state"
)

// DIFactory creates executor instances using dependency injection
type DIFactory struct {
	container di.Container
}

// NewDIFactory creates a new DI factory
func NewDIFactory(container di.Container) *DIFactory {
	return &DIFactory{
		container: container,
	}
}

// CreateWorkflowExecutor creates a workflow executor with dependencies from DI container
func (f *DIFactory) CreateWorkflowExecutor(config Config) (*WorkflowExecutor, error) {
	// Resolve dependencies from container
	reg, err := f.resolveRegistry()
	if err != nil {
		return nil, err
	}

	stateMgmt, err := f.resolveStateManagement()
	if err != nil {
		return nil, err
	}

	primitiveProvider, err := f.resolvePrimitiveProvider()
	if err != nil {
		return nil, err
	}

	// Create executor with resolved dependencies
	executor := NewWorkflowExecutor(reg, stateMgmt, config)

	// Set primitive provider if available
	if primitiveProvider != nil {
		// We need to modify the executor to use the primitive provider
		// For now, we'll create a wrapper or modify the executor
		// This is a placeholder for the actual integration
		_ = primitiveProvider
	}

	return executor, nil
}

// resolveRegistry resolves the registry from DI container
func (f *DIFactory) resolveRegistry() (registry.Registry, error) {
	instance, err := f.container.Resolve((*registry.Registry)(nil))
	if err != nil {
		// Fall back to default in-memory registry
		return registry.NewInMemoryRegistry(), nil
	}
	return instance.(registry.Registry), nil
}

// resolveStateManagement resolves state management from DI container
func (f *DIFactory) resolveStateManagement() (state.StateManagement, error) {
	instance, err := f.container.Resolve((*state.StateManagement)(nil))
	if err != nil {
		// Fall back to default in-memory state
		return state.NewInMemoryState(), nil
	}
	return instance.(state.StateManagement), nil
}

// resolvePrimitiveProvider resolves primitive provider from DI container
func (f *DIFactory) resolvePrimitiveProvider() (*di.PrimitiveProvider, error) {
	instance, err := f.container.Resolve((*di.PrimitiveProvider)(nil))
	if err != nil {
		// Create a new primitive provider using the container
		provider := di.NewPrimitiveProvider(f.container)
		// Register it for future use
		f.container.RegisterInstance((*di.PrimitiveProvider)(nil), provider)
		return provider, nil
	}
	return instance.(*di.PrimitiveProvider), nil
}

// RegisterExecutorDependencies registers all executor dependencies with the container
func RegisterExecutorDependencies(container di.Container) error {
	// Note: Registry is registered in registerCoreServices in main.go
	// We don't register it here to avoid overwriting the HTTP registry

	// Register state management
	err := container.RegisterFactory((*state.StateManagement)(nil), func(c di.Container) (interface{}, error) {
		return state.NewInMemoryState(), nil
	}, di.Singleton)
	if err != nil {
		return err
	}

	// Register primitive provider
	err = container.RegisterFactory((*di.PrimitiveProvider)(nil), func(c di.Container) (interface{}, error) {
		return di.NewPrimitiveProvider(c), nil
	}, di.Singleton)
	if err != nil {
		return err
	}

	// Register executor config
	err = container.RegisterFactory((*Config)(nil), func(c di.Container) (interface{}, error) {
		return DefaultConfig(), nil
	}, di.Singleton)
	if err != nil {
		return err
	}

	// Register workflow executor factory
	err = container.RegisterFactory((*DIFactory)(nil), func(c di.Container) (interface{}, error) {
		return NewDIFactory(c), nil
	}, di.Singleton)
	if err != nil {
		return err
	}

	return nil
}

// DefaultConfig returns the default executor configuration
func DefaultConfig() Config {
	return Config{
		WorkerCount:            5,
		QueuePollInterval:      100 * 1000000, // 100ms in nanoseconds
		MaxRetries:             3,
		RetryDelay:             1 * 1000000000,      // 1s in nanoseconds
		ExecutionTimeout:       5 * 60 * 1000000000, // 5m in nanoseconds
		StepTimeout:            30 * 1000000000,     // 30s in nanoseconds
		EnableMetrics:          true,
		EnableTracing:          false,
		MaxConcurrentWorkflows: 10,
	}
}

// InitializeContainer initializes a DI container with all executor dependencies
func InitializeContainer() (di.Container, error) {
	container := di.New()

	// Register primitive services
	primitiveConfig := &primitive.Config{
		EchoEnabled: true,
	}

	// Initialize global primitives
	if err := primitive.Init(primitiveConfig); err != nil {
		return nil, err
	}

	// Register primitive services with DI
	if err := di.RegisterPrimitiveServices(container, primitiveConfig); err != nil {
		return nil, err
	}

	// Register executor dependencies
	if err := RegisterExecutorDependencies(container); err != nil {
		return nil, err
	}

	// Start the container
	if err := container.Start(); err != nil {
		return nil, err
	}

	return container, nil
}
