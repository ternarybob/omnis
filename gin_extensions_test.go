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

	t.Run("Basic Fluent Interface", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/test", func(c *gin.Context) {
			log := arbor.GetLogger().WithPrefix("TestHandler")

			// Test the fluent interface: omnis.C(c).WithLogger(log).JSON(200, data)
			C(c).WithLogger(log).JSON(http.StatusOK, gin.H{
				"message": "fluent interface",
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
		// Response is wrapped in APIResponse format, actual data is in result field
		result, ok := response["result"].(map[string]interface{})
		assert.True(t, ok, "result should be a map")
		assert.Equal(t, "fluent interface", result["message"])
		assert.Equal(t, "success", result["status"])
	})

	t.Run("Convenience Methods", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/success", func(c *gin.Context) {
			log := arbor.GetLogger().WithPrefix("SuccessHandler")
			C(c).WithLogger(log).Success(gin.H{"result": "ok"})
		})

		r.GET("/created", func(c *gin.Context) {
			log := arbor.GetLogger().WithPrefix("CreatedHandler")
			C(c).WithLogger(log).Created(gin.H{"id": 123})
		})

		r.GET("/bad-request", func(c *gin.Context) {
			log := arbor.GetLogger().WithPrefix("BadRequestHandler")
			C(c).WithLogger(log).BadRequest(gin.H{"error": "invalid input"})
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

	t.Run("Without Logger", func(t *testing.T) {
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
			C(c).JSON(http.StatusOK, gin.H{"message": "no logger"})
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
		assert.Equal(t, "no logger", result["message"])
	})

	t.Run("Enhanced Response", func(t *testing.T) {
		r := gin.New()
		r.Use(SetCorrelationID())
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/test", func(c *gin.Context) {
			log := arbor.GetLogger().WithPrefix("EnhancedHandler")

			// Test enhanced response (full omnis wrapper)
			C(c).WithLogger(log).Enhanced(http.StatusOK, gin.H{
				"message": "enhanced response",
			})
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

func TestGinExtensionsChaining(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Method Chaining", func(t *testing.T) {
		r := gin.New()
		r.Use(JSONMiddlewareWithDefaults())

		r.GET("/test", func(c *gin.Context) {
			log := arbor.GetLogger().WithPrefix("ChainHandler")

			// Test that chaining works properly
			ext := C(c)
			extWithLogger := ext.WithLogger(log)

			// Should be the same instance for fluent chaining
			assert.Equal(t, ext, extWithLogger)

			extWithLogger.JSON(http.StatusOK, gin.H{"chained": true})
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
		assert.Equal(t, true, result["chained"])
	})
}
