You are a Senior API Architect.
Generate a PRODUCTION-READY OpenAPI 3.1 specification in YAML.
Do not output pseudo-specs. Output valid, complete OpenAPI content.

Inputs:
1) Project Context Pack
2) Approved PRD
3) Approved Functional Spec
4) Approved Technical Architecture
5) Constraints: REST + URI versioning + JWT

Project Context Pack:
```
Please answer the questions below (short bullets are fine).
1) Vision, goals, success metrics
- What is the product in one sentence?
The product is used by customers to purchase DTH plans. The users will first choose the service provider, then the plan and will finally make the payment. The payment page will have Google reCaptcha V3. There will be a customer facing web application where they will choose the plan and make payment, and a web application for the admin who will do CRUD for plans and service providers. The admin will got through the service requests, make the necessary changes and change the status of the service request.
- What business problem are we solving, and why now?
This is for a very small niche audience. This product is for a small village where a single DTH operator will accepts orders from the people in the village. Once he gets the orders he will manually activate the plans for the relevant customer. They have been needing an application like this for quite some time.
- Top 3 business goals for the first 6–12 months?
Having 30 to 50 active customers using the webapp. The web application should not have any glitches. Keeping the cost as low as possible
- What metrics define success (e.g., activation %, conversion %, revenue, retention, SLA)?
If people continue to use the website. activation 20%, retention 25%
2) Users and journeys
- Primary persona(s): who uses it most, and what do they need done?
The customer will use the website to purchase DTH plans.
- Secondary persona(s): who else matters?
The admin will manage the application. he will handle CRUD for providers, and plans. He will manage the service requests.
- What are the top 3 critical user journeys (end-to-end steps + desired outcome)?
Customer Plan purchase flow: When the Customers reach the landing page, they will see a screen to select the service provider. Then the customer will choose the plan. Then the customer will make the payment
Admin CRUD DTH providers: Admin will login to the application. He will choose the option on te top NavBar to switch to the DTH providers list page. He will do CRUD operations for the providers
Admin CRUD DTH Plans: Admin will login to the application. He will choose the option on te top NavBar to switch to the DTH Plans list page. He will do CRUD operations for the Plans
- Any high-risk or high-friction journey we should optimize first?
No
3) Scope boundaries
- What is explicitly in-scope for v1?
All of the above mentioned requirements are part of V1
- What is explicitly out-of-scope for v1?
No strong focus on performance
- What is “nice-to-have” if time permits?
Everything is required
4) Functional requirements and business rules
- Core features/modules required for launch?
All the above mentioned are required for launch
- Key business rules or validations (pricing, approvals, eligibility, limits, workflows)?
No
- Any role-based permissions needed (admin, manager, end-user, support)?
No RBAC in the same application. There is one admin web app and another end-user web app.
- Are there required notifications (email/SMS/push/webhooks)?
We need SMS.
5) Non-functional requirements
- Performance targets (p95 response time, page load, throughput)?
Not a high priority as of now
- Reliability targets (uptime %, RTO/RPO, backup expectations)?
The uptime should be 90%, need scripts to take manual DB backups
- Security requirements (SSO/MFA, encryption, secrets handling, audit trails)?
Not required
- Compliance/privacy requirements (GDPR, HIPAA, SOC2, data residency, retention)?
Not required
- Accessibility target (WCAG level, supported assistive tech)?
Not required
6) Data model and integrations
- What are the core business entities (e.g., User, Account, Order, Subscription)?
User - id, name, phone number, created date, updated date, status, role (customer, admin), password
Providers - id, name, created date, updated date, image_url
Plan - id, provider_id, name, description, price, discount, is_active
Service Request - id, user_id, status (Pending, payment failed, payment success, blocked, completed), plan_id
- Any expected data volumes (records/day, file sizes, growth)?
Not required
- System-of-record(s): where does truth live for key entities?
- Required integrations (payments, CRM, ERP, identity, analytics, messaging)?
payment integration is required
- Import/export needs (CSV, APIs, scheduled syncs, webhooks)?
CSV upload for providers, and plans
7) API and platform interface
- Preferred API style: REST, GraphQL, or mixed?
REST is preferred
- Auth model: session, JWT, OAuth2/OIDC, API keys, service-to-service auth?
JWT auth
- Versioning strategy (URI versioning, header-based, schema evolution)?
URI versioning
- Any public API or partner-facing developer platform requirements?
Not required
8) Technical constraints and delivery
- Preferred cloud/provider and managed services constraints?
Oracle cloud is preferred
- Required tech stack (frontend, backend, DB, queue, cache)?
Frontend: NextJS with Static site generation
backend: Go 1.26 + Gin 1.10 API server + Gorm
DB: PostgreSQL
Message broker & cache: redis
Background tasks: asynq
- Team composition/size (engineering, product, design, QA, DevOps)?
Solo developer
- Release cadence target (weekly, biweekly, continuous)?
Need not be tracked
- Budget/hosting constraints and expected runway?
Need not be tracked
9) Environments and release strategy
- Required environments (local/dev/stage/prod + preview envs)?
Single prod environment for now
- Deployment strategy (blue/green, canary, rolling)?
Need not be tracked
- Test strategy gates (unit/integration/e2e, performance, security scans)?
unit testing for the backend code alone
- Data strategy per environment (seed data, masked prod snapshots)?
Create seed commands to create an admin user with a string hardcoded credentials
Need a system to run migration scripts based on releases.
10) Observability and operations
- Logging/metrics/tracing expectations (tooling preferences)?
I am going to run the Gin webserver as a systemd service with journalctl for logs. No other tracking is required
- Alerting model (on-call hours, PagerDuty/Slack escalation, severity levels)?
Need not be tracked
- Operational dashboards needed for product + engineering?
Need not be tracked
- Incident response expectations (runbooks, postmortems, SLAs/SLOs)?
Need not be tracked
```

