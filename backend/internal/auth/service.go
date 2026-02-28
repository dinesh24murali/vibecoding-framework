package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dinesh/vibecoding-framework/backend/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	FindAdminByUsername(ctx context.Context, username string) (*users.User, error)
}

type Service struct {
	users  userRepository
	tokens *TokenManager
}

type LoginInput struct {
	Identifier string
	Password   string
}

type LoginResult struct {
	AccessToken      string
	TokenType        string
	ExpiresIn        int64
	RefreshToken     string
	RefreshExpiresIn int64
	User             LoginUser
}

type LoginUser struct {
	ID          string
	Name        string
	PhoneNumber string
	Role        string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

func NewService(usersRepo userRepository, tokenManager *TokenManager) *Service {
	return &Service{users: usersRepo, tokens: tokenManager}
}

func (s *Service) LoginAdmin(ctx context.Context, input LoginInput) (*LoginResult, error) {
	user, err := s.users.FindAdminByUsername(ctx, input.Identifier)
	if err != nil {
		return nil, fmt.Errorf("fetch admin user: %w", err)
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	access, refresh, accessTTL, refreshTTL, err := s.tokens.GenerateAdminTokens(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &LoginResult{
		AccessToken:      access,
		TokenType:        "Bearer",
		ExpiresIn:        accessTTL,
		RefreshToken:     refresh,
		RefreshExpiresIn: refreshTTL,
		User: LoginUser{
			ID:          strconv.FormatInt(user.ID, 10),
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
			Status:      user.Status,
			CreatedAt:   user.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   user.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}
