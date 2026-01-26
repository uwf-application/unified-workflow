package primitive

import (
	"fmt"
	"time"
)

// ProxyLogger defines the logging interface for proxies
type ProxyLogger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}

// DefaultProxyLogger is a simple logger implementation for proxies
type DefaultProxyLogger struct{}

func (l *DefaultProxyLogger) Debug(msg string, fields map[string]interface{}) {
	fmt.Printf("[DEBUG] %s: %v\n", msg, fields)
}

func (l *DefaultProxyLogger) Info(msg string, fields map[string]interface{}) {
	fmt.Printf("[INFO] %s: %v\n", msg, fields)
}

func (l *DefaultProxyLogger) Warn(msg string, fields map[string]interface{}) {
	fmt.Printf("[WARN] %s: %v\n", msg, fields)
}

func (l *DefaultProxyLogger) Error(msg string, fields map[string]interface{}) {
	fmt.Printf("[ERROR] %s: %v\n", msg, fields)
}

// StorageProxy is a proxy layer for storage operations
type StorageProxy struct {
	executor StorageService
	logger   ProxyLogger
}

// NewStorageProxy creates a new StorageProxy
func NewStorageProxy(executor StorageService) *StorageProxy {
	return &StorageProxy{
		executor: executor,
		logger:   &DefaultProxyLogger{},
	}
}

// Save implements StorageService.Save with proxy functionality
func (p *StorageProxy) Save(data interface{}) (interface{}, error) {
	startTime := time.Now()
	operation := "Save"

	// Log the operation start
	p.logger.Info("Executing storage operation", map[string]interface{}{
		"operation": operation,
		"data_type": fmt.Sprintf("%T", data),
	})

	// Validate input
	if data == nil {
		err := fmt.Errorf("data cannot be nil")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return nil, err
	}

	// Execute the operation
	result, err := p.executor.Save(data)
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Storage operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return nil, fmt.Errorf("storage.Save failed: %w", err)
	}

	// Log success
	p.logger.Info("Storage operation completed", map[string]interface{}{
		"operation": operation,
		"duration":  duration.String(),
		"result":    result,
	})

	return result, nil
}

// Get implements StorageService.Get with proxy functionality
func (p *StorageProxy) Get(id string) (interface{}, error) {
	startTime := time.Now()
	operation := "Get"

	// Log the operation start
	p.logger.Info("Executing storage operation", map[string]interface{}{
		"operation": operation,
		"id":        id,
	})

	// Validate input
	if id == "" {
		err := fmt.Errorf("id cannot be empty")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return nil, err
	}

	// Execute the operation
	result, err := p.executor.Get(id)
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Storage operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return nil, fmt.Errorf("storage.Get failed: %w", err)
	}

	// Log success
	p.logger.Info("Storage operation completed", map[string]interface{}{
		"operation":  operation,
		"duration":   duration.String(),
		"has_result": result != nil,
	})

	return result, nil
}

// Delete implements StorageService.Delete with proxy functionality
func (p *StorageProxy) Delete(id string) error {
	startTime := time.Now()
	operation := "Delete"

	// Log the operation start
	p.logger.Info("Executing storage operation", map[string]interface{}{
		"operation": operation,
		"id":        id,
	})

	// Validate input
	if id == "" {
		err := fmt.Errorf("id cannot be empty")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return err
	}

	// Execute the operation
	err := p.executor.Delete(id)
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Storage operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return fmt.Errorf("storage.Delete failed: %w", err)
	}

	// Log success
	p.logger.Info("Storage operation completed", map[string]interface{}{
		"operation": operation,
		"duration":  duration.String(),
	})

	return nil
}

// List implements StorageService.List with proxy functionality
func (p *StorageProxy) List() ([]interface{}, error) {
	startTime := time.Now()
	operation := "List"

	// Log the operation start
	p.logger.Info("Executing storage operation", map[string]interface{}{
		"operation": operation,
	})

	// Execute the operation
	result, err := p.executor.List()
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Storage operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return nil, fmt.Errorf("storage.List failed: %w", err)
	}

	// Log success
	p.logger.Info("Storage operation completed", map[string]interface{}{
		"operation":  operation,
		"duration":   duration.String(),
		"item_count": len(result),
	})

	return result, nil
}

