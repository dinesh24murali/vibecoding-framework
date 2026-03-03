# Prompt: Generate Production Technical Architecture

## Instructions

Copy everything below the `---` line into your LLM. Replace the placeholders with your intake answers, PRD, and functional spec.

---

You are a senior software architect designing a production-grade technical architecture.

Here is the project context:

**Intake Answers:**
```
{{INTAKE_ANSWERS}}
```

**Product Requirements Document (PRD):**
```
{{PRD}}
```

**Functional Specification:**
```
{{FUNCTIONAL_SPEC}}
```

Generate a Production Technical Architecture document with the following structure. Favor practical, proven technology choices over cutting-edge. Optimize for the team size and timeline specified in the intake.

## Required Sections

### 1. Architecture Overview
- High-level system description (1 paragraph)
- Architecture style (monolith, modular monolith, microservices, serverless, hybrid) with justification
- ASCII system diagram showing all major components and their connections:
```
[Client] --> [API Gateway] --> [App Server] --> [Database]
                                    |
                                    v
                               [Cache Layer]
```

### 2. Technology Stack

| Layer | Technology | Version | Justification |
|-------|-----------|---------|---------------|
| Frontend | | | |
| Backend | | | |
| Database | | | |
| Cache | | | |
| Queue | | | |
| Search | | | |

- Justify each choice against the constraints (team size, budget, timeline)
- Note any alternatives considered and why they were rejected

### 3. Application Architecture

#### Frontend Architecture
- Framework and rendering strategy (CSR, SSR, SSG, ISR)
- State management approach
- Folder structure
- Key libraries and their purposes
- Build and bundle strategy

#### Backend Architecture
- Framework and language
- Project structure / layering (controller → service → repository, etc.)
- Middleware pipeline
- Dependency injection approach (if applicable)
- Background job processing

#### Data Architecture
- Database schema overview (key entities and relationships as ERD or table)
- Indexing strategy for key queries
- Migration strategy
- Data seeding approach for development

### 4. API Design
- API style (REST, GraphQL, gRPC) with justification
- Versioning strategy
- Authentication mechanism (JWT, session, API key) and token lifecycle
- Rate limiting strategy
- Request/response conventions (envelope format, error format)
- Pagination strategy

### 5. Infrastructure Architecture

#### Environments
| Environment | Purpose | Infra | URL Pattern |
|------------|---------|-------|-------------|
| Local | Development | Docker Compose | localhost:3000 |
| Staging | QA/Testing | | |
| Production | Live | | |

#### Deployment
- CI/CD pipeline design
- Deployment strategy (blue/green, rolling, canary)
- Container strategy (Docker images, registries)
- Infrastructure as Code tool (Terraform, Pulumi, CDK, none)

#### Networking
- DNS and domain setup
- CDN strategy
- Load balancing approach
- SSL/TLS certificate management

### 6. Security Architecture
- Authentication flow (diagram)
- Authorization model (RBAC, ABAC, ACL)
- Secret management (env vars, vault, cloud secrets manager)
- Input validation strategy
- CORS policy
- Security headers
- Data encryption (at rest, in transit)
- Dependency vulnerability scanning

### 7. Observability
- Logging strategy (structured logging, log levels, aggregation)
- Monitoring and alerting (metrics to track, alerting thresholds)
- Distributed tracing (if applicable)
- Error tracking (Sentry, Bugsnag, etc.)
- Health check endpoints

### 8. Scalability & Performance
- Caching strategy (what to cache, TTLs, invalidation)
- Database connection pooling
- Horizontal scaling approach
- Performance budgets (page load time, API response time)
- CDN and static asset optimization

### 9. Disaster Recovery & Reliability
- Backup strategy (frequency, retention, testing)
- Recovery Point Objective (RPO) and Recovery Time Objective (RTO)
- Failover strategy
- Circuit breaker patterns (if applicable)
- Data retention and archival policy

### 10. Cost Estimation
| Service | Tier | Monthly Cost (Est.) |
|---------|------|-------------------|

- Total estimated monthly cost
- Cost scaling projections at 10x and 100x users

### 11. Technical Debt & Trade-offs
- Shortcuts taken for MVP and their future cost
- Known limitations
- Recommended post-MVP improvements

## Output Format
- Use markdown
- Use ASCII diagrams (no image references)
- Be specific about versions and configurations
- Target length: 2,500-5,000 words
