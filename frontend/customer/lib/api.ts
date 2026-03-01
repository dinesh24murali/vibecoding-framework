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

export type PaymentResponse = {
  service_request: ServiceRequest
  payment: {
    gateway: "razorpay"
    order_id: string
    amount: number
    currency: "INR"
    status: "created" | "authorized" | "captured" | "failed"
  }
}

export type ErrorEnvelope = {
  error?: {
    code?: string
    message?: string
  }
}

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1"

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {})
    }
  })

  if (!response.ok) {
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

  return (await response.json()) as T
}

export async function listCustomerProviders(): Promise<Provider[]> {
  const data = await request<{ items: Provider[] }>("/customer/providers")
  return data.items
}

export async function listCustomerPlansByProvider(providerId: string): Promise<Plan[]> {
  const data = await request<{ items: Plan[] }>(`/customer/providers/${providerId}/plans`)
  return data.items
}

export async function createServiceRequestAndPayment(input: {
  providerId: number
  planId: number
  customerName: string
  customerPhone: string
  recaptchaToken: string
  idempotencyKey: string
}): Promise<PaymentResponse> {
  return request<PaymentResponse>("/customer/checkout/service-requests", {
    method: "POST",
    headers: {
      "Idempotency-Key": input.idempotencyKey
    },
    body: JSON.stringify({
      provider_id: input.providerId,
      plan_id: input.planId,
      customer: {
        name: input.customerName,
        phone_number: input.customerPhone
      },
      recaptcha_token: input.recaptchaToken
    })
  })
}

export async function retryPaymentForServiceRequest(input: {
  serviceRequestId: string
  recaptchaToken: string
  idempotencyKey: string
}): Promise<PaymentResponse> {
  return request<PaymentResponse>(`/customer/service-requests/${input.serviceRequestId}/retry-payment`, {
    method: "POST",
    headers: {
      "Idempotency-Key": input.idempotencyKey
    },
    body: JSON.stringify({
      recaptcha_token: input.recaptchaToken
    })
  })
}

export async function getCustomerPaymentStatus(serviceRequestId: string, phoneNumber: string): Promise<ServiceRequest> {
  const data = await request<{ service_request: ServiceRequest }>(
    `/customer/service-requests/${serviceRequestId}/payment-status?phone_number=${encodeURIComponent(phoneNumber)}`
  )
  return data.service_request
}
