# Sample Intake: TaskFlow (Task Management App)

This is a filled-out example to show how the intake questionnaire works.

---

## REQUIRED: Core Information

### 1. App Identity
- **App Name**: TaskFlow
- **One-line description**: A collaborative task management app for small teams that tracks work from idea to done.
- **Problem it solves**: Small teams (3-10 people) use scattered tools — Slack threads, spreadsheets, sticky notes — to track work. TaskFlow gives them one simple board to plan, assign, and track tasks without the complexity of Jira.

### 2. Target Users
- **Primary user type**: Small team leads and project managers at startups (10-50 employees)
- **Secondary users**: Individual contributors on those teams (developers, designers, marketers)
- **Rough user scale**: 1,000 users within 6 months, 10,000 within 18 months

### 3. Core Features
1. **Kanban board** — Drag-and-drop cards across customizable columns (e.g., To Do, In Progress, Done)
2. **Task creation & assignment** — Create tasks with title, description, assignee, due date, labels, and priority
3. **Team workspaces** — Create a workspace, invite members via email, manage roles (admin/member)
4. **Activity feed** — See a chronological log of all changes on a board (who moved what, who commented)
5. **Comments & mentions** — Comment on tasks, @mention teammates to notify them

### 4. Tech Preferences
- **Platform**: Web app
- **Frontend preference**: Next.js
- **Backend preference**: Python/FastAPI
- **Database preference**: PostgreSQL
- **Hosting preference**: AWS (but keep costs low, start with minimal infra)

### 5. Constraints
- **Budget range**: $50/mo for infrastructure
- **Timeline**: MVP in 6 weeks
- **Team size**: Solo developer

---

## OPTIONAL DEEP-DIVE: Business Context

### 6. Business Model
- **Revenue model**: Freemium — free for up to 3 users per workspace, $8/user/month for more
- **Pricing tiers**: Free (3 users, 1 workspace), Pro ($8/user/mo, unlimited workspaces, file attachments)
- **Key competitors**: Trello, Linear, Asana, Notion
- **Differentiator**: Simpler than Linear/Asana, more focused than Notion, faster setup than any of them. Zero learning curve.

### 7. Success Metrics
- **How do you measure success?**: Weekly active workspaces, tasks completed per week
- **MVP success criteria**: 5 teams actively using it weekly, average >10 tasks completed per team per week

---

## OPTIONAL DEEP-DIVE: Users & Flows

### 8. User Personas

**Persona 1:**
- Name/Role: Sarah, Startup CTO (team of 6)
- Goals: See what everyone is working on at a glance, unblock people quickly
- Frustrations: Jira is way too complex, Trello doesn't show enough context, Slack threads get lost

**Persona 2:**
- Name/Role: Marcus, Frontend Developer
- Goals: Know exactly what to work on next, update status without friction
- Frustrations: Forgets to update task status because it's too many clicks, hates context-switching to a project tool

### 9. Critical User Flows

**Flow 1: Team lead sets up a workspace and invites the team**
1. Sign up with email
2. Create a workspace ("Acme Engineering")
3. Create a board ("Sprint 12")
4. Invite 3 team members via email
5. Teammates click invite link, sign up, land in the workspace

**Flow 2: Developer picks up and completes a task**
1. Open the board
2. See tasks in "To Do" column assigned to them
3. Drag a task to "In Progress"
4. Add a comment with a question
5. When done, drag to "Done"

---

## OPTIONAL DEEP-DIVE: Technical Requirements

### 10. Integrations
- **Third-party services**: None for MVP. Post-MVP: Slack notifications, GitHub PR linking.
- **External APIs consumed**: None for MVP
- **APIs exposed**: REST API — will eventually build a public API for integrations

### 11. Authentication & Authorization
- **Auth method**: Email/password for MVP, Google OAuth post-MVP
- **Role types**: Workspace Admin, Member
- **Multi-tenancy**: Multi-tenant with data isolation per workspace

### 12. Data & Compliance
- **Sensitive data handled**: Email addresses and names (PII)
- **Compliance requirements**: GDPR basics (data export, account deletion)
- **Data residency**: No specific requirements
- **Backup/recovery**: Daily database backups, 30-day retention

---

## Anything Else?

Real-time updates on the Kanban board would be amazing for MVP but can be post-MVP if it adds too much complexity. WebSocket support for live card movement across clients.
