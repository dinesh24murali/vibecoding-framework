# Technical Architecture Document (TAD) - Village DTH Plan Purchase Platform

## Metadata
- **Version:** 1.0
- **Date:** 2026-02-27
- **Owner:** Solutions Architecture
- **Linked PRD:** v1.0
- **Linked FSD:** v1.0
- **Status:** Approved-for-Build

## Architecture Overview (C4 level-1/2 textual)

### C4 Level 1 - System Context
The platform consists of:
- A **Customer Web App** for provider/plan selection and payment initiation.
- An **Admin Web App** for provider/plan CRUD, CSV imports, and service request status operations.
- A **Backend API** exposing REST endpoints under `/api/v1`.
- A **PostgreSQL database** as system of record.
- **Redis + Asynq worker** for asynchronous jobs (SMS retries, reconciliation jobs).
- External services: **Razorpay** (payments), **Google reCAPTCHA v3**, **SMS provider**.

Primary users are village customers (guest checkout) and a single operator admin.

### C4 Level 2 - Containers and Responsibilities
1. **Customer Frontend (Next.js SSG export)**
   - Static shell + runtime API calls for dynamic data.
   - Handles provider/plan browsing, checkout details, reCAPTCHA token generation.
2. **Admin Frontend (Next.js SSG export)**
   - Static shell + JWT-based API access.
   - Handles CRUD, CSV upload, and service request status transitions.
3. **API Service (Go 1.26 + Gin + GORM)**
   - Auth, catalog, checkout, payment, request lifecycle, CSV processing.
   - Enforces business rules, validation, status transition rules.
4. **Worker Service (Asynq consumer)**
   - SMS dispatch on completion, retries, and deferred payment reconciliation.
5. **PostgreSQL**
   - Authoritative store for users/providers/plans/service requests/payment and audit tables.
6. **Redis**
   - Asynq queue backend, retry scheduling, DLQ state.
7. **External Integrations**
   - Razorpay APIs/webhooks, Google reCAPTCHA verify API, SMS gateway API.

---

## Design Principles and Constraints
- **Low-cost first:** Single production environment and minimal infrastructure complexity.
- **Operational clarity:** systemd-managed Go services + journalctl logs.
- **Explicit boundaries:** Separate customer/admin frontends; shared API backend.
- **Contract-first API:** REST with URI versioning (`/api/v1`) and OpenAPI alignment.
- **Idempotent external workflows:** Required for payment callbacks and SMS retries.
- **Data integrity over convenience:** Strict validation, soft-delete semantics, unique constraints.
- **Solo-developer maintainability:** Conventional project layout, minimal moving parts.

---

## System Components

### Customer frontend
- Next.js static export (SSG build artifacts).
- Pages:
  - Provider list
  - Plan list by provider
  - Checkout/payment initiation
  - Payment result page
- Runtime behavior:
  - Fetches provider/plan/service request APIs.
  - Obtains reCAPTCHA v3 token before payment initiation.
  - Displays exact failure message for low reCAPTCHA score: `Please try after some time`.

### Admin frontend
- Next.js static export (SSG).
- Pages:
  - Login
  - Providers CRUD + CSV upload
  - Plans CRUD + CSV upload
  - Service requests list/detail/status update
- Stores JWT in secure browser storage strategy (HttpOnly cookie preferred via API-set cookie; fallback: in-memory/session storage with strict expiration).

### API service
- Gin routes under `/api/v1`.
- Modules:
  - `auth` (admin login/JWT)
  - `providers` CRUD + CSV upsert
  - `plans` CRUD + CSV upsert
  - `checkout/payments` (Razorpay order creation, callback verification, retries)
  - `service_requests` list/update/transition enforcement
  - `captcha` verification proxy/service
- GORM with soft-delete on providers/plans/users/service requests where applicable.

### DB
- PostgreSQL as sole system of record.
- Enforces constraints and partial unique indexes to allow reuse after soft-delete.

### Redis/Asynq workers
- Redis hosts job queue.
- Worker consumes:
  - SMS send jobs
  - SMS retry jobs
  - Payment reconciliation jobs (for delayed callback scenarios)

### Payment integration
- Razorpay selected.
- API service creates order and validates callback signatures.
- Idempotency key strategy:
  - Internal key per payment attempt.
  - Callback dedupe by `(gateway_payment_id, gateway_order_id)` unique constraint.
