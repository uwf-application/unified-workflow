package primitive

import (
	"fmt"
	"sync"
	"unified-workflow/internal/primitive/services/antifraud/models"
	serviceclientsantifraud "unified-workflow/internal/serviceclients/antifraud"
)

// Primitive is the global struct that provides access to all services
// This enables the syntax: primitive.Default.Storage.Save(data)
type Primitive struct {
	Storage   StorageService
	Echo      EchoService
	HTTP      HTTPService
	DB        DatabaseService
	Antifraud AntifraudService
}

var (
	// Default is the global primitive instance
	Default  *Primitive
	initOnce sync.Once
)

// Init initializes the global primitive instance with service implementations
// This should be called once at application startup
func Init(config *Config) error {
	var initErr error
	initOnce.Do(func() {
		if config == nil {
			config = &Config{}
		}

		// Initialize services
		storageService, err := initStorageService(config)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize storage service: %w", err)
			return
		}

		echoService, err := initEchoService(config)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize echo service: %w", err)
			return
		}

		// Initialize antifraud service
		antifraudService, err := initAntifraudService(config)
		if err != nil {
			initErr = fmt.Errorf("failed to initialize antifraud service: %w", err)
			return
		}

		// Create the global instance
		Default = &Primitive{
			Storage:   storageService,
			Echo:      echoService,
			Antifraud: antifraudService,
			// HTTP and DB services can be added later
		}
	})

	return initErr
}

// Config holds configuration for primitive initialization
type Config struct {
	// Storage configuration
	StorageEndpoint string
	StorageTimeout  int

	// Echo configuration
	EchoEnabled bool

	// HTTP configuration
	HTTPBaseURL string
	HTTPTimeout int

	// Database configuration
	DBConnectionString string
	DBMaxConnections   int

	// Antifraud configuration
	AntifraudAPIKey                  string
	AntifraudAPIHost                 string
	AntifraudTimeout                 int
	AntifraudEnabled                 bool
	AntifraudMaxRetries              int
	AntifraudCircuitBreakerEnabled   bool
	AntifraudCircuitBreakerThreshold int
	AntifraudCircuitBreakerTimeout   int
}

// initStorageService initializes the storage service with proxy
func initStorageService(config *Config) (StorageService, error) {
	// Create the actual storage client
	storageClient := &storageClientImpl{
		endpoint: config.StorageEndpoint,
		timeout:  config.StorageTimeout,
	}

	// Wrap with proxy for logging and error handling
	storageProxy := NewStorageProxy(storageClient)
	return storageProxy, nil
}

// initEchoService initializes the echo service with proxy
func initEchoService(config *Config) (EchoService, error) {
	// Create the actual echo client
	echoClient := &echoClientImpl{
		enabled: config.EchoEnabled,
	}

	// Wrap with proxy
	echoProxy := NewEchoProxy(echoClient)
	return echoProxy, nil
}

// storageClientImpl is a simple implementation of StorageService
type storageClientImpl struct {
	endpoint string
	timeout  int
}

func (c *storageClientImpl) Save(data interface{}) (interface{}, error) {
	// Implementation would connect to actual storage
	return fmt.Sprintf("Saved data to %s: %v", c.endpoint, data), nil
}

func (c *storageClientImpl) Get(id string) (interface{}, error) {
	return fmt.Sprintf("Data from %s with ID: %s", c.endpoint, id), nil
}

func (c *storageClientImpl) Delete(id string) error {
	// Implementation would delete from storage
	return nil
}

func (c *storageClientImpl) List() ([]interface{}, error) {
	return []interface{}{"item1", "item2", "item3"}, nil
}

func (c *storageClientImpl) Update(id string, data interface{}) (interface{}, error) {
	return fmt.Sprintf("Updated %s at %s: %v", id, c.endpoint, data), nil
}

// echoClientImpl is a simple implementation of EchoService
type echoClientImpl struct {
	enabled bool
}

func (c *echoClientImpl) Echo(message string) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("echo service is disabled")
	}
	return fmt.Sprintf("Echo: %s", message), nil
}

func (c *echoClientImpl) Reverse(message string) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("echo service is disabled")
	}
	runes := []rune(message)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func (c *echoClientImpl) UpperCase(message string) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("echo service is disabled")
	}
	// Simple uppercase implementation
	var result []rune
	for _, r := range message {
		if r >= 'a' && r <= 'z' {
			result = append(result, r-32)
		} else {
			result = append(result, r)
		}
	}
	return string(result), nil
}

func (c *echoClientImpl) LowerCase(message string) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("echo service is disabled")
	}
	// Simple lowercase implementation
	var result []rune
	for _, r := range message {
		if r >= 'A' && r <= 'Z' {
			result = append(result, r+32)
		} else {
			result = append(result, r)
		}
	}
	return string(result), nil
}

// MustInit initializes the global primitive instance or panics on error
func MustInit(config *Config) {
	if err := Init(config); err != nil {
		panic(fmt.Sprintf("Failed to initialize primitive: %v", err))
	}
}

