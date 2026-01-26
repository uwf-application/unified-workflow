package model

// ChildStep represents a granular execution unit with hooks
// Similar to Java's ChildStep class in workflow-common
type ChildStep struct {
	name         string
	requestHook  func(context interface{}, data interface{}) interface{}
	responseHook func(context interface{}, data interface{}) interface{}
	validateHook func(response interface{}) error
}

// NewChildStep creates a new ChildStep
func NewChildStep(name string, requestHook, responseHook func(context interface{}, data interface{}) interface{}, validateHook func(response interface{}) error) *ChildStep {
	return &ChildStep{
		name:         name,
		requestHook:  requestHook,
		responseHook: responseHook,
		validateHook: validateHook,
	}
}

// GetName returns the name of the child step
func (cs *ChildStep) GetName() string {
	return cs.name
}

// GetRequestHook returns the request hook function
func (cs *ChildStep) GetRequestHook() func(context interface{}, data interface{}) interface{} {
	return cs.requestHook
}

// GetResponseHook returns the response hook function
func (cs *ChildStep) GetResponseHook() func(context interface{}, data interface{}) interface{} {
	return cs.responseHook
}

// GetValidateHook returns the validate hook function
func (cs *ChildStep) GetValidateHook() func(response interface{}) error {
	return cs.validateHook
}
