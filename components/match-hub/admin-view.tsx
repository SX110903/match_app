"use client"

import { useState, useEffect, useCallback } from "react"
import { Shield, User, SnowflakeIcon, Star, Coins, ChevronDown, ChevronUp, RefreshCw, ArrowLeft } from "lucide-react"
import { apiClient } from "@/lib/api-client"
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

const VIP_LABELS = ["Sin VIP", "Bronce", "Plata", "Oro", "Platino", "Diamante"]
const VIP_COLORS = [
  "bg-gray-100 text-gray-700",
  "bg-amber-100 text-amber-700",
  "bg-slate-200 text-slate-700",
  "bg-yellow-100 text-yellow-700",
  "bg-sky-100 text-sky-700",
  "bg-cyan-100 text-cyan-700",
]

function UserRow({ adminUser, onRefresh }: { adminUser: AdminUser; onRefresh: () => void }) {
  const [expanded, setExpanded] = useState(false)
  const [loading, setLoading] = useState(false)
  const [creditDelta, setCreditDelta] = useState(0)
  const [vipLevel, setVipLevel] = useState(adminUser.vip_level)

  const act = async (action: string, body: object) => {
    setLoading(true)
    try {
      await apiClient(`/api/v1/admin/${action}`, { method: 'POST', body })
      onRefresh()
    } catch {
      //
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bg-card border border-border rounded-xl overflow-hidden mb-2">
      <div
        className="flex items-center gap-3 p-3 cursor-pointer"
        onClick={() => setExpanded((v) => !v)}
      >
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
          {/* Freeze/Unfreeze */}
          <div className="flex gap-2">
            <Button
              size="sm"
              variant={adminUser.is_frozen ? "outline" : "secondary"}
              className="flex-1 text-xs"
              disabled={loading}
              onClick={() => act(adminUser.is_frozen ? "unfreeze" : "freeze", { user_id: adminUser.id })}
            >
              <SnowflakeIcon className="w-3.5 h-3.5 mr-1" />
              {adminUser.is_frozen ? "Descongelar" : "Congelar"}
            </Button>
            <Button
              size="sm"
              variant={adminUser.is_admin ? "destructive" : "secondary"}
              className="flex-1 text-xs"
              disabled={loading}
              onClick={() => act("role", { user_id: adminUser.id, is_admin: !adminUser.is_admin })}
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
              onClick={() => act("vip", { user_id: adminUser.id, vip_level: vipLevel })}
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
              onClick={() => act("credits", { user_id: adminUser.id, delta: creditDelta })}
            >
              Aplicar
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}

export function AdminView({ onClose }: { onClose?: () => void }) {
  const [users, setUsers] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState("")

  const loadUsers = useCallback(async () => {
    setLoading(true)
    try {
      const data = await apiClient<AdminUser[]>('/api/v1/admin/users')
      setUsers(data ?? [])
    } catch {
      //
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { loadUsers() }, [loadUsers])

  const filtered = users.filter((u) =>
    u.name.toLowerCase().includes(search.toLowerCase()) ||
    u.email.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="flex flex-col h-full">
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
          <p className="text-2xl font-bold text-blue-500">{users.filter((u) => u.is_frozen).length}</p>
          <p className="text-xs text-muted-foreground">Congelados</p>
        </div>
        <div className="text-center">
          <p className="text-2xl font-bold text-yellow-500">{users.filter((u) => u.vip_level > 0).length}</p>
          <p className="text-xs text-muted-foreground">VIP</p>
        </div>
      </div>

      {/* Search */}
      <div className="px-4 py-2 bg-card border-b border-border/50">
        <input
          placeholder="Buscar usuario..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
        />
      </div>

      {/* User list */}
      <div className="flex-1 overflow-y-auto px-4 pt-3 pb-24">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
          </div>
        ) : filtered.length === 0 ? (
          <p className="text-center text-muted-foreground text-sm py-8">
            {search ? "No se encontraron usuarios" : "No hay usuarios"}
          </p>
        ) : (
          filtered.map((u) => (
            <UserRow key={u.id} adminUser={u} onRefresh={loadUsers} />
          ))
        )}
      </div>
    </div>
  )
}
