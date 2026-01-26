package child_steps

import (
	"unified-workflow/internal/common/model"
)

// GetExampleChildSteps returns example child steps
func GetExampleChildSteps() []*model.ChildStep {
	return []*model.ChildStep{
		createValidationChildStep(),
		createTransformationChildStep(),
		createNotificationChildStep(),
		createLoggingChildStep(),
	}
}

// createValidationChildStep creates a validation child step
func createValidationChildStep() *model.ChildStep {
	return model.NewChildStep(
		"validate-input",
		func(context interface{}, data interface{}) interface{} {
			// Request hook: validate input data
			return map[string]interface{}{
				"validated": true,
				"timestamp": "2024-01-01T00:00:00Z",
			}
		},
		func(context interface{}, data interface{}) interface{} {
			// Response hook: process validation result
			return map[string]interface{}{
				"validation_passed": true,
				"message":           "Input validation successful",
			}
		},
		func(response interface{}) error {
			// Validate hook: check if validation passed
			// In a real implementation, this would check the response
			return nil
		},
	)
}

// createTransformationChildStep creates a data transformation child step
func createTransformationChildStep() *model.ChildStep {
	return model.NewChildStep(
		"transform-data",
		func(context interface{}, data interface{}) interface{} {
			// Request hook: prepare data for transformation
			return map[string]interface{}{
				"transformation_type": "json_to_xml",
				"source_format":       "json",
				"target_format":       "xml",
			}
		},
		func(context interface{}, data interface{}) interface{} {
			// Response hook: process transformed data
			return map[string]interface{}{
				"transformed": true,
				"data_size":   1024,
				"format":      "xml",
			}
		},
		func(response interface{}) error {
			// Validate hook: check if transformation was successful
			// In a real implementation, this would check the response
			return nil
		},
	)
}

// createNotificationChildStep creates a notification child step
func createNotificationChildStep() *model.ChildStep {
	return model.NewChildStep(
		"send-notification",
		func(context interface{}, data interface{}) interface{} {
			// Request hook: prepare notification
			return map[string]interface{}{
				"notification_type": "email",
				"recipient":         "user@example.com",
				"subject":           "Workflow Notification",
			}
		},
		func(context interface{}, data interface{}) interface{} {
			// Response hook: process notification result
			return map[string]interface{}{
				"sent":            true,
				"notification_id": "notif-12345",
				"timestamp":       "2024-01-01T00:00:00Z",
			}
		},
		func(response interface{}) error {
			// Validate hook: check if notification was sent
			// In a real implementation, this would check the response
			return nil
		},
	)
}

// createLoggingChildStep creates a logging child step
func createLoggingChildStep() *model.ChildStep {
	return model.NewChildStep(
		"log-execution",
		func(context interface{}, data interface{}) interface{} {
			// Request hook: prepare log data
			return map[string]interface{}{
				"log_level": "info",
				"operation": "workflow_execution",
				"timestamp": "2024-01-01T00:00:00Z",
			}
		},
		func(context interface{}, data interface{}) interface{} {
			// Response hook: process log result
			return map[string]interface{}{
				"logged":  true,
				"log_id":  "log-67890",
				"message": "Execution logged successfully",
			}
		},
		func(response interface{}) error {
			// Validate hook: check if logging was successful
			// In a real implementation, this would check the response
			return nil
		},
	)
}
