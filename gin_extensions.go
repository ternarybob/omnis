// -----------------------------------------------------------------------
// Gin Context Extensions
// Extends gin.Context with enhanced JSON methods and logging
// Created: 2025-08-27
// -----------------------------------------------------------------------

package omnis

import (
	"github.com/gin-gonic/gin"
	"github.com/ternarybob/arbor"
)

// ExtendedContext extends gin.Context with additional methods
type ExtendedContext struct {
	*gin.Context
}

// WithLogger creates a logger-enhanced context for structured JSON responses
// Usage: omnis.WithLogger(c, log).StructuredJson(200, data)
func WithLogger(c *gin.Context, logger arbor.ILogger) *ContextExtension {
	// Set the logger in context for the interceptor to use
	c.Set("request_logger", logger)

	return &ContextExtension{
		Context: c,
		logger:  logger,
	}
}

// ContextExtension wraps gin.Context to provide enhanced methods
type ContextExtension struct {
	*gin.Context
	logger arbor.ILogger
}

// StructuredJson renders a JSON response with automatic logging and middleware enhancement
// This is the primary method for structured JSON responses with logging
func (ce *ContextExtension) StructuredJson(code int, obj interface{}) {
	// Use the standard Gin JSON method - will be intercepted by our middleware
	ce.Context.JSON(code, obj)
}

