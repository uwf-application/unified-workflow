package di

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// MetricsCollector collects performance metrics for the DI container
type MetricsCollector interface {
	// Record resolution metrics
	RecordResolution(componentType string, duration time.Duration, success bool)

	// Record registration metrics
	RecordRegistration(componentType string, lifecycle Lifecycle)

	// Record scope creation metrics
	RecordScopeCreation(duration time.Duration)

	// Record pool operations
	RecordPoolGet(componentType string, duration time.Duration, success bool)
	RecordPoolPut(componentType string, duration time.Duration)

	// Get metrics snapshot
	GetMetrics() *ContainerMetrics

	// Reset metrics
	Reset()
}

// ContainerMetrics represents DI container performance metrics
type ContainerMetrics struct {
	// Timestamp of metrics collection
	Timestamp time.Time `json:"timestamp"`

	// Resolution metrics
	TotalResolutions      int64         `json:"total_resolutions"`
	SuccessfulResolutions int64         `json:"successful_resolutions"`
	FailedResolutions     int64         `json:"failed_resolutions"`
	AverageResolutionTime time.Duration `json:"average_resolution_time"`
	MaxResolutionTime     time.Duration `json:"max_resolution_time"`
	MinResolutionTime     time.Duration `json:"min_resolution_time"`

	// Registration metrics
	TotalRegistrations int64               `json:"total_registrations"`
	ByLifecycle        map[Lifecycle]int64 `json:"by_lifecycle"`

	// Scope metrics
	TotalScopesCreated int64         `json:"total_scopes_created"`
	AverageScopeTime   time.Duration `json:"average_scope_time"`
	MaxScopeTime       time.Duration `json:"max_scope_time"`

	// Pool metrics
	TotalPoolGets      int64         `json:"total_pool_gets"`
	SuccessfulPoolGets int64         `json:"successful_pool_gets"`
	FailedPoolGets     int64         `json:"failed_pool_gets"`
	TotalPoolPuts      int64         `json:"total_pool_puts"`
	AveragePoolGetTime time.Duration `json:"average_pool_get_time"`

	// Component-specific metrics
	ComponentMetrics map[string]*ComponentMetrics `json:"component_metrics"`

	// Performance indicators
	ResolutionRatePerSecond float64 `json:"resolution_rate_per_second"`
	ErrorRate               float64 `json:"error_rate"`
	PoolHitRate             float64 `json:"pool_hit_rate"`
}

// ComponentMetrics represents metrics for a specific component type
type ComponentMetrics struct {
	TotalResolutions      int64         `json:"total_resolutions"`
	SuccessfulResolutions int64         `json:"successful_resolutions"`
	FailedResolutions     int64         `json:"failed_resolutions"`
	TotalResolutionTime   time.Duration `json:"total_resolution_time"`
	AverageResolutionTime time.Duration `json:"average_resolution_time"`
	MaxResolutionTime     time.Duration `json:"max_resolution_time"`
	MinResolutionTime     time.Duration `json:"min_resolution_time"`

	TotalPoolGets      int64         `json:"total_pool_gets"`
	SuccessfulPoolGets int64         `json:"successful_pool_gets"`
	FailedPoolGets     int64         `json:"failed_pool_gets"`
	TotalPoolPuts      int64         `json:"total_pool_puts"`
	TotalPoolGetTime   time.Duration `json:"total_pool_get_time"`
	AveragePoolGetTime time.Duration `json:"average_pool_get_time"`
}

// defaultMetricsCollector is the default implementation of MetricsCollector
type defaultMetricsCollector struct {
	mu sync.RWMutex

	// Resolution metrics
	totalResolutions      atomic.Int64
	successfulResolutions atomic.Int64
	failedResolutions     atomic.Int64
	totalResolutionTime   atomic.Int64 // nanoseconds
	maxResolutionTime     atomic.Int64 // nanoseconds
	minResolutionTime     atomic.Int64 // nanoseconds

	// Registration metrics
	totalRegistrations atomic.Int64
	byLifecycle        map[Lifecycle]*atomic.Int64

	// Scope metrics
	totalScopesCreated atomic.Int64
	totalScopeTime     atomic.Int64 // nanoseconds
	maxScopeTime       atomic.Int64 // nanoseconds

	// Pool metrics
	totalPoolGets      atomic.Int64
	successfulPoolGets atomic.Int64
	failedPoolGets     atomic.Int64
	totalPoolPuts      atomic.Int64
	totalPoolGetTime   atomic.Int64 // nanoseconds

	// Component-specific metrics
	componentMetrics map[string]*componentMetrics

	// Start time for rate calculation
	startTime time.Time
}

