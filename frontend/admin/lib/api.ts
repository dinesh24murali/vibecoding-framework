import { clearAccessToken, getAccessToken } from "./auth"

type ErrorEnvelope = {
  error?: {
    message?: string
  }
}

export type Provider = {
  id: string
  name: string
  image_url?: string | null
  created_at: string
  updated_at: string
}

export type Plan = {
  id: string
  provider_id: string
  name: string
  description: string
  price: number
  discount: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export type ServiceRequest = {
  id: string
  user_id: string
  plan_id: string
  status: "pending" | "payment_failed" | "payment_success" | "blocked" | "completed"
  created_at: string
  updated_at: string
}

export type ServiceRequestDetails = ServiceRequest & {
  user: {
    id: string
    name: string
    phone_number: string
    role: string
    status: string
    created_at: string
    updated_at: string
  }
  plan: {
    id: string
    provider_id: string
    name: string
    description: string
    price: number
    discount: number
    is_active: boolean
    created_at: string
    updated_at: string
  }
  provider: {
    id: string
    name: string
    image_url?: string | null
    created_at: string
    updated_at: string
  }
}

export type CsvUploadResult = {
  total_rows: number
  success_rows: number
  failed_rows: number
  errors: Array<{ row: number; field?: string; issue: string }>
}

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1"

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const token = getAccessToken()
  const headers = new Headers(init?.headers ?? {})
  if (!headers.has("Content-Type") && !(init?.body instanceof FormData)) {
    headers.set("Content-Type", "application/json")
  }
  if (token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers
  })

  if (!response.ok) {
    if (response.status === 401) {
      clearAccessToken()
    }
    let message = `Request failed with status ${response.status}`
    try {
      const parsed = (await response.json()) as ErrorEnvelope
      if (parsed.error?.message) {
        message = parsed.error.message
      }
    } catch {
      // noop
    }
    throw new Error(message)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export async function adminLogin(identifier: string, password: string): Promise<{ access_token: string }> {
  return request<{ access_token: string }>("/auth/admin/login", {
    method: "POST",
    body: JSON.stringify({ identifier, password })
  })
}

export async function listProviders(): Promise<Provider[]> {
  const data = await request<{ items: Provider[] }>("/admin/providers?page=1&page_size=100")
  return data.items
}

export async function createProvider(payload: { name: string; image_url?: string | null }): Promise<Provider> {
  return request<Provider>("/admin/providers", {
    method: "POST",
    body: JSON.stringify(payload)
  })
}

export async function updateProvider(providerId: string, payload: { name?: string; image_url?: string | null }): Promise<Provider> {
  return request<Provider>(`/admin/providers/${providerId}`, {
    method: "PATCH",
    body: JSON.stringify(payload)
  })
}

export async function deleteProvider(providerId: string): Promise<void> {
  await request<void>(`/admin/providers/${providerId}`, { method: "DELETE" })
}

export async function uploadProvidersCsv(file: File): Promise<CsvUploadResult> {
  const formData = new FormData()
  formData.append("file", file)
  return request<CsvUploadResult>("/admin/providers/csv-upload", {
    method: "POST",
    body: formData
  })
}

export async function listPlans(): Promise<Plan[]> {
  const data = await request<{ items: Plan[] }>("/admin/plans?page=1&page_size=100")
  return data.items
}

export async function createPlan(payload: {
  provider_id: number
  name: string
  description: string
  price: number
  discount: number
  is_active: boolean
}): Promise<Plan> {
  return request<Plan>("/admin/plans", {
    method: "POST",
    body: JSON.stringify(payload)
  })
}

export async function updatePlan(planId: string, payload: Partial<Omit<Plan, "id" | "created_at" | "updated_at">>): Promise<Plan> {
  return request<Plan>(`/admin/plans/${planId}`, {
    method: "PATCH",
    body: JSON.stringify(payload)
  })
}

export async function deletePlan(planId: string): Promise<void> {
  await request<void>(`/admin/plans/${planId}`, { method: "DELETE" })
}

export async function uploadPlansCsv(file: File): Promise<CsvUploadResult> {
  const formData = new FormData()
  formData.append("file", file)
  return request<CsvUploadResult>("/admin/plans/csv-upload", {
    method: "POST",
    body: formData
  })
}

export async function listServiceRequests(): Promise<ServiceRequest[]> {
  const data = await request<{ items: ServiceRequest[] }>("/admin/service-requests?page=1&page_size=100")
  return data.items
}

export async function getServiceRequestById(id: string): Promise<ServiceRequestDetails> {
  return request<ServiceRequestDetails>(`/admin/service-requests/${id}`)
}

export async function updateServiceRequestStatus(id: string, status: ServiceRequest["status"]): Promise<ServiceRequest> {
  return request<ServiceRequest>(`/admin/service-requests/${id}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status })
  })
}