// IsInitialized returns true if the global primitive instance has been initialized
func IsInitialized() bool {
	return Default != nil
}

// initAntifraudService initializes the antifraud service with proxy
func initAntifraudService(config *Config) (AntifraudService, error) {
	// Create antifraud client config
	afConfig := models.ClientConfig{
		APIKey:                  config.AntifraudAPIKey,
		Host:                    config.AntifraudAPIHost,
		Timeout:                 config.AntifraudTimeout,
		Enabled:                 config.AntifraudEnabled,
		MaxRetries:              config.AntifraudMaxRetries,
		CircuitBreakerEnabled:   config.AntifraudCircuitBreakerEnabled,
		CircuitBreakerThreshold: config.AntifraudCircuitBreakerThreshold,
		CircuitBreakerTimeout:   config.AntifraudCircuitBreakerTimeout,
	}

	// Create the actual antifraud client from serviceclients
	antifraudClient, err := serviceclientsantifraud.NewClient(afConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create antifraud client: %w", err)
	}

	// Wrap with proxy for logging, metrics, and circuit breaking
	antifraudProxy := serviceclientsantifraud.NewProxy(antifraudClient, afConfig)

	// Create adapter to convert between concrete types and interface{}
	return &antifraudAdapter{
		service: antifraudProxy,
		config:  config,
	}, nil
}

// antifraudAdapter adapts the concrete antifraud service to the primitive.AntifraudService interface
type antifraudAdapter struct {
	service interface {
		StoreTransaction(models.AF_Transaction) error
		ValidateTransactionByAML(models.AF_Transaction) (models.ServiceResolution, error)
		ValidateTransactionByFC(models.AF_Transaction) (models.ServiceResolution, error)
		ValidateTransactionByML(models.AF_Transaction) (models.ServiceResolution, error)
		StoreServiceResolution(models.ServiceResolution) error
		AddTransactionServiceCheck(models.ServiceResolution) error
		FinalizeTransaction(models.AF_Transaction) (models.FinalResolution, error)
		StoreFinalResolution(models.FinalResolution) error
		HealthCheck() (bool, error)
		GetConfig() models.ClientConfig
	}
	config *Config
}

func (a *antifraudAdapter) StoreTransaction(afTransaction interface{}) error {
	// Convert interface{} to models.AF_Transaction
	tx, ok := afTransaction.(models.AF_Transaction)
	if !ok {
		// Try to convert from map or other types
		return fmt.Errorf("invalid transaction type: %T", afTransaction)
	}
	return a.service.StoreTransaction(tx)
}

func (a *antifraudAdapter) ValidateTransactionByAML(afTransaction interface{}) (interface{}, error) {
	tx, ok := afTransaction.(models.AF_Transaction)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type: %T", afTransaction)
	}
	return a.service.ValidateTransactionByAML(tx)
}

func (a *antifraudAdapter) ValidateTransactionByFC(afTransaction interface{}) (interface{}, error) {
	tx, ok := afTransaction.(models.AF_Transaction)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type: %T", afTransaction)
	}
	return a.service.ValidateTransactionByFC(tx)
}

func (a *antifraudAdapter) ValidateTransactionByML(afTransaction interface{}) (interface{}, error) {
	tx, ok := afTransaction.(models.AF_Transaction)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type: %T", afTransaction)
	}
	return a.service.ValidateTransactionByML(tx)
}

func (a *antifraudAdapter) StoreServiceResolution(resolution interface{}) error {
	res, ok := resolution.(models.ServiceResolution)
	if !ok {
		return fmt.Errorf("invalid resolution type: %T", resolution)
	}
	return a.service.StoreServiceResolution(res)
}

func (a *antifraudAdapter) AddTransactionServiceCheck(resolution interface{}) error {
	res, ok := resolution.(models.ServiceResolution)
	if !ok {
		return fmt.Errorf("invalid resolution type: %T", resolution)
	}
	return a.service.AddTransactionServiceCheck(res)
}

func (a *antifraudAdapter) FinalizeTransaction(afTransaction interface{}) (interface{}, error) {
	tx, ok := afTransaction.(models.AF_Transaction)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type: %T", afTransaction)
	}
	return a.service.FinalizeTransaction(tx)
}

func (a *antifraudAdapter) StoreFinalResolution(resolution interface{}) error {
	res, ok := resolution.(models.FinalResolution)
	if !ok {
		return fmt.Errorf("invalid resolution type: %T", resolution)
	}
	return a.service.StoreFinalResolution(res)
}

func (a *antifraudAdapter) HealthCheck() (bool, error) {
	return a.service.HealthCheck()
}

func (a *antifraudAdapter) GetConfig() interface{} {
	return a.config
}

// ResetForTesting resets the global instance (for testing only)
func ResetForTesting() {
	Default = nil
	initOnce = sync.Once{}
}
