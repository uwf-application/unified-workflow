package config

import (
	"time"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Queue    QueueConfig    `yaml:"queue"`
	Executor ExecutorConfig `yaml:"executor"`
	Logging  LoggingConfig  `yaml:"logging"`
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
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	// For now, return default configuration
	// In a real implementation, this would load from config file and env vars
	return DefaultConfig(), nil
}
