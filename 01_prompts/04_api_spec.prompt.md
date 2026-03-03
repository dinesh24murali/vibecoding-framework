# Prompt: Generate OpenAPI Specification

## Instructions

Copy everything below the `---` line into your LLM. Replace the placeholders with your previously generated documents.

---

You are a senior backend engineer writing an OpenAPI 3.1 specification.

Here is the project context:

**Intake Answers:**
```
{{INTAKE_ANSWERS}}
```

**Functional Specification:**
```
{{FUNCTIONAL_SPEC}}
```

**Technical Architecture:**
```
{{TECH_ARCHITECTURE}}
```

Generate a complete OpenAPI 3.1 specification in YAML format. This spec should be valid and importable into tools like Swagger UI, Postman, or Stoplight.

## Requirements

### General
- OpenAPI version: 3.1.0
- Include `info` block with title, description, version, contact
- Include `servers` block for local, staging, and production
- Use `tags` to group endpoints by feature/resource

### Paths
For EVERY feature in the functional spec that requires an API endpoint:

- Use RESTful conventions (plural nouns, proper HTTP methods)
- Include all CRUD operations where applicable
- Specify `operationId` for each endpoint
- Include `summary` and `description`
- Define path parameters, query parameters, and request bodies with:
  - Type and format
  - Required/optional
  - Validation constraints (minLength, maxLength, pattern, enum, minimum, maximum)
  - Example values
- Define all response codes (200, 201, 400, 401, 403, 404, 409, 422, 500)
- Use `$ref` for shared schemas

### Schemas (Components)
- Define a schema for each data entity
- Include all fields with types, formats, and descriptions
- Mark required fields
- Include example values
- Use composition (`allOf`, `oneOf`) where appropriate
- Define common schemas:
  - `Error` (standard error response envelope)
  - `PaginatedResponse` (generic pagination wrapper)
  - `SuccessResponse` (generic success wrapper)

### Security
- Define security schemes (Bearer JWT, API Key, OAuth2 — whatever fits the architecture)
- Apply security globally and override per-endpoint where needed

### Standard Patterns
Include these standard endpoints:
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`
- `GET /health` (public, no auth)
- `GET /health/ready` (detailed, admin only)

### Pagination Convention
Use consistent pagination across all list endpoints:
```yaml
parameters:
  - name: page
    in: query
    schema:
      type: integer
      minimum: 1
      default: 1
  - name: per_page
    in: query
    schema:
      type: integer
      minimum: 1
      maximum: 100
      default: 20
  - name: sort_by
    in: query
    schema:
      type: string
  - name: sort_order
    in: query
    schema:
      type: string
      enum: [asc, desc]
      default: desc
```

### Error Response Convention
```yaml
Error:
  type: object
  required: [error]
  properties:
    error:
      type: object
      required: [code, message]
      properties:
        code:
          type: string
          description: Machine-readable error code
          example: "VALIDATION_ERROR"
        message:
          type: string
          description: Human-readable error message
          example: "The email field is required"
        details:
          type: array
          items:
            type: object
            properties:
              field:
                type: string
              message:
                type: string
```

## Output Format
- Valid OpenAPI 3.1.0 YAML
- Well-organized with comments separating major sections
- All `$ref` paths must resolve correctly
- Target: complete and importable — not a skeleton
