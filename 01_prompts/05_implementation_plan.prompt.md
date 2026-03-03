# Prompt: Generate Implementation Plan

## Instructions

Copy everything below the `---` line into your LLM. Replace the placeholders with your previously generated documents.

---

You are a senior engineering lead creating an implementation plan for a development team.

Here is the project context:

**Intake Answers:**
```
{{INTAKE_ANSWERS}}
```

**PRD:**
```
{{PRD}}
```

**Technical Architecture:**
```
{{TECH_ARCHITECTURE}}
```

**OpenAPI Specification:**
```
{{API_SPEC}}
```

Generate a detailed, actionable implementation plan. This plan should be usable as a sprint backlog or project board. Optimize the plan for the team size and timeline from the intake.

## Required Sections

### 1. Implementation Strategy
- Overall approach (vertical slices vs. horizontal layers)
- Branching strategy (trunk-based, gitflow, GitHub flow)
- Code review process
- Definition of done for each task

### 2. Phase Breakdown

Break the entire build into phases. For each phase:

#### Phase N: [Phase Name] (Week X - Week Y)
**Goal:** One-sentence goal for this phase
**Exit Criteria:** What must be true to move to the next phase

##### Tasks

| # | Task | Description | Depends On | Estimated Effort | Priority |
|---|------|------------|------------|-----------------|----------|
| N.1 | | | | S/M/L/XL | P0/P1/P2 |

Typical phases:
- **Phase 0**: Project setup (repo, CI/CD, dev environment, base config)
- **Phase 1**: Data layer (database schema, migrations, models, seed data)
- **Phase 2**: Auth (registration, login, token management, middleware)
- **Phase 3-N**: Feature implementation (one phase per major feature or feature group)
- **Phase N+1**: Integration & polish (cross-feature integration, error handling, loading states)
- **Phase N+2**: Testing & QA (test coverage, manual QA, bug fixes)
- **Phase N+3**: Deployment & launch (staging deploy, production setup, monitoring, launch)

### 3. Task Details

For each task in the phases above, provide:

#### Task N.X: [Task Name]
- **Files to create/modify**: List specific file paths based on the architecture's folder structure
- **Key implementation notes**: Technical guidance, patterns to follow, gotchas
- **Acceptance criteria**:
  - [ ] Specific, testable criteria
- **Testing requirements**: What tests to write (unit, integration, e2e)

### 4. Dependency Graph

Show which tasks block other tasks:
```
Phase 0 (Setup)
  └─→ Phase 1 (Data Layer)
       ├─→ Phase 2 (Auth)
       │    └─→ Phase 3 (Feature A)
       │         └─→ Phase 5 (Integration)
       └─→ Phase 4 (Feature B)
            └─→ Phase 5 (Integration)
```

### 5. Risk Register

| Risk | Probability | Impact | Mitigation | Contingency |
|------|------------|--------|------------|-------------|

### 6. Testing Strategy
- Unit testing approach (what to test, coverage target)
- Integration testing approach
- E2E testing approach (tools, critical paths to cover)
- Performance testing plan
- Manual QA checklist for each phase

### 7. Milestone Checkpoints

| Milestone | Target Date | Criteria | Demo-able? |
|-----------|------------|----------|------------|
| Dev environment working | | CI passes, local dev runs | Yes |
| Auth complete | | Can register, login, access protected routes | Yes |
| Core feature complete | | Primary user flow works end-to-end | Yes |
| MVP ready | | All P0 features, deployed to staging | Yes |
| Production launch | | Monitoring active, rollback plan tested | Yes |

### 8. Post-Launch Plan
- Monitoring checklist for first 48 hours
- Performance baseline to establish
- Feedback collection plan
- First iteration priorities

## Output Format
- Use markdown with clear hierarchy
- Use tables for task lists
- Use ASCII diagrams for dependency graphs
- Effort estimates: S (< 2hrs), M (2-4hrs), L (4-8hrs), XL (1-2 days)
- Target length: 2,000-4,000 words
