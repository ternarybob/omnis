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

### API Response Structure

All responses include metadata and request correlation:

```json
{
  "version": "1.0.0",
  "name": "my-api", 
  "support": "support@mycompany.com",
  "status": 200,
  "scope": "DEV",
  "correlationid": "550e8400-e29b-41d4-a716-446655440000",
  "request": {
    "url": "/users",
    "method": "GET"
  },
  "log": {
    "001": "INF|10:30:45.123|my-api|Processing request",
    "002": "INF|10:30:45.124|my-api|Found 2 users"
  },
  "result": ["alice", "bob"]
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