// Update implements StorageService.Update with proxy functionality
func (p *StorageProxy) Update(id string, data interface{}) (interface{}, error) {
	startTime := time.Now()
	operation := "Update"

	// Log the operation start
	p.logger.Info("Executing storage operation", map[string]interface{}{
		"operation": operation,
		"id":        id,
		"data_type": fmt.Sprintf("%T", data),
	})

	// Validate input
	if id == "" {
		err := fmt.Errorf("id cannot be empty")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return nil, err
	}

	if data == nil {
		err := fmt.Errorf("data cannot be nil")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return nil, err
	}

	// Execute the operation
	result, err := p.executor.Update(id, data)
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Storage operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return nil, fmt.Errorf("storage.Update failed: %w", err)
	}

	// Log success
	p.logger.Info("Storage operation completed", map[string]interface{}{
		"operation": operation,
		"duration":  duration.String(),
		"result":    result,
	})

	return result, nil
}

// EchoProxy is a proxy layer for echo operations
type EchoProxy struct {
	executor EchoService
	logger   ProxyLogger
}

// NewEchoProxy creates a new EchoProxy
func NewEchoProxy(executor EchoService) *EchoProxy {
	return &EchoProxy{
		executor: executor,
		logger:   &DefaultProxyLogger{},
	}
}

// Echo implements EchoService.Echo with proxy functionality
func (p *EchoProxy) Echo(message string) (string, error) {
	startTime := time.Now()
	operation := "Echo"

	// Log the operation start
	p.logger.Info("Executing echo operation", map[string]interface{}{
		"operation": operation,
		"message":   message,
	})

	// Validate input
	if message == "" {
		err := fmt.Errorf("message cannot be empty")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return "", err
	}

	// Execute the operation
	result, err := p.executor.Echo(message)
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Echo operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return "", fmt.Errorf("echo.Echo failed: %w", err)
	}

	// Log success
	p.logger.Info("Echo operation completed", map[string]interface{}{
		"operation": operation,
		"duration":  duration.String(),
		"result":    result,
	})

	return result, nil
}

// Reverse implements EchoService.Reverse with proxy functionality
func (p *EchoProxy) Reverse(message string) (string, error) {
	startTime := time.Now()
	operation := "Reverse"

	// Log the operation start
	p.logger.Info("Executing echo operation", map[string]interface{}{
		"operation": operation,
		"message":   message,
	})

	// Validate input
	if message == "" {
		err := fmt.Errorf("message cannot be empty")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return "", err
	}

	// Execute the operation
	result, err := p.executor.Reverse(message)
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Echo operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return "", fmt.Errorf("echo.Reverse failed: %w", err)
	}

	// Log success
	p.logger.Info("Echo operation completed", map[string]interface{}{
		"operation": operation,
		"duration":  duration.String(),
		"result":    result,
	})

	return result, nil
}

// UpperCase implements EchoService.UpperCase with proxy functionality
func (p *EchoProxy) UpperCase(message string) (string, error) {
	startTime := time.Now()
	operation := "UpperCase"

	// Log the operation start
	p.logger.Info("Executing echo operation", map[string]interface{}{
		"operation": operation,
		"message":   message,
	})

	// Validate input
	if message == "" {
		err := fmt.Errorf("message cannot be empty")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return "", err
	}

	// Execute the operation
	result, err := p.executor.UpperCase(message)
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Echo operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return "", fmt.Errorf("echo.UpperCase failed: %w", err)
	}

	// Log success
	p.logger.Info("Echo operation completed", map[string]interface{}{
		"operation": operation,
		"duration":  duration.String(),
		"result":    result,
	})

	return result, nil
}

// LowerCase implements EchoService.LowerCase with proxy functionality
func (p *EchoProxy) LowerCase(message string) (string, error) {
	startTime := time.Now()
	operation := "LowerCase"

	// Log the operation start
	p.logger.Info("Executing echo operation", map[string]interface{}{
		"operation": operation,
		"message":   message,
	})

	// Validate input
	if message == "" {
		err := fmt.Errorf("message cannot be empty")
		p.logger.Error("Validation failed", map[string]interface{}{
			"operation": operation,
			"error":     err.Error(),
		})
		return "", err
	}

	// Execute the operation
	result, err := p.executor.LowerCase(message)
	duration := time.Since(startTime)

	if err != nil {
		// Log error
		p.logger.Error("Echo operation failed", map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
			"error":     err.Error(),
		})
		return "", fmt.Errorf("echo.LowerCase failed: %w", err)
	}

	// Log success
	p.logger.Info("Echo operation completed", map[string]interface{}{
		"operation": operation,
		"duration":  duration.String(),
		"result":    result,
	})

	return result, nil
}
