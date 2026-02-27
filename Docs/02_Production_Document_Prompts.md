## 1) Prompt Name
**Production PRD Generator (Traceable to Context Pack)**

### Copy-paste Prompt
```text
You are a Principal Product Manager + Product Writer.
Generate a PRODUCTION-READY Product Requirements Document (PRD) for the product described below.
Do NOT produce a draft. Do NOT leave placeholders except in the explicit "Open Questions" section.

Inputs:
1) Project Context Pack (authoritative source)
2) Any stated constraints in this prompt
3) If context is missing, capture it under Assumptions and Open Questions (do not block output)

Project Context Pack:
<<<PASTE PROJECT CONTEXT PACK HERE>>>

Output format (Markdown):
- Title + Metadata (version, date, owner, status=Approved-for-Build)
- Executive Summary
- Problem Statement and Opportunity
- Product Vision and Goals
- Non-Goals (Out of Scope)
- Personas
  - Primary persona
  - Secondary persona
- Critical User Journeys (end-to-end, with entry/exit criteria)
- Feature Requirements (FR-001...)
  - Each FR includes: description, rationale, priority (P0/P1/P2), acceptance criteria
- Business Rules
- Success Metrics
  - Baseline, target, measurement method, reporting cadence
- Risks and Mitigations
- Assumptions (explicit)
- Open Questions (explicit, owner + due date suggestion)
- Traceability Matrix
  - Columns: PRD ID, Source Context Item, Rationale, Downstream Artifact
- Decision Log
  - Columns: Decision ID, Date, Decision, Options Considered, Chosen Option, Rationale, Impact

Validation checklist (must pass before finalizing):
- Every requirement is uniquely ID’d and testable.
- Every requirement maps back to context input (traceability present).
- Scope is explicit (in/out) and consistent with context.
- No contradiction between goals, metrics, and features.
- Assumptions are explicit and non-hidden.
- Open questions are actionable and assigned.
- Content is production-ready and implementation-usable today.

Quality bar:
- Clear, unambiguous language; no generic filler.
- Requirements are atomic, verifiable, and prioritized.
- Journeys are realistic for this exact DTH village-operator use case.
- Incorporates constraints: low-cost, solo developer, Oracle Cloud, separate admin and customer web apps.
- Includes reCAPTCHA v3 on payment flow and SMS requirement.
- Uses normative language: MUST/SHOULD/MAY.
```

### Expected Output Structure
- Full PRD with `FR-*` requirements, measurable success metrics, risks, assumptions/open questions, and a traceability matrix from context to PRD items.

### Common Failure Modes to avoid
- Vague requirements (“system should be user-friendly”).
- Missing IDs, missing acceptance criteria, or no priority.
- Ignoring explicit constraints (solo dev, Oracle, low-cost).
- No mapping from context to PRD requirements.


---

## 2) Prompt Name
**Production Functional Specification Generator (Traceable to PRD + Context)**

### Copy-paste Prompt
```text
You are a Senior Systems Analyst and Functional Architect.
Generate a PRODUCTION-READY Functional Specification (FSD) from the inputs below.
Output must be build-ready, not a draft.

Inputs:
1) Project Context Pack
2) Approved PRD Markdown
3) Any additional constraints in this prompt

Project Context Pack:
<<<PASTE CONTEXT PACK>>>

Approved PRD:
<<<PASTE PRD>>>

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
```

### Expected Output Structure
- A fully detailed FSD with actor behavior, flow logic, validations, state transitions, error handling, and strict traceability to PRD/context.

### Common Failure Modes to avoid
- Only rewriting PRD in different words.
- Missing alternate flows and failure cases.
- No field-level validation or lifecycle definition.
- No explicit traceability mapping.


---

## 3) Prompt Name
**Production Technical Architecture Generator (Traceable to PRD + FSD)**

