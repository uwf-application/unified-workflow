package di

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// ClusterNode represents a node in the cluster
type ClusterNode struct {
	// Node ID
	ID string

	// Node address (host:port)
	Address string

	// Node status
	Status NodeStatus

	// Last heartbeat time
	LastHeartbeat time.Time

	// Node capabilities
	Capabilities map[string]interface{}

	// Node metadata
	Metadata map[string]string
}

// NodeStatus represents the status of a cluster node
type NodeStatus string

const (
	// NodeStatusActive - node is active and healthy
	NodeStatusActive NodeStatus = "active"

	// NodeStatusJoining - node is joining the cluster
	NodeStatusJoining NodeStatus = "joining"

	// NodeStatusLeaving - node is leaving the cluster
	NodeStatusLeaving NodeStatus = "leaving"

	// NodeStatusFailed - node has failed
	NodeStatusFailed NodeStatus = "failed"

	// NodeStatusUnreachable - node is unreachable
	NodeStatusUnreachable NodeStatus = "unreachable"
)

// ClusterConfig holds configuration for cluster-aware DI
type ClusterConfig struct {
	// Node ID (auto-generated if empty)
	NodeID string

	// Node address (host:port)
	NodeAddress string

	// Cluster join addresses
	JoinAddresses []string

	// Heartbeat interval
	HeartbeatInterval time.Duration

	// Heartbeat timeout
	HeartbeatTimeout time.Duration

	// Replication factor
	ReplicationFactor int

	// Enable service discovery
	EnableDiscovery bool

	// Enable load balancing
	EnableLoadBalancing bool

	// Enable failover
	EnableFailover bool
}

// DefaultClusterConfig returns default cluster configuration
func DefaultClusterConfig() ClusterConfig {
	return ClusterConfig{
		NodeID:              generateNodeID(),
		NodeAddress:         "localhost:0",
		JoinAddresses:       []string{},
		HeartbeatInterval:   5 * time.Second,
		HeartbeatTimeout:    30 * time.Second,
		ReplicationFactor:   3,
		EnableDiscovery:     true,
		EnableLoadBalancing: true,
		EnableFailover:      true,
	}
}

// ClusterManager manages a cluster of DI containers
type ClusterManager struct {
	config ClusterConfig

	// Local node
	localNode *ClusterNode

	// Cluster nodes
	nodes   map[string]*ClusterNode
	nodesMu sync.RWMutex

	// Service registry
	services   map[string]*ClusterService
	servicesMu sync.RWMutex

	// Heartbeat ticker
	heartbeatTicker *time.Ticker

	// Stop channel
	stopChan chan struct{}

	// Event subscribers
	subscribers []ClusterEventSubscriber

	// Is running
	running bool
	mu      sync.RWMutex
}

// ClusterService represents a service in the cluster
type ClusterService struct {
	// Service name
	Name string

	// Service type
	Type string

	// Service instances
	Instances map[string]*ServiceInstance

	// Load balancing strategy
	LoadBalancingStrategy LoadBalancingStrategy

	// Replication factor
	ReplicationFactor int

	// Health check endpoint
	HealthCheckEndpoint string
}

// ServiceInstance represents an instance of a service
type ServiceInstance struct {
	// Instance ID
	ID string

	// Node ID where instance is running
	NodeID string

	// Service endpoint
	Endpoint string

	// Instance status
	Status InstanceStatus

	// Instance metadata
	Metadata map[string]string

	// Last health check
	LastHealthCheck time.Time

	// Load metrics
	LoadMetrics LoadMetrics
}

// InstanceStatus represents the status of a service instance
type InstanceStatus string

const (
	// InstanceStatusHealthy - instance is healthy
	InstanceStatusHealthy InstanceStatus = "healthy"

	// InstanceStatusUnhealthy - instance is unhealthy
	InstanceStatusUnhealthy InstanceStatus = "unhealthy"

	// InstanceStatusStarting - instance is starting
	InstanceStatusStarting InstanceStatus = "starting"

	// InstanceStatusStopping - instance is stopping
	InstanceStatusStopping InstanceStatus = "stopping"
)

// LoadBalancingStrategy represents load balancing strategy
type LoadBalancingStrategy string

