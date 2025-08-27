// -----------------------------------------------------------------------
// Gin Extensions Tests
// Created: 2025-08-27
// -----------------------------------------------------------------------

package omnis

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/ternarybob/arbor"
)

func TestGinExtensions(t *testing.T) {
	gin.SetMode(gin.TestMode)

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
		// Just verify the response isn't empty
		assert.NotEmpty(t, w.Body.String())
	})
}
