## Implementation Plan - Village DTH Platform

### Metadata
- **Version:** 1.0
- **Date:** 2026-02-27
- **Owner:** Solo Developer (Staff Engineer role)
- **Input Artifacts:** PRD v1.0, FSD v1.0, TAD v1.0, OpenAPI v1.0 (`/api/v1`)
- **Execution Mode:** Single-prod, low-cost, production-ready

## Delivery assumptions and sequencing constraints
- Build order is **contract-first**: OpenAPI -> DB/migrations -> domain logic -> API handlers -> frontends -> ops hardening.
- Payment, status transitions, and async SMS are the highest-risk flows and must be completed before admin/customer UI finalization.
- CSV uploads are implemented in two phases: synchronous v1 endpoint + async-ready architecture hooks (for future larger files).
- JWT auth (10-day access, 3-month refresh) is implemented now, but refresh endpoint remains out-of-scope per decision.
- All admin endpoints require `BearerAuth`; customer endpoints are public with strict validation and reCAPTCHA checks.
- `users.phone_number` uniqueness includes soft-deleted users (global uniqueness).
- App-layer enforcement for max 100 plans per provider is mandatory; DB trigger is not included.

## Milestones (M1, M2, ...)
- **M1 Foundation & Contract Lock:** project skeleton, config, OpenAPI sync, migrations scaffold.
- **M2 Auth + Catalog Core:** admin auth, providers/plans CRUD, customer browse endpoints.
- **M3 Checkout & Payments:** service request creation, reCAPTCHA verify, Razorpay order/callback, retry audit.
- **M4 Admin Ops & Async:** service request transitions, SMS on `completed`, CSV upload + async preparation.
- **M5 Frontend Completion:** customer and admin apps wired to APIs.
- **M6 Production Readiness:** tests, backups, systemd, deployment docs, release cutover.

## Task Table

