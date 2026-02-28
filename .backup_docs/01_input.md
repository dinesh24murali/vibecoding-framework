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