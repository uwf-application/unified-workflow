package payment

// ProcessPaymentInput represents input for processing a payment
type ProcessPaymentInput struct {
	TransactionID string                 `json:"transaction_id"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	CustomerID    string                 `json:"customer_id"`
	MerchantID    string                 `json:"merchant_id"`
	PaymentMethod string                 `json:"payment_method"` // "credit_card", "debit_card", "bank_transfer", "wallet"
	Description   string                 `json:"description,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessPaymentOutput represents output from processing a payment
type ProcessPaymentOutput struct {
	TransactionID     string                 `json:"transaction_id"`
	PaymentID         string                 `json:"payment_id"`
	Status            string                 `json:"status"` // "pending", "completed", "failed", "refunded"
	AuthorizationCode string                 `json:"authorization_code,omitempty"`
	ProcessedAt       string                 `json:"processed_at"`
	FailureReason     string                 `json:"failure_reason,omitempty"`
	GatewayResponse   map[string]interface{} `json:"gateway_response,omitempty"`
}

// RefundPaymentInput represents input for refunding a payment
type RefundPaymentInput struct {
	PaymentID     string  `json:"payment_id"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Reason        string  `json:"reason,omitempty"`
	PartialRefund bool    `json:"partial_refund"`
}

// RefundPaymentOutput represents output from refunding a payment
type RefundPaymentOutput struct {
	RefundID       string  `json:"refund_id"`
	PaymentID      string  `json:"payment_id"`
	Status         string  `json:"status"` // "pending", "completed", "failed"
	RefundedAmount float64 `json:"refunded_amount"`
	RefundedAt     string  `json:"refunded_at"`
	FailureReason  string  `json:"failure_reason,omitempty"`
}

// GetPaymentStatusInput represents input for getting payment status
type GetPaymentStatusInput struct {
	PaymentID     string `json:"payment_id"`
	TransactionID string `json:"transaction_id,omitempty"`
}

// GetPaymentStatusOutput represents output from getting payment status
type GetPaymentStatusOutput struct {
	PaymentID     string                 `json:"payment_id"`
	TransactionID string                 `json:"transaction_id"`
	Status        string                 `json:"status"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	Details       map[string]interface{} `json:"details,omitempty"`
}

// PaymentService interface defines payment processing operations
type PaymentService interface {
	// ProcessPayment processes a payment transaction
	ProcessPayment(input ProcessPaymentInput) (ProcessPaymentOutput, error)

	// RefundPayment refunds a payment
	RefundPayment(input RefundPaymentInput) (RefundPaymentOutput, error)

	// GetPaymentStatus gets the status of a payment
	GetPaymentStatus(input GetPaymentStatusInput) (GetPaymentStatusOutput, error)

	// CapturePayment captures an authorized payment
	CapturePayment(input CapturePaymentInput) (CapturePaymentOutput, error)

	// VoidPayment voids a payment before settlement
	VoidPayment(input VoidPaymentInput) (VoidPaymentOutput, error)
}

// CapturePaymentInput represents input for capturing a payment
type CapturePaymentInput struct {
	PaymentID     string  `json:"payment_id"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
}

// CapturePaymentOutput represents output from capturing a payment
type CapturePaymentOutput struct {
	CaptureID      string  `json:"capture_id"`
	PaymentID      string  `json:"payment_id"`
	Status         string  `json:"status"`
	CapturedAmount float64 `json:"captured_amount"`
	CapturedAt     string  `json:"captured_at"`
}

// VoidPaymentInput represents input for voiding a payment
type VoidPaymentInput struct {
	PaymentID     string `json:"payment_id"`
	TransactionID string `json:"transaction_id"`
	Reason        string `json:"reason,omitempty"`
}

// VoidPaymentOutput represents output from voiding a payment
type VoidPaymentOutput struct {
	VoidID    string `json:"void_id"`
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	VoidedAt  string `json:"voided_at"`
}
