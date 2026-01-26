package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DIConfig holds DI container configuration
type DIConfig struct {
	// Container settings
	PoolSize                 int  `yaml:"pool_size" json:"pool_size"`
	EnableMetrics            bool `yaml:"enable_metrics" json:"enable_metrics"`
	PreAllocatePools         bool `yaml:"pre_allocate_pools" json:"pre_allocate_pools"`
	MaxConcurrentResolutions int  `yaml:"max_concurrent_resolutions" json:"max_concurrent_resolutions"`

	// Lifecycle defaults
	DefaultLifecycle string `yaml:"default_lifecycle" json:"default_lifecycle"` // singleton, transient, scoped, pooled

	// Component configurations
	Components map[string]ComponentConfig `yaml:"components" json:"components"`
}

// ComponentConfig holds configuration for individual components
type ComponentConfig struct {
	Lifecycle string                 `yaml:"lifecycle" json:"lifecycle"`
	Factory   string                 `yaml:"factory" json:"factory"` // function name or path
	Params    map[string]interface{} `yaml:"params" json:"params"`
	Enabled   bool                   `yaml:"enabled" json:"enabled"`
	Priority  int                    `yaml:"priority" json:"priority"` // load order
}

// DefaultDIConfig returns the default DI configuration
func DefaultDIConfig() *DIConfig {
	return &DIConfig{
		PoolSize:                 1000,
		EnableMetrics:            true,
		PreAllocatePools:         true,
		MaxConcurrentResolutions: 10000,
		DefaultLifecycle:         "singleton",
		Components: map[string]ComponentConfig{
			"registry": {
				Lifecycle: "singleton",
				Factory:   "NewInMemoryRegistry",
				Enabled:   true,
				Priority:  1,
			},
			"state": {
				Lifecycle: "singleton",
				Factory:   "NewInMemoryState",
				Enabled:   true,
				Priority:  2,
			},
			"queue": {
				Lifecycle: "singleton",
				Factory:   "createQueueService",
				Enabled:   true,
				Priority:  3,
				Params: map[string]interface{}{
					"type": "in-memory",
				},
			},
			"executor": {
				Lifecycle: "singleton",
				Factory:   "CreateWorkflowExecutor",
				Enabled:   true,
				Priority:  4,
			},
			"primitives": {
				Lifecycle: "singleton",
				Factory:   "RegisterPrimitiveServices",
				Enabled:   true,
				Priority:  5,
			},
		},
	}
}

// LoadDIConfig loads DI configuration from file
func LoadDIConfig(configPath string) (*DIConfig, error) {
	// Try to load from file
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err == nil {
			var config DIConfig
			if err := yaml.Unmarshal(data, &config); err != nil {
				return nil, fmt.Errorf("failed to parse DI config: %w", err)
			}
			return &config, nil
		}
	}

	// Fall back to default
	return DefaultDIConfig(), nil
}

// LoadDIConfigFromEnv loads DI configuration from environment
func LoadDIConfigFromEnv() (*DIConfig, error) {
	config := DefaultDIConfig()

	// Override from environment variables
	if val := os.Getenv("DI_POOL_SIZE"); val != "" {
		if size, err := parseInt(val); err == nil {
			config.PoolSize = size
		}
	}

	if val := os.Getenv("DI_ENABLE_METRICS"); val != "" {
		config.EnableMetrics = strings.ToLower(val) == "true"
	}

	if val := os.Getenv("DI_PRE_ALLOCATE_POOLS"); val != "" {
		config.PreAllocatePools = strings.ToLower(val) == "true"
	}

	if val := os.Getenv("DI_DEFAULT_LIFECYCLE"); val != "" {
		config.DefaultLifecycle = val
	}

	return config, nil
}

// GetConfigPath returns the DI configuration file path
func GetConfigPath() string {
	// Check environment variable
	if path := os.Getenv("DI_CONFIG_PATH"); path != "" {
		return path
	}

	// Check common locations
	locations := []string{
		"./config/di.yaml",
		"./config/di.yml",
		"./di.yaml",
		"./di.yml",
		"/etc/unified-workflow/di.yaml",
	}

	for _, location := range locations {
		if _, err := os.Stat(location); err == nil {
			return location
		}
	}

	return ""
}

// MergeConfigs merges multiple DI configurations
func MergeConfigs(configs ...*DIConfig) *DIConfig {
	if len(configs) == 0 {
		return DefaultDIConfig()
	}

	result := &DIConfig{
		PoolSize:                 configs[0].PoolSize,
		EnableMetrics:            configs[0].EnableMetrics,
		PreAllocatePools:         configs[0].PreAllocatePools,
		MaxConcurrentResolutions: configs[0].MaxConcurrentResolutions,
		DefaultLifecycle:         configs[0].DefaultLifecycle,
		Components:               make(map[string]ComponentConfig),
	}

	// Merge components from all configs
	for _, config := range configs {
		if config == nil {
			continue
		}

		// Update scalar values (last config wins)
		if config.PoolSize > 0 {
			result.PoolSize = config.PoolSize
		}
		result.EnableMetrics = config.EnableMetrics
		result.PreAllocatePools = config.PreAllocatePools
		if config.MaxConcurrentResolutions > 0 {
			result.MaxConcurrentResolutions = config.MaxConcurrentResolutions
		}
		if config.DefaultLifecycle != "" {
			result.DefaultLifecycle = config.DefaultLifecycle
		}

		// Merge components
		for name, component := range config.Components {
			result.Components[name] = component
		}
	}

	return result
}

// SaveConfig saves DI configuration to file
func (c *DIConfig) SaveConfig(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the DI configuration
func (c *DIConfig) Validate() error {
	if c.PoolSize <= 0 {
		return fmt.Errorf("pool_size must be positive")
	}

	if c.MaxConcurrentResolutions <= 0 {
		return fmt.Errorf("max_concurrent_resolutions must be positive")
	}

	validLifecycles := map[string]bool{
		"singleton": true,
		"transient": true,
		"scoped":    true,
		"pooled":    true,
	}

	if !validLifecycles[strings.ToLower(c.DefaultLifecycle)] {
		return fmt.Errorf("invalid default_lifecycle: %s", c.DefaultLifecycle)
	}

	for name, component := range c.Components {
		if !validLifecycles[strings.ToLower(component.Lifecycle)] {
			return fmt.Errorf("invalid lifecycle for component %s: %s", name, component.Lifecycle)
		}

		if component.Factory == "" {
			return fmt.Errorf("factory must be specified for component %s", name)
		}
	}

	return nil
}

// GetComponentConfig returns configuration for a specific component
func (c *DIConfig) GetComponentConfig(name string) (ComponentConfig, bool) {
	config, exists := c.Components[name]
	return config, exists
}

// SetComponentConfig sets configuration for a component
func (c *DIConfig) SetComponentConfig(name string, config ComponentConfig) {
	if c.Components == nil {
		c.Components = make(map[string]ComponentConfig)
	}
	c.Components[name] = config
}

// EnableComponent enables a component
func (c *DIConfig) EnableComponent(name string) {
	if config, exists := c.Components[name]; exists {
		config.Enabled = true
		c.Components[name] = config
	}
}

// DisableComponent disables a component
func (c *DIConfig) DisableComponent(name string) {
	if config, exists := c.Components[name]; exists {
		config.Enabled = false
		c.Components[name] = config
	}
}

// Helper function to parse integer from string
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
