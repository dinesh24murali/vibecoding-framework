package auth

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/users"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthTestRouter(repo userRepository) *gin.Engine {
	tokenManager := NewTokenManager("access-secret", "refresh-secret", time.Hour, 24*time.Hour)
	service := NewService(repo, tokenManager)
	handler := NewHandler(service)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/auth/admin/login", handler.AdminLogin)
	return r
}

func TestAdminLoginValidationError(t *testing.T) {
	r := setupAuthTestRouter(&fakeUserRepo{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/admin/login", bytes.NewBufferString(`{"identifier":"admin"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnprocessableEntity)
	}
}

func TestAdminLoginInvalidCredentials(t *testing.T) {
	r := setupAuthTestRouter(&fakeUserRepo{user: nil})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/admin/login", bytes.NewBufferString(`{"identifier":"admin","password":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnauthorized)
	}
}

func TestAdminLoginInternalError(t *testing.T) {
	r := setupAuthTestRouter(&fakeUserRepo{err: errors.New("db down")})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/admin/login", bytes.NewBufferString(`{"identifier":"admin","password":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusInternalServerError)
	}
}

func TestAdminLoginSuccess(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	repo := &fakeUserRepo{user: &users.User{
		ID:           21,
		Name:         "admin",
		PhoneNumber:  "9000000001",
		PasswordHash: string(hash),
		Role:         "admin",
		Status:       "active",
		CreatedAt:    time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC),
	}}

	r := setupAuthTestRouter(repo)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/admin/login", bytes.NewBufferString(`{"identifier":"admin","password":"StrongPassword123!"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
}