Approved PRD:
```
# Product Requirements Document (PRD) - Village DTH Plan Purchase Platform

## Metadata
- **Version:** 1.0
- **Date:** 2026-02-27
- **Owner:** Product Owner / Solo Developer
- **Status:** Approved-for-Build

## Executive Summary
This product digitizes DTH plan ordering for a single-operator village market through two web applications: a customer-facing purchase app and an admin operations app. Customers select a provider, choose a plan, and complete payment with Google reCAPTCHA v3 protection. Admin manages providers and plans, reviews service requests, and updates request status based on manual activation workflow.  
The v1 objective is a stable, low-cost system that supports 30-50 active users, achieves at least 20% activation and 25% retention, and runs reliably in a single production environment.

## Problem Statement and Opportunity
The village DTH operator currently manages orders manually without a dedicated digital ordering flow. This creates friction for customers, inconsistent request tracking, and operational overhead for the admin.  
A lightweight web platform creates a single source of truth for plans, providers, and service requests while reducing manual coordination and enabling repeat usage. The opportunity is to provide a fit-for-purpose, low-cost product with sufficient reliability (90% uptime target) for a niche but recurring local demand.

## Product Vision and Goals

### Vision
Enable village customers to easily purchase DTH plans online while giving the local operator a simple operational console to manage catalog and activation requests.

### 6-12 Month Goals
1. Reach **30-50 active customers**.
2. Deliver a stable product with minimal user-facing glitches.
3. Keep operating and development costs low for a solo-developer model.

### Product Objectives
- Customers MUST complete plan purchase in a simple 3-step flow: provider -> plan -> payment.
- Admin MUST fully manage providers, plans, and service request statuses.
- Platform MUST support payment integration and SMS notifications.
- Platform MUST enforce reCAPTCHA v3 verification on payment submission.

## Non-Goals (Out of Scope)
- Advanced performance optimization and aggressive latency SLAs.
- Multi-tenant or multi-operator support.
- Public partner/developer APIs.
- Enterprise compliance programs (GDPR/HIPAA/SOC2) for v1.
- Advanced analytics dashboards, alerting systems, or incident automation.
- Multi-environment release orchestration (dev/stage/prod matrix) in v1.

## Personas

### Primary Persona - Customer
- Village end-user purchasing DTH plans.
- Needs a clear, low-friction purchase journey.
- Success = completes payment and gets service request confirmation.

### Secondary Persona - Admin Operator
- Single local operator managing business operations.
- Needs CRUD for providers/plans and service request status management.
- Success = accurate catalog + complete request tracking for manual activation.

## Critical User Journeys (end-to-end, with entry/exit criteria)

### CJ-001 Customer Plan Purchase
- **Entry criteria:** Customer accesses customer web app landing page.
- **Main flow:** Select provider -> select plan -> enter required details -> pass reCAPTCHA v3 check -> initiate payment -> receive success/failure outcome -> service request recorded with correct status.
- **Exit criteria (success):** Service request created with payment success status and confirmation shown; SMS dispatched (if configured).
- **Exit criteria (failure):** Service request recorded with payment failed status and retry path shown.

### CJ-002 Admin Provider Management
- **Entry criteria:** Admin authenticated in admin app.
- **Main flow:** Navigate to Providers -> create/edit/delete/list providers; optional CSV upload for bulk updates.
- **Exit criteria:** Provider catalog reflects intended changes and is available to customer app.

### CJ-003 Admin Plan Management
- **Entry criteria:** Admin authenticated in admin app.
- **Main flow:** Navigate to Plans -> create/edit/delete/list plans linked to providers; optional CSV upload for bulk updates.
- **Exit criteria:** Active plan catalog reflects intended changes and is available to customer app.

### CJ-004 Admin Service Request Operations
- **Entry criteria:** Admin authenticated and viewing service requests.
- **Main flow:** Review request details -> update status (Pending, payment failed, payment success, blocked, completed) according to manual activation progress.
- **Exit criteria:** Service request lifecycle accurately reflects latest operational state.

## Feature Requirements (FR-001...)

### FR-001 Customer Provider Discovery
- **Description:** System MUST display a list of available DTH providers to customers.
- **Rationale:** Provider selection is the first purchase step.
- **Priority:** P0
- **Acceptance criteria:**
  - Customer landing page shows active providers with name and image.
  - Only providers created in admin app are shown.
  - Empty-state message appears when no providers exist.

### FR-002 Customer Plan Listing by Provider
- **Description:** System MUST show plans filtered by selected provider.
- **Rationale:** Prevents invalid plan selection and simplifies UX.
- **Priority:** P0
- **Acceptance criteria:**
  - Selecting a provider returns only plans for that provider.
  - Only active plans (`is_active=true`) appear for customers.
  - Plan card includes name, description, price, discount.

### FR-003 Customer Payment Initiation
- **Description:** System MUST allow customer to initiate payment for a selected plan.
- **Rationale:** Core transaction path.
- **Priority:** P0
- **Acceptance criteria:**
  - Payment request is created only after plan selection.
  - Service request record is created/updated with associated user and plan.
  - Payment attempt outcome is persisted.

### FR-004 Payment Page reCAPTCHA v3 Enforcement
- **Description:** Payment submission MUST validate Google reCAPTCHA v3 before processing.
- **Rationale:** Reduce abusive/bot payment attempts.
- **Priority:** P0
- **Acceptance criteria:**
  - Token is required on payment submit.
  - Backend verification failure blocks payment initiation.
  - User receives actionable error on verification failure.

### FR-005 Customer Account/Identity Capture
- **Description:** System MUST persist customer user record with required fields.
- **Rationale:** Needed for service request tracking and operations.
- **Priority:** P0
- **Acceptance criteria:**
  - User entity stores `id, name, phone_number, status, role, password, created_at, updated_at`.
  - Role for customer records is `customer`.
  - Phone number format validation is applied.

### FR-006 Admin Authentication
- **Description:** Admin app MUST support login via JWT-based authentication.
- **Rationale:** Protect operational functions.
- **Priority:** P0
- **Acceptance criteria:**
  - Valid credentials return JWT token.
  - Invalid credentials are rejected with standard error response.
  - Protected admin endpoints require valid JWT.

### FR-007 Provider CRUD (Admin)
- **Description:** Admin MUST create, read, update, and delete providers.
- **Rationale:** Maintain purchasable provider catalog.
- **Priority:** P0
- **Acceptance criteria:**
  - CRUD endpoints/UI are available in admin app.
  - Provider fields include `id, name, image_url, created_at, updated_at`.
  - Changes are visible in customer provider list after persistence.

### FR-008 Plan CRUD (Admin)
- **Description:** Admin MUST create, read, update, and delete plans.
- **Rationale:** Maintain purchasable plan catalog.
- **Priority:** P0
- **Acceptance criteria:**
  - CRUD endpoints/UI are available in admin app.
  - Plan fields include `id, provider_id, name, description, price, discount, is_active`.
  - `provider_id` MUST reference an existing provider.

### FR-009 Service Request Management (Admin)
- **Description:** Admin MUST view and update service request status.
- **Rationale:** Supports manual activation workflow.
- **Priority:** P0
- **Acceptance criteria:**
  - Admin can list and inspect service requests.
  - Status updates limited to allowed values: Pending, payment failed, payment success, blocked, completed.
  - Status changes are persisted with updated timestamp.

### FR-010 SMS Notifications
- **Description:** System SHOULD send SMS notifications for key request/payment events.
- **Rationale:** Improves customer communication in low-support environment.
- **Priority:** P1
- **Acceptance criteria:**
  - SMS trigger events are configurable for payment success/failure.
  - Failure to send SMS is logged without blocking main transaction.
  - Message delivery attempt result is recorded.

### FR-011 CSV Upload for Providers
- **Description:** Admin MUST upload provider data via CSV.
- **Rationale:** Faster bulk onboarding and updates.
- **Priority:** P1
- **Acceptance criteria:**
  - CSV file parsing validates required provider fields.
  - Invalid rows are reported with row-level errors.
  - Valid rows are persisted.

### FR-012 CSV Upload for Plans
- **Description:** Admin MUST upload plan data via CSV.
- **Rationale:** Faster catalog management.
- **Priority:** P1
- **Acceptance criteria:**
  - CSV parsing validates required plan fields and provider linkage.
  - Invalid rows are rejected with row-level feedback.
  - Valid rows are persisted.

### FR-013 REST API Standards
- **Description:** Backend MUST expose REST APIs with URI versioning.
- **Rationale:** Consistent contract and maintainability.
- **Priority:** P0
- **Acceptance criteria:**
  - All endpoints are under `/api/v1`.
  - API responses use a consistent success/error envelope.
  - Endpoint contracts are documented in OpenAPI.

### FR-014 Data Migration and Seed Support
- **Description:** System MUST support release-based DB migrations and admin seed command.
- **Rationale:** Safe schema evolution and operability.
- **Priority:** P0
- **Acceptance criteria:**
  - Migration mechanism supports ordered release execution.
  - Seed command creates admin user with predefined credentials (as currently requested).
  - Migration and seed operations are executable from CLI/scripts.

### FR-015 Reliability and Backup Operations
- **Description:** System MUST support manual DB backup scripts and target 90% uptime.
- **Rationale:** Minimum operational resilience for production use.
- **Priority:** P0
- **Acceptance criteria:**
  - Backup script exists and runs successfully against production DB.
  - Restore procedure is documented and manually testable.
  - Service uptime is monitored at basic availability level.

### FR-016 Logging and Runtime Operations
- **Description:** Backend MUST run as systemd service with logs available via journalctl.
- **Rationale:** Matches stated operations model.
- **Priority:** P0
- **Acceptance criteria:**
  - Service unit starts/stops/restarts API process.
  - Application logs are visible via `journalctl`.
  - Critical errors are logged with request context.

### FR-017 Backend Unit Test Coverage Baseline
- **Description:** Backend MUST include unit tests for core business logic.
- **Rationale:** Required quality gate in context.
- **Priority:** P0
- **Acceptance criteria:**
  - Unit tests exist for auth, plan/provider CRUD logic, and service-request transitions.
  - Tests run in CI/local command with pass/fail output.
  - No launch blocker defects in tested core flows.

## Business Rules
1. Customer purchase flow order MUST be provider selection -> plan selection -> payment.
2. Customer app and admin app are separate applications; RBAC across one shared UI is out of scope.
3. Service request statuses MUST be constrained to: Pending, payment failed, payment success, blocked, completed.
4. Customer-facing plan listing MUST exclude inactive plans.
5. Plan MUST reference a valid provider (`provider_id` integrity).
6. Payment flow MUST enforce reCAPTCHA v3 verification before initiating payment.
7. Admin functions MUST require valid JWT authentication.
8. CSV imports MUST validate input and report row-level errors for invalid data.
9. SMS failures MUST NOT block payment/service-request persistence.
10. Manual operator action remains part of fulfillment; “completed” status indicates manual activation completion.

## Success Metrics

| Metric ID | Metric | Baseline | Target | Measurement Method | Reporting Cadence |
|---|---|---:|---:|---|---|
| SM-001 | Active customers | 0 (digital baseline) | 30-50 within 6-12 months | Count unique active customer users/month | Monthly |
| SM-002 | Activation rate | Unknown pre-launch | >=20% | New users completing first successful purchase / new users | Monthly |
| SM-003 | Retention rate | Unknown pre-launch | >=25% | Returning customers in defined cohort window | Monthly |
| SM-004 | Purchase flow completion | Unknown pre-launch | Baseline to be established in month 1 | Funnel: provider -> plan -> payment success | Weekly |
| SM-005 | Service availability | Not currently measured | >=90% uptime | Service health availability checks/log-based tracking | Monthly |
| SM-006 | Operational stability | No benchmark | Near-zero critical glitches | Count production incidents causing purchase/admin blocking | Monthly |

## Risks and Mitigations
- **Payment integration failure risk:** Use idempotent request handling, retry-safe callbacks, and manual reconciliation runbook.
- **SMS delivery unreliability:** Make SMS non-blocking, log failures, and enable resend workflow.
- **Single-developer capacity risk:** Prioritize P0 scope only; defer non-critical enhancements.
- **Hardcoded seed credential risk:** Restrict to initial bootstrap only; force immediate credential rotation post-deploy.
- **Limited observability risk:** Standardize structured logs and minimal health checks despite no full monitoring stack.
- **Data quality risk from CSV:** Enforce schema validation with row-level error reporting and safe partial processing.
- **Single production environment risk:** Maintain tested backup/restore scripts and pre-release checklist.

## Assumptions (explicit)
1. One business operator (single admin organization) is sufficient for v1.
2. Customer authentication/identity approach is minimal and acceptable for this audience.
3. Payment provider supports required API flow from India/village usage context.
4. SMS cost and provider availability are acceptable within budget.
5. Low scale (30-50 active users) allows single production deployment without complex scaling.
6. Backend unit tests are the only mandatory test gate for launch.
7. No formal legal/compliance requirements apply in v1 beyond basic secure handling.
8. Manual fulfillment by operator continues after payment; product does not automate DTH activation.

## Open Questions (explicit, with owner + due date suggestion)
| OQ ID | Question | Owner | Due Date (Suggested) | Impact if Unresolved |
|---|---|---|---|---|
| OQ-001 | Which payment gateway will be used (provider, callback model, settlement workflow)? | Product Owner | 2026-03-05 | Blocks payment API and reconciliation design |
| OQ-002 | What exact SMS events and templates are required (success, failure, pending)? | Product Owner | 2026-03-05 | Delays notification implementation and UX messaging |
| OQ-003 | What customer authentication model is required (guest checkout vs login)? | Product Owner | 2026-03-06 | Affects user lifecycle, JWT scope, and UX |
| OQ-004 | What are acceptable reCAPTCHA v3 score thresholds and fallback behavior? | Engineering Owner | 2026-03-06 | Affects fraud controls vs false positives |
| OQ-005 | Are deletions soft-delete or hard-delete for providers/plans with historical requests? | Engineering Owner | 2026-03-06 | Impacts DB model and data integrity |
| OQ-006 | What are CSV file limits and duplicate-handling rules? | Product Owner | 2026-03-07 | Impacts import performance and data correctness |
| OQ-007 | What exactly constitutes “active customer” for KPI reporting? | Product Owner | 2026-03-07 | Required for metric correctness |
| OQ-008 | What credential rotation policy applies after seeded admin creation? | Engineering Owner | 2026-03-05 | Security and operational risk |

### Answers:
OQ-001: Razorpay
OQ-002: Send SMS only when the admin has marked a Service Requests a completed. Here is the template: `Your plan <plan name> is now active. Here are the benefits: <Plan description>`
OQ-003: Guest checkout. But we need reCaptcha for payment page.
OQ-004: reCAPTCHA v3 score thresholds be either equal or more than 0.5. If failed then throw an error in the screen `Please try after some time`
OQ-005: The deletions are soft-delete GORM has built in features to handle this
OQ-006: Two plans for a the same provider cannot have the same name. Two providers cannot have the same name. There cannot be more than 100 plans for a provider
OQ-007: An active customer is the one who does recurring payments at the end of their plans
OQ-008: The credentials rotation policy is not needed it can be manually handled

## Traceability Matrix
| PRD ID | Source Context Item | Rationale | Downstream Artifact |
|---|---|---|---|
| FR-001 | Customer chooses provider first | Enables step 1 of purchase funnel | FSD customer journey flow, OpenAPI provider list endpoint |
| FR-002 | Customer chooses plan after provider | Ensures valid provider-plan linkage | FSD validation rules, OpenAPI plan list by provider |
| FR-003 | Customer makes payment | Core business transaction | Payment service design, OpenAPI payment endpoint |
| FR-004 | Payment page has reCAPTCHA v3 | Bot-abuse control requirement | Security architecture, backend verification handler |
| FR-005 | User entity fields provided | Core data capture for requests | DB schema, OpenAPI User schema |
| FR-006 | JWT auth preferred | Admin access control | Auth middleware, OpenAPI security scheme |
| FR-007 | Admin CRUD providers | Catalog operations requirement | Admin UI modules, provider CRUD endpoints |
| FR-008 | Admin CRUD plans | Catalog operations requirement | Admin UI modules, plan CRUD endpoints |
| FR-009 | Admin manages service requests/status | Manual fulfillment workflow | Service request module, status transition logic |
| FR-010 | SMS required | Communication requirement | Async job + SMS integration spec |
| FR-011 | CSV upload providers | Bulk admin operations | CSV parser/import endpoint |
| FR-012 | CSV upload plans | Bulk admin operations | CSV parser/import endpoint |
| FR-013 | REST + URI versioning | API standardization requirement | `/api/v1` contract, OpenAPI versioned paths |
| FR-014 | Migration scripts by release + seed admin | Operability requirement | Migration runner + seed command |
| FR-015 | 90% uptime + manual DB backups | Reliability baseline | Backup scripts, ops runbook |
| FR-016 | systemd + journalctl logging | Runtime operations model | Deployment unit file + logging format |
| FR-017 | Backend unit tests only | Required test gate | Test suite plan, CI/local test command |
| SM-001..006 | Business goals + metrics stated | KPI-driven product management | Analytics/reporting checklist |
| CJ-001..004 | Stated critical journeys | End-to-end behavior definition | FSD flows, QA scenarios |

## Decision Log
| Decision ID | Date | Decision | Options Considered | Chosen Option | Rationale | Impact |
|---|---|---|---|---|---|---|
| D-001 | 2026-02-27 | Product split into two web apps | Single app with RBAC; separate apps | Separate customer and admin apps | Matches explicit context, simpler UX separation | Clearer boundaries, simpler role model |
| D-002 | 2026-02-27 | API style and versioning | GraphQL; REST unversioned; REST versioned | REST with URI versioning (`/api/v1`) | Explicit preference in context | Stable contract management |
| D-003 | 2026-02-27 | Auth model | Session; OAuth2; JWT | JWT | Explicit context requirement | Consistent admin API security |
| D-004 | 2026-02-27 | Payment fraud protection | None; CAPTCHA v2; reCAPTCHA v3 | reCAPTCHA v3 | Explicit payment-page requirement | Added backend verification dependency |
| D-005 | 2026-02-27 | Infrastructure baseline | Multi-env CI/CD; single prod | Single production environment | Matches current delivery constraints | Lower overhead, higher operational risk |
| D-006 | 2026-02-27 | Test gate for v1 | Full pyramid; backend unit only | Backend unit tests mandatory | Explicit constraint from context | Faster delivery, reduced test coverage |
| D-007 | 2026-02-27 | Reliability posture | High SLA; minimal target | 90% uptime + manual backups | Explicit target and ops preference | Basic operability with low complexity |
| D-008 | 2026-02-27 | Admin bootstrapping | Manual DB insert; seed command | Seed command with initial admin credentials | Explicit setup requirement | Faster setup, credential management risk |
```

