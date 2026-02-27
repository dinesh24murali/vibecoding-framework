# Functional Specification (FSD) - Village DTH Plan Purchase Platform

## Metadata
- **Version:** 1.0
- **Author:** Senior Systems Analyst / Functional Architect
- **Date:** 2026-02-27
- **Linked PRD Version:** 1.0 (`Docs/03_Production_PRD.md`)
- **Status:** Approved-for-Build

## Functional Scope Summary
This FSD defines build-ready functional behavior for:
1. Customer web app (provider selection -> plan selection -> payment with reCAPTCHA v3).
2. Admin web app (JWT login, provider CRUD, plan CRUD, service request status management, CSV uploads).
3. Shared backend services (REST `/api/v1`, Razorpay integration, SMS on completion, soft-delete, migrations/seed support).

In scope is v1 delivery for a low-scale, single-operator deployment with a solo developer and one production environment.  
Out of scope includes advanced performance optimization, multi-tenant behavior, and enterprise compliance programs.

## Actors and Permissions Model

### Actor A1 - Customer (Guest Checkout)
- Accesses customer app without mandatory account login.
- Can browse providers and plans.
- Can initiate payment for selected plan.
- Cannot access admin functions.

### Actor A2 - Admin
- Authenticates through admin login endpoint.
- Receives JWT token and accesses protected admin APIs.
- Can perform CRUD on providers and plans.
- Can upload CSV for providers and plans.
- Can view service requests and update status.

### Permission Model
- Customer app endpoints: public read + controlled write for checkout/payment.
- Admin endpoints: JWT-required.
- No in-app multi-role RBAC matrix; role separation is by separate apps and protected admin API surface.

## Functional Decomposition by Module

### Customer App
- Provider list display.
- Provider-specific plan list display.
- Plan selection and checkout entry.
- Payment submit with reCAPTCHA v3 token.
- Payment result display (success/failure).
- Service request tracking acknowledgement.

### Admin App
- Admin login/session handling with JWT.
- Provider management UI and actions.
- Plan management UI and actions.
- Service request list/detail/status update.
- CSV upload screens for provider and plan bulk operations.

### Shared Backend Services
- REST API (`/api/v1`) and uniform error envelopes.
- Data persistence (PostgreSQL via GORM, soft-delete).
- Razorpay payment order/callback/status handling.
- reCAPTCHA v3 verification endpoint/service.
- SMS dispatch trigger when service request marked `completed`.
- Migration runner and admin seed command.
- Logging to journalctl-compatible output.

---

## Detailed Functional Requirements (FS-001...)

### FS-001 Customer Provider Discovery
- **Linked PRD:** FR-001
- **Trigger:** Customer lands on customer home page.
- **Preconditions:** Backend is reachable; provider data exists or empty state allowed.
- **Main Flow:**
  1. Customer app requests provider list.
  2. Backend returns active/non-deleted providers.
  3. UI renders provider cards (name, image).
- **Alternate Flows:**
  - A1: No providers available -> show empty-state message.
  - A2: Provider image missing -> show fallback placeholder image.
- **Postconditions:** Customer can select a provider or exit.
- **Error handling:** On API failure, show retryable error and “Try again”.
- **Business validations:** Provider names MUST be unique among non-deleted records.

### FS-002 Customer Plan Discovery by Provider
- **Linked PRD:** FR-002
- **Trigger:** Customer selects a provider.
- **Preconditions:** Selected provider exists and is non-deleted.
- **Main Flow:**
  1. Customer app requests plans filtered by provider.
  2. Backend returns active, non-deleted plans for provider.
  3. UI renders plan cards (name, description, price, discount).
- **Alternate Flows:**
  - A1: Provider has no active plans -> show provider-specific empty state.
- **Postconditions:** Customer can choose one plan to proceed.
- **Error handling:** Invalid provider ID -> `404`; API failure -> retry prompt.
- **Business validations:**
  - Plan `provider_id` MUST reference existing provider.
  - Max plans per provider MUST NOT exceed 100 (across active + inactive, excluding soft-deleted if business decides otherwise; see OQ).

### FS-003 Checkout and Service Request Creation
- **Linked PRD:** FR-003, FR-005, FR-009
- **Trigger:** Customer clicks “Proceed to Payment” from selected plan.
- **Preconditions:** Provider and plan are valid and active.
- **Main Flow:**
  1. Capture customer details (name, phone).
  2. Create or upsert customer record (role=`customer`).
  3. Create service request in `Pending` status with user_id and plan_id.
  4. Create payment order with gateway metadata.
