package subscribers

import (
	"github.com/jagreehal/autolemetry-go"
)

func init() {
	// Register the queue factory with autolemetry to avoid import cycles.
	// This allows autolemetry.Init() to create queues when subscribers are provided.
	autolemetry.RegisterQueueFactory(func(cfg *autolemetry.Config, subscribers []autolemetry.Subscriber) autolemetry.EventTracker {
		subs := make([]Subscriber, len(subscribers))
		for i, s := range subscribers {
			subs[i] = s.(Subscriber)
		}
		qc := QueueConfig{
			QueueSize:        cfg.EventQueueSize,
			FlushInterval:    cfg.EventFlushInterval,
			CircuitThreshold: cfg.EventCBThreshold,
			BackoffMin:       cfg.EventBackoffMin,
			BackoffMax:       cfg.EventBackoffMax,
			CircuitReset:     cfg.EventCBReset,
		}
		return NewQueueWithConfig(qc, subs...)
	})
}
