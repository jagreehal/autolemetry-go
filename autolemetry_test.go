package autolemetry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jagreehal/autolemetry-go"
)

func TestInit_Basic(t *testing.T) {
	cleanup, err := autolemetry.Init(context.Background(),
		autolemetry.WithService("test-service"),
	)
	require.NoError(t, err)
	require.NotNil(t, cleanup)
	defer cleanup()

	// Verify tracer provider is set
	// (Check global otel.GetTracerProvider())
	// Note: We can't directly access GetTracerProvider without exposing it
	// This test verifies Init() completes without error
}

func TestInit_WithCustomEndpoint(t *testing.T) {
	cleanup, err := autolemetry.Init(context.Background(),
		autolemetry.WithService("test"),
		autolemetry.WithEndpoint("custom:4318"),
	)
	require.NoError(t, err)
	defer cleanup()
}

func TestInit_WithGRPCProtocol(t *testing.T) {
	cleanup, err := autolemetry.Init(context.Background(),
		autolemetry.WithService("test"),
		autolemetry.WithProtocol(autolemetry.ProtocolGRPC),
	)
	require.NoError(t, err)
	defer cleanup()
}
