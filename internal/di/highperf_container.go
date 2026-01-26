package di

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// highPerfContainer is a high-performance dependency injection container
type highPerfContainer struct {
	mu           sync.RWMutex
	singletons   *sync.Map // typeKey -> instance
	factories    *sync.Map // typeKey -> cachedFactory
	scopedCache  *sync.Map // scopeID -> *sync.Map(typeKey -> instance)
	registry     map[typeKey]*registration
	pools        map[reflect.Type]*sync.Pool
	scopes       map[string]*scopeImpl
	scopeCounter *atomic.Int64
	config       Config
	started      bool
	stopped      bool
}

// registration holds registration information
type registration struct {
	provider    Provider
	lifecycle   Lifecycle
	pool        *sync.Pool
	poolSize    int
	instance    interface{} // For singletons
	initialized bool
}

// cachedFactory holds a factory function with caching
type cachedFactory struct {
	factory   FactoryFunc
	lifecycle Lifecycle
	created   time.Time
	lastUsed  time.Time
	useCount  int64
}

// scopeImpl implements the Scope interface
type scopeImpl struct {
	id        string
	container *highPerfContainer
	instances *sync.Map // typeKey -> instance
	disposed  bool
	created   time.Time
}

// ID returns the scope ID
func (s *scopeImpl) ID() string {
	return s.id
}

// Resolve resolves a dependency within this scope
func (s *scopeImpl) Resolve(key interface{}) (interface{}, error) {
	if s.disposed {
		return nil, ErrScopeDisposed
	}
	return s.container.ResolveScoped(key, s)
}

// Dispose disposes the scope and all scoped instances
func (s *scopeImpl) Dispose() {
	if s.disposed {
		return
	}
	s.disposed = true

	// Clean up scoped instances
	s.instances.Range(func(key, value interface{}) bool {
		// Call cleanup if the value has a Dispose method
		if disposer, ok := value.(interface{ Dispose() }); ok {
			disposer.Dispose()
		}
		s.instances.Delete(key)
		return true
	})

	// Remove from container
	s.container.mu.Lock()
	delete(s.container.scopes, s.id)
	s.container.mu.Unlock()
}

// initializePools pre-allocates object pools
func (c *highPerfContainer) initializePools() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, reg := range c.registry {
		if reg.lifecycle == Pooled {
			pool := &sync.Pool{
				New: func() interface{} {
					instance, err := reg.provider.Create(c)
					if err != nil {
						return nil
					}
					return instance
				},
			}

			// Pre-allocate pool items
			for i := 0; i < c.config.PoolSize; i++ {
				instance := pool.New()
				if instance != nil {
					pool.Put(instance)
				}
			}

			reg.pool = pool
			c.pools[key.Type] = pool
		}
	}
}

// Register registers a dependency with the container
func (c *highPerfContainer) Register(key interface{}, provider Provider, lifecycle Lifecycle) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return errors.New("cannot register after container has started")
	}

	keyType, err := getType(key)
	if err != nil {
		return err
	}

	reg := &registration{
		provider:  provider,
		lifecycle: lifecycle,
		poolSize:  c.config.PoolSize,
	}

	c.registry[typeKey{keyType}] = reg

	// If it's a singleton and we have an instance, store it
	if lifecycle == Singleton {
		instance, err := provider.Create(c)
		if err != nil {
			return fmt.Errorf("failed to create singleton instance: %w", err)
		}
		reg.instance = instance
		reg.initialized = true
		c.singletons.Store(typeKey{keyType}, instance)
	}

	return nil
}

// RegisterSingleton registers a singleton instance
func (c *highPerfContainer) RegisterSingleton(key, instance interface{}) error {
	_, err := getType(key)
	if err != nil {
		return err
	}

	provider := ProviderFunc(func(Container) (interface{}, error) {
		return instance, nil
	})

	return c.Register(key, provider, Singleton)
}

// RegisterFactory registers a factory function
func (c *highPerfContainer) RegisterFactory(key interface{}, factory FactoryFunc, lifecycle Lifecycle) error {
	provider := ProviderFunc(func(c Container) (interface{}, error) {
		return factory(c)
	})

	return c.Register(key, provider, lifecycle)
}

// RegisterInstance registers an existing instance
func (c *highPerfContainer) RegisterInstance(key, instance interface{}) error {
	return c.RegisterSingleton(key, instance)
}

// Resolve resolves a dependency
func (c *highPerfContainer) Resolve(key interface{}) (interface{}, error) {
	return c.resolve(key, nil)
}

// MustResolve resolves a dependency or panics
func (c *highPerfContainer) MustResolve(key interface{}) interface{} {
	instance, err := c.Resolve(key)
	if err != nil {
		panic(fmt.Sprintf("failed to resolve dependency: %v", err))
	}
	return instance
}

