# Vibe Coding Document Generator

A prompt framework that generates the complete documentation set needed to build an application from scratch.

## How It Works

```
You fill out the Intake    →    Feed into Generator Prompts    →    Get production-ready docs
(00_intake/)                    (01_prompts/)                       (02_outputs/)
```

### Step-by-step

1. **Fill out the intake** — Copy `00_intake/intake_questionnaire.md`, answer the required questions (and optionally the deep-dive sections)
2. **Generate docs in order** — Each prompt in `01_prompts/` is designed to be pasted into an LLM along with your intake answers. Follow the numbered order:
   - `01_prd.prompt.md` → Production PRD
   - `02_functional_spec.prompt.md` → Functional Specifications
   - `03_tech_architecture.prompt.md` → Production Technical Architecture
   - `04_api_spec.prompt.md` → OpenAPI Specification
   - `05_implementation_plan.prompt.md` → Implementation Plan
   - `06_dev_setup.prompt.md` → Developer Setup Guide
3. **Each prompt builds on the previous** — Later prompts reference earlier generated documents, so generate them in order and include prior outputs as context.

### Document Dependency Chain

```
Intake Answers
  └─→ PRD (what & why)
       └─→ Functional Spec (detailed behavior)
            └─→ Tech Architecture (how it's built)
                 └─→ OpenAPI Spec (API contracts)
                      └─→ Implementation Plan (execution order)
                           └─→ Dev Setup Guide (how to run it)
```

## Directory Structure

```
vibe_template/
├── README.md                          # You are here
├── 00_intake/
│   └── intake_questionnaire.md        # Fill this out first
├── 01_prompts/
│   ├── 01_prd.prompt.md               # → Production PRD
│   ├── 02_functional_spec.prompt.md   # → Functional Specifications
│   ├── 03_tech_architecture.prompt.md # → Technical Architecture
│   ├── 04_api_spec.prompt.md          # → OpenAPI Specification
│   ├── 05_implementation_plan.prompt.md # → Implementation Plan
│   └── 06_dev_setup.prompt.md         # → Developer Setup Guide
├── 02_outputs/                        # Save generated docs here
└── examples/
    └── sample_intake_filled.md        # Worked example (task management app)
```

## Tips

- **Minimum viable intake**: Only the "Required" sections need answers. The more you fill out, the better the output.
- **Iterate**: If a generated document doesn't feel right, refine your intake answers and re-run that prompt.
- **Context stacking**: When running prompt N, include the outputs from prompts 1 through N-1 as context for best results.
- **Any LLM works**: These prompts are model-agnostic. Use Claude, GPT, Gemini, or any capable LLM.
