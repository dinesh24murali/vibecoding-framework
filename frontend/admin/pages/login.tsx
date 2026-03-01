import Head from "next/head"
import { useState } from "react"
import { useRouter } from "next/router"
import { adminLogin } from "../lib/api"
import { setAccessToken } from "../lib/auth"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export default function LoginPage() {
  const router = useRouter()
  const [identifier, setIdentifier] = useState("")
  const [password, setPassword] = useState("")
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState("")

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setSubmitting(true)
    setError("")
    try {
      const data = await adminLogin(identifier.trim(), password)
      setAccessToken(data.access_token)
      await router.push("/providers")
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <>
      <Head>
        <title>Village DTH Admin - Login</title>
      </Head>
      <main className="grid min-h-screen place-items-center bg-[linear-gradient(150deg,#f2f7ff,#fff3e3_55%,#edf8f0)] px-4 py-8">
        <Card className="w-full max-w-md bg-white/90">
          <CardHeader>
            <CardTitle>Admin Login</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="space-y-4" onSubmit={onSubmit}>
              <div className="space-y-2">
                <Label htmlFor="identifier">Username</Label>
                <Input id="identifier" value={identifier} onChange={(e) => setIdentifier(e.target.value)} required />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>

              {error && <p className="text-sm text-destructive">{error}</p>}

              <Button type="submit" disabled={submitting} className="w-full">
                {submitting ? "Signing in..." : "Sign in"}
              </Button>
            </form>
          </CardContent>
        </Card>
      </main>
    </>
  )
}

export async function getStaticProps() {
  return { props: {} }
}