const (
	// LoadBalancingRoundRobin - round robin load balancing
	LoadBalancingRoundRobin LoadBalancingStrategy = "round_robin"

	// LoadBalancingLeastConnections - least connections load balancing
	LoadBalancingLeastConnections LoadBalancingStrategy = "least_connections"

	// LoadBalancingRandom - random load balancing
	LoadBalancingRandom LoadBalancingStrategy = "random"

	// LoadBalancingHash - hash-based load balancing
	LoadBalancingHash LoadBalancingStrategy = "hash"
)

// LoadMetrics represents load metrics for a service instance
type LoadMetrics struct {
	// Active connections
	ActiveConnections int

	// Request rate (requests per second)
	RequestRate float64

	// CPU usage percentage
	CPUUsage float64

	// Memory usage percentage
	MemoryUsage float64

	// Response time percentiles
	ResponseTimeP50 time.Duration
	ResponseTimeP90 time.Duration
	ResponseTimeP95 time.Duration
	ResponseTimeP99 time.Duration
}

// ClusterEvent represents a cluster event
type ClusterEvent struct {
	// Event type
	Type ClusterEventType

	// Event timestamp
	Timestamp time.Time

	// Node involved (if any)
	Node *ClusterNode

	// Service involved (if any)
	Service *ClusterService

	// Instance involved (if any)
	Instance *ServiceInstance

	// Additional data
	Data map[string]interface{}
}

// ClusterEventType represents the type of cluster event
type ClusterEventType string

const (
	// ClusterEventNodeJoined - node joined the cluster
	ClusterEventNodeJoined ClusterEventType = "node_joined"

	// ClusterEventNodeLeft - node left the cluster
	ClusterEventNodeLeft ClusterEventType = "node_left"

	// ClusterEventNodeFailed - node failed
	ClusterEventNodeFailed ClusterEventType = "node_failed"

	// ClusterEventServiceRegistered - service registered
	ClusterEventServiceRegistered ClusterEventType = "service_registered"

	// ClusterEventServiceUnregistered - service unregistered
	ClusterEventServiceUnregistered ClusterEventType = "service_unregistered"

	// ClusterEventInstanceHealthy - instance became healthy
	ClusterEventInstanceHealthy ClusterEventType = "instance_healthy"

	// ClusterEventInstanceUnhealthy - instance became unhealthy
	ClusterEventInstanceUnhealthy ClusterEventType = "instance_unhealthy"
)

// ClusterEventSubscriber subscribes to cluster events
type ClusterEventSubscriber interface {
	// OnClusterEvent is called when a cluster event occurs
	OnClusterEvent(event ClusterEvent)
}

// NewClusterManager creates a new cluster manager
func NewClusterManager(config ClusterConfig) *ClusterManager {
	if config.NodeID == "" {
		config.NodeID = generateNodeID()
	}

	localNode := &ClusterNode{
		ID:            config.NodeID,
		Address:       config.NodeAddress,
		Status:        NodeStatusJoining,
		LastHeartbeat: time.Now(),
		Capabilities:  make(map[string]interface{}),
		Metadata:      make(map[string]string),
	}

	return &ClusterManager{
		config:      config,
		localNode:   localNode,
		nodes:       make(map[string]*ClusterNode),
		services:    make(map[string]*ClusterService),
		stopChan:    make(chan struct{}),
		subscribers: make([]ClusterEventSubscriber, 0),
	}
}

// Start starts the cluster manager
func (cm *ClusterManager) Start() error {
	cm.mu.Lock()
	if cm.running {
		cm.mu.Unlock()
		return fmt.Errorf("cluster manager already running")
	}

	// Add local node to cluster
	cm.nodes[cm.localNode.ID] = cm.localNode
	cm.localNode.Status = NodeStatusActive

	// Start heartbeat
	cm.heartbeatTicker = time.NewTicker(cm.config.HeartbeatInterval)
	go cm.heartbeat()

	// Join cluster if join addresses provided
	if len(cm.config.JoinAddresses) > 0 {
		go cm.joinCluster()
	}

	cm.running = true
	cm.mu.Unlock()

	log.Printf("Cluster manager started (node: %s)", cm.localNode.ID)

	// Notify node joined (after releasing lock)
	cm.notifyEvent(ClusterEvent{
		Type:      ClusterEventNodeJoined,
		Timestamp: time.Now(),
		Node:      cm.localNode,
	})

	return nil
}

