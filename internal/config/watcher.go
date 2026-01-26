package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// ConfigChange represents a configuration change
type ConfigChange struct {
	// Path of the changed config file
	Path string

	// Old configuration (before change)
	OldConfig *Config

	// New configuration (after change)
	NewConfig *Config

	// Timestamp of change
	Timestamp time.Time

	// Error if config loading failed
	Error error
}

// ConfigWatcher watches configuration files for changes
type ConfigWatcher struct {
	// Paths to watch
	paths []string

	// Current configuration
	currentConfig *Config

	// Change subscribers
	subscribers []ConfigChangeSubscriber

	// Mutex for thread safety
	mu sync.RWMutex

	// Stop channel
	stopChan chan struct{}

	// Is running
	running bool

	// File modification times
	fileModTimes map[string]time.Time
}

// ConfigChangeSubscriber subscribes to configuration changes
type ConfigChangeSubscriber interface {
	// OnConfigChange is called when configuration changes
	OnConfigChange(change ConfigChange)
}

// ConfigChangeFunc is a function that handles configuration changes
type ConfigChangeFunc func(change ConfigChange)

// NewConfigWatcher creates a new configuration watcher
func NewConfigWatcher(configPaths ...string) (*ConfigWatcher, error) {
	cw := &ConfigWatcher{
		paths:        configPaths,
		stopChan:     make(chan struct{}),
		subscribers:  make([]ConfigChangeSubscriber, 0),
		fileModTimes: make(map[string]time.Time),
	}

	// Load initial configuration
	if err := cw.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load initial config: %w", err)
	}

	return cw, nil
}

// Start starts watching for configuration changes
func (cw *ConfigWatcher) Start() error {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	if cw.running {
		return fmt.Errorf("watcher already running")
	}

	// Initialize file modification times
	for _, path := range cw.paths {
		if info, err := os.Stat(path); err == nil {
			cw.fileModTimes[path] = info.ModTime()
		}
	}

	cw.running = true

	// Start watching in goroutine
	go cw.watch()

	log.Println("Configuration watcher started")
	return nil
}

// Stop stops watching for configuration changes
func (cw *ConfigWatcher) Stop() error {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	if !cw.running {
		return fmt.Errorf("watcher not running")
	}

	close(cw.stopChan)
	cw.running = false

	log.Println("Configuration watcher stopped")
	return nil
}

// watch watches for file system changes using polling
func (cw *ConfigWatcher) watch() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cw.checkForChanges()

		case <-cw.stopChan:
			return
		}
	}
}

// checkForChanges checks if any config files have changed
func (cw *ConfigWatcher) checkForChanges() {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	for _, path := range cw.paths {
		info, err := os.Stat(path)
		if err != nil {
			// File might have been deleted
			if _, exists := cw.fileModTimes[path]; exists {
				delete(cw.fileModTimes, path)
				log.Printf("Config file removed: %s", path)
			}
			continue
		}

		newModTime := info.ModTime()
		oldModTime, exists := cw.fileModTimes[path]

		if !exists || newModTime.After(oldModTime) {
			// File has changed
			cw.fileModTimes[path] = newModTime
			cw.reloadConfig(path)
		}
	}
}

// reloadConfig reloads configuration from file
func (cw *ConfigWatcher) reloadConfig(filePath string) {
	oldConfig := cw.currentConfig
	newConfig, err := cw.loadConfigFromFile(filePath)

	change := ConfigChange{
		Path:      filePath,
		OldConfig: oldConfig,
		NewConfig: newConfig,
		Timestamp: time.Now(),
		Error:     err,
	}

	// Update current config if loaded successfully
	if err == nil && newConfig != nil {
		cw.currentConfig = newConfig
		log.Printf("Configuration reloaded successfully from %s", filePath)
	} else if err != nil {
		log.Printf("Failed to reload configuration from %s: %v", filePath, err)
	}

	// Notify subscribers
	cw.notifySubscribers(change)
}

// loadConfig loads configuration from all watched files
func (cw *ConfigWatcher) loadConfig() error {
	// Try to load from each config file
	for _, path := range cw.paths {
		config, err := cw.loadConfigFromFile(path)
		if err == nil && config != nil {
			cw.currentConfig = config
			log.Printf("Loaded configuration from %s", path)
			return nil
		}
	}

	return fmt.Errorf("failed to load configuration from any watched file")
}

// loadConfigFromFile loads configuration from a specific file
func (cw *ConfigWatcher) loadConfigFromFile(filePath string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", filePath)
	}

	// Use the existing LoadConfig function
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return config, nil
}