// resolve is the internal resolution method
func (c *highPerfContainer) resolve(key interface{}, scope *scopeImpl) (interface{}, error) {
	keyType, err := getType(key)
	if err != nil {
		return nil, err
	}

	typeKey := typeKey{keyType}

	// Fast path: check singletons first (lock-free)
	if instance, ok := c.singletons.Load(typeKey); ok {
		return instance, nil
	}

	// Check if we have a registration
	c.mu.RLock()
	reg, registered := c.registry[typeKey]
	c.mu.RUnlock()

	if !registered {
		return nil, ErrNotRegistered
	}

	// Handle based on lifecycle
	switch reg.lifecycle {
	case Singleton:
		// Double-check locking for singleton initialization
		c.mu.Lock()
		if !reg.initialized {
			instance, err := reg.provider.Create(c)
			if err != nil {
				c.mu.Unlock()
				return nil, err
			}
			reg.instance = instance
			reg.initialized = true
			c.singletons.Store(typeKey, instance)
		}
		instance := reg.instance
		c.mu.Unlock()
		return instance, nil

	case Transient:
		return reg.provider.Create(c)

	case Scoped:
		if scope == nil {
			// Create a temporary scope for this resolution
			tempScope := c.BeginScope()
			defer tempScope.Dispose()
			return c.ResolveScoped(key, tempScope)
		}

		// Check if already resolved in this scope
		if instance, ok := scope.instances.Load(typeKey); ok {
			return instance, nil
		}

		// Create new instance for this scope
		instance, err := reg.provider.Create(c)
		if err != nil {
			return nil, err
		}

		scope.instances.Store(typeKey, instance)
		return instance, nil

	case Pooled:
		if reg.pool == nil {
			// Initialize pool on first use
			c.mu.Lock()
			if reg.pool == nil {
				reg.pool = &sync.Pool{
					New: func() interface{} {
						instance, err := reg.provider.Create(c)
						if err != nil {
							return nil
						}
						return instance
					},
				}
				c.pools[keyType] = reg.pool
			}
			c.mu.Unlock()
		}

		instance := reg.pool.Get()
		if instance == nil {
			return nil, ErrPoolExhausted
		}

		// Return to pool when done (caller should use defer)
		return &pooledInstance{
			instance: instance,
			pool:     reg.pool,
		}, nil

	default:
		return nil, fmt.Errorf("unknown lifecycle: %v", reg.lifecycle)
	}
}

// ResolveScoped resolves a dependency within a scope
func (c *highPerfContainer) ResolveScoped(key interface{}, scope Scope) (interface{}, error) {
	scopeImpl, ok := scope.(*scopeImpl)
	if !ok {
		return nil, errors.New("invalid scope type")
	}

	return c.resolve(key, scopeImpl)
}

// BeginScope creates a new scope
func (c *highPerfContainer) BeginScope() Scope {
	scopeID := fmt.Sprintf("scope-%d-%d", time.Now().UnixNano(), c.scopeCounter.Add(1))

	scope := &scopeImpl{
		id:        scopeID,
		container: c,
		instances: &sync.Map{},
		created:   time.Now(),
	}

	c.mu.Lock()
	c.scopes[scopeID] = scope
	c.mu.Unlock()

	return scope
}

// Has checks if a dependency is registered
func (c *highPerfContainer) Has(key interface{}) bool {
	keyType, err := getType(key)
	if err != nil {
		return false
	}

	c.mu.RLock()
	_, exists := c.registry[typeKey{keyType}]
	c.mu.RUnlock()

	return exists
}

// Clear clears all registrations
func (c *highPerfContainer) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.singletons = &sync.Map{}
	c.factories = &sync.Map{}
	c.scopedCache = &sync.Map{}
	c.registry = make(map[typeKey]*registration)
	c.pools = make(map[reflect.Type]*sync.Pool)

	// Dispose all scopes
	for _, scope := range c.scopes {
		scope.Dispose()
	}
	c.scopes = make(map[string]*scopeImpl)

	c.started = false
	c.stopped = false
}

// Start starts the container
func (c *highPerfContainer) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return errors.New("container already started")
	}

	c.started = true

	// Initialize all singletons
	for key, reg := range c.registry {
		if reg.lifecycle == Singleton && !reg.initialized {
			instance, err := reg.provider.Create(c)
			if err != nil {
				return fmt.Errorf("failed to initialize singleton %v: %w", key.Type, err)
			}
			reg.instance = instance
			reg.initialized = true
			c.singletons.Store(key, instance)
		}
	}

	return nil
}

// Stop stops the container
func (c *highPerfContainer) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopped {
		return errors.New("container already stopped")
	}

	c.stopped = true

	// Dispose all scopes
	for _, scope := range c.scopes {
		scope.Dispose()
	}

	// Clear all pools
	for _, pool := range c.pools {
		// Drain the pool
		for {
			item := pool.Get()
			if item == nil {
				break
			}
			// Call cleanup if the item has a Dispose method
			if disposer, ok := item.(interface{ Dispose() }); ok {
				disposer.Dispose()
			}
		}
	}

	return nil
}

// HealthCheck performs health checks on registered components
func (c *highPerfContainer) HealthCheck() map[string]HealthStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := make(map[string]HealthStatus)

	for key, reg := range c.registry {
		typeName := key.Type.String()

		// Check if component is healthy
		healthy := true
		errorMsg := ""

		if reg.lifecycle == Singleton && reg.initialized {
			// Check if singleton implements HealthChecker
			if checker, ok := reg.instance.(interface{ HealthCheck() error }); ok {
				if err := checker.HealthCheck(); err != nil {
					healthy = false
					errorMsg = err.Error()
				}
			}
		}

		status[typeName] = HealthStatus{
			Healthy: healthy,
			Status:  "ok",
			Error:   errorMsg,
		}
	}

	return status
}

// pooledInstance wraps an instance from a pool
type pooledInstance struct {
	instance interface{}
	pool     *sync.Pool
}

// Release returns the instance to the pool
func (p *pooledInstance) Release() {
	if p.pool != nil && p.instance != nil {
		p.pool.Put(p.instance)
		p.instance = nil
		p.pool = nil
	}
}

// getType extracts the reflect.Type from an interface
func getType(key interface{}) (reflect.Type, error) {
	if key == nil {
		return nil, ErrInvalidType
	}

	switch v := key.(type) {
	case reflect.Type:
		return v, nil
	default:
		typ := reflect.TypeOf(key)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		return typ, nil
	}
}
