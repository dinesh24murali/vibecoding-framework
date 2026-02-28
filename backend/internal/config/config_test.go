package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("API_BASE_PATH", "")
	t.Setenv("OPENAPI_SPEC_PATH", "")
	t.Setenv("MIGRATIONS_DIR", "")
	t.Setenv("API_PORT", "")
	t.Setenv("WORKER_HEALTH_PORT", "")
	t.Setenv("SHUTDOWN_TIMEOUT", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_SSLMODE", "")
	t.Setenv("SEED_ADMIN_USERNAME", "")
	t.Setenv("SEED_ADMIN_PHONE", "")
	t.Setenv("SEED_ADMIN_PASSWORD", "")
	t.Setenv("JWT_ACCESS_SECRET", "")
	t.Setenv("JWT_REFRESH_SECRET", "")
	t.Setenv("JWT_ACCESS_TTL", "")
	t.Setenv("JWT_REFRESH_TTL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AppEnv != "local" {
		t.Fatalf("AppEnv = %q, want %q", cfg.AppEnv, "local")
	}

	if cfg.APIBasePath != "/api/v1" {
		t.Fatalf("APIBasePath = %q, want %q", cfg.APIBasePath, "/api/v1")
	}

	if cfg.OpenAPISpecPath != "api/openapi/openapi.yaml" {
		t.Fatalf("OpenAPISpecPath = %q, want %q", cfg.OpenAPISpecPath, "api/openapi/openapi.yaml")
	}

	if cfg.MigrationsDir != "db/migrations" {
		t.Fatalf("MigrationsDir = %q, want %q", cfg.MigrationsDir, "db/migrations")
	}

	if cfg.APIPort != 8080 {
		t.Fatalf("APIPort = %d, want %d", cfg.APIPort, 8080)
	}

	if cfg.WorkerHealthPort != 8081 {
		t.Fatalf("WorkerHealthPort = %d, want %d", cfg.WorkerHealthPort, 8081)
	}

	if cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want %s", cfg.ShutdownTimeout, 10*time.Second)
	}

	if cfg.DB.Host != "127.0.0.1" {
		t.Fatalf("DB.Host = %q, want %q", cfg.DB.Host, "127.0.0.1")
	}

	if cfg.DB.Port != 5432 {
		t.Fatalf("DB.Port = %d, want %d", cfg.DB.Port, 5432)
	}

	if cfg.DB.Name != "dth_app" {
		t.Fatalf("DB.Name = %q, want %q", cfg.DB.Name, "dth_app")
	}

	if cfg.DB.User != "dth_user" {
		t.Fatalf("DB.User = %q, want %q", cfg.DB.User, "dth_user")
	}

	if cfg.DB.Password != "dth_pass" {
		t.Fatalf("DB.Password = %q, want %q", cfg.DB.Password, "dth_pass")
	}

	if cfg.DB.SSLMode != "disable" {
		t.Fatalf("DB.SSLMode = %q, want %q", cfg.DB.SSLMode, "disable")
	}

	if cfg.SeedAdmin.Username != "admin" {
		t.Fatalf("SeedAdmin.Username = %q, want %q", cfg.SeedAdmin.Username, "admin")
	}

	if cfg.SeedAdmin.Phone != "9000000001" {
		t.Fatalf("SeedAdmin.Phone = %q, want %q", cfg.SeedAdmin.Phone, "9000000001")
	}

	if cfg.SeedAdmin.Password != "ChangeThisImmediately123!" {
		t.Fatalf("SeedAdmin.Password = %q, want %q", cfg.SeedAdmin.Password, "ChangeThisImmediately123!")
	}

	if cfg.JWT.AccessSecret != "replace_with_strong_secret_at_least_32_chars" {
		t.Fatalf("JWT.AccessSecret = %q, want default", cfg.JWT.AccessSecret)
	}

	if cfg.JWT.RefreshSecret != "replace_with_strong_secret_at_least_32_chars" {
		t.Fatalf("JWT.RefreshSecret = %q, want default", cfg.JWT.RefreshSecret)
	}

	if cfg.JWT.AccessTTL != 240*time.Hour {
		t.Fatalf("JWT.AccessTTL = %s, want %s", cfg.JWT.AccessTTL, 240*time.Hour)
	}

	if cfg.JWT.RefreshTTL != 2160*time.Hour {
		t.Fatalf("JWT.RefreshTTL = %s, want %s", cfg.JWT.RefreshTTL, 2160*time.Hour)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("API_BASE_PATH", "/api/test")
	t.Setenv("OPENAPI_SPEC_PATH", "./api/openapi/test.yaml")
	t.Setenv("MIGRATIONS_DIR", "./db/test-migrations")
	t.Setenv("API_PORT", "9090")
	t.Setenv("WORKER_HEALTH_PORT", "9091")
	t.Setenv("SHUTDOWN_TIMEOUT", "30s")
	t.Setenv("DB_HOST", "db.internal")
	t.Setenv("DB_PORT", "15432")
	t.Setenv("DB_NAME", "dth_prod")
	t.Setenv("DB_USER", "dth_admin")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("SEED_ADMIN_USERNAME", "super-admin")
	t.Setenv("SEED_ADMIN_PHONE", "9000000099")
	t.Setenv("SEED_ADMIN_PASSWORD", "UltraSecure123!")
	t.Setenv("JWT_ACCESS_SECRET", "access-secret")
	t.Setenv("JWT_REFRESH_SECRET", "refresh-secret")
	t.Setenv("JWT_ACCESS_TTL", "24h")
	t.Setenv("JWT_REFRESH_TTL", "720h")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AppEnv != "production" {
		t.Fatalf("AppEnv = %q, want %q", cfg.AppEnv, "production")
	}

	if cfg.APIBasePath != "/api/test" {
		t.Fatalf("APIBasePath = %q, want %q", cfg.APIBasePath, "/api/test")
	}

	if cfg.OpenAPISpecPath != "./api/openapi/test.yaml" {
		t.Fatalf("OpenAPISpecPath = %q, want %q", cfg.OpenAPISpecPath, "./api/openapi/test.yaml")
	}

	if cfg.MigrationsDir != "./db/test-migrations" {
		t.Fatalf("MigrationsDir = %q, want %q", cfg.MigrationsDir, "./db/test-migrations")
	}

	if cfg.APIPort != 9090 {
		t.Fatalf("APIPort = %d, want %d", cfg.APIPort, 9090)
	}

	if cfg.WorkerHealthPort != 9091 {
		t.Fatalf("WorkerHealthPort = %d, want %d", cfg.WorkerHealthPort, 9091)
	}

	if cfg.ShutdownTimeout != 30*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want %s", cfg.ShutdownTimeout, 30*time.Second)
	}

	if cfg.DB.Host != "db.internal" {
		t.Fatalf("DB.Host = %q, want %q", cfg.DB.Host, "db.internal")
	}

	if cfg.DB.Port != 15432 {
		t.Fatalf("DB.Port = %d, want %d", cfg.DB.Port, 15432)
	}

	if cfg.DB.Name != "dth_prod" {
		t.Fatalf("DB.Name = %q, want %q", cfg.DB.Name, "dth_prod")
	}

	if cfg.DB.User != "dth_admin" {
		t.Fatalf("DB.User = %q, want %q", cfg.DB.User, "dth_admin")
	}

	if cfg.DB.Password != "secret" {
		t.Fatalf("DB.Password = %q, want %q", cfg.DB.Password, "secret")
	}

	if cfg.DB.SSLMode != "require" {
		t.Fatalf("DB.SSLMode = %q, want %q", cfg.DB.SSLMode, "require")
	}

	if cfg.SeedAdmin.Username != "super-admin" {
		t.Fatalf("SeedAdmin.Username = %q, want %q", cfg.SeedAdmin.Username, "super-admin")
	}

	if cfg.SeedAdmin.Phone != "9000000099" {
		t.Fatalf("SeedAdmin.Phone = %q, want %q", cfg.SeedAdmin.Phone, "9000000099")
	}

	if cfg.SeedAdmin.Password != "UltraSecure123!" {
		t.Fatalf("SeedAdmin.Password = %q, want %q", cfg.SeedAdmin.Password, "UltraSecure123!")
	}

	if cfg.JWT.AccessSecret != "access-secret" {
		t.Fatalf("JWT.AccessSecret = %q, want %q", cfg.JWT.AccessSecret, "access-secret")
	}

	if cfg.JWT.RefreshSecret != "refresh-secret" {
		t.Fatalf("JWT.RefreshSecret = %q, want %q", cfg.JWT.RefreshSecret, "refresh-secret")
	}

	if cfg.JWT.AccessTTL != 24*time.Hour {
		t.Fatalf("JWT.AccessTTL = %s, want %s", cfg.JWT.AccessTTL, 24*time.Hour)
	}

	if cfg.JWT.RefreshTTL != 720*time.Hour {
		t.Fatalf("JWT.RefreshTTL = %s, want %s", cfg.JWT.RefreshTTL, 720*time.Hour)
	}
}

func TestDBConfigDSN(t *testing.T) {
	db := DBConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "dth_app",
		User:     "dth_user",
		Password: "dth_pass",
		SSLMode:  "disable",
	}

	got := db.DSN()
	want := "host=localhost port=5432 user=dth_user password=dth_pass dbname=dth_app sslmode=disable"
	if got != want {
		t.Fatalf("DSN() = %q, want %q", got, want)
	}
}

func TestLoadFromDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	envData := "APP_ENV=dev\nAPI_PORT=8181\nDB_HOST=env-db\n"
	if err := os.WriteFile(envPath, []byte(envData), 0o600); err != nil {
		t.Fatalf("write .env file: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current directory: %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir to temp dir: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AppEnv != "dev" {
		t.Fatalf("AppEnv = %q, want %q", cfg.AppEnv, "dev")
	}

	if cfg.APIPort != 8181 {
		t.Fatalf("APIPort = %d, want %d", cfg.APIPort, 8181)
	}

	if cfg.DB.Host != "env-db" {
		t.Fatalf("DB.Host = %q, want %q", cfg.DB.Host, "env-db")
	}
}
