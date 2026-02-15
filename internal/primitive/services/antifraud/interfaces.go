package antifraud

import "unified-workflow/internal/primitive/services/antifraud/models"

// AntifraudService interface defines antifraud operations
// Note: No explicit ctx/wfdata parameters - these will be handled by wrappers
type AntifraudService interface {
	// StoreTransaction stores a transaction in the antifraud system
	StoreTransaction(afTransaction models.AF_Transaction) error

	// ValidateTransactionByAML validates a transaction using the AML service
	ValidateTransactionByAML(afTransaction models.AF_Transaction) (models.ServiceResolution, error)

	// ValidateTransactionByFC validates a transaction using the FC service
	ValidateTransactionByFC(afTransaction models.AF_Transaction) (models.ServiceResolution, error)

	// ValidateTransactionByML validates a transaction using the ML service
	ValidateTransactionByML(afTransaction models.AF_Transaction) (models.ServiceResolution, error)

	// StoreServiceResolution stores the resolution from a service check (AML, FC, LST)
	StoreServiceResolution(resolution models.ServiceResolution) error

	// AddTransactionServiceCheck adds a completed service check resolution to the transaction aggregation process
	AddTransactionServiceCheck(resolution models.ServiceResolution) error

	// FinalizeTransaction finalizes the transaction validation process and retrieves the final resolution
	FinalizeTransaction(afTransaction models.AF_Transaction) (models.FinalResolution, error)

	// StoreFinalResolution stores the final resolution of the transaction
	StoreFinalResolution(resolution models.FinalResolution) error

	// HealthCheck checks the health of the antifraud service
	HealthCheck() (bool, error)

	// GetConfig returns the current configuration
	GetConfig() models.ClientConfig
}
