package autolemetry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/jagreehal/autolemetry-go"
	autolemetrytesting "github.com/jagreehal/autolemetry-go/testing"
)

func TestSpan_SetAttribute(t *testing.T) {
	exporter, cleanup := autolemetrytesting.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	_, span := autolemetry.Start(ctx, "test-span")
	span.SetAttribute("string", "value")
	span.SetAttribute("int", 42)
	span.SetAttribute("int64", int64(100))
	span.SetAttribute("float64", 3.14)
	span.SetAttribute("bool", true)
	span.End()

	// SimpleSpanProcessor exports synchronously, so spans should be available immediately
	spans := exporter.GetSpans()
	assert.GreaterOrEqual(t, len(spans), 1)
}

func TestSpan_AddEvent(t *testing.T) {
	exporter, cleanup := autolemetrytesting.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	_, span := autolemetry.Start(ctx, "test-span")
	span.AddEvent("test-event", attribute.String("key", "value"))
	span.End()

	spans := exporter.GetSpans()
	assert.GreaterOrEqual(t, len(spans), 1)
}

func TestSpan_SetStatus(t *testing.T) {
	exporter, cleanup := autolemetrytesting.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	_, span := autolemetry.Start(ctx, "test-span")
	span.SetStatus(codes.Ok, "success")
	span.End()

	spans := exporter.GetSpans()
	assert.GreaterOrEqual(t, len(spans), 1)
}

func TestSpan_RecordError(t *testing.T) {
	exporter, cleanup := autolemetrytesting.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	_, span := autolemetry.Start(ctx, "test-span")
	err := assert.AnError
	span.RecordError(err)
	span.End()

	spans := exporter.GetSpans()
	assert.GreaterOrEqual(t, len(spans), 1)
}
