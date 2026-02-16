package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server              ServerConfig              `yaml:"server"`
	Queue               QueueConfig               `yaml:"queue"`
	Executor            ExecutorConfig            `yaml:"executor"`
	Logging             LoggingConfig             `yaml:"logging"`
	DependencyInjection DependencyInjectionConfig `yaml:"dependency_injection"`
	Services            ServicesConfig            `yaml:"services"`
	Clients             ClientsConfig             `yaml:"clients"`
	Primitives          PrimitivesConfig          `yaml:"primitives"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int `yaml:"port"`
}

// QueueConfig represents queue configuration
type QueueConfig struct {
	Type string     `yaml:"type"`
	NATS NATSConfig `yaml:"nats"`
}

// NATSConfig represents NATS JetStream configuration
type NATSConfig struct {
	URLs           []string      `yaml:"urls"`
	StreamName     string        `yaml:"stream_name"`
	SubjectPrefix  string        `yaml:"subject_prefix"`
	DurableName    string        `yaml:"durable_name"`
	MaxReconnects  int           `yaml:"max_reconnects"`
	ReconnectWait  time.Duration `yaml:"reconnect_wait"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
}

// ExecutorConfig represents executor configuration
type ExecutorConfig struct {
	WorkerCount            int           `yaml:"worker_count"`
	QueuePollInterval      time.Duration `yaml:"queue_poll_interval"`
	MaxRetries             int           `yaml:"max_retries"`
	RetryDelay             time.Duration `yaml:"retry_delay"`
	ExecutionTimeout       time.Duration `yaml:"execution_timeout"`
	StepTimeout            time.Duration `yaml:"step_timeout"`
	EnableMetrics          bool          `yaml:"enable_metrics"`
	EnableTracing          bool          `yaml:"enable_tracing"`
	MaxConcurrentWorkflows int           `yaml:"max_concurrent_workflows"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DependencyInjectionConfig represents DI container configuration
type DependencyInjectionConfig struct {
	PoolSize                 int    `yaml:"pool_size"`
	EnableMetrics            bool   `yaml:"enable_metrics"`
	PreAllocatePools         bool   `yaml:"pre_allocate_pools"`
	MaxConcurrentResolutions int    `yaml:"max_concurrent_resolutions"`
	DefaultLifecycle         string `yaml:"default_lifecycle"`
}

// ServicesConfig represents external services configuration
type ServicesConfig struct {
	Registry RegistryServiceConfig `yaml:"registry"`
}

// RegistryServiceConfig represents registry service configuration
type RegistryServiceConfig struct {
	URL string `yaml:"url"`
}

// ClientsConfig represents client configurations
type ClientsConfig struct {
	Antifraud AntifraudClientConfig `yaml:"antifraud"`
	SDK       SDKClientConfig       `yaml:"sdk"`
}

// AntifraudClientConfig represents antifraud client configuration
type AntifraudClientConfig struct {
	APIKey                  string `yaml:"api_key"`
	Host                    string `yaml:"host"`
	Timeout                 int    `yaml:"timeout"`
	Enabled                 bool   `yaml:"enabled"`
	MaxRetries              int    `yaml:"max_retries"`
	CircuitBreakerEnabled   bool   `yaml:"circuit_breaker_enabled"`
	CircuitBreakerThreshold int    `yaml:"circuit_breaker_threshold"`
	CircuitBreakerTimeout   int    `yaml:"circuit_breaker_timeout"`
}

// SDKClientConfig represents SDK client configuration
type SDKClientConfig struct {
	WorkflowAPIEndpoint string `yaml:"workflow_api_endpoint"`
}

// PrimitivesConfig represents primitive services configuration
type PrimitivesConfig struct {
	EchoEnabled bool `yaml:"echo_enabled"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Queue: QueueConfig{
			Type: "in-memory",
			NATS: NATSConfig{
				URLs:           []string{"nats://localhost:4222"},
				StreamName:     "workflow-execution",
				SubjectPrefix:  "workflow",
				DurableName:    "workflow-consumer",
				MaxReconnects:  5,
				ReconnectWait:  2 * time.Second,
				ConnectTimeout: 5 * time.Second,
			},
		},
		Executor: ExecutorConfig{
			WorkerCount:            5,
			QueuePollInterval:      1 * time.Second,
			MaxRetries:             3,
			RetryDelay:             5 * time.Second,
			ExecutionTimeout:       5 * time.Minute,
			StepTimeout:            30 * time.Second,
			EnableMetrics:          true,
			EnableTracing:          false,
			MaxConcurrentWorkflows: 10,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		DependencyInjection: DependencyInjectionConfig{
			PoolSize:                 1000,
			EnableMetrics:            true,
			PreAllocatePools:         true,
			MaxConcurrentResolutions: 10000,
			DefaultLifecycle:         "singleton",
		},
		Services: ServicesConfig{
			Registry: RegistryServiceConfig{
				URL: "http://registry-service:8080",
			},
		},
		Clients: ClientsConfig{
			Antifraud: AntifraudClientConfig{
				APIKey:                  "",
				Host:                    "https://api.antifraudservice.com/v1",
				Timeout:                 30,
				Enabled:                 false,
				MaxRetries:              3,
				CircuitBreakerEnabled:   true,
				CircuitBreakerThreshold: 5,
				CircuitBreakerTimeout:   60,
			},
			SDK: SDKClientConfig{
				WorkflowAPIEndpoint: "http://localhost:8080",
			},
		},
		Primitives: PrimitivesConfig{
			EchoEnabled: true,
		},
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	// Try to load from config file
	configPath := getConfigPath()
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err == nil {
			var config Config
			if err := yaml.Unmarshal(data, &config); err != nil {
				return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
			}
			// Apply environment variable overrides
			applyEnvOverrides(&config)
			return &config, nil
		}
	}

	// Fall back to default configuration
	config := DefaultConfig()
	// Apply environment variable overrides
	applyEnvOverrides(config)
	return config, nil
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	// Check environment variable
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// Check common locations
	locations := []string{
		"./config.yaml",
		"./config.yml",
		"./config/config.yaml",
		"./config/config.yml",
		"/etc/unified-workflow/config.yaml",
	}

	for _, location := range locations {
		if _, err := os.Stat(location); err == nil {
			return location
		}
	}

	return ""
}

