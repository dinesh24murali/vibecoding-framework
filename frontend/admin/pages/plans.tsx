import Head from "next/head"
import { useEffect, useMemo, useState } from "react"
import {
  createPlan,
  deletePlan,
  listPlans,
  type Plan,
  updatePlan,
  uploadPlansCsv
} from "../lib/api"
import { isLoggedIn } from "../lib/auth"
import { useRouter } from "next/router"
import { AdminPageLayout } from "@/components/admin-page-layout"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Badge } from "@/components/ui/badge"
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

export default function PlansPage() {
  const router = useRouter()
  const [plans, setPlans] = useState<Plan[]>([])
  const [form, setForm] = useState({
    provider_id: "",
    name: "",
    description: "",
    price: "",
    discount: "",
    is_active: true
  })
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState("")
  const [deleteTarget, setDeleteTarget] = useState<Plan | null>(null)

  const sortedPlans = useMemo(() => [...plans].sort((a, b) => a.name.localeCompare(b.name)), [plans])

  const refresh = async () => {
    setPlans(await listPlans())
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
      await createPlan({
        provider_id: Number(form.provider_id),
        name: form.name.trim(),
        description: form.description.trim(),
        price: Number(form.price),
        discount: Number(form.discount),
        is_active: form.is_active
      })
      setForm({ provider_id: "", name: "", description: "", price: "", discount: "", is_active: true })
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create plan")
    } finally {
      setBusy(false)
    }
  }

  const onToggleActive = async (plan: Plan) => {
    setBusy(true)
    setError("")
    try {
      await updatePlan(plan.id, { is_active: !plan.is_active })
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update plan")
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
      await deletePlan(deleteTarget.id)
      setDeleteTarget(null)
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete plan")
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
      await uploadPlansCsv(file)
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
        <title>Village DTH Admin - Plans</title>
      </Head>

      <AdminPageLayout title="Plans">
        <Card className="mx-auto mb-4 w-full max-w-6xl bg-white/90">
          <CardHeader>
            <CardTitle>Add plan</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <form className="grid grid-cols-1 gap-2 md:grid-cols-4" onSubmit={onCreate}>
              <Input
                placeholder="Provider ID"
                value={form.provider_id}
                onChange={(e) => setForm((prev) => ({ ...prev, provider_id: e.target.value }))}
                required
              />
              <Input
                placeholder="Plan name"
                value={form.name}
                onChange={(e) => setForm((prev) => ({ ...prev, name: e.target.value }))}
                required
              />
              <Input
                placeholder="Description"
                value={form.description}
                onChange={(e) => setForm((prev) => ({ ...prev, description: e.target.value }))}
                required
              />
              <Input
                placeholder="Price"
                value={form.price}
                onChange={(e) => setForm((prev) => ({ ...prev, price: e.target.value }))}
                required
              />
              <Input
                placeholder="Discount"
                value={form.discount}
                onChange={(e) => setForm((prev) => ({ ...prev, discount: e.target.value }))}
                required
              />
              <label className="flex items-center gap-2 rounded-md border bg-background px-3 py-2 text-sm font-medium">
                <input
                  type="checkbox"
                  checked={form.is_active}
                  onChange={(e) => setForm((prev) => ({ ...prev, is_active: e.target.checked }))}
                />
                Active
              </label>
              <Button type="submit" disabled={busy}>
                Save
              </Button>
            </form>

            <div className="space-y-2">
              <Label htmlFor="plans-csv">Upload CSV</Label>
              <Input id="plans-csv" type="file" accept=".csv,text/csv" onChange={onUpload} disabled={busy} />
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
                  <TableHead>Plan</TableHead>
                  <TableHead>Provider</TableHead>
                  <TableHead>Net Price</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-56">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {sortedPlans.map((plan) => (
                  <TableRow key={plan.id}>
                    <TableCell className="font-medium">{plan.name}</TableCell>
                    <TableCell>{plan.provider_id}</TableCell>
                    <TableCell>Rs {(plan.price - plan.discount).toFixed(2)}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className={plan.is_active ? "bg-emerald-100 text-emerald-700" : "bg-zinc-100 text-zinc-700"}>
                        {plan.is_active ? "Active" : "Inactive"}
                      </Badge>
                    </TableCell>
                    <TableCell className="space-x-2">
                      <Button type="button" variant="outline" onClick={() => onToggleActive(plan)} disabled={busy}>
                        Toggle active
                      </Button>
                      <Button type="button" variant="destructive" onClick={() => setDeleteTarget(plan)} disabled={busy}>
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

        <Dialog open={Boolean(deleteTarget)} onOpenChange={(open) => !open && setDeleteTarget(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Delete plan</DialogTitle>
              <DialogDescription>
                This removes <strong>{deleteTarget?.name}</strong> permanently.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setDeleteTarget(null)}>
                Cancel
              </Button>
              <Button type="button" variant="destructive" onClick={onDeleteConfirm} disabled={busy}>
                Delete plan
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
