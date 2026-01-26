package di

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

// Lifecycle defines the lifecycle of a registered dependency
type Lifecycle int

const (
	// Transient creates a new instance each time it's resolved
	Transient Lifecycle = iota

	// Singleton creates a single instance shared across the container
	Singleton

	// Scoped creates an instance per scope (e.g., per request)
	Scoped

	// Pooled uses an object pool for instances
	Pooled
)

// Container is the dependency injection container interface
type Container interface {
	// Registration methods
	Register(interface{}, Provider, Lifecycle) error
	RegisterSingleton(interface{}, interface{}) error
	RegisterFactory(interface{}, FactoryFunc, Lifecycle) error
	RegisterInstance(interface{}, interface{}) error

	// Resolution methods
	Resolve(interface{}) (interface{}, error)
	MustResolve(interface{}) interface{}

	// Scoped resolution
	BeginScope() Scope
	ResolveScoped(interface{}, Scope) (interface{}, error)

	// Lifecycle management
	Start() error
	Stop() error
	HealthCheck() map[string]HealthStatus

	// Utility methods
	Has(interface{}) bool
	Clear()
}

// Scope represents a dependency scope (e.g., request scope)
type Scope interface {
	ID() string
	Resolve(interface{}) (interface{}, error)
	Dispose()
}

// Provider creates an instance of a dependency
type Provider interface {
	Create(Container) (interface{}, error)
}

// FactoryFunc is a function that creates an instance
type FactoryFunc func(Container) (interface{}, error)

// HealthStatus represents the health status of a component
type HealthStatus struct {
	Healthy bool   `json:"healthy"`
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
}

// ProviderFunc adapts a function to the Provider interface
type ProviderFunc func(Container) (interface{}, error)

func (f ProviderFunc) Create(c Container) (interface{}, error) {
	return f(c)
}

// Errors
var (
	ErrNotRegistered      = errors.New("dependency not registered")
	ErrInvalidType        = errors.New("invalid type")
	ErrCircularDependency = errors.New("circular dependency detected")
	ErrScopeDisposed      = errors.New("scope has been disposed")
	ErrPoolExhausted      = errors.New("object pool exhausted")
)

// TypeKey is used as a map key for type registration
type typeKey struct {
	reflect.Type
}

// New creates a new high-performance DI container
func New() Container {
	return &highPerfContainer{
		singletons:   &sync.Map{},
		factories:    &sync.Map{},
		scopedCache:  &sync.Map{},
		registry:     make(map[typeKey]*registration),
		pools:        make(map[reflect.Type]*sync.Pool),
		scopes:       make(map[string]*scopeImpl),
		scopeCounter: new(atomic.Int64),
	}
}

// NewWithConfig creates a new container with configuration
func NewWithConfig(config Config) Container {
	container := &highPerfContainer{
		singletons:   &sync.Map{},
		factories:    &sync.Map{},
		scopedCache:  &sync.Map{},
		registry:     make(map[typeKey]*registration),
		pools:        make(map[reflect.Type]*sync.Pool),
		scopes:       make(map[string]*scopeImpl),
		scopeCounter: new(atomic.Int64),
		config:       config,
	}

	// Pre-allocate pools if configured
	if config.PreAllocatePools {
		container.initializePools()
	}

	return container
}

// Config holds container configuration
type Config struct {
	// PreAllocatePools pre-allocates object pools on startup
	PreAllocatePools bool `json:"pre_allocate_pools" yaml:"pre_allocate_pools"`

	// PoolSize is the size of object pools
	PoolSize int `json:"pool_size" yaml:"pool_size"`

	// EnableMetrics enables performance metrics collection
	EnableMetrics bool `json:"enable_metrics" yaml:"enable_metrics"`

	// MaxConcurrentResolutions limits concurrent resolutions
	MaxConcurrentResolutions int `json:"max_concurrent_resolutions" yaml:"max_concurrent_resolutions"`
}

// DefaultConfig returns the default container configuration
func DefaultConfig() Config {
	return Config{
		PreAllocatePools:         true,
		PoolSize:                 100,
		EnableMetrics:            true,
		MaxConcurrentResolutions: 10000,
	}
}
