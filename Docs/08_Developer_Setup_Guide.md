# Developer Setup Guide - Village DTH Platform

## Metadata
- **Version:** 1.0
- **Date:** 2026-02-27
- **Owner:** DevEx / Solo Developer
- **Applies To:** `backend` (Go + Gin + GORM), `frontend/customer` (Next.js SSG), `frontend/admin` (Next.js SSG), Redis/Asynq, PostgreSQL
- **Source of Truth:** TAD v1.0, OpenAPI v1.0, Implementation Plan v1.0
- **Environment Target:** Local development with production-like behavior where practical

---

## Prerequisites (exact versions/tools)

Install the following before starting:

- `git` >= 2.40
- `go` = 1.26.x
- `node` = 22.x LTS
- `npm` = 10.x (or `pnpm` 9.x if repo standardizes; examples below use npm)
- `postgresql` = 16.x (server + `psql` client)
- `redis` = 7.x
- `make` >= 4.x
- `curl` >= 8.x
- `jq` >= 1.7
- `openssl` >= 3.x
- Optional but recommended:
  - `docker` + `docker compose` (for local DB/Redis)
  - `golangci-lint` >= 1.60
  - `swagger-cli` / `redocly` for OpenAPI validation

Verify:
```bash
go version
node -v
npm -v
psql --version
redis-server --version
make --version
```

---

## Repository Structure Overview

Expected structure (aligned with implementation plan):

```text
backend/
  cmd/
    api/
    worker/
    migrate/
    seed-admin/
  internal/
    auth/
    providers/
    plans/
    checkout/
    payments/
    servicerequests/
    notifications/
    csv/
    middleware/
    router/
    config/
    errors/
  db/
    migrations/
  api/
    openapi/
  tests/
    contract/

frontend/
  customer/
  admin/

infra/
  scripts/
  systemd/
  nginx/
  deploy/
  runbooks/
  release/
```

---

## Local Setup

### 1) Clone/install

```bash
git clone <repo-url>
cd <repo-folder>

# backend deps
cd backend && go mod tidy

# frontend deps
cd ../frontend/customer && npm install
cd ../admin && npm install
```

### 2) Env vars (`.env.example` contract)

Create these files:
- `backend/.env`
- `frontend/customer/.env.local`
- `frontend/admin/.env.local`

#### `backend/.env.example`
```env
# App
APP_ENV=local
API_PORT=8080
LOG_LEVEL=debug
API_BASE_PATH=/api/v1

# Database
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=dth_app
DB_USER=dth_user
DB_PASSWORD=dth_pass
DB_SSLMODE=disable

# Redis / Asynq
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0

# Auth
JWT_ACCESS_SECRET=replace_with_strong_secret_at_least_32_chars
JWT_REFRESH_SECRET=replace_with_strong_secret_at_least_32_chars
JWT_ACCESS_TTL=240h        # 10 days
JWT_REFRESH_TTL=2160h      # 90 days

# Seed admin
SEED_ADMIN_USERNAME=admin
SEED_ADMIN_PHONE=9000000001
SEED_ADMIN_PASSWORD=ChangeThisImmediately123!

# Razorpay
RAZORPAY_KEY_ID=rzp_test_xxx
RAZORPAY_KEY_SECRET=xxx
RAZORPAY_WEBHOOK_SECRET=xxx

# SMS provider
SMS_PROVIDER=your_provider
SMS_API_URL=https://sms.example.com/send
SMS_API_KEY=xxx
SMS_SENDER_ID=DTHAPP

# reCAPTCHA
RECAPTCHA_SITE_KEY=xxx
RECAPTCHA_SECRET_KEY=xxx
RECAPTCHA_MIN_SCORE=0.5

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

#### `frontend/customer/.env.local.example`
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_RECAPTCHA_SITE_KEY=xxx
```

#### `frontend/admin/.env.local.example`
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

### 3) DB setup (PostgreSQL)

```sql
CREATE ROLE dth_user WITH LOGIN PASSWORD 'dth_pass';
CREATE DATABASE dth_app OWNER dth_user;
GRANT ALL PRIVILEGES ON DATABASE dth_app TO dth_user;
```

Quick check:
```bash
psql "postgresql://dth_user:dth_pass@127.0.0.1:5432/dth_app?sslmode=disable" -c "select now();"
```

### 4) Redis setup

Run local Redis:
```bash
redis-server --port 6379
```
Health check:
```bash
redis-cli ping
# PONG
```

### 5) Migration commands

From `backend/`:

```bash
go run ./cmd/migrate up
go run ./cmd/migrate status
```

Rollback basics:
```bash
go run ./cmd/migrate down 1
go run ./cmd/migrate up
```

### 6) Seed command (admin user)

From `backend/`:

```bash
go run ./cmd/seed-admin
```

