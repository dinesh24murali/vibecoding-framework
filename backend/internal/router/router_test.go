package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine, err := New("/api/v1", "production", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("root healthz", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
		}

		if got := resp.Header().Get("X-Request-ID"); got == "" {
			t.Fatalf("X-Request-ID header is empty")
		}
	})

	t.Run("api health", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
		}

		if got := resp.Header().Get("X-Request-ID"); got == "" {
			t.Fatalf("X-Request-ID header is empty")
		}
	})
}
