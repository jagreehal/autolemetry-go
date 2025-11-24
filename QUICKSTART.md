# Quick Start Guide

Get started with autolemetry-go in 5 minutes.

## Installation

```bash
go get github.com/jagreehal/autolemetry-go
```

## 1. Initialize

```go
package main

import (
    "context"
    "log"
    "github.com/jagreehal/autolemetry-go"
)

func main() {
cleanup, err := autolemetry.Init(context.Background(),
    autolemetry.WithService("my-service"),
    autolemetry.WithEndpoint("http://localhost:4318"),
    // Optional: OTLP vendor preset without extra SDKs
    // autolemetry.WithBackend("datadog"),
)
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()

    // Your application code here
}
```

## 2. Add Tracing

### Option A: Using Start()

```go
func ProcessOrder(ctx context.Context, orderID string) error {
    ctx, span := autolemetry.Start(ctx, "ProcessOrder")
    defer span.End()

    span.SetAttribute("order.id", orderID)
    // ... your code ...
    return nil
}
```

### Option B: Using Trace()

```go
func GetUser(ctx context.Context, userID string) (*User, error) {
    return autolemetry.Trace(ctx, "GetUser", func(ctx context.Context, span autolemetry.Span) (*User, error) {
        span.SetAttribute("user.id", userID)
        return db.FindUser(ctx, userID)
    })
}
```

## 3. Add HTTP Middleware

```go
import "github.com/jagreehal/autolemetry-go/middleware"

mux := http.NewServeMux()
mux.HandleFunc("/users", handleUsers)

handler := middleware.HTTPMiddleware("my-service")(mux)
http.ListenAndServe(":8080", handler)
```

## 4. Add Production Features (Optional)

```go
cleanup, err := autolemetry.Init(context.Background(),
    autolemetry.WithService("my-service"),
    autolemetry.WithEndpoint("http://localhost:4318"),

    // Production hardening
    autolemetry.WithAdaptiveSampler(...),
    autolemetry.WithRateLimit(100, 200),
    autolemetry.WithCircuitBreaker(5, 3, 30*time.Second),
    autolemetry.WithPIIRedaction(...),
)
```

## 5. Add Analytics Events (Optional)

```go
import (
    "github.com/jagreehal/autolemetry-go"
    "github.com/jagreehal/autolemetry-go/subscribers"
)

cleanup, _ := autolemetry.Init(context.Background(),
    autolemetry.WithService("my-service"),
    autolemetry.WithSubscribers(
        subscribers.NewPostHogSubscriber("your-api-key"),
    ),
    autolemetry.WithEventQueue(2000, 500*time.Millisecond, 5),
    autolemetry.WithEventBackoff(100*time.Millisecond, 5*time.Second, 10*time.Second),
)
defer cleanup()

ctx, span := autolemetry.Start(ctx, "userAction")
defer span.End()
autolemetry.Track(ctx, "user_signed_up", map[string]any{
    "user_id": "123",
})
```

## 6. Emit Metrics with Trace Correlation (Optional)

```go
m := autolemetry.Meter()
m.Counter(ctx, "orders.created", 1, map[string]any{"region": "iad"})
m.Histogram(ctx, "orders.duration_ms", float64(time.Since(start).Milliseconds()), nil)
```

## Logging with Trace Context

- Slog: wrap your handler with `logging.NewTraceHandler(...)`.
- Zap: add `logging.TraceFields(ctx)...` to log calls.

## Next Steps

- See [examples/](examples/) for complete examples
- Read [README.md](README.md) for full documentation
- Check [ARCHITECTURE.md](ARCHITECTURE.md) for design details
