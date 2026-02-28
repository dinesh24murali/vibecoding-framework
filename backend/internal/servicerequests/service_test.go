package servicerequests

import (
	"context"
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newServiceRequestsService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, phone_number TEXT NOT NULL, role TEXT NOT NULL, status TEXT NOT NULL, created_at DATETIME, updated_at DATETIME);`).Error; err != nil {
		t.Fatalf("create users: %v", err)
	}
	if err := db.Exec(`CREATE TABLE providers (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, image_url TEXT, created_at DATETIME, updated_at DATETIME);`).Error; err != nil {
		t.Fatalf("create providers: %v", err)
	}
	if err := db.Exec(`CREATE TABLE plans (id INTEGER PRIMARY KEY AUTOINCREMENT, provider_id INTEGER NOT NULL, name TEXT NOT NULL, description TEXT NOT NULL, price REAL NOT NULL, discount REAL NOT NULL, is_active BOOLEAN NOT NULL, created_at DATETIME, updated_at DATETIME);`).Error; err != nil {
		t.Fatalf("create plans: %v", err)
	}
	if err := db.Exec(`CREATE TABLE service_requests (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, plan_id INTEGER NOT NULL, status TEXT NOT NULL, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME);`).Error; err != nil {
		t.Fatalf("create service_requests: %v", err)
	}

	if err := db.Exec(`INSERT INTO users (id, name, phone_number, role, status, created_at, updated_at) VALUES (1, 'Asha', '9000000002', 'customer', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`).Error; err != nil {
		t.Fatalf("seed users: %v", err)
	}
	if err := db.Exec(`INSERT INTO providers (id, name, image_url, created_at, updated_at) VALUES (1, 'Dish TV', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`).Error; err != nil {
		t.Fatalf("seed providers: %v", err)
	}
	if err := db.Exec(`INSERT INTO plans (id, provider_id, name, description, price, discount, is_active, created_at, updated_at) VALUES (1, 1, 'Family 199', 'Base plan', 199, 0, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`).Error; err != nil {
		t.Fatalf("seed plans: %v", err)
	}
	if err := db.Exec(`INSERT INTO service_requests (id, user_id, plan_id, status, created_at, updated_at, deleted_at) VALUES (1, 1, 1, 'payment_success', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL), (2, 1, 1, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL)`).Error; err != nil {
		t.Fatalf("seed service requests: %v", err)
	}

	return NewService(db), db
}

func TestUpdateStatusRejectsInvalidTransition(t *testing.T) {
	svc, _ := newServiceRequestsService(t)

	_, err := svc.UpdateStatus(context.Background(), 2, StatusCompleted)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("UpdateStatus() error = %v, want %v", err, ErrInvalidTransition)
	}
}

func TestUpdateStatusSuccess(t *testing.T) {
	svc, _ := newServiceRequestsService(t)

	updated, err := svc.UpdateStatus(context.Background(), 1, StatusCompleted)
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	if updated.Status != StatusCompleted {
		t.Fatalf("status = %q, want %q", updated.Status, StatusCompleted)
	}
}

func TestListAdminFiltersByStatus(t *testing.T) {
	svc, _ := newServiceRequestsService(t)

	items, total, err := svc.ListAdmin(context.Background(), ListAdminInput{Page: 1, PageSize: 10, Status: StatusPending})
	if err != nil {
		t.Fatalf("ListAdmin() error = %v", err)
	}
	if total != 1 {
		t.Fatalf("total = %d, want 1", total)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].Status != StatusPending {
		t.Fatalf("status = %q, want %q", items[0].Status, StatusPending)
	}
}
