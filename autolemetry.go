package autolemetry

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	otlpmetricgrpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otlpmetrichttp "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/jagreehal/autolemetry-go/internal/exporters"
)

// EventTracker is an interface for tracking analytics events.
// This avoids import cycles by not importing the analytics package directly.
type EventTracker interface {
	Track(ctx context.Context, event string, properties map[string]any)
	Shutdown(ctx context.Context) error
}

var (
	globalTracker   EventTracker
	globalTrackerMu sync.RWMutex
	queueFactory    func(cfg *Config, subscribers []Subscriber) EventTracker
)

// RegisterQueueFactory registers a function to create analytics queues.
// This is called by the analytics package to avoid import cycles.
// Users should not call this directly.
func RegisterQueueFactory(factory func(cfg *Config, subscribers []Subscriber) EventTracker) {
	queueFactory = factory
}

// Init initializes autolemetry with OpenTelemetry SDK using functional options.
// Returns a cleanup function that should be called on shutdown.
//
// This is the recommended way to initialize autolemetry for most users.
//
// Example:
//
//	cleanup, err := autolemetry.Init(ctx,
//	    autolemetry.WithService("my-service"),
//	    autolemetry.WithEndpoint("http://localhost:4318"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer cleanup()
func Init(ctx context.Context, opts ...Option) (func(), error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return initWithConfig(ctx, cfg)
}

// InitWithConfig initializes autolemetry with a custom Config.
// This provides advanced users with full control over configuration.
//
// Most users should use Init() with functional options instead.
//
// Example:
//
//	cfg := autolemetry.DefaultConfig()
//	cfg.ServiceName = "my-service"
//	cfg.Endpoint = "custom:4318"
//	cfg.Sampler = trace.AlwaysSample() // Custom sampler
//	cleanup, err := autolemetry.InitWithConfig(ctx, cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer cleanup()
func InitWithConfig(ctx context.Context, cfg *Config) (func(), error) {
	return initWithConfig(ctx, cfg)
}

func initWithConfig(ctx context.Context, cfg *Config) (func(), error) {
	applyEnvOverrides(cfg)
	applyBackendPreset(cfg)

	if cfg.Debug == nil {
		enabled := ShouldEnableDebug(nil)
		cfg.Debug = &enabled
	}

	if *cfg.Debug {
		EnableDebug()
	} else {
		DisableDebug()
	}

	// Build resource attributes
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
		resource.WithFromEnv(), // Discover resource from OTEL_RESOURCE_ATTRIBUTES env
		resource.WithProcess(), // Add process attributes
		resource.WithOS(),      // Add OS attributes
		resource.WithHost(),    // Add host attributes
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	exportersList, err := buildExporters(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := setupMetrics(ctx, res, cfg); err != nil {
		return nil, err
	}

	processors := make([]trace.SpanProcessor, 0, len(exportersList)+len(cfg.SpanProcessors))
	processors = append(processors, cfg.SpanProcessors...)
	for _, exp := range exportersList {
		processors = append(processors, trace.NewBatchSpanProcessor(exp,
			trace.WithBatchTimeout(cfg.BatchTimeout),
			trace.WithMaxQueueSize(cfg.MaxQueueSize),
			trace.WithMaxExportBatchSize(cfg.MaxExportBatchSize),
		))
	}

	// Use AlwaysSample in debug mode to see all spans
	sampler := cfg.Sampler
	if IsDebugEnabled() {
		sampler = trace.AlwaysSample()
		debugPrint("Using AlwaysSample sampler for debug mode")
	}

	// Create tracer provider
	providerOpts := []trace.TracerProviderOption{
		trace.WithResource(res),
		trace.WithSampler(sampler),
	}
	for _, processor := range processors {
		providerOpts = append(providerOpts, trace.WithSpanProcessor(processor))
	}

	tp := trace.NewTracerProvider(providerOpts...)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global production hardening features
	if cfg.RateLimiter != nil {
		setGlobalRateLimiter(cfg.RateLimiter)
	}
	if cfg.CircuitBreaker != nil {
		setGlobalCircuitBreaker(cfg.CircuitBreaker)
	}
	if cfg.PIIRedactor != nil {
		setGlobalPIIRedactor(cfg.PIIRedactor)
	}

	// Create global analytics queue if subscribers are provided
	if len(cfg.Subscribers) > 0 && queueFactory != nil {
		globalTrackerMu.Lock()
		globalTracker = queueFactory(cfg, cfg.Subscribers)
		globalTrackerMu.Unlock()
	}

	// Return cleanup function
	cleanup := func() {
		_ = tp.Shutdown(context.Background())

		// Shutdown global analytics queue if it exists
		globalTrackerMu.Lock()
		if globalTracker != nil {
			_ = globalTracker.Shutdown(context.Background())
			globalTracker = nil
		}
		globalTrackerMu.Unlock()
	}

	return cleanup, nil
}

// Track sends an analytics event to the global queue (if configured).
// This is a convenience function that uses the queue created during Init().
// If no subscribers were provided during Init(), this function does nothing.
//
// Example:
//
//	autolemetry.Track(ctx, "user_signed_up", map[string]any{
//	    "user_id": "123",
//	    "plan":    "premium",
//	})
func Track(ctx context.Context, event string, properties map[string]any) {
	globalTrackerMu.RLock()
	tracker := globalTracker
	globalTrackerMu.RUnlock()

	if tracker != nil {
		tracker.Track(ctx, event, properties)
	}
}

func setupMetrics(ctx context.Context, res *resource.Resource, cfg *Config) error {
	if !cfg.MetricsEnabled {
		return nil
	}

	exportersList := cfg.MetricExporters
	if len(exportersList) == 0 {
		if cfg.Endpoint == "" {
			// No exporter configured and no endpoint provided; skip metrics setup.
			return nil
		}

		exp, err := newOTLPMetricsExporter(ctx, cfg)
		if err != nil {
			return fmt.Errorf("failed to create metrics exporter: %w", err)
		}
		exportersList = append(exportersList, exp)
	}

	providerOpts := []sdkmetric.Option{sdkmetric.WithResource(res)}
	for _, exp := range exportersList {
		providerOpts = append(providerOpts, sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp, sdkmetric.WithInterval(cfg.MetricInterval))))
	}

	mp := sdkmetric.NewMeterProvider(providerOpts...)
	otel.SetMeterProvider(mp)
	return nil
}

func newOTLPMetricsExporter(ctx context.Context, cfg *Config) (sdkmetric.Exporter, error) {
	if cfg.Protocol == ProtocolHTTP {
		httpOpts := []otlptmetricOption{
			otlptmetricWithEndpoint(cfg.Endpoint),
			otlptmetricWithHeaders(cfg.OTLPHeaders),
			otlptmetricWithTimeout(cfg.BatchTimeout + 5*time.Second),
		}
		if cfg.Insecure {
			httpOpts = append(httpOpts, otlptmetricWithInsecure())
		}
		return otlpmetrichttp.New(ctx, httpOpts...)
	}

	grpcOpts := []otlpgmetricOption{
		otlpgmetricWithEndpoint(cfg.Endpoint),
		otlpgmetricWithHeaders(cfg.OTLPHeaders),
		otlpgmetricWithTimeout(cfg.BatchTimeout + 5*time.Second),
	}
	if cfg.Insecure {
		grpcOpts = append(grpcOpts, otlpmetricgrpc.WithInsecure())
	}
	return otlpmetricgrpc.New(ctx, grpcOpts...)
}

// Aliases to keep option slices readable.
type (
	otlptmetricOption = otlpmetrichttp.Option
	otlpgmetricOption = otlpmetricgrpc.Option
)

func otlptmetricWithEndpoint(e string) otlptmetricOption { return otlpmetrichttp.WithEndpoint(e) }
func otlptmetricWithHeaders(h map[string]string) otlptmetricOption {
	return otlpmetrichttp.WithHeaders(h)
}
func otlptmetricWithTimeout(d time.Duration) otlptmetricOption { return otlpmetrichttp.WithTimeout(d) }
func otlptmetricWithInsecure() otlptmetricOption               { return otlpmetrichttp.WithInsecure() }