// Subscribe subscribes to configuration changes
func (cw *ConfigWatcher) Subscribe(subscriber ConfigChangeSubscriber) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.subscribers = append(cw.subscribers, subscriber)
}

// SubscribeFunc subscribes a function to configuration changes
func (cw *ConfigWatcher) SubscribeFunc(fn ConfigChangeFunc) {
	cw.Subscribe(&funcSubscriber{fn: fn})
}

// Unsubscribe unsubscribes from configuration changes
func (cw *ConfigWatcher) Unsubscribe(subscriber ConfigChangeSubscriber) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	for i, sub := range cw.subscribers {
		if sub == subscriber {
			cw.subscribers = append(cw.subscribers[:i], cw.subscribers[i+1:]...)
			break
		}
	}
}

// notifySubscribers notifies all subscribers of a configuration change
func (cw *ConfigWatcher) notifySubscribers(change ConfigChange) {
	cw.mu.RLock()
	subscribers := make([]ConfigChangeSubscriber, len(cw.subscribers))
	copy(subscribers, cw.subscribers)
	cw.mu.RUnlock()

	for _, subscriber := range subscribers {
		go func(s ConfigChangeSubscriber) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic in config change subscriber: %v", r)
				}
			}()
			s.OnConfigChange(change)
		}(subscriber)
	}
}

// GetCurrentConfig returns the current configuration
func (cw *ConfigWatcher) GetCurrentConfig() *Config {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	return cw.currentConfig
}

// IsRunning returns true if the watcher is running
func (cw *ConfigWatcher) IsRunning() bool {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	return cw.running
}

// funcSubscriber wraps a function as a ConfigChangeSubscriber
type funcSubscriber struct {
	fn ConfigChangeFunc
}

// OnConfigChange calls the wrapped function
func (f *funcSubscriber) OnConfigChange(change ConfigChange) {
	f.fn(change)
}

// HotReloadConfig provides hot-reloadable configuration
type HotReloadConfig struct {
	watcher *ConfigWatcher
	mu      sync.RWMutex
	config  *Config
}

// NewHotReloadConfig creates a new hot-reloadable configuration
func NewHotReloadConfig(configPaths ...string) (*HotReloadConfig, error) {
	watcher, err := NewConfigWatcher(configPaths...)
	if err != nil {
		return nil, err
	}

	hrc := &HotReloadConfig{
		watcher: watcher,
		config:  watcher.GetCurrentConfig(),
	}

	// Subscribe to config changes
	watcher.SubscribeFunc(func(change ConfigChange) {
		if change.Error == nil && change.NewConfig != nil {
			hrc.mu.Lock()
			hrc.config = change.NewConfig
			hrc.mu.Unlock()
			log.Println("Hot-reloaded configuration updated")
		}
	})

	return hrc, nil
}

// Start starts the hot-reload configuration
func (h *HotReloadConfig) Start() error {
	return h.watcher.Start()
}

// Stop stops the hot-reload configuration
func (h *HotReloadConfig) Stop() error {
	return h.watcher.Stop()
}

// Get returns the current configuration
func (h *HotReloadConfig) Get() *Config {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config
}

// GetWithFallback returns the current configuration or a fallback
func (h *HotReloadConfig) GetWithFallback(fallback *Config) *Config {
	config := h.Get()
	if config == nil {
		return fallback
	}
	return config
}

// WatchForChanges watches for specific configuration changes
func (h *HotReloadConfig) WatchForChanges(callback func(oldConfig, newConfig *Config)) {
	h.watcher.SubscribeFunc(func(change ConfigChange) {
		if change.Error == nil {
			callback(change.OldConfig, change.NewConfig)
		}
	})
}

// ConfigAware interface for components that need to know about config changes
type ConfigAware interface {
	// OnConfigChange is called when configuration changes
	OnConfigChange(oldConfig, newConfig *Config)
}

// Simple example of using hot-reload configuration
func ExampleHotReload() {
	// Create hot-reload configuration
	hotConfig, err := NewHotReloadConfig("./config.yaml", "./config.local.yaml")
	if err != nil {
		log.Fatalf("Failed to create hot-reload config: %v", err)
	}

	// Start watching for changes
	if err := hotConfig.Start(); err != nil {
		log.Fatalf("Failed to start hot-reload: %v", err)
	}
	defer hotConfig.Stop()

	// Use configuration
	config := hotConfig.Get()
	log.Printf("Using configuration: %+v", config)

	// Keep running
	select {}
}
