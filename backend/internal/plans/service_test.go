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
