"use client"

import { useState, useEffect, useCallback } from "react"
import { ArrowLeft, ShoppingBag, Star, Clock, CheckCircle2, Lock } from "lucide-react"
import { apiClient, APIError } from "@/lib/api-client"
import { useAuth } from "@/lib/auth-context"
import { Button } from "@/components/ui/button"
import { VERIFY_BADGE_COST } from "@/lib/constants"

interface ShopItem {
  item_type: string
  item_value: number
  cost: number
  name: string
  benefits: string
}

interface ShopTransaction {
  id: string
  item_type: string
  item_value: number
  cost: number
  created_at: string
}

const VIP_LABELS = ["", "Bronce", "Plata", "Oro", "Platino", "Diamante"]
const VIP_COLORS = [
  "",
  "from-amber-500 to-amber-700",
  "from-slate-400 to-slate-600",
  "from-yellow-400 to-yellow-600",
  "from-sky-400 to-sky-600",
  "from-cyan-300 to-cyan-600",
]

function Toast({ message, type, onDismiss }: { message: string; type: "success" | "error"; onDismiss: () => void }) {
  useEffect(() => {
    const t = setTimeout(onDismiss, 3500)
    return () => clearTimeout(t)
  }, [onDismiss])
  return (
    <div className={`fixed top-4 left-1/2 -translate-x-1/2 z-50 px-4 py-2 rounded-xl shadow-lg text-sm font-medium text-white max-w-[360px] text-center ${type === "success" ? "bg-green-500" : "bg-red-500"}`}>
      {message}
    </div>
  )
}

