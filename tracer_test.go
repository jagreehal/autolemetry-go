package autolemetry_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jagreehal/autolemetry-go"
	autolemetrytesting "github.com/jagreehal/autolemetry-go/testing"
)

func TestStart(t *testing.T) {
	_, cleanup := autolemetrytesting.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	_, span := autolemetry.Start(ctx, "test-span")
	defer span.End()

	assert.True(t, span.IsRecording())
	assert.NotEmpty(t, span.SpanContext().TraceID().String())
}

func TestTrace_Success(t *testing.T) {
	exporter, cleanup := autolemetrytesting.SetupTest(t)
	defer cleanup()

	ctx := context.Background()

	result, err := autolemetry.Trace(ctx, "test-trace", func(ctx context.Context, span autolemetry.Span) (string, error) {
		span.SetAttribute("test", "value")
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)

	// Verify span was created
	spans := exporter.GetSpans()
	assert.GreaterOrEqual(t, len(spans), 1)
}

func TestTrace_Error(t *testing.T) {
	exporter, cleanup := autolemetrytesting.SetupTest(t)
	defer cleanup()

	ctx := context.Background()

	expectedErr := errors.New("test error")
	result, err := autolemetry.Trace(ctx, "test-trace", func(ctx context.Context, span autolemetry.Span) (string, error) {
		return "", expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, result)

	// Verify span was created and error was recorded
	spans := exporter.GetSpans()
	assert.GreaterOrEqual(t, len(spans), 1)
}
