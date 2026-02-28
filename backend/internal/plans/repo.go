package plans

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/common"
	"gorm.io/gorm"
)

type Repository struct {
	*common.GormRepository[Plan]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{GormRepository: common.NewGormRepository[Plan](db)}
}

func (r *Repository) ProviderExists(ctx context.Context, providerID int64) (bool, error) {
	var count int64
	err := r.Db.WithContext(ctx).Table("providers").Where("id = ? AND deleted_at IS NULL", providerID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check provider exists: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) ListAdmin(ctx context.Context, page int, pageSize int, providerID *int64, isActive *bool, sort string) ([]Plan, int64, error) {
	db := r.Db.WithContext(ctx).Model(&Plan{}).Where("deleted_at IS NULL")
	if providerID != nil {
		db = db.Where("provider_id = ?", *providerID)
	}
	if isActive != nil {
		db = db.Where("is_active = ?", *isActive)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count plans: %w", err)
	}

	var items []Plan
	if err := db.Order(normalizeSort(sort)).Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list plans: %w", err)
	}

	return items, total, nil
}

func (r *Repository) ListCustomerByProvider(ctx context.Context, providerID int64) ([]Plan, error) {
	var items []Plan
	err := r.Db.WithContext(ctx).
		Where("provider_id = ? AND is_active = TRUE AND deleted_at IS NULL", providerID).
		Order("name ASC").
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("list customer plans: %w", err)
	}

	return items, nil
}

func (r *Repository) FindActiveByID(ctx context.Context, id int64) (*Plan, error) {
	var plan Plan
	err := r.Db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("find plan by id: %w", err)
	}

	return &plan, nil
}

func (r *Repository) FindActiveByProviderAndName(ctx context.Context, providerID int64, name string) (*Plan, error) {
	var plan Plan
	err := r.Db.WithContext(ctx).
		Where("provider_id = ? AND LOWER(name) = LOWER(?) AND deleted_at IS NULL", providerID, strings.TrimSpace(name)).
		First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("find plan by provider and name: %w", err)
	}

	return &plan, nil
}

func (r *Repository) ExistsActiveByProviderAndName(ctx context.Context, providerID int64, name string, excludeID *int64) (bool, error) {
	query := r.Db.WithContext(ctx).Model(&Plan{}).
		Where("provider_id = ? AND LOWER(name) = LOWER(?) AND deleted_at IS NULL", providerID, strings.TrimSpace(name))
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check plan uniqueness: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) ActivePlanCountByProvider(ctx context.Context, providerID int64, excludeID *int64) (int64, error) {
	query := r.Db.WithContext(ctx).Model(&Plan{}).Where("provider_id = ? AND deleted_at IS NULL", providerID)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count plans by provider: %w", err)
	}

	return count, nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int64) (bool, error) {
	now := time.Now().UTC()
	result := r.Db.WithContext(ctx).Model(&Plan{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": &now, "updated_at": now})
	if result.Error != nil {
		return false, fmt.Errorf("soft delete plan: %w", result.Error)
	}

	return result.RowsAffected > 0, nil
}

func normalizeSort(sort string) string {
	switch sort {
	case "updated_at_asc":
		return "updated_at ASC"
	case "price_asc":
		return "price ASC"
	case "price_desc":
		return "price DESC"
	default:
		return "updated_at DESC"
	}
}
