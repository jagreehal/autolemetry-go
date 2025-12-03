// Package middleware redirects to github.com/jagreehal/autotel-go/middleware.
//
// This package has been renamed. Please update your imports:
//
//	// Old
//	import "github.com/jagreehal/autolemetry-go/middleware"
//
//	// New
//	import "github.com/jagreehal/autotel-go/middleware"
package middleware

import (
	"github.com/jagreehal/autotel-go/middleware"
)

// Re-export all middleware functions
var (
	HTTPMiddleware         = middleware.HTTPMiddleware
	HTTPMiddlewareWithOptions = middleware.HTTPMiddlewareWithOptions
	GinMiddleware         = middleware.GinMiddleware
	GRPCServerHandler     = middleware.GRPCServerHandler
	GRPCClientHandler     = middleware.GRPCClientHandler
)

