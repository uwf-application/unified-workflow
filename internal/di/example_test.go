package di_test

import (
	"testing"
	"time"

	"unified-workflow/internal/di"
)

// TestService is a simple service for testing
type TestService interface {
	DoWork() string
}

// testServiceImpl implements TestService
type testServiceImpl struct {
	id string
}

func (s *testServiceImpl) DoWork() string {
	return "Work done by " + s.id
}

// HealthChecker is a health check interface
type HealthChecker interface {
	HealthCheck() error
}

// healthyService implements both TestService and HealthChecker
type healthyService struct {
	id string
}

func (s *healthyService) DoWork() string {
	return "Healthy work by " + s.id
}

func (s *healthyService) HealthCheck() error {
	return nil // Always healthy
}

func TestBasicDI(t *testing.T) {
	container := di.New()

	// Register a singleton
	err := container.RegisterSingleton((*TestService)(nil), &testServiceImpl{id: "singleton"})
	if err != nil {
		t.Fatalf("Failed to register singleton: %v", err)
	}

	// Start the container
	err = container.Start()
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Resolve the singleton
	instance, err := container.Resolve((*TestService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve singleton: %v", err)
	}

	service := instance.(TestService)
	if service.DoWork() != "Work done by singleton" {
		t.Errorf("Unexpected result: %s", service.DoWork())
	}

	// Resolve again should return same instance
	instance2, err := container.Resolve((*TestService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve singleton second time: %v", err)
	}

	if instance != instance2 {
		t.Error("Singleton should return same instance")
	}

	// Stop the container
	err = container.Stop()
	if err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}
}

func TestTransientLifecycle(t *testing.T) {
	container := di.New()

	// Register a transient service
	err := container.RegisterFactory((*TestService)(nil), func(c di.Container) (interface{}, error) {
		return &testServiceImpl{id: time.Now().String()}, nil
	}, di.Transient)

	if err != nil {
		t.Fatalf("Failed to register transient: %v", err)
	}

	// Start the container
	err = container.Start()
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Resolve multiple times should return different instances
	instance1, err := container.Resolve((*TestService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve first instance: %v", err)
	}

	instance2, err := container.Resolve((*TestService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve second instance: %v", err)
	}

	if instance1 == instance2 {
		t.Error("Transient should return different instances")
	}

	// Stop the container
	err = container.Stop()
	if err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}
}

func TestScopedLifecycle(t *testing.T) {
	container := di.New()

	// Register a scoped service
	counter := 0
	err := container.RegisterFactory((*TestService)(nil), func(c di.Container) (interface{}, error) {
		counter++
		return &testServiceImpl{id: "scoped-" + string(rune(counter))}, nil
	}, di.Scoped)

	if err != nil {
		t.Fatalf("Failed to register scoped: %v", err)
	}

	// Start the container
	err = container.Start()
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Create a scope
	scope := container.BeginScope()
	defer scope.Dispose()

	// Resolve within scope - should create new instance
	instance1, err := scope.Resolve((*TestService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve in scope: %v", err)
	}

	// Resolve again in same scope - should return same instance
	instance2, err := scope.Resolve((*TestService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve second time in scope: %v", err)
	}

	if instance1 != instance2 {
		t.Error("Scoped should return same instance within same scope")
	}

	// Create another scope
	scope2 := container.BeginScope()
	defer scope2.Dispose()

	// Resolve in different scope - should create new instance
	instance3, err := scope2.Resolve((*TestService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve in second scope: %v", err)
	}

	if instance1 == instance3 {
		t.Error("Scoped should return different instances in different scopes")
	}

	// Stop the container
	err = container.Stop()
	if err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}
}

func TestPooledLifecycle(t *testing.T) {
	container := di.NewWithConfig(di.DefaultConfig())

	// Register a pooled service
	counter := 0
	err := container.RegisterFactory((*TestService)(nil), func(c di.Container) (interface{}, error) {
		counter++
		return &testServiceImpl{id: "pooled-" + string(rune(counter))}, nil
	}, di.Pooled)

	if err != nil {
		t.Fatalf("Failed to register pooled: %v", err)
	}

	// Start the container
	err = container.Start()
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Resolve pooled instance
	instance, err := container.Resolve((*TestService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve pooled: %v", err)
	}

	// Should be a pooledInstance wrapper
	pooled, ok := instance.(interface{ Release() })
	if !ok {
		t.Fatal("Pooled instance should have Release method")
	}

	// Use the service
	service := instance.(TestService)
	result := service.DoWork()
	if result == "" {
		t.Error("Pooled service should work")
	}

	// Release back to pool
	pooled.Release()

	// Health check
	health := container.HealthCheck()
	if len(health) == 0 {
		t.Error("Health check should return status")
	}

	// Stop the container
	err = container.Stop()
	if err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	container := di.New()

	// Register a healthy service
	err := container.RegisterSingleton((*TestService)(nil), &healthyService{id: "healthy"})
	if err != nil {
		t.Fatalf("Failed to register healthy service: %v", err)
	}

	// Start the container
	err = container.Start()
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Check health
	health := container.HealthCheck()
	if len(health) == 0 {
		t.Error("Health check should return status")
	}

	// Stop the container
	err = container.Stop()
	if err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}
}

func BenchmarkSingletonResolution(b *testing.B) {
	container := di.New()

	err := container.RegisterSingleton((*TestService)(nil), &testServiceImpl{id: "benchmark"})
	if err != nil {
		b.Fatalf("Failed to register: %v", err)
	}

	err = container.Start()
	if err != nil {
		b.Fatalf("Failed to start: %v", err)
	}
	defer container.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		instance, err := container.Resolve((*TestService)(nil))
		if err != nil {
			b.Fatalf("Failed to resolve: %v", err)
		}
		_ = instance.(TestService).DoWork()
	}
}

func BenchmarkTransientResolution(b *testing.B) {
	container := di.New()

	err := container.RegisterFactory((*TestService)(nil), func(c di.Container) (interface{}, error) {
		return &testServiceImpl{id: "transient"}, nil
	}, di.Transient)

	if err != nil {
		b.Fatalf("Failed to register: %v", err)
	}

	err = container.Start()
	if err != nil {
		b.Fatalf("Failed to start: %v", err)
	}
	defer container.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		instance, err := container.Resolve((*TestService)(nil))
		if err != nil {
			b.Fatalf("Failed to resolve: %v", err)
		}
		_ = instance.(TestService).DoWork()
	}
}
