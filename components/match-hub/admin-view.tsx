"use client"

import { useState, useEffect, useCallback } from "react"
import { Shield, SnowflakeIcon, Star, Coins, ChevronDown, ChevronUp, RefreshCw, ArrowLeft, Trash2, ClipboardList, Users } from "lucide-react"
import { apiClient, APIError } from "@/lib/api-client"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"

interface AdminUser {
  id: string
  email: string
  name: string
  is_admin: boolean
  is_frozen: boolean
  vip_level: number
  credits: number
}

interface AuditEntry {
  id: string
  admin_id: string
  target_id?: string
  action: string
  details?: string
  created_at: string
}

const VIP_LABELS = ["Sin VIP", "Bronce", "Plata", "Oro", "Platino", "Diamante"]
const VIP_COLORS = [
  "bg-gray-100 text-gray-700",
  "bg-amber-100 text-amber-700",
  "bg-slate-200 text-slate-700",
  "bg-yellow-100 text-yellow-700",
  "bg-sky-100 text-sky-700",
  "bg-cyan-100 text-cyan-700",
]

function Toast({ message, type, onDismiss }: { message: string; type: "success" | "error"; onDismiss: () => void }) {
  useEffect(() => {
    const t = setTimeout(onDismiss, 3000)
    return () => clearTimeout(t)
  }, [onDismiss])
  return (
    <div className={`fixed top-4 left-1/2 -translate-x-1/2 z-50 px-4 py-2 rounded-xl shadow-lg text-sm font-medium text-white max-w-[360px] text-center ${type === "success" ? "bg-green-500" : "bg-red-500"}`}>
      {message}
    </div>
  )
}

function ConfirmModal({ message, onConfirm, onCancel }: { message: string; onConfirm: () => void; onCancel: () => void }) {
  return (
    <div className="fixed inset-0 z-40 bg-black/50 flex items-center justify-center p-4">
      <div className="bg-card rounded-2xl p-6 max-w-sm w-full shadow-xl">
        <p className="text-sm text-card-foreground mb-6">{message}</p>
        <div className="flex gap-3">
          <Button size="sm" variant="outline" className="flex-1" onClick={onCancel}>Cancelar</Button>
          <Button size="sm" variant="destructive" className="flex-1" onClick={onConfirm}>Confirmar</Button>
        </div>
      </div>
    </div>
  )
}

