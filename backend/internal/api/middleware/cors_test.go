package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORS(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		origin         string
		expectedStatus int
		expectedHeader string
		isPreflight    bool
	}{
		{
			name:           "Simple Request - Allowed Origin",
			method:         "GET",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			expectedHeader: "http://localhost:3000",
			isPreflight:    false,
		},
		{
			name:           "Simple Request - Disallowed Origin",
			method:         "GET",
			origin:         "http://evil.com",
			expectedStatus: http.StatusForbidden,
			expectedHeader: "",
			isPreflight:    false,
		},
		{
			name:           "Preflight Request - Allowed Origin",
			method:         "OPTIONS",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusNoContent,
			expectedHeader: "http://localhost:3000",
			isPreflight:    true,
		},
		{
			name:           "Preflight Request - Disallowed Origin",
			method:         "OPTIONS",
			origin:         "http://evil.com",
			expectedStatus: http.StatusForbidden,
			expectedHeader: "",
			isPreflight:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()
			// Apply the CORS middleware
			router.Use(CORS())

			// Define a dummy route to handle the request
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})
			router.OPTIONS("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			req, _ := http.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			if tt.isPreflight {
				req.Header.Set("Access-Control-Request-Method", "GET")
				req.Header.Set("Access-Control-Request-Headers", "Content-Type")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			actualHeader := w.Header().Get("Access-Control-Allow-Origin")
			if actualHeader != tt.expectedHeader {
				t.Errorf("Expected Access-Control-Allow-Origin header %q, got %q", tt.expectedHeader, actualHeader)
			}
		})
	}
}
