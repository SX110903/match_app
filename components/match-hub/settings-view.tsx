"use client"

import { useState, useEffect } from "react"
import { X, Bell, Lock, Search, ChevronRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/lib/auth-context"
import { apiClient, APIError } from "@/lib/api-client"
import { AVATAR_BASE } from "@/lib/constants"

interface NotificationSettings {
  new_matches: boolean
  new_messages: boolean
  news_updates: boolean
  marketing: boolean
}

interface PrivacySettings {
  show_online_status: boolean
  show_last_seen: boolean
  show_distance: boolean
  incognito_mode: boolean
}

interface SearchPrefs {
  min_age: number
  max_age: number
  max_distance: number
}

function Toggle({ checked, onChange }: { checked: boolean; onChange: (v: boolean) => void }) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      onClick={() => onChange(!checked)}
      className={`relative w-11 h-6 rounded-full transition-colors ${checked ? "bg-primary" : "bg-muted"}`}
    >
      <span
        className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow transition-transform ${checked ? "translate-x-5" : "translate-x-0"}`}
      />
    </button>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="mb-6">
      <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider px-4 mb-2">{title}</h3>
      <div className="bg-card rounded-xl overflow-hidden border border-border">{children}</div>
    </div>
  )
}

function SettingRow({
  label,
  description,
  children,
  last,
}: {
  label: string
  description?: string
  children: React.ReactNode
  last?: boolean
}) {
  return (
    <div className={`flex items-center justify-between px-4 py-3 ${!last ? "border-b border-border" : ""}`}>
      <div className="flex-1 mr-4">
        <p className="text-sm font-medium text-card-foreground">{label}</p>
        {description && <p className="text-xs text-muted-foreground mt-0.5">{description}</p>}
      </div>
      {children}
    </div>
  )
}

export function SettingsView({ onClose }: { onClose: () => void }) {
  const { user } = useAuth()

  const [notifications, setNotifications] = useState<NotificationSettings>({
    new_matches: true,
    new_messages: true,
    news_updates: false,
    marketing: false,
  })
  const [privacy, setPrivacy] = useState<PrivacySettings>({
    show_online_status: true,
    show_last_seen: true,
    show_distance: true,
    incognito_mode: false,
  })
  const [searchPrefs, setSearchPrefs] = useState<SearchPrefs>({
    min_age: 18,
    max_age: 45,
    max_distance: 100,
  })
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState("")
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      try {
        const [notif, priv] = await Promise.all([
          apiClient<NotificationSettings>("/api/v1/settings/notifications"),
          apiClient<PrivacySettings>("/api/v1/settings/privacy"),
        ])
        if (notif) setNotifications(notif)
        if (priv) setPrivacy(priv)
      } catch {
        // use defaults
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const saveNotifications = async (updated: NotificationSettings) => {
    setNotifications(updated)
    setSaving(true)
    try {
      await apiClient("/api/v1/settings/notifications", { method: "PUT", body: updated })
      setSaved("notifications")
      setTimeout(() => setSaved(""), 2000)
    } catch { /* silent */ }
    finally { setSaving(false) }
  }

  const savePrivacy = async (updated: PrivacySettings) => {
    setPrivacy(updated)
    setSaving(true)
    try {
      await apiClient("/api/v1/settings/privacy", { method: "PUT", body: updated })
      setSaved("privacy")
      setTimeout(() => setSaved(""), 2000)
    } catch { /* silent */ }
    finally { setSaving(false) }
  }

  if (loading) {
    return (
      <div className="fixed inset-0 bg-background z-50 flex items-center justify-center">
        <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  return (
    <div className="fixed inset-0 bg-background z-50 flex flex-col">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-border bg-card">
        <h1 className="text-lg font-bold text-card-foreground">Configuración</h1>
        <button onClick={onClose} className="p-2 text-muted-foreground hover:text-card-foreground rounded-full">
          <X className="w-5 h-5" />
        </button>
      </div>

      <div className="flex-1 overflow-y-auto p-4">
        {/* Account info */}
        {user && (
          <div className="flex items-center gap-3 mb-6 p-4 bg-card rounded-xl border border-border">
            <div className="w-12 h-12 rounded-full overflow-hidden bg-secondary">
              <img
                src={user.photos?.[0]?.url ?? `${AVATAR_BASE}?u=${user.id}`}
                alt={user.name}
                className="w-full h-full object-cover"
                referrerPolicy="no-referrer"
              />
            </div>
            <div>
              <p className="font-semibold text-card-foreground">{user.name}</p>
              <p className="text-sm text-muted-foreground">{user.email}</p>
            </div>
          </div>
        )}

        {/* Search preferences */}
        <Section title="Preferencias de búsqueda">
          <SettingRow label="Edad mínima" description={`${searchPrefs.min_age} años`}>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setSearchPrefs((p) => ({ ...p, min_age: Math.max(18, p.min_age - 1) }))}
                className="w-7 h-7 rounded-full bg-secondary text-card-foreground flex items-center justify-center text-sm font-bold"
              >
                −
              </button>
              <span className="w-8 text-center text-sm font-medium text-card-foreground">{searchPrefs.min_age}</span>
              <button
                onClick={() => setSearchPrefs((p) => ({ ...p, min_age: Math.min(p.max_age - 1, p.min_age + 1) }))}
                className="w-7 h-7 rounded-full bg-secondary text-card-foreground flex items-center justify-center text-sm font-bold"
              >
                +
              </button>
            </div>
          </SettingRow>
          <SettingRow label="Edad máxima" description={`${searchPrefs.max_age} años`}>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setSearchPrefs((p) => ({ ...p, max_age: Math.max(p.min_age + 1, p.max_age - 1) }))}
                className="w-7 h-7 rounded-full bg-secondary text-card-foreground flex items-center justify-center text-sm font-bold"
              >
                −
              </button>
              <span className="w-8 text-center text-sm font-medium text-card-foreground">{searchPrefs.max_age}</span>
              <button
                onClick={() => setSearchPrefs((p) => ({ ...p, max_age: Math.min(99, p.max_age + 1) }))}
                className="w-7 h-7 rounded-full bg-secondary text-card-foreground flex items-center justify-center text-sm font-bold"
              >
                +
              </button>
            </div>
          </SettingRow>
          <SettingRow label="Distancia máxima" description={`${searchPrefs.max_distance} km`} last>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setSearchPrefs((p) => ({ ...p, max_distance: Math.max(5, p.max_distance - 5) }))}
                className="w-7 h-7 rounded-full bg-secondary text-card-foreground flex items-center justify-center text-sm font-bold"
              >
                −
              </button>
              <span className="w-10 text-center text-sm font-medium text-card-foreground">{searchPrefs.max_distance}</span>
              <button
                onClick={() => setSearchPrefs((p) => ({ ...p, max_distance: Math.min(500, p.max_distance + 5) }))}
                className="w-7 h-7 rounded-full bg-secondary text-card-foreground flex items-center justify-center text-sm font-bold"
              >
                +
              </button>
            </div>
          </SettingRow>
        </Section>

        {/* Notifications */}
        <Section title={`Notificaciones${saved === "notifications" ? " ✓" : ""}`}>
          <SettingRow label="Nuevos matches" description="Alerta cuando alguien te da like de vuelta">
            <Toggle
              checked={notifications.new_matches}
              onChange={(v) => saveNotifications({ ...notifications, new_matches: v })}
            />
          </SettingRow>
          <SettingRow label="Nuevos mensajes" description="Alerta cuando recibes un mensaje">
            <Toggle
              checked={notifications.new_messages}
              onChange={(v) => saveNotifications({ ...notifications, new_messages: v })}
            />
          </SettingRow>
          <SettingRow label="Noticias y actualizaciones" description="Artículos y novedades de la plataforma">
            <Toggle
              checked={notifications.news_updates}
              onChange={(v) => saveNotifications({ ...notifications, news_updates: v })}
            />
          </SettingRow>
          <SettingRow label="Marketing" description="Ofertas y promociones especiales" last>
            <Toggle
              checked={notifications.marketing}
              onChange={(v) => saveNotifications({ ...notifications, marketing: v })}
            />
          </SettingRow>
        </Section>

        {/* Privacy */}
        <Section title={`Privacidad${saved === "privacy" ? " ✓" : ""}`}>
          <SettingRow label="Estado en línea" description="Mostrar cuando estás activo">
            <Toggle
              checked={privacy.show_online_status}
              onChange={(v) => savePrivacy({ ...privacy, show_online_status: v })}
            />
          </SettingRow>
          <SettingRow label="Última vez visto" description="Mostrar tu última conexión">
            <Toggle
              checked={privacy.show_last_seen}
              onChange={(v) => savePrivacy({ ...privacy, show_last_seen: v })}
            />
          </SettingRow>
          <SettingRow label="Mostrar distancia" description="Mostrar a qué distancia estás de otros usuarios">
            <Toggle
              checked={privacy.show_distance}
              onChange={(v) => savePrivacy({ ...privacy, show_distance: v })}
            />
          </SettingRow>
          <SettingRow label="Modo incógnito" description="Tu perfil no aparecerá en búsquedas" last>
            <Toggle
              checked={privacy.incognito_mode}
              onChange={(v) => savePrivacy({ ...privacy, incognito_mode: v })}
            />
          </SettingRow>
        </Section>

        {/* Account */}
        <Section title="Cuenta">
          <SettingRow label="Cambiar contraseña" last>
            <ChevronRight className="w-4 h-4 text-muted-foreground" />
          </SettingRow>
        </Section>

        {/* App info */}
        <div className="text-center text-xs text-muted-foreground mt-4 pb-8">
          <p>MatchHub v1.0.0</p>
          <p className="mt-1">Hecho con ❤️ para conectar personas</p>
        </div>
      </div>
    </div>
  )
}
