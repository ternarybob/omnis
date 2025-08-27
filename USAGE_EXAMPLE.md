# JSON Renderer Middleware Usage Examples

The JSON renderer middleware provides two approaches for enhanced JSON responses with integrated logging:

1. **Automatic Interception**: Transparently enhances all `c.JSON()` calls
2. **Fluent Interface**: Explicit control with `omnis.JSON(c).WithLogger(log).Success(data)`

## Basic Setup

### Simple Setup (Recommended for most cases)
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/ternarybob/omnis"
    "github.com/ternarybob/arbor"
)

func main() {
    r := gin.New()
    
    // Basic setup with service configuration
    config := &omnis.ServiceConfig{
        Name:    "pexa-mock-api",
        Version: "0.0.2",
        Scope:   "DEV",
    }
    
    // Add correlation ID middleware first
    r.Use(omnis.SetCorrelationID())
    
    // Add JSON renderer middleware
    r.Use(omnis.JSONMiddleware(config))
    
    r.GET("/api/data", func(c *gin.Context) {
        log := arbor.GetLogger().WithPrefix("DataHandler")
        
        data := gin.H{
            "message": "Hello World",
            "items": []string{"item1", "item2"},
        }
        
        // Approach 1: Automatic interception (recommended)
        omnis.WithLogger(c, log)
        c.JSON(200, data)  // Standard Gin - automatically enhanced
        
        // Approach 2: Fluent interface (explicit control)
        // omnis.JSON(c).WithLogger(log).Simple(200, data)
    })
    
    r.Run(":8080")
}
```

### Advanced Setup with Full Configuration
```go
func main() {
    r := gin.New()
    
    // Create logger that will be used as default
    logger := arbor.GetLogger().WithPrefix("API")
    
    config := &omnis.ServiceConfig{
        Name:    "pexa-mock-api",
        Version: "0.0.2", 
        Scope:   "DEV",
    }
    
    // Advanced configuration
    jsonConfig := &omnis.JSONRendererConfig{
        ServiceConfig:     config,
        DefaultLogger:     logger,
        EnablePrettyPrint: true, // Always pretty print
    }
    
    r.Use(omnis.SetCorrelationID())
    r.Use(omnis.JSONMiddlewareWithConfig(jsonConfig))
    
    // Handlers will automatically use the default logger and config
    r.GET("/api/users", func(c *gin.Context) {
        users := []gin.H{
            {"id": 1, "name": "John Doe"},
            {"id": 2, "name": "Jane Smith"},
        }
        
        // Uses default logger and config from middleware
        omnis.JSON(c).Success(users)
    })
    
    r.Run(":8080")
}
```

## Usage Patterns

### Automatic Interception (Recommended)

The middleware automatically intercepts all `c.JSON()` calls and enhances them:

```go
r.GET("/api/users", func(c *gin.Context) {
    log := arbor.GetLogger().WithPrefix("UsersHandler")
    
    // Set logger for automatic logging
    omnis.WithLogger(c, log)
    
    // Standard Gin JSON response - automatically enhanced with:
    // - Automatic logging of response details
    // - Pretty printing in development mode
    // - Response size tracking
    c.JSON(200, gin.H{
        "users": []string{"user1", "user2"},
        "count": 2,
    })
})