- **Alternate Flows:**
  - A1: Existing customer by phone -> reuse existing user record.
  - A2: Duplicate in-flight request for same user+plan -> return existing pending request reference (idempotent behavior).
- **Postconditions:** Service request exists and is linked to payment context.
- **Error handling:** Validation failures return field-specific errors; payment order creation failure sets `payment failed` if order cannot be established.
- **Business validations:** Phone number format required; plan must remain active at order creation time.

### FS-004 reCAPTCHA v3 Payment Gate
- **Linked PRD:** FR-004
- **Trigger:** Customer submits payment action.
- **Preconditions:** Service request exists in `Pending`.
- **Main Flow:**
  1. Client sends reCAPTCHA token.
  2. Backend verifies token with Google.
  3. If score >= 0.5, continue payment processing.
- **Alternate Flows:**
  - A1: Score < 0.5 or verification failure -> block action.
- **Postconditions:** Payment proceeds only if verification passes.
- **Error handling:** Show exact user message: `Please try after some time`.
- **Business validations:** reCAPTCHA token MUST be present and non-expired.

### FS-005 Payment Processing and Callback Handling (Razorpay)
- **Linked PRD:** FR-003, FR-013
- **Trigger:** Payment initiation and gateway callback.
- **Preconditions:** Valid service request and payment order.
- **Main Flow:**
  1. Backend creates Razorpay order.
  2. Customer completes payment on gateway.
  3. Backend receives callback/verification payload.
  4. Verify signature and payment status.
  5. Update service request status to `payment success` or `payment failed`.
- **Alternate Flows:**
  - A1: Callback delayed -> polling/status API reconciles.
  - A2: Duplicate callback -> ignore duplicate update via idempotency key/signature+order tracking.
  - A3: Callback signature invalid -> reject and log security event.
- **Postconditions:** Service request reflects final payment outcome until admin fulfillment.
- **Error handling:** Payment uncertain states are logged and flagged for admin review.
- **Business validations:** Only authorized status transitions from payment events are allowed.

### FS-006 Admin Authentication (JWT)
- **Linked PRD:** FR-006, FR-013
- **Trigger:** Admin submits credentials on login page.
- **Preconditions:** Admin user exists and is active.
- **Main Flow:**
  1. Validate credentials.
  2. Issue JWT with configured expiry and claims.
  3. Admin app stores token securely for session.
- **Alternate Flows:**
  - A1: Invalid credentials -> reject login.
  - A2: Expired token on API call -> return unauthorized and force re-login.
- **Postconditions:** Authenticated admin can access protected modules.
- **Error handling:** Standardized `401` and `403` responses.
- **Business validations:** Only users with role `admin` can authenticate to admin endpoints.

### FS-007 Provider CRUD
- **Linked PRD:** FR-007
- **Trigger:** Admin performs create/read/update/delete action.
- **Preconditions:** Admin authenticated.
- **Main Flow:**
  1. Admin submits provider form data.
  2. Backend validates uniqueness and required fields.
  3. Persist changes.
- **Alternate Flows:**
  - A1: Duplicate provider name -> reject with conflict.
  - A2: Soft-delete requested with linked plans -> allow soft-delete; linked plans become inaccessible by business policy.
- **Postconditions:** Provider catalog reflects latest non-deleted entries.
- **Error handling:** Validation conflict and server errors are surfaced with actionable messages.
- **Business validations:** Provider `name` MUST be unique (case-insensitive recommended).

### FS-008 Plan CRUD
- **Linked PRD:** FR-008
- **Trigger:** Admin performs plan create/read/update/delete action.
- **Preconditions:** Admin authenticated; referenced provider exists.
- **Main Flow:**
  1. Admin submits plan details.
  2. Backend validates provider, uniqueness, limit constraints.
  3. Persist plan and publish availability state.
- **Alternate Flows:**
  - A1: Duplicate plan name under same provider -> reject.
  - A2: Attempt to exceed 100 plans/provider -> reject.
- **Postconditions:** Plan catalog updated and reflected in customer app if `is_active=true`.
- **Error handling:** `409` for duplicates/limits; `404` for missing provider.
- **Business validations:**
  - Unique constraint: (`provider_id`, `plan_name`) for non-deleted records.
  - Max 100 plans per provider.

