package auth

import (
	"testing"
	"time"
)

func TestTokenManagerGenerateAndVerifyAccessToken(t *testing.T) {
	manager := NewTokenManager("access-secret", "refresh-secret", time.Hour, 24*time.Hour)

	access, refresh, accessTTL, refreshTTL, err := manager.GenerateAdminTokens(42, "admin")
	if err != nil {
		t.Fatalf("GenerateAdminTokens() error = %v", err)
	}

	if access == "" {
		t.Fatalf("access token is empty")
	}

	if refresh == "" {
		t.Fatalf("refresh token is empty")
	}

	if accessTTL != int64(time.Hour.Seconds()) {
		t.Fatalf("accessTTL = %d, want %d", accessTTL, int64(time.Hour.Seconds()))
	}

	if refreshTTL != int64((24 * time.Hour).Seconds()) {
		t.Fatalf("refreshTTL = %d, want %d", refreshTTL, int64((24 * time.Hour).Seconds()))
	}

	claims, err := manager.VerifyAccessToken(access)
	if err != nil {
		t.Fatalf("VerifyAccessToken() error = %v", err)
	}

	if claims.RegisteredClaims.Subject != "42" {
		t.Fatalf("subject = %q, want %q", claims.RegisteredClaims.Subject, "42")
	}

	if claims.Role != "admin" {
		t.Fatalf("role = %q, want %q", claims.Role, "admin")
	}
}

func TestTokenManagerVerifyAccessTokenFailsForWrongSecret(t *testing.T) {
	issuer := NewTokenManager("access-secret", "refresh-secret", time.Hour, 24*time.Hour)
	access, _, _, _, err := issuer.GenerateAdminTokens(7, "admin")
	if err != nil {
		t.Fatalf("GenerateAdminTokens() error = %v", err)
	}

	verifier := NewTokenManager("wrong-secret", "refresh-secret", time.Hour, 24*time.Hour)
	_, err = verifier.VerifyAccessToken(access)
	if err == nil {
		t.Fatalf("VerifyAccessToken() error = nil, want non-nil")
	}
}
