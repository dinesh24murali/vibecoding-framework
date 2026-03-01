import Head from "next/head"
import { useEffect, useState } from "react"
import {
  getServiceRequestById,
  listServiceRequests,
  type ServiceRequest,
  type ServiceRequestDetails,
  updateServiceRequestStatus
} from "../lib/api"
import { isLoggedIn } from "../lib/auth"
import { useRouter } from "next/router"
import { AdminPageLayout } from "@/components/admin-page-layout"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

const allowedStatuses: ServiceRequest["status"][] = ["pending", "payment_failed", "payment_success", "blocked", "completed"]

function statusTone(status: ServiceRequest["status"]) {
  switch (status) {
    case "pending":
      return "bg-amber-100 text-amber-700"
    case "payment_failed":
      return "bg-red-100 text-red-700"
    case "payment_success":
    case "completed":
      return "bg-emerald-100 text-emerald-700"
    case "blocked":
      return "bg-rose-100 text-rose-700"
    default:
      return ""
  }
}

export default function ServiceRequestsPage() {
  const router = useRouter()
  const [items, setItems] = useState<ServiceRequest[]>([])
  const [selected, setSelected] = useState<ServiceRequestDetails | null>(null)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState("")

  const refresh = async () => {
    setItems(await listServiceRequests())
  }

  useEffect(() => {
    if (!isLoggedIn()) {
      void router.replace("/login")
      return
    }
    refresh().catch((err: Error) => setError(err.message))
  }, [router])

  const selectById = async (id: string) => {
    setBusy(true)
    setError("")
    try {
      setSelected(await getServiceRequestById(id))
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch details")
    } finally {
      setBusy(false)
    }
  }

  const onStatusChange = async (id: string, nextStatus: ServiceRequest["status"]) => {
    setBusy(true)
    setError("")
    try {
      await updateServiceRequestStatus(id, nextStatus)
      await refresh()
      await selectById(id)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update status")
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      <Head>
        <title>Village DTH Admin - Service Requests</title>
      </Head>

      <AdminPageLayout title="Service Requests">
        <section className="mx-auto grid w-full max-w-6xl grid-cols-1 gap-4 lg:grid-cols-2">
          <Card className="bg-white/90">
            <CardHeader>
              <CardTitle>Requests</CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>ID</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="w-28">View</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {items.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell className="font-mono text-xs">{item.id}</TableCell>
                      <TableCell>
                        <Badge className={statusTone(item.status)} variant="outline">
                          {item.status}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <Button type="button" variant="outline" onClick={() => selectById(item.id)} disabled={busy}>
                          Open
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          <Card className="bg-white/90">
            <CardHeader>
              <CardTitle>Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3 text-sm">
              {!selected && <p className="text-muted-foreground">Select a request to view details.</p>}
              {selected && (
                <>
                  <p>
                    <strong>ID:</strong> {selected.id}
                  </p>
                  <p className="flex items-center gap-2">
                    <strong>Status:</strong>
                    <Badge className={statusTone(selected.status)} variant="outline">
                      {selected.status}
                    </Badge>
                  </p>
                  <p>
                    <strong>Customer:</strong> {selected.user.name} ({selected.user.phone_number})
                  </p>
                  <p>
                    <strong>Plan:</strong> {selected.plan.name}
                  </p>

                  <div className="space-y-2">
                    <p className="text-sm font-medium">Update status</p>
                    <Select
                      value={selected.status}
                      onValueChange={(value) => onStatusChange(selected.id, value as ServiceRequest["status"])}
                      disabled={busy}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select status" />
                      </SelectTrigger>
                      <SelectContent>
                        {allowedStatuses.map((status) => (
                          <SelectItem key={status} value={status}>
                            {status}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </section>

        {error && (
          <div className="mx-auto mt-4 w-full max-w-6xl">
            <Alert variant="destructive">
              <AlertTitle>Action failed</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          </div>
        )}
      </AdminPageLayout>
    </>
  )
}

export async function getStaticProps() {
  return { props: {} }
}
