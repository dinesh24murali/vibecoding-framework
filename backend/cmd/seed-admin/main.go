package main

import (
	"context"
	"log"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/config"
	"github.com/dinesh/vibecoding-framework/backend/internal/db"
	"github.com/dinesh/vibecoding-framework/backend/internal/seed"
	"github.com/dinesh/vibecoding-framework/backend/internal/users"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gormDB, err := db.OpenGorm(cfg.DB)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("extract sql db: %v", err)
	}

	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Printf("close db error: %v", closeErr)
		}
	}()

	seeder := seed.NewAdminSeeder(users.NewRepository(gormDB))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	created, userID, err := seeder.Seed(ctx, seed.AdminSeedInput{
		Username: cfg.SeedAdmin.Username,
		Phone:    cfg.SeedAdmin.Phone,
		Password: cfg.SeedAdmin.Password,
	})
	if err != nil {
		log.Fatalf("seed admin failed: %v", err)
	}

	if created {
		log.Printf("seed admin created: id=%d phone=%s", userID, cfg.SeedAdmin.Phone)
		return
	}

	log.Printf("seed admin skipped: already exists for phone=%s", cfg.SeedAdmin.Phone)
}
