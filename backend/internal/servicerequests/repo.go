package servicerequests

import (
	"context"
	"fmt"

	"github.com/dinesh/vibecoding-framework/backend/internal/common"
	"gorm.io/gorm"
)

type Repository struct {
	*common.GormRepository[ServiceRequest]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{GormRepository: common.NewGormRepository[ServiceRequest](db)}
}

func (r *Repository) CreatePending(ctx context.Context, userID int64, planID int64) (*ServiceRequest, error) {
	req := ServiceRequest{
		UserID: userID,
		PlanID: planID,
		Status: "pending",
	}

	if err := gorm.G[ServiceRequest](r.Db.WithContext(ctx)).Create(ctx, &req); err != nil {
		return nil, fmt.Errorf("create service request: %w", err)
	}

	return &req, nil
}
