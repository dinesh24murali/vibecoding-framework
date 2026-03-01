import { useEffect } from "react"
import { useRouter } from "next/router"
import { isLoggedIn } from "../lib/auth"

export default function HomePage() {
  const router = useRouter()

  useEffect(() => {
    if (isLoggedIn()) {
      void router.replace("/providers")
      return
    }
    void router.replace("/login")
  }, [router])

  return null
}

export async function getStaticProps() {
  return { props: {} }
}