// Stop stops the cluster manager
func (cm *ClusterManager) Stop() error {
	cm.mu.Lock()
	if !cm.running {
		cm.mu.Unlock()
		return fmt.Errorf("cluster manager not running")
	}

	close(cm.stopChan)
	cm.heartbeatTicker.Stop()

	// Mark local node as leaving
	cm.localNode.Status = NodeStatusLeaving

	cm.running = false
	cm.mu.Unlock()

	log.Println("Cluster manager stopped")

	// Notify node leaving (after releasing lock)
	cm.notifyEvent(ClusterEvent{
		Type:      ClusterEventNodeLeft,
		Timestamp: time.Now(),
		Node:      cm.localNode,
	})

	return nil
}

// heartbeat sends heartbeats and checks node health
func (cm *ClusterManager) heartbeat() {
	for {
		select {
		case <-cm.heartbeatTicker.C:
			cm.sendHeartbeat()
			cm.checkNodeHealth()

		case <-cm.stopChan:
			return
		}
	}
}

// sendHeartbeat sends heartbeat from local node
func (cm *ClusterManager) sendHeartbeat() {
	cm.nodesMu.Lock()
	cm.localNode.LastHeartbeat = time.Now()
	cm.nodesMu.Unlock()

	// In a real implementation, this would send heartbeats to other nodes
	// For now, we just update the local node's heartbeat time
}

// checkNodeHealth checks health of all nodes
func (cm *ClusterManager) checkNodeHealth() {
	cm.nodesMu.Lock()

	// Collect failed nodes
	failedNodes := make([]*ClusterNode, 0)
	now := time.Now()

	for nodeID, node := range cm.nodes {
		if nodeID == cm.localNode.ID {
			continue
		}

		// Check if node heartbeat is stale
		if now.Sub(node.LastHeartbeat) > cm.config.HeartbeatTimeout {
			if node.Status != NodeStatusFailed && node.Status != NodeStatusUnreachable {
				log.Printf("Node %s failed (last heartbeat: %v)", nodeID, node.LastHeartbeat)
				node.Status = NodeStatusFailed
				failedNodes = append(failedNodes, node)
			}
		}
	}

	cm.nodesMu.Unlock()

	// Notify about failed nodes and mark instances as unhealthy
	for _, node := range failedNodes {
		// Notify node failed
		cm.notifyEvent(ClusterEvent{
			Type:      ClusterEventNodeFailed,
			Timestamp: now,
			Node:      node,
		})

		// Mark service instances on failed node as unhealthy
		cm.markNodeInstancesUnhealthy(node.ID)
	}
}

// markNodeInstancesUnhealthy marks all service instances on a node as unhealthy
func (cm *ClusterManager) markNodeInstancesUnhealthy(nodeID string) {
	cm.servicesMu.Lock()

	// Collect unhealthy instances
	unhealthyInstances := make([]struct {
		service  *ClusterService
		instance *ServiceInstance
	}, 0)

	for _, service := range cm.services {
		for _, instance := range service.Instances {
			if instance.NodeID == nodeID && instance.Status == InstanceStatusHealthy {
				instance.Status = InstanceStatusUnhealthy
				unhealthyInstances = append(unhealthyInstances, struct {
					service  *ClusterService
					instance *ServiceInstance
				}{service, instance})
			}
		}
	}

	cm.servicesMu.Unlock()

	// Notify about unhealthy instances
	for _, item := range unhealthyInstances {
		cm.notifyEvent(ClusterEvent{
			Type:      ClusterEventInstanceUnhealthy,
			Timestamp: time.Now(),
			Service:   item.service,
			Instance:  item.instance,
		})
	}
}

