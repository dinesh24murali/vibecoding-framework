package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("API_BASE_PATH", "")
	t.Setenv("OPENAPI_SPEC_PATH", "")
	t.Setenv("API_PORT", "")
	t.Setenv("WORKER_HEALTH_PORT", "")
	t.Setenv("SHUTDOWN_TIMEOUT", "")

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

	if cfg.APIPort != 8080 {
		t.Fatalf("APIPort = %d, want %d", cfg.APIPort, 8080)
	}

	if cfg.WorkerHealthPort != 8081 {
		t.Fatalf("WorkerHealthPort = %d, want %d", cfg.WorkerHealthPort, 8081)
	}

	if cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want %s", cfg.ShutdownTimeout, 10*time.Second)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("API_BASE_PATH", "/api/test")
	t.Setenv("OPENAPI_SPEC_PATH", "./api/openapi/test.yaml")
	t.Setenv("API_PORT", "9090")
	t.Setenv("WORKER_HEALTH_PORT", "9091")
	t.Setenv("SHUTDOWN_TIMEOUT", "30s")

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

	if cfg.APIPort != 9090 {
		t.Fatalf("APIPort = %d, want %d", cfg.APIPort, 9090)
	}

	if cfg.WorkerHealthPort != 9091 {
		t.Fatalf("WorkerHealthPort = %d, want %d", cfg.WorkerHealthPort, 9091)
	}

	if cfg.ShutdownTimeout != 30*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want %s", cfg.ShutdownTimeout, 30*time.Second)
	}
}
