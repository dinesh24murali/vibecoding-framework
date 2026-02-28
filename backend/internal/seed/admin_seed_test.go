package seed

import (
	"context"
	"testing"

	"github.com/dinesh/vibecoding-framework/backend/internal/users"
)

type fakeAdminRepo struct {
	existsValue bool
	createdID   int64
	createCalls int
}

func (f *fakeAdminRepo) AdminExistsByPhone(_ context.Context, _ string) (bool, error) {
	return f.existsValue, nil
}

func (f *fakeAdminRepo) CreateAdmin(_ context.Context, _ users.CreateAdminInput) (int64, error) {
	f.createCalls++
	return f.createdID, nil
}

func TestAdminSeederSeedSkipsWhenExists(t *testing.T) {
	repo := &fakeAdminRepo{existsValue: true, createdID: 9}
	seeder := NewAdminSeeder(repo)

	created, id, err := seeder.Seed(context.Background(), AdminSeedInput{
		Username: "admin",
		Phone:    "9000000001",
		Password: "StrongPassword123!",
	})
	if err != nil {
		t.Fatalf("Seed() error = %v", err)
	}

	if created {
		t.Fatalf("created = true, want false")
	}

	if id != 0 {
		t.Fatalf("id = %d, want 0", id)
	}

	if repo.createCalls != 0 {
		t.Fatalf("createCalls = %d, want 0", repo.createCalls)
	}
}

func TestAdminSeederSeedCreatesWhenMissing(t *testing.T) {
	repo := &fakeAdminRepo{existsValue: false, createdID: 11}
	seeder := NewAdminSeeder(repo)

	created, id, err := seeder.Seed(context.Background(), AdminSeedInput{
		Username: "admin",
		Phone:    "9000000001",
		Password: "StrongPassword123!",
	})
	if err != nil {
		t.Fatalf("Seed() error = %v", err)
	}

	if !created {
		t.Fatalf("created = false, want true")
	}

	if id != 11 {
		t.Fatalf("id = %d, want 11", id)
	}

	if repo.createCalls != 1 {
		t.Fatalf("createCalls = %d, want 1", repo.createCalls)
	}
}
