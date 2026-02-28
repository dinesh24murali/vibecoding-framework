package notifications

import (
	"context"
	"errors"
	"testing"

	"github.com/dinesh/vibecoding-framework/backend/internal/queue"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeSMSClient struct {
	err error
}

func (f fakeSMSClient) Send(_ context.Context, _ string, _ string) error {
	return f.err
}

func newWorkerDeps(t *testing.T) (*Worker, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.Exec(`CREATE TABLE sms_notifications (id INTEGER PRIMARY KEY AUTOINCREMENT, service_request_id INTEGER NOT NULL, phone_number TEXT NOT NULL, template TEXT NOT NULL, status TEXT NOT NULL, attempt_count INTEGER NOT NULL DEFAULT 0, last_error TEXT, sent_at DATETIME, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP);`).Error; err != nil {
		t.Fatalf("create sms_notifications: %v", err)
	}

	q := queue.NewAsynq(db)
	w := NewWorker(q, fakeSMSClient{})
	return w, db
}

func TestWorkerMarksSentOnSuccessfulDispatch(t *testing.T) {
	w, db := newWorkerDeps(t)
	w.client = fakeSMSClient{err: nil}

	if err := db.Exec(`INSERT INTO sms_notifications (service_request_id, phone_number, template, status, attempt_count) VALUES (1, '9000000002', 'msg', 'queued', 0)`).Error; err != nil {
		t.Fatalf("seed sms_notifications: %v", err)
	}

	if err := w.processOnce(context.Background()); err != nil {
		t.Fatalf("processOnce() error = %v", err)
	}

	var status string
	if err := db.Table("sms_notifications").Select("status").Where("id = 1").Take(&status).Error; err != nil {
		t.Fatalf("load notification: %v", err)
	}
	if status != "sent" {
		t.Fatalf("status = %q, want sent", status)
	}
}

func TestWorkerMarksFailedOnDispatchError(t *testing.T) {
	w, db := newWorkerDeps(t)
	w.client = fakeSMSClient{err: errors.New("provider unavailable")}

	if err := db.Exec(`INSERT INTO sms_notifications (service_request_id, phone_number, template, status, attempt_count) VALUES (1, '9000000002', 'msg', 'queued', 0)`).Error; err != nil {
		t.Fatalf("seed sms_notifications: %v", err)
	}

	if err := w.processOnce(context.Background()); err != nil {
		t.Fatalf("processOnce() error = %v", err)
	}

	var row struct {
		Status       string `gorm:"column:status"`
		AttemptCount int    `gorm:"column:attempt_count"`
	}
	if err := db.Table("sms_notifications").Select("status, attempt_count").Where("id = 1").Take(&row).Error; err != nil {
		t.Fatalf("load notification: %v", err)
	}
	if row.Status != "failed" {
		t.Fatalf("status = %q, want failed", row.Status)
	}
	if row.AttemptCount != 1 {
		t.Fatalf("attempt_count = %d, want 1", row.AttemptCount)
	}
}