type componentMetrics struct {
	mu sync.RWMutex

	totalResolutions      int64
	successfulResolutions int64
	failedResolutions     int64
	totalResolutionTime   int64 // nanoseconds
	maxResolutionTime     int64 // nanoseconds
	minResolutionTime     int64 // nanoseconds

	totalPoolGets      int64
	successfulPoolGets int64
	failedPoolGets     int64
	totalPoolPuts      int64
	totalPoolGetTime   int64 // nanoseconds
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() MetricsCollector {
	collector := &defaultMetricsCollector{
		byLifecycle:      make(map[Lifecycle]*atomic.Int64),
		componentMetrics: make(map[string]*componentMetrics),
		startTime:        time.Now(),
	}

	// Initialize lifecycle counters
	for _, lifecycle := range []Lifecycle{Singleton, Transient, Scoped, Pooled} {
		collector.byLifecycle[lifecycle] = &atomic.Int64{}
	}

	// Initialize min resolution time to a large value
	collector.minResolutionTime.Store(int64(time.Hour))

	return collector
}

// RecordResolution records a resolution operation
func (c *defaultMetricsCollector) RecordResolution(componentType string, duration time.Duration, success bool) {
	c.totalResolutions.Add(1)

	durationNs := int64(duration)
	c.totalResolutionTime.Add(durationNs)

	// Update max resolution time
	for {
		currentMax := c.maxResolutionTime.Load()
		if durationNs <= currentMax {
			break
		}
		if c.maxResolutionTime.CompareAndSwap(currentMax, durationNs) {
			break
		}
	}

	// Update min resolution time
	for {
		currentMin := c.minResolutionTime.Load()
		if durationNs >= currentMin {
			break
		}
		if c.minResolutionTime.CompareAndSwap(currentMin, durationNs) {
			break
		}
	}

	if success {
		c.successfulResolutions.Add(1)
	} else {
		c.failedResolutions.Add(1)
	}

	// Update component-specific metrics
	c.mu.Lock()
	compMetrics, exists := c.componentMetrics[componentType]
	if !exists {
		compMetrics = &componentMetrics{
			minResolutionTime: int64(time.Hour),
		}
		c.componentMetrics[componentType] = compMetrics
	}
	c.mu.Unlock()

	compMetrics.mu.Lock()
	compMetrics.totalResolutions++
	compMetrics.totalResolutionTime += durationNs

	if durationNs > compMetrics.maxResolutionTime {
		compMetrics.maxResolutionTime = durationNs
	}
	if durationNs < compMetrics.minResolutionTime {
		compMetrics.minResolutionTime = durationNs
	}

	if success {
		compMetrics.successfulResolutions++
	} else {
		compMetrics.failedResolutions++
	}
	compMetrics.mu.Unlock()
}

// RecordRegistration records a registration operation
func (c *defaultMetricsCollector) RecordRegistration(componentType string, lifecycle Lifecycle) {
	c.totalRegistrations.Add(1)

	if counter, exists := c.byLifecycle[lifecycle]; exists {
		counter.Add(1)
	}
}

// RecordScopeCreation records a scope creation operation
func (c *defaultMetricsCollector) RecordScopeCreation(duration time.Duration) {
	c.totalScopesCreated.Add(1)

	durationNs := int64(duration)
	c.totalScopeTime.Add(durationNs)

	// Update max scope time
	for {
		currentMax := c.maxScopeTime.Load()
		if durationNs <= currentMax {
			break
		}
		if c.maxScopeTime.CompareAndSwap(currentMax, durationNs) {
			break
		}
	}
}

// RecordPoolGet records a pool get operation
func (c *defaultMetricsCollector) RecordPoolGet(componentType string, duration time.Duration, success bool) {
	c.totalPoolGets.Add(1)

	durationNs := int64(duration)
	c.totalPoolGetTime.Add(durationNs)

	if success {
		c.successfulPoolGets.Add(1)
	} else {
		c.failedPoolGets.Add(1)
	}

	// Update component-specific metrics
	c.mu.Lock()
	compMetrics, exists := c.componentMetrics[componentType]
	if !exists {
		compMetrics = &componentMetrics{}
		c.componentMetrics[componentType] = compMetrics
	}
	c.mu.Unlock()

	compMetrics.mu.Lock()
	compMetrics.totalPoolGets++
	compMetrics.totalPoolGetTime += durationNs
	if success {
		compMetrics.successfulPoolGets++
	} else {
		compMetrics.failedPoolGets++
	}
	compMetrics.mu.Unlock()
}

// RecordPoolPut records a pool put operation
func (c *defaultMetricsCollector) RecordPoolPut(componentType string, duration time.Duration) {
	c.totalPoolPuts.Add(1)

	// Update component-specific metrics
	c.mu.Lock()
	compMetrics, exists := c.componentMetrics[componentType]
	if !exists {
		compMetrics = &componentMetrics{}
		c.componentMetrics[componentType] = compMetrics
	}
	c.mu.Unlock()

	compMetrics.mu.Lock()
	compMetrics.totalPoolPuts++
	compMetrics.mu.Unlock()
}

// GetMetrics returns a snapshot of current metrics
func (c *defaultMetricsCollector) GetMetrics() *ContainerMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now()
	totalResolutions := c.totalResolutions.Load()
	successfulResolutions := c.successfulResolutions.Load()
	failedResolutions := c.failedResolutions.Load()
	totalResolutionTime := c.totalResolutionTime.Load()
	maxResolutionTime := c.maxResolutionTime.Load()
	minResolutionTime := c.minResolutionTime.Load()

	totalPoolGets := c.totalPoolGets.Load()
	successfulPoolGets := c.successfulPoolGets.Load()
	failedPoolGets := c.failedPoolGets.Load()
	totalPoolPuts := c.totalPoolPuts.Load()
	totalPoolGetTime := c.totalPoolGetTime.Load()

	// Calculate averages
	var avgResolutionTime time.Duration
	if totalResolutions > 0 {
		avgResolutionTime = time.Duration(totalResolutionTime / totalResolutions)
	}

	var avgPoolGetTime time.Duration
	if totalPoolGets > 0 {
		avgPoolGetTime = time.Duration(totalPoolGetTime / totalPoolGets)
	}

	var avgScopeTime time.Duration
	totalScopes := c.totalScopesCreated.Load()
	if totalScopes > 0 {
		avgScopeTime = time.Duration(c.totalScopeTime.Load() / totalScopes)
	}

	// Calculate rates
	elapsed := now.Sub(c.startTime).Seconds()
	var resolutionRate float64
	if elapsed > 0 {
		resolutionRate = float64(totalResolutions) / elapsed
	}

	var errorRate float64
	if totalResolutions > 0 {
		errorRate = float64(failedResolutions) / float64(totalResolutions)
	}

	var poolHitRate float64
	if totalPoolGets > 0 {
		poolHitRate = float64(successfulPoolGets) / float64(totalPoolGets)
	}

	// Build lifecycle map
	byLifecycle := make(map[Lifecycle]int64)
	for lifecycle, counter := range c.byLifecycle {
		byLifecycle[lifecycle] = counter.Load()
	}

	// Build component metrics
	componentMetrics := make(map[string]*ComponentMetrics)
	for compType, compMetrics := range c.componentMetrics {
		compMetrics.mu.RLock()
		var compAvgResolutionTime time.Duration
		if compMetrics.totalResolutions > 0 {
			compAvgResolutionTime = time.Duration(compMetrics.totalResolutionTime / compMetrics.totalResolutions)
		}

		var compAvgPoolGetTime time.Duration
		if compMetrics.totalPoolGets > 0 {
			compAvgPoolGetTime = time.Duration(compMetrics.totalPoolGetTime / compMetrics.totalPoolGets)
		}

		componentMetrics[compType] = &ComponentMetrics{
			TotalResolutions:      compMetrics.totalResolutions,
			SuccessfulResolutions: compMetrics.successfulResolutions,
			FailedResolutions:     compMetrics.failedResolutions,
			TotalResolutionTime:   time.Duration(compMetrics.totalResolutionTime),
			AverageResolutionTime: compAvgResolutionTime,
			MaxResolutionTime:     time.Duration(compMetrics.maxResolutionTime),
			MinResolutionTime:     time.Duration(compMetrics.minResolutionTime),

			TotalPoolGets:      compMetrics.totalPoolGets,
			SuccessfulPoolGets: compMetrics.successfulPoolGets,
			FailedPoolGets:     compMetrics.failedPoolGets,
			TotalPoolPuts:      compMetrics.totalPoolPuts,
			TotalPoolGetTime:   time.Duration(compMetrics.totalPoolGetTime),
			AveragePoolGetTime: compAvgPoolGetTime,
		}
		compMetrics.mu.RUnlock()
	}

	return &ContainerMetrics{
		Timestamp:               now,
		TotalResolutions:        totalResolutions,
		SuccessfulResolutions:   successfulResolutions,
		FailedResolutions:       failedResolutions,
		AverageResolutionTime:   avgResolutionTime,
		MaxResolutionTime:       time.Duration(maxResolutionTime),
		MinResolutionTime:       time.Duration(minResolutionTime),
		TotalRegistrations:      c.totalRegistrations.Load(),
		ByLifecycle:             byLifecycle,
		TotalScopesCreated:      totalScopes,
		AverageScopeTime:        avgScopeTime,
		MaxScopeTime:            time.Duration(c.maxScopeTime.Load()),
		TotalPoolGets:           totalPoolGets,
		SuccessfulPoolGets:      successfulPoolGets,
		FailedPoolGets:          failedPoolGets,
		TotalPoolPuts:           totalPoolPuts,
		AveragePoolGetTime:      avgPoolGetTime,
		ComponentMetrics:        componentMetrics,
		ResolutionRatePerSecond: resolutionRate,
		ErrorRate:               errorRate,
		PoolHitRate:             poolHitRate,
	}
}

