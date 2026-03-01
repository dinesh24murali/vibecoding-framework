import Head from "next/head"
import Link from "next/link"
import { useRouter } from "next/router"
import { useEffect, useMemo, useState } from "react"
import { listCustomerPlansByProvider, type Plan } from "../lib/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"

export default function ProviderPlansPage() {
  const router = useRouter()
  const providerId = useMemo(() => {
    const raw = router.query.id
    return Array.isArray(raw) ? raw[0] : raw
  }, [router.query.id])

  const [plans, setPlans] = useState<Plan[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    if (!providerId) {
      return
    }

    setLoading(true)
    listCustomerPlansByProvider(providerId)
      .then(setPlans)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false))
  }, [providerId])

  return (
    <>
      <Head>
        <title>Village DTH - Plans</title>
      </Head>
      <main className="min-h-screen bg-[linear-gradient(180deg,#eef6ff,#fef9eb)] px-4 py-8">
        <div className="mx-auto mb-4 flex max-w-5xl items-center justify-between gap-3">
          <Button asChild variant="outline">
            <Link href="/">Back to providers</Link>
          </Button>
          <h1 className="text-3xl font-semibold tracking-tight">Plans</h1>
        </div>

        {loading && <p className="mx-auto max-w-5xl text-sm text-muted-foreground">Loading plans...</p>}
        {error && <p className="mx-auto mb-4 max-w-5xl text-sm text-destructive">{error}</p>}

        <section className="mx-auto grid max-w-5xl grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {plans.map((plan) => (
            <Card key={plan.id} className="bg-white/90">
              <CardHeader>
                <CardTitle className="text-xl">{plan.name}</CardTitle>
                <CardDescription>{plan.description}</CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold text-primary">Rs {Math.max(plan.price - plan.discount, 0.01).toFixed(2)}</p>
              </CardContent>
              <CardFooter>
                <Button
                  type="button"
                  className="w-full"
                  onClick={() =>
                    router.push(`/checkout?providerId=${encodeURIComponent(plan.provider_id)}&planId=${encodeURIComponent(plan.id)}`)
                  }
                >
                  Choose this plan
                </Button>
              </CardFooter>
            </Card>
          ))}
        </section>
      </main>
    </>
  )
}

export async function getStaticProps() {
  return { props: {} }
}
