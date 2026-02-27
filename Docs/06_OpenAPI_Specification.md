```yaml
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

### Assumptions
- Customer checkout is guest-first; customer identity is captured via name+strict-unique phone number.
- Canonical API status values use snake_case (`payment_success`) while UI may render human-readable labels.
- Razorpay callback endpoint shape may vary by event; `payload` is intentionally flexible.
- CSV uploads are API-based multipart file operations with upsert semantics.

### Open Questions
- Should `/customer/service-requests/{id}/payment-status` require additional tokenized proof beyond phone number?
  Answer: No phone number should be fine
- Is token refresh endpoint required in v1 (currently omitted; login re-auth is fallback)?
  Answer: No required
- Do admin CSV uploads need asynchronous processing for larger files later?
  Answer: yes
- Should `users.phone_number` uniqueness include soft-deleted users or active-only scope?
  Answer: `users.phone_number` uniqueness should include soft-deleted users

### Traceability Matrix (OperationId/Schema -> PRD ID/FS ID/Arch ID)
| API Element | PRD ID | FS ID | Arch ID |
|---|---|---|---|
| `adminLogin`, `BearerAuth` | FR-006, FR-013 | FS-006 | AE-003, AE-002 |
| `listCustomerProviders` + `Provider` | FR-001 | FS-001 | AE-001 |
| `listCustomerPlansByProvider` + `Plan` | FR-002 | FS-002 | AE-001 |
| `createServiceRequestAndPayment`, `retryPaymentForServiceRequest`, `PaymentRequest`, `PaymentResponse` | FR-003, FR-004, FR-005 | FS-003, FS-004, FS-005 | AE-004, AE-005, AE-010 |
| `handleRazorpayCallback` | FR-003, FR-013 | FS-005 | AE-005 |
| `listAdminProviders`, `createProvider`, `updateProvider`, `deleteProvider` | FR-007 | FS-007 | AE-008, AE-009 |
| `uploadProvidersCsv`, `CsvUploadResult` | FR-011 | FS-011 | AE-008 |
| `listAdminPlans`, `createPlan`, `updatePlan`, `deletePlan` | FR-008 | FS-008 | AE-008, AE-009 |
| `uploadPlansCsv`, `CsvUploadResult` | FR-012 | FS-012 | AE-008 |
| `listServiceRequests`, `getServiceRequestById`, `updateServiceRequestStatus`, `ServiceRequest` | FR-009 | FS-009 | AE-006 |
| `ServiceRequestStatusUpdateRequest` (completion transition behavior) | FR-010 | FS-010 | AE-007 |
| `ErrorResponse` | FR-013 | FS-005, FS-006 | AE-002 |

### Risks
- Webhook signature mismatch or callback timing issues can lead to delayed status convergence.
- Guest payment-status polling with phone number may expose enumeration risk if not rate-limited.
- App-layer-only max-plan constraint can drift under concurrent writes unless transaction checks are strict.
- Manual credential handling and long token TTLs increase blast radius if admin token is leaked.

### Decision Log
| Decision ID | Date | Decision | Options Considered | Chosen Option | Rationale | Impact |
|---|---|---|---|---|---|---|
| API-D1 | 2026-02-27 | API versioning strategy | Header versioning vs URI versioning | URI `/api/v1` | Aligns approved PRD/TAD | Clear contract evolution |
| API-D2 | 2026-02-27 | Auth model for admin APIs | Session vs JWT | Bearer JWT | Explicit constraint | Stateless auth middleware |
| API-D3 | 2026-02-27 | Payment gateway contract | Generic abstraction only vs concrete provider endpoint | Razorpay callback + payment operations | Resolved OQ | Implementation-ready integrations |
| API-D4 | 2026-02-27 | Payment idempotency | Optional vs required header | Required `Idempotency-Key` for payment-create/retry | Prevent duplicate charges | Safer retries |
| API-D5 | 2026-02-27 | Service request retry model | New request each retry vs same request | Same request + new payment attempt | Approved FSD decision | Better traceability |
| API-D6 | 2026-02-27 | CSV ingestion style | Out-of-band job only vs sync multipart endpoint | Sync multipart endpoint with row-level results | Low-scale + solo-dev practicality | Faster admin feedback |
