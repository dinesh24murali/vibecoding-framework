import Head from "next/head"
import { useEffect, useMemo, useState } from "react"
import {
  createProvider,
  deleteProvider,
  listProviders,
  updateProvider,
  uploadProvidersCsv,
  type Provider
} from "../lib/api"
import { isLoggedIn } from "../lib/auth"
import { useRouter } from "next/router"
import { AdminPageLayout } from "@/components/admin-page-layout"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

export default function ProvidersPage() {
  const router = useRouter()
  const [providers, setProviders] = useState<Provider[]>([])
  const [name, setName] = useState("")
  const [imageURL, setImageURL] = useState("")
  const [error, setError] = useState("")
  const [busy, setBusy] = useState(false)

  const [renameTarget, setRenameTarget] = useState<Provider | null>(null)
  const [renameValue, setRenameValue] = useState("")
  const [deleteTarget, setDeleteTarget] = useState<Provider | null>(null)

  const sortedProviders = useMemo(() => [...providers].sort((a, b) => a.name.localeCompare(b.name)), [providers])

  const refresh = async () => {
    const data = await listProviders()
    setProviders(data)
  }

  useEffect(() => {
    if (!isLoggedIn()) {
      void router.replace("/login")
      return
    }
    refresh().catch((err: Error) => setError(err.message))
  }, [router])

  const onCreate = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setBusy(true)
    setError("")
    try {
      await createProvider({ name: name.trim(), image_url: imageURL.trim() || null })
      setName("")
      setImageURL("")
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create provider")
    } finally {
      setBusy(false)
    }
  }

  const onRenameOpen = (provider: Provider) => {
    setRenameTarget(provider)
    setRenameValue(provider.name)
  }

  const onRenameSubmit = async () => {
    if (!renameTarget) {
      return
    }
    const next = renameValue.trim()
    if (!next || next === renameTarget.name) {
      setRenameTarget(null)
      return
    }

    setBusy(true)
    setError("")
    try {
      await updateProvider(renameTarget.id, { name: next })
      setRenameTarget(null)
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update provider")
    } finally {
      setBusy(false)
    }
  }

  const onDeleteConfirm = async () => {
    if (!deleteTarget) {
      return
    }
    setBusy(true)
    setError("")
    try {
      await deleteProvider(deleteTarget.id)
      setDeleteTarget(null)
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete provider")
    } finally {
      setBusy(false)
    }
  }

  const onUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) {
      return
    }
    setBusy(true)
    setError("")
    try {
      await uploadProvidersCsv(file)
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to upload CSV")
    } finally {
      setBusy(false)
      event.target.value = ""
    }
  }

  return (
    <>
      <Head>
        <title>Village DTH Admin - Providers</title>
      </Head>

      <AdminPageLayout title="Providers">
        <Card className="mx-auto mb-4 w-full max-w-6xl bg-white/90">
          <CardHeader>
            <CardTitle>Add provider</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <form onSubmit={onCreate} className="grid grid-cols-1 gap-2 sm:grid-cols-[1fr_1fr_auto]">
              <Input placeholder="Provider name" value={name} onChange={(e) => setName(e.target.value)} required />
              <Input placeholder="Image URL (optional)" value={imageURL} onChange={(e) => setImageURL(e.target.value)} />
              <Button type="submit" disabled={busy}>
                Save
              </Button>
            </form>

            <div className="space-y-2">
              <Label htmlFor="providers-csv">Upload CSV</Label>
              <Input id="providers-csv" type="file" accept=".csv,text/csv" onChange={onUpload} disabled={busy} />
            </div>
          </CardContent>
        </Card>

        <Card className="mx-auto w-full max-w-6xl bg-white/90">
          <CardHeader>
            <CardTitle>List</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>ID</TableHead>
                  <TableHead className="w-48">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {sortedProviders.map((provider) => (
                  <TableRow key={provider.id}>
                    <TableCell className="font-medium">{provider.name}</TableCell>
                    <TableCell className="font-mono text-xs text-muted-foreground">{provider.id}</TableCell>
                    <TableCell className="space-x-2">
                      <Button type="button" variant="outline" onClick={() => onRenameOpen(provider)} disabled={busy}>
                        Rename
                      </Button>
                      <Button type="button" variant="destructive" onClick={() => setDeleteTarget(provider)} disabled={busy}>
                        Delete
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            {error && (
              <Alert variant="destructive">
                <AlertTitle>Action failed</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}
          </CardContent>
        </Card>

        <Dialog open={Boolean(renameTarget)} onOpenChange={(open) => !open && setRenameTarget(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Rename provider</DialogTitle>
              <DialogDescription>Update provider name and save changes.</DialogDescription>
            </DialogHeader>
            <div className="space-y-2">
              <Label htmlFor="rename-provider">Provider name</Label>
              <Input id="rename-provider" value={renameValue} onChange={(e) => setRenameValue(e.target.value)} autoFocus />
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setRenameTarget(null)}>
                Cancel
              </Button>
              <Button type="button" onClick={onRenameSubmit} disabled={busy}>
                Save changes
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <Dialog open={Boolean(deleteTarget)} onOpenChange={(open) => !open && setDeleteTarget(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Delete provider</DialogTitle>
              <DialogDescription>
                This removes <strong>{deleteTarget?.name}</strong> permanently.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setDeleteTarget(null)}>
                Cancel
              </Button>
              <Button type="button" variant="destructive" onClick={onDeleteConfirm} disabled={busy}>
                Delete provider
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </AdminPageLayout>
    </>
  )
}

export async function getStaticProps() {
  return { props: {} }
}