// No omnis.WithLogger call - uses default logger from middleware
r.GET("/api/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy"})  // Still enhanced automatically
})
```

### 1. Simple JSON Responses (No Omnis Wrapper)
```go
r.GET("/simple", func(c *gin.Context) {
    log := arbor.GetLogger().WithPrefix("SimpleHandler")
    
    data := gin.H{"status": "ok", "message": "Simple response"}
    
    // Direct JSON output
    omnis.JSON(c).WithLogger(log).Simple(200, data)
    
    // Indented JSON (for debugging)
    omnis.JSON(c).WithLogger(log).IndentedSimple(200, data)
})
```

### 2. Omnis-Wrapped Responses (Full API Response)
```go
r.GET("/api/data", func(c *gin.Context) {
    log := arbor.GetLogger().WithPrefix("APIHandler")
    
    data := gin.H{"users": []string{"user1", "user2"}}
    
    // Full omnis response with logging and correlation tracking
    omnis.JSON(c).WithLogger(log).Response(200, data)
})
```

### 3. Convenience Methods
```go
r.GET("/api/resource", func(c *gin.Context) {
    log := arbor.GetLogger().WithPrefix("ResourceHandler")
    
    // Success (200)
    omnis.JSON(c).WithLogger(log).Success(gin.H{"data": "success"})
    
    // Created (201)
    omnis.JSON(c).WithLogger(log).Created(gin.H{"id": 123})
    
    // Bad Request (400)
    omnis.JSON(c).WithLogger(log).BadRequest(gin.H{"error": "invalid input"})
    
    // Not Found (404)
    omnis.JSON(c).WithLogger(log).NotFound(gin.H{"error": "resource not found"})
    
    // Internal Server Error (500)
    omnis.JSON(c).WithLogger(log).InternalServerError(gin.H{"error": "server error"})
})
```

### 4. Error Handling
```go
r.GET("/api/process", func(c *gin.Context) {
    log := arbor.GetLogger().WithPrefix("ProcessHandler")
    
    result, err := someProcessingFunction()
    
    if err != nil {
        // Return error response
        omnis.JSON(c).WithLogger(log).Error(500, err)
        return
    }
    
    // Return success with potential error info
    omnis.JSON(c).WithLogger(log).ResultWithError(200, result, nil)
})
```

## Integration with Existing PEXA Mock API

### Update main.go
```go
func main() {
    // ... existing setup ...
    
    // Replace existing middleware setup with:
    config := configService.ToOmnisServiceConfig()
    logger := common.GetLogger().WithPrefix("API")
    
    jsonConfig := &omnis.JSONRendererConfig{
        ServiceConfig:     config,
        DefaultLogger:     logger,
        EnablePrettyPrint: configService.GetConfig().Server.GinMode == "debug",
    }
    
    r.Use(omnis.SetCorrelationID())
    r.Use(omnis.JSONMiddlewareWithConfig(jsonConfig))
    
    // ... rest of setup ...
}
```

### Update Handlers
```go
// Before
func (h *UserHandlerService) UserHandler(c *gin.Context) {
    userInfo := h.getUserInfoFromToken(c)
    if userInfo == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    c.JSON(http.StatusOK, userInfo)
}

// After
func (h *UserHandlerService) UserHandler(c *gin.Context) {
    log := common.GetLogger().WithPrefix("UserHandler")
    
    userInfo := h.getUserInfoFromToken(c)
    if userInfo == nil {
        omnis.JSON(c).WithLogger(log).Unauthorized(gin.H{"error": "Unauthorized"})
        return
    }
    
    omnis.JSON(c).WithLogger(log).Success(userInfo)
}
```

## Features

- **Fluent Interface**: Chain methods for clean, readable code
- **Automatic Logging**: Integrated with arbor logger for correlation tracking
- **Configuration Support**: Service metadata included in responses
- **Development Mode**: Automatic pretty printing in development
- **Error Handling**: Comprehensive error response support
- **Performance**: Minimal overhead, only logs when logger is provided
- **Compatibility**: Works with or without middleware setup

## Response Formats

### Simple Response
```json
{
  "message": "Hello World",
  "status": "success"
}
```

### Omnis Response (with correlation tracking)
```json
{
  "result": {
    "message": "Hello World",
    "status": "success"
  },
  "name": "pexa-mock-api",
  "version": "0.0.2+build.go1.24.4.20250827.105800",
  "scope": "DEV",
  "status": 200,
  "correlationid": "550e8400-e29b-41d4-a716-446655440000",
  "request": {
    "url": "/api/data"
  },
  "log": {
    "001": "INF|Processing request started",
    "002": "DBG|Data retrieved successfully"
  }
}
```