package main

import (
	"fmt"
	"os"
	"time"

	"unified-workflow/internal/primitive"

	"github.com/google/uuid"
)

func main() {
	// Example 1: Using global primitive
	fmt.Println("=== Example 1: Using Global Primitive ===")
	exampleGlobalPrimitive()

	// Example 2: Full transaction validation flow
	fmt.Println("\n=== Example 2: Full Transaction Validation Flow ===")
	exampleFullTransactionFlow()

	// Example 3: Error handling
	fmt.Println("\n=== Example 3: Error Handling ===")
	exampleErrorHandling()
}

func exampleGlobalPrimitive() {
	// Initialize primitive with antifraud configuration
	config := &primitive.Config{
		AntifraudAPIKey:                  os.Getenv("ANTIFRAUD_API_KEY"),
		AntifraudAPIHost:                 os.Getenv("ANTIFRAUD_API_HOST"),
		AntifraudTimeout:                 30,
		AntifraudEnabled:                 true,
		AntifraudMaxRetries:              3,
		AntifraudCircuitBreakerEnabled:   true,
		AntifraudCircuitBreakerThreshold: 5,
		AntifraudCircuitBreakerTimeout:   60,
	}

	// Initialize the global primitive
	err := primitive.Init(config)
	if err != nil {
		fmt.Printf("Failed to initialize primitive: %v\n", err)
		return
	}

	// Check if antifraud service is available
	if primitive.Default.Antifraud == nil {
		fmt.Println("Antifraud service is not available")
		return
	}

	// Create a sample transaction
	transaction := createSampleTransaction()

	// Store transaction
	fmt.Println("Storing transaction...")
	err = primitive.Default.Antifraud.StoreTransaction(transaction)
	if err != nil {
		fmt.Printf("Failed to store transaction: %v\n", err)
		return
	}
	fmt.Println("Transaction stored successfully")

	// Validate with AML service
	fmt.Println("\nValidating transaction with AML service...")
	amlResult, err := primitive.Default.Antifraud.ValidateTransactionByAML(transaction)
	if err != nil {
		fmt.Printf("AML validation failed: %v\n", err)
		return
	}
	fmt.Printf("AML Result: %v\n", amlResult)
}

func exampleFullTransactionFlow() {
	// For this example, we'll use environment variables
	// In production, you would load these from config
	apiKey := os.Getenv("ANTIFRAUD_API_KEY")
	host := os.Getenv("ANTIFRAUD_API_HOST")

	if apiKey == "" || host == "" {
		fmt.Println("Please set ANTIFRAUD_API_KEY and ANTIFRAUD_API_HOST environment variables")
		fmt.Println("Example:")
		fmt.Println("  export ANTIFRAUD_API_KEY='your-api-key'")
		fmt.Println("  export ANTIFRAUD_API_HOST='https://api.antifraud.example.com'")
		return
	}

	// Initialize primitive
	config := &primitive.Config{
		AntifraudAPIKey:     apiKey,
		AntifraudAPIHost:    host,
		AntifraudTimeout:    30,
		AntifraudEnabled:    true,
		AntifraudMaxRetries: 3,
	}

	err := primitive.Init(config)
	if err != nil {
		fmt.Printf("Failed to initialize primitive: %v\n", err)
		return
	}

	// Create transaction
	transaction := createSampleTransaction()

	fmt.Println("Starting full transaction validation flow...")

	// Step 1: Store transaction
	fmt.Println("\n1. Storing transaction...")
	err = primitive.Default.Antifraud.StoreTransaction(transaction)
	if err != nil {
		fmt.Printf("❌ Failed to store transaction: %v\n", err)
		return
	}
	fmt.Println("✅ Transaction stored")

	// Step 2: Validate with AML
	fmt.Println("\n2. Validating with AML service...")
	amlResult, err := primitive.Default.Antifraud.ValidateTransactionByAML(transaction)
	if err != nil {
		fmt.Printf("❌ AML validation failed: %v\n", err)
		return
	}
	fmt.Printf("✅ AML Result: %v\n", amlResult)

	// Step 3: Store AML resolution
	fmt.Println("\n3. Storing AML resolution...")
	err = primitive.Default.Antifraud.StoreServiceResolution(amlResult)
	if err != nil {
		fmt.Printf("❌ Failed to store AML resolution: %v\n", err)
		return
	}
	fmt.Println("✅ AML resolution stored")

	// Step 4: Add to transaction check
	fmt.Println("\n4. Adding AML check to transaction...")
	err = primitive.Default.Antifraud.AddTransactionServiceCheck(amlResult)
	if err != nil {
		fmt.Printf("❌ Failed to add AML check: %v\n", err)
		return
	}
	fmt.Println("✅ AML check added")

	// Step 5: Validate with FC
	fmt.Println("\n5. Validating with FC service...")
	fcResult, err := primitive.Default.Antifraud.ValidateTransactionByFC(transaction)
	if err != nil {
		fmt.Printf("❌ FC validation failed: %v\n", err)
		return
	}
	fmt.Printf("✅ FC Result: %v\n", fcResult)

	// Step 6: Store FC resolution
	fmt.Println("\n6. Storing FC resolution...")
	err = primitive.Default.Antifraud.StoreServiceResolution(fcResult)
	if err != nil {
		fmt.Printf("❌ Failed to store FC resolution: %v\n", err)
		return
	}
	fmt.Println("✅ FC resolution stored")

	// Step 7: Add FC check
	fmt.Println("\n7. Adding FC check to transaction...")
	err = primitive.Default.Antifraud.AddTransactionServiceCheck(fcResult)
	if err != nil {
		fmt.Printf("❌ Failed to add FC check: %v\n", err)
		return
	}
	fmt.Println("✅ FC check added")

	// Step 8: Finalize transaction
	fmt.Println("\n8. Finalizing transaction...")
	finalResult, err := primitive.Default.Antifraud.FinalizeTransaction(transaction)
	if err != nil {
		fmt.Printf("❌ Failed to finalize transaction: %v\n", err)
		return
	}
	fmt.Printf("✅ Final Result: %v\n", finalResult)

	// Step 9: Store final resolution
	fmt.Println("\n9. Storing final resolution...")
	err = primitive.Default.Antifraud.StoreFinalResolution(finalResult)
	if err != nil {
		fmt.Printf("❌ Failed to store final resolution: %v\n", err)
		return
	}
	fmt.Println("✅ Final resolution stored")

	fmt.Println("\n🎉 Transaction validation completed successfully!")
}

