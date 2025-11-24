package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/jagreehal/autolemetry-go"
	"github.com/jagreehal/autolemetry-go/logging"
)

func main() {
	// Initialize autolemetry with debug mode
	cleanup, err := autolemetry.Init(context.Background(),
		autolemetry.WithService("logging-example"),
		autolemetry.WithEndpoint("http://localhost:4318"),
		autolemetry.WithDebug(true), // Enable debug mode
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Create logger with trace context handler
	logger := slog.New(logging.NewTraceHandler(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	))

	// Example 1: Logging within a span
	ctx := context.Background()
	ctx, span := autolemetry.Start(ctx, "processOrder")
	defer span.End()

	span.SetAttribute("order.id", "12345")
	logger.InfoContext(ctx, "Processing order", slog.String("order_id", "12345"))

	// Example 2: Manual trace context enrichment
	ctx2, span2 := autolemetry.Start(ctx, "sendEmail")
	defer span2.End()

	attrs := logging.WithTraceContext(ctx2)
	// Convert []slog.Attr to []any for slog.Group
	attrsAny := make([]any, len(attrs))
	for i, attr := range attrs {
		attrsAny[i] = attr
	}
	logger.InfoContext(ctx2, "Sending email",
		slog.String("to", "user@example.com"),
		slog.Group("trace", attrsAny...),
	)

	log.Println("Example completed - check logs for trace context")
}
