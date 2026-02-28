package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string
	APIBasePath      string
	OpenAPISpecPath  string
	MigrationsDir    string
	APIPort          int
	WorkerHealthPort int
	ShutdownTimeout  time.Duration
	DB               DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

func Load() (Config, error) {
	if err := loadEnvFiles(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:           getEnv("APP_ENV", "local"),
		APIBasePath:      getEnv("API_BASE_PATH", "/api/v1"),
		OpenAPISpecPath:  getEnv("OPENAPI_SPEC_PATH", "api/openapi/openapi.yaml"),
		MigrationsDir:    getEnv("MIGRATIONS_DIR", "db/migrations"),
		APIPort:          getEnvInt("API_PORT", 8080),
		WorkerHealthPort: getEnvInt("WORKER_HEALTH_PORT", 8081),
		ShutdownTimeout:  getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnvInt("DB_PORT", 5432),
			Name:     getEnv("DB_NAME", "dth_app"),
			User:     getEnv("DB_USER", "dth_user"),
			Password: getEnv("DB_PASSWORD", "dth_pass"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	if cfg.APIPort <= 0 {
		return Config{}, fmt.Errorf("invalid API_PORT: %d", cfg.APIPort)
	}

	if cfg.WorkerHealthPort <= 0 {
		return Config{}, fmt.Errorf("invalid WORKER_HEALTH_PORT: %d", cfg.WorkerHealthPort)
	}

	if cfg.OpenAPISpecPath == "" {
		return Config{}, fmt.Errorf("invalid OPENAPI_SPEC_PATH: empty")
	}

	if cfg.MigrationsDir == "" {
		return Config{}, fmt.Errorf("invalid MIGRATIONS_DIR: empty")
	}

	if cfg.DB.Host == "" {
		return Config{}, fmt.Errorf("invalid DB_HOST: empty")
	}

	if cfg.DB.Port <= 0 {
		return Config{}, fmt.Errorf("invalid DB_PORT: %d", cfg.DB.Port)
	}

	if cfg.DB.Name == "" {
		return Config{}, fmt.Errorf("invalid DB_NAME: empty")
	}

	if cfg.DB.User == "" {
		return Config{}, fmt.Errorf("invalid DB_USER: empty")
	}

	if cfg.DB.SSLMode == "" {
		return Config{}, fmt.Errorf("invalid DB_SSLMODE: empty")
	}

	if cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %s", cfg.ShutdownTimeout)
	}

	return cfg, nil
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvInt(key string, defaultValue int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return defaultValue
	}

	return value
}

func loadEnvFiles() error {
	paths := []string{".env", filepath.Join("backend", ".env")}

	for _, path := range paths {
		if err := godotenv.Load(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return fmt.Errorf("load env file %q: %w", path, err)
		}

		return nil
	}

	return nil
}
