// -----------------------------------------------------------------------
// Gin Extensions Tests
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

func TestGinExtensions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Basic JSON Interception", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/test", func(c *gin.Context) {
			// Standard c.JSON call - automatically intercepted
			c.JSON(http.StatusOK, gin.H{
				"message": "intercepted",
				"status":  "success",
			})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Without a logger, middleware passes through raw JSON
		assert.Equal(t, "intercepted", response["message"])
		assert.Equal(t, "success", response["status"])
	})

	t.Run("Standard JSON Responses", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/success", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"result": "ok"})
		})

		r.GET("/created", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"id": 123})
		})

		r.GET("/bad-request", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		})

		// Test Success (200)
		req, _ := http.NewRequest("GET", "/success", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Test Created (201)
		req, _ = http.NewRequest("GET", "/created", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Test Bad Request (400)
		req, _ = http.NewRequest("GET", "/bad-request", nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Middleware with Default Logger", func(t *testing.T) {
		r := gin.New()
		config := &ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
			Scope:   "DEV",
		}

		defaultLogger := arbor.GetLogger().WithPrefix("DefaultLogger")
		r.Use(JSONMiddlewareWithConfig(&JSONRendererConfig{
			ServiceConfig:     config,
			DefaultLogger:     defaultLogger,
			EnablePrettyPrint: true,
		}))

		r.GET("/test", func(c *gin.Context) {
			// No logger set - should use default from middleware
			c.JSON(http.StatusOK, gin.H{"message": "with default logger"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Response is wrapped in APIResponse format, actual data is in result field
		result, ok := response["result"].(map[string]interface{})
		assert.True(t, ok, "result should be a map")
		assert.Equal(t, "with default logger", result["message"])
	})

	t.Run("Enhanced Response Structure", func(t *testing.T) {
		r := gin.New()
		r.Use(SetCorrelationID())
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/test", func(c *gin.Context) {
			// Standard JSON response
			c.JSON(http.StatusOK, gin.H{
				"message": "enhanced response",
			})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Without a logger, middleware passes through raw JSON
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "enhanced response", response["message"])
	})
}

func TestJSONMiddlewareBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Response Wrapping", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/test", func(c *gin.Context) {
			// Standard JSON response
			c.JSON(http.StatusOK, gin.H{"data": "test"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Without a logger, middleware passes through raw JSON
		assert.Equal(t, "test", response["data"])
	})
}
