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
	if !strings.Contains(contentType, "application/json") {
		return w.ResponseWriter.Write(data)
	}

	w.written = true

	// Get logger from context if available (set by handlers using omnis.REQUEST_LOGGER)
	var logger arbor.ILogger
	if loggerInterface, exists := w.context.Get(REQUEST_LOGGER); exists {
		if requestLogger, ok := loggerInterface.(arbor.ILogger); ok {
			logger = requestLogger
		}
	}

	// Fall back to default logger if available
	if logger == nil && w.config != nil {
		logger = w.config.DefaultLogger
	}

	// If no logger available at all, skip processing and pass through
	skipProcessing := logger == nil
	if skipProcessing {
		return w.ResponseWriter.Write(data)
	}

	// Log the response
	logger.Debug().
		Int("status_code", w.context.Writer.Status()).
		Str("response_size", fmt.Sprintf("%d bytes", len(data))).
		Msg("JSON response intercepted")

	// Parse the JSON to potentially enhance it
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		// If we can't parse it, just pass it through
		return w.ResponseWriter.Write(data)
	}

	// Check if this is already an APIResponse (to avoid double-wrapping)
	if apiResp, ok := jsonData.(map[string]interface{}); ok {
		if _, hasVersion := apiResp["version"]; hasVersion {
			if _, hasName := apiResp["name"]; hasName {
				if _, hasResult := apiResp["result"]; hasResult {
					// Already wrapped, just pretty print if needed
					var output []byte
					var writeErr error
					if w.config != nil && (w.config.EnablePrettyPrint || w.isDevelopmentMode()) {
						output, writeErr = json.MarshalIndent(jsonData, "", "  ")
					} else {
						output, writeErr = json.Marshal(jsonData)
					}
					if writeErr != nil {
						return w.ResponseWriter.Write(data)
					}
					return w.ResponseWriter.Write(output)
				}
			}
		}
	}

	// Wrap the response in APIResponse format
	apiResponse := ApiResponse{
		Version: "1.0.0",
		Build:   "",
		Name:    "",
		Status:  w.context.Writer.Status(),
		Scope:   "",
		Result:  jsonData,
	}

	// Add service config if available
	if w.config != nil && w.config.ServiceConfig != nil {
		apiResponse.Version = w.config.ServiceConfig.Version
		apiResponse.Build = w.config.ServiceConfig.Build
		apiResponse.Name = w.config.ServiceConfig.Name
		apiResponse.Scope = w.config.ServiceConfig.Scope
		// Support field can be set via configuration or left empty
	}

	// Get correlation ID from context
	if correlationID, exists := w.context.Get("correlation-id"); exists {
		if id, ok := correlationID.(string); ok {
			apiResponse.CorrelationId = id
		}
	}

	// Check if this is an error response (typically has "error" field)
	if errResp, ok := jsonData.(map[string]interface{}); ok {
		if errMsg, hasError := errResp["error"]; hasError {
			// Move error to the error field and clear result
			apiResponse.Error = fmt.Sprintf("%v", errMsg)
			apiResponse.Result = nil
		}
	}

	// Marshal the wrapped response
	var output []byte
	var writeErr error

	if w.config != nil && (w.config.EnablePrettyPrint || w.isDevelopmentMode()) {
		output, writeErr = json.MarshalIndent(apiResponse, "", "  ")
	} else {
		output, writeErr = json.Marshal(apiResponse)
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

// JSONMiddlewareWithDefaults creates middleware with default configuration
// Usage: router.Use(omnis.JSONMiddlewareWithDefaults())
func JSONMiddlewareWithDefaults() gin.HandlerFunc {
	return JSONMiddleware(nil)
}

// Note: Old JSONRenderer methods removed. Functionality replaced by:
// 1. Automatic JSON interception (transparent enhancement of c.JSON calls)
// 2. Gin context extensions: omnis.C(c).WithLogger(log).JSON(200, data)
