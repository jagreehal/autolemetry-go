// Package subscribers redirects to github.com/jagreehal/autotel-go/subscribers.
//
// This package has been renamed. Please update your imports:
//
//	// Old
//	import "github.com/jagreehal/autolemetry-go/subscribers"
//
//	// New
//	import "github.com/jagreehal/autotel-go/subscribers"
package subscribers

import (
	"github.com/jagreehal/autotel-go/subscribers"
)

// Re-export all subscriber functions
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
	WithAmplitudeDeviceID   = subscribers.WithAmplitudeDeviceID
	WithAmplitudeTimeout   = subscribers.WithAmplitudeTimeout
	WithWebhookHeaders     = subscribers.WithWebhookHeaders
	WithWebhookTimeout     = subscribers.WithWebhookTimeout
)
