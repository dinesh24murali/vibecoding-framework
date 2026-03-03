# Prompt: Generate Developer Setup Guide

## Instructions

Copy everything below the `---` line into your LLM. Replace the placeholders with your previously generated documents.

---

You are a senior developer writing a Developer Setup Guide for a new team member who needs to get the project running locally from scratch.

Here is the project context:

**Intake Answers:**
```
{{INTAKE_ANSWERS}}
```

**Technical Architecture:**
```
{{TECH_ARCHITECTURE}}
```

**Implementation Plan:**
```
{{IMPLEMENTATION_PLAN}}
```

Generate a complete developer setup guide. Assume the reader is a competent developer but has never seen this project before. Every command should be copy-pasteable. Nothing should require guesswork.

## Required Sections

### 1. Prerequisites

List everything needed before starting, with exact version requirements and install commands:

| Tool | Minimum Version | Install Command (macOS) | Install Command (Linux) |
|------|----------------|------------------------|------------------------|
| Node.js | | | |
| Python | | | |
| Docker | | | |
| ... | | | |

Include verification commands:
```bash
node --version  # Expected: vX.Y.Z+
```

### 2. Repository Setup

```bash
# Step-by-step commands from clone to running
git clone <repo-url>
cd <project-name>
```

### 3. Environment Configuration

- List every environment variable needed
- Provide a `.env.example` with all variables, defaults, and comments
- Explain how to obtain any API keys or secrets needed for local dev
- Distinguish between required and optional variables

```env
# === Required ===
DATABASE_URL=postgresql://user:pass@localhost:5432/dbname  # Local PostgreSQL
JWT_SECRET=your-local-dev-secret-change-in-production

# === Optional ===
SENTRY_DSN=  # Leave empty for local dev
```

### 4. Database Setup

Step-by-step database setup:
```bash
# Create database
# Run migrations
# Seed sample data
```

Include how to:
- Reset the database
- Create new migrations
- View the database (recommended GUI tool)

### 5. Running the Application

#### Development Mode
```bash
# Start all services
# What URL to open
# Expected output / how to verify it's working
```

#### With Docker Compose (if applicable)
```bash
docker compose up -d
# What services are running and on what ports
```

#### Port Map
| Service | Port | URL |
|---------|------|-----|
| Frontend | 3000 | http://localhost:3000 |
| Backend API | 8000 | http://localhost:8000 |
| Database | 5432 | - |
| Redis | 6379 | - |

### 6. Running Tests

```bash
# Run all tests
# Run specific test file
# Run with coverage
# Run only unit tests
# Run only integration tests
```

### 7. Common Development Tasks

#### Adding a new API endpoint
1. Step-by-step process
2. Files to create/modify
3. How to test it

#### Adding a new database migration
```bash
# Commands with explanation
```

#### Adding a new frontend page/component
1. Step-by-step process

#### Updating dependencies
```bash
# Commands with explanation
```

### 8. Code Style & Linting

- Linter configuration and how to run it
- Formatter configuration and how to run it
- Pre-commit hooks setup
- Editor/IDE recommended extensions

```bash
# Format code
# Run linter
# Run type checker
```

### 9. Project Structure

```
project-root/
├── src/
│   ├── api/           # Explain what goes here
│   ├── models/        # Explain what goes here
│   ...
├── tests/
├── docker/
├── docs/
└── ...
```

### 10. Debugging

- How to attach a debugger (VS Code launch config, etc.)
- How to view logs
- How to inspect the database
- How to test API endpoints (curl examples, Postman collection, etc.)
- Common error messages and their solutions

### 11. Troubleshooting

| Problem | Cause | Solution |
|---------|-------|---------|
| Port already in use | Another process on port X | `lsof -i :X` then kill |
| Database connection refused | PostgreSQL not running | `brew services start postgresql` |
| ... | | |

### 12. Useful Commands Cheat Sheet

```bash
# Start everything
# Stop everything
# View logs
# Reset database
# Run tests
# Deploy to staging
```

## Output Format
- Use markdown
- Every command must be in a fenced code block with the correct language tag
- Test every command sequence mentally for correctness — no broken pipes
- Include expected output where it helps verify success
- Target length: 1,500-3,000 words