function UserRow({
  adminUser,
  onRefresh,
  onToast,
}: {
  adminUser: AdminUser
  onRefresh: () => void
  onToast: (msg: string, type: "success" | "error") => void
}) {
  const [expanded, setExpanded] = useState(false)
  const [loading, setLoading] = useState(false)
  const [creditDelta, setCreditDelta] = useState(0)
  const [vipLevel, setVipLevel] = useState(adminUser.vip_level)
  const [confirm, setConfirm] = useState<{ msg: string; fn: () => void } | null>(null)

  const act = async (action: string, body: object, successMsg?: string) => {
    setLoading(true)
    try {
      await apiClient(`/api/v1/admin/${action}`, { method: "POST", body })
      onToast(successMsg ?? "Acción aplicada", "success")
      onRefresh()
    } catch (e) {
      const msg = e instanceof APIError ? e.message : "Error desconocido"
      onToast(msg, "error")
    } finally {
      setLoading(false)
    }
  }

  const deleteUser = async () => {
    setLoading(true)
    try {
      await apiClient(`/api/v1/admin/users/${adminUser.id}`, { method: "DELETE" })
      onToast("Usuario eliminado", "success")
      onRefresh()
    } catch (e) {
      const msg = e instanceof APIError ? e.message : "Error al eliminar"
      onToast(msg, "error")
    } finally {
      setLoading(false)
    }
  }

  const withConfirm = (msg: string, fn: () => void) => setConfirm({ msg, fn })

  return (
    <>
      {confirm && (
        <ConfirmModal
          message={confirm.msg}
          onConfirm={() => { setConfirm(null); confirm.fn() }}
          onCancel={() => setConfirm(null)}
        />
      )}
      <div className="bg-card border border-border rounded-xl overflow-hidden mb-2">
        <div className="flex items-center gap-3 p-3 cursor-pointer" onClick={() => setExpanded((v) => !v)}>
          <div className="w-10 h-10 rounded-full bg-secondary flex items-center justify-center flex-shrink-0">
            <span className="text-sm font-semibold text-secondary-foreground">
              {adminUser.name?.[0]?.toUpperCase() ?? "?"}
            </span>
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-1.5 flex-wrap">
              <p className="text-sm font-medium text-card-foreground truncate">{adminUser.name}</p>
              {adminUser.is_admin && (
                <Badge variant="destructive" className="text-[10px] px-1.5 py-0">Admin</Badge>
              )}
              {adminUser.vip_level > 0 && (
                <span className={`text-[10px] font-medium px-1.5 py-0.5 rounded-full ${VIP_COLORS[adminUser.vip_level]}`}>
                  {VIP_LABELS[adminUser.vip_level]}
                </span>
              )}
              {adminUser.is_frozen && (
                <span className="text-[10px] font-medium px-1.5 py-0.5 rounded-full bg-blue-100 text-blue-700">Congelado</span>
              )}
            </div>
            <p className="text-xs text-muted-foreground truncate">{adminUser.email}</p>
          </div>
          <div className="flex items-center gap-2 flex-shrink-0">
            <span className="text-xs text-muted-foreground">{adminUser.credits} cr</span>
            {expanded ? <ChevronUp className="w-4 h-4 text-muted-foreground" /> : <ChevronDown className="w-4 h-4 text-muted-foreground" />}
          </div>
        </div>

        {expanded && (
          <div className="px-3 pb-3 border-t border-border/50 pt-3 space-y-3">
            {/* Freeze/Unfreeze + Role */}
            <div className="flex gap-2">
              <Button
                size="sm"
                variant={adminUser.is_frozen ? "outline" : "secondary"}
                className="flex-1 text-xs"
                disabled={loading}
                onClick={() => act(adminUser.is_frozen ? "unfreeze" : "freeze", { user_id: adminUser.id },
                  adminUser.is_frozen ? "Usuario descongelado" : "Usuario congelado")}
              >
                <SnowflakeIcon className="w-3.5 h-3.5 mr-1" />
                {adminUser.is_frozen ? "Descongelar" : "Congelar"}
              </Button>
              <Button
                size="sm"
                variant={adminUser.is_admin ? "destructive" : "secondary"}
                className="flex-1 text-xs"
                disabled={loading}
                onClick={() =>
                  withConfirm(
                    adminUser.is_admin ? "¿Quitar rol admin?" : "¿Dar rol admin?",
                    () => act("role", { user_id: adminUser.id, is_admin: !adminUser.is_admin },
                      adminUser.is_admin ? "Admin revocado" : "Admin concedido")
                  )
                }
              >
                <Shield className="w-3.5 h-3.5 mr-1" />
                {adminUser.is_admin ? "Quitar admin" : "Dar admin"}
              </Button>
            </div>

            {/* VIP level */}
            <div className="flex items-center gap-2">
              <Star className="w-4 h-4 text-yellow-500 flex-shrink-0" />
              <select
                value={vipLevel}
                onChange={(e) => setVipLevel(Number(e.target.value))}
                className="flex-1 text-xs border border-border rounded-lg px-2 py-1.5 bg-background focus:outline-none focus:ring-1 focus:ring-primary"
              >
                {VIP_LABELS.map((label, i) => (
                  <option key={i} value={i}>{label}</option>
                ))}
              </select>
              <Button
                size="sm"
                className="text-xs bg-primary hover:bg-primary/90 px-3"
                disabled={loading || vipLevel === adminUser.vip_level}
                onClick={() => act("vip", { user_id: adminUser.id, vip_level: vipLevel }, `VIP → ${VIP_LABELS[vipLevel]}`)}
              >
                Guardar
              </Button>
            </div>

            {/* Credits */}
            <div className="flex items-center gap-2">
              <Coins className="w-4 h-4 text-yellow-500 flex-shrink-0" />
              <input
                type="number"
                value={creditDelta}
                onChange={(e) => setCreditDelta(Number(e.target.value))}
                placeholder="±créditos"
                className="flex-1 text-xs border border-border rounded-lg px-2 py-1.5 bg-background focus:outline-none focus:ring-1 focus:ring-primary"
              />
              <Button
                size="sm"
                className="text-xs bg-primary hover:bg-primary/90 px-3"
                disabled={loading || creditDelta === 0}
                onClick={() => act("credits", { user_id: adminUser.id, delta: creditDelta }, `Créditos ajustados: ${creditDelta > 0 ? "+" : ""}${creditDelta}`)}
              >
                Aplicar
              </Button>
            </div>

            {/* Delete */}
            <Button
              size="sm"
              variant="destructive"
              className="w-full text-xs"
              disabled={loading}
              onClick={() => withConfirm(`¿Eliminar cuenta de ${adminUser.name}? Esta acción no se puede deshacer.`, deleteUser)}
            >
              <Trash2 className="w-3.5 h-3.5 mr-1" />
              Eliminar cuenta
            </Button>
          </div>
        )}
      </div>
    </>
  )
}

