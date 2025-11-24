package debugutil

import (
	"github.com/jagreehal/autolemetry-go"
)

// Print proxies to autolemetry.debugPrint to avoid import cycles.
func Print(format string, args ...any) {
	autolemetry.DebugPrintf(format, args...)
}
