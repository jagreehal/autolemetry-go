package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jagreehal/autolemetry-go"
	"github.com/jagreehal/autolemetry-go/middleware"
	"github.com/jagreehal/autolemetry-go/redaction"
	"github.com/jagreehal/autolemetry-go/sampling"
)

func main() {
	// Initialize autolemetry with production hardening features
	cleanup, err := autolemetry.Init(context.Background(),
		autolemetry.WithService("production-service"),
		autolemetry.WithServiceVersion("1.0.0"),
		autolemetry.WithEnvironment("production"),
		autolemetry.WithEndpoint("http://localhost:4318"),

		// Production hardening
		autolemetry.WithAdaptiveSampler(
			sampling.WithBaselineRate(0.1), // 10% baseline sampling
			sampling.WithErrorRate(1.0),    // 100% error sampling
		),
		autolemetry.WithRateLimit(100, 200),                  // 100 spans/sec, burst of 200
		autolemetry.WithCircuitBreaker(5, 3, 30*time.Second), // 5 failures, 3 successes, 30s timeout
		autolemetry.WithPIIRedaction(
			redaction.WithAllowlistKeys("user_id", "order_id"), // Allow these keys
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Create HTTP server with tracing middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/users", handleUsers)

	handler := middleware.HTTPMiddleware("production-service")(mux)

	log.Println("Starting production-ready server on :8080")
	log.Println("Features enabled:")
	log.Println("  - Adaptive sampling (10% baseline, 100% errors)")
	log.Println("  - Rate limiting (100 spans/sec)")
	log.Println("  - Circuit breaker protection")
	log.Println("  - PII redaction")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	_, span := autolemetry.Start(r.Context(), "handleUsers")
	defer span.End()

	// These will be redacted (PII)
	span.SetAttribute("user.email", "user@example.com")
	span.SetAttribute("user.phone", "555-123-4567")

	// These will NOT be redacted (allowlisted)
	span.SetAttribute("user_id", "user-123")
	span.SetAttribute("order_id", "order-456")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`[{"id":"user-123","name":"Alice"}]`))
}