export function ShopView({ onClose }: { onClose?: () => void }) {
  const { user, refreshUser } = useAuth()
  const [items, setItems] = useState<ShopItem[]>([])
  const [transactions, setTransactions] = useState<ShopTransaction[]>([])
  const [activeTab, setActiveTab] = useState<"shop" | "history">("shop")
  const [loading, setLoading] = useState(true)
  const [purchasing, setPurchasing] = useState<number | null>(null)
  const [verifying, setVerifying] = useState(false)
  const [toast, setToast] = useState<{ msg: string; type: "success" | "error" } | null>(null)

  const showToast = useCallback((msg: string, type: "success" | "error") => {
    setToast({ msg, type })
  }, [])

  useEffect(() => {
    async function load() {
      setLoading(true)
      try {
        const [shopItems, txs] = await Promise.all([
          apiClient<ShopItem[]>("/api/v1/shop/items"),
          apiClient<ShopTransaction[]>("/api/v1/shop/transactions"),
        ])
        setItems(shopItems ?? [])
        setTransactions(txs ?? [])
      } catch { /* silent */ } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const handlePurchase = async (item: ShopItem) => {
    setPurchasing(item.item_value)
    try {
      await apiClient("/api/v1/shop/purchase", {
        method: "POST",
        body: { item_type: item.item_type, item_value: item.item_value },
      })
      showToast(`¡${item.name} activado!`, "success")
      await refreshUser()
      // Refresh transactions
      const txs = await apiClient<ShopTransaction[]>("/api/v1/shop/transactions")
      setTransactions(txs ?? [])
    } catch (e) {
      const msg = e instanceof APIError ? e.message : "Error en la compra"
      if (msg.includes("invalid purchase") || msg.includes("credits")) {
        showToast("Créditos insuficientes o nivel inválido", "error")
      } else {
        showToast(msg, "error")
      }
    } finally {
      setPurchasing(null)
    }
  }

  const handleVerify = async () => {
    setVerifying(true)
    try {
      await apiClient("/api/v1/users/me/verify", { method: "POST" })
      showToast("¡Badge Verificado activado!", "success")
      await refreshUser()
    } catch (e) {
      const msg = e instanceof APIError ? e.message : "Error"
      if (e instanceof APIError && e.status === 402) {
        showToast(`Necesitas ${VERIFY_BADGE_COST.toLocaleString()} créditos`, "error")
      } else if (e instanceof APIError && e.status === 409) {
        showToast("Ya tienes verificación gubernamental", "error")
      } else {
        showToast(msg, "error")
      }
    } finally {
      setVerifying(false)
    }
  }

  const currentVip = user?.vip_level ?? 0
  const credits = user?.credits ?? 0

  const getItemStatus = (item: ShopItem): "current" | "purchase" | "no_credits" | "level_skip" => {
    if (item.item_value === currentVip) return "current"
    if (item.item_value !== currentVip + 1) return "level_skip"
    if (credits < item.cost) return "no_credits"
    return "purchase"
  }

  const formatDate = (s: string) =>
    new Date(s).toLocaleString("es-ES", { day: "2-digit", month: "2-digit", year: "2-digit", hour: "2-digit", minute: "2-digit" })

  return (
    <div className="flex flex-col h-full bg-background">
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
          <ShoppingBag className="w-5 h-5 text-primary" />
          <h1 className="text-lg font-bold text-card-foreground">Tienda</h1>
        </div>
        <div className="flex items-center gap-1 bg-amber-100 text-amber-700 px-2 py-1 rounded-full">
          <span className="text-xs font-bold">{credits.toLocaleString()} cr</span>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-border bg-card">
        <button
          onClick={() => setActiveTab("shop")}
          className={`flex-1 py-2.5 text-xs font-medium transition-colors ${activeTab === "shop" ? "text-primary border-b-2 border-primary" : "text-muted-foreground"}`}
        >
          VIP & Extras
        </button>
        <button
          onClick={() => setActiveTab("history")}
          className={`flex-1 py-2.5 text-xs font-medium transition-colors flex items-center justify-center gap-1 ${activeTab === "history" ? "text-primary border-b-2 border-primary" : "text-muted-foreground"}`}
        >
          <Clock className="w-3.5 h-3.5" />
          Historial
        </button>
      </div>

      {loading ? (
        <div className="flex-1 flex items-center justify-center">
          <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
        </div>
      ) : activeTab === "shop" ? (
        <div className="flex-1 overflow-y-auto px-4 py-4 pb-24 space-y-3">
          {/* VIP items */}
          <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wide mb-2">Niveles VIP</p>
          {items.filter((i) => i.item_type === "vip_upgrade").map((item) => {
            const status = getItemStatus(item)
            return (
              <div
                key={item.item_value}
                className={`rounded-xl p-4 border ${status === "current" ? "border-primary bg-primary/5" : "border-border bg-card"}`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className={`w-10 h-10 rounded-full bg-gradient-to-br ${VIP_COLORS[item.item_value]} flex items-center justify-center`}>
                      <Star className="w-5 h-5 text-white" />
                    </div>
                    <div>
                      <p className="font-semibold text-card-foreground text-sm">{item.name}</p>
                      <p className="text-xs text-muted-foreground">{item.benefits}</p>
                    </div>
                  </div>
                  <div className="flex flex-col items-end gap-1 ml-2">
                    <span className="text-xs font-bold text-amber-600">{item.cost.toLocaleString()} cr</span>
                    {status === "current" && (
                      <span className="flex items-center gap-1 text-[11px] text-primary font-medium">
                        <CheckCircle2 className="w-3.5 h-3.5" /> Actual
                      </span>
                    )}
                    {status === "level_skip" && (
                      <span className="text-[11px] text-muted-foreground flex items-center gap-1">
                        <Lock className="w-3 h-3" /> Sube primero al nivel {item.item_value - 1}
                      </span>
                    )}
                    {(status === "purchase" || status === "no_credits") && (
                      <Button
                        size="sm"
                        className="text-xs h-7 px-3"
                        disabled={status === "no_credits" || purchasing === item.item_value}
                        onClick={() => handlePurchase(item)}
                      >
                        {purchasing === item.item_value ? (
                          <div className="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                        ) : status === "no_credits" ? (
                          "Sin créditos"
                        ) : (
                          "Comprar"
                        )}
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            )
          })}

          {/* Badge verificado */}
          <div className="border-t border-border/50 pt-3">
            <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wide mb-2">Badge Verificado</p>
            <div className="rounded-xl p-4 border border-border bg-card">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center">
                    <span className="text-white font-bold text-sm">✓</span>
                  </div>
                  <div>
                    <p className="font-semibold text-card-foreground text-sm">Badge Verificado</p>
                    <p className="text-xs text-muted-foreground">Marca de confianza en tu perfil</p>
                  </div>
                </div>
                <div className="flex flex-col items-end gap-1 ml-2">
                  <span className="text-xs font-bold text-amber-600">{VERIFY_BADGE_COST.toLocaleString()} cr</span>
                  {user?.badge === "verified" || user?.badge === "verified_gov" ? (
                    <span className="flex items-center gap-1 text-[11px] text-primary font-medium">
                      <CheckCircle2 className="w-3.5 h-3.5" /> Activo
                    </span>
                  ) : (
                    <Button
                      size="sm"
                      className="text-xs h-7 px-3"
                      disabled={verifying || credits < VERIFY_BADGE_COST}
                      onClick={handleVerify}
                    >
                      {verifying ? (
                        <div className="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      ) : credits < VERIFY_BADGE_COST ? (
                        "Sin créditos"
                      ) : (
                        "Verificar"
                      )}
                    </Button>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      ) : (
        <div className="flex-1 overflow-y-auto px-4 py-4 pb-24">
          {transactions.length === 0 ? (
            <p className="text-center text-muted-foreground text-sm py-8">Sin transacciones aún</p>
          ) : (
            transactions.map((tx) => (
              <div key={tx.id} className="flex items-center justify-between py-3 border-b border-border/40">
                <div>
                  <p className="text-sm font-medium text-card-foreground">
                    {tx.item_type === "vip_upgrade" ? `VIP ${VIP_LABELS[tx.item_value] ?? tx.item_value}` : tx.item_type}
                  </p>
                  <p className="text-xs text-muted-foreground">{formatDate(tx.created_at)}</p>
                </div>
                <span className="text-sm font-semibold text-destructive">-{tx.cost.toLocaleString()} cr</span>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  )
}