Approved Functional Spec:
```
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

```

Approved Technical Architecture:
```
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
```

Output format:
1) OpenAPI YAML only (fenced code block with `yaml`)
2) Then a short Markdown appendix with:
   - Assumptions
   - Open Questions
   - Traceability Matrix (OperationId/Schema -> PRD ID/FS ID/Arch ID)
   - Risks
   - Decision Log table

Spec requirements:
- OpenAPI version: 3.1.0
- Base path includes `/api/v1`
- Security scheme: Bearer JWT
- Consistent error schema and HTTP status conventions
- Paths for:
  - Auth/login (admin)
  - Providers CRUD (+ CSV upload if API-based)
  - Plans CRUD (+ CSV upload if API-based)
  - Customer browse providers/plans
  - Create payment/service request flow
  - Payment status/callback handling (if applicable)
  - Service request listing/update for admin
- Components:
  - Schemas for User, Provider, Plan, ServiceRequest, PaymentRequest/Response, CSV upload result, Error
- Include:
  - Request validation constraints
  - Enum definitions for status values
  - Pagination/filter params where needed
  - Idempotency header for payment-related create operation
  - OperationId for every endpoint
  - Examples for key request/response payloads

Validation checklist:
- YAML parses as valid OpenAPI 3.1.
- Every operation links to a functional or product requirement.
- Schemas align with stated entities and statuses.
- Security requirements applied consistently.
- Error responses standardized across endpoints.
- Assumptions and unresolved questions are explicit.

Quality bar:
- Contract can be directly used to scaffold server/client.
- No ambiguous field types or missing required constraints.
- Production-quality naming, status codes, and schema hygiene.