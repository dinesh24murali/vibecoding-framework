package queue

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const maxFailureAttempts = 4

type SMSNotification struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement"`
	ServiceRequestID int64      `gorm:"column:service_request_id"`
	PhoneNumber      string     `gorm:"column:phone_number"`
	Template         string     `gorm:"column:template"`
	Status           string     `gorm:"column:status"`
	AttemptCount     int        `gorm:"column:attempt_count"`
	LastError        *string    `gorm:"column:last_error"`
	SentAt           *time.Time `gorm:"column:sent_at"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
}

func (SMSNotification) TableName() string {
	return "sms_notifications"
}

type Asynq struct {
	db *gorm.DB
}

func NewAsynq(db *gorm.DB) *Asynq {
	return &Asynq{db: db}
}

func (q *Asynq) EnqueueCompleted(ctx context.Context, tx *gorm.DB, serviceRequestID int64, phoneNumber string, template string) error {
	record := SMSNotification{
		ServiceRequestID: serviceRequestID,
		PhoneNumber:      phoneNumber,
		Template:         template,
		Status:           "queued",
		AttemptCount:     0,
	}

	if err := tx.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("enqueue completed sms: %w", err)
	}

	return nil
}

func (q *Asynq) ListPendingSMS(ctx context.Context, limit int) ([]SMSNotification, error) {
	if limit <= 0 {
		limit = 25
	}

	var rows []SMSNotification
	err := q.db.WithContext(ctx).
		Where("status = ? OR (status = ? AND attempt_count < ?)", "queued", "failed", maxFailureAttempts).
		Order("created_at ASC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list pending sms notifications: %w", err)
	}

	return rows, nil
}

func (q *Asynq) MarkSent(ctx context.Context, id int64) error {
	now := time.Now().UTC()
	if err := q.db.WithContext(ctx).
		Model(&SMSNotification{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     "sent",
			"sent_at":    now,
			"updated_at": now,
			"last_error": nil,
		}).Error; err != nil {
		return fmt.Errorf("mark sms as sent: %w", err)
	}

	return nil
}

func (q *Asynq) MarkFailed(ctx context.Context, id int64, errText string) error {
	now := time.Now().UTC()
	if err := q.db.WithContext(ctx).
		Model(&SMSNotification{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":        "failed",
			"attempt_count": gorm.Expr("attempt_count + 1"),
			"last_error":    errText,
			"updated_at":    now,
		}).Error; err != nil {
		return fmt.Errorf("mark sms as failed: %w", err)
	}

	return nil
}
