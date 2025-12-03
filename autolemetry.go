// Package autolemetry redirects to github.com/jagreehal/autotel-go.
//
// This package has been renamed to autotel-go.
// Please update your imports:
//
//	// Old
//	import "github.com/jagreehal/autolemetry-go"
//
//	// New
//	import "github.com/jagreehal/autotel-go"
//
// This package will continue to work as a redirect, but we recommend updating to the new import path.
package autolemetry

import (
	"context"

	"github.com/jagreehal/autotel-go"
	"github.com/jagreehal/autotel-go/middleware"
	"github.com/jagreehal/autotel-go/subscribers"
	"github.com/jagreehal/autotel-go/logging"
)

// Re-export all main functions
var (
	Init             = autotel.Init
	InitWithConfig   = autotel.InitWithConfig
	Start            = autotel.Start
	Track            = autotel.Track
	Meter            = autotel.Meter
	SetAttribute     = autotel.SetAttribute
	SetAttributes    = autotel.SetAttributes
	AddEvent         = autotel.AddEvent
	RecordError      = autotel.RecordError
	GetTraceID       = autotel.GetTraceID
	GetSpanID        = autotel.GetSpanID
	GetVersion       = autotel.GetVersion
	GetOperationName = autotel.GetOperationName
	IsTracingEnabled = autotel.IsTracingEnabled
	SetDuration      = autotel.SetDuration
	SetHTTPRequestAttributes = autotel.SetHTTPRequestAttributes
	AddEventWithAttributes   = autotel.AddEventWithAttributes
)

// Trace is a generic function wrapper for autotel.Trace
func Trace[T any](ctx context.Context, name string, fn func(context.Context, Span) (T, error)) (T, error) {
	return autotel.Trace(ctx, name, fn)
}

// TraceNoError is a generic function wrapper for autotel.TraceNoError
func TraceNoError[T any](ctx context.Context, name string, fn func(context.Context, Span) T) T {
	return autotel.TraceNoError(ctx, name, fn)
}

// TraceVoid is a wrapper for autotel.TraceVoid
func TraceVoid(ctx context.Context, name string, fn func(context.Context, Span) error) error {
	return autotel.TraceVoid(ctx, name, fn)
}

// TraceFunc is a wrapper for autotel.TraceFunc
func TraceFunc(ctx context.Context, name string, fn any) any {
	return autotel.TraceFunc(ctx, name, fn)
}

// Re-export types
type (
	Span         = autotel.Span
	Config       = autotel.Config
	Option       = autotel.Option
	Protocol     = autotel.Protocol
	EventTracker = autotel.EventTracker
	TraceContext = autotel.TraceContext
	Metric       = autotel.Metric
)

// Re-export constants
const (
	ProtocolHTTP = autotel.ProtocolHTTP
	ProtocolGRPC = autotel.ProtocolGRPC
	Version      = autotel.Version
)

// Re-export middleware functions
var (
	HTTPMiddleware         = middleware.HTTPMiddleware
	HTTPMiddlewareWithOptions = middleware.HTTPMiddlewareWithOptions
	GinMiddleware         = middleware.GinMiddleware
	GRPCServerHandler     = middleware.GRPCServerHandler
	GRPCClientHandler     = middleware.GRPCClientHandler
)

// Re-export subscriber types and functions
var (
	NewPostHogSubscriber   = subscribers.NewPostHogSubscriber
	NewMixpanelSubscriber  = subscribers.NewMixpanelSubscriber
	NewAmplitudeSubscriber = subscribers.NewAmplitudeSubscriber
	NewWebhookSubscriber   = subscribers.NewWebhookSubscriber
	NewQueue              = subscribers.NewQueue
	NewQueueWithConfig    = subscribers.NewQueueWithConfig
	NewInMemorySubscriber = subscribers.NewInMemorySubscriber
)

// Re-export subscriber types
type (
	Subscriber        = subscribers.Subscriber
	PostHogSubscriber = subscribers.PostHogSubscriber
	MixpanelSubscriber = subscribers.MixpanelSubscriber
	AmplitudeSubscriber = subscribers.AmplitudeSubscriber
	WebhookSubscriber = subscribers.WebhookSubscriber
	InMemorySubscriber = subscribers.InMemorySubscriber
	InMemoryEvent     = subscribers.InMemoryEvent
	Queue            = subscribers.Queue
	QueueConfig      = subscribers.QueueConfig
)

// Re-export subscriber option types
type (
	PostHogOption   = subscribers.PostHogOption
	MixpanelOption  = subscribers.MixpanelOption
	AmplitudeOption = subscribers.AmplitudeOption
	WebhookOption   = subscribers.WebhookOption
)

// Re-export subscriber option functions
var (
	WithPostHogHost        = subscribers.WithPostHogHost
	WithPostHogDistinctID  = subscribers.WithPostHogDistinctID
	WithPostHogTimeout     = subscribers.WithPostHogTimeout
	WithMixpanelHost       = subscribers.WithMixpanelHost
	WithMixpanelAPISecret  = subscribers.WithMixpanelAPISecret
	WithMixpanelDistinctID = subscribers.WithMixpanelDistinctID
	WithMixpanelTimeout    = subscribers.WithMixpanelTimeout
	WithAmplitudeHost      = subscribers.WithAmplitudeHost
	WithAmplitudeUserID    = subscribers.WithAmplitudeUserID
	WithAmplitudeDeviceID  = subscribers.WithAmplitudeDeviceID
	WithAmplitudeTimeout   = subscribers.WithAmplitudeTimeout
	WithWebhookHeaders     = subscribers.WithWebhookHeaders
	WithWebhookTimeout     = subscribers.WithWebhookTimeout
)

// Re-export logging functions
var (
	NewTraceHandler  = logging.NewTraceHandler
	WithTraceContext = logging.WithTraceContext
	TraceFields      = logging.TraceFields
)

// Re-export logging types
type (
	TraceHandler = logging.TraceHandler
)

// Re-export option functions
var (
	WithService            = autotel.WithService
	WithServiceVersion     = autotel.WithServiceVersion
	WithEnvironment        = autotel.WithEnvironment
	WithEndpoint           = autotel.WithEndpoint
	WithProtocol           = autotel.WithProtocol
	WithSampler            = autotel.WithSampler
	WithInsecure           = autotel.WithInsecure
	WithRateLimit          = autotel.WithRateLimit
	WithCircuitBreaker     = autotel.WithCircuitBreaker
	WithPIIRedaction       = autotel.WithPIIRedaction
	WithAdaptiveSampler    = autotel.WithAdaptiveSampler
	WithDebug              = autotel.WithDebug
	WithSubscribers        = autotel.WithSubscribers
	WithBackend            = autotel.WithBackend
	WithOTLPHeaders        = autotel.WithOTLPHeaders
	WithSpanExporters      = autotel.WithSpanExporters
	WithSpanProcessors     = autotel.WithSpanProcessors
	WithBatchTimeout       = autotel.WithBatchTimeout
	WithMaxQueueSize       = autotel.WithMaxQueueSize
	WithMaxExportBatchSize = autotel.WithMaxExportBatchSize
	WithEventQueue         = autotel.WithEventQueue
	WithEventBackoff      = autotel.WithEventBackoff
	WithEventRetry        = autotel.WithEventRetry
	WithMetrics           = autotel.WithMetrics
	WithMetricExporters   = autotel.WithMetricExporters
	WithMetricInterval    = autotel.WithMetricInterval
)