- Payment outcome updates service request to `payment success` or `payment failed`.

### SMS integration
- Trigger only when service request transitions to `completed`.
- Template: `Your plan <plan name> is now active. Here are the benefits: <Plan description>`.
- Asynchronous, non-blocking, retry-enabled.

### reCAPTCHA verification
- Backend verifies token with Google endpoint.
- Threshold policy: `score >= 0.5` pass.
- On fail: block payment processing and return user-facing message.

---

## Data Architecture

### Logical schema overview
Core tables:
- `users`
  - `id, name, phone_number, password_hash, role, status, created_at, updated_at, deleted_at`
- `providers`
  - `id, name, image_url, created_at, updated_at, deleted_at`
- `plans`
  - `id, provider_id, name, description, price, discount, is_active, created_at, updated_at, deleted_at`
- `service_requests`
  - `id, user_id, plan_id, status, created_at, updated_at, deleted_at`
- `payment_attempts`
  - `id, service_request_id, gateway_order_id, gateway_payment_id, amount, status, idempotency_key, callback_payload, created_at, updated_at`
- `service_request_retry_audit`
  - `id, service_request_id, retry_no, reason, created_at`
- `sms_notifications`
  - `id, service_request_id, phone_number, template, status, attempt_count, last_error, sent_at, created_at, updated_at`
- `schema_migrations`
  - migration version tracking

### Key indexes/constraints
- `providers`: unique index on `lower(name)` where `deleted_at is null`.
- `plans`: unique index on `(provider_id, lower(name))` where `deleted_at is null`.
- `plans`: check/logic to enforce max 100 plans per provider (application-level guard + DB trigger optional; app-level mandatory).
- `users`: index on `phone_number` (unique per active user policy; finalize in OpenAPI/impl).
- `service_requests`: indexes on `(status, created_at desc)` and `(user_id, created_at desc)`.
- `payment_attempts`: unique `(gateway_order_id)`, unique `(gateway_payment_id)` when not null.
- Foreign keys: plan->provider, service_request->user/plan, payment_attempt->service_request, sms->service_request.

### Migration strategy per release
- Versioned SQL (or GORM migrator wrapper) in ordered release folders.
- Deploy step:
  1. Backup DB.
  2. Run pending migrations.
  3. Verify migration checksum/version.
  4. Start/restart services.
- Migration failures are fail-fast and block release completion.

### Seed strategy (admin seed command)
- `seed-admin` command creates default admin only if absent.
- Seed is idempotent and logs whether admin was created/skipped.
- Default credentials allowed for bootstrap; manual change required immediately post-deploy.

---

## API Architecture

### Versioning strategy (`/api/v1`)
- All endpoints namespaced under `/api/v1`.
- Breaking changes require `/api/v2` rather than silent mutation.

### AuthN/AuthZ model
- Admin login issues JWT access token.
- JWT required for all admin CRUD/status/CSV endpoints.
- Customer endpoints are public for browse and checkout creation (with validation and reCAPTCHA enforcement).
- Authorization is endpoint-segment based (admin vs customer), not in-app RBAC complexity.

### Error envelope standard
All non-2xx responses return:
```json
{
  "error": {
    "code": "STRING_CODE",
    "message": "Human readable message",
    "details": [
      {"field": "field_name", "issue": "validation_issue"}
    ],
    "request_id": "uuid-or-correlation-id"
  }
}
```
Common codes:
- `AUTH_INVALID_CREDENTIALS`, `AUTH_UNAUTHORIZED`
- `VALIDATION_ERROR`, `RESOURCE_NOT_FOUND`, `CONFLICT_DUPLICATE`
- `PAYMENT_FAILED`, `PAYMENT_VERIFICATION_FAILED`
- `CAPTCHA_FAILED`
- `INTERNAL_ERROR`

---

## Async Processing Architecture

### Job types
- `sms.send.completed`
- `sms.retry.completed`
- `payment.reconcile.pending`

### Retry policy
- SMS: 3 attempts with exponential backoff (1m, 5m, 15m).
- Reconciliation: periodic retries up to configured max attempts.
- Permanent failure goes to dead-letter queue.

