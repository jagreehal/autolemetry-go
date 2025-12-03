# ⚠️ This package has moved

This package has been renamed to `github.com/jagreehal/autotel-go`.

## Migration

Update your imports:

```go
// Old
import "github.com/jagreehal/autolemetry-go"

// New  
import "github.com/jagreehal/autotel-go"
```

This package will continue to work as a redirect, but we recommend updating to the new import path.

## New Package

👉 **[github.com/jagreehal/autotel-go](https://github.com/jagreehal/autotel-go)**

All functionality remains the same, just with a new name.

## Quick Start

```bash
go get github.com/jagreehal/autotel-go
```

```go
import "github.com/jagreehal/autotel-go"

func main() {
    cleanup, err := autotel.Init(context.Background(),
        autotel.WithService("my-service"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()
}
```

See the [autotel-go README](https://github.com/jagreehal/autotel-go) for complete documentation.
