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

	t.Run("Basic JSON Response", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())
		
		r.GET("/test", func(c *gin.Context) {
			data := gin.H{"message": "Hello World", "status": "success"}
			JSON(c).Simple(http.StatusOK, data)
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
			JSON(c).Simple(http.StatusOK, data)
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
			JSON(c).Success(gin.H{"result": "ok"})
		})

		r.GET("/bad-request", func(c *gin.Context) {
			JSON(c).BadRequest(gin.H{"error": "invalid input"})
		})

		r.GET("/not-found", func(c *gin.Context) {
			JSON(c).NotFound(gin.H{"error": "resource not found"})
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
			JSON(c).WithLogger(logger).Response(http.StatusOK, data)
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
			JSON(c).Simple(http.StatusOK, data)
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