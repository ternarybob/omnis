package omnis

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type gincorrelation struct {
	ctx gin.Context
}

func SetCorrelationID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if correlation ID already exists in context
		correlationID := ctx.GetString(CORRELATION_ID_KEY)
		
		// If not in context, check for X-Correlation-ID header
		if correlationID == "" {
			correlationID = ctx.GetHeader("X-Correlation-ID")
		}
		
		// If still empty, generate a new UUID
		if correlationID == "" {
			uuidValue, err := uuid.NewRandom()
			if err != nil {
				// Fallback to a simple UUID string if generation fails
				correlationID = uuid.New().String()
			} else {
				correlationID = uuidValue.String()
			}
		}
		
		// Set correlation ID in context
		ctx.Set(CORRELATION_ID_KEY, correlationID)
		
		// Set correlation ID in response headers (both formats for compatibility)
		ctx.Header("X-Correlation-ID", correlationID)
		ctx.Header(CORRELATION_ID_KEY, correlationID)
		
		// Continue to next middleware
		ctx.Next()
	}
}

// GetCorrelationID retrieves the correlation ID from the gin context
// Returns the correlation ID or "unknown" if not found
func GetCorrelationID(c *gin.Context) string {
	if c == nil {
		return "unknown"
	}
	
	correlationID := c.GetString(CORRELATION_ID_KEY)
	if correlationID != "" {
		return correlationID
	}
	
	// Fallback: check headers
	correlationID = c.GetHeader("X-Correlation-ID")
	if correlationID != "" {
		return correlationID
	}
	
	correlationID = c.GetHeader(CORRELATION_ID_KEY)
	if correlationID != "" {
		return correlationID
	}
	
	return "unknown"
}

// GetCorrelationIDOrGenerate retrieves the correlation ID or generates one if not found
func GetCorrelationIDOrGenerate(c *gin.Context) string {
	correlationID := GetCorrelationID(c)
	if correlationID == "unknown" {
		// Generate a new correlation ID
		uuidValue, err := uuid.NewRandom()
		if err != nil {
			// Fallback to simple UUID generation
			correlationID = uuid.New().String()
		} else {
			correlationID = uuidValue.String()
		}
		
		// Set it in context and headers if context is available
		if c != nil {
			c.Set(CORRELATION_ID_KEY, correlationID)
			c.Header("X-Correlation-ID", correlationID)
			c.Header(CORRELATION_ID_KEY, correlationID)
		}
	}
	return correlationID
}
