package fraud

// CheckSuspiciousTransactionInput represents input for checking suspicious transactions
type CheckSuspiciousTransactionInput struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	CustomerID    string  `json:"customer_id,omitempty"`
	MerchantID    string  `json:"merchant_id,omitempty"`
	Location      string  `json:"location,omitempty"`
	IPAddress     string  `json:"ip_address,omitempty"`
}

// CheckSuspiciousTransactionOutput represents output from checking suspicious transactions
type CheckSuspiciousTransactionOutput struct {
	IsSuspicious bool    `json:"is_suspicious"`
	RiskScore    int     `json:"risk_score"`
	Reason       string  `json:"reason,omitempty"`
	Action       string  `json:"action"`     // "allow", "review", "block"
	Confidence   float64 `json:"confidence"` // 0.0 to 1.0
}

// FlagTransactionInput represents input for flagging a transaction
type FlagTransactionInput struct {
	TransactionID string `json:"transaction_id"`
	Reason        string `json:"reason"`
	Severity      string `json:"severity"` // "low", "medium", "high", "critical"
	Notes         string `json:"notes,omitempty"`
}

// FlagTransactionOutput represents output from flagging a transaction
type FlagTransactionOutput struct {
	FlagID         string `json:"flag_id"`
	Success        bool   `json:"success"`
	Message        string `json:"message,omitempty"`
	FlaggedAt      string `json:"flagged_at"`
	ReviewRequired bool   `json:"review_required"`
}

// GetRiskScoreInput represents input for getting risk score
type GetRiskScoreInput struct {
	CustomerID string `json:"customer_id"`
	TimeRange  string `json:"time_range,omitempty"` // "24h", "7d", "30d", "all"
}

// GetRiskScoreOutput represents output from getting risk score
type GetRiskScoreOutput struct {
	CustomerID  string       `json:"customer_id"`
	RiskScore   int          `json:"risk_score"`
	RiskLevel   string       `json:"risk_level"` // "low", "medium", "high", "critical"
	Factors     []RiskFactor `json:"factors"`
	LastUpdated string       `json:"last_updated"`
}

// RiskFactor represents a factor contributing to risk score
type RiskFactor struct {
	Factor     string  `json:"factor"`
	Score      int     `json:"score"`
	Confidence float64 `json:"confidence"`
	Details    string  `json:"details,omitempty"`
}

// FraudService interface defines fraud detection operations
// Note: No explicit ctx/wfdata parameters - these will be handled by wrappers
type FraudService interface {
	// CheckSuspiciousTransaction checks if a transaction is suspicious
	CheckSuspiciousTransaction(input CheckSuspiciousTransactionInput) (CheckSuspiciousTransactionOutput, error)

	// FlagTransaction flags a transaction for review
	FlagTransaction(input FlagTransactionInput) (FlagTransactionOutput, error)

	// GetRiskScore gets the risk score for a customer
	GetRiskScore(input GetRiskScoreInput) (GetRiskScoreOutput, error)

	// BatchCheck checks multiple transactions at once
	BatchCheck(inputs []CheckSuspiciousTransactionInput) ([]CheckSuspiciousTransactionOutput, error)
}