func exampleErrorHandling() {
	// Initialize with disabled antifraud service
	config := &primitive.Config{
		AntifraudEnabled: false,
	}

	err := primitive.Init(config)
	if err != nil {
		fmt.Printf("Failed to initialize primitive: %v\n", err)
		return
	}

	// Try to use disabled service
	transaction := createSampleTransaction()

	fmt.Println("Attempting to use disabled antifraud service...")
	err = primitive.Default.Antifraud.StoreTransaction(transaction)
	if err != nil {
		fmt.Printf("Expected error (service disabled): %v\n", err)
	} else {
		fmt.Println("Unexpected: Service should be disabled")
	}

	// Test health check
	fmt.Println("\nChecking service health...")
	healthy, err := primitive.Default.Antifraud.HealthCheck()
	if err != nil {
		fmt.Printf("Health check error: %v\n", err)
	} else {
		fmt.Printf("Service healthy: %v\n", healthy)
	}
}

func createSampleTransaction() interface{} {
	now := time.Now()

	// Create a transaction structure that matches what the antifraud service expects
	// In a real implementation, this would use the actual types from the antifraud package
	return map[string]interface{}{
		"af_id":       uuid.NewString(),
		"af_add_date": now.Format(time.RFC3339Nano),
		"transaction": map[string]interface{}{
			"id":                   uuid.NewString(),
			"type":                 "deposit",
			"date":                 now.Format(time.RFC3339Nano),
			"amount":               "100000",
			"currency":             "KZT",
			"client_id":            uuid.NewString(),
			"client_name":          "John Smith",
			"client_pan":           "111111******1111",
			"client_cvv":           "111",
			"client_card_holder":   "JOHN SMITH",
			"client_phone":         "+77007007070",
			"merchant_terminal_id": "00000001",
			"channel":              "E-com",
			"location_ip":          "192.168.0.1",
		},
	}
}

// Helper function to print configuration
func printConfig(config *primitive.Config) {
	fmt.Println("Antifraud Configuration:")
	fmt.Printf("  API Key: %s\n", maskAPIKey(config.AntifraudAPIKey))
	fmt.Printf("  Host: %s\n", config.AntifraudAPIHost)
	fmt.Printf("  Timeout: %d seconds\n", config.AntifraudTimeout)
	fmt.Printf("  Enabled: %v\n", config.AntifraudEnabled)
	fmt.Printf("  Max Retries: %d\n", config.AntifraudMaxRetries)
	fmt.Printf("  Circuit Breaker Enabled: %v\n", config.AntifraudCircuitBreakerEnabled)
	fmt.Printf("  Circuit Breaker Threshold: %d\n", config.AntifraudCircuitBreakerThreshold)
	fmt.Printf("  Circuit Breaker Timeout: %d seconds\n", config.AntifraudCircuitBreakerTimeout)
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}
