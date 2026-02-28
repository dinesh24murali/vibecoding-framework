You are a Senior Developer Experience (DevEx) Engineer.
Generate a PRODUCTION-READY Developer Setup Guide for this project.
This is not a quickstart draft; it must be complete and operational.

Inputs:
1) Project Context Pack
2) Approved Technical Architecture
3) Approved OpenAPI
4) Approved Implementation Plan

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

Approved Technical Architecture:
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

Approved OpenAPI:
```
openapi: 3.1.0
info:
  title: Village DTH Platform API
  version: 1.0.0
  summary: REST API for customer purchase flow and admin operations
  description: |
    Production API contract for DTH plan purchase and service-request operations.
    Versioned via URI namespace `/api/v1`.
jsonSchemaDialect: https://json-schema.org/draft/2020-12/schema
servers:
  - url: /api/v1
    description: Versioned API base path
tags:
  - name: Auth
  - name: Customer
  - name: Payments
  - name: Admin Providers
  - name: Admin Plans
  - name: Admin Service Requests
security: []
paths:
  /auth/admin/login:
    post:
      tags: [Auth]
      operationId: adminLogin
      summary: Admin login and token issuance
      description: Authenticate admin and return access/refresh tokens.
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/AdminLoginRequest"
            examples:
              default:
                value:
                  identifier: "admin"
                  password: "StrongPassword123!"
      responses:
        "200":
          description: Login success
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/AdminLoginResponse"
              examples:
                success:
                  value:
                    access_token: "eyJhbGciOiJIUzI1NiIsInR5cCI..."
                    token_type: "Bearer"
                    expires_in: 864000
                    refresh_token: "rt_abc123"
                    refresh_expires_in: 7776000
                    user:
                      id: "usr_1"
                      name: "Village Admin"
                      phone_number: "9000000001"
                      role: "admin"
                      status: "active"
                      created_at: "2026-02-27T10:00:00Z"
                      updated_at: "2026-02-27T10:00:00Z"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "422":
          $ref: "#/components/responses/ValidationError"
        "500":
          $ref: "#/components/responses/InternalError"

  /customer/providers:
    get:
      tags: [Customer]
      operationId: listCustomerProviders
      summary: List providers for customer app
      description: Returns non-deleted providers.
      security: []
      responses:
        "200":
          description: Provider list
          content:
            application/json:
              schema:
                type: object
                required: [items]
                properties:
                  items:
                    type: array
                    items:
                      $ref: "#/components/schemas/Provider"
              examples:
                default:
                  value:
                    items:
                      - id: "pro_1"
                        name: "Airtel DTH"
                        image_url: "https://cdn.example.com/providers/airtel.png"
                        created_at: "2026-02-27T10:00:00Z"
                        updated_at: "2026-02-27T10:00:00Z"
        "500":
          $ref: "#/components/responses/InternalError"

  /customer/providers/{providerId}/plans:
    get:
      tags: [Customer]
      operationId: listCustomerPlansByProvider
      summary: List active plans by provider for customer app
      security: []
      parameters:
        - $ref: "#/components/parameters/ProviderIdPath"
      responses:
        "200":
          description: Plans by provider
          content:
            application/json:
              schema:
                type: object
                required: [items]
                properties:
                  items:
                    type: array
                    items:
                      $ref: "#/components/schemas/Plan"
              examples:
                default:
                  value:
                    items:
                      - id: "pln_1"
                        provider_id: "pro_1"
                        name: "Monthly Saver"
                        description: "100+ channels, 30 days validity"
                        price: 299
                        discount: 20
                        is_active: true
                        created_at: "2026-02-27T10:00:00Z"
                        updated_at: "2026-02-27T10:00:00Z"
        "404":
          $ref: "#/components/responses/NotFoundError"
        "500":
          $ref: "#/components/responses/InternalError"

  /customer/checkout/service-requests:
    post:
      tags: [Customer, Payments]
      operationId: createServiceRequestAndPayment
      summary: Create service request and initiate payment
      description: Creates/uses customer record, creates service request, creates Razorpay order.
      security: []
      parameters:
        - $ref: "#/components/parameters/IdempotencyKeyHeader"
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/PaymentRequest"
            examples:
              default:
                value:
                  provider_id: "pro_1"
                  plan_id: "pln_1"
                  customer:
                    name: "Kumar"
                    phone_number: "9000000002"
                  recaptcha_token: "03AFcWeA..."
      responses:
        "201":
          description: Service request and payment order created
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/PaymentResponse"
              examples:
                created:
                  value:
                    service_request:
                      id: "sr_1"
                      user_id: "usr_2"
                      plan_id: "pln_1"
                      status: "pending"
                      created_at: "2026-02-27T10:10:00Z"
                      updated_at: "2026-02-27T10:10:00Z"
                    payment:
                      gateway: "razorpay"
                      order_id: "order_Qx1Ab2C3"
                      amount: 279
                      currency: "INR"
                      status: "created"
        "400":
          $ref: "#/components/responses/BadRequestError"
        "409":
          $ref: "#/components/responses/ConflictError"
        "422":
          $ref: "#/components/responses/ValidationError"
        "429":
          $ref: "#/components/responses/TooManyRequestsError"
        "500":
          $ref: "#/components/responses/InternalError"

  /customer/service-requests/{serviceRequestId}/retry-payment:
    post:
      tags: [Customer, Payments]
      operationId: retryPaymentForServiceRequest
      summary: Retry payment using the same service request
      description: Reuses existing service request and creates a new payment attempt.
      security: []
      parameters:
        - $ref: "#/components/parameters/ServiceRequestIdPath"
        - $ref: "#/components/parameters/IdempotencyKeyHeader"
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [recaptcha_token]
              properties:
                recaptcha_token:
                  type: string
                  minLength: 10
      responses:
        "200":
          description: Retry payment order created
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/PaymentResponse"
        "404":
          $ref: "#/components/responses/NotFoundError"
        "409":
          $ref: "#/components/responses/ConflictError"
        "422":
          $ref: "#/components/responses/ValidationError"
        "500":
          $ref: "#/components/responses/InternalError"

  /customer/service-requests/{serviceRequestId}/payment-status:
    get:
      tags: [Customer, Payments]
      operationId: getCustomerPaymentStatus
      summary: Poll payment status for service request
      security: []
      parameters:
        - $ref: "#/components/parameters/ServiceRequestIdPath"
        - name: phone_number
          in: query
          required: true
          schema:
            type: string
            pattern: "^[6-9][0-9]{9}$"
      responses:
        "200":
          description: Payment/service-request status
          content:
            application/json:
              schema:
                type: object
                required: [service_request]
                properties:
                  service_request:
                    $ref: "#/components/schemas/ServiceRequest"
        "404":
          $ref: "#/components/responses/NotFoundError"
        "422":
          $ref: "#/components/responses/ValidationError"
        "500":
          $ref: "#/components/responses/InternalError"

  /payments/razorpay/callback:
    post:
      tags: [Payments]
      operationId: handleRazorpayCallback
      summary: Razorpay callback/webhook endpoint
      description: Verifies gateway signature and updates payment/service request status.
      security: []
      parameters:
        - name: X-Razorpay-Signature
          in: header
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/RazorpayCallbackRequest"
      responses:
        "200":
          description: Callback processed
          content:
            application/json:
              schema:
                type: object
                required: [acknowledged]
                properties:
                  acknowledged:
                    type: boolean
                    const: true
        "400":
          $ref: "#/components/responses/BadRequestError"
        "401":
          description: Invalid gateway signature
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ErrorResponse"
        "409":
          description: Duplicate callback safely ignored
          content:
            application/json:
              schema:
                type: object
                required: [acknowledged, duplicate]
                properties:
                  acknowledged:
                    type: boolean
                  duplicate:
                    type: boolean
        "500":
          $ref: "#/components/responses/InternalError"

  /admin/providers:
    get:
      tags: [Admin Providers]
      operationId: listAdminProviders
      summary: List providers (admin)
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/PageQuery"
        - $ref: "#/components/parameters/PageSizeQuery"
        - name: search
          in: query
          schema:
            type: string
            maxLength: 100
        - name: sort
          in: query
          schema:
            type: string
            enum: [updated_at_desc, updated_at_asc, name_asc, name_desc]
            default: updated_at_desc
      responses:
        "200":
          description: Paginated providers
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/PaginatedProvidersResponse"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "500":
          $ref: "#/components/responses/InternalError"
    post:
      tags: [Admin Providers]
      operationId: createProvider
      summary: Create provider
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/ProviderCreateRequest"
      responses:
        "201":
          description: Provider created
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Provider"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "409":
          $ref: "#/components/responses/ConflictError"
        "422":
          $ref: "#/components/responses/ValidationError"
        "500":
          $ref: "#/components/responses/InternalError"

  /admin/providers/{providerId}:
    get:
      tags: [Admin Providers]
      operationId: getProviderById
      summary: Get provider by id
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/ProviderIdPath"
      responses:
        "200":
          description: Provider
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Provider"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"
    patch:
      tags: [Admin Providers]
      operationId: updateProvider
      summary: Update provider
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/ProviderIdPath"
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/ProviderUpdateRequest"
      responses:
        "200":
          description: Provider updated
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Provider"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"
        "409":
          $ref: "#/components/responses/ConflictError"
        "422":
          $ref: "#/components/responses/ValidationError"
    delete:
      tags: [Admin Providers]
      operationId: deleteProvider
      summary: Soft-delete provider
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/ProviderIdPath"
      responses:
        "204":
          description: Provider soft-deleted
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"

  /admin/providers/csv-upload:
    post:
      tags: [Admin Providers]
      operationId: uploadProvidersCsv
      summary: Upsert providers from CSV
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          multipart/form-data:
            schema:
              type: object
              required: [file]
              properties:
                file:
                  type: string
                  format: binary
      responses:
        "200":
          description: CSV processed
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/CsvUploadResult"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "422":
          $ref: "#/components/responses/ValidationError"

  /admin/plans:
    get:
      tags: [Admin Plans]
      operationId: listAdminPlans
      summary: List plans (admin)
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/PageQuery"
        - $ref: "#/components/parameters/PageSizeQuery"
        - name: provider_id
          in: query
          schema:
            type: string
        - name: is_active
          in: query
          schema:
            type: boolean
        - name: sort
          in: query
          schema:
            type: string
            enum: [updated_at_desc, updated_at_asc, price_asc, price_desc]
            default: updated_at_desc
      responses:
        "200":
          description: Paginated plans
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/PaginatedPlansResponse"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "500":
          $ref: "#/components/responses/InternalError"
    post:
      tags: [Admin Plans]
      operationId: createPlan
      summary: Create plan
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/PlanCreateRequest"
      responses:
        "201":
          description: Plan created
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Plan"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"
        "409":
          $ref: "#/components/responses/ConflictError"
        "422":
          $ref: "#/components/responses/ValidationError"

  /admin/plans/{planId}:
    get:
      tags: [Admin Plans]
      operationId: getPlanById
      summary: Get plan by id
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/PlanIdPath"
      responses:
        "200":
          description: Plan
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Plan"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"
    patch:
      tags: [Admin Plans]
      operationId: updatePlan
      summary: Update plan
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/PlanIdPath"
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/PlanUpdateRequest"
      responses:
        "200":
          description: Plan updated
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Plan"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"
        "409":
          $ref: "#/components/responses/ConflictError"
        "422":
          $ref: "#/components/responses/ValidationError"
    delete:
      tags: [Admin Plans]
      operationId: deletePlan
      summary: Soft-delete plan
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/PlanIdPath"
      responses:
        "204":
          description: Plan soft-deleted
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"

  /admin/plans/csv-upload:
    post:
      tags: [Admin Plans]
      operationId: uploadPlansCsv
      summary: Upsert plans from CSV
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          multipart/form-data:
            schema:
              type: object
              required: [file]
              properties:
                file:
                  type: string
                  format: binary
      responses:
        "200":
          description: CSV processed
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/CsvUploadResult"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "422":
          $ref: "#/components/responses/ValidationError"

  /admin/service-requests:
    get:
      tags: [Admin Service Requests]
      operationId: listServiceRequests
      summary: List service requests
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/PageQuery"
        - $ref: "#/components/parameters/PageSizeQuery"
        - name: status
          in: query
          schema:
            $ref: "#/components/schemas/ServiceRequestStatus"
        - name: provider_id
          in: query
          schema:
            type: string
        - name: from_date
          in: query
          schema:
            type: string
            format: date
        - name: to_date
          in: query
          schema:
            type: string
            format: date
        - name: sort
          in: query
          schema:
            type: string
            enum: [created_at_desc, created_at_asc, updated_at_desc, updated_at_asc]
            default: created_at_desc
      responses:
        "200":
          description: Paginated service requests
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/PaginatedServiceRequestsResponse"
        "401":
          $ref: "#/components/responses/UnauthorizedError"

  /admin/service-requests/{serviceRequestId}:
    get:
      tags: [Admin Service Requests]
      operationId: getServiceRequestById
      summary: Get service request details
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/ServiceRequestIdPath"
      responses:
        "200":
          description: Service request
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ServiceRequestDetails"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"

  /admin/service-requests/{serviceRequestId}/status:
    patch:
      tags: [Admin Service Requests]
      operationId: updateServiceRequestStatus
      summary: Update service request status
      description: Enforces allowed transitions; triggers SMS when transitioning to completed.
      security:
        - BearerAuth: []
      parameters:
        - $ref: "#/components/parameters/ServiceRequestIdPath"
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/ServiceRequestStatusUpdateRequest"
            examples:
              toCompleted:
                value:
                  status: "completed"
      responses:
        "200":
          description: Status updated
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ServiceRequest"
        "401":
          $ref: "#/components/responses/UnauthorizedError"
        "404":
          $ref: "#/components/responses/NotFoundError"
        "422":
          description: Invalid state transition
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ErrorResponse"

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  parameters:
    ProviderIdPath:
      name: providerId
      in: path
      required: true
      schema:
        type: string
        minLength: 1
    PlanIdPath:
      name: planId
      in: path
      required: true
      schema:
        type: string
        minLength: 1
    ServiceRequestIdPath:
      name: serviceRequestId
      in: path
      required: true
      schema:
        type: string
        minLength: 1
    PageQuery:
      name: page
      in: query
      schema:
        type: integer
        minimum: 1
        default: 1
    PageSizeQuery:
      name: page_size
      in: query
      schema:
        type: integer
        minimum: 1
        maximum: 100
        default: 20
    IdempotencyKeyHeader:
      name: Idempotency-Key
      in: header
      required: true
      schema:
        type: string
        minLength: 8
        maxLength: 128
      description: Unique key to prevent duplicate payment operations.

  responses:
    BadRequestError:
      description: Bad request
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/ErrorResponse"
    UnauthorizedError:
      description: Unauthorized
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/ErrorResponse"
    NotFoundError:
      description: Resource not found
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/ErrorResponse"
    ConflictError:
      description: Conflict / duplicate / idempotency conflict
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/ErrorResponse"
    ValidationError:
      description: Validation error
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/ErrorResponse"
    TooManyRequestsError:
      description: Too many requests
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/ErrorResponse"
    InternalError:
      description: Internal server error
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/ErrorResponse"

  schemas:
    User:
      type: object
      required: [id, name, phone_number, role, status, created_at, updated_at]
      properties:
        id:
          type: string
          examples: ["usr_1"]
        name:
          type: string
          minLength: 1
          maxLength: 120
        phone_number:
          type: string
          pattern: "^[6-9][0-9]{9}$"
        role:
          type: string
          enum: [customer, admin]
        status:
          type: string
          enum: [active, inactive, blocked]
        password:
          type: string
          minLength: 8
          writeOnly: true
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    Provider:
      type: object
      required: [id, name, created_at, updated_at]
      properties:
        id:
          type: string
        name:
          type: string
          minLength: 1
          maxLength: 120
        image_url:
          type: string
          format: uri
          nullable: true
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    Plan:
      type: object
      required: [id, provider_id, name, description, price, discount, is_active, created_at, updated_at]
      properties:
        id:
          type: string
        provider_id:
          type: string
        name:
          type: string
          minLength: 1
          maxLength: 120
        description:
          type: string
          minLength: 1
          maxLength: 1000
        price:
          type: number
          minimum: 0.01
          multipleOf: 0.01
        discount:
          type: number
          minimum: 0
          multipleOf: 0.01
        is_active:
          type: boolean
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    ServiceRequestStatus:
      type: string
      enum: [pending, payment_failed, payment_success, blocked, completed]
      description: Canonical API values for service request state.

    ServiceRequest:
      type: object
      required: [id, user_id, plan_id, status, created_at, updated_at]
      properties:
        id:
          type: string
        user_id:
          type: string
        plan_id:
          type: string
        status:
          $ref: "#/components/schemas/ServiceRequestStatus"
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    ServiceRequestDetails:
      allOf:
        - $ref: "#/components/schemas/ServiceRequest"
        - type: object
          properties:
            user:
              $ref: "#/components/schemas/User"
            plan:
              $ref: "#/components/schemas/Plan"
            provider:
              $ref: "#/components/schemas/Provider"

    PaymentRequest:
      type: object
      required: [provider_id, plan_id, customer, recaptcha_token]
      properties:
        provider_id:
          type: string
        plan_id:
          type: string
        customer:
          type: object
          required: [name, phone_number]
          properties:
            name:
              type: string
              minLength: 1
              maxLength: 120
            phone_number:
              type: string
              pattern: "^[6-9][0-9]{9}$"
        recaptcha_token:
          type: string
          minLength: 10

    PaymentResponse:
      type: object
      required: [service_request, payment]
      properties:
        service_request:
          $ref: "#/components/schemas/ServiceRequest"
        payment:
          type: object
          required: [gateway, order_id, amount, currency, status]
          properties:
            gateway:
              type: string
              enum: [razorpay]
            order_id:
              type: string
            amount:
              type: number
              minimum: 0.01
              multipleOf: 0.01
            currency:
              type: string
              enum: [INR]
            status:
              type: string
              enum: [created, authorized, captured, failed]

    RazorpayCallbackRequest:
      type: object
      required: [event, payload]
      properties:
        event:
          type: string
          examples: ["payment.captured"]
        payload:
          type: object
          additionalProperties: true

    ProviderCreateRequest:
      type: object
      required: [name]
      properties:
        name:
          type: string
          minLength: 1
          maxLength: 120
        image_url:
          type: string
          format: uri
          nullable: true

    ProviderUpdateRequest:
      type: object
      minProperties: 1
      properties:
        name:
          type: string
          minLength: 1
          maxLength: 120
        image_url:
          type: string
          format: uri
          nullable: true

    PlanCreateRequest:
      type: object
      required: [provider_id, name, description, price, discount, is_active]
      properties:
        provider_id:
          type: string
        name:
          type: string
          minLength: 1
          maxLength: 120
        description:
          type: string
          minLength: 1
          maxLength: 1000
        price:
          type: number
          minimum: 0.01
          multipleOf: 0.01
        discount:
          type: number
          minimum: 0
          multipleOf: 0.01
        is_active:
          type: boolean

    PlanUpdateRequest:
      type: object
      minProperties: 1
      properties:
        provider_id:
          type: string
        name:
          type: string
          minLength: 1
          maxLength: 120
        description:
          type: string
          minLength: 1
          maxLength: 1000
        price:
          type: number
          minimum: 0.01
          multipleOf: 0.01
        discount:
          type: number
          minimum: 0
          multipleOf: 0.01
        is_active:
          type: boolean

    ServiceRequestStatusUpdateRequest:
      type: object
      required: [status]
      properties:
        status:
          $ref: "#/components/schemas/ServiceRequestStatus"

    CsvUploadResult:
      type: object
      required: [total_rows, success_rows, failed_rows, errors]
      properties:
        total_rows:
          type: integer
          minimum: 0
        success_rows:
          type: integer
          minimum: 0
        failed_rows:
          type: integer
          minimum: 0
        errors:
          type: array
          items:
            $ref: "#/components/schemas/CsvRowError"
      examples:
        partial:
          value:
            total_rows: 3
            success_rows: 2
            failed_rows: 1
            errors:
              - row: 3
                field: "name"
                issue: "duplicate_plan_name_for_provider"

    CsvRowError:
      type: object
      required: [row, issue]
      properties:
        row:
          type: integer
          minimum: 1
        field:
          type: string
        issue:
          type: string

    PaginatedProvidersResponse:
      type: object
      required: [items, page, page_size, total]
      properties:
        items:
          type: array
          items:
            $ref: "#/components/schemas/Provider"
        page:
          type: integer
        page_size:
          type: integer
        total:
          type: integer

    PaginatedPlansResponse:
      type: object
      required: [items, page, page_size, total]
      properties:
        items:
          type: array
          items:
            $ref: "#/components/schemas/Plan"
        page:
          type: integer
        page_size:
          type: integer
        total:
          type: integer

    PaginatedServiceRequestsResponse:
      type: object
      required: [items, page, page_size, total]
      properties:
        items:
          type: array
          items:
            $ref: "#/components/schemas/ServiceRequest"
        page:
          type: integer
        page_size:
          type: integer
        total:
          type: integer

    AdminLoginRequest:
      type: object
      required: [identifier, password]
      properties:
        identifier:
          type: string
          minLength: 1
          maxLength: 120
          description: Username or registered phone number
        password:
          type: string
          minLength: 8

    AdminLoginResponse:
      type: object
      required: [access_token, token_type, expires_in, refresh_token, refresh_expires_in, user]
      properties:
        access_token:
          type: string
        token_type:
          type: string
          enum: [Bearer]
        expires_in:
          type: integer
          description: Access token TTL in seconds
          examples: [864000]
        refresh_token:
          type: string
        refresh_expires_in:
          type: integer
          description: Refresh token TTL in seconds
          examples: [7776000]
        user:
          $ref: "#/components/schemas/User"

    ErrorDetail:
      type: object
      required: [field, issue]
      properties:
        field:
          type: string
        issue:
          type: string

    Error:
      type: object
      required: [code, message, request_id]
      properties:
        code:
          type: string
        message:
          type: string
        details:
          type: array
          items:
            $ref: "#/components/schemas/ErrorDetail"
        request_id:
          type: string

    ErrorResponse:
      type: object
      required: [error]
      properties:
        error:
          $ref: "#/components/schemas/Error"
```

