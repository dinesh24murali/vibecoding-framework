package contract

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/dinesh/vibecoding-framework/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

const testSpec = `openapi: 3.0.3
info:
  title: Test API
  version: 1.0.0
servers:
  - url: /api/v1
paths:
  /auth/admin/login:
    post:
      operationId: adminLogin
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AdminLoginRequest'
      responses:
        '200':
          description: ok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AdminLoginResponse'
        '422':
          description: validation
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
components:
  schemas:
    AdminLoginRequest:
      type: object
      required: [identifier, password]
      properties:
        identifier:
          type: string
        password:
          type: string
    User:
      type: object
      required: [id, name, phone_number, role, status, created_at, updated_at]
      properties:
        id:
          type: string
        name:
          type: string
        phone_number:
          type: string
        role:
          type: string
        status:
          type: string
        created_at:
          type: string
        updated_at:
          type: string
    AdminLoginResponse:
      type: object
      required: [access_token, token_type, expires_in, refresh_token, refresh_expires_in, user]
      properties:
        access_token:
          type: string
        token_type:
          type: string
        expires_in:
          type: integer
        refresh_token:
          type: string
        refresh_expires_in:
          type: integer
        user:
          $ref: '#/components/schemas/User'
    ErrorDetail:
      type: object
      required: [field, issue]
      properties:
        field:
          type: string
        issue:
          type: string
    Error:
      type: object
      required: [code, message]
      properties:
        code:
          type: string
        message:
          type: string
        details:
          type: array
          items:
            $ref: '#/components/schemas/ErrorDetail'
    ErrorResponse:
      type: object
      required: [error]
      properties:
        error:
          $ref: '#/components/schemas/Error'
`

func writeSpecFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	specPath := filepath.Join(dir, "openapi.yaml")
	if err := os.WriteFile(specPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write spec file: %v", err)
	}

	return specPath
}

func TestNew_DisabledInProduction(t *testing.T) {
	validator, err := New("production", "")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if validator.enabled {
		t.Fatalf("validator.enabled = %v, want %v", validator.enabled, false)
	}
}

func TestNew_RequiresValidSpecInNonProd(t *testing.T) {
	_, err := New("local", "./does-not-exist.yaml")
	if err == nil {
		t.Fatalf("New() error = nil, want non-nil")
	}
}

func TestMiddleware_RejectsInvalidCriticalRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	specPath := writeSpecFile(t, testSpec)

	validator, err := New("local", specPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	engine := gin.New()
	engine.Use(middleware.RequestID())
	engine.Use(validator.Middleware())
	engine.POST("/api/v1/auth/admin/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	body := []byte(`{"identifier":"admin"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}

func TestMiddleware_AllowsValidCriticalRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	specPath := writeSpecFile(t, testSpec)

	validator, err := New("local", specPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	engine := gin.New()
	engine.Use(middleware.RequestID())
	engine.Use(validator.Middleware())
	engine.POST("/api/v1/auth/admin/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"access_token":       "token",
			"token_type":         "Bearer",
			"expires_in":         864000,
			"refresh_token":      "refresh",
			"refresh_expires_in": 7776000,
			"user": gin.H{
				"id":           "usr_1",
				"name":         "admin",
				"phone_number": "9000000001",
				"role":         "admin",
				"status":       "active",
				"created_at":   "2026-01-01T00:00:00Z",
				"updated_at":   "2026-01-01T00:00:00Z",
			},
		})
	})

	body := []byte(`{"identifier":"admin","password":"StrongPassword123!"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	engine.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}

	if got := resp.Header().Get("X-Contract-Response-Invalid"); got != "" {
		t.Fatalf("X-Contract-Response-Invalid = %q, want empty", got)
	}
}

func TestShouldEnable(t *testing.T) {
	if shouldEnable("production") {
		t.Fatalf("shouldEnable(production) = true, want false")
	}

	if shouldEnable("prod") {
		t.Fatalf("shouldEnable(prod) = true, want false")
	}

	if !shouldEnable("local") {
		t.Fatalf("shouldEnable(local) = false, want true")
	}
}
