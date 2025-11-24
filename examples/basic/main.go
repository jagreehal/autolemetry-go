package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jagreehal/autolemetry-go"
)

func main() {
	// Initialize autolemetry with debug mode
	cleanup, err := autolemetry.Init(context.Background(),
		autolemetry.WithService("basic-example"),
		autolemetry.WithDebug(true), // Enable verbose console logging
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Example 1: Using Start() with defer
	ctx := context.Background()
	ctx, span := autolemetry.Start(ctx, "ProcessOrder")
	defer span.End()

	span.SetAttribute("order.id", "12345")
	span.SetAttribute("order.amount", 99.99)

	// Example 2: Using Trace() helper
	result, err := autolemetry.Trace(ctx, "GetUser", func(ctx context.Context, span autolemetry.Span) (string, error) {
		span.SetAttribute("user.id", "user-123")
		return "user-data", nil
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}

	fmt.Println("Example completed successfully!")
}
