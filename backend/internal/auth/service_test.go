package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	user *users.User
	err  error
}

func (f *fakeUserRepo) FindAdminByUsername(_ context.Context, _ string) (*users.User, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.user, nil
}

func TestServiceLoginAdminSuccess(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("StrongPassword123!"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	repo := &fakeUserRepo{user: &users.User{
		ID:           11,
		Name:         "admin",
		PhoneNumber:  "9000000001",
		PasswordHash: string(hash),
		Role:         "admin",
		Status:       "active",
		CreatedAt:    time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC),
	}}

	service := NewService(repo, NewTokenManager("access-secret", "refresh-secret", time.Hour, 24*time.Hour))
	result, err := service.LoginAdmin(context.Background(), LoginInput{
		Identifier: "admin",
		Password:   "StrongPassword123!",
	})
	if err != nil {
		t.Fatalf("LoginAdmin() error = %v", err)
	}

	if result.TokenType != "Bearer" {
		t.Fatalf("TokenType = %q, want %q", result.TokenType, "Bearer")
	}

	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatalf("tokens should not be empty")
	}

	if result.User.ID != "11" {
		t.Fatalf("user id = %q, want %q", result.User.ID, "11")
	}
}

func TestServiceLoginAdminInvalidCredentials(t *testing.T) {
	service := NewService(&fakeUserRepo{user: nil}, NewTokenManager("access-secret", "refresh-secret", time.Hour, 24*time.Hour))

	_, err := service.LoginAdmin(context.Background(), LoginInput{Identifier: "admin", Password: "bad"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("LoginAdmin() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestServiceLoginAdminRepoError(t *testing.T) {
	service := NewService(&fakeUserRepo{err: errors.New("db down")}, NewTokenManager("access-secret", "refresh-secret", time.Hour, 24*time.Hour))

	_, err := service.LoginAdmin(context.Background(), LoginInput{Identifier: "admin", Password: "x"})
	if err == nil {
		t.Fatalf("LoginAdmin() error = nil, want non-nil")
	}
}
