package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine, err := New("/api/v1", "production", "", []string{"http://localhost:3000", "http://localhost:3001"}, Dependencies{})
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

	t.Run("no route uses standard error envelope", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/unknown", nil)
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusNotFound)
		}

		body := resp.Body.String()
		if !strings.Contains(body, "\"error\":") {
			t.Fatalf("body missing error envelope: %s", body)
		}
		if !strings.Contains(body, "\"code\":\"RESOURCE_NOT_FOUND\"") {
			t.Fatalf("body missing not found code: %s", body)
		}
		if !strings.Contains(body, "\"request_id\":") {
			t.Fatalf("body missing request id: %s", body)
		}
	})

	t.Run("method not allowed uses standard error envelope", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/health", nil)
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusMethodNotAllowed)
		}

		body := resp.Body.String()
		if !strings.Contains(body, "\"error\":") {
			t.Fatalf("body missing error envelope: %s", body)
		}
		if !strings.Contains(body, "\"code\":\"METHOD_NOT_ALLOWED\"") {
			t.Fatalf("body missing method code: %s", body)
		}
		if !strings.Contains(body, "\"request_id\":") {
			t.Fatalf("body missing request id: %s", body)
		}
	})
}
