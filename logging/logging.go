// Package logging redirects to github.com/jagreehal/autotel-go/logging.
//
// This package has been renamed. Please update your imports:
//
//	// Old
//	import "github.com/jagreehal/autolemetry-go/logging"
//
//	// New
//	import "github.com/jagreehal/autotel-go/logging"
package logging

import (
	"github.com/jagreehal/autotel-go/logging"
)

// Re-export all logging functions
var (
	NewTraceHandler  = logging.NewTraceHandler
	WithTraceContext = logging.WithTraceContext
	TraceFields      = logging.TraceFields
)

// Re-export logging types
type (
	TraceHandler = logging.TraceHandler
)
