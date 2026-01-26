package di

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Trace represents a single trace in a distributed system
type Trace struct {
	// Trace ID (unique across the system)
	TraceID string

	// Parent span ID (empty for root spans)
	ParentSpanID string

	// Span ID (unique within trace)
	SpanID string

	// Operation name
	Operation string

	// Component name
	Component string

	// Start time
	StartTime time.Time

	// End time (zero if not finished)
	EndTime time.Time

	// Duration (calculated when EndTime is set)
	Duration time.Duration

	// Tags (key-value pairs for filtering/grouping)
	Tags map[string]string

	// Logs (structured log entries)
	Logs []TraceLog

	// Error message (if any)
	Error string

	// Status code (e.g., HTTP status, gRPC code)
	StatusCode int

	// Metadata (additional structured data)
	Metadata map[string]interface{}
}

// TraceLog represents a log entry within a trace
type TraceLog struct {
	Timestamp time.Time
	Message   string
	Fields    map[string]interface{}
}

// Span represents an active span in a trace
type Span struct {
	trace    *Trace
	tracer   *Tracer
	children []*Span
	mu       sync.RWMutex
}

// Tracer manages distributed tracing
type Tracer struct {
	// Service name
	ServiceName string

	// Environment (prod, staging, dev, etc.)
	Environment string

	// Trace storage
	traces map[string]*Trace
	spans  map[string]*Span
	mu     sync.RWMutex

	// Sampling rate (0.0 to 1.0)
	SamplingRate float64

	// Exporters for sending traces to external systems
	exporters []TraceExporter

	// Enable/disable tracing
	Enabled bool
}

// TraceExporter exports traces to external systems
type TraceExporter interface {
	// Export exports a trace
	Export(trace *Trace) error

	// Shutdown gracefully shuts down the exporter
	Shutdown() error
}

// NewTracer creates a new tracer
func NewTracer(serviceName, environment string) *Tracer {
	return &Tracer{
		ServiceName:  serviceName,
		Environment:  environment,
		traces:       make(map[string]*Trace),
		spans:        make(map[string]*Span),
		SamplingRate: 1.0, // Sample all traces by default
		Enabled:      true,
	}
}

// StartTrace starts a new trace
func (t *Tracer) StartTrace(operation, component string) *Span {
	if !t.Enabled {
		return &Span{tracer: t}
	}

	// Check sampling
	if !t.shouldSample() {
		return &Span{tracer: t}
	}

	traceID := generateID()
	spanID := generateID()

	trace := &Trace{
		TraceID:   traceID,
		SpanID:    spanID,
		Operation: operation,
		Component: component,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
		Logs:      make([]TraceLog, 0),
		Metadata:  make(map[string]interface{}),
	}

	// Add default tags
	trace.Tags["service"] = t.ServiceName
	trace.Tags["environment"] = t.Environment
	trace.Tags["component"] = component

	span := &Span{
		trace:  trace,
		tracer: t,
	}

	t.mu.Lock()
	t.traces[traceID] = trace
	t.spans[spanID] = span
	t.mu.Unlock()

	return span
}

// StartSpan starts a new span within an existing trace
func (t *Tracer) StartSpan(parentSpan *Span, operation, component string) *Span {
	if !t.Enabled || parentSpan == nil || parentSpan.trace == nil {
		return &Span{tracer: t}
	}

	spanID := generateID()

	trace := &Trace{
		TraceID:      parentSpan.trace.TraceID,
		ParentSpanID: parentSpan.trace.SpanID,
		SpanID:       spanID,
		Operation:    operation,
		Component:    component,
		StartTime:    time.Now(),
		Tags:         make(map[string]string),
		Logs:         make([]TraceLog, 0),
		Metadata:     make(map[string]interface{}),
	}

	// Copy parent tags
	for k, v := range parentSpan.trace.Tags {
		trace.Tags[k] = v
	}
	trace.Tags["component"] = component

	span := &Span{
		trace:  trace,
		tracer: t,
	}

	// Add to parent
	parentSpan.mu.Lock()
	parentSpan.children = append(parentSpan.children, span)
	parentSpan.mu.Unlock()

	t.mu.Lock()
	t.spans[spanID] = span
	t.mu.Unlock()

	return span
}

