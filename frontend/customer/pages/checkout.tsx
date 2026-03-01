import Head from "next/head"
import Link from "next/link"
import { useRouter } from "next/router"
import { useEffect, useMemo, useState } from "react"
import {
  createServiceRequestAndPayment,
  listCustomerPlansByProvider,
  type PaymentResponse,
  type Plan
} from "../lib/api"
import { getRecaptchaToken, loadRecaptchaScript } from "../lib/recaptcha"
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

export default function CheckoutPage() {
  const router = useRouter()
  const providerId = useMemo(() => {
    const raw = router.query.providerId
    return Array.isArray(raw) ? raw[0] : raw
  }, [router.query.providerId])
  const planId = useMemo(() => {
    const raw = router.query.planId
    return Array.isArray(raw) ? raw[0] : raw
  }, [router.query.planId])

  const [plan, setPlan] = useState<Plan | null>(null)
  const [name, setName] = useState("")
  const [phoneNumber, setPhoneNumber] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")
  const [lastResponse, setLastResponse] = useState<PaymentResponse | null>(null)

  useEffect(() => {
    loadRecaptchaScript()
  }, [])

  useEffect(() => {
    if (!providerId || !planId) {
      return
    }
    setError("")
    listCustomerPlansByProvider(providerId)
      .then((plans) => {
        const selectedPlan = plans.find((item) => String(item.id) === String(planId)) ?? null
        setPlan(selectedPlan)
        if (!selectedPlan) {
          setError("Selected plan is unavailable. Please choose a plan again.")
        }
      })
      .catch((err: Error) => setError(err.message))
  }, [providerId, planId])

  const payable = plan ? Math.max(plan.price - plan.discount, 0.01) : 0

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!providerId || !planId) {
      setError("providerId and planId are required")
      return
    }

    setError("")
    setSubmitting(true)
    try {
      const recaptchaToken = await getRecaptchaToken("checkout_submit")
      const response = await createServiceRequestAndPayment({
        providerId: Number(providerId),
        planId: Number(planId),
        customerName: name.trim(),
        customerPhone: phoneNumber.trim(),
        recaptchaToken,
        idempotencyKey: createIdempotencyKey()
      })
      setLastResponse(response)
      await router.push(
        `/payment-status?id=${encodeURIComponent(response.service_request.id)}&phone_number=${encodeURIComponent(phoneNumber.trim())}`
      )
    } catch (err) {
      setError(err instanceof Error ? err.message : "Checkout failed")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <>
      <Head>
        <title>Village DTH - Checkout</title>
      </Head>
      <main className="min-h-screen bg-[radial-gradient(80rem_50rem_at_90%_10%,#dff7ff,#fff2dc_45%,#edf9ee_100%)] px-4 py-8">
        <div className="mx-auto grid max-w-6xl grid-cols-1 gap-4 md:grid-cols-2">
          <Card className="bg-white/90">
            <CardHeader>
              <Button asChild variant="outline" className="w-fit">
                <Link href={providerId ? `/providers?id=${encodeURIComponent(providerId)}` : "/"}>Back to plans</Link>
              </Button>
              <CardTitle className="text-3xl">Checkout</CardTitle>
              {plan ? (
                <CardDescription>{plan.description}</CardDescription>
              ) : (
                <CardDescription>Select a plan to continue.</CardDescription>
              )}
            </CardHeader>
            {plan && (
              <CardContent className="space-y-2">
                <p className="text-lg font-semibold">{plan.name}</p>
                <p className="text-3xl font-bold text-primary">Rs {payable.toFixed(2)}</p>
              </CardContent>
            )}
          </Card>

          <Card className="bg-white/90">
            <CardHeader>
              <CardTitle>Customer details</CardTitle>
            </CardHeader>
            <CardContent>
              <form className="space-y-4" onSubmit={onSubmit}>
                <div className="space-y-2">
                  <Label htmlFor="name">Full name</Label>
                  <Input id="name" value={name} onChange={(e) => setName(e.target.value)} required maxLength={120} />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="phone">Phone number</Label>
                  <Input
                    id="phone"
                    value={phoneNumber}
                    onChange={(e) => setPhoneNumber(e.target.value)}
                    pattern="^[6-9][0-9]{9}$"
                    required
                  />
                </div>

                {error && <p className="text-sm text-destructive">{error}</p>}

                <Button type="submit" className="w-full" disabled={submitting || !plan}>
                  {submitting ? "Creating order..." : "Proceed to payment"}
                </Button>

                {lastResponse && (
                  <p className="text-sm text-emerald-700">
                    Order {lastResponse.payment.order_id} created. You are being redirected to status page.
                  </p>
                )}
              </form>
            </CardContent>
          </Card>
        </div>
      </main>
    </>
  )
}

export async function getStaticProps() {
  return { props: {} }
}
