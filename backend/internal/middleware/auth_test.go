package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdminAuthUnauthorizedWithoutBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin/ping", AdminAuth(func(_ string) (string, string, error) {
		return "", "", errors.New("no token")
	}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/ping", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnauthorized)
	}
}

func TestAdminAuthSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin/ping", AdminAuth(func(token string) (string, string, error) {
		if token != "valid-token" {
			return "", "", errors.New("invalid")
		}

		return "123", "admin", nil
	}), func(c *gin.Context) {
		if c.GetString("admin_user_id") != "123" {
			c.Status(http.StatusInternalServerError)
			return
		}

		if c.GetString("admin_role") != "admin" {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/ping", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
}
