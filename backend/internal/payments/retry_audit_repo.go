package payments

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type RetryAuditRepo struct {
	db *gorm.DB
}

func NewRetryAuditRepo(db *gorm.DB) *RetryAuditRepo {
	return &RetryAuditRepo{db: db}
}

func (r *RetryAuditRepo) NextRetryNo(ctx context.Context, serviceRequestID int64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("service_request_retry_audit").Where("service_request_id = ?", serviceRequestID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count retry audits: %w", err)
	}

	return int(count) + 1, nil
}

func (r *RetryAuditRepo) Insert(ctx context.Context, serviceRequestID int64, retryNo int, reason string) error {
	payload := map[string]any{
		"service_request_id": serviceRequestID,
		"retry_no":           retryNo,
		"reason":             reason,
	}
	if err := r.db.WithContext(ctx).Table("service_request_retry_audit").Create(payload).Error; err != nil {
		return fmt.Errorf("insert retry audit: %w", err)
	}

	return nil
}
