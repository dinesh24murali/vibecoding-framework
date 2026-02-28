package payments

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/dinesh/vibecoding-framework/backend/internal/captcha"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeCaptchaVerifier struct {
	err error
}

func (f fakeCaptchaVerifier) Verify(_ context.Context, _ string) error {
	return f.err
}

func newPaymentsService(t *testing.T, captchaErr error) (*Service, *gorm.DB, *RazorpayClient) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, phone_number TEXT NOT NULL);`).Error; err != nil {
		t.Fatalf("create users: %v", err)
	}
	if err := db.Exec(`CREATE TABLE plans (id INTEGER PRIMARY KEY AUTOINCREMENT, price REAL NOT NULL, discount REAL NOT NULL);`).Error; err != nil {
		t.Fatalf("create plans: %v", err)
	}
	if err := db.Exec(`CREATE TABLE service_requests (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, plan_id INTEGER NOT NULL, status TEXT NOT NULL, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME);`).Error; err != nil {
		t.Fatalf("create service_requests: %v", err)
	}
	if err := db.Exec(`CREATE TABLE payment_attempts (id INTEGER PRIMARY KEY AUTOINCREMENT, service_request_id INTEGER NOT NULL, gateway_order_id TEXT NOT NULL, gateway_payment_id TEXT, amount REAL NOT NULL, status TEXT NOT NULL, idempotency_key TEXT NOT NULL, callback_payload TEXT);`).Error; err != nil {
		t.Fatalf("create payment_attempts: %v", err)
	}
	if err := db.Exec(`CREATE TABLE service_request_retry_audit (id INTEGER PRIMARY KEY AUTOINCREMENT, service_request_id INTEGER NOT NULL, retry_no INTEGER NOT NULL, reason TEXT NOT NULL);`).Error; err != nil {
		t.Fatalf("create retry_audit: %v", err)
	}

	if err := db.Exec(`INSERT INTO users (id, phone_number) VALUES (1, '9000000002')`).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if err := db.Exec(`INSERT INTO plans (id, price, discount) VALUES (1, 300, 20)`).Error; err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	if err := db.Exec(`INSERT INTO service_requests (id, user_id, plan_id, status, created_at, updated_at, deleted_at) VALUES (1, 1, 1, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL)`).Error; err != nil {
		t.Fatalf("seed service request: %v", err)
	}

	client := NewRazorpayClient("rzp_test_x", "secret", "webhook_secret")
	svc := NewService(db, client, NewRetryAuditRepo(db), fakeCaptchaVerifier{err: captchaErr})
	return svc, db, client
}

func TestRetryPaymentSuccess(t *testing.T) {
	svc, db, _ := newPaymentsService(t, nil)

	result, err := svc.RetryPayment(context.Background(), 1, "token", "idem-1")
	if err != nil {
		t.Fatalf("RetryPayment() error = %v", err)
	}

	if result.Payment["status"] != "created" {
		t.Fatalf("payment status = %v, want created", result.Payment["status"])
	}

	var count int64
	if err := db.Table("service_request_retry_audit").Where("service_request_id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count retry audit: %v", err)
	}
	if count != 1 {
		t.Fatalf("retry audit count = %d, want 1", count)
	}
}

func TestRetryPaymentCaptchaRejected(t *testing.T) {
	svc, _, _ := newPaymentsService(t, captcha.ErrCaptchaFailed)

	_, err := svc.RetryPayment(context.Background(), 1, "token", "idem-1")
	if !errors.Is(err, ErrCaptchaRejected) {
		t.Fatalf("RetryPayment() error = %v, want %v", err, ErrCaptchaRejected)
	}
}

func TestHandleCallbackSuccessAndDuplicate(t *testing.T) {
	svc, db, client := newPaymentsService(t, nil)
	if err := db.Exec(`INSERT INTO payment_attempts (id, service_request_id, gateway_order_id, amount, status, idempotency_key) VALUES (1, 1, 'order_123', 280, 'created', 'idem-1')`).Error; err != nil {
		t.Fatalf("seed payment attempt: %v", err)
	}

	raw := []byte(`{"payload":{"payment":{"entity":{"order_id":"order_123","id":"pay_123","status":"captured"}}}}`)
	sig := sign(client.webhookSecret, raw)

	duplicate, err := svc.HandleCallback(context.Background(), raw, sig, map[string]any{
		"payload": map[string]any{
			"payment": map[string]any{
				"entity": map[string]any{
					"order_id": "order_123",
					"id":       "pay_123",
					"status":   "captured",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("HandleCallback() error = %v", err)
	}
	if duplicate {
		t.Fatalf("duplicate = true, want false")
	}

	duplicate, err = svc.HandleCallback(context.Background(), raw, sig, map[string]any{
		"payload": map[string]any{
			"payment": map[string]any{
				"entity": map[string]any{
					"order_id": "order_123",
					"id":       "pay_123",
					"status":   "captured",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("HandleCallback() duplicate error = %v", err)
	}
	if !duplicate {
		t.Fatalf("duplicate = false, want true")
	}
}

func TestHandleCallbackInvalidSignature(t *testing.T) {
	svc, _, _ := newPaymentsService(t, nil)
	_, err := svc.HandleCallback(context.Background(), []byte("{}"), "bad", map[string]any{})
	if !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("HandleCallback() error = %v, want %v", err, ErrInvalidSignature)
	}
}

func sign(secret string, raw []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}
