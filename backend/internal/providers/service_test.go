package providers

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&Provider{}); err != nil {
		t.Fatalf("auto migrate provider: %v", err)
	}

	repo := NewRepository(db)
	return NewService(repo), db
}

func TestServiceCreateAndDuplicate(t *testing.T) {
	service, _ := newTestService(t)

	created, err := service.Create(context.Background(), CreateInput{Name: "  Airtel DTH  "})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.Name != "Airtel DTH" {
		t.Fatalf("name = %q, want %q", created.Name, "Airtel DTH")
	}

	_, err = service.Create(context.Background(), CreateInput{Name: "airtel dth"})
	if err != ErrProviderDuplicate {
		t.Fatalf("Create() duplicate error = %v, want %v", err, ErrProviderDuplicate)
	}
}

func TestServiceCreateValidation(t *testing.T) {
	service, _ := newTestService(t)

	_, err := service.Create(context.Background(), CreateInput{Name: "   "})
	if err != ErrProviderValidation {
		t.Fatalf("Create() validation error = %v, want %v", err, ErrProviderValidation)
	}
}

func TestServiceUpdateAndDelete(t *testing.T) {
	service, _ := newTestService(t)

	p, err := service.Create(context.Background(), CreateInput{Name: "Sun Direct"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	newName := "Sun Direct HD"
	updated, err := service.Update(context.Background(), p.ID, UpdateInput{Name: &newName})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updated.Name != "Sun Direct HD" {
		t.Fatalf("updated name = %q", updated.Name)
	}

	if err := service.Delete(context.Background(), p.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = service.GetByID(context.Background(), p.ID)
	if err != ErrProviderNotFound {
		t.Fatalf("GetByID() after delete error = %v, want %v", err, ErrProviderNotFound)
	}
}
