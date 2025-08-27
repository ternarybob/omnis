// -----------------------------------------------------------------------
// Gin Context Extensions
// Extends gin.Context with enhanced JSON methods and logging
// Created: 2025-08-27
// -----------------------------------------------------------------------

package omnis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ternarybob/arbor"
)

// ExtendedContext extends gin.Context with additional methods
type ExtendedContext struct {
	*gin.Context
}

// LoggerChain provides fluent logger chaining for gin.Context
// Usage: omnis.LoggerChain(c, log).JSON(200, data)
// Alternative: omnis.Chain(c).WithLogger(log).JSON(200, data)
func LoggerChain(c *gin.Context, logger arbor.ILogger) *ContextExtension {
	// Set the logger in context for the interceptor to use
	c.Set("request_logger", logger)

	return &ContextExtension{
		Context: c,
		logger:  logger,
	}
}

// Chain provides an alternative entry point for fluent chaining
// Usage: omnis.Chain(c).WithLogger(log).JSON(200, data)
func Chain(c *gin.Context) *ContextExtension {
	return &ContextExtension{
		Context: c,
		logger:  nil,
	}
}

// ContextExtension wraps gin.Context to provide enhanced methods
type ContextExtension struct {
	*gin.Context
	logger arbor.ILogger
}

// C creates an enhanced context with Gin extensions
// Usage: omnis.C(c).WithLogger(log).JSON(200, data)
func C(c *gin.Context) *ContextExtension {
	return &ContextExtension{
		Context: c,
		logger:  nil,
	}
}

// WithLogger sets a logger for fluent chaining
// Usage: omnis.C(c).WithLogger(log).JSON(200, data)
func (ce *ContextExtension) WithLogger(logger arbor.ILogger) *ContextExtension {
	ce.logger = logger
	return ce
}

// JSON renders a JSON response with automatic logging and enhancement
// This method provides the fluent interface: c.WithLogger(log).JSON(200, data)
func (ce *ContextExtension) JSON(code int, obj interface{}) {
	// Set the logger in context for the interceptor to use
	if ce.logger != nil {
		ce.Context.Set("request_logger", ce.logger)
	}

	// Use the standard Gin JSON method - will be intercepted by our middleware
	ce.Context.JSON(code, obj)
}

// IndentedJSON renders an indented JSON response with automatic logging
func (ce *ContextExtension) IndentedJSON(code int, obj interface{}) {
	// Set the logger in context for the interceptor to use
	if ce.logger != nil {
		ce.Context.Set("request_logger", ce.logger)
	}

	// Use the standard Gin IndentedJSON method
	ce.Context.IndentedJSON(code, obj)
}

// Success is a convenience method for 200 OK responses
func (ce *ContextExtension) Success(obj interface{}) {
	ce.JSON(http.StatusOK, obj)
}

// Created is a convenience method for 201 Created responses
func (ce *ContextExtension) Created(obj interface{}) {
	ce.JSON(http.StatusCreated, obj)
}

// BadRequest is a convenience method for 400 Bad Request responses
func (ce *ContextExtension) BadRequest(obj interface{}) {
	ce.JSON(http.StatusBadRequest, obj)
}

// Unauthorized is a convenience method for 401 Unauthorized responses
func (ce *ContextExtension) Unauthorized(obj interface{}) {
	ce.JSON(http.StatusUnauthorized, obj)
}

// NotFound is a convenience method for 404 Not Found responses
func (ce *ContextExtension) NotFound(obj interface{}) {
	ce.JSON(http.StatusNotFound, obj)
}

// InternalServerError is a convenience method for 500 Internal Server Error responses
func (ce *ContextExtension) InternalServerError(obj interface{}) {
	ce.JSON(http.StatusInternalServerError, obj)
}

// Enhanced renders enhanced JSON response with omnis wrapper (full API response)
func (ce *ContextExtension) Enhanced(code int, obj interface{}) {
	render := RenderService(ce.Context)

	if ce.logger != nil {
		render = render.WithLogger(ce.logger)
	}

	render.AsResult(code, obj)
}

// Error renders an error response with omnis wrapper
func (ce *ContextExtension) Error(code int, err interface{}) {
	render := RenderService(ce.Context)

	if ce.logger != nil {
		render = render.WithLogger(ce.logger)
	}

	render.AsError(code, err)
}
