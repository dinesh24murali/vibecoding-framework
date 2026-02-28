package checkout

import (
	"context"
	"errors"
	"testing"

	"github.com/dinesh/vibecoding-framework/backend/internal/captcha"
	"github.com/dinesh/vibecoding-framework/backend/internal/plans"
	"github.com/dinesh/vibecoding-framework/backend/internal/servicerequests"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeCaptchaVerifier struct {
	err error
}

type fakeOrderCreator struct{}

func (fakeOrderCreator) CreateOrder(_ context.Context, _ float64, _ string, _ string) (string, error) {
	return "order_test_1", nil
}

func (f fakeCaptchaVerifier) Verify(_ context.Context, _ string) error {
	return f.err
}

func newCheckoutService(t *testing.T, verifier captchaVerifier) *Service {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, phone_number TEXT NOT NULL UNIQUE, role TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'active', password_hash TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME);`).Error; err != nil {
		t.Fatalf("create users table: %v", err)
	}
	if err := db.Exec(`CREATE TABLE providers (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, deleted_at DATETIME);`).Error; err != nil {
		t.Fatalf("create providers table: %v", err)
	}
	if err := db.AutoMigrate(&plans.Plan{}, &servicerequests.ServiceRequest{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	if err := db.Exec(`CREATE TABLE payment_attempts (id INTEGER PRIMARY KEY AUTOINCREMENT, service_request_id INTEGER NOT NULL, gateway_order_id TEXT NOT NULL, amount REAL NOT NULL, status TEXT NOT NULL, idempotency_key TEXT NOT NULL);`).Error; err != nil {
		t.Fatalf("create payment_attempts: %v", err)
	}

	if err := db.Exec(`INSERT INTO providers(id, name) VALUES (1, 'Airtel')`).Error; err != nil {
		t.Fatalf("seed provider: %v", err)
	}
	if err := db.Exec(`INSERT INTO plans(id, provider_id, name, description, price, discount, is_active, created_at, updated_at) VALUES (1, 1, 'Monthly', 'Base', 300, 20, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`).Error; err != nil {
		t.Fatalf("seed plan: %v", err)
	}

	return NewService(db, plans.NewRepository(db), servicerequests.NewRepository(db), verifier, fakeOrderCreator{})
}

func TestCreateServiceRequestAndPaymentSuccess(t *testing.T) {
	service := newCheckoutService(t, fakeCaptchaVerifier{})

	result, err := service.CreateServiceRequestAndPayment(context.Background(), CreateCheckoutInput{
		ProviderID:     1,
		PlanID:         1,
		CustomerName:   "Kumar",
		CustomerPhone:  "9000000002",
		RecaptchaToken: "token-token",
		IdempotencyKey: "idem-1",
	})
	if err != nil {
		t.Fatalf("CreateServiceRequestAndPayment() error = %v", err)
	}

	if result.ServiceRequest.Status != "pending" {
		t.Fatalf("status = %q, want pending", result.ServiceRequest.Status)
	}

	if result.Payment.Gateway != "razorpay" {
		t.Fatalf("gateway = %q, want razorpay", result.Payment.Gateway)
	}
}

func TestCreateServiceRequestAndPaymentNotFound(t *testing.T) {
	service := newCheckoutService(t, fakeCaptchaVerifier{})

	_, err := service.CreateServiceRequestAndPayment(context.Background(), CreateCheckoutInput{
		ProviderID:     2,
		PlanID:         1,
		CustomerName:   "Kumar",
		CustomerPhone:  "9000000002",
		RecaptchaToken: "token-token",
		IdempotencyKey: "idem-1",
	})
	if err != ErrNotFound {
		t.Fatalf("error = %v, want %v", err, ErrNotFound)
	}
}

func TestCreateServiceRequestAndPaymentCaptchaRejected(t *testing.T) {
	service := newCheckoutService(t, fakeCaptchaVerifier{err: captcha.ErrCaptchaFailed})

	_, err := service.CreateServiceRequestAndPayment(context.Background(), CreateCheckoutInput{
		ProviderID:     1,
		PlanID:         1,
		CustomerName:   "Kumar",
		CustomerPhone:  "9000000002",
		RecaptchaToken: "token-token",
		IdempotencyKey: "idem-1",
	})
	if !errors.Is(err, ErrCaptchaFailed) {
		t.Fatalf("error = %v, want %v", err, ErrCaptchaFailed)
	}
}