| task_id | title | status | priority | owner_role | design_pattern | api_mapping (operationId(s) + schema(s) + security requirement mapping) | files_to_create | files_to_modify | dependencies | acceptance_criteria | test_plan | risk_notes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| T001 | Backend project bootstrap and module boundaries | done | P0 | solo-dev | Layered + Dependency Injection | Schemas: `ErrorResponse`, `User`; Security: `BearerAuth` | `backend/cmd/api/main.go`, `backend/cmd/worker/main.go`, `backend/internal/{config,router,middleware,errors}/` | `backend/go.mod` | - | API and worker start with env config and health wiring stubs | unit | Poor structure early causes rework |
| T002 | OpenAPI contract import and server contract guards | done | P0 | solo-dev | Contract-First | All operationIds; Schemas: all; Security: `BearerAuth` | `backend/api/openapi/openapi.yaml`, `backend/internal/http/contract/validator.go` | `backend/internal/router/router.go` | T001 | Runtime request/response validation enabled for critical endpoints in non-prod mode | contract | Contract drift if skipped |
| T003 | Base migrations + core schema setup | done | P0 | solo-dev | Migration Runner | Schemas: `User`, `Provider`, `Plan`, `ServiceRequest`, `PaymentRequest`, `PaymentResponse`, `CsvUploadResult` | `backend/db/migrations/0001_init.sql`, `backend/db/migrations/0002_indexes.sql`, `backend/db/migrations/0003_payment_sms_retry.sql` | `backend/internal/db/db.go` | T001 | Tables/indexes exist, constraints match PRD/FSD/TAD | integration | Bad schema locks later behavior |
| T004 | Migration runner + release execution command | done | P0 | solo-dev | Command Pattern | Schemas: all persisted entities | `backend/cmd/migrate/main.go`, `backend/internal/migrate/runner.go` | `backend/Makefile` | T003 | `migrate up` runs ordered scripts with failure stop | integration | Release failures without safe migration |
| T005 | Seed admin command (idempotent) | done | P0 | solo-dev | Repository | Operation: `adminLogin`; Schemas: `AdminLoginRequest`,`AdminLoginResponse`,`User`; Security: `BearerAuth` | `backend/cmd/seed-admin/main.go`, `backend/internal/seed/admin_seed.go` | `backend/internal/users/repo.go` | T003 | Creates admin only if absent; logs create/skip | integration | Hardcoded creds operational risk |
| T006 | JWT auth and admin login endpoint | in_progress | P0 | solo-dev | Middleware + Strategy | Operation: `adminLogin`; Schemas: `AdminLoginRequest`,`AdminLoginResponse`,`ErrorResponse`; Security: `BearerAuth` | `backend/internal/auth/{handler,service,jwt.go}.go`, `backend/internal/middleware/auth.go` | `backend/internal/router/router.go` | T001,T005 | Login returns 10-day access and 3-month refresh token payload; admin routes protected | unit+integration+contract | Token misconfig blocks admin ops |
| T007 | Providers module CRUD + soft delete | todo | P0 | solo-dev | Repository | Operations: `listAdminProviders`,`createProvider`,`getProviderById`,`updateProvider`,`deleteProvider`,`listCustomerProviders`; Schemas: `Provider`,`ProviderCreateRequest`,`ProviderUpdateRequest`; Security: `BearerAuth` | `backend/internal/providers/{model,repo,service,handler}.go` | `backend/internal/router/router.go`, `backend/db/migrations/0002_indexes.sql` | T003,T006 | CRUD works, customer list excludes soft-deleted; unique name enforced | unit+integration+contract | Duplicate handling race conditions |
| T008 | Plans module CRUD + constraints (max 100/provider) | todo | P0 | solo-dev | Repository + Guard Clause | Operations: `listAdminPlans`,`createPlan`,`getPlanById`,`updatePlan`,`deletePlan`,`listCustomerPlansByProvider`; Schemas: `Plan`,`PlanCreateRequest`,`PlanUpdateRequest`; Security: `BearerAuth` | `backend/internal/plans/{model,repo,service,handler}.go` | `backend/internal/router/router.go` | T003,T006,T007 | Unique per provider and max 100/provider app-layer rule enforced | unit+integration+contract | Concurrent writes can bypass app check if not transactional |
| T009 | Customer checkout and service-request creation | todo | P0 | solo-dev | Service Layer + Idempotency | Operation: `createServiceRequestAndPayment`; Schemas: `PaymentRequest`,`PaymentResponse`,`ServiceRequest`,`User`,`ErrorResponse`; Security: none | `backend/internal/checkout/{handler,service}.go`, `backend/internal/servicerequests/{model,repo}.go` | `backend/internal/router/router.go` | T003,T007,T008 | Creates/upserts customer, creates pending service request, returns payment order payload | unit+integration+contract | Duplicate request/charge risk |
| T010 | reCAPTCHA v3 verification integration | todo | P0 | solo-dev | Adapter | Operations: `createServiceRequestAndPayment`,`retryPaymentForServiceRequest`; Schemas: `PaymentRequest`,`ErrorResponse` | `backend/internal/captcha/{client,service}.go` | `backend/internal/checkout/service.go` | T009 | Reject score <0.5 with exact error message | unit+integration | False positives may reduce conversion |
| T011 | Razorpay adapter, callback, and payment retry flow | todo | P0 | solo-dev | Adapter + Saga-lite | Operations: `handleRazorpayCallback`,`retryPaymentForServiceRequest`,`getCustomerPaymentStatus`; Schemas: `PaymentResponse`,`RazorpayCallbackRequest`,`ServiceRequest`,`ErrorResponse` | `backend/internal/payments/{razorpay_client,handler,service}.go`, `backend/internal/payments/retry_audit_repo.go` | `backend/internal/router/router.go`, `backend/internal/checkout/service.go` | T009,T010 | Callback signature verified, idempotent updates, same service request reused for retries | unit+integration+contract | Webhook replay and delayed callbacks |
| T012 | Service request state machine + admin status endpoint | todo | P0 | solo-dev | State Machine | Operations: `listServiceRequests`,`getServiceRequestById`,`updateServiceRequestStatus`,`getCustomerPaymentStatus`; Schemas: `ServiceRequest`,`ServiceRequestStatus`,`ServiceRequestStatusUpdateRequest`; Security: `BearerAuth` | `backend/internal/servicerequests/state_machine.go`, `backend/internal/servicerequests/handler.go` | `backend/internal/router/router.go` | T011,T006 | Only valid transitions accepted; invalid transitions return 422 | unit+integration+contract | Incorrect transitions corrupt lifecycle |
| T013 | SMS async jobs on completed status | todo | P0 | solo-dev | Outbox-lite + Worker | Operation: `updateServiceRequestStatus`; Schemas: `ServiceRequest`,`ErrorResponse`; Security: `BearerAuth` | `backend/internal/notifications/{sms_client,worker,jobs}.go`, `backend/internal/queue/asynq.go` | `backend/internal/servicerequests/service.go`, `backend/cmd/worker/main.go` | T012 | SMS enqueued only on transition to `completed`; retries 3x backoff | unit+integration | Notification failures without retries |
| T014 | Providers CSV upload (upsert) + async-ready abstraction | todo | P1 | solo-dev | Strategy | Operation: `uploadProvidersCsv`; Schemas: `CsvUploadResult`,`Provider`; Security: `BearerAuth` | `backend/internal/csv/providers_importer.go`, `backend/internal/csv/import_job_interface.go` | `backend/internal/providers/handler.go`, `backend/internal/router/router.go` | T007,T006 | Row-level validation + partial success response; upsert mode enabled | unit+integration+contract | CSV quality issues |
| T015 | Plans CSV upload (upsert) + async-ready abstraction | todo | P1 | solo-dev | Strategy | Operation: `uploadPlansCsv`; Schemas: `CsvUploadResult`,`Plan`; Security: `BearerAuth` | `backend/internal/csv/plans_importer.go` | `backend/internal/plans/handler.go`, `backend/internal/router/router.go` | T008,T006,T014 | Enforces provider existence, unique name/provider, max-100 rule | unit+integration+contract | Large file processing later |
| T016 | Customer frontend flows (provider->plan->checkout->status) | todo | P0 | solo-dev | Presenter + API Client | Operations: `listCustomerProviders`,`listCustomerPlansByProvider`,`createServiceRequestAndPayment`,`retryPaymentForServiceRequest`,`getCustomerPaymentStatus`; Schemas: `Provider`,`Plan`,`PaymentRequest`,`PaymentResponse`,`ServiceRequest` | `frontend/customer/pages/{index.tsx,providers/[id].tsx,checkout.tsx,payment-status/[id].tsx}`, `frontend/customer/lib/api.ts`, `frontend/customer/lib/recaptcha.ts` | `frontend/customer/package.json` | T007,T008,T009,T010,T011 | Customer can complete full purchase flow with reCAPTCHA and retry support | integration+contract (mock) | UX confusion in payment failures |
| T017 | Admin frontend flows (auth, CRUD, CSV, service requests) | todo | P0 | solo-dev | Container-Component + API Client | Operations: `adminLogin`,`listAdminProviders`,`createProvider`,`updateProvider`,`deleteProvider`,`uploadProvidersCsv`,`listAdminPlans`,`createPlan`,`updatePlan`,`deletePlan`,`uploadPlansCsv`,`listServiceRequests`,`getServiceRequestById`,`updateServiceRequestStatus`; Schemas: related admin schemas; Security: `BearerAuth` | `frontend/admin/pages/{login.tsx,providers.tsx,plans.tsx,service-requests.tsx}`, `frontend/admin/lib/{api.ts,auth.ts}` | `frontend/admin/package.json` | T006,T007,T008,T012,T014,T015 | Admin can perform all required operations end-to-end | integration+contract (mock) | Token handling bugs block admin |
| T018 | Logging, error envelope, and request correlation | todo | P0 | solo-dev | Middleware | Schemas: `ErrorResponse`,`Error`; Operations: all API operationIds | `backend/internal/middleware/{request_id,logging,error_handler}.go` | `backend/internal/router/router.go` | T001,T002 | All endpoints return standard error envelope; logs include request_id | unit+integration | Harder debugging if inconsistent |
| T019 | Backup/restore scripts + runbooks | todo | P0 | solo-dev | Scripted Procedure | Schemas: persisted core schemas (`User`,`Provider`,`Plan`,`ServiceRequest`) | `infra/scripts/{backup_db.sh,restore_db.sh}`, `infra/runbooks/backup_restore.md` | `README.md` | T003 | Backup script creates timestamped dump; restore procedure verified | integration (ops drill) | Data loss risk if untested |
| T020 | systemd + Nginx deployment assets | todo | P0 | solo-dev | Adapter (infra boundary) | Operations exposed under `/api/v1` including `adminLogin`; Security: `BearerAuth` | `infra/systemd/{dth-api.service,dth-worker.service}`, `infra/nginx/dth.conf`, `infra/deploy/deploy.sh` | `infra/README.md` | T016,T017,T018 | Services start/restart correctly; Nginx routes frontends and `/api/v1` | integration | Misrouting breaks prod |
| T021 | Backend unit tests (mandatory gate) | todo | P0 | solo-dev | Test Pyramid (Unit-heavy) | Operations/schemas covered: auth, providers, plans, payments, service requests, csv, error | `backend/internal/**/**/*_test.go` | `backend/Makefile` | T006,T007,T008,T009,T010,T011,T012,T013,T014,T015,T018 | Core logic unit tests passing | unit | Insufficient test depth |
| T022 | API contract tests against OpenAPI | todo | P0 | solo-dev | Consumer-Driven Contract | All operationIds + schemas + `BearerAuth` | `backend/tests/contract/{openapi_contract_test.go,fixtures/}.go` | `backend/api/openapi/openapi.yaml` | T021,T020 | Contract tests pass for required endpoints and status codes | contract+integration | Spec/impl drift |
| T023 | Release checklist and cutover execution | todo | P0 | solo-dev | Checklist | Operations: all production endpoints; Schemas: migration-related entities | `infra/release/release_checklist.md`, `infra/release/cutover.md` | `README.md` | T004,T005,T019,T020,T021,T022 | Release steps reproducible; rollback steps validated per milestone | integration (dry run) | Failed launch without runbook |