### Copy-paste Prompt
```text
You are a Principal Solutions Architect.
Create a PRODUCTION-READY Technical Architecture Document (TAD) for implementation.
Do not output a draft. Make architecture decisions explicit.

Inputs:
1) Project Context Pack
2) Approved PRD
3) Approved Functional Specification
4) Tech constraints: Next.js (SSG) frontend, Go 1.26 + Gin + GORM backend, PostgreSQL, Redis, Asynq, Oracle Cloud, JWT auth, URI versioning, single prod env

Project Context Pack:
<<<PASTE CONTEXT PACK>>>

Approved PRD:
<<<PASTE PRD>>>

Approved Functional Spec:
<<<PASTE FSD>>>

Output format (Markdown):
- Title + Metadata
- Architecture Overview (C4 level-1/2 textual)
- Design Principles and Constraints
- System Components
  - Customer frontend
  - Admin frontend
  - API service
  - DB
  - Redis/Asynq workers
  - Payment integration
  - SMS integration
  - reCAPTCHA verification
- Data Architecture
  - Logical schema overview
  - Key indexes/constraints
  - Migration strategy per release
  - Seed strategy (admin seed command)
- API Architecture
  - Versioning strategy (`/api/v1`)
  - AuthN/AuthZ model
  - Error envelope standard
- Async Processing Architecture
  - Job types, retry policy, idempotency, dead-letter handling
- Security Architecture
  - JWT handling, password storage, secret management, transport security
- Reliability & Operations
  - 90% uptime target approach
  - Backup/restore manual scripts strategy
  - systemd + journalctl logging model
- Deployment Topology (single prod)
- Testing Strategy (backend unit focus + critical API contract tests)
- Traceability Matrix
  - Columns: Architecture Element ID, PRD ID, FS ID, Rationale, Verification Method
- Assumptions
- Open Questions
- Risks and Tradeoffs
- Decision Log table

Validation checklist:
- All architecture choices map to PRD/FSD requirements.
- Includes concrete handling for payment, SMS, CSV upload, and service request states.
- Defines migration + seed process and release-time migration execution.
- Defines idempotency and failure handling for external integrations.
- Explicit assumptions and unresolved questions included.
- No contradiction with stated stack and constraints.

Quality bar:
- Actionable and implementation-specific, not conceptual fluff.
- Explicit interfaces, boundaries, and responsibilities.
- Production-conscious despite lightweight NFR posture.
- Suitable for handoff to a solo engineer for immediate build.
```

### Expected Output Structure
- End-to-end architecture with component boundaries, data/API/async/security/ops details and traceability to PRD/FSD.

### Common Failure Modes to avoid
- High-level architecture with no implementation details.
- Missing migration/seed and backup strategy.
- No idempotency/retry policy for payment/SMS.
- Ignoring traceability to PRD/FSD.


---

## 4) Prompt Name
**OpenAPI 3.1 Specification Generator (Traceable to PRD + FSD + Architecture)**

### Copy-paste Prompt
```text
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
<<<PASTE CONTEXT PACK>>>

Approved PRD:
<<<PASTE PRD>>>

Approved Functional Spec:
<<<PASTE FSD>>>

Approved Technical Architecture:
<<<PASTE TAD>>>

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
```

### Expected Output Structure
- Valid OpenAPI YAML plus a concise appendix containing assumptions/open questions/traceability/risks/decision log.

### Common Failure Modes to avoid
- Invalid OpenAPI syntax or inconsistent schema refs.
- Missing `operationId`, security, or error schema.
- No enum constraints for service request status.
- No traceability mapping from operations/schemas to upstream docs.


---

## 5) Prompt Name
**Implementation Plan Generator (Strictly Mapped to OpenAPI + Schemas + Security)**

### Copy-paste Prompt
```text
You are a Staff Engineer and Delivery Planner.
Generate a PRODUCTION-READY implementation plan, not a draft backlog.

Inputs:
1) Project Context Pack
2) Approved PRD
3) Approved Functional Spec
4) Approved Technical Architecture
5) Approved OpenAPI spec (authoritative API contract)

Project Context Pack:
<<<PASTE CONTEXT PACK>>>

Approved PRD:
<<<PASTE PRD>>>

Approved Functional Spec:
<<<PASTE FSD>>>

Approved Technical Architecture:
<<<PASTE TAD>>>

Approved OpenAPI:
<<<PASTE OPENAPI YAML>>>

Output format (Markdown):
- Title + Metadata
- Delivery assumptions and sequencing constraints
- Milestones (M1, M2, ...)
- Task Table (REQUIRED COLUMNS):
  - task_id
  - title
  - status (todo | in_progress | blocked | done)
  - priority (P0/P1/P2)
  - owner_role (solo-dev default acceptable)
  - design_pattern (e.g., Repository, Adapter, Strategy, Factory, Middleware, Saga-lite, Outbox-lite)
  - api_mapping (operationId(s) + schema(s) + security requirement mapping)
  - files_to_create (explicit paths)
  - files_to_modify (explicit paths)
  - dependencies (task_id list)
  - acceptance_criteria (testable)
  - test_plan (unit/integration/contract)
  - risk_notes
- Folder/File Blueprint
  - Backend
  - Frontend customer
  - Frontend admin
  - Infra/scripts/migrations
- Migration & Seed Execution Plan
- Rollback Plan per milestone
- Traceability Matrix
  - task_id -> PRD ID -> FS ID -> Arch ID -> OpenAPI operationId/schema/security
- Assumptions
- Open Questions
- Risks
- Decision Log table

Mandatory constraints:
- Every implementation task MUST map to one or more OpenAPI operations OR schemas OR security schemes.
- Every task MUST include explicit files/folders to create/modify.
- Use only the declared stack and deployment model.
- Include tasks for payment integration, SMS, reCAPTCHA v3 validation, CSV import, migrations, seed admin command, logging via journalctl-compatible output.
- Include test tasks at minimum for backend unit tests and API contract checks.

Validation checklist:
- No orphan tasks without upstream requirement and API mapping.
- `status` field present and valid for every task.
- `design_pattern` specified for every task.
- Files/folders are concrete and plausible.
- Dependency graph is coherent (no impossible ordering).
- Assumptions and unresolved questions explicitly listed.

Quality bar:
- Execution-ready for immediate coding.
- Granular enough for daily progress tracking.
- Strict traceability across all artifacts.
- Realistic for solo developer throughput.
```

