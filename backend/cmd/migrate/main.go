package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/config"
	"github.com/dinesh/vibecoding-framework/backend/internal/db"
	"github.com/dinesh/vibecoding-framework/backend/internal/migrate"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: go run ./cmd/migrate up")
	}

	command := os.Args[1]
	if command != "up" {
		log.Fatalf("unsupported migrate command %q: only 'up' is supported", command)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	database, err := db.Open(cfg.DB)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			log.Printf("close db error: %v", closeErr)
		}
	}()

	runner := migrate.NewRunner(database, cfg.MigrationsDir)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := runner.Up(ctx); err != nil {
		log.Fatalf("migrate up failed: %v", err)
	}

	fmt.Println("migrate up completed")
}
