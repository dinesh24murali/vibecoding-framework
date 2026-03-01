import Head from "next/head"
import { useEffect, useState } from "react"
import { useRouter } from "next/router"
import { listCustomerProviders, type Provider } from "../lib/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function HomePage() {
  const router = useRouter()
  const [providers, setProviders] = useState<Provider[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    listCustomerProviders()
      .then(setProviders)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  return (
    <>
      <Head>
        <title>Village DTH - Choose Provider</title>
      </Head>
      <main className="min-h-screen bg-[radial-gradient(75rem_50rem_at_10%_0%,#fff4da,#f1f8ff_45%,#edf6ef_100%)] px-4 py-10">
        <section className="mx-auto max-w-5xl text-center">
          <p className="text-xs uppercase tracking-[0.16em] text-amber-700">Village DTH</p>
          <h1 className="mt-2 text-4xl font-bold tracking-tight sm:text-5xl">Pick your TV provider</h1>
          <p className="mt-2 text-muted-foreground">Browse active providers and move to plan selection in one tap.</p>
        </section>

        {loading && <p className="mx-auto mt-6 max-w-5xl text-sm text-muted-foreground">Loading providers...</p>}
        {error && <p className="mx-auto mt-6 max-w-5xl text-sm text-destructive">{error}</p>}

        <section className="mx-auto mt-7 grid max-w-5xl grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {providers.map((provider) => (
            <Card key={provider.id} className="bg-white/90">
              <CardHeader>
                <CardTitle className="text-xl">{provider.name}</CardTitle>
                <CardDescription>See all active plans for this provider.</CardDescription>
              </CardHeader>
              <CardContent>
                <Button
                  type="button"
                  className="w-full"
                  onClick={() => router.push(`/providers?id=${encodeURIComponent(provider.id)}`)}
                >
                  View plans
                </Button>
              </CardContent>
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
