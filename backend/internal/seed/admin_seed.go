package seed

import (
	"context"
	"fmt"

	"github.com/dinesh/vibecoding-framework/backend/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type AdminSeedInput struct {
	Username string
	Phone    string
	Password string
}

type adminRepository interface {
	AdminExistsByPhone(ctx context.Context, phoneNumber string) (bool, error)
	CreateAdmin(ctx context.Context, input users.CreateAdminInput) (int64, error)
}

type AdminSeeder struct {
	repo adminRepository
}

func NewAdminSeeder(repo adminRepository) *AdminSeeder {
	return &AdminSeeder{repo: repo}
}

func (s *AdminSeeder) Seed(ctx context.Context, input AdminSeedInput) (created bool, userID int64, err error) {
	if input.Username == "" {
		return false, 0, fmt.Errorf("seed username is required")
	}

	if input.Phone == "" {
		return false, 0, fmt.Errorf("seed phone is required")
	}

	if input.Password == "" {
		return false, 0, fmt.Errorf("seed password is required")
	}

	exists, err := s.repo.AdminExistsByPhone(ctx, input.Phone)
	if err != nil {
		return false, 0, err
	}

	if exists {
		return false, 0, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, 0, fmt.Errorf("hash password: %w", err)
	}

	userID, err = s.repo.CreateAdmin(ctx, users.CreateAdminInput{
		Username:     input.Username,
		PhoneNumber:  input.Phone,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return false, 0, err
	}

	return true, userID, nil
}
