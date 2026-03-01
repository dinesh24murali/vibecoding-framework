package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoggingMiddlewareIncludesRequestIDAndStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var logs bytes.Buffer

	engine := gin.New()
	engine.Use(RequestID())
	engine.Use(LoggingWithWriter(&logs))
	engine.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req.Header.Set(RequestIDHeader, "req-123")
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusNoContent)
	}

	line := logs.String()
	if !strings.Contains(line, "request_id=req-123") {
		t.Fatalf("log line missing request id: %q", line)
	}

	if !strings.Contains(line, "status=204") {
		t.Fatalf("log line missing status: %q", line)
	}
}
