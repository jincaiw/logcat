package handlers

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestUsesHTTPS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		req      *http.Request
		expected bool
	}{
		{
			name:     "plain http",
			req:      httptest.NewRequest(http.MethodGet, "http://example.test", nil),
			expected: false,
		},
		{
			name: "tls request",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "https://example.test", nil)
				req.TLS = &tls.ConnectionState{}
				return req
			}(),
			expected: true,
		},
		{
			name: "x forwarded proto",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "http://example.test", nil)
				req.Header.Set("X-Forwarded-Proto", "https")
				return req
			}(),
			expected: true,
		},
		{
			name: "forwarded header",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "http://example.test", nil)
				req.Header.Set("Forwarded", "for=192.0.2.10;proto=https;host=example.test")
				return req
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = tt.req
			if got := requestUsesHTTPS(c); got != tt.expected {
				t.Fatalf("requestUsesHTTPS() = %v, want %v", got, tt.expected)
			}
		})
	}
}