### FS-009 Service Request Operations
- **Linked PRD:** FR-009
- **Trigger:** Admin views or updates service request.
- **Preconditions:** Admin authenticated; request exists.
- **Main Flow:**
  1. Admin lists requests with filters.
  2. Admin opens request detail.
  3. Admin updates status as fulfillment progresses.
- **Alternate Flows:**
  - A1: Invalid transition attempt -> reject with transition rule error.
- **Postconditions:** Status and updated timestamp are persisted.
- **Error handling:** Missing request -> `404`; invalid transition -> `422`.
- **Business validations:** Transition rules enforced (see State Model section).

### FS-010 SMS Notification on Completion
- **Linked PRD:** FR-010
- **Trigger:** Admin sets service request status to `completed`.
- **Preconditions:** Service request has valid user phone and linked plan data.
- **Main Flow:**
  1. Status update to `completed` committed.
  2. Async SMS job enqueued.
  3. SMS sent using template:  
     `Your plan <plan name> is now active. Here are the benefits: <Plan description>`
  4. Delivery outcome logged.
- **Alternate Flows:**
  - A1: SMS provider timeout/failure -> retry based on policy; do not rollback completed status.
- **Postconditions:** Customer notified when possible; audit trail retained.
- **Error handling:** Failed sends recorded with reason and retry state.
- **Business validations:** SMS is sent only on transition to `completed`.

### FS-011 CSV Upload - Providers
- **Linked PRD:** FR-011
- **Trigger:** Admin uploads provider CSV.
- **Preconditions:** Admin authenticated; file format valid CSV.
- **Main Flow:**
  1. Parse CSV header and rows.
  2. Validate required fields and uniqueness constraints.
  3. Insert/update valid rows.
  4. Return import summary.
- **Alternate Flows:**
  - A1: Row-level validation errors -> row rejected, import continues.
  - A2: Entirely invalid header -> reject file.
- **Postconditions:** Valid provider rows are persisted.
- **Error handling:** Response includes processed/success/failed counts and row errors.
- **Business validations:** Duplicate provider names rejected.

### FS-012 CSV Upload - Plans
- **Linked PRD:** FR-012
- **Trigger:** Admin uploads plan CSV.
- **Preconditions:** Admin authenticated; referenced provider resolvable.
- **Main Flow:**
  1. Parse and validate rows.
  2. Validate provider existence, duplicate names per provider, and max plans/provider.
  3. Persist valid rows.
  4. Return row-level report.
- **Alternate Flows:**
  - A1: Unknown provider -> row rejected.
  - A2: Duplicate plan name under provider -> row rejected.
  - A3: Provider over 100 plans -> row rejected.
- **Postconditions:** Partial success allowed with explicit report.
- **Error handling:** Structured row error payload with line number and reason.
- **Business validations:** Uniqueness and limit rules strictly enforced.

### FS-013 Migration and Seed Operations
- **Linked PRD:** FR-014
- **Trigger:** Deployment/release operation.
- **Preconditions:** DB connectivity and migration files available.
- **Main Flow:**
  1. Execute pending migrations in ordered sequence.
  2. Record migration version in migrations table.
  3. Run seed command for initial admin if absent.
- **Alternate Flows:**
  - A1: Migration fails mid-run -> stop, mark failed migration, require manual intervention.
- **Postconditions:** Schema aligned to release; admin bootstrap available.
- **Error handling:** Fail-fast with detailed logs.
- **Business validations:** Seed should be idempotent for admin creation.

### FS-014 Logging, Backup, and Availability Support
- **Linked PRD:** FR-015, FR-016
- **Trigger:** Runtime operations and scheduled/manual maintenance.
- **Preconditions:** Service installed under systemd; DB credentials configured.
- **Main Flow:**
  1. API logs structured output to stdout/stderr for journalctl capture.
  2. Manual backup script runs PostgreSQL dump.
  3. Restore script verifies successful recoverability.
- **Alternate Flows:**
  - A1: Backup failure -> non-zero exit + log critical event.
- **Postconditions:** Basic operations readiness for 90% uptime target.
- **Error handling:** Operational scripts emit actionable failure messages.
- **Business validations:** Backup and restore must be executable by operator-admin process.

---

## State Models