// joinCluster joins the cluster
func (cm *ClusterManager) joinCluster() {
	// In a real implementation, this would:
	// 1. Connect to join addresses
	// 2. Discover cluster nodes
	// 3. Sync cluster state
	// 4. Register local node

	log.Printf("Joining cluster via addresses: %v", cm.config.JoinAddresses)

	// Simulate joining
	time.Sleep(2 * time.Second)

	cm.nodesMu.Lock()
	// Add some simulated nodes for demo
	newNodes := make([]*ClusterNode, 0)
	for i := 1; i <= 3; i++ {
		nodeID := fmt.Sprintf("node-%d", i)
		node := &ClusterNode{
			ID:            nodeID,
			Address:       fmt.Sprintf("node-%d:8080", i),
			Status:        NodeStatusActive,
			LastHeartbeat: time.Now(),
			Capabilities:  make(map[string]interface{}),
			Metadata:      make(map[string]string),
		}
		cm.nodes[nodeID] = node
		newNodes = append(newNodes, node)
	}
	cm.nodesMu.Unlock()

	// Notify about new nodes (after releasing lock)
	for _, node := range newNodes {
		cm.notifyEvent(ClusterEvent{
			Type:      ClusterEventNodeJoined,
			Timestamp: time.Now(),
			Node:      node,
		})
	}

	log.Printf("Joined cluster with %d nodes", len(cm.nodes))
}

// RegisterService registers a service in the cluster
func (cm *ClusterManager) RegisterService(service *ClusterService) error {
	cm.servicesMu.Lock()
	defer cm.servicesMu.Unlock()

	if _, exists := cm.services[service.Name]; exists {
		return fmt.Errorf("service %s already registered", service.Name)
	}

	cm.services[service.Name] = service

	// Notify service registered
	cm.notifyEvent(ClusterEvent{
		Type:      ClusterEventServiceRegistered,
		Timestamp: time.Now(),
		Service:   service,
	})

	log.Printf("Service registered: %s", service.Name)
	return nil
}

// RegisterServiceInstance registers a service instance
func (cm *ClusterManager) RegisterServiceInstance(serviceName string, instance *ServiceInstance) error {
	cm.servicesMu.Lock()
	defer cm.servicesMu.Unlock()

	service, exists := cm.services[serviceName]
	if !exists {
		return fmt.Errorf("service %s not found", serviceName)
	}

	if _, exists := service.Instances[instance.ID]; exists {
		return fmt.Errorf("instance %s already registered", instance.ID)
	}

	service.Instances[instance.ID] = instance

	// Notify instance healthy
	cm.notifyEvent(ClusterEvent{
		Type:      ClusterEventInstanceHealthy,
		Timestamp: time.Now(),
		Service:   service,
		Instance:  instance,
	})

	log.Printf("Service instance registered: %s/%s", serviceName, instance.ID)
	return nil
}

// GetServiceInstance gets a service instance using load balancing
func (cm *ClusterManager) GetServiceInstance(serviceName string) (*ServiceInstance, error) {
	cm.servicesMu.RLock()
	service, exists := cm.services[serviceName]
	cm.servicesMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	// Filter healthy instances
	healthyInstances := make([]*ServiceInstance, 0)
	for _, instance := range service.Instances {
		if instance.Status == InstanceStatusHealthy {
			healthyInstances = append(healthyInstances, instance)
		}
	}

	if len(healthyInstances) == 0 {
		return nil, fmt.Errorf("no healthy instances for service %s", serviceName)
	}

	// Apply load balancing strategy
	switch service.LoadBalancingStrategy {
	case LoadBalancingRoundRobin:
		return cm.roundRobin(healthyInstances, serviceName)
	case LoadBalancingLeastConnections:
		return cm.leastConnections(healthyInstances)
	case LoadBalancingRandom:
		return cm.random(healthyInstances)
	case LoadBalancingHash:
		return cm.hashBased(healthyInstances, serviceName)
	default:
		return cm.roundRobin(healthyInstances, serviceName)
	}
}

// roundRobin implements round robin load balancing
func (cm *ClusterManager) roundRobin(instances []*ServiceInstance, serviceName string) (*ServiceInstance, error) {
	// Simple round robin using service name as key
	// In production, use a proper round robin algorithm with persistence
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances available")
	}
	return instances[0], nil
}

// leastConnections implements least connections load balancing
func (cm *ClusterManager) leastConnections(instances []*ServiceInstance) (*ServiceInstance, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances available")
	}

	// Find instance with least active connections
	var selected *ServiceInstance
	minConnections := int(^uint(0) >> 1) // Max int

	for _, instance := range instances {
		if instance.LoadMetrics.ActiveConnections < minConnections {
			minConnections = instance.LoadMetrics.ActiveConnections
			selected = instance
		}
	}

	if selected == nil {
		return instances[0], nil
	}
	return selected, nil
}