Expected:
- First run: admin created
- Subsequent run: no-op/idempotent skip

---

## Runbook Commands

### Backend run/test/lint

From `backend/`:

```bash
# Run API
go run ./cmd/api

# Run worker
go run ./cmd/worker

# Unit tests
go test ./... -v

# Race checks (recommended)
go test ./... -race

# Lint (if configured)
golangci-lint run

# Contract tests
go test ./tests/contract -v
```

### Frontend customer run/build

From `frontend/customer/`:

```bash
npm run dev       # local
npm run build     # production build
npm run start     # serve built app
```

### Frontend admin run/build

From `frontend/admin/`:

```bash
npm run dev
npm run build
npm run start
```

### Worker process (Asynq) run

From `backend/`:
```bash
go run ./cmd/worker
```

---

## Integration Setup

### Payment sandbox config (Razorpay)

- Required backend env:
  - `RAZORPAY_KEY_ID`
  - `RAZORPAY_KEY_SECRET`
  - `RAZORPAY_WEBHOOK_SECRET`
- Configure webhook URL in Razorpay dashboard:
  - `http://localhost:8080/api/v1/payments/razorpay/callback` (via tunnel for remote callbacks)
- Verify callback signature is enforced in local logs.
- Use test cards/methods from Razorpay sandbox docs.

### SMS provider config

- Required backend env:
  - `SMS_PROVIDER`
  - `SMS_API_URL`
  - `SMS_API_KEY`
  - `SMS_SENDER_ID`
- Trigger behavior:
  - SMS only when service request transitions to `completed`.
- Confirm retry behavior:
  - 3 attempts with backoff.

### reCAPTCHA v3 config

- Required:
  - `RECAPTCHA_SITE_KEY` (frontend customer)
  - `RECAPTCHA_SECRET_KEY` (backend)
  - `RECAPTCHA_MIN_SCORE=0.5`
- Failure UX string must be exactly:
  - `Please try after some time`

---

## API Contract Workflow

### OpenAPI validation

From `backend/` (example tools):
```bash
npx @redocly/cli lint api/openapi/openapi.yaml
# or
npx swagger-cli validate api/openapi/openapi.yaml
```

### Codegen (optional)

If used:
```bash
npx @openapitools/openapi-generator-cli generate \
  -i api/openapi/openapi.yaml \
  -g typescript-fetch \
  -o ../frontend/shared/api-client
```

### Contract testing

- Validate implemented endpoints against:
  - operation IDs
  - required auth
  - request/response schema shape
  - status code matrix
- Mandatory coverage:
  - `adminLogin`
  - `createServiceRequestAndPayment`
  - `handleRazorpayCallback`
  - `updateServiceRequestStatus`
  - `uploadProvidersCsv`
  - `uploadPlansCsv`

---

## Troubleshooting Guide (top 10)

1) **API fails to start (`connection refused` DB)**
- Check `DB_*` env values and PostgreSQL process.
- Run connectivity check with `psql`.

2) **Redis connection errors in worker**
- Ensure Redis is running on `REDIS_ADDR`.
- Verify password/db index.

3) **`401` on admin APIs after login**
- Confirm `Authorization: Bearer <access_token>`.
- Check token expiry and `JWT_ACCESS_SECRET` consistency.

4) **reCAPTCHA always fails**
- Ensure site/secret key pair belong to same project.
- Verify domain config includes localhost.
- Confirm `RECAPTCHA_MIN_SCORE=0.5`.

5) **Razorpay callback not reaching local**
- Use ngrok/cloudflared for tunnel.
- Update webhook URL in dashboard.
- Confirm route `/api/v1/payments/razorpay/callback`.

6) **CSV upload returns many row errors**
- Validate headers and field formats.
- Check provider existence before plan import.
- Ensure unique constraints respected.

7) **Max plans/provider enforcement surprises**
- Count active + inactive non-deleted plans before insert.
- App-layer enforcement is expected behavior.

8) **Migration fails midway**
- Inspect migration status table.
- Fix SQL and rerun from failed point.
- Restore from backup if irreversible partial state.

9) **Seed-admin does not create user**
- Verify env values and DB permissions.
- Check idempotency logic if user already exists.

10) **No SMS sent on completion**
- Ensure status transition actually reached `completed`.
- Check worker logs and SMS credentials.
- Verify queue consumer is running.

---

## Production-like Local Checks before merge

- Run full backend tests:
  - `go test ./... -v`
- Run contract tests:
  - `go test ./tests/contract -v`
- Validate OpenAPI:
  - lint/validate passes
- Fresh DB test:
  - drop/recreate DB -> migrate up -> seed-admin -> smoke login
- Checkout flow test:
  - provider -> plan -> payment initiation with reCAPTCHA token
- Status transition test:
  - `payment_success -> completed` triggers SMS job
- CSV test:
  - valid + mixed-invalid file cases for providers/plans
