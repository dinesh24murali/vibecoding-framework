# Prompt: Generate Production PRD

## Instructions

Copy everything below the `---` line into your LLM. Replace `{{INTAKE_ANSWERS}}` with your filled-out intake questionnaire.

---

You are a senior product manager writing a Production Product Requirements Document (PRD).

Below are the intake answers describing an application to be built:

```
{{INTAKE_ANSWERS}}
```

Generate a comprehensive PRD with the following structure. Be specific and actionable — avoid vague language. Every requirement should be testable or measurable.

## Required PRD Sections

### 1. Executive Summary
- Product name and one-paragraph overview
- Problem statement (backed by the intake answers)
- Proposed solution summary
- Target launch timeline

### 2. Vision & Strategic Alignment
- Product vision (where this is heading in 1-2 years)
- Strategic goals this product serves
- Key assumptions and hypotheses to validate

### 3. Target Audience
- Primary and secondary user segments
- User personas (expand from intake if provided, create if not)
- Market size estimation (TAM/SAM/SOM if applicable)

### 4. Goals & Success Metrics
- Primary KPIs with specific targets (e.g., "500 DAU within 3 months of launch")
- Secondary metrics
- North star metric

### 5. Scope Definition
#### In Scope (MVP)
- Numbered list of features with brief descriptions
- Priority ranking (P0 = must-have, P1 = should-have, P2 = nice-to-have)

#### Out of Scope (Post-MVP)
- Features explicitly deferred and why

### 6. User Stories
For each core feature, write user stories in the format:
> As a [user type], I want to [action] so that [benefit].

Include acceptance criteria for each story.

### 7. Constraints & Assumptions
- Technical constraints (from intake tech preferences)
- Budget constraints
- Timeline constraints
- Team constraints
- Key assumptions that if wrong would change the plan

### 8. Risks & Mitigations
| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|

### 9. Release Strategy
- MVP definition (what ships first)
- Phase 2 additions
- Phase 3 vision
- Go-to-market approach (if applicable)

### 10. Open Questions
- List any ambiguities or decisions that still need to be made

## Output Format
- Use markdown
- Be specific — no "TBD" or "to be determined" without flagging it as an open question
- Target length: 1,500-3,000 words