### Expected Output Structure
- Milestone-based, dependency-aware execution plan with task-level status, design pattern, explicit file paths, and strict mapping to OpenAPI operations/schemas/security.

### Common Failure Modes to avoid
- Generic task list without file-level impact.
- Missing `status` or `design_pattern`.
- No strict mapping to OpenAPI operations/schemas/security.
- No dependency order or no rollback/testing plan.


---

## 6) Prompt Name
**Production Developer Setup Guide Generator (Traceable to Architecture + Implementation Plan)**

### Copy-paste Prompt
```text
You are a Senior Developer Experience (DevEx) Engineer.
Generate a PRODUCTION-READY Developer Setup Guide for this project.
This is not a quickstart draft; it must be complete and operational.

Inputs:
1) Project Context Pack
2) Approved Technical Architecture
3) Approved OpenAPI
4) Approved Implementation Plan

Project Context Pack:
<<<PASTE CONTEXT PACK>>>

Approved Technical Architecture:
<<<PASTE TAD>>>

Approved OpenAPI:
<<<PASTE OPENAPI>>>

Approved Implementation Plan:
<<<PASTE IMPLEMENTATION PLAN>>>

Output format (Markdown):
- Title + Metadata
- Prerequisites (exact versions/tools)
- Repository Structure Overview
- Local Setup
  - Clone/install
  - Env vars (`.env.example` contract)
  - DB setup (PostgreSQL)
  - Redis setup
  - Migration commands
  - Seed command (admin user)
- Runbook Commands
  - Backend run/test/lint
  - Frontend customer run/build
  - Frontend admin run/build
  - Worker process (Asynq) run
- Integration Setup
  - Payment sandbox config
  - SMS provider config
  - reCAPTCHA v3 config
- API Contract Workflow
  - OpenAPI validation, codegen (if used), contract testing
- Troubleshooting Guide (top 10 likely issues + fixes)
- Production-like Local Checks before merge
- Security and Secrets Handling Guidelines
- Traceability Matrix
  - Setup step -> Architecture section -> Implementation task_id -> OpenAPI relevance
- Assumptions
- Open Questions
- Risks
- Decision Log table

Validation checklist:
- Steps are executable end-to-end on a fresh machine.
- Commands and env vars are concrete and consistent with architecture.
- Includes migration + seed + rollback verification basics.
- Includes integration setup for payment/SMS/reCAPTCHA.
- Every major setup section traces to architecture and implementation tasks.
- Assumptions and unresolved questions explicit.

Quality bar:
- Zero ambiguity in critical steps.
- Suitable for onboarding a new engineer with minimal assistance.
- Operationally safe defaults and clear troubleshooting.
- Production-minded despite single-prod environment.
```

### Expected Output Structure
- Complete setup/run/troubleshooting handbook with concrete commands, env contracts, and traceability to architecture and implementation tasks.

### Common Failure Modes to avoid
- Incomplete setup steps requiring tribal knowledge.
- Missing env var definitions or integration configuration.
- No migration/seed workflow.
- No traceability back to architecture and implementation tasks.

---

If you want, I can also generate a **recommended execution order** and a **document dependency graph** so these prompts are run in the safest sequence with minimal rework.