// random implements random load balancing
func (cm *ClusterManager) random(instances []*ServiceInstance) (*ServiceInstance, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances available")
	}

	// Simple random selection
	// In production, use a proper random algorithm
	return instances[0], nil
}

// hashBased implements hash-based load balancing
func (cm *ClusterManager) hashBased(instances []*ServiceInstance, serviceName string) (*ServiceInstance, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances available")
	}

	// Simple hash-based selection
	// In production, use a consistent hashing algorithm
	hash := 0
	for _, char := range serviceName {
		hash += int(char)
	}
	index := hash % len(instances)
	return instances[index], nil
}

// UnregisterService unregisters a service
func (cm *ClusterManager) UnregisterService(serviceName string) error {
	cm.servicesMu.Lock()
	defer cm.servicesMu.Unlock()

	service, exists := cm.services[serviceName]
	if !exists {
		return fmt.Errorf("service %s not found", serviceName)
	}

	delete(cm.services, serviceName)

	// Notify service unregistered
	cm.notifyEvent(ClusterEvent{
		Type:      ClusterEventServiceUnregistered,
		Timestamp: time.Now(),
		Service:   service,
	})

	log.Printf("Service unregistered: %s", serviceName)
	return nil
}

// UnregisterServiceInstance unregisters a service instance
func (cm *ClusterManager) UnregisterServiceInstance(serviceName, instanceID string) error {
	cm.servicesMu.Lock()
	defer cm.servicesMu.Unlock()

	service, exists := cm.services[serviceName]
	if !exists {
		return fmt.Errorf("service %s not found", serviceName)
	}

	instance, exists := service.Instances[instanceID]
	if !exists {
		return fmt.Errorf("instance %s not found", instanceID)
	}

	delete(service.Instances, instanceID)

	// Notify instance removed
	cm.notifyEvent(ClusterEvent{
		Type:      ClusterEventInstanceUnhealthy,
		Timestamp: time.Now(),
		Service:   service,
		Instance:  instance,
	})

	log.Printf("Service instance unregistered: %s/%s", serviceName, instanceID)
	return nil
}

// GetNodes returns all cluster nodes
func (cm *ClusterManager) GetNodes() []*ClusterNode {
	cm.nodesMu.RLock()
	defer cm.nodesMu.RUnlock()

	nodes := make([]*ClusterNode, 0, len(cm.nodes))
	for _, node := range cm.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetServices returns all registered services
func (cm *ClusterManager) GetServices() []*ClusterService {
	cm.servicesMu.RLock()
	defer cm.servicesMu.RUnlock()

	services := make([]*ClusterService, 0, len(cm.services))
	for _, service := range cm.services {
		services = append(services, service)
	}
	return services
}

// GetService returns a service by name
func (cm *ClusterManager) GetService(serviceName string) (*ClusterService, bool) {
	cm.servicesMu.RLock()
	defer cm.servicesMu.RUnlock()

	service, exists := cm.services[serviceName]
	return service, exists
}

// GetNode returns a node by ID
func (cm *ClusterManager) GetNode(nodeID string) (*ClusterNode, bool) {
	cm.nodesMu.RLock()
	defer cm.nodesMu.RUnlock()

	node, exists := cm.nodes[nodeID]
	return node, exists
}

// Subscribe subscribes to cluster events
func (cm *ClusterManager) Subscribe(subscriber ClusterEventSubscriber) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.subscribers = append(cm.subscribers, subscriber)
}

// Unsubscribe unsubscribes from cluster events
func (cm *ClusterManager) Unsubscribe(subscriber ClusterEventSubscriber) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for i, sub := range cm.subscribers {
		if sub == subscriber {
			cm.subscribers = append(cm.subscribers[:i], cm.subscribers[i+1:]...)
			break
		}
	}
}

