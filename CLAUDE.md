# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Test
```bash
go build .                    # Build the library
go test ./...                 # Run all tests
go test -v ./...              # Run tests with verbose output
go test -run TestName ./...   # Run specific test
```

### Module Management
```bash
go mod tidy                   # Clean up dependencies
go mod download               # Download dependencies
```

## Architecture Overview

This is a Go web framework middleware library built on top of Gin that provides high-level integrations for web services. It's part of the ternarybob ecosystem alongside:
- **arbor** - Structured logging system (../arbor)
- **satus** - Configuration and status management (../satus) 
- **funktion** - Core utility functions

### Core Components

#### Configuration System
- Uses `satus.AppConfig` for centralized configuration via `config.yml`
- Supports environment-based scoping (DEV/PRD)
- Configuration includes service metadata, database connections, and logging settings

#### Middleware Stack
- **Correlation ID**: `SetCorrelationID()` - Automatic request tracking with UUID generation
- **Error Handling**: Comprehensive error response middleware
- **Headers**: Standard HTTP header management
- **Recovery**: Panic recovery with proper logging
- **Static Files**: Static content serving

#### Response System
- `IRenderService` interface for structured API responses
- `ApiResponse` struct with metadata (version, support, correlation ID, logs)
- Generic typed responses with `ApiTypedResponse[T]` and `ApiPagedResponse[T]`
- Environment-aware error details (stack traces in DEV only)

#### Logging Integration
- Dual logging system: `phuslu/log` and `arbor` integration
- Memory-based log correlation using correlation IDs
- Request-scoped log retrieval and inclusion in API responses
- Configurable log levels via satus configuration

### Key Interfaces
- `IRenderService` - Response rendering abstraction
- `ICorrelationService` - Request correlation management

### Development Notes
- Uses Go 1.24 with Gin web framework
- Local development uses module replacements for ternarybob dependencies
- Test files follow Go naming conventions (`*_test.go`)
- Environment-aware behavior based on `cfg.Service.Scope`