func otlpgmetricWithEndpoint(e string) otlpgmetricOption { return otlpmetricgrpc.WithEndpoint(e) }
func otlpgmetricWithHeaders(h map[string]string) otlpgmetricOption {
	return otlpmetricgrpc.WithHeaders(h)
}
func otlpgmetricWithTimeout(d time.Duration) otlpgmetricOption { return otlpmetricgrpc.WithTimeout(d) }

// buildExporters builds the list of span exporters respecting presets and custom exporters.
func buildExporters(ctx context.Context, cfg *Config) ([]trace.SpanExporter, error) {
	exportersList := make([]trace.SpanExporter, 0, 4)
	exportersList = append(exportersList, cfg.SpanExporters...)

	// Base OTLP exporter (only when endpoint configured)
	if cfg.Endpoint != "" {
		exp, err := newOTLPExporter(ctx, cfg)
		if err != nil {
			return nil, err
		}
		exportersList = append(exportersList, exp)
	}

	// Console exporter for debug ergonomics
	if IsDebugEnabled() {
		exportersList = append(exportersList, exporters.NewConsoleExporter())
	}

	return exportersList, nil
}

func newOTLPExporter(ctx context.Context, cfg *Config) (trace.SpanExporter, error) {
	if cfg.Protocol == ProtocolHTTP {
		httpOpts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(cfg.Endpoint),
			otlptracehttp.WithHeaders(cfg.OTLPHeaders),
			otlptracehttp.WithTimeout(cfg.BatchTimeout + 5*time.Second),
		}
		if cfg.Insecure {
			httpOpts = append(httpOpts, otlptracehttp.WithInsecure())
		}
		return otlptracehttp.New(ctx, httpOpts...)
	}

	grpcOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithHeaders(cfg.OTLPHeaders),
		otlptracegrpc.WithTimeout(cfg.BatchTimeout + 5*time.Second),
	}
	if cfg.Insecure {
		grpcOpts = append(grpcOpts, otlptracegrpc.WithInsecure())
	}
	return otlptracegrpc.New(ctx, grpcOpts...)
}

// applyBackendPreset adjusts config for common vendors while staying OTLP-first.
func applyBackendPreset(cfg *Config) {
	switch strings.ToLower(cfg.BackendPreset) {
	case "datadog", "dd":
		if cfg.Endpoint == "" || cfg.Endpoint == "localhost:4318" {
			cfg.Endpoint = "api.datadoghq.com:443"
		}
		cfg.Protocol = ProtocolGRPC
		cfg.Insecure = false
		ensureHeaders(cfg)
		if key := os.Getenv("DD_API_KEY"); key != "" {
			cfg.OTLPHeaders["DD-API-KEY"] = key
		}
		cfg.OTLPHeaders["X-Datadog-Origin"] = "otlp"
	case "honeycomb", "hny":
		if cfg.Endpoint == "" || cfg.Endpoint == "localhost:4318" {
			cfg.Endpoint = "api.honeycomb.io:443"
		}
		cfg.Protocol = ProtocolHTTP
		cfg.Insecure = false
		ensureHeaders(cfg)
		if key := os.Getenv("HONEYCOMB_API_KEY"); key != "" {
			cfg.OTLPHeaders["x-honeycomb-team"] = key
		}
		if dataset := os.Getenv("HONEYCOMB_DATASET"); dataset != "" {
			cfg.OTLPHeaders["x-honeycomb-dataset"] = dataset
		}
	case "grafana", "grafana-cloud", "grafana_cloud":
		if cfg.Endpoint == "" || cfg.Endpoint == "localhost:4318" {
			cfg.Endpoint = "otlp-gateway-prod.grafana.net:443"
		}
		cfg.Protocol = ProtocolGRPC
		cfg.Insecure = false
		ensureHeaders(cfg)
		if key := os.Getenv("GRAFANA_OTLP_API_KEY"); key != "" {
			cfg.OTLPHeaders["Authorization"] = "Bearer " + key
		}
	default:
		// OTLP defaults already set
	}
}

func ensureHeaders(cfg *Config) {
	if cfg.OTLPHeaders == nil {
		cfg.OTLPHeaders = make(map[string]string)
	}
}
