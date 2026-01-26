package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"unified-workflow/internal/di"
)

// DatabaseService represents a database service that can be distributed
type DatabaseService interface {
	Query(ctx context.Context, query string) ([]map[string]interface{}, error)
	Insert(ctx context.Context, data map[string]interface{}) error
}

// PaymentService represents a payment processing service
type PaymentService interface {
	ProcessPayment(ctx context.Context, amount float64, currency string) (string, error)
	RefundPayment(ctx context.Context, paymentID string) error
}

// DatabaseServiceImpl implements DatabaseService
type DatabaseServiceImpl struct {
	nodeID string
}

func NewDatabaseServiceImpl(nodeID string) *DatabaseServiceImpl {
	return &DatabaseServiceImpl{nodeID: nodeID}
}

func (d *DatabaseServiceImpl) Query(ctx context.Context, query string) ([]map[string]interface{}, error) {
	log.Printf("[%s] Executing query: %s", d.nodeID, query)
	// Simulate database query
	time.Sleep(10 * time.Millisecond)
	return []map[string]interface{}{
		{"id": 1, "name": "Alice", "balance": 100.50},
		{"id": 2, "name": "Bob", "balance": 200.75},
	}, nil
}

func (d *DatabaseServiceImpl) Insert(ctx context.Context, data map[string]interface{}) error {
	log.Printf("[%s] Inserting data: %v", d.nodeID, data)
	// Simulate database insert
	time.Sleep(20 * time.Millisecond)
	return nil
}

// PaymentServiceImpl implements PaymentService
type PaymentServiceImpl struct {
	nodeID string
}

func NewPaymentServiceImpl(nodeID string) *PaymentServiceImpl {
	return &PaymentServiceImpl{nodeID: nodeID}
}

func (p *PaymentServiceImpl) ProcessPayment(ctx context.Context, amount float64, currency string) (string, error) {
	log.Printf("[%s] Processing payment: %.2f %s", p.nodeID, amount, currency)
	// Simulate payment processing
	time.Sleep(50 * time.Millisecond)
	return fmt.Sprintf("pay_%d", time.Now().UnixNano()), nil
}

func (p *PaymentServiceImpl) RefundPayment(ctx context.Context, paymentID string) error {
	log.Printf("[%s] Refunding payment: %s", p.nodeID, paymentID)
	// Simulate refund processing
	time.Sleep(30 * time.Millisecond)
	return nil
}

// WorkflowService uses cluster-aware DI to resolve services
type WorkflowService struct {
	container *di.ClusterAwareContainer
}

func NewWorkflowService(container *di.ClusterAwareContainer) *WorkflowService {
	return &WorkflowService{container: container}
}

func (w *WorkflowService) ExecuteWorkflow(ctx context.Context) error {
	log.Println("Starting workflow execution...")

	// Get cluster manager from container
	clusterManager := w.container.GetClusterManager()

	// Get database service instance from cluster
	dbInstance, err := clusterManager.GetServiceInstance("database")
	if err != nil {
		return fmt.Errorf("failed to get database service instance: %w", err)
	}

	log.Printf("Using database service instance: %s on node %s", dbInstance.ID, dbInstance.NodeID)

	// In a real implementation, we would:
	// 1. Create a client connection to the instance endpoint
	// 2. Execute the query remotely
	// 3. Handle network errors and retries

	// For this demo, we'll simulate using the local instance
	localNode := clusterManager.GetLocalNode()
	if dbInstance.NodeID == localNode.ID {
		// Use local implementation
		dbService := NewDatabaseServiceImpl(localNode.ID)
		results, err := dbService.Query(ctx, "SELECT * FROM users")
		if err != nil {
			return fmt.Errorf("database query failed: %w", err)
		}
		log.Printf("Query results: %v", results)
	} else {
		// Simulate remote call
		log.Printf("Simulating remote database query to instance %s at %s", dbInstance.ID, dbInstance.Endpoint)
		time.Sleep(50 * time.Millisecond)
		log.Printf("Remote query completed successfully")
	}

	// Get payment service instance from cluster
	paymentInstance, err := clusterManager.GetServiceInstance("payment")
	if err != nil {
		return fmt.Errorf("failed to get payment service instance: %w", err)
	}

	log.Printf("Using payment service instance: %s on node %s", paymentInstance.ID, paymentInstance.NodeID)

	// Simulate payment processing
	if paymentInstance.NodeID == localNode.ID {
		// Use local implementation
		paymentService := NewPaymentServiceImpl(localNode.ID)
		paymentID, err := paymentService.ProcessPayment(ctx, 100.0, "USD")
		if err != nil {
			return fmt.Errorf("payment processing failed: %w", err)
		}
		log.Printf("Payment processed successfully: %s", paymentID)
	} else {
		// Simulate remote call
		log.Printf("Simulating remote payment processing to instance %s at %s", paymentInstance.ID, paymentInstance.Endpoint)
		time.Sleep(80 * time.Millisecond)
		log.Printf("Remote payment processed successfully: pay_simulated_%d", time.Now().UnixNano())
	}

	log.Println("Workflow completed successfully")
	return nil
}

// ClusterNode represents a node in our example cluster
type ClusterNode struct {
	id        string
	address   string
	container di.Container
}

func NewClusterNode(id, address string) *ClusterNode {
	container := di.NewWithConfig(di.DefaultConfig())
	return &ClusterNode{
		id:        id,
		address:   address,
		container: container,
	}
}

