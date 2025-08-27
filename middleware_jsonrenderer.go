// -----------------------------------------------------------------------
// JSON Renderer Middleware
// Provides fluent interface for JSON responses with logging context
// Created: 2025-08-27
// -----------------------------------------------------------------------

package omnis

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ternarybob/arbor"
)

// JSONRendererConfig holds configuration for the JSON renderer middleware
type JSONRendererConfig struct {
	ServiceConfig   *ServiceConfig  // Service configuration
	DefaultLogger   arbor.ILogger   // Default logger to use if none specified
	EnablePrettyPrint bool          // Enable pretty printing in development
}

// JSONRenderer provides a fluent interface for rendering JSON responses with logging
type JSONRenderer struct {
	ctx           *gin.Context
	logger        arbor.ILogger
	config        *ServiceConfig
	defaultLogger arbor.ILogger
	enablePretty  bool
}

// JSONMiddleware creates middleware that enables fluent JSON rendering with logging context
// Usage: router.Use(omnis.JSONMiddleware(config))
func JSONMiddleware(config *ServiceConfig) gin.HandlerFunc {
	return JSONMiddlewareWithConfig(&JSONRendererConfig{
		ServiceConfig: config,
	})
}

// JSONMiddlewareWithConfig creates middleware with full configuration options
// Usage: router.Use(omnis.JSONMiddlewareWithConfig(&omnis.JSONRendererConfig{...}))
func JSONMiddlewareWithConfig(config *JSONRendererConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store JSON renderer in context for later use
		renderer := &JSONRenderer{
			ctx:           c,
			config:        config.ServiceConfig,
			defaultLogger: config.DefaultLogger,
			enablePretty:  config.EnablePrettyPrint,
		}
		c.Set("json_renderer", renderer)
		c.Next()
	}
}

// JSONMiddlewareWithDefaults creates middleware with default configuration
// Usage: router.Use(omnis.JSONMiddlewareWithDefaults())
func JSONMiddlewareWithDefaults() gin.HandlerFunc {
	return JSONMiddleware(nil)
}

// JSON creates a new JSONRenderer from the gin context
// Usage: omnis.JSON(c).WithLogger(log).Response(200, data)
func JSON(c *gin.Context) *JSONRenderer {
	// Try to get existing renderer from context first
	if renderer, exists := c.Get("json_renderer"); exists {
		if jr, ok := renderer.(*JSONRenderer); ok {
			return jr
		}
	}
	
	// Create new renderer if not in middleware
	return &JSONRenderer{
		ctx: c,
	}
}

// WithLogger sets the logger for this JSON renderer
func (j *JSONRenderer) WithLogger(logger arbor.ILogger) *JSONRenderer {
	j.logger = logger
	return j
}

// WithConfig sets the service configuration for this JSON renderer
func (j *JSONRenderer) WithConfig(config *ServiceConfig) *JSONRenderer {
	j.config = config
	return j
}

// Response renders a successful JSON response with the provided data
// This integrates with the existing omnis RenderService for consistent formatting
func (j *JSONRenderer) Response(code int, data interface{}) {
	render := RenderService(j.ctx)
	
	if j.logger != nil {
		render = render.WithLogger(j.logger)
	}
	
	if j.config != nil {
		render = render.WithConfig(j.config)
	}
	
	render.AsResult(code, data)
}

// Error renders an error JSON response
func (j *JSONRenderer) Error(code int, err interface{}) {
	render := RenderService(j.ctx)
	
	if j.logger != nil {
		render = render.WithLogger(j.logger)
	}
	
	if j.config != nil {
		render = render.WithConfig(j.config)
	}
	
	render.AsError(code, err)
}

// ResultWithError renders a JSON response with both data and error information
func (j *JSONRenderer) ResultWithError(code int, data interface{}, err error) {
	render := RenderService(j.ctx)
	
	if j.logger != nil {
		render = render.WithLogger(j.logger)
	}
	
	if j.config != nil {
		render = render.WithConfig(j.config)
	}
	
	render.AsResultWithError(code, data, err)
}

// Simple renders a simple JSON response without the omnis wrapper
// This is for cases where you want direct JSON output
func (j *JSONRenderer) Simple(code int, data interface{}) {
	logger := j.getEffectiveLogger()
	if logger != nil {
		// Log the response for debugging
		logger.Debug().
			Int("status_code", code).
			Str("response_data", fmt.Sprintf("%+v", data)).
			Msg("Simple JSON response")
	}
	
	// Use pretty print if enabled and in development scope
	if j.enablePretty || j.isDevelopmentScope() {
		j.ctx.IndentedJSON(code, data)
	} else {
		j.ctx.JSON(code, data)
	}
}

// IndentedSimple renders an indented JSON response (for development)
func (j *JSONRenderer) IndentedSimple(code int, data interface{}) {
	logger := j.getEffectiveLogger()
	if logger != nil {
		// Log the response for debugging
		logger.Debug().
			Int("status_code", code).
			Str("response_data", fmt.Sprintf("%+v", data)).
			Msg("Indented JSON response")
	}
	
	j.ctx.IndentedJSON(code, data)
}

// getEffectiveLogger returns the logger to use (explicit logger takes precedence)
func (j *JSONRenderer) getEffectiveLogger() arbor.ILogger {
	if j.logger != nil {
		return j.logger
	}
	return j.defaultLogger
}

// isDevelopmentScope checks if we're in development mode
func (j *JSONRenderer) isDevelopmentScope() bool {
	if j.config == nil {
		return true // Default to development if no config
	}
	scope := j.config.Scope
	return scope == "" || scope == "DEV" || scope == "DEVELOPMENT"
}

// Success is a convenience method for 200 OK responses
func (j *JSONRenderer) Success(data interface{}) {
	j.Response(http.StatusOK, data)
}

// Created is a convenience method for 201 Created responses
func (j *JSONRenderer) Created(data interface{}) {
	j.Response(http.StatusCreated, data)
}

// BadRequest is a convenience method for 400 Bad Request responses
func (j *JSONRenderer) BadRequest(message interface{}) {
	j.Error(http.StatusBadRequest, message)
}

// Unauthorized is a convenience method for 401 Unauthorized responses
func (j *JSONRenderer) Unauthorized(message interface{}) {
	j.Error(http.StatusUnauthorized, message)
}

// NotFound is a convenience method for 404 Not Found responses
func (j *JSONRenderer) NotFound(message interface{}) {
	j.Error(http.StatusNotFound, message)
}

// InternalServerError is a convenience method for 500 Internal Server Error responses
func (j *JSONRenderer) InternalServerError(message interface{}) {
	j.Error(http.StatusInternalServerError, message)
}