- Build all apps:
  - backend binary build + both frontend production builds
- Backup/restore dry run:
  - run backup script + restore to temp DB
- Verify logs:
  - structured error envelope + request_id correlation present

---

## Security and Secrets Handling Guidelines

- Never commit `.env`, secrets, tokens, or webhook secrets.
- Use strong random values for JWT secrets (>=32 chars).
- Seed admin credentials are bootstrap-only; change immediately post-seed.
- Keep access to env files restricted (`chmod 600`).
- Store prod secrets in systemd `EnvironmentFile` (e.g., `/etc/dth-app/env`) with least privilege.
- Enforce HTTPS in production via Nginx.
- Redact secrets from logs; only log request IDs and non-sensitive metadata.
- Keep Razorpay and SMS keys separate for local/test/prod.

---

## Traceability Matrix

| Setup Step | Architecture section | Implementation task_id | OpenAPI relevance |
|---|---|---|---|
| Install tools and bootstrap backend/frontend | Design Principles, Containers | T001 | all operation execution |
| OpenAPI lint/validation | API Architecture | T002, T022 | all `operationId`s + schemas |
| DB create + migrations | Data Architecture | T003, T004 | entity schemas (`User`,`Provider`,`Plan`,`ServiceRequest`) |
| Seed admin | Seed strategy | T005 | `adminLogin`, `AdminLoginRequest/Response` |
| JWT auth local verification | Security Architecture | T006 | `BearerAuth`, admin operations |
| Providers/Plans CRUD checks | System Components/API | T007, T008 | `createProvider`,`createPlan`, list/update/delete ops |
| Checkout + reCAPTCHA + payment create | Payment + reCAPTCHA | T009, T010 | `createServiceRequestAndPayment`, `PaymentRequest/Response` |
| Razorpay callback handling | Payment integration | T011 | `handleRazorpayCallback` |
| Service-request state transition checks | State model/API | T012 | `updateServiceRequestStatus`, `ServiceRequestStatusUpdateRequest` |
| SMS worker and retry validation | Async Processing | T013 | status update flow tied to `updateServiceRequestStatus` |
| CSV provider/plan imports | CSV import pipeline | T014, T015 | `uploadProvidersCsv`,`uploadPlansCsv`,`CsvUploadResult` |
| Frontend customer/admin run | Customer/Admin Frontends | T016, T017 | customer + admin operation consumption |
| Error envelope/logging checks | Error envelope + Operations | T018 | `ErrorResponse` contract consistency |
| Backup/restore scripts | Reliability & Operations | T019 | persistence guarantees for API-backed data |
| systemd + Nginx local simulation | Deployment Topology | T020 | `/api/v1` routing + auth paths |
| Unit and contract gate before merge | Testing Strategy | T021, T022 | schema and operation compliance |
| Release dry run and rollback drill | Migration/Release Ops | T023 | prod endpoint readiness |

---

## Assumptions
- Repository layout follows the implementation plan blueprint.
- Local machine can run PostgreSQL and Redis concurrently.
- Razorpay/SMS sandbox credentials are available for integration tests.
- Contract tests exist or will be added under `backend/tests/contract`.
- Refresh token endpoint remains intentionally absent in v1.

## Open Questions
- For async CSV future behavior, define exact threshold and whether response should switch to `202 Accepted`.
- Confirm whether admin username format has character restrictions (for client-side validation).
- Confirm if local HTTPS is required for reCAPTCHA depending on domain policy.

## Risks
- Local setup may pass while webhook callbacks fail without tunnel/network setup.
- Long-lived admin tokens increase risk if local storage handling is weak.
- App-layer max-plan checks need transaction-safe implementation to avoid race conditions.
- Single-node assumptions can hide deployment-time SPOF issues.
- Manual backup/restore not exercised regularly can fail when needed.

## Decision Log

| Decision ID | Date | Decision | Options Considered | Chosen Option | Rationale | Impact |
|---|---|---|---|---|---|---|
| DX-001 | 2026-02-27 | Local stack strategy | Native-only vs Docker-only vs hybrid | Hybrid supported, native-first docs | Least friction for solo dev | Faster onboarding |
| DX-002 | 2026-02-27 | Auth refresh in setup | Include now vs omit | Omit refresh endpoint usage | Matches approved scope | Simpler local flow |
| DX-003 | 2026-02-27 | CSV local behavior | Async now vs sync now | Sync now, async-ready notes | Matches implementation plan | Clear v1 behavior |
| DX-004 | 2026-02-27 | Secrets management local | Plain env in docs vs managed secret tooling | `.env` local + strict hygiene | Practical for solo setup | Requires discipline |
| DX-005 | 2026-02-27 | Verification gate | Basic run only vs full pre-merge checks | Full production-like local checks | Reduce late-stage failures | More upfront effort |