// StartSpanFromContext starts a span from context
func (t *Tracer) StartSpanFromContext(ctx context.Context, operation, component string) *Span {
	if !t.Enabled {
		return &Span{tracer: t}
	}

	// Try to get parent span from context
	if parentSpan, ok := ctx.Value(spanContextKey).(*Span); ok {
		return t.StartSpan(parentSpan, operation, component)
	}

	// No parent in context, start new trace
	return t.StartTrace(operation, component)
}

// End finishes a span
func (s *Span) End() {
	if s.trace == nil || !s.tracer.Enabled {
		return
	}

	s.trace.EndTime = time.Now()
	s.trace.Duration = s.trace.EndTime.Sub(s.trace.StartTime)

	// Export trace if this is the root span
	if s.trace.ParentSpanID == "" {
		s.tracer.exportTrace(s.trace)
	}

	// Remove from active spans
	s.tracer.mu.Lock()
	delete(s.tracer.spans, s.trace.SpanID)
	s.tracer.mu.Unlock()
}

// EndWithError finishes a span with an error
func (s *Span) EndWithError(err error) {
	if s.trace == nil || !s.tracer.Enabled {
		return
	}

	if err != nil {
		s.trace.Error = err.Error()
		s.trace.Tags["error"] = "true"
	}

	s.End()
}

// SetTag sets a tag on the span
func (s *Span) SetTag(key, value string) {
	if s.trace == nil || !s.tracer.Enabled {
		return
	}

	s.trace.Tags[key] = value
}

// SetTags sets multiple tags on the span
func (s *Span) SetTags(tags map[string]string) {
	if s.trace == nil || !s.tracer.Enabled {
		return
	}

	for k, v := range tags {
		s.trace.Tags[k] = v
	}
}

// Log adds a log entry to the span
func (s *Span) Log(message string, fields map[string]interface{}) {
	if s.trace == nil || !s.tracer.Enabled {
		return
	}

	log := TraceLog{
		Timestamp: time.Now(),
		Message:   message,
		Fields:    fields,
	}

	s.trace.Logs = append(s.trace.Logs, log)
}

// SetMetadata sets metadata on the span
func (s *Span) SetMetadata(key string, value interface{}) {
	if s.trace == nil || !s.tracer.Enabled {
		return
	}

	s.trace.Metadata[key] = value
}

// SetStatusCode sets the status code on the span
func (s *Span) SetStatusCode(code int) {
	if s.trace == nil || !s.tracer.Enabled {
		return
	}

	s.trace.StatusCode = code
}

// Context returns a context with the span
func (s *Span) Context(ctx context.Context) context.Context {
	if s.trace == nil || !s.tracer.Enabled {
		return ctx
	}

	return context.WithValue(ctx, spanContextKey, s)
}

// GetTrace returns the trace
func (s *Span) GetTrace() *Trace {
	if s.trace == nil {
		return nil
	}
	return s.trace
}

// GetChildren returns child spans
func (s *Span) GetChildren() []*Span {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.children
}

// shouldSample determines if a trace should be sampled
func (t *Tracer) shouldSample() bool {
	if t.SamplingRate >= 1.0 {
		return true
	}
	if t.SamplingRate <= 0.0 {
		return false
	}
	// Simple random sampling
	// In production, use a proper sampling algorithm
	return time.Now().UnixNano()%100 < int64(t.SamplingRate*100)
}

// exportTrace exports a trace to all exporters
func (t *Tracer) exportTrace(trace *Trace) {
	t.mu.RLock()
	exporters := t.exporters
	t.mu.RUnlock()

	for _, exporter := range exporters {
		go func(e TraceExporter) {
			if err := e.Export(trace); err != nil {
				// Log export error (in production, use proper logging)
				fmt.Printf("Failed to export trace: %v\n", err)
			}
		}(exporter)
	}
}

// AddExporter adds a trace exporter
func (t *Tracer) AddExporter(exporter TraceExporter) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.exporters = append(t.exporters, exporter)
}

