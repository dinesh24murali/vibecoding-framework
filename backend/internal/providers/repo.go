package providers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/common"
	"gorm.io/gorm"
)

type Repository struct {
	*common.GormRepository[Provider]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{GormRepository: common.NewGormRepository[Provider](db)}
}

func (r *Repository) ListAdmin(ctx context.Context, page int, pageSize int, search string, sort string) ([]Provider, int64, error) {
	db := r.Db.WithContext(ctx).Model(&Provider{}).Where("deleted_at IS NULL")
	if strings.TrimSpace(search) != "" {
		db = db.Where("name ILIKE ?", "%"+strings.TrimSpace(search)+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count providers: %w", err)
	}

	var items []Provider
	if err := db.Order(normalizeSort(sort)).Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list providers: %w", err)
	}

	return items, total, nil
}

func (r *Repository) ListCustomerActive(ctx context.Context) ([]Provider, error) {
	var items []Provider
	err := r.Db.WithContext(ctx).Where("deleted_at IS NULL").Order("name ASC").Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("list customer providers: %w", err)
	}

	return items, nil
}

func (r *Repository) FindActiveByID(ctx context.Context, id int64) (*Provider, error) {
	var provider Provider
	err := r.Db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&provider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("find provider by id: %w", err)
	}

	return &provider, nil
}

func (r *Repository) FindActiveByName(ctx context.Context, name string) (*Provider, error) {
	var provider Provider
	err := r.Db.WithContext(ctx).Where("LOWER(name) = LOWER(?) AND deleted_at IS NULL", strings.TrimSpace(name)).First(&provider).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("find provider by name: %w", err)
	}

	return &provider, nil
}

func (r *Repository) ExistsActiveByName(ctx context.Context, name string, excludeID *int64) (bool, error) {
	query := r.Db.WithContext(ctx).Model(&Provider{}).Where("LOWER(name) = LOWER(?) AND deleted_at IS NULL", strings.TrimSpace(name))
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check provider name uniqueness: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) SoftDelete(ctx context.Context, id int64) (bool, error) {
	now := time.Now().UTC()
	result := r.Db.WithContext(ctx).Model(&Provider{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": &now, "updated_at": now})

	if result.Error != nil {
		return false, fmt.Errorf("soft delete provider: %w", result.Error)
	}

	return result.RowsAffected > 0, nil
}

func normalizeSort(sort string) string {
	switch sort {
	case "updated_at_asc":
		return "updated_at ASC"
	case "name_asc":
		return "name ASC"
	case "name_desc":
		return "name DESC"
	default:
		return "updated_at DESC"
	}
}
