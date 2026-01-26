package primitive

// StorageService defines the interface for storage operations
type StorageService interface {
	// Save saves data to storage
	Save(data interface{}) (interface{}, error)

	// Get retrieves data from storage by ID
	Get(id string) (interface{}, error)

	// Delete removes data from storage by ID
	Delete(id string) error

	// List lists all items in storage
	List() ([]interface{}, error)

	// Update updates existing data in storage
	Update(id string, data interface{}) (interface{}, error)
}

// EchoService defines the interface for echo operations
type EchoService interface {
	// Echo returns the input message as-is
	Echo(message string) (string, error)

	// Reverse returns the input message reversed
	Reverse(message string) (string, error)

	// UpperCase returns the input message in uppercase
	UpperCase(message string) (string, error)

	// LowerCase returns the input message in lowercase
	LowerCase(message string) (string, error)
}

// HTTPService defines the interface for HTTP operations
type HTTPService interface {
	// Get performs an HTTP GET request
	Get(url string, headers map[string]string) (interface{}, error)

	// Post performs an HTTP POST request
	Post(url string, body interface{}, headers map[string]string) (interface{}, error)

	// Put performs an HTTP PUT request
	Put(url string, body interface{}, headers map[string]string) (interface{}, error)

	// Delete performs an HTTP DELETE request
	Delete(url string, headers map[string]string) (interface{}, error)
}

// DatabaseService defines the interface for database operations
type DatabaseService interface {
	// Query executes a database query
	Query(query string, params map[string]interface{}) ([]map[string]interface{}, error)

	// Execute executes a database command (INSERT, UPDATE, DELETE)
	Execute(command string, params map[string]interface{}) (int64, error)

	// Transaction executes operations within a transaction
	Transaction(operations func(DatabaseService) error) error
}

// ServiceRegistry provides access to all services
type ServiceRegistry interface {
	// Storage returns the storage service
	Storage() StorageService

	// Echo returns the echo service
	Echo() EchoService

	// HTTP returns the HTTP service
	HTTP() HTTPService

	// Database returns the database service
	Database() DatabaseService

	// RegisterService registers a custom service
	RegisterService(name string, service interface{}) error

	// GetService retrieves a service by name
	GetService(name string) (interface{}, error)
}
