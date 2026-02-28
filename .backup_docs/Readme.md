# Vibe Coding Documentation Framework

This framework is a prompt-driven documentation factory for production applications.

## Framework

- Use a 3-step flow: `Project Intake` -> `Prompt Generator` -> `Document Generator`.
- The `Prompt Generator` creates tailored prompts for each document.
- The `Document Generator` outputs final docs with strict production constraints.
- Keep one shared `Source of Truth` block (product goals, users, NFRs, API style, stack) and inject it into every prompt to avoid drift.
- Generate docs in this order: `PRD` -> `Functional Spec` -> `Technical Architecture` -> `OpenAPI` -> `Implementation Plan` -> `Developer Setup`.

```text
Order matters: each later document must cite and reconcile earlier documents.
```

## 1) Project Intake Prompt (run first)

```markdown
You are a Senior Product+Engineering Discovery Assistant.

Goal: Build a complete production-ready documentation baseline for a new application.

Ask me targeted questions to collect:
1) Product vision, business goals, success metrics
2) Primary/secondary user personas and critical journeys
3) Scope: in-scope/out-of-scope
4) Functional requirements and business rules
5) Non-functional requirements (performance, reliability, security, compliance, privacy, accessibility)
6) Data model expectations and integrations
7) API style (REST/GraphQL), auth model, versioning strategy
8) Technical constraints (cloud, stack, team size, release cadence, budget)
9) Environment strategy (local/dev/stage/prod)
10) Observability and operations expectations

Output:
- A normalized "Project Context Pack" in Markdown
- "Assumptions" section
- "Open Questions" section
- "Risks" section
- "Decision Log" table
```
- Paste the output to a file under @Docs

## 2) Prompt Generator (meta-prompt to generate prompts for each document)

```markdown
You are a Prompt Architect.
Using the Project Context Pack below, generate 6 high-quality prompts that will each produce one production-grade document:
1) Production PRD
2) Functional Specification
3) Production Technical Architecture
4) OpenAPI Specification
5) Implementation Plan
6) Developer Setup Guide

Requirements for generated prompts:
- Each prompt must enforce traceability to previous docs.
- Each prompt must include "Inputs", "Output format", "Validation checklist", "Quality bar".
- Each prompt must force explicit assumptions and unresolved questions.
- Each prompt must output production-ready content, not drafts.
- Implementation Plan prompt must include:
  - task `status` field
  - design pattern per task
  - explicit files/folders to create/modify
  - strict mapping to OpenAPI operations/schemas/security

Return format:
- Section per document:
  - Prompt Name
  - Copy-paste Prompt
  - Expected Output Structure
  - Common Failure Modes to avoid

Project Context Pack:
<<<PASTE CONTEXT PACK>>>
```

## 3) Document Generator Prompts (ready-to-use templates)

### A) Production PRD Prompt

```markdown
You are a Principal Product Manager.
Generate a Production PRD using the context and constraints.

Must include:
- Executive summary
- Problem statement
- Goals / non-goals
- Personas
- User journeys
- Feature requirements (prioritized)
- Acceptance criteria (testable)
- NFRs (latency, uptime, scalability, security, accessibility)
- Risks and mitigations
- Milestones and release slices
- KPI framework and instrumentation events

Output quality bar:
- No ambiguity in requirements
- Every feature has measurable acceptance criteria
- Conflicts called out with decisions proposed

Inputs:
- Project Context Pack
```

### B) Functional Specification Prompt

```markdown
You are a Senior Business Analyst + Systems Analyst.
Generate a Functional Specification aligned to the PRD.

Must include:
- Functional decomposition by module
- Detailed use cases (preconditions, triggers, main flow, alternate flows, errors)
- Business rules and validation rules
- State transitions and lifecycle rules
- Role/permission behavior
- Edge-case behavior matrix
- Requirement-to-use-case traceability table

Constraints:
- Every requirement must map back to PRD item IDs.
- No implementation details unless needed for deterministic behavior.

Inputs:
- PRD
- Project Context Pack
```

### C) Production Technical Architecture Prompt

