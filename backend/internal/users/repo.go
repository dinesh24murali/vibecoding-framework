package users

import (
	"context"
	"fmt"

	"github.com/dinesh/vibecoding-framework/backend/internal/common"
	"gorm.io/gorm"
)

type Repository struct {
	*common.GormRepository[User]
}

type CreateAdminInput struct {
	Username     string
	PhoneNumber  string
	PasswordHash string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		GormRepository: common.NewGormRepository[User](db),
	}
}

func (r *Repository) AdminExistsByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	count, err := gorm.G[User](r.Db.WithContext(ctx)).
		Where("phone_number = ? AND role = ?", phoneNumber, "admin").
		Count(ctx, "id")
	if err != nil {
		return false, fmt.Errorf("query admin by phone: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) CreateAdmin(ctx context.Context, input CreateAdminInput) (int64, error) {
	admin := User{
		Name:         input.Username,
		PhoneNumber:  input.PhoneNumber,
		PasswordHash: input.PasswordHash,
		Role:         "admin",
		Status:       "active",
	}

	err := gorm.G[User](r.Db.WithContext(ctx)).Create(ctx, &admin)
	if err != nil {
		return 0, fmt.Errorf("insert admin user: %w", err)
	}

	return admin.ID, nil
}

func (r *Repository) FindAdminByUsername(ctx context.Context, username string) (*User, error) {
	user, err := gorm.G[User](r.Db.WithContext(ctx)).
		Where("name = ? AND role = ?", username, "admin").
		Take(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("query admin by username: %w", err)
	}

	return &user, nil
}
