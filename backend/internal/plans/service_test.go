package plans

import (
	"context"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type providerSeed struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
}

func (providerSeed) TableName() string { return "providers" }

func newPlanTestService(t *testing.T) *Service {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.Exec(`CREATE TABLE providers (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, deleted_at DATETIME);`).Error; err != nil {
		t.Fatalf("create providers table: %v", err)
	}

	if err := db.AutoMigrate(&Plan{}); err != nil {
		t.Fatalf("auto migrate plans: %v", err)
	}

	for i := 1; i <= 2; i++ {
		if err := db.Exec("INSERT INTO providers (id, name) VALUES (?, ?)", i, fmt.Sprintf("Provider %d", i)).Error; err != nil {
			t.Fatalf("seed provider %d: %v", i, err)
		}
	}

	return NewService(NewRepository(db))
}

func TestPlanCreateAndDuplicate(t *testing.T) {
	service := newPlanTestService(t)

	created, err := service.Create(context.Background(), CreateInput{
		ProviderID:  1,
		Name:        "Starter",
		Description: "Base plan",
		Price:       299,
		Discount:    10,
		IsActive:    true,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID == 0 {
		t.Fatalf("created ID should be set")
	}

	_, err = service.Create(context.Background(), CreateInput{
		ProviderID:  1,
		Name:        "starter",
		Description: "Duplicate",
		Price:       199,
		Discount:    0,
		IsActive:    true,
	})
	if err != ErrPlanDuplicate {
		t.Fatalf("duplicate error = %v, want %v", err, ErrPlanDuplicate)
	}
}

func TestPlanProviderLimit(t *testing.T) {
	service := newPlanTestService(t)
	for i := 0; i < 100; i++ {
		_, err := service.Create(context.Background(), CreateInput{
			ProviderID:  1,
			Name:        fmt.Sprintf("P-%d", i),
			Description: "D",
			Price:       100,
			Discount:    0,
			IsActive:    true,
		})
		if err != nil {
			t.Fatalf("seed create %d error: %v", i, err)
		}
	}

	_, err := service.Create(context.Background(), CreateInput{
		ProviderID:  1,
		Name:        "overflow",
		Description: "D",
		Price:       100,
		Discount:    0,
		IsActive:    true,
	})
	if err != ErrPlanLimitExceeded {
		t.Fatalf("overflow error = %v, want %v", err, ErrPlanLimitExceeded)
	}
}

func TestPlanUpsertByProviderAndName(t *testing.T) {
	service := newPlanTestService(t)

	err := service.UpsertByProviderAndName(context.Background(), CreateInput{
		ProviderID:  1,
		Name:        "Family Pack",
		Description: "Initial",
		Price:       250,
		Discount:    20,
		IsActive:    true,
	})
	if err != nil {
		t.Fatalf("UpsertByProviderAndName() create error = %v", err)
	}

	err = service.UpsertByProviderAndName(context.Background(), CreateInput{
		ProviderID:  1,
		Name:        "Family Pack",
		Description: "Updated",
		Price:       275,
		Discount:    25,
		IsActive:    false,
	})
	if err != nil {
		t.Fatalf("UpsertByProviderAndName() update error = %v", err)
	}

	plans, _, err := service.ListAdmin(context.Background(), 1, 20, nil, nil, "")
	if err != nil {
		t.Fatalf("ListAdmin() error = %v", err)
	}
	if len(plans) != 1 {
		t.Fatalf("plan count = %d, want 1", len(plans))
	}
	if plans[0].Description != "Updated" || plans[0].Price != 275 || plans[0].Discount != 25 || plans[0].IsActive {
		t.Fatalf("unexpected upserted plan: %+v", plans[0])
	}
}
