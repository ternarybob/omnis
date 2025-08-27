# omnis

[![Go Reference](https://pkg.go.dev/badge/github.com/ternarybob/omnis.svg)](https://pkg.go.dev/github.com/ternarybob/omnis)
[![Go Report Card](https://goreportcard.com/badge/github.com/ternarybob/omnis)](https://goreportcard.com/report/github.com/ternarybob/omnis)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

High-level web framework integrations and middleware for Go web applications built on Gin.

## Installation

```bash
go get github.com/ternarybob/omnis@latest
```

## Features

- **Configurable Render Service**: Structured API responses with correlation ID tracking
- **Memory Log Integration**: Request-scoped log retrieval using correlation IDs  
- **Custom Logger Support**: Inject your own arbor logger instances
- **Middleware Collection**: Headers, static files, correlation ID, error handling, and recovery
- **Environment-Aware**: Different behavior for DEV/PRD environments
- **No External Dependencies**: Self-contained configuration (removed satus dependency)

## Quick Start

### Basic Usage

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/ternarybob/omnis"
)

func main() {
    r := gin.New()
    
    // Configure your service
    config := &omnis.ServiceConfig{
        Version: "1.0.0",
        Name:    "my-api",
        Support: "support@mycompany.com", 
        Scope:   "DEV",
    }
    
    // Add middleware
    r.Use(omnis.SetCorrelationID())
    r.Use(omnis.SetHeaders(config))
    
    // Use render service
    r.GET("/users", func(c *gin.Context) {
        users := []string{"alice", "bob"}
        omnis.RenderService(c).WithConfig(config).AsResult(200, users)
    })
    
    r.Run(":8080")
}
```

### With Custom Logger

```go
// Configure your arbor logger
logger := arbor.Logger().
    WithConsoleWriter(models.WriterConfiguration{
        Type: models.LogWriterTypeConsole,
    }).
    WithMemoryWriter(models.WriterConfiguration{
        Type: models.LogWriterTypeMemory,
    }).
    WithPrefix("my-api")

// Use with render service for log correlation
r.GET("/data", func(c *gin.Context) {
    // Your business logic here
    logger.Info().Msg("Processing data request")
    
    data := map[string]interface{}{"result": "success"}
    
    // Render with custom logger - memory logs will be included in response
    omnis.RenderService(c).
        WithConfig(config).
        WithLogger(logger).
        AsResult(200, data)
})
```

### Advanced Configuration Examples

#### Complete Service Setup
```go
func setupAPI() *gin.Engine {
    r := gin.New()
    
    // Service configuration
    config := &omnis.ServiceConfig{
        Name:    "pexa-mock-api",
        Version: "0.0.2",
        Support: "support@pexa.com",
        Scope:   "DEV",
    }
    
    // Configure logger with memory storage for correlation
    logger := arbor.Logger().
        WithPrefix("API").
        WithConsoleWriter(models.WriterConfiguration{
            Type: models.LogWriterTypeConsole,
        }).
        WithMemoryWriter(models.WriterConfiguration{
            Type: models.LogWriterTypeMemory,
        })
    
    // Add middleware stack
    r.Use(omnis.SetCorrelationID())
    r.Use(omnis.SetHeaders(config))
    r.Use(omnis.StaticRequests(config, []string{"assets/", "favicon.ico"}))
    
    return r
}
```

#### Multiple Response Types
```go
// Success responses
r.GET("/users", func(c *gin.Context) {
    users := []gin.H{
        {"id": 1, "name": "John Doe"},
        {"id": 2, "name": "Jane Smith"},
    }
    
    omnis.RenderService(c).
        WithConfig(config).
        WithLogger(logger).
        AsResult(200, users)
})

// Error responses with stack traces (DEV mode only)
r.GET("/error-demo", func(c *gin.Context) {
    err := errors.New("demonstration error")
    
    omnis.RenderService(c).
        WithConfig(config).
        WithLogger(logger).
        AsResultWithError(500, nil, err)
})

// Model responses (merge into existing struct)
r.GET("/profile", func(c *gin.Context) {
    profile := &UserProfile{
        ID:   123,
        Name: "John Doe",
    }
    
    omnis.RenderService(c).
        WithConfig(config).
        WithLogger(logger).
        AsModel(200, profile)
})
```

## Examples

### Error Handling with Stack Traces (DEV mode)

```go
r.GET("/error", func(c *gin.Context) {
    err := errors.New("something went wrong")
    
    omnis.RenderService(c).
        WithConfig(config).
        AsResultWithError(500, nil, err)
})
```

### API Response Formats

#### Standard API Response (with correlation tracking)
All responses include metadata, correlation tracking, and memory logs:

```json
{
  "result": {
    "users": ["alice", "bob"]
  },
  "name": "pexa-mock-api",
  "version": "0.0.2+build.go1.24.20250827.105800",
  "support": "support@pexa.com",
  "status": 200,
  "scope": "DEV",
  "correlationid": "550e8400-e29b-41d4-a716-446655440000",
  "request": {
    "url": "/users"
  },
  "log": {
    "001": "INF|10:30:45.123|API|Processing request started",
    "002": "DBG|10:30:45.124|API|Found 2 users",
    "003": "INF|10:30:45.125|API|Request completed"
  }
}
```

#### Error Response with Stack Trace (DEV mode only)
```json
{
  "result": null,
  "name": "pexa-mock-api",
  "version": "0.0.2+build.go1.24.20250827.105800",
  "status": 500,
  "scope": "DEV",
  "correlationid": "550e8400-e29b-41d4-a716-446655440000",
  "error": "demonstration error",
  "stack": [
    "main.errorHandler()",
    "  /app/handlers.go:45",
    "github.com/gin-gonic/gin.(*Context).Next()",
    "  /go/pkg/mod/github.com/gin-gonic/gin@v1.10.1/context.go:174"
  ],
  "request": {
    "url": "/error-demo"
  },
  "log": {
    "001": "ERR|10:30:50.123|API|Error occurred: demonstration error"
  }
}
```

## Documentation

Full documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/ternarybob/omnis).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## API Reference

### ServiceConfig

Configure your service metadata:

```go
type ServiceConfig struct {
    Version string  // Service version (e.g., "1.0.0")
    Name    string  // Service name (e.g., "my-api") 
    Support string  // Support contact (e.g., "support@company.com")
    Scope   string  // Environment scope ("DEV", "PRD", etc.)
}
```

### Render Service Methods

```go
service := omnis.RenderService(ctx)

// Configuration and logger injection
service.WithConfig(config)           // Set service configuration
service.WithLogger(logger)           // Set custom arbor logger

// Response methods
service.AsResult(200, data)          // Success response with data
service.AsError(500, err)            // Error response  
service.AsResultWithError(200, data, err) // Success with error details (DEV only)
service.AsModel(200, &modelStruct)   // Response merged into existing model
```

### Available Middleware

```go
// Correlation ID tracking
r.Use(omnis.SetCorrelationID())

// Service headers (x-t3b-app, x-t3b-version)  
r.Use(omnis.SetHeaders(config))

// Static file handling with cache control
r.Use(omnis.StaticRequests(config, []string{"assets/", "favicon.ico"}))

// Additional middleware available:
// - Error handler middleware
// - Recovery middleware  
// - CORS headers middleware
```

## Migration Guide

### Updating Existing Applications

To integrate omnis into existing Gin applications:

#### 1. Update main.go setup
```go
// Before: Standard Gin setup
func main() {
    r := gin.New()
    
    r.GET("/users", func(c *gin.Context) {
        users := getUserList()
        c.JSON(http.StatusOK, users)
    })
}

// After: With omnis integration
func main() {
    r := gin.New()
    
    config := &omnis.ServiceConfig{
        Name:    "my-api",
        Version: "1.0.0",
        Support: "support@company.com",
        Scope:   "DEV",
    }
    
    logger := arbor.Logger().WithPrefix("API")
    
    // Add omnis middleware
    r.Use(omnis.SetCorrelationID())
    r.Use(omnis.SetHeaders(config))
    
    r.GET("/users", func(c *gin.Context) {
        users := getUserList()
        
        // Enhanced response with correlation tracking
        omnis.RenderService(c).
            WithConfig(config).
            WithLogger(logger).
            AsResult(200, users)
    })
}
```

#### 2. Update Handler Methods
```go
// Before: Direct JSON responses
func UserHandler(c *gin.Context) {
    userInfo := getUserFromToken(c)
    if userInfo == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    c.JSON(http.StatusOK, userInfo)
}

// After: With omnis structured responses
func UserHandler(c *gin.Context) {
    logger := arbor.Logger().WithPrefix("UserHandler")
    config := getServiceConfig() // Your config source
    
    userInfo := getUserFromToken(c)
    if userInfo == nil {
        omnis.RenderService(c).
            WithConfig(config).
            WithLogger(logger).
            AsError(401, errors.New("unauthorized"))
        return
    }
    
    omnis.RenderService(c).
        WithConfig(config).
        WithLogger(logger).
        AsResult(200, userInfo)
}
```

## Related Libraries

This library is part of the ternarybob ecosystem:

- [funktion](https://github.com/ternarybob/funktion) - Core utility functions
- [arbor](https://github.com/ternarybob/arbor) - Structured logging system  
- [omnis](https://github.com/ternarybob/omnis) - Web framework integrations

## Breaking Changes

### v1.0.22+ 
- **Removed satus dependency**: Use `ServiceConfig` instead of `satus.AppConfig`
- **Updated middleware signatures**: Pass `*ServiceConfig` instead of `*satus.AppConfig`
- **Configuration injection**: Use `WithConfig()` method on render service