### Service Request Lifecycle
**States:**
- `Pending`
- `payment failed`
- `payment success`
- `blocked`
- `completed`

**Allowed transitions:**
1. `Pending` -> `payment success` (gateway verified success)
2. `Pending` -> `payment failed` (gateway failure/verification failure after attempt)
3. `payment success` -> `blocked` (admin operational hold)
4. `payment success` -> `completed` (manual activation done)
5. `blocked` -> `payment success` (unblock to ready state)
6. `blocked` -> `completed` (activation completed after review)

**Terminal guidance:**
- `completed` is terminal for v1.
- `payment failed` is terminal for a given attempt; new payment creates new attempt context (same or new request per implementation choice).

**Invalid transitions (examples):**
- `Pending` -> `completed` (not allowed)
- `payment failed` -> `completed` (not allowed without new successful payment)
- `completed` -> any other state (not allowed)

---

## UI/UX Functional Behavior (not visual design)

### Field-level Validation Rules
- **Provider:** `name` required, unique; `image_url` optional but must be valid URL if present.
- **Plan:** `provider_id` required; `name` required and unique within provider; `price` required and > 0; `discount` >= 0; `is_active` boolean.
- **Customer checkout:** `name` required; `phone_number` required and pattern-validated.
- **Admin login:** username/phone + password required.
- **reCAPTCHA token:** required on payment submit.

### Form Behaviors
- Inline validation on blur and on submit.
- Submit buttons disabled during in-flight requests.
- Server validation errors mapped to specific fields where possible.
- Destructive admin actions require confirmation modal.
- Status updates show optimistic lock prevention (latest state check before save).

### Pagination/Sorting/Filtering Rules
- **Providers admin list:** pagination required; sort by updated date desc default; search by name.
- **Plans admin list:** pagination required; filters by provider and active status; sort by updated date desc.
- **Service requests admin list:** pagination required; filters by status, provider, date range; sort by created date desc.
- **Customer lists:** no pagination required at current expected scale; MAY add if dataset grows.

---

## Notifications Behavior (SMS events, retries, failure handling)
- SMS event is triggered **only when service request transitions to `completed`**.
- Template: `Your plan <plan name> is now active. Here are the benefits: <Plan description>`.
- SMS send is asynchronous and non-blocking to core status update.
- Retry policy: 3 attempts with exponential backoff (e.g., 1m, 5m, 15m).
- After final failure, mark notification status as failed; keep request `completed`.
- All send attempts and outcomes are logged with request ID correlation.

---

## CSV Upload Functional Rules

### Providers Upload
- Required columns: `name`, optional `image_url`.
- Header validation mandatory.
- Duplicate provider names in file or DB rejected row-by-row.
- Partial success enabled; import response contains:
  - total rows
  - successful rows
  - failed rows
  - row-level error list

### Plans Upload
- Required columns: `provider_name` or `provider_id` (implementation decision), `name`, `description`, `price`, `discount`, `is_active`.
- Each row validates:
  - provider exists
  - plan name unique within provider
  - provider plan count does not exceed 100
- Partial success enabled with row-level error reporting.

### Validation and Partial Failure Handling
- File-level hard fail: invalid CSV syntax or invalid headers.
- Row-level fail: field format, duplicates, missing provider, business limit violation.
- Persist valid rows in transactional batches (batch-level rollback on DB failure).

---

## Traceability Matrix

| FS ID | PRD ID | Context Source | Test Case Hint |
|---|---|---|---|
| FS-001 | FR-001 | Customer selects provider first | Verify provider list/empty state |
| FS-002 | FR-002 | Plan chosen after provider | Verify provider-filtered active plans |
| FS-003 | FR-003, FR-005, FR-009 | Checkout + service request entity model | Verify request created and linked to user/plan |
| FS-004 | FR-004 | reCAPTCHA v3 required with threshold >=0.5 | Verify fail message on low score |
| FS-005 | FR-003, FR-013 | Razorpay integration + REST contract | Verify callback signature/idempotency |
| FS-006 | FR-006, FR-013 | JWT auth model | Verify token issuance and protected route access |
| FS-007 | FR-007 | Admin provider CRUD requirement | Verify create/update/soft-delete/duplicate reject |
| FS-008 | FR-008 | Admin plan CRUD + constraints | Verify unique per provider + max 100 plans |
| FS-009 | FR-009 | Admin service request management | Verify valid/invalid status transitions |
| FS-010 | FR-010 | SMS required + resolved OQ behavior | Verify SMS only on `completed` transition |
| FS-011 | FR-011 | Provider CSV upload requirement | Verify partial success and row errors |
| FS-012 | FR-012 | Plan CSV upload requirement | Verify validation and constraint enforcement |
| FS-013 | FR-014 | Migration + seed required | Verify ordered migrations + idempotent seed |
| FS-014 | FR-015, FR-016 | 90% uptime, backups, journalctl logs | Verify backup/restore scripts and service logs |
| FS-ALL | FR-017 | Backend unit test gate | Verify unit tests for core logic |