// Reset resets all metrics
func (c *defaultMetricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalResolutions.Store(0)
	c.successfulResolutions.Store(0)
	c.failedResolutions.Store(0)
	c.totalResolutionTime.Store(0)
	c.maxResolutionTime.Store(0)
	c.minResolutionTime.Store(int64(time.Hour))

	c.totalRegistrations.Store(0)
	for _, counter := range c.byLifecycle {
		counter.Store(0)
	}

	c.totalScopesCreated.Store(0)
	c.totalScopeTime.Store(0)
	c.maxScopeTime.Store(0)

	c.totalPoolGets.Store(0)
	c.successfulPoolGets.Store(0)
	c.failedPoolGets.Store(0)
	c.totalPoolPuts.Store(0)
	c.totalPoolGetTime.Store(0)

	c.componentMetrics = make(map[string]*componentMetrics)
	c.startTime = time.Now()
}

// String returns a human-readable representation of metrics
func (m *ContainerMetrics) String() string {
	return fmt.Sprintf(
		"DI Container Metrics (collected at %s):\n"+
			"  Resolutions: %d total, %.2f/sec, %.1f%% error rate\n"+
			"  Resolution Time: avg=%v, min=%v, max=%v\n"+
			"  Registrations: %d total\n"+
			"  Scopes: %d created, avg creation time=%v\n"+
			"  Pool: gets=%d (%.1f%% hit), puts=%d, avg get time=%v\n"+
			"  Components: %d tracked",
		m.Timestamp.Format(time.RFC3339),
		m.TotalResolutions,
		m.ResolutionRatePerSecond,
		m.ErrorRate*100,
		m.AverageResolutionTime,
		m.MinResolutionTime,
		m.MaxResolutionTime,
		m.TotalRegistrations,
		m.TotalScopesCreated,
		m.AverageScopeTime,
		m.TotalPoolGets,
		m.PoolHitRate*100,
		m.TotalPoolPuts,
		m.AveragePoolGetTime,
		len(m.ComponentMetrics),
	)
}

// GetComponentMetrics returns metrics for a specific component
func (m *ContainerMetrics) GetComponentMetrics(componentType string) (*ComponentMetrics, bool) {
	metrics, exists := m.ComponentMetrics[componentType]
	return metrics, exists
}

// GetHealthStatus returns a health status based on metrics
func (m *ContainerMetrics) GetHealthStatus() string {
	if m.ErrorRate > 0.1 { // More than 10% error rate
		return "unhealthy"
	}
	if m.ResolutionRatePerSecond < 1 && m.TotalResolutions > 100 {
		return "degraded"
	}
	return "healthy"
}
