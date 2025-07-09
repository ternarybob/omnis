package omnis

import (
	"github.com/gin-gonic/gin"
)

// RequestScopedLoggerMiddleware creates a fresh logger instance for each request
// with correlation ID pre-configured. This eliminates the need for correlation ID
// lifecycle management and cleanup - the logger instance dies with the request.
//
// This is the cleanest architectural approach:
// - No global state pollution
// - No cleanup required
// - No correlation ID leaking between requests
// - Handlers get a dedicated logger instance
//
// Usage:
//
//	router.Use(omnis.SetCorrelationID())
//	router.Use(omnis.RequestScopedLoggerMiddleware(func() interface{} {
//	    return arbor.NewLogger() // Create fresh instance
//	}))
//	// In handlers: logger := c.MustGet("logger")
func RequestScopedLoggerMiddleware(createLogger func() interface{}) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		correlationID := c.GetString(CORRELATION_ID_KEY)

		if createLogger != nil {
			// Create a fresh logger instance for this request
			requestLogger := createLogger()

			// Set correlation ID on the request-scoped logger
			if correlationID != "" {
				if arborLogger, ok := requestLogger.(interface {
					WithCorrelationId(string) interface{}
				}); ok {
					// Configure the logger with correlation ID
					requestLogger = arborLogger.WithCorrelationId(correlationID)
				}
			}

			// Store the configured logger in Gin context
			c.Set("logger", requestLogger)

			// Optional: also store as "arbor" for backward compatibility
			c.Set("arbor", requestLogger)
		}

		// Process request - handlers use c.MustGet("logger")
		c.Next()

		// No cleanup needed - logger instance dies with request context
	})
}