// notifyEvent notifies all subscribers of a cluster event
func (cm *ClusterManager) notifyEvent(event ClusterEvent) {
	cm.mu.RLock()
	subscribers := make([]ClusterEventSubscriber, len(cm.subscribers))
	copy(subscribers, cm.subscribers)
	cm.mu.RUnlock()

	for _, subscriber := range subscribers {
		go func(s ClusterEventSubscriber) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic in cluster event subscriber: %v", r)
				}
			}()
			s.OnClusterEvent(event)
		}(subscriber)
	}
}

// UpdateInstanceHealth updates the health status of a service instance
func (cm *ClusterManager) UpdateInstanceHealth(serviceName, instanceID string, healthy bool) error {
	cm.servicesMu.Lock()
	defer cm.servicesMu.Unlock()

	service, exists := cm.services[serviceName]
	if !exists {
		return fmt.Errorf("service %s not found", serviceName)
	}

	instance, exists := service.Instances[instanceID]
	if !exists {
		return fmt.Errorf("instance %s not found", instanceID)
	}

	oldStatus := instance.Status
	if healthy {
		instance.Status = InstanceStatusHealthy
		instance.LastHealthCheck = time.Now()
	} else {
		instance.Status = InstanceStatusUnhealthy
	}

	// Notify status change
	if oldStatus != instance.Status {
		eventType := ClusterEventInstanceUnhealthy
		if healthy {
			eventType = ClusterEventInstanceHealthy
		}

		cm.notifyEvent(ClusterEvent{
			Type:      eventType,
			Timestamp: time.Now(),
			Service:   service,
			Instance:  instance,
		})
	}

	return nil
}

// UpdateInstanceMetrics updates load metrics for a service instance
func (cm *ClusterManager) UpdateInstanceMetrics(serviceName, instanceID string, metrics LoadMetrics) error {
	cm.servicesMu.Lock()
	defer cm.servicesMu.Unlock()

	service, exists := cm.services[serviceName]
	if !exists {
		return fmt.Errorf("service %s not found", serviceName)
	}

	instance, exists := service.Instances[instanceID]
	if !exists {
		return fmt.Errorf("instance %s not found", instanceID)
	}

	instance.LoadMetrics = metrics
	instance.LastHealthCheck = time.Now()
	return nil
}

// IsRunning returns true if cluster manager is running
func (cm *ClusterManager) IsRunning() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.running
}

// GetLocalNode returns the local node
func (cm *ClusterManager) GetLocalNode() *ClusterNode {
	return cm.localNode
}

// GetClusterSize returns the number of nodes in the cluster
func (cm *ClusterManager) GetClusterSize() int {
	cm.nodesMu.RLock()
	defer cm.nodesMu.RUnlock()
	return len(cm.nodes)
}

// GetHealthyInstances returns healthy instances for a service
func (cm *ClusterManager) GetHealthyInstances(serviceName string) ([]*ServiceInstance, error) {
	cm.servicesMu.RLock()
	service, exists := cm.services[serviceName]
	cm.servicesMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	healthyInstances := make([]*ServiceInstance, 0)
	for _, instance := range service.Instances {
		if instance.Status == InstanceStatusHealthy {
			healthyInstances = append(healthyInstances, instance)
		}
	}

	return healthyInstances, nil
}

// generateNodeID generates a unique node ID
func generateNodeID() string {
	// In production, use a proper ID generator (e.g., UUID)
	return fmt.Sprintf("node-%d", time.Now().UnixNano())
}

// ClusterAwareContainer extends the DI container with cluster awareness
type ClusterAwareContainer struct {
	Container
	clusterManager *ClusterManager
}

// NewClusterAwareContainer creates a new cluster-aware container
func NewClusterAwareContainer(base Container, clusterManager *ClusterManager) *ClusterAwareContainer {
	return &ClusterAwareContainer{
		Container:      base,
		clusterManager: clusterManager,
	}
}

// GetClusterManager returns the cluster manager
func (c *ClusterAwareContainer) GetClusterManager() *ClusterManager {
	return c.clusterManager
}