### Idempotency
- Job payload includes deterministic dedupe key.
- Worker checks existing `sms_notifications`/`payment_attempts` state before processing.
- Callback processing uses gateway identifiers as unique dedupe anchors.

### Dead-letter handling
- Failed jobs moved to DLQ with full error context.
- Admin/operator runbook includes manual replay procedure.

---

## Security Architecture

### JWT handling
- Signed with strong secret (`HS256`) or keypair (`RS256` preferred if manageable).
- Expiry: short-lived access token (e.g., 8-12h for admin ops).
- Token validation middleware on all admin routes.
- Token revocation list is out of scope for v1; rely on expiry + re-login.

### Password storage
- Store only `password_hash` using bcrypt/argon2id (bcrypt acceptable for v1).
- Never store plaintext passwords after seed bootstrap step.
- Seed credential change endpoint/process required operationally.

### Secret management
- Secrets via systemd `EnvironmentFile` on host (`/etc/dth-app/env`), permission `600`.
- No secrets in repo.
- Rotate manually in v1 with documented restart sequence.

### Transport security
- HTTPS termination at reverse proxy (Nginx/Caddy) with TLS certificates.
- Force HTTP->HTTPS redirect.
- Secure headers on admin app (CSP, X-Frame-Options, etc. baseline hardening).

---

## Reliability & Operations

### 90% uptime target approach
- Single VM with systemd auto-restart for API and worker.
- Health endpoint for API readiness/liveness.
- Simple external uptime check (cron/curl or lightweight monitor).

### Backup/restore manual scripts strategy
- Daily manual or cron-triggered `pg_dump` script with timestamped artifacts.
- Retain backups per simple retention policy (e.g., 7-14 days).
- Monthly restore drill to validate recoverability.

### systemd + journalctl logging model
- API and worker as separate systemd services.
- Structured logs (JSON lines preferred) with `request_id`, route, status, duration, error code.
- `journalctl -u dth-api` and `journalctl -u dth-worker` as primary diagnostics.

---

## Deployment Topology (single prod)
Recommended low-cost topology on Oracle Cloud:
- **One OCI Compute VM** (single-node) running:
  - Reverse proxy (TLS + static file serving for both frontends)
  - Go API service (systemd)
  - Go worker service (systemd)
  - PostgreSQL
  - Redis
- Network:
  - Public ingress only via HTTPS (80/443).
  - DB/Redis bound to localhost/private network only.
- Tradeoff accepted: single point of failure, consistent with cost and 90% uptime target.

---

## Testing Strategy (backend unit focus + critical API contract tests)
- **Mandatory:**
  - Unit tests for auth, service-request transitions, provider/plan validations, CSV validators, payment signature verifier.
- **Critical contract tests:**
  - Auth login + protected route behavior.
  - Payment initiation + callback verification path.
  - Service request status transition API.
  - CSV upload partial-failure response contract.
- **Operational verification tests:**
  - Migration run on clean DB.
  - Seed-admin idempotency.
  - Backup script dry run.

---

## Traceability Matrix

| Architecture Element ID | PRD ID | FS ID | Rationale | Verification Method |
|---|---|---|---|---|
| AE-001 Split frontends (customer/admin) | FR-001..FR-009 | FS module decomposition | Matches separated personas and permissions | E2E navigation + auth access tests |
| AE-002 REST API `/api/v1` | FR-013 | FS-005, FS-006 | Contract consistency and versioning | OpenAPI lint + integration tests |
| AE-003 JWT admin auth | FR-006 | FS-006 | Secure admin-only operations | Unit + integration auth tests |
| AE-004 reCAPTCHA verification service | FR-004 | FS-004 | Bot mitigation before payment | Contract test with pass/fail token |
| AE-005 Razorpay adapter + callback verifier | FR-003 | FS-005 | Payment reliability and status integrity | Signature verification tests + replay test |
| AE-006 Service request state engine | FR-009 | FS-009 + state model | Enforce valid lifecycle transitions | Transition unit tests |
| AE-007 SMS async worker | FR-010 | FS-010 | Non-blocking completion notifications | Queue integration + retry tests |
| AE-008 CSV import pipeline | FR-011, FR-012 | FS-011, FS-012 | Bulk admin operations with row-level errors | CSV fixture tests |
| AE-009 Soft-delete + partial unique indexes | FR-007, FR-008 | FS-007, FS-008 | Reusable names after soft delete | DB migration + constraint tests |
| AE-010 Retry audit table | FR-003, FR-009 | FS-003, FS-009 | Same-request retry traceability | DB assertion in retry flow tests |
| AE-011 Migration + seed framework | FR-014 | FS-013 | Release-safe schema and bootstrap | Deploy script integration test |
| AE-012 Backup/restore scripts | FR-015 | FS-014 | Basic resilience for single node | Backup success + restore drill |
| AE-013 systemd + journalctl runtime | FR-016 | FS-014 | Operational diagnostics and restartability | Service restart/log visibility checks |
| AE-014 Backend unit + contract test gate | FR-017 | FS-ALL | Minimum quality baseline | CI/local test command pass |

