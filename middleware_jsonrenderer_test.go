// -----------------------------------------------------------------------
// JSON Renderer Middleware Tests
// Created: 2025-08-27
// -----------------------------------------------------------------------

package omnis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/ternarybob/arbor"
)

func TestJSONRenderer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("JSON Interception with Standard c.JSON", func(t *testing.T) {
		r := gin.New()
		logger := arbor.GetLogger().WithPrefix("TestHandler")
		config := &ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
			Scope:   "DEV",
		}

		r.Use(JSONMiddlewareWithConfig(&JSONRendererConfig{
			ServiceConfig:     config,
			DefaultLogger:     logger,
			EnablePrettyPrint: true,
		}))

		r.GET("/test", func(c *gin.Context) {
			// Standard Gin c.JSON call - should be intercepted and enhanced
			c.JSON(http.StatusOK, gin.H{"message": "Intercepted", "status": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Should be pretty-printed JSON due to development mode
		body := w.Body.String()
		assert.Contains(t, body, "Intercepted")
		assert.Contains(t, body, "success")
		// Check that it's pretty printed (contains newlines and spaces)
		assert.Contains(t, body, "\n")
	})

	t.Run("JSON Interception with WithLogger", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/test", func(c *gin.Context) {
			log := arbor.GetLogger().WithPrefix("TestHandler")
			// Set logger in context, then use standard c.JSON
			WithLogger(c, log)
			c.JSON(http.StatusOK, gin.H{"message": "With Logger", "count": 42})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "With Logger", response["message"])
		assert.Equal(t, float64(42), response["count"]) // JSON unmarshals numbers as float64
	})

	t.Run("Basic JSON Response", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/test", func(c *gin.Context) {
			// Standard c.JSON call - automatically intercepted
			c.JSON(http.StatusOK, gin.H{"message": "Hello World", "status": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Hello World", response["message"])
		assert.Equal(t, "success", response["status"])
	})

	t.Run("JSON Response with Logger", func(t *testing.T) {
		r := gin.New()
		r.Use(SetCorrelationID())

		logger := arbor.GetLogger().WithPrefix("TestHandler")
		config := &ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
			Scope:   "TEST",
		}

		r.Use(JSONMiddlewareWithConfig(&JSONRendererConfig{
			ServiceConfig:     config,
			DefaultLogger:     logger,
			EnablePrettyPrint: true,
		}))

		r.GET("/test", func(c *gin.Context) {
			data := gin.H{"message": "Test with logger"}
			c.JSON(http.StatusOK, data)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Convenience Methods", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/success", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"result": "ok"})
		})

		r.GET("/bad-request", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		})

		r.GET("/not-found", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		})

		// Test Success
		req, _ := http.NewRequest("GET", "/success", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Test BadRequest
		req, _ = http.NewRequest("GET", "/bad-request", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test NotFound
		req, _ = http.NewRequest("GET", "/not-found", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Integration with RenderService", func(t *testing.T) {
		r := gin.New()
		r.Use(SetCorrelationID())

		logger := arbor.GetLogger().WithPrefix("TestHandler")
		config := &ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
			Scope:   "TEST",
		}

		r.Use(JSONMiddleware(config))

		r.GET("/test", func(c *gin.Context) {
			data := gin.H{"message": "Integration test"}
			// Use the render service directly for full omnis response
			RenderService(c).WithLogger(logger).WithConfig(config).AsResult(http.StatusOK, data)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Should contain the omnis response structure
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "result")
		assert.Contains(t, response, "name")
		assert.Contains(t, response, "version")
	})
}

func TestJSONRendererWithoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Direct JSON Usage", func(t *testing.T) {
		r := gin.New()

		r.GET("/test", func(c *gin.Context) {
			data := gin.H{"message": "Direct usage"}
			c.JSON(http.StatusOK, data)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Direct usage", response["message"])
	})
}