## Folder/File Blueprint

### Backend
- `backend/cmd/api/main.go`
- `backend/cmd/worker/main.go`
- `backend/cmd/migrate/main.go`
- `backend/cmd/seed-admin/main.go`
- `backend/internal/auth/*`
- `backend/internal/providers/*`
- `backend/internal/plans/*`
- `backend/internal/checkout/*`
- `backend/internal/payments/*`
- `backend/internal/servicerequests/*`
- `backend/internal/notifications/*`
- `backend/internal/csv/*`
- `backend/internal/middleware/*`
- `backend/internal/router/*`
- `backend/internal/errors/*`
- `backend/api/openapi/openapi.yaml`
- `backend/tests/contract/*`

### Frontend customer
- `frontend/customer/pages/index.tsx`
- `frontend/customer/pages/providers/[id].tsx`
- `frontend/customer/pages/checkout.tsx`
- `frontend/customer/pages/payment-status/[id].tsx`
- `frontend/customer/lib/api.ts`
- `frontend/customer/lib/recaptcha.ts`

### Frontend admin
- `frontend/admin/pages/login.tsx`
- `frontend/admin/pages/providers.tsx`
- `frontend/admin/pages/plans.tsx`
- `frontend/admin/pages/service-requests.tsx`
- `frontend/admin/lib/api.ts`
- `frontend/admin/lib/auth.ts`

