Perfect — I’ll run this as a structured discovery and then produce a normalized Project Context Pack with Assumptions, Open Questions, Risks, and a Decision Log.
Please answer the questions below (short bullets are fine).
1) Vision, goals, success metrics
- What is the product in one sentence?
- What business problem are we solving, and why now?
- Top 3 business goals for the first 6–12 months?
- What metrics define success (e.g., activation %, conversion %, revenue, retention, SLA)?
2) Users and journeys
- Primary persona(s): who uses it most, and what do they need done?
- Secondary persona(s): who else matters?
- What are the top 3 critical user journeys (end-to-end steps + desired outcome)?
- Any high-risk or high-friction journey we should optimize first?
3) Scope boundaries
- What is explicitly in-scope for v1?
- What is explicitly out-of-scope for v1?
- What is “nice-to-have” if time permits?
4) Functional requirements and business rules
- Core features/modules required for launch?
- Key business rules or validations (pricing, approvals, eligibility, limits, workflows)?
- Any role-based permissions needed (admin, manager, end-user, support)?
- Are there required notifications (email/SMS/push/webhooks)?
5) Non-functional requirements
- Performance targets (p95 response time, page load, throughput)?
- Reliability targets (uptime %, RTO/RPO, backup expectations)?
- Security requirements (SSO/MFA, encryption, secrets handling, audit trails)?
- Compliance/privacy requirements (GDPR, HIPAA, SOC2, data residency, retention)?
- Accessibility target (WCAG level, supported assistive tech)?
6) Data model and integrations
- What are the core business entities (e.g., User, Account, Order, Subscription)?
- Any expected data volumes (records/day, file sizes, growth)?
- System-of-record(s): where does truth live for key entities?
- Required integrations (payments, CRM, ERP, identity, analytics, messaging)?
- Import/export needs (CSV, APIs, scheduled syncs, webhooks)?
7) API and platform interface
- Preferred API style: REST, GraphQL, or mixed?
- Auth model: session, JWT, OAuth2/OIDC, API keys, service-to-service auth?
- Versioning strategy (URI versioning, header-based, schema evolution)?
- Any public API or partner-facing developer platform requirements?
8) Technical constraints and delivery
- Preferred cloud/provider and managed services constraints?
- Required tech stack (frontend, backend, DB, queue, cache)?
- Team composition/size (engineering, product, design, QA, DevOps)?
- Release cadence target (weekly, biweekly, continuous)?
- Budget/hosting constraints and expected runway?
9) Environments and release strategy
- Required environments (local/dev/stage/prod + preview envs)?
- Deployment strategy (blue/green, canary, rolling)?
- Test strategy gates (unit/integration/e2e, performance, security scans)?
- Data strategy per environment (seed data, masked prod snapshots)?
10) Observability and operations
- Logging/metrics/tracing expectations (tooling preferences)?
- Alerting model (on-call hours, PagerDuty/Slack escalation, severity levels)?
- Operational dashboards needed for product + engineering?
- Incident response expectations (runbooks, postmortems, SLAs/SLOs)?