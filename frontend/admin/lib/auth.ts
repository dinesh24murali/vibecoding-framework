const TOKEN_KEY = "dth_admin_access_token"

export function getAccessToken(): string {
  if (typeof window === "undefined") {
    return ""
  }
  return window.localStorage.getItem(TOKEN_KEY) ?? ""
}

export function setAccessToken(token: string): void {
  if (typeof window === "undefined") {
    return
  }
  window.localStorage.setItem(TOKEN_KEY, token)
}

export function clearAccessToken(): void {
  if (typeof window === "undefined") {
    return
  }
  window.localStorage.removeItem(TOKEN_KEY)
}

export function isLoggedIn(): boolean {
  return getAccessToken().length > 0
}