### Infra/scripts/migrations
- `backend/db/migrations/*.sql`
- `infra/scripts/backup_db.sh`
- `infra/scripts/restore_db.sh`
- `infra/systemd/dth-api.service`
- `infra/systemd/dth-worker.service`
- `infra/nginx/dth.conf`
- `infra/deploy/deploy.sh`
- `infra/release/release_checklist.md`
- `infra/release/cutover.md`

## Migration & Seed Execution Plan
1. Put app in maintenance mode (if needed for schema-breaking releases).
2. Run DB backup: `infra/scripts/backup_db.sh`.
3. Execute migrations: `backend/cmd/migrate up`.
4. Verify migration version table and key constraints.
5. Run admin seed once: `backend/cmd/seed-admin`.
6. Restart services: `systemctl restart dth-api dth-worker`.
7. Run smoke tests for `adminLogin`, provider list, checkout create, callback endpoint.
8. Record release artifact version and migration versions in release log.

## Rollback Plan per milestone
- **M1 rollback:** revert code, drop newly created non-critical tables only in non-prod; in prod restore DB backup if migration applied.
- **M2 rollback:** disable new routes via router flags; restore previous binary; DB remains backward-compatible.
- **M3 rollback:** disable payment create/retry routes; keep service requests readable; restore pre-M3 DB backup if schema incompatible.
- **M4 rollback:** disable worker services and CSV endpoints; manual admin ops continue via existing CRUD.
- **M5 rollback:** redeploy previous frontend static bundles while keeping backend stable.
- **M6 rollback:** full release rollback by restoring previous binaries/config and last verified DB backup.