// applyEnvOverrides applies environment variable overrides to the configuration
func applyEnvOverrides(config *Config) {
	// Server configuration
	if val := os.Getenv("SERVER_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.Server.Port = port
		}
	}

	// Queue configuration
	if val := os.Getenv("QUEUE_TYPE"); val != "" {
		config.Queue.Type = val
	}

	// Registry service URL
	if val := os.Getenv("REGISTRY_SERVICE_URL"); val != "" {
		config.Services.Registry.URL = val
	}

	// Antifraud client configuration
	if val := os.Getenv("ANTIFRAUD_API_KEY"); val != "" {
		config.Clients.Antifraud.APIKey = val
	}
	if val := os.Getenv("ANTIFRAUD_HOST"); val != "" {
		config.Clients.Antifraud.Host = val
	}
	if val := os.Getenv("ANTIFRAUD_ENABLED"); val != "" {
		config.Clients.Antifraud.Enabled = strings.ToLower(val) == "true"
	}

	// SDK client configuration
	if val := os.Getenv("SDK_WORKFLOW_API_ENDPOINT"); val != "" {
		config.Clients.SDK.WorkflowAPIEndpoint = val
	}

	// DI configuration
	if val := os.Getenv("DI_POOL_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil {
			config.DependencyInjection.PoolSize = size
		}
	}
	if val := os.Getenv("DI_ENABLE_METRICS"); val != "" {
		config.DependencyInjection.EnableMetrics = strings.ToLower(val) == "true"
	}

	// Primitives configuration
	if val := os.Getenv("PRIMITIVES_ECHO_ENABLED"); val != "" {
		config.Primitives.EchoEnabled = strings.ToLower(val) == "true"
	}
}
