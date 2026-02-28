package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv           string
	APIBasePath      string
	OpenAPISpecPath  string
	APIPort          int
	WorkerHealthPort int
	ShutdownTimeout  time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:           getEnv("APP_ENV", "local"),
		APIBasePath:      getEnv("API_BASE_PATH", "/api/v1"),
		OpenAPISpecPath:  getEnv("OPENAPI_SPEC_PATH", "api/openapi/openapi.yaml"),
		APIPort:          getEnvInt("API_PORT", 8080),
		WorkerHealthPort: getEnvInt("WORKER_HEALTH_PORT", 8081),
		ShutdownTimeout:  getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
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

	if cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %s", cfg.ShutdownTimeout)
	}

	return cfg, nil
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