const ACTION_LABELS: Record<string, string> = {
  freeze_user: "Congelar",
  unfreeze_user: "Descongelar",
  set_vip: "VIP",
  adjust_credits: "Créditos",
  grant_admin: "Dar admin",
  revoke_admin: "Quitar admin",
  delete_user: "Eliminar",
}

function AuditTab() {
  const [logs, setLogs] = useState<AuditEntry[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setLoading(true)
    apiClient<AuditEntry[]>("/api/v1/admin/audit-log")
      .then((data) => setLogs(data ?? []))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const formatDate = (s: string) =>
    new Date(s).toLocaleString("es-ES", { day: "2-digit", month: "2-digit", hour: "2-digit", minute: "2-digit" })

  return (
    <div className="flex-1 overflow-y-auto px-4 pt-3 pb-24">
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
        </div>
      ) : logs.length === 0 ? (
        <p className="text-center text-muted-foreground text-sm py-8">Sin acciones registradas</p>
      ) : (
        logs.map((log) => (
          <div key={log.id} className="flex items-start gap-3 py-2 border-b border-border/40">
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <span className="text-xs font-semibold text-card-foreground">{ACTION_LABELS[log.action] ?? log.action}</span>
                {log.details && (
                  <span className="text-[10px] text-muted-foreground font-mono">{log.details}</span>
                )}
              </div>
              <p className="text-[10px] text-muted-foreground truncate">
                {log.target_id ? `→ ${log.target_id.slice(0, 8)}...` : ""}
              </p>
            </div>
            <span className="text-[10px] text-muted-foreground flex-shrink-0">{formatDate(log.created_at)}</span>
          </div>
        ))
      )}
    </div>
  )
}

type AdminTab = "users" | "audit"

