package model

// WorkflowData represents shared workflow data (mutable)
// Similar to Java's WorkflowData class in workflow-primitive
type WorkflowData interface {
	// Get returns a value by key
	Get(key string) interface{}

	// Put sets a value by key
	Put(key string, value interface{})

	// Remove removes a value by key
	Remove(key string)

	// Contains checks if a key exists
	Contains(key string) bool

	// Size returns the number of entries
	Size() int

	// Clear removes all entries
	Clear()

	// ToMap returns a copy of the data as a map
	ToMap() map[string]interface{}

	// DeepCopy creates a deep copy of the workflow data
	DeepCopy() WorkflowData

	// Merge merges another WorkflowData into this one
	Merge(other WorkflowData)

	// GetString returns a string value by key
	GetString(key string) (string, bool)

	// GetInt returns an int value by key
	GetInt(key string) (int, bool)

	// GetBool returns a bool value by key
	GetBool(key string) (bool, bool)

	// GetFloat returns a float64 value by key
	GetFloat(key string) (float64, bool)

	// GetMap returns a map value by key
	GetMap(key string) (map[string]interface{}, bool)

	// GetSlice returns a slice value by key
	GetSlice(key string) ([]interface{}, bool)
}

// WorkflowDataImpl implements the WorkflowData interface
type WorkflowDataImpl struct {
	data map[string]interface{}
}

// NewWorkflowData creates a new workflow data instance
func NewWorkflowData() *WorkflowDataImpl {
	return &WorkflowDataImpl{
		data: make(map[string]interface{}),
	}
}

// Get returns a value by key
func (wd *WorkflowDataImpl) Get(key string) interface{} {
	return wd.data[key]
}

// Put sets a value by key
func (wd *WorkflowDataImpl) Put(key string, value interface{}) {
	wd.data[key] = value
}

// Remove removes a value by key
func (wd *WorkflowDataImpl) Remove(key string) {
	delete(wd.data, key)
}

// Contains checks if a key exists
func (wd *WorkflowDataImpl) Contains(key string) bool {
	_, exists := wd.data[key]
	return exists
}

// Size returns the number of entries
func (wd *WorkflowDataImpl) Size() int {
	return len(wd.data)
}

// Clear removes all entries
func (wd *WorkflowDataImpl) Clear() {
	wd.data = make(map[string]interface{})
}

// ToMap returns a copy of the data as a map
func (wd *WorkflowDataImpl) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range wd.data {
		result[k] = v
	}
	return result
}

// DeepCopy creates a deep copy of the workflow data
func (wd *WorkflowDataImpl) DeepCopy() *WorkflowDataImpl {
	copy := NewWorkflowData()
	for k, v := range wd.data {
		// Simple copy - for complex nested structures, a more sophisticated copy would be needed
		copy.Put(k, v)
	}
	return copy
}

// Merge merges another WorkflowData into this one
func (wd *WorkflowDataImpl) Merge(other *WorkflowDataImpl) {
	for k, v := range other.data {
		wd.Put(k, v)
	}
}

// GetString returns a string value by key
func (wd *WorkflowDataImpl) GetString(key string) (string, bool) {
	value, exists := wd.data[key]
	if !exists {
		return "", false
	}
	str, ok := value.(string)
	return str, ok
}

// GetInt returns an int value by key
func (wd *WorkflowDataImpl) GetInt(key string) (int, bool) {
	value, exists := wd.data[key]
	if !exists {
		return 0, false
	}
	// Handle different numeric types
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// GetBool returns a bool value by key
func (wd *WorkflowDataImpl) GetBool(key string) (bool, bool) {
	value, exists := wd.data[key]
	if !exists {
		return false, false
	}
	b, ok := value.(bool)
	return b, ok
}

// GetFloat returns a float64 value by key
func (wd *WorkflowDataImpl) GetFloat(key string) (float64, bool) {
	value, exists := wd.data[key]
	if !exists {
		return 0, false
	}
	// Handle different numeric types
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

// GetMap returns a map value by key
func (wd *WorkflowDataImpl) GetMap(key string) (map[string]interface{}, bool) {
	value, exists := wd.data[key]
	if !exists {
		return nil, false
	}
	m, ok := value.(map[string]interface{})
	return m, ok
}

// GetSlice returns a slice value by key
func (wd *WorkflowDataImpl) GetSlice(key string) ([]interface{}, bool) {
	value, exists := wd.data[key]
	if !exists {
		return nil, false
	}
	s, ok := value.([]interface{})
	return s, ok
}
