package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Capture log output
	var logOutput bytes.Buffer
	logrus.SetOutput(&logOutput)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Create a new Gin router and apply the middleware
	router := gin.New()
	router.Use(Logger())

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	logString := logOutput.String()
	assert.Contains(t, logString, `"method":"GET"`)
	assert.Contains(t, logString, `"path":"/test"`)
	assert.Contains(t, logString, `"status":200`)
	assert.Contains(t, logString, `"msg":"HTTP request"`)
}