---

## Assumptions
1. Customer checkout remains guest-first; customer records are created/upserted during checkout.
2. Razorpay provides required APIs for order creation, verification, and callback handling.
3. Soft-delete semantics via GORM are acceptable for provider/plan entities and history retention.
4. Admin credential rotation is manual and operationally acceptable for v1.
5. Single production deployment can meet expected user volume (30-50 active users).
6. SMS provider supports required template and regional delivery routes.
7. No additional compliance controls are required beyond baseline secure implementation.

## Open Questions
| OQ ID | Question | Owner | Due Date (Suggested) | Impact |
|---|---|---|---|---|
| OQ-FSD-001 | For plan/provider uniqueness, should soft-deleted names be reusable immediately? | Engineering Owner | 2026-03-03 | Affects unique index strategy |
| OQ-FSD-002 | For `payment failed`, should retry reuse same service request or create a new request per attempt? | Product Owner | 2026-03-03 | Affects reporting and reconciliation |
| OQ-FSD-003 | For guest checkout, is OTP verification of phone required now or deferred? | Product Owner | 2026-03-04 | Affects fraud and identity confidence |
| OQ-FSD-004 | CSV import mode: upsert existing rows or strict create-only for v1? | Product Owner | 2026-03-03 | Affects admin operations and data safety |

### Answers:

OQ-FSD-001: yes, soft-deleted names should be reusable immediately
OQ-FSD-002: Use the same service request for retry. Maintain a retry audit table
OQ-FSD-003: OTP verification of phone is deferred
OQ-FSD-004: upsert existing rows

## Risks
- Payment callback ambiguity may create inconsistent status without robust idempotency.
- Manual credential handling increases admin account compromise risk.
- Single-env deployment reduces rollback safety.
- CSV bulk imports can introduce catalog inconsistency if validation is weak.
- Minimal observability may slow incident diagnosis.
- Guest checkout may allow low-quality user records unless phone validation is strict.

## Decision Log

| Decision ID | Date | Decision | Options Considered | Chosen Option | Rationale | Impact |
|---|---|---|---|---|---|---|
| FD-001 | 2026-02-27 | Payment gateway selection | Razorpay, PayU, Cashfree | Razorpay | Resolved from PRD OQ-001 | Enables concrete payment integration scope |
| FD-002 | 2026-02-27 | SMS trigger event | On payment success, on completed, both | On `completed` only | Resolved from PRD OQ-002 | Aligns notification with actual activation |
| FD-003 | 2026-02-27 | Customer auth model | Login required, guest checkout | Guest checkout + reCAPTCHA | Resolved from PRD OQ-003 | Faster funnel, simpler UX |
| FD-004 | 2026-02-27 | reCAPTCHA threshold and fail UX | 0.3/0.5/0.7 and custom errors | >= 0.5 and fixed error message | Resolved from PRD OQ-004 | Clear anti-bot policy |
| FD-005 | 2026-02-27 | Delete semantics | Hard delete, soft delete | Soft delete (GORM) | Resolved from PRD OQ-005 | Preserves historical integrity |
| FD-006 | 2026-02-27 | Catalog constraints | No limits, bounded limits | Unique names + max 100 plans/provider | Resolved from PRD OQ-006 | Prevents catalog sprawl/inconsistency |
| FD-007 | 2026-02-27 | Active customer KPI definition | One-time, recurring, login-based | Recurring payer | Resolved from PRD OQ-007 | Clarifies retention metric logic |
| FD-008 | 2026-02-27 | Credential rotation policy | Scheduled automated, manual ad hoc | Manual handling in v1 | Resolved from PRD OQ-008 | Lower complexity, higher security risk |