```markdown
You are a Staff Architect.
Generate a production-ready Technical Architecture document aligned to PRD + Functional Spec.

Must include:
- Architecture style and rationale
- C4-level breakdown (context/container/component)
- Service boundaries and responsibilities
- Data architecture (primary stores, caching, indexing, retention)
- Security architecture (authn/authz, secrets, encryption, threat model summary)
- Reliability strategy (SLOs, retries, circuit breakers, idempotency, DR)
- Deployment architecture (envs, CI/CD, rollback)
- Observability architecture (logs, metrics, traces, alerting)
- Capacity planning assumptions
- Tradeoff decisions with ADR-style notes

Output:
- Include explicit technology choices + why
- Include risk register and scaling path
```

### D) OpenAPI Specification Prompt

```markdown
You are a Senior API Architect.
Generate an OpenAPI 3.1 specification in valid YAML for the system.

Must include:
- servers, tags, paths, operations
- request/response schemas
- validation constraints
- pagination/filter/sort conventions
- standard error model
- auth/security schemes
- idempotency and rate-limit headers where needed
- examples for major operations

Quality constraints:
- Consistent naming conventions
- Reusable components/schemas
- Operation IDs are stable and descriptive
- Spec supports implementation and client generation

Also output:
- Endpoint-to-requirement traceability table
- Breaking-change policy notes
```

### E) Implementation Plan Prompt (custom requirements baked in)

```markdown
You are an Engineering Delivery Lead.
Generate a production implementation plan strictly aligned to:
1) PRD
2) Functional Specification
3) Technical Architecture
4) OpenAPI Spec

Output as a phased task plan with dependencies.

Each task MUST include this schema:
- task_id: string
- title: string
- objective: string
- status: one of [todo, in_progress, blocked, in_review, done]
- priority: one of [P0, P1, P2]
- phase: string
- depends_on: [task_id...]
- mapped_requirements: [PRD/FS IDs]
- mapped_openapi:
  - operation_ids: [string...]
  - schemas: [string...]
  - security_schemes: [string...]
- design_pattern: string (e.g., Repository, Strategy, CQRS, Adapter, Saga, Factory, Hexagonal)
- implementation_details:
  - files_to_create: [path...]
  - files_to_modify: [path...]
  - folders_to_create: [path...]
- test_plan:
  - unit
  - integration
  - contract
  - e2e
- definition_of_done: [checklist]
- estimate: {story_points: int, ideal_days: number}
- risks: [string...]

Critical rules:
- No task without explicit file/folder impact.
- No backend/API task without OpenAPI operation_id mapping.
- Include DB migrations and rollback tasks where applicable.
- Include security, observability, and performance tasks.
- Include release/readiness and runbook tasks.

Final sections:
- Dependency graph (textual)
- Critical path
- Sprint slicing recommendation
```

### F) Developer Setup Guide Prompt

```markdown
You are a Senior Developer Experience Engineer.
Generate a production-grade Developer Setup Guide for onboarding in <30 minutes.

Must include:
- Prerequisites (versions pinned)
- Local environment setup
- Secrets/config management pattern
- Database/bootstrap steps
- Running app/services locally
- Running tests/lint/type checks
- Seeding data and fixtures
- Debugging tips
- Common issues + fixes
- CI parity checks
- Branching/commit/PR conventions

Output constraints:
- Commands must be copy-paste ready
- Works for clean machine setup
- Includes verification checkpoints after each major step
```

## 4) Orchestrator Prompt (generate all docs in sequence)

```markdown
You are a Documentation Orchestrator.
Generate all documents in this exact sequence:
1) PRD
2) Functional Spec
3) Technical Architecture
4) OpenAPI 3.1
5) Implementation Plan
6) Developer Setup Guide

Rules:
- Each document must reference IDs from prior documents.
- If conflicts appear, stop and emit "Conflict Report" + proposed resolution, then continue with resolved assumptions.
- Implementation Plan tasks must map to OpenAPI operationIds/schemas and include status, design pattern, and explicit files/folders.
- Produce a final "Traceability Matrix" linking:
  - PRD -> Functional Spec -> Architecture -> OpenAPI -> Implementation Tasks

Inputs:
<<<PROJECT CONTEXT PACK>>>
```
