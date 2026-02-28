You are a Senior Systems Analyst and Functional Architect.
Generate a PRODUCTION-READY Functional Specification (FSD) from the inputs below.
Output must be build-ready, not a draft.

Inputs:
1) Project Context Pack
2) Approved PRD Markdown
3) Any additional constraints in this prompt

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

Output format (Markdown):
- Title + Metadata (version, author, linked PRD version)
- Functional Scope Summary
- Actors and Permissions Model
- Functional Decomposition by Module
  - Customer app
  - Admin app
  - Shared backend services
- Detailed Functional Requirements (FS-001...)
  - Trigger
  - Preconditions
  - Main Flow
  - Alternate Flows
  - Postconditions
  - Error handling
  - Business validations
- State Models
  - Service Request lifecycle (Pending, payment failed, payment success, blocked, completed)
- UI/UX Functional Behavior (not visual design)
  - Field-level validation rules
  - Form behaviors
  - Pagination/sorting/filtering rules (where applicable)
- Notifications Behavior (SMS events, retries, failure handling)
- CSV Upload Functional Rules
  - Providers upload
  - Plans upload
  - Validation and partial failure handling
- Traceability Matrix
  - Columns: FS ID, PRD ID, Context Source, Test Case Hint
- Assumptions
- Open Questions
- Risks
- Decision Log table

Validation checklist:
- Every FS item links to at least one PRD item.
- All critical journeys are fully decomposed into main + alternate flows.
- Error cases are explicit for payment, CSV upload, auth, and admin CRUD.
- Service Request state transitions are complete and non-contradictory.
- Assumptions and unresolved questions are explicit.
- No requirement conflicts with PRD/context.

Quality bar:
- Function-level precision; no vague prose.
- Directly implementable by engineering and testable by QA.
- Explicit handling for real-world failures (payment callback failure, duplicate requests, invalid CSV rows).
- Practical for a solo developer delivery model.