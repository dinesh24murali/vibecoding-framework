package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(CORS([]string{"http://localhost:3000", "http://localhost:3001"}))
	engine.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("sets CORS headers for allowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
		}
		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
			t.Fatalf("allow origin = %q, want %q", got, "http://localhost:3000")
		}
	})

	t.Run("does not set CORS headers for disallowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("Origin", "http://evil.local")
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
		}
		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Fatalf("allow origin = %q, want empty", got)
		}
	})

	t.Run("preflight returns no content", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/ok", nil)
		req.Header.Set("Origin", "http://localhost:3001")
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		if resp.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", resp.Code, http.StatusNoContent)
		}
	})
}