### Open Questions
- Should `/customer/service-requests/{id}/payment-status` require additional tokenized proof beyond phone number?
  Answer: No phone number should be fine
- Is token refresh endpoint required in v1 (currently omitted; login re-auth is fallback)?
  Answer: No required
- Do admin CSV uploads need asynchronous processing for larger files later?
  Answer: yes
- Should `users.phone_number` uniqueness include soft-deleted users or active-only scope?
  Answer: `users.phone_number` uniqueness should include soft-deleted users

Approved Implementation Plan:
```
## Implementation Plan - Village DTH Platform

### Metadata
- **Version:** 1.0
- **Date:** 2026-02-27
- **Owner:** Solo Developer (Staff Engineer role)
- **Input Artifacts:** PRD v1.0, FSD v1.0, TAD v1.0, OpenAPI v1.0 (`/api/v1`)
- **Execution Mode:** Single-prod, low-cost, production-ready

## Delivery assumptions and sequencing constraints
- Build order is **contract-first**: OpenAPI -> DB/migrations -> domain logic -> API handlers -> frontends -> ops hardening.
- Payment, status transitions, and async SMS are the highest-risk flows and must be completed before admin/customer UI finalization.
- CSV uploads are implemented in two phases: synchronous v1 endpoint + async-ready architecture hooks (for future larger files).
- JWT auth (10-day access, 3-month refresh) is implemented now, but refresh endpoint remains out-of-scope per decision.
- All admin endpoints require `BearerAuth`; customer endpoints are public with strict validation and reCAPTCHA checks.
- `users.phone_number` uniqueness includes soft-deleted users (global uniqueness).
- App-layer enforcement for max 100 plans per provider is mandatory; DB trigger is not included.

## Milestones (M1, M2, ...)
- **M1 Foundation & Contract Lock:** project skeleton, config, OpenAPI sync, migrations scaffold.
- **M2 Auth + Catalog Core:** admin auth, providers/plans CRUD, customer browse endpoints.
- **M3 Checkout & Payments:** service request creation, reCAPTCHA verify, Razorpay order/callback, retry audit.
- **M4 Admin Ops & Async:** service request transitions, SMS on `completed`, CSV upload + async preparation.
- **M5 Frontend Completion:** customer and admin apps wired to APIs.
- **M6 Production Readiness:** tests, backups, systemd, deployment docs, release cutover.

## Task Table

| task_id | title | status | priority | owner_role | design_pattern | api_mapping (operationId(s) + schema(s) + security requirement mapping) | files_to_create | files_to_modify | dependencies | acceptance_criteria | test_plan | risk_notes |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| T001 | Backend project bootstrap and module boundaries | todo | P0 | solo-dev | Layered + Dependency Injection | Schemas: `ErrorResponse`, `User`; Security: `BearerAuth` | `backend/cmd/api/main.go`, `backend/cmd/worker/main.go`, `backend/internal/{config,router,middleware,errors}/` | `backend/go.mod` | - | API and worker start with env config and health wiring stubs | unit | Poor structure early causes rework |
| T002 | OpenAPI contract import and server contract guards | todo | P0 | solo-dev | Contract-First | All operationIds; Schemas: all; Security: `BearerAuth` | `backend/api/openapi/openapi.yaml`, `backend/internal/http/contract/validator.go` | `backend/internal/router/router.go` | T001 | Runtime request/response validation enabled for critical endpoints in non-prod mode | contract | Contract drift if skipped |
| T003 | Base migrations + core schema setup | todo | P0 | solo-dev | Migration Runner | Schemas: `User`, `Provider`, `Plan`, `ServiceRequest`, `PaymentRequest`, `PaymentResponse`, `CsvUploadResult` | `backend/db/migrations/0001_init.sql`, `backend/db/migrations/0002_indexes.sql`, `backend/db/migrations/0003_payment_sms_retry.sql` | `backend/internal/db/db.go` | T001 | Tables/indexes exist, constraints match PRD/FSD/TAD | integration | Bad schema locks later behavior |
| T004 | Migration runner + release execution command | todo | P0 | solo-dev | Command Pattern | Schemas: all persisted entities | `backend/cmd/migrate/main.go`, `backend/internal/migrate/runner.go` | `backend/Makefile` | T003 | `migrate up` runs ordered scripts with failure stop | integration | Release failures without safe migration |
| T005 | Seed admin command (idempotent) | todo | P0 | solo-dev | Repository | Operation: `adminLogin`; Schemas: `AdminLoginRequest`,`AdminLoginResponse`,`User`; Security: `BearerAuth` | `backend/cmd/seed-admin/main.go`, `backend/internal/seed/admin_seed.go` | `backend/internal/users/repo.go` | T003 | Creates admin only if absent; logs create/skip | integration | Hardcoded creds operational risk |
| T006 | JWT auth and admin login endpoint | todo | P0 | solo-dev | Middleware + Strategy | Operation: `adminLogin`; Schemas: `AdminLoginRequest`,`AdminLoginResponse`,`ErrorResponse`; Security: `BearerAuth` | `backend/internal/auth/{handler,service,jwt.go}.go`, `backend/internal/middleware/auth.go` | `backend/internal/router/router.go` | T001,T005 | Login returns 10-day access and 3-month refresh token payload; admin routes protected | unit+integration+contract | Token misconfig blocks admin ops |
| T007 | Providers module CRUD + soft delete | todo | P0 | solo-dev | Repository | Operations: `listAdminProviders`,`createProvider`,`getProviderById`,`updateProvider`,`deleteProvider`,`listCustomerProviders`; Schemas: `Provider`,`ProviderCreateRequest`,`ProviderUpdateRequest`; Security: `BearerAuth` | `backend/internal/providers/{model,repo,service,handler}.go` | `backend/internal/router/router.go`, `backend/db/migrations/0002_indexes.sql` | T003,T006 | CRUD works, customer list excludes soft-deleted; unique name enforced | unit+integration+contract | Duplicate handling race conditions |
| T008 | Plans module CRUD + constraints (max 100/provider) | todo | P0 | solo-dev | Repository + Guard Clause | Operations: `listAdminPlans`,`createPlan`,`getPlanById`,`updatePlan`,`deletePlan`,`listCustomerPlansByProvider`; Schemas: `Plan`,`PlanCreateRequest`,`PlanUpdateRequest`; Security: `BearerAuth` | `backend/internal/plans/{model,repo,service,handler}.go` | `backend/internal/router/router.go` | T003,T006,T007 | Unique per provider and max 100/provider app-layer rule enforced | unit+integration+contract | Concurrent writes can bypass app check if not transactional |
| T009 | Customer checkout and service-request creation | todo | P0 | solo-dev | Service Layer + Idempotency | Operation: `createServiceRequestAndPayment`; Schemas: `PaymentRequest`,`PaymentResponse`,`ServiceRequest`,`User`,`ErrorResponse`; Security: none | `backend/internal/checkout/{handler,service}.go`, `backend/internal/servicerequests/{model,repo}.go` | `backend/internal/router/router.go` | T003,T007,T008 | Creates/upserts customer, creates pending service request, returns payment order payload | unit+integration+contract | Duplicate request/charge risk |
| T010 | reCAPTCHA v3 verification integration | todo | P0 | solo-dev | Adapter | Operations: `createServiceRequestAndPayment`,`retryPaymentForServiceRequest`; Schemas: `PaymentRequest`,`ErrorResponse` | `backend/internal/captcha/{client,service}.go` | `backend/internal/checkout/service.go` | T009 | Reject score <0.5 with exact error message | unit+integration | False positives may reduce conversion |
| T011 | Razorpay adapter, callback, and payment retry flow | todo | P0 | solo-dev | Adapter + Saga-lite | Operations: `handleRazorpayCallback`,`retryPaymentForServiceRequest`,`getCustomerPaymentStatus`; Schemas: `PaymentResponse`,`RazorpayCallbackRequest`,`ServiceRequest`,`ErrorResponse` | `backend/internal/payments/{razorpay_client,handler,service}.go`, `backend/internal/payments/retry_audit_repo.go` | `backend/internal/router/router.go`, `backend/internal/checkout/service.go` | T009,T010 | Callback signature verified, idempotent updates, same service request reused for retries | unit+integration+contract | Webhook replay and delayed callbacks |
| T012 | Service request state machine + admin status endpoint | todo | P0 | solo-dev | State Machine | Operations: `listServiceRequests`,`getServiceRequestById`,`updateServiceRequestStatus`,`getCustomerPaymentStatus`; Schemas: `ServiceRequest`,`ServiceRequestStatus`,`ServiceRequestStatusUpdateRequest`; Security: `BearerAuth` | `backend/internal/servicerequests/state_machine.go`, `backend/internal/servicerequests/handler.go` | `backend/internal/router/router.go` | T011,T006 | Only valid transitions accepted; invalid transitions return 422 | unit+integration+contract | Incorrect transitions corrupt lifecycle |
| T013 | SMS async jobs on completed status | todo | P0 | solo-dev | Outbox-lite + Worker | Operation: `updateServiceRequestStatus`; Schemas: `ServiceRequest`,`ErrorResponse`; Security: `BearerAuth` | `backend/internal/notifications/{sms_client,worker,jobs}.go`, `backend/internal/queue/asynq.go` | `backend/internal/servicerequests/service.go`, `backend/cmd/worker/main.go` | T012 | SMS enqueued only on transition to `completed`; retries 3x backoff | unit+integration | Notification failures without retries |
| T014 | Providers CSV upload (upsert) + async-ready abstraction | todo | P1 | solo-dev | Strategy | Operation: `uploadProvidersCsv`; Schemas: `CsvUploadResult`,`Provider`; Security: `BearerAuth` | `backend/internal/csv/providers_importer.go`, `backend/internal/csv/import_job_interface.go` | `backend/internal/providers/handler.go`, `backend/internal/router/router.go` | T007,T006 | Row-level validation + partial success response; upsert mode enabled | unit+integration+contract | CSV quality issues |
| T015 | Plans CSV upload (upsert) + async-ready abstraction | todo | P1 | solo-dev | Strategy | Operation: `uploadPlansCsv`; Schemas: `CsvUploadResult`,`Plan`; Security: `BearerAuth` | `backend/internal/csv/plans_importer.go` | `backend/internal/plans/handler.go`, `backend/internal/router/router.go` | T008,T006,T014 | Enforces provider existence, unique name/provider, max-100 rule | unit+integration+contract | Large file processing later |
| T016 | Customer frontend flows (provider->plan->checkout->status) | todo | P0 | solo-dev | Presenter + API Client | Operations: `listCustomerProviders`,`listCustomerPlansByProvider`,`createServiceRequestAndPayment`,`retryPaymentForServiceRequest`,`getCustomerPaymentStatus`; Schemas: `Provider`,`Plan`,`PaymentRequest`,`PaymentResponse`,`ServiceRequest` | `frontend/customer/pages/{index.tsx,providers/[id].tsx,checkout.tsx,payment-status/[id].tsx}`, `frontend/customer/lib/api.ts`, `frontend/customer/lib/recaptcha.ts` | `frontend/customer/package.json` | T007,T008,T009,T010,T011 | Customer can complete full purchase flow with reCAPTCHA and retry support | integration+contract (mock) | UX confusion in payment failures |
| T017 | Admin frontend flows (auth, CRUD, CSV, service requests) | todo | P0 | solo-dev | Container-Component + API Client | Operations: `adminLogin`,`listAdminProviders`,`createProvider`,`updateProvider`,`deleteProvider`,`uploadProvidersCsv`,`listAdminPlans`,`createPlan`,`updatePlan`,`deletePlan`,`uploadPlansCsv`,`listServiceRequests`,`getServiceRequestById`,`updateServiceRequestStatus`; Schemas: related admin schemas; Security: `BearerAuth` | `frontend/admin/pages/{login.tsx,providers.tsx,plans.tsx,service-requests.tsx}`, `frontend/admin/lib/{api.ts,auth.ts}` | `frontend/admin/package.json` | T006,T007,T008,T012,T014,T015 | Admin can perform all required operations end-to-end | integration+contract (mock) | Token handling bugs block admin |
| T018 | Logging, error envelope, and request correlation | todo | P0 | solo-dev | Middleware | Schemas: `ErrorResponse`,`Error`; Operations: all API operationIds | `backend/internal/middleware/{request_id,logging,error_handler}.go` | `backend/internal/router/router.go` | T001,T002 | All endpoints return standard error envelope; logs include request_id | unit+integration | Harder debugging if inconsistent |
| T019 | Backup/restore scripts + runbooks | todo | P0 | solo-dev | Scripted Procedure | Schemas: persisted core schemas (`User`,`Provider`,`Plan`,`ServiceRequest`) | `infra/scripts/{backup_db.sh,restore_db.sh}`, `infra/runbooks/backup_restore.md` | `README.md` | T003 | Backup script creates timestamped dump; restore procedure verified | integration (ops drill) | Data loss risk if untested |
| T020 | systemd + Nginx deployment assets | todo | P0 | solo-dev | Adapter (infra boundary) | Operations exposed under `/api/v1` including `adminLogin`; Security: `BearerAuth` | `infra/systemd/{dth-api.service,dth-worker.service}`, `infra/nginx/dth.conf`, `infra/deploy/deploy.sh` | `infra/README.md` | T016,T017,T018 | Services start/restart correctly; Nginx routes frontends and `/api/v1` | integration | Misrouting breaks prod |
| T021 | Backend unit tests (mandatory gate) | todo | P0 | solo-dev | Test Pyramid (Unit-heavy) | Operations/schemas covered: auth, providers, plans, payments, service requests, csv, error | `backend/internal/**/**/*_test.go` | `backend/Makefile` | T006,T007,T008,T009,T010,T011,T012,T013,T014,T015,T018 | Core logic unit tests passing | unit | Insufficient test depth |
| T022 | API contract tests against OpenAPI | todo | P0 | solo-dev | Consumer-Driven Contract | All operationIds + schemas + `BearerAuth` | `backend/tests/contract/{openapi_contract_test.go,fixtures/}.go` | `backend/api/openapi/openapi.yaml` | T021,T020 | Contract tests pass for required endpoints and status codes | contract+integration | Spec/impl drift |
| T023 | Release checklist and cutover execution | todo | P0 | solo-dev | Checklist | Operations: all production endpoints; Schemas: migration-related entities | `infra/release/release_checklist.md`, `infra/release/cutover.md` | `README.md` | T004,T005,T019,T020,T021,T022 | Release steps reproducible; rollback steps validated per milestone | integration (dry run) | Failed launch without runbook |

## Folder/File Blueprint

### Backend
- `backend/cmd/api/main.go`
- `backend/cmd/worker/main.go`
- `backend/cmd/migrate/main.go`
- `backend/cmd/seed-admin/main.go`
- `backend/internal/auth/*`
- `backend/internal/providers/*`
- `backend/internal/plans/*`
- `backend/internal/checkout/*`
- `backend/internal/payments/*`
- `backend/internal/servicerequests/*`
- `backend/internal/notifications/*`
- `backend/internal/csv/*`
- `backend/internal/middleware/*`
- `backend/internal/router/*`
- `backend/internal/errors/*`
- `backend/api/openapi/openapi.yaml`
- `backend/tests/contract/*`

### Frontend customer
- `frontend/customer/pages/index.tsx`
- `frontend/customer/pages/providers/[id].tsx`
- `frontend/customer/pages/checkout.tsx`
- `frontend/customer/pages/payment-status/[id].tsx`
- `frontend/customer/lib/api.ts`
- `frontend/customer/lib/recaptcha.ts`

### Frontend admin
- `frontend/admin/pages/login.tsx`
- `frontend/admin/pages/providers.tsx`
- `frontend/admin/pages/plans.tsx`
- `frontend/admin/pages/service-requests.tsx`
- `frontend/admin/lib/api.ts`
- `frontend/admin/lib/auth.ts`

### Infra/scripts/migrations
- `backend/db/migrations/*.sql`
- `infra/scripts/backup_db.sh`
- `infra/scripts/restore_db.sh`
- `infra/systemd/dth-api.service`
- `infra/systemd/dth-worker.service`
- `infra/nginx/dth.conf`
- `infra/deploy/deploy.sh`
- `infra/release/release_checklist.md`
- `infra/release/cutover.md`

## Migration & Seed Execution Plan
1. Put app in maintenance mode (if needed for schema-breaking releases).
2. Run DB backup: `infra/scripts/backup_db.sh`.
3. Execute migrations: `backend/cmd/migrate up`.
4. Verify migration version table and key constraints.
5. Run admin seed once: `backend/cmd/seed-admin`.
6. Restart services: `systemctl restart dth-api dth-worker`.
7. Run smoke tests for `adminLogin`, provider list, checkout create, callback endpoint.
8. Record release artifact version and migration versions in release log.

## Rollback Plan per milestone
- **M1 rollback:** revert code, drop newly created non-critical tables only in non-prod; in prod restore DB backup if migration applied.
- **M2 rollback:** disable new routes via router flags; restore previous binary; DB remains backward-compatible.
- **M3 rollback:** disable payment create/retry routes; keep service requests readable; restore pre-M3 DB backup if schema incompatible.
- **M4 rollback:** disable worker services and CSV endpoints; manual admin ops continue via existing CRUD.
- **M5 rollback:** redeploy previous frontend static bundles while keeping backend stable.
- **M6 rollback:** full release rollback by restoring previous binaries/config and last verified DB backup.

## Traceability Matrix
| task_id | PRD ID | FS ID | Arch ID | OpenAPI operationId/schema/security |
|---|---|---|---|---|
| T001 | FR-013, FR-016 | FS-014 | AE-002, AE-013 | `ErrorResponse`, `BearerAuth` |
| T002 | FR-013 | FS-005, FS-006 | AE-002 | all operationIds + all schemas + `BearerAuth` |
| T003 | FR-005, FR-007, FR-008, FR-009 | FS-003, FS-007, FS-008, FS-009 | AE-009, AE-010 | `User`,`Provider`,`Plan`,`ServiceRequest`,`PaymentResponse` |
| T004 | FR-014 | FS-013 | AE-011 | core persistence schemas |
| T005 | FR-014, FR-006 | FS-013, FS-006 | AE-011, AE-003 | `adminLogin`,`AdminLoginResponse`,`User`,`BearerAuth` |
| T006 | FR-006, FR-013 | FS-006 | AE-003, AE-002 | `adminLogin`,`AdminLoginRequest`,`AdminLoginResponse`,`BearerAuth` |
| T007 | FR-001, FR-007 | FS-001, FS-007 | AE-008, AE-009 | `listCustomerProviders`,`listAdminProviders`,`createProvider`,`getProviderById`,`updateProvider`,`deleteProvider`,`Provider*`,`BearerAuth` |
| T008 | FR-002, FR-008 | FS-002, FS-008 | AE-008, AE-009 | `listCustomerPlansByProvider`,`listAdminPlans`,`createPlan`,`getPlanById`,`updatePlan`,`deletePlan`,`Plan*`,`BearerAuth` |
| T009 | FR-003, FR-005 | FS-003 | AE-010 | `createServiceRequestAndPayment`,`PaymentRequest`,`PaymentResponse`,`ServiceRequest` |
| T010 | FR-004 | FS-004 | AE-004 | `createServiceRequestAndPayment`,`retryPaymentForServiceRequest`,`PaymentRequest`,`ErrorResponse` |
| T011 | FR-003, FR-013 | FS-005 | AE-005, AE-010 | `handleRazorpayCallback`,`retryPaymentForServiceRequest`,`getCustomerPaymentStatus`,`RazorpayCallbackRequest`,`PaymentResponse`,`ServiceRequest` |
| T012 | FR-009 | FS-009 | AE-006 | `listServiceRequests`,`getServiceRequestById`,`updateServiceRequestStatus`,`ServiceRequestStatus*`,`BearerAuth` |
| T013 | FR-010 | FS-010 | AE-007 | `updateServiceRequestStatus`,`ServiceRequest`,`BearerAuth` |
| T014 | FR-011 | FS-011 | AE-008 | `uploadProvidersCsv`,`CsvUploadResult`,`Provider`,`BearerAuth` |
| T015 | FR-012 | FS-012 | AE-008 | `uploadPlansCsv`,`CsvUploadResult`,`Plan`,`BearerAuth` |
| T016 | FR-001..FR-004 | FS-001..FS-005 | AE-001, AE-004, AE-005 | customer operationIds + `Provider`,`Plan`,`PaymentRequest`,`PaymentResponse`,`ServiceRequest` |
| T017 | FR-006..FR-012 | FS-006..FS-012 | AE-001, AE-003, AE-006, AE-008 | admin operationIds + `BearerAuth` + admin schemas |
| T018 | FR-013, FR-016 | FS-005, FS-006, FS-014 | AE-002, AE-013 | all operationIds + `ErrorResponse` |
| T019 | FR-015 | FS-014 | AE-012 | persistence schemas (`User`,`Provider`,`Plan`,`ServiceRequest`) |
| T020 | FR-016, FR-013 | FS-014 | AE-013, AE-002 | `/api/v1` operations + `BearerAuth` enforcement path |
| T021 | FR-017 | FS-ALL | AE-014 | all critical operationIds/schemas |
| T022 | FR-013, FR-017 | FS-005, FS-006 | AE-002, AE-014 | all operationIds + all schemas + `BearerAuth` |
| T023 | FR-014, FR-015, FR-016 | FS-013, FS-014 | AE-011, AE-012, AE-013 | release-critical operations and persistence schemas |

## Assumptions
- Solo developer executes tasks sequentially with minimal parallel streams.
- Existing repository can host `backend/`, `frontend/customer`, and `frontend/admin` structure.
- Razorpay and SMS credentials are available before M3.
- OTP remains deferred; guest checkout is acceptable for v1.
- Refresh token is returned by login response but no dedicated refresh endpoint in v1.

## Open Questions
- For async CSV future mode, should same endpoints return `202 + job_id` when file exceeds threshold?
    Answer: No, it is not required
- Should admin login support username only or username+phone equally from day one?
    Answer: username
- Do we need customer-facing service-request history page in v1 or only payment-status polling?
    Answer: No, we don't need customer facing service request history page

## Risks
- Payment callback timing and retries can create status drift if idempotency is incomplete.
- Long admin token validity increases exposure if token storage is weak.
- App-layer max-plan constraint can fail under race unless DB transaction locks are used.
- Single-node deployment increases outage impact.
- Manual credential and backup procedures can fail without routine drills.

## Decision Log
| Decision ID | Date | Decision | Options Considered | Chosen Option | Rationale | Impact |
|---|---|---|---|---|---|---|
| IMP-001 | 2026-02-27 | Build sequence | UI-first vs contract-first | Contract-first | Reduces rework and integration mismatches | Faster stable delivery |
| IMP-002 | 2026-02-27 | CSV architecture | Sync-only vs async-only vs sync+async-ready | Sync now + async-ready abstraction | Matches low scale now and future requirement | Minimal v1 overhead |
| IMP-003 | 2026-02-27 | Payment retry model | New request vs same request | Same service request + retry audit | Approved FSD decision | Clean traceability |
| IMP-004 | 2026-02-27 | Constraint enforcement | DB trigger vs app-layer for max plans | App-layer only | Approved TAD decision | Simpler migrations, higher concurrency care |
| IMP-005 | 2026-02-27 | Deployment topology | Managed multi-node vs single OCI VM | Single VM | Cost and solo-dev constraints | Accept SPOF risk |
| IMP-006 | 2026-02-27 | Auth refresh endpoint | Include now vs defer | Defer | Explicit OpenAPI decision | Simpler v1 auth surface |

```

Output format (Markdown):
- Title + Metadata
- Prerequisites (exact versions/tools)
- Repository Structure Overview
- Local Setup
  - Clone/install
  - Env vars (`.env.example` contract)
  - DB setup (PostgreSQL)
  - Redis setup
  - Migration commands
  - Seed command (admin user)
- Runbook Commands
  - Backend run/test/lint
  - Frontend customer run/build
  - Frontend admin run/build
  - Worker process (Asynq) run
- Integration Setup
  - Payment sandbox config
  - SMS provider config
  - reCAPTCHA v3 config
- API Contract Workflow
  - OpenAPI validation, codegen (if used), contract testing
- Troubleshooting Guide (top 10 likely issues + fixes)
- Production-like Local Checks before merge
- Security and Secrets Handling Guidelines
- Traceability Matrix
  - Setup step -> Architecture section -> Implementation task_id -> OpenAPI relevance
- Assumptions
- Open Questions
- Risks
- Decision Log table

Validation checklist:
- Steps are executable end-to-end on a fresh machine.
- Commands and env vars are concrete and consistent with architecture.
- Includes migration + seed + rollback verification basics.
- Includes integration setup for payment/SMS/reCAPTCHA.
- Every major setup section traces to architecture and implementation tasks.
- Assumptions and unresolved questions explicit.

Quality bar:
- Zero ambiguity in critical steps.
- Suitable for onboarding a new engineer with minimal assistance.
- Operationally safe defaults and clear troubleshooting.
- Production-minded despite single-prod environment.