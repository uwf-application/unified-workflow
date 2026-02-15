package models

// AF_Transaction represents antifraud transaction (matching SDK)
type AF_Transaction struct {
	AF_Id       string      `json:"af_id"`
	AF_AddDate  string      `json:"af_add_date"`
	Transaction Transaction `json:"transaction"`
}

// Transaction represents the actual transaction data
type Transaction struct {
	Id                 string `json:"id"`
	Type               string `json:"type"`
	Date               string `json:"date"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	ClientId           string `json:"client_id"`
	ClientName         string `json:"client_name"`
	ClientPAN          string `json:"client_pan"`
	ClientCVV          string `json:"client_cvv"`
	ClientCardHolder   string `json:"client_card_holder"`
	ClientPhone        string `json:"client_phone"`
	MerchantTerminalId string `json:"merchant_terminal_id"`
	Channel            string `json:"channel"`
	LocationIp         string `json:"location_ip"`
}

// ServiceResolution represents service validation result
type ServiceResolution struct {
	ServiceName string `json:"service_name"`
	Resolution  string `json:"resolution"`
	Score       int    `json:"score"`
	Details     string `json:"details"`
}

// FinalResolution represents final transaction validation result
type FinalResolution struct {
	TransactionId string   `json:"transaction_id"`
	FinalDecision string   `json:"final_decision"`
	RiskScore     int      `json:"risk_score"`
	Reasons       []string `json:"reasons"`
}

// ClientConfig represents configuration for the antifraud client
type ClientConfig struct {
	APIKey string `json:"api_key"`
	Host   string `json:"host"`
	// Timeout in seconds
	Timeout int `json:"timeout"`
	// Enable/disable the service
	Enabled bool `json:"enabled"`
	// Max retries for failed requests
	MaxRetries int `json:"max_retries"`
	// Circuit breaker configuration
	CircuitBreakerEnabled   bool `json:"circuit_breaker_enabled"`
	CircuitBreakerThreshold int  `json:"circuit_breaker_threshold"`
	CircuitBreakerTimeout   int  `json:"circuit_breaker_timeout"`
}