export function AdminView({ onClose }: { onClose?: () => void }) {
  const [users, setUsers] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState("")
  const [filter, setFilter] = useState<"all" | "frozen" | "admin">("all")
  const [activeTab, setActiveTab] = useState<AdminTab>("users")
  const [toast, setToast] = useState<{ msg: string; type: "success" | "error" } | null>(null)

  const showToast = useCallback((msg: string, type: "success" | "error") => {
    setToast({ msg, type })
  }, [])

  const loadUsers = useCallback(async () => {
    setLoading(true)
    try {
      const data = await apiClient<AdminUser[]>("/api/v1/admin/users")
      setUsers(data ?? [])
    } catch (e) {
      const msg = e instanceof APIError ? e.message : "Error cargando usuarios"
      showToast(msg, "error")
    } finally {
      setLoading(false)
    }
  }, [showToast])

  useEffect(() => { loadUsers() }, [loadUsers])

  const filtered = users.filter((u) => {
    const matchSearch =
      u.name.toLowerCase().includes(search.toLowerCase()) ||
      u.email.toLowerCase().includes(search.toLowerCase())
    const matchFilter =
      filter === "all" ||
      (filter === "frozen" && u.is_frozen) ||
      (filter === "admin" && u.is_admin)
    return matchSearch && matchFilter
  })

  const adminCount = users.filter((u) => u.is_admin).length
  const frozenCount = users.filter((u) => u.is_frozen).length

  return (
    <div className="flex flex-col h-full">
      {toast && (
        <Toast message={toast.msg} type={toast.type} onDismiss={() => setToast(null)} />
      )}

      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-border bg-card sticky top-0 z-10">
        <div className="flex items-center gap-2">
          {onClose && (
            <button onClick={onClose} className="mr-1 text-muted-foreground hover:text-card-foreground">
              <ArrowLeft className="w-5 h-5" />
            </button>
          )}
          <Shield className="w-5 h-5 text-primary" />
          <h1 className="text-lg font-bold text-card-foreground">Panel Admin</h1>
        </div>
        <button onClick={loadUsers} className="text-muted-foreground hover:text-card-foreground">
          <RefreshCw className={`w-4 h-4 ${loading ? "animate-spin" : ""}`} />
        </button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-3 px-4 py-3 border-b border-border/50 bg-card">
        <div className="text-center">
          <p className="text-2xl font-bold text-primary">{users.length}</p>
          <p className="text-xs text-muted-foreground">Usuarios</p>
        </div>
        <div className="text-center">
          <p className="text-2xl font-bold text-blue-500">{frozenCount}</p>
          <p className="text-xs text-muted-foreground">Congelados</p>
        </div>
        <div className="text-center">
          <p className="text-2xl font-bold text-yellow-500">{adminCount}</p>
          <p className="text-xs text-muted-foreground">Admins</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-border bg-card">
        <button
          onClick={() => setActiveTab("users")}
          className={`flex-1 flex items-center justify-center gap-1.5 py-2.5 text-xs font-medium transition-colors ${activeTab === "users" ? "text-primary border-b-2 border-primary" : "text-muted-foreground"}`}
        >
          <Users className="w-3.5 h-3.5" />
          Usuarios
        </button>
        <button
          onClick={() => setActiveTab("audit")}
          className={`flex-1 flex items-center justify-center gap-1.5 py-2.5 text-xs font-medium transition-colors ${activeTab === "audit" ? "text-primary border-b-2 border-primary" : "text-muted-foreground"}`}
        >
          <ClipboardList className="w-3.5 h-3.5" />
          Auditoría
        </button>
      </div>

      {activeTab === "audit" ? (
        <AuditTab />
      ) : (
        <>
          {/* Search + filter */}
          <div className="px-4 py-2 bg-card border-b border-border/50 space-y-2">
            <input
              placeholder="Buscar usuario..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
            />
            <div className="flex gap-2">
              {(["all", "frozen", "admin"] as const).map((f) => (
                <button
                  key={f}
                  onClick={() => setFilter(f)}
                  className={`text-[11px] px-2.5 py-1 rounded-full border transition-colors ${filter === f ? "bg-primary text-primary-foreground border-primary" : "border-border text-muted-foreground"}`}
                >
                  {f === "all" ? "Todos" : f === "frozen" ? "Congelados" : "Admins"}
                </button>
              ))}
              <span className="ml-auto text-[11px] text-muted-foreground self-center">
                {filtered.length} resultado{filtered.length !== 1 ? "s" : ""}
              </span>
            </div>
          </div>

          {/* User list */}
          <div className="flex-1 overflow-y-auto px-4 pt-3 pb-24">
            {loading ? (
              <div className="flex items-center justify-center py-12">
                <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
              </div>
            ) : filtered.length === 0 ? (
              <p className="text-center text-muted-foreground text-sm py-8">
                {search || filter !== "all" ? "No se encontraron usuarios" : "No hay usuarios"}
              </p>
            ) : (
              filtered.map((u) => (
                <UserRow key={u.id} adminUser={u} onRefresh={loadUsers} onToast={showToast} />
              ))
            )}
          </div>
        </>
      )}
    </div>
  )
}
