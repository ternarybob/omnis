// -----------------------------------------------------------------------
// JSON Renderer Middleware
// Provides fluent interface for JSON responses with logging context
// Created: 2025-08-27
// -----------------------------------------------------------------------

package omnis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ternarybob/arbor"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// JSONRendererConfig holds configuration for the JSON renderer middleware
type JSONRendererConfig struct {
	ServiceConfig     *ServiceConfig // Service configuration
	DefaultLogger     arbor.ILogger  // Default logger to use if none specified
	EnablePrettyPrint bool           // Enable pretty printing in development
}

// Note: JSONRenderer struct removed - functionality replaced by:
// 1. Automatic JSON interception via jsonResponseInterceptor
// 2. Gin context extensions via omnis.C(c).WithLogger(log).JSON()

// JSONMiddleware creates middleware that enables fluent JSON rendering with logging context
// Usage: router.Use(omnis.JSONMiddleware(config))
func JSONMiddleware(config *ServiceConfig) gin.HandlerFunc {
	return JSONMiddlewareWithConfig(&JSONRendererConfig{
		ServiceConfig: config,
	})
}

// JSONMiddlewareWithConfig creates middleware with full configuration options
// This middleware intercepts all c.JSON() calls and enhances them with logging
// Usage: router.Use(omnis.JSONMiddlewareWithConfig(&omnis.JSONRendererConfig{...}))
func JSONMiddlewareWithConfig(config *JSONRendererConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create custom response writer that intercepts JSON responses
		originalWriter := c.Writer
		interceptor := &jsonResponseInterceptor{
			ResponseWriter: originalWriter,
			context:        c,
			config:         config,
		}
		c.Writer = interceptor

		// Note: No longer storing JSONRenderer in context
		// Functionality moved to Gin extensions and automatic interception
		c.Next()
	}
}

// jsonResponseInterceptor intercepts JSON responses and enhances them
type jsonResponseInterceptor struct {
	gin.ResponseWriter
	context *gin.Context
	config  *JSONRendererConfig
	written bool
}

// Write intercepts the response and processes JSON content
func (w *jsonResponseInterceptor) Write(data []byte) (int, error) {
	if w.written {
		return w.ResponseWriter.Write(data)
	}

	// Check if this is a JSON response
	contentType := w.Header().Get("Content-Type")
	
	// Debug: Log what we're intercepting
	if w.config != nil && w.config.DefaultLogger != nil {
		w.config.DefaultLogger.Debug().
			Str("content_type", contentType).
			Str("data_preview", string(data)[:min(100, len(data))]).
			Msg("JSON interceptor: Checking response")
	}
	
	if !strings.Contains(contentType, "application/json") {
		return w.ResponseWriter.Write(data)
	}

	w.written = true

	// Get logger from context if available (set by handlers)
	var logger arbor.ILogger
	if loggerInterface, exists := w.context.Get("request_logger"); exists {
		if requestLogger, ok := loggerInterface.(arbor.ILogger); ok {
			logger = requestLogger
		}
	}

	// Fall back to default logger
	if logger == nil && w.config != nil {
		logger = w.config.DefaultLogger
	}

	// Log the response if logger is available
	if logger != nil {
		logger.Debug().
			Int("status_code", w.context.Writer.Status()).
			Str("response_size", fmt.Sprintf("%d bytes", len(data))).
			Msg("JSON response intercepted")
	}

	// Parse the JSON to potentially enhance it
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		// If we can't parse it, just pass it through
		return w.ResponseWriter.Write(data)
	}

	// For simple c.JSON() calls, we can enhance the response here if desired
	// For now, just log and pass through with pretty printing if enabled
	var output []byte
	var writeErr error

	if w.config != nil && (w.config.EnablePrettyPrint || w.isDevelopmentMode()) {
		output, writeErr = json.MarshalIndent(jsonData, "", "  ")
	} else {
		output, writeErr = json.Marshal(jsonData)
	}

	if writeErr != nil {
		return w.ResponseWriter.Write(data) // Fall back to original
	}

	return w.ResponseWriter.Write(output)
}

// WriteHeader captures the status code
func (w *jsonResponseInterceptor) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

// isDevelopmentMode checks if we're in development mode
func (w *jsonResponseInterceptor) isDevelopmentMode() bool {
	if w.config == nil || w.config.ServiceConfig == nil {
		return true // Default to development if no config
	}
	scope := w.config.ServiceConfig.Scope
	return scope == "" || scope == "DEV" || scope == "DEVELOPMENT"
}

// WithLogger stores a logger in the context for the JSON interceptor to use
// Usage: omnis.WithLogger(c, logger) followed by c.JSON(200, data)
func WithLogger(c *gin.Context, logger arbor.ILogger) {
	c.Set("request_logger", logger)
}

// JSONMiddlewareWithDefaults creates middleware with default configuration
// Usage: router.Use(omnis.JSONMiddlewareWithDefaults())
func JSONMiddlewareWithDefaults() gin.HandlerFunc {
	return JSONMiddleware(nil)
}

// Note: Old JSONRenderer methods removed. Functionality replaced by:
// 1. Automatic JSON interception (transparent enhancement of c.JSON calls)
// 2. Gin context extensions: omnis.C(c).WithLogger(log).JSON(200, data)
