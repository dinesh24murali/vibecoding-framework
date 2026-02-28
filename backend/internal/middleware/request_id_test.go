package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID())
	engine.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("generates request id when missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
		}

		if got := resp.Header().Get("X-Request-ID"); got == "" {
			t.Fatalf("X-Request-ID header is empty")
		}
	})

	t.Run("preserves incoming request id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("X-Request-ID", "test-request-id")
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
		}

		if got := resp.Header().Get("X-Request-ID"); got != "test-request-id" {
			t.Fatalf("X-Request-ID = %q, want %q", got, "test-request-id")
		}
	})
}

func TestRecoveryMiddlewareIncludesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID())
	engine.Use(Recovery())
	engine.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set("X-Request-ID", "panic-request-id")
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusInternalServerError)
	}

	body := resp.Body.String()
	if body == "" {
		t.Fatalf("response body is empty")
	}

	if want := "panic-request-id"; !strings.Contains(body, want) {
		t.Fatalf("response body does not contain request id %q", want)
	}
}