## Traceability Matrix
| task_id | PRD ID | FS ID | Arch ID | OpenAPI operationId/schema/security |
|---|---|---|---|---|
| T001 | FR-013, FR-016 | FS-014 | AE-002, AE-013 | `ErrorResponse`, `BearerAuth` |
| T002 | FR-013 | FS-005, FS-006 | AE-002 | all operationIds + all schemas + `BearerAuth` |
| T003 | FR-005, FR-007, FR-008, FR-009 | FS-003, FS-007, FS-008, FS-009 | AE-009, AE-010 | `User`,`Provider`,`Plan`,`ServiceRequest`,`PaymentResponse` |
| T004 | FR-014 | FS-013 | AE-011 | core persistence schemas |
| T005 | FR-014, FR-006 | FS-013, FS-006 | AE-011, AE-003 | `adminLogin`,`AdminLoginResponse`,`User`,`BearerAuth` |
| T006 | FR-006, FR-013 | FS-006 | AE-003, AE-002 | `adminLogin`,`AdminLoginRequest`,`AdminLoginResponse`,`BearerAuth` |
| T007 | FR-001, FR-007 | FS-001, FS-007 | AE-008, AE-009 | `listCustomerProviders`,`listAdminProviders`,`createProvider`,`getProviderById`,`updateProvider`,`deleteProvider`,`Provider*`,`BearerAuth` |
| T008 | FR-002, FR-008 | FS-002, FS-008 | AE-008, AE-009 | `listCustomerPlansByProvider`,`listAdminPlans`,`createPlan`,`getPlanById`,`updatePlan`,`deletePlan`,`Plan*`,`BearerAuth` |
| T009 | FR-003, FR-005 | FS-003 | AE-010 | `createServiceRequestAndPayment`,`PaymentRequest`,`PaymentResponse`,`ServiceRequest` |
| T010 | FR-004 | FS-004 | AE-004 | `createServiceRequestAndPayment`,`retryPaymentForServiceRequest`,`PaymentRequest`,`ErrorResponse` |
| T011 | FR-003, FR-013 | FS-005 | AE-005, AE-010 | `handleRazorpayCallback`,`retryPaymentForServiceRequest`,`getCustomerPaymentStatus`,`RazorpayCallbackRequest`,`PaymentResponse`,`ServiceRequest` |
| T012 | FR-009 | FS-009 | AE-006 | `listServiceRequests`,`getServiceRequestById`,`updateServiceRequestStatus`,`ServiceRequestStatus*`,`BearerAuth` |
| T013 | FR-010 | FS-010 | AE-007 | `updateServiceRequestStatus`,`ServiceRequest`,`BearerAuth` |
| T014 | FR-011 | FS-011 | AE-008 | `uploadProvidersCsv`,`CsvUploadResult`,`Provider`,`BearerAuth` |
| T015 | FR-012 | FS-012 | AE-008 | `uploadPlansCsv`,`CsvUploadResult`,`Plan`,`BearerAuth` |
| T016 | FR-001..FR-004 | FS-001..FS-005 | AE-001, AE-004, AE-005 | customer operationIds + `Provider`,`Plan`,`PaymentRequest`,`PaymentResponse`,`ServiceRequest` |
| T017 | FR-006..FR-012 | FS-006..FS-012 | AE-001, AE-003, AE-006, AE-008 | admin operationIds + `BearerAuth` + admin schemas |
| T018 | FR-013, FR-016 | FS-005, FS-006, FS-014 | AE-002, AE-013 | all operationIds + `ErrorResponse` |
| T019 | FR-015 | FS-014 | AE-012 | persistence schemas (`User`,`Provider`,`Plan`,`ServiceRequest`) |
| T020 | FR-016, FR-013 | FS-014 | AE-013, AE-002 | `/api/v1` operations + `BearerAuth` enforcement path |
| T021 | FR-017 | FS-ALL | AE-014 | all critical operationIds/schemas |
| T022 | FR-013, FR-017 | FS-005, FS-006 | AE-002, AE-014 | all operationIds + all schemas + `BearerAuth` |
| T023 | FR-014, FR-015, FR-016 | FS-013, FS-014 | AE-011, AE-012, AE-013 | release-critical operations and persistence schemas |

