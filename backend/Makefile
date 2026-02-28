GO ?= go

.PHONY: help run-api run-worker migrate-up seed-admin test test-v tidy fmt

help:
	@printf "Available targets:\n"
	@printf "  make run-api      # Run API server\n"
	@printf "  make run-worker   # Run worker health server\n"
	@printf "  make migrate-up   # Run DB migrations\n"
	@printf "  make seed-admin   # Seed default admin\n"
	@printf "  make test         # Run backend tests\n"
	@printf "  make test-v       # Run backend tests (verbose)\n"
	@printf "  make tidy         # Sync Go module dependencies\n"
	@printf "  make fmt          # Format Go code\n"

run-api:
	$(GO) run ./cmd/api

run-worker:
	$(GO) run ./cmd/worker

migrate-up:
	$(GO) run ./cmd/migrate up

seed-admin:
	$(GO) run ./cmd/seed-admin

test:
	$(GO) test ./...

test-v:
	$(GO) test ./... -v

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...
