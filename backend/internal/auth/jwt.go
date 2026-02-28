package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type AccessClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}

func NewTokenManager(accessSecret string, refreshSecret string, accessTTL time.Duration, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (m *TokenManager) GenerateAdminTokens(userID int64, role string) (string, string, int64, int64, error) {
	now := time.Now()
	accessExpiry := now.Add(m.accessTTL)
	refreshExpiry := now.Add(m.refreshTTL)

	subject := strconv.FormatInt(userID, 10)
	accessClaims := AccessClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
		},
	}

	refreshClaims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(m.accessSecret)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("sign access token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(m.refreshSecret)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("sign refresh token: %w", err)
	}

	return accessToken, refreshToken, int64(m.accessTTL.Seconds()), int64(m.refreshTTL.Seconds()), nil
}

func (m *TokenManager) VerifyAccessToken(token string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return m.accessSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse access token: %w", err)
	}

	if !parsed.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return claims, nil
}
