package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorHandlerUsesAPIErrorEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID())
	engine.Use(ErrorHandler())
	engine.GET("/boom", func(c *gin.Context) {
		AbortWithAPIError(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid input")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	req.Header.Set(RequestIDHeader, "req-err-1")
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusBadRequest)
	}

	body := resp.Body.String()
	if !strings.Contains(body, "\"code\":\"VALIDATION_ERROR\"") {
		t.Fatalf("body missing code: %s", body)
	}
	if !strings.Contains(body, "\"request_id\":\"req-err-1\"") {
		t.Fatalf("body missing request id: %s", body)
	}
}

func TestErrorHandlerFallsBackToInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(RequestID())
	engine.Use(ErrorHandler())
	engine.GET("/boom", func(c *gin.Context) {
		_ = c.Error(http.ErrBodyNotAllowed)
		c.Abort()
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusInternalServerError)
	}

	body := resp.Body.String()
	if !strings.Contains(body, "\"code\":\"INTERNAL_ERROR\"") {
		t.Fatalf("body missing internal error code: %s", body)
	}
}
