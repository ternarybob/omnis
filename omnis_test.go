package omnis

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestDefaultLogger tests the default logger creation
func TestDefaultLogger(t *testing.T) {
	logger := defaultLogger()

	// Check that the logger is properly configured
	if logger.Writer == nil {
		t.Error("Expected default logger writer to be non-nil")
	}
	
	// Test that logger can be used without panicking
	logger.Info().Msg("Test default logger creation")
}

// TestWarnLogger tests the warning logger creation
func TestWarnLogger(t *testing.T) {
	logger := warnLogger()

	// Check that the logger is properly configured
	if logger.Writer == nil {
		t.Error("Expected warn logger writer to be non-nil")
	}

	// Test that logger can be used (we can't directly check level in current interface)
	// but we can test that it doesn't panic and works
	logger.Info().Msg("Test warn logger creation")
}

// TestSetCorrelationID tests setting a correlation ID middleware
func TestSetCorrelationID(t *testing.T) {
	// Setup a Gin engine with the middleware
	r := gin.New()
	r.Use(SetCorrelationID())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello")
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}

	correlationID := w.Header().Get(CORRELATION_ID_KEY)

	if correlationID == "" {
		t.Error("Expected correlation ID to be set in the header")
	}
}

// Test that constants are correctly defined
func TestConstants(t *testing.T) {
	expectedConstants := map[string]string{
		CORRELATION_ID_KEY: "correlationid",
	}

	for constant, expected := range expectedConstants {
		if constant != expected {
			t.Errorf("Expected constant %s to be %q, got %q", expected, expected, constant)
		}
	}
}

// TestConfig tests that config is initialized
func TestConfigInitialization(t *testing.T) {
	if cfg == nil {
		t.Error("Expected config to be initialized")
	}
}
