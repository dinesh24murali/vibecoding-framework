declare global {
  interface Window {
    grecaptcha?: {
      ready: (fn: () => void) => void
      execute: (siteKey: string, options: { action: string }) => Promise<string>
    }
  }
}

const SITE_KEY = process.env.NEXT_PUBLIC_RECAPTCHA_SITE_KEY ?? ""

export async function getRecaptchaToken(action: string): Promise<string> {
  if (typeof window === "undefined") {
    return `server-${action}-${Date.now()}`
  }

  if (!SITE_KEY || !window.grecaptcha) {
    return `local-${action}-${Date.now()}`
  }

  return new Promise<string>((resolve, reject) => {
    window.grecaptcha?.ready(() => {
      window.grecaptcha
        ?.execute(SITE_KEY, { action })
        .then(resolve)
        .catch(reject)
    })
  })
}

export function loadRecaptchaScript(): void {
  if (typeof window === "undefined" || !SITE_KEY) {
    return
  }

  if (document.querySelector("script[data-recaptcha='customer']")) {
    return
  }

  const script = document.createElement("script")
  script.src = `https://www.google.com/recaptcha/api.js?render=${encodeURIComponent(SITE_KEY)}`
  script.async = true
  script.defer = true
  script.dataset.recaptcha = "customer"
  document.head.appendChild(script)
}
