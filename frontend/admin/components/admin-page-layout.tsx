import Link from "next/link"
import { useRouter } from "next/router"
import { ReactNode } from "react"
import { clearAccessToken } from "@/lib/auth"
import { Button } from "@/components/ui/button"

type AdminPageLayoutProps = {
  title: string
  children: ReactNode
}

export function AdminPageLayout({ title, children }: AdminPageLayoutProps) {
  const router = useRouter()

  return (
    <main className="min-h-screen bg-[linear-gradient(160deg,#f2f8ff,#fff5e8_50%,#edf9f0)] px-4 py-6">
      <header className="mx-auto mb-4 flex w-full max-w-6xl flex-wrap items-center justify-between gap-3">
        <h1 className="text-3xl font-semibold tracking-tight">{title}</h1>
        <nav className="flex flex-wrap items-center gap-2">
          <Button asChild variant="outline">
            <Link href="/providers">Providers</Link>
          </Button>
          <Button asChild variant="outline">
            <Link href="/plans">Plans</Link>
          </Button>
          <Button asChild variant="outline">
            <Link href="/service-requests">Service requests</Link>
          </Button>
          <Button
            type="button"
            variant="secondary"
            onClick={() => {
              clearAccessToken()
              void router.push("/login")
            }}
          >
            Logout
          </Button>
        </nav>
      </header>

      {children}
    </main>
  )
}