---

## Assumptions
1. Single-node Oracle Cloud deployment is acceptable for v1 cost constraints.
2. Next.js apps are deployed as static exports and consume backend APIs at runtime.
3. Razorpay webhook/callback endpoints can be reached from production network.
4. SMS provider supports required template and retryable API semantics.
5. Manual operational procedures (backup, credential rotation) are acceptable in v1.
6. OTP is deferred; guest checkout remains approved.

## Open Questions
| ID | Question | Owner | Suggested Due | Impact |
|---|---|---|---|---|
| OQ-TAD-001 | Final JWT expiry duration and refresh-token policy for admin app? | Engineering | 2026-03-03 | Affects security/usability tradeoff |
| OQ-TAD-002 | Should max-100-plans constraint be hard DB trigger or app-layer only? | Engineering | 2026-03-03 | Affects data consistency guarantees |
| OQ-TAD-003 | Preferred reverse proxy on VM (Nginx vs Caddy)? | Engineering | 2026-03-02 | Affects TLS automation and ops effort |
| OQ-TAD-004 | Customer phone uniqueness: strict unique vs multi-account allowed? | Product | 2026-03-04 | Affects user model and checkout behavior |

### Answers:
OQ-TAD-001: The access token can be valid for 10 days, and the refresh token can be valid for 3 months for the admin
OQ-TAD-002: max-100-plans constraint be a app-layer only
OQ-TAD-003: Preferred reverse proxy on VM is Nginx
OQ-TAD-004: Customer phone uniqueness should be strict unique

## Risks and Tradeoffs
- **Single-node architecture risk:** cheapest and simplest, but introduces SPOF.
- **Manual ops tradeoff:** low complexity now, higher human-error exposure.
- **Limited observability tradeoff:** faster setup, slower root-cause analysis.
- **Guest checkout tradeoff:** lower friction, weaker identity assurance.
- **Soft-delete reuse tradeoff:** good UX/admin flexibility, requires careful index design.

## Decision Log

| Decision ID | Date | Decision | Options Considered | Chosen Option | Rationale | Impact |
|---|---|---|---|---|---|---|
| AD-001 | 2026-02-27 | Deployment model | Multi-node managed stack vs single VM | Single OCI VM topology | Lowest cost + solo-dev operability | Accepts SPOF risk |
| AD-002 | 2026-02-27 | Frontend deployment style | SSR runtime vs static export + API | Next.js static export + API runtime fetch | Simpler hosting and lower cost | API handles dynamic logic |
| AD-003 | 2026-02-27 | Payment integration pattern | Sync-only vs callback + reconcile | Callback verification + reconcile job | Better resilience to delayed callbacks | Extra worker complexity |
| AD-004 | 2026-02-27 | SMS processing mode | Inline API call vs async queue | Async via Asynq | Non-blocking and retry-friendly | Requires Redis/worker ops |
| AD-005 | 2026-02-27 | Deletion strategy | Hard delete vs soft delete | Soft delete with reusable names | Preserves history + matches FSD decision | Requires partial unique indexes |
| AD-006 | 2026-02-27 | Retry handling for failed payments | New request per retry vs same request + audit | Same request + retry audit table | Matches approved FSD | Simpler user tracking |
| AD-007 | 2026-02-27 | CAPTCHA enforcement threshold | Variable thresholds | Fixed `>=0.5` | Approved FSD rule | Deterministic behavior |
| AD-008 | 2026-02-27 | Auth boundary model | Shared app RBAC vs split apps + JWT | Split apps + JWT-protected admin APIs | Clear persona separation | Simpler authorization model |
