package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/config"
	_ "github.com/lib/pq"
)

func Open(cfg config.DBConfig) (*sql.DB, error) {
	database, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	database.SetMaxOpenConns(20)
	database.SetMaxIdleConns(10)
	database.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.PingContext(ctx); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return database, nil
}
