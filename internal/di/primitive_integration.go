package di

import (
	"unified-workflow/internal/primitive"
)

// PrimitiveProvider provides access to primitive services through DI
type PrimitiveProvider struct {
	container Container
}

// NewPrimitiveProvider creates a new primitive provider
func NewPrimitiveProvider(container Container) *PrimitiveProvider {
	return &PrimitiveProvider{
		container: container,
	}
}

// Storage returns the storage service
func (p *PrimitiveProvider) Storage() primitive.StorageService {
	instance, err := p.container.Resolve((*primitive.StorageService)(nil))
	if err != nil {
		// Fall back to global primitive if not registered
		return primitive.Default.Storage
	}
	return instance.(primitive.StorageService)
}

// Echo returns the echo service
func (p *PrimitiveProvider) Echo() primitive.EchoService {
	instance, err := p.container.Resolve((*primitive.EchoService)(nil))
	if err != nil {
		// Fall back to global primitive if not registered
		return primitive.Default.Echo
	}
	return instance.(primitive.EchoService)
}

// HTTP returns the HTTP service
func (p *PrimitiveProvider) HTTP() primitive.HTTPService {
	instance, err := p.container.Resolve((*primitive.HTTPService)(nil))
	if err != nil {
		// Fall back to global primitive if not registered
		return primitive.Default.HTTP
	}
	return instance.(primitive.HTTPService)
}

// DB returns the database service
func (p *PrimitiveProvider) DB() primitive.DatabaseService {
	instance, err := p.container.Resolve((*primitive.DatabaseService)(nil))
	if err != nil {
		// Fall back to global primitive if not registered
		return primitive.Default.DB
	}
	return instance.(primitive.DatabaseService)
}

// Antifraud returns the antifraud service
func (p *PrimitiveProvider) Antifraud() primitive.AntifraudService {
	instance, err := p.container.Resolve((*primitive.AntifraudService)(nil))
	if err != nil {
		// Fall back to global primitive if not registered
		return primitive.Default.Antifraud
	}
	return instance.(primitive.AntifraudService)
}

// RegisterPrimitiveServices registers all primitive services with the container
func RegisterPrimitiveServices(container Container, config *primitive.Config) error {
	// Register storage service
	err := container.RegisterFactory((*primitive.StorageService)(nil), func(c Container) (interface{}, error) {
		return primitive.Default.Storage, nil
	}, Singleton)
	if err != nil {
		return err
	}

	// Register echo service
	err = container.RegisterFactory((*primitive.EchoService)(nil), func(c Container) (interface{}, error) {
		return primitive.Default.Echo, nil
	}, Singleton)
	if err != nil {
		return err
	}

	// Register HTTP service (if available)
	err = container.RegisterFactory((*primitive.HTTPService)(nil), func(c Container) (interface{}, error) {
		return primitive.Default.HTTP, nil
	}, Singleton)
	if err != nil {
		return err
	}

	// Register DB service (if available)
	err = container.RegisterFactory((*primitive.DatabaseService)(nil), func(c Container) (interface{}, error) {
		return primitive.Default.DB, nil
	}, Singleton)
	if err != nil {
		return err
	}

	// Register antifraud service
	err = container.RegisterFactory((*primitive.AntifraudService)(nil), func(c Container) (interface{}, error) {
		return primitive.Default.Antifraud, nil
	}, Singleton)
	if err != nil {
		return err
	}

	return nil
}
