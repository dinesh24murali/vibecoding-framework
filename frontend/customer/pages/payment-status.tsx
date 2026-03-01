import Head from "next/head"
import Link from "next/link"
import { useRouter } from "next/router"
import { useEffect, useMemo, useState } from "react"
import { getCustomerPaymentStatus, retryPaymentForServiceRequest, type ServiceRequest } from "../lib/api"
import { getRecaptchaToken, loadRecaptchaScript } from "../lib/recaptcha"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

function createIdempotencyKey() {
  if (typeof crypto !== "undefined" && crypto.randomUUID) {
    return crypto.randomUUID()
  }
  return `idem-${Date.now()}`
}

function getStatusClass(status: ServiceRequest["status"]) {
  switch (status) {
    case "payment_failed":
      return "bg-red-100 text-red-700"
    case "pending":
      return "bg-amber-100 text-amber-700"
    case "payment_success":
    case "completed":
      return "bg-emerald-100 text-emerald-700"
    case "blocked":
      return "bg-rose-100 text-rose-700"
    default:
      return ""
  }
}

export default function PaymentStatusPage() {
  const router = useRouter()
  const serviceRequestId = useMemo(() => {
    const raw = router.query.id
    return Array.isArray(raw) ? raw[0] : raw
  }, [router.query.id])

  const [phoneNumber, setPhoneNumber] = useState("")
  const [status, setStatus] = useState<ServiceRequest | null>(null)
  const [loading, setLoading] = useState(false)
  const [retrying, setRetrying] = useState(false)
  const [error, setError] = useState("")

  useEffect(() => {
    loadRecaptchaScript()
  }, [])

  useEffect(() => {
    const rawPhone = router.query.phone_number
    const phone = Array.isArray(rawPhone) ? rawPhone[0] : rawPhone
    if (phone) {
      setPhoneNumber(phone)
    }
  }, [router.query.phone_number])

  const refresh = async () => {
    if (!serviceRequestId || !phoneNumber.trim()) {
      setError("Phone number is required")
      return
    }
    setError("")
    setLoading(true)
    try {
      const nextStatus = await getCustomerPaymentStatus(serviceRequestId, phoneNumber.trim())
      setStatus(nextStatus)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch payment status")
    } finally {
      setLoading(false)
    }
  }

  const retry = async () => {
    if (!serviceRequestId) {
      setError("Service request id is required")
      return
    }
    setRetrying(true)
    setError("")
    try {
      const recaptchaToken = await getRecaptchaToken("retry_payment")
      await retryPaymentForServiceRequest({
        serviceRequestId,
        recaptchaToken,
        idempotencyKey: createIdempotencyKey()
      })
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Retry failed")
    } finally {
      setRetrying(false)
    }
  }

  return (
    <>
      <Head>
        <title>Village DTH - Payment Status</title>
      </Head>
      <main className="grid min-h-screen place-items-center bg-[linear-gradient(145deg,#f5f0ff,#f1fcf5_45%,#eef8ff)] px-4 py-8">
        <Card className="w-full max-w-xl bg-white/90">
          <CardHeader>
            <Button asChild variant="outline" className="w-fit">
              <Link href="/">Home</Link>
            </Button>
            <CardTitle>Payment status</CardTitle>
            <CardDescription>Track your checkout using service request id and phone number.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="phone">Phone number</Label>
              <Input
                id="phone"
                value={phoneNumber}
                onChange={(event) => setPhoneNumber(event.target.value)}
                pattern="^[6-9][0-9]{9}$"
                placeholder="10-digit number"
              />
            </div>

            <div className="flex flex-wrap gap-2">
              <Button type="button" onClick={refresh} disabled={loading}>
                {loading ? "Checking..." : "Check status"}
              </Button>
              {status?.status === "payment_failed" && (
                <Button type="button" onClick={retry} disabled={retrying} variant="secondary">
                  {retrying ? "Retrying..." : "Retry payment"}
                </Button>
              )}
            </div>

            {status && (
              <div className="rounded-md border bg-white p-4 text-sm">
                <p>
                  <strong>Service request:</strong> {status.id}
                </p>
                <p className="mt-2 flex items-center gap-2">
                  <strong>Status:</strong>
                  <Badge className={getStatusClass(status.status)} variant="outline">
                    {status.status}
                  </Badge>
                </p>
                <p className="mt-2">
                  <strong>Updated at:</strong> {new Date(status.updated_at).toLocaleString()}
                </p>
              </div>
            )}

            {error && <p className="text-sm text-destructive">{error}</p>}
          </CardContent>
        </Card>
      </main>
    </>
  )
}

export async function getStaticProps() {
  return { props: {} }
}