// GetTrace returns a trace by ID
func (t *Tracer) GetTrace(traceID string) (*Trace, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	trace, exists := t.traces[traceID]
	return trace, exists
}

// GetActiveTraces returns all active traces
func (t *Tracer) GetActiveTraces() []*Trace {
	t.mu.RLock()
	defer t.mu.RUnlock()

	traces := make([]*Trace, 0, len(t.traces))
	for _, trace := range t.traces {
		traces = append(traces, trace)
	}
	return traces
}

// GetActiveSpans returns all active spans
func (t *Tracer) GetActiveSpans() []*Span {
	t.mu.RLock()
	defer t.mu.RUnlock()

	spans := make([]*Span, 0, len(t.spans))
	for _, span := range t.spans {
		spans = append(spans, span)
	}
	return spans
}

// SetSamplingRate sets the sampling rate
func (t *Tracer) SetSamplingRate(rate float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if rate < 0.0 {
		rate = 0.0
	}
	if rate > 1.0 {
		rate = 1.0
	}
	t.SamplingRate = rate
}

// Enable enables tracing
func (t *Tracer) Enable() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Enabled = true
}

// Disable disables tracing
func (t *Tracer) Disable() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Enabled = false
}

// Shutdown shuts down the tracer and all exporters
func (t *Tracer) Shutdown() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var errs []error
	for _, exporter := range t.exporters {
		if err := exporter.Shutdown(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to shutdown exporters: %v", errs)
	}
	return nil
}

// ConsoleExporter exports traces to console (for debugging)
type ConsoleExporter struct{}

// Export exports a trace to console
func (e *ConsoleExporter) Export(trace *Trace) error {
	fmt.Printf("[TRACE] %s %s %s %v\n",
		trace.TraceID,
		trace.Operation,
		trace.Component,
		trace.Duration)
	return nil
}

// Shutdown shuts down the console exporter
func (e *ConsoleExporter) Shutdown() error {
	return nil
}

// NewConsoleExporter creates a new console exporter
func NewConsoleExporter() TraceExporter {
	return &ConsoleExporter{}
}

// context key for storing span in context
type contextKey string

const spanContextKey = contextKey("span")

// generateID generates a unique ID
func generateID() string {
	// In production, use a proper ID generator (e.g., UUID)
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

// TraceMiddleware provides tracing middleware for HTTP handlers
type TraceMiddleware struct {
	tracer *Tracer
}

// NewTraceMiddleware creates a new trace middleware
func NewTraceMiddleware(tracer *Tracer) *TraceMiddleware {
	return &TraceMiddleware{
		tracer: tracer,
	}
}

// Wrap wraps an HTTP handler with tracing
func (m *TraceMiddleware) Wrap(handler http.Handler, operation string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start trace
		span := m.tracer.StartTrace(operation, "http")
		defer span.End()

		// Add trace ID to response headers
		trace := span.GetTrace()
		if trace != nil {
			w.Header().Set("X-Trace-ID", trace.TraceID)
		}

		// Set tags
		span.SetTag("http.method", r.Method)
		span.SetTag("http.url", r.URL.String())
		span.SetTag("http.user_agent", r.UserAgent())

		// Create response writer to capture status code
		rw := &traceResponseWriter{
			ResponseWriter: w,
			span:           span,
		}

		// Create context with span
		ctx := span.Context(r.Context())

		// Execute handler
		handler.ServeHTTP(rw, r.WithContext(ctx))

		// Set status code tag
		span.SetStatusCode(rw.statusCode)
		span.SetTag("http.status_code", fmt.Sprintf("%d", rw.statusCode))

		if rw.statusCode >= 400 {
			span.SetTag("error", "true")
		}
	})
}

// traceResponseWriter captures the status code for tracing
type traceResponseWriter struct {
	http.ResponseWriter
	span       *Span
	statusCode int
}

// WriteHeader captures the status code
func (rw *traceResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes data and captures status if not set
func (rw *traceResponseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

// GetTracer returns the tracer
func (m *TraceMiddleware) GetTracer() *Tracer {
	return m.tracer
}
