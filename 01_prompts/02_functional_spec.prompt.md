# Prompt: Generate Functional Specifications

## Instructions

Copy everything below the `---` line into your LLM. Replace `{{INTAKE_ANSWERS}}` with your intake and `{{PRD}}` with the previously generated PRD.

---

You are a senior product designer and business analyst writing a Functional Specification document.

Here is the project context:

**Intake Answers:**
```
{{INTAKE_ANSWERS}}
```

**Product Requirements Document (PRD):**
```
{{PRD}}
```

Generate a detailed Functional Specification with the following structure. This document should be precise enough that a developer can implement features without ambiguity.

## Required Sections

### 1. Overview
- Document purpose and audience
- Relationship to the PRD
- Conventions used (terminology, notation)

### 2. Feature Specifications

For EACH feature listed in the PRD's scope, create a detailed specification block:

#### Feature: [Feature Name]
- **Description**: What it does in 2-3 sentences
- **User roles involved**: Which user types interact with this feature
- **Preconditions**: What must be true before this feature can be used
- **Trigger**: What initiates this feature (user action, system event, scheduled)

##### User Flow
Step-by-step flow in numbered format:
1. User does X
2. System responds with Y
3. ...

##### Input/Output Specification
| Field | Type | Required | Validation Rules | Example |
|-------|------|----------|-----------------|---------|

##### Business Rules
- Numbered list of business rules that govern behavior
- Include edge cases

##### State Transitions
If the feature involves status changes (e.g., order: pending → confirmed → shipped):
```
[State A] --trigger--> [State B] --trigger--> [State C]
```

##### Error Handling
| Error Condition | System Behavior | User Message |
|----------------|-----------------|--------------|

##### Acceptance Criteria
- [ ] Given [context], when [action], then [expected result]
- [ ] ...

### 3. Cross-Feature Interactions
- How features depend on or interact with each other
- Shared state or data flows between features
- Conflict resolution (when two features could produce contradictory results)

### 4. Navigation & Information Architecture
- Page/screen inventory with descriptions
- Navigation hierarchy (site map or app structure)
- Entry points and landing behavior

### 5. Notification Specifications
For each notification type:
| Trigger | Channel | Recipient | Content Summary | Frequency Limit |
|---------|---------|-----------|----------------|-----------------|

### 6. Search & Filtering
- What is searchable
- Filter options per list/table view
- Sort options and defaults
- Pagination behavior

### 7. Permissions Matrix
| Action | Admin | User | Viewer | Guest |
|--------|-------|------|--------|-------|

### 8. Data Display Rules
- Date/time formats
- Number formats (currency, percentages)
- Empty states (what to show when there's no data)
- Loading states
- Truncation rules for long text

### 9. Glossary
- Domain-specific terms and their definitions

## Output Format
- Use markdown with clear heading hierarchy
- Use tables for structured data
- Use checklists for acceptance criteria
- Target length: 2,000-5,000 words depending on feature count