## Assumptions
- Solo developer executes tasks sequentially with minimal parallel streams.
- Existing repository can host `backend/`, `frontend/customer`, and `frontend/admin` structure.
- Razorpay and SMS credentials are available before M3.
- OTP remains deferred; guest checkout is acceptable for v1.
- Refresh token is returned by login response but no dedicated refresh endpoint in v1.

## Open Questions
- For async CSV future mode, should same endpoints return `202 + job_id` when file exceeds threshold?
    Answer: No, it is not required
- Should admin login support username only or username+phone equally from day one?
    Answer: username
- Do we need customer-facing service-request history page in v1 or only payment-status polling?
    Answer: No, we don't need customer facing service request history page

## Risks
- Payment callback timing and retries can create status drift if idempotency is incomplete.
- Long admin token validity increases exposure if token storage is weak.
- App-layer max-plan constraint can fail under race unless DB transaction locks are used.
- Single-node deployment increases outage impact.
- Manual credential and backup procedures can fail without routine drills.

## Decision Log
| Decision ID | Date | Decision | Options Considered | Chosen Option | Rationale | Impact |
|---|---|---|---|---|---|---|
| IMP-001 | 2026-02-27 | Build sequence | UI-first vs contract-first | Contract-first | Reduces rework and integration mismatches | Faster stable delivery |
| IMP-002 | 2026-02-27 | CSV architecture | Sync-only vs async-only vs sync+async-ready | Sync now + async-ready abstraction | Matches low scale now and future requirement | Minimal v1 overhead |
| IMP-003 | 2026-02-27 | Payment retry model | New request vs same request | Same service request + retry audit | Approved FSD decision | Clean traceability |
| IMP-004 | 2026-02-27 | Constraint enforcement | DB trigger vs app-layer for max plans | App-layer only | Approved TAD decision | Simpler migrations, higher concurrency care |
| IMP-005 | 2026-02-27 | Deployment topology | Managed multi-node vs single OCI VM | Single VM | Cost and solo-dev constraints | Accept SPOF risk |
| IMP-006 | 2026-02-27 | Auth refresh endpoint | Include now vs defer | Defer | Explicit OpenAPI decision | Simpler v1 auth surface |