func (n *ClusterNode) RegisterServices() {
	// Register local implementations
	n.container.Register(func() DatabaseService {
		return NewDatabaseServiceImpl(n.id)
	}, di.ProviderFunc(func(c di.Container) (interface{}, error) {
		return NewDatabaseServiceImpl(n.id), nil
	}), di.Singleton)

	n.container.Register(func() PaymentService {
		return NewPaymentServiceImpl(n.id)
	}, di.ProviderFunc(func(c di.Container) (interface{}, error) {
		return NewPaymentServiceImpl(n.id), nil
	}), di.Singleton)
}

func main() {
	log.Println("Starting cluster-aware DI example...")

	// Create cluster configuration
	clusterConfig := di.DefaultClusterConfig()
	clusterConfig.NodeAddress = "localhost:8080"
	clusterConfig.JoinAddresses = []string{"node1:8080", "node2:8080", "node3:8080"}
	clusterConfig.HeartbeatInterval = 2 * time.Second
	clusterConfig.HeartbeatTimeout = 10 * time.Second
	clusterConfig.ReplicationFactor = 2

	// Create cluster manager
	clusterManager := di.NewClusterManager(clusterConfig)

	// Start cluster manager
	if err := clusterManager.Start(); err != nil {
		log.Fatalf("Failed to start cluster manager: %v", err)
	}
	defer clusterManager.Stop()

	// Create local node
	localNode := NewClusterNode(clusterConfig.NodeID, clusterConfig.NodeAddress)
	localNode.RegisterServices()

	// Create cluster-aware container
	clusterContainer := di.NewClusterAwareContainer(localNode.container, clusterManager)

	// Register services in the cluster
	if err := clusterContainer.RegisterService("database", "postgres"); err != nil {
		log.Printf("Failed to register database service: %v", err)
	}

	if err := clusterContainer.RegisterService("payment", "stripe"); err != nil {
		log.Printf("Failed to register payment service: %v", err)
	}

	// Register local instances with capabilities
	dbCapabilities := map[string]interface{}{
		"version":     "1.2.0",
		"engine":      "postgresql",
		"pool_size":   100,
		"read_only":   false,
		"replication": true,
	}

	if err := clusterContainer.RegisterLocalInstance("database", dbCapabilities); err != nil {
		log.Printf("Failed to register database instance: %v", err)
	}

	paymentCapabilities := map[string]interface{}{
		"version":     "2.1.0",
		"provider":    "stripe",
		"currencies":  []string{"USD", "EUR", "GBP"},
		"max_amount":  10000.0,
		"refund_days": 30,
	}

	if err := clusterContainer.RegisterLocalInstance("payment", paymentCapabilities); err != nil {
		log.Printf("Failed to register payment instance: %v", err)
	}

	// Subscribe to cluster events using custom subscriber
	clusterManager.Subscribe(&clusterEventSubscriber{})

	// Create workflow service
	workflowService := NewWorkflowService(clusterContainer)

	// Simulate multiple workflow executions
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		log.Printf("\n=== Workflow Execution %d ===", i)
		if err := workflowService.ExecuteWorkflow(ctx); err != nil {
			log.Printf("Workflow failed: %v", err)
		}

		// Simulate cluster changes between executions
		if i == 2 {
			log.Println("\n🔧 Simulating node failure...")
			// In a real scenario, this would be detected by heartbeat timeout
			time.Sleep(3 * time.Second)
		}

		if i == 4 {
			log.Println("\n🔧 Simulating new node joining...")
			// In a real scenario, a new node would join the cluster
			time.Sleep(2 * time.Second)
		}

		time.Sleep(1 * time.Second)
	}

	// Demonstrate cluster information
	log.Println("\n=== Cluster Information ===")
	nodes := clusterManager.GetNodes()
	log.Printf("Total nodes in cluster: %d", len(nodes))
	for _, node := range nodes {
		log.Printf("  - Node %s: %s (Status: %s)", node.ID, node.Address, node.Status)
	}

	services := clusterManager.GetServices()
	log.Printf("\nTotal services registered: %d", len(services))
	for _, service := range services {
		log.Printf("  - Service %s (%s):", service.Name, service.Type)
		log.Printf("    Load Balancing: %s", service.LoadBalancingStrategy)
		log.Printf("    Replication Factor: %d", service.ReplicationFactor)
		log.Printf("    Instances: %d", len(service.Instances))
		for _, instance := range service.Instances {
			log.Printf("      * %s on node %s (%s)", instance.ID, instance.NodeID, instance.Status)
			log.Printf("        Endpoint: %s", instance.Endpoint)
			log.Printf("        Last check: %v", instance.LastHealthCheck)
		}
	}

	// Demonstrate load balancing
	log.Println("\n=== Load Balancing Demo ===")
	for i := 1; i <= 3; i++ {
		instance, err := clusterManager.GetServiceInstance("database")
		if err != nil {
			log.Printf("Failed to get database instance: %v", err)
		} else {
			log.Printf("Request %d: Load balancer selected instance %s", i, instance.ID)
		}
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("\n✅ Cluster-aware DI example completed successfully!")
}

// clusterEventSubscriber implements di.ClusterEventSubscriber
type clusterEventSubscriber struct{}

func (c *clusterEventSubscriber) OnClusterEvent(event di.ClusterEvent) {
	log.Printf("📡 Cluster Event: %s", event.Type)
	if event.Node != nil {
		log.Printf("   Node: %s (%s)", event.Node.ID, event.Node.Status)
	}
	if event.Service != nil {
		log.Printf("   Service: %s", event.Service.Name)
	}
	if event.Instance != nil {
		log.Printf("   Instance: %s (%s)", event.Instance.ID, event.Instance.Status)
	}
}

// Helper function to create FuncClusterEventSubscriber
type FuncClusterEventSubscriber struct {
	Fn func(event di.ClusterEvent)
}

func (f *FuncClusterEventSubscriber) OnClusterEvent(event di.ClusterEvent) {
	f.Fn(event)
}