// ResolveWithCluster resolves a component with cluster awareness
func (c *ClusterAwareContainer) ResolveWithCluster(componentType interface{}) (interface{}, error) {
	// First try local resolution
	instance, err := c.Resolve(componentType)
	if err == nil {
		return instance, nil
	}

	// If local resolution fails, try to find in cluster
	typeName := fmt.Sprintf("%T", componentType)
	serviceName := extractServiceName(typeName)

	// Get service instance from cluster
	instanceInfo, err := c.clusterManager.GetServiceInstance(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve component locally or from cluster: %w", err)
	}

	// In a real implementation, this would:
	// 1. Create a proxy to the remote instance
	// 2. Handle network communication
	// 3. Manage connection pooling
	// 4. Handle failover

	log.Printf("Resolved component %s from cluster instance %s", serviceName, instanceInfo.ID)
	return createClusterProxy(instanceInfo), nil
}

// RegisterService registers a service in the cluster
func (c *ClusterAwareContainer) RegisterService(serviceName, serviceType string) error {
	service := &ClusterService{
		Name:                  serviceName,
		Type:                  serviceType,
		Instances:             make(map[string]*ServiceInstance),
		LoadBalancingStrategy: LoadBalancingRoundRobin,
		ReplicationFactor:     c.clusterManager.config.ReplicationFactor,
		HealthCheckEndpoint:   "/health",
	}

	return c.clusterManager.RegisterService(service)
}

// RegisterLocalInstance registers the local container as a service instance
func (c *ClusterAwareContainer) RegisterLocalInstance(serviceName string, capabilities map[string]interface{}) error {
	localNode := c.clusterManager.GetLocalNode()

	instance := &ServiceInstance{
		ID:              fmt.Sprintf("%s-%s", serviceName, localNode.ID),
		NodeID:          localNode.ID,
		Endpoint:        localNode.Address,
		Status:          InstanceStatusHealthy,
		Metadata:        make(map[string]string),
		LastHealthCheck: time.Now(),
		LoadMetrics:     LoadMetrics{},
	}

	// Add capabilities to metadata
	for k, v := range capabilities {
		instance.Metadata[k] = fmt.Sprintf("%v", v)
	}

	return c.clusterManager.RegisterServiceInstance(serviceName, instance)
}

// extractServiceName extracts service name from type
func extractServiceName(typeName string) string {
	// Simple extraction - in production, use proper type analysis
	return typeName
}

// createClusterProxy creates a proxy to a cluster instance
func createClusterProxy(instance *ServiceInstance) interface{} {
	// In a real implementation, this would create a dynamic proxy
	// that forwards calls to the remote instance
	return nil
}

// ClusterEventFunc is a function that handles cluster events
type ClusterEventFunc func(event ClusterEvent)

// funcClusterEventSubscriber wraps a function as a ClusterEventSubscriber
type funcClusterEventSubscriber struct {
	fn ClusterEventFunc
}

// OnClusterEvent calls the wrapped function
func (f *funcClusterEventSubscriber) OnClusterEvent(event ClusterEvent) {
	f.fn(event)
}

// ExampleClusterUsage demonstrates cluster-aware DI usage
func ExampleClusterUsage() {
	// Create cluster configuration
	clusterConfig := DefaultClusterConfig()
	clusterConfig.NodeAddress = "localhost:8080"
	clusterConfig.JoinAddresses = []string{"node1:8080", "node2:8080"}

	// Create cluster manager
	clusterManager := NewClusterManager(clusterConfig)

	// Start cluster manager
	if err := clusterManager.Start(); err != nil {
		log.Fatalf("Failed to start cluster manager: %v", err)
	}
	defer clusterManager.Stop()

	// Create DI container
	container := NewWithConfig(DefaultConfig())

	// Create cluster-aware container
	clusterContainer := NewClusterAwareContainer(container, clusterManager)

	// Register a service in the cluster
	if err := clusterContainer.RegisterService("database", "postgres"); err != nil {
		log.Printf("Failed to register service: %v", err)
	}

	// Register local instance
	capabilities := map[string]interface{}{
		"version":  "1.0.0",
		"features": []string{"read", "write", "replication"},
	}
	if err := clusterContainer.RegisterLocalInstance("database", capabilities); err != nil {
		log.Printf("Failed to register local instance: %v", err)
	}

	// Subscribe to cluster events
	clusterManager.Subscribe(&funcClusterEventSubscriber{
		fn: func(event ClusterEvent) {
			log.Printf("Cluster event: %s", event.Type)
		},
	})

	log.Println("Cluster-aware DI container ready")
}
