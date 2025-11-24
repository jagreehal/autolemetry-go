package autolemetry_test

import (
	"context"
	"fmt"
	"log"

	"github.com/jagreehal/autolemetry-go"
)

// Example demonstrates basic usage of autolemetry.
func Example() {
	ctx := context.Background()

	// Initialize autolemetry
	cleanup, err := autolemetry.Init(ctx,
		autolemetry.WithService("example-service"),
		autolemetry.WithEndpoint("localhost:4318"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Create a traced operation
	_, span := autolemetry.Start(ctx, "example-operation")
	defer span.End()

	span.SetAttribute("example.key", "example-value")

	fmt.Println("Tracing initialized")
	// Output: Tracing initialized
}

// Example_init demonstrates initialization with various options.
func Example_init() {
	ctx := context.Background()

	cleanup, err := autolemetry.Init(ctx,
		autolemetry.WithService("my-service"),
		autolemetry.WithServiceVersion("1.0.0"),
		autolemetry.WithEnvironment("production"),
		autolemetry.WithEndpoint("localhost:4318"),
		autolemetry.WithProtocol(autolemetry.ProtocolHTTP),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Println("Initialized with custom options")
	// Output: Initialized with custom options
}

// Example_start demonstrates creating spans with Start.
func Example_start() {
	ctx := context.Background()

	// Start a new span
	_, span := autolemetry.Start(ctx, "database-query")
	defer span.End()

	// Set attributes
	span.SetAttribute("db.system", "postgresql")
	span.SetAttribute("db.operation", "SELECT")
	span.SetAttribute("db.table", "users")

	fmt.Println("Span created with attributes")
	// Output: Span created with attributes
}

// Example_trace demonstrates using the Trace helper for automatic error handling.
func Example_trace() {
	ctx := context.Background()

	// Use Trace helper with automatic error handling
	result, err := autolemetry.Trace(ctx, "process-data",
		func(ctx context.Context, span autolemetry.Span) (string, error) {
			span.SetAttribute("operation", "processing")
			// Your business logic here
			return "processed-data", nil
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:", result)
	// Output: Result: processed-data
}

// Example_withRateLimit demonstrates rate limiting configuration.
func Example_withRateLimit() {
	ctx := context.Background()

	cleanup, err := autolemetry.Init(ctx,
		autolemetry.WithService("rate-limited-service"),
		autolemetry.WithRateLimit(100, 200), // 100 spans/sec, burst of 200
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Println("Rate limiting enabled")
	// Output: Rate limiting enabled
}

// Example_withCircuitBreaker demonstrates circuit breaker configuration.
func Example_withCircuitBreaker() {
	ctx := context.Background()

	cleanup, err := autolemetry.Init(ctx,
		autolemetry.WithService("resilient-service"),
		autolemetry.WithCircuitBreaker(
			5,  // failure threshold
			3,  // success threshold
			30, // timeout in seconds
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Println("Circuit breaker enabled")
	// Output: Circuit breaker enabled
}

// Example_withDebug demonstrates debug mode for development.
func Example_withDebug() {
	ctx := context.Background()

	cleanup, err := autolemetry.Init(ctx,
		autolemetry.WithService("debug-service"),
		autolemetry.WithDebug(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Println("Debug mode enabled")
	// Output: Debug mode enabled
}

// Example_nestedSpans demonstrates creating nested spans.
func Example_nestedSpans() {
	ctx := context.Background()

	// Parent span
	ctx, parentSpan := autolemetry.Start(ctx, "parent-operation")
	defer parentSpan.End()

	parentSpan.SetAttribute("level", "parent")

	// Child span
	_, childSpan := autolemetry.Start(ctx, "child-operation")
	childSpan.SetAttribute("level", "child")
	childSpan.End()

	fmt.Println("Nested spans created")
	// Output: Nested spans created
}

// Example_errorHandling demonstrates automatic error recording.
func Example_errorHandling() {
	ctx := context.Background()

	_, span := autolemetry.Start(ctx, "risky-operation")
	defer span.End()

	// Simulate an error
	err := fmt.Errorf("something went wrong")
	if err != nil {
		span.RecordError(err) // Automatically sets span status to ERROR
	}

	fmt.Println("Error recorded in span")
	// Output: Error recorded in span
}

// Example_traceNoError demonstrates Trace helper for operations without errors.
func Example_traceNoError() {
	ctx := context.Background()

	result := autolemetry.TraceNoError(ctx, "calculation",
		func(ctx context.Context, span autolemetry.Span) int {
			span.SetAttribute("operation", "add")
			return 42 + 8
		},
	)

	fmt.Println("Result:", result)
	// Output: Result: 50
}

// Example_traceVoid demonstrates Trace helper for operations that return only errors.
func Example_traceVoid() {
	ctx := context.Background()

	err := autolemetry.TraceVoid(ctx, "write-operation",
		func(ctx context.Context, span autolemetry.Span) error {
			span.SetAttribute("operation", "write")
			// Your operation here
			return nil
		},
	)

	if err == nil {
		fmt.Println("Operation completed successfully")
	}
	// Output: Operation completed successfully
}

// Example_initWithConfig demonstrates advanced configuration using Config directly.
// This is for advanced users who need full control over configuration.
func Example_initWithConfig() {
	ctx := context.Background()

	// Get default config and customize it
	cfg := autolemetry.DefaultConfig()
	cfg.ServiceName = "advanced-service"
	cfg.ServiceVersion = "2.0.0"
	cfg.Environment = "production"
	cfg.Endpoint = "otel-collector:4318"
	cfg.Protocol = autolemetry.ProtocolHTTP
	cfg.Insecure = false

	// Initialize with custom config
	cleanup, err := autolemetry.InitWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Println("Initialized with custom Config")
	// Output: Initialized with custom Config
}
