"use client"

import { useState } from "react"
import { Trash2, LogOut, Edit2, X, Plus, Settings, Shield, Grid3X3, MapPin, Briefcase } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { useAuth } from "@/lib/auth-context"
import { apiClient, APIError } from "@/lib/api-client"
import { AVATAR_BASE } from "@/lib/constants"
import { BadgeIcon, type BadgeType } from "@/components/match-hub/badge"

const VIP_LABELS = ["", "Bronce", "Plata", "Oro", "Platino", "Diamante"]
const VIP_COLORS = [
  "",
  "text-amber-600",
  "text-slate-400",
  "text-yellow-500",
  "text-sky-400",
  "text-cyan-300",
]

function EditProfileModal({ onClose, onSaved }: { onClose: () => void; onSaved: () => Promise<void> }) {
  const { user } = useAuth()
  const [name, setName] = useState(user?.name ?? "")
  const [bio, setBio] = useState(user?.bio ?? "")
  const [occupation, setOccupation] = useState(user?.occupation ?? "")
  const [location, setLocation] = useState(user?.location ?? "")
  const [interests, setInterests] = useState<string[]>(user?.interests ?? [])
  const [newInterest, setNewInterest] = useState("")
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState("")

  const addInterest = () => {
    const trimmed = newInterest.trim()
    if (trimmed && !interests.includes(trimmed) && interests.length < 10) {
      setInterests((prev) => [...prev, trimmed])
      setNewInterest("")
    }
  }

  const removeInterest = (i: string) => setInterests((prev) => prev.filter((x) => x !== i))

  const handleSave = async () => {
    if (!name.trim()) { setError("El nombre es obligatorio"); return }
    setSaving(true)
    setError("")
    try {
      await apiClient("/api/v1/users/me", {
        method: "PUT",
        body: {
          name: name.trim(),
          bio: bio.trim() || null,
          occupation: occupation.trim() || null,
          location: location.trim() || null,
          interests,
        },
      })
      await onSaved()
      onClose()
    } catch (err) {
      setError(err instanceof APIError ? err.message : "Error al guardar")
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/60 flex items-end sm:items-center justify-center z-50 p-4">
      <div className="bg-card rounded-2xl w-full max-w-md max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-4 border-b border-border sticky top-0 bg-card z-10">
          <h3 className="font-semibold text-card-foreground">Editar perfil</h3>
          <button onClick={onClose} className="text-muted-foreground hover:text-card-foreground">
            <X className="w-5 h-5" />
          </button>
        </div>
        <div className="p-4 space-y-3">
          {error && <p className="text-sm text-destructive">{error}</p>}
          <div>
            <label className="text-xs font-medium text-muted-foreground mb-1 block">Nombre*</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
            />
          </div>
          <div>
            <label className="text-xs font-medium text-muted-foreground mb-1 block">Bio</label>
            <textarea
              value={bio}
              onChange={(e) => setBio(e.target.value.slice(0, 500))}
              rows={3}
              className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary resize-none"
              placeholder="Cuéntanos algo sobre ti..."
            />
          </div>
          <div>
            <label className="text-xs font-medium text-muted-foreground mb-1 block">Ocupación</label>
            <input
              value={occupation}
              onChange={(e) => setOccupation(e.target.value)}
              className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
              placeholder="Tu trabajo o profesión"
            />
          </div>
          <div>
            <label className="text-xs font-medium text-muted-foreground mb-1 block">Ubicación</label>
            <input
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
              placeholder="Ciudad, País"
            />
          </div>
          <div>
            <label className="text-xs font-medium text-muted-foreground mb-1 block">Intereses</label>
            <div className="flex flex-wrap gap-1.5 mb-2">
              {interests.map((interest) => (
                <span
                  key={interest}
                  className="flex items-center gap-1 text-xs bg-secondary text-secondary-foreground px-2 py-0.5 rounded-full"
                >
                  {interest}
                  <button onClick={() => removeInterest(interest)} className="hover:text-destructive">
                    <X className="w-3 h-3" />
                  </button>
                </span>
              ))}
            </div>
            <div className="flex gap-2">
              <input
                value={newInterest}
                onChange={(e) => setNewInterest(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && addInterest()}
                placeholder="Añadir interés..."
                className="flex-1 border border-border rounded-lg px-3 py-1.5 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
              />
              <Button size="sm" variant="outline" onClick={addInterest} disabled={!newInterest.trim()}>
                <Plus className="w-3.5 h-3.5" />
              </Button>
            </div>
          </div>
        </div>
        <div className="p-4 pt-0 flex gap-2">
          <Button variant="outline" className="flex-1" onClick={onClose}>
            Cancelar
          </Button>
          <Button className="flex-1 bg-primary hover:bg-primary/90" onClick={handleSave} disabled={saving}>
            {saving ? "Guardando..." : "Guardar"}
          </Button>
        </div>
      </div>
    </div>
  )
}

function AddPhotoModal({
  onClose,
  onAdd,
}: {
  onClose: () => void
  onAdd: (url: string) => Promise<void>
}) {
  const [photoUrl, setPhotoUrl] = useState("")
  const [error, setError] = useState("")
  const [adding, setAdding] = useState(false)

  const handleAdd = async () => {
    if (!photoUrl.trim()) return
    setError("")
    setAdding(true)
    try {
      await onAdd(photoUrl.trim())
      onClose()
    } catch (err) {
      if (err instanceof APIError) {
        if (err.status === 400) setError("URL inválida. Usa un enlace directo a imagen.")
        else if (err.status === 409) setError("Máximo 6 fotos permitidas")
        else setError("Error al agregar la foto")
      } else {
        setError("Error al agregar la foto")
      }
    } finally {
      setAdding(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/60 flex items-end sm:items-center justify-center z-50 p-4">
      <div className="bg-card rounded-2xl w-full max-w-sm">
        <div className="flex items-center justify-between p-4 border-b border-border">
          <h3 className="font-semibold text-card-foreground">Añadir foto</h3>
          <button onClick={onClose} className="text-muted-foreground hover:text-card-foreground">
            <X className="w-5 h-5" />
          </button>
        </div>
        <div className="p-4 space-y-3">
          <p className="text-xs text-muted-foreground">Pega un enlace directo a una imagen (jpg, png, webp)</p>
          <input
            value={photoUrl}
            onChange={(e) => setPhotoUrl(e.target.value)}
            placeholder="https://i.imgur.com/..."
            autoFocus
            className="w-full border border-border rounded-lg px-3 py-2 text-sm bg-background focus:outline-none focus:ring-1 focus:ring-primary"
          />
          {error && <p className="text-xs text-destructive">{error}</p>}
        </div>
        <div className="p-4 pt-0 flex gap-2">
          <Button variant="outline" className="flex-1" onClick={onClose}>
            Cancelar
          </Button>
          <Button
            className="flex-1 bg-primary hover:bg-primary/90"
            onClick={handleAdd}
            disabled={adding || !photoUrl.trim()}
          >
            {adding ? "Añadiendo..." : "Añadir"}
          </Button>
        </div>
      </div>
    </div>
  )
}

export function ProfileView({
  onOpenSettings,
  onOpenAdmin,
  onOpenShop,
}: {
  onOpenSettings?: () => void
  onOpenAdmin?: () => void
  onOpenShop?: () => void
}) {
  const { user, logout, refreshUser } = useAuth()
  const [showAddPhoto, setShowAddPhoto] = useState(false)
  const [deletingId, setDeletingId] = useState<string | null>(null)
  const [showEditProfile, setShowEditProfile] = useState(false)

  if (!user) return null

  const photos = user.photos ?? []
  const mainPhoto = photos[0]?.url ?? `${AVATAR_BASE}?u=${user.id}`

  const profileCompleteness = Math.min(
    100,
    [user.bio, user.occupation, user.location, photos.length > 0, (user.interests?.length ?? 0) > 0].filter(
      Boolean
    ).length * 20
  )

  const handleAddPhoto = async (url: string) => {
    await apiClient("/api/v1/users/me/photos", { method: "POST", body: { url } })
    await refreshUser()
  }

  const handleDeletePhoto = async (photoId: string) => {
    setDeletingId(photoId)
    try {
      await apiClient(`/api/v1/users/me/photos/${photoId}`, { method: "DELETE" })
      await refreshUser()
    } catch {
      //
    } finally {
      setDeletingId(null)
    }
  }

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Top bar */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-border bg-background sticky top-0 z-10">
        <span className="font-bold text-foreground text-base truncate max-w-[60%]">{user.name}</span>
        <div className="flex items-center gap-0.5">
          {onOpenAdmin && (
            <button
              onClick={onOpenAdmin}
              className="p-2 rounded-full hover:bg-muted transition-colors"
              aria-label="Panel Admin"
            >
              <Shield className="w-5 h-5 text-destructive" />
            </button>
          )}
          {onOpenSettings && (
            <button
              onClick={onOpenSettings}
              className="p-2 rounded-full hover:bg-muted transition-colors"
              aria-label="Configuración"
            >
              <Settings className="w-5 h-5 text-foreground" />
            </button>
          )}
          <button
            onClick={logout}
            className="p-2 rounded-full hover:bg-muted transition-colors"
            aria-label="Cerrar sesión"
          >
            <LogOut className="w-5 h-5 text-foreground" />
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto pb-24">
        {/* ── Instagram header ── */}
        <div className="px-4 pt-5 pb-4">
          {/* Row: avatar + stats */}
          <div className="flex items-center gap-5 mb-4">
            {/* Profile pic */}
            <div className="relative flex-shrink-0">
              <div className="w-20 h-20 rounded-full overflow-hidden ring-2 ring-primary ring-offset-2 ring-offset-background">
                <img
                  src={mainPhoto}
                  alt={user.name}
                  className="w-full h-full object-cover"
                  referrerPolicy="no-referrer"
                />
              </div>
            </div>

            {/* Stats */}
            <div className="flex-1 flex items-center justify-around text-center">
              <div>
                <p className="text-base font-bold text-foreground">{photos.length}</p>
                <p className="text-xs text-muted-foreground">Fotos</p>
              </div>
              <div>
                <p className="text-base font-bold text-foreground">{user.credits}</p>
                <p className="text-xs text-muted-foreground">Créditos</p>
              </div>
              <div>
                <p className="text-base font-bold text-foreground">{profileCompleteness}%</p>
                <p className="text-xs text-muted-foreground">Perfil</p>
              </div>
            </div>
          </div>

          {/* Name + badges */}
          <div className="mb-1 flex items-center gap-2 flex-wrap">
            <p className="font-bold text-foreground text-sm flex items-center gap-1">
              {user.name}
              <BadgeIcon badge={user.badge as BadgeType | undefined} />
            </p>
            {user.age && <p className="text-sm text-muted-foreground">{user.age} años</p>}
            {user.vip_level > 0 && (
              <span className={`text-xs font-bold ${VIP_COLORS[user.vip_level]}`}>
                ★ {VIP_LABELS[user.vip_level]}
              </span>
            )}
            {user.is_admin && (
              <Badge variant="destructive" className="text-xs px-1.5 py-0">
                Admin
              </Badge>
            )}
          </div>

          {/* Occupation */}
          {user.occupation && (
            <p className="text-xs text-muted-foreground flex items-center gap-1 mb-0.5">
              <Briefcase className="w-3 h-3" />
              {user.occupation}
            </p>
          )}

          {/* Location */}
          {user.location && (
            <p className="text-xs text-muted-foreground flex items-center gap-1 mb-1">
              <MapPin className="w-3 h-3" />
              {user.location}
            </p>
          )}

          {/* Bio */}
          {user.bio && (
            <p className="text-sm text-foreground leading-snug mt-1 mb-2 whitespace-pre-wrap">{user.bio}</p>
          )}

          {/* Interests */}
          {user.interests && user.interests.length > 0 && (
            <div className="flex flex-wrap gap-1.5 mt-2 mb-3">
              {user.interests.map((interest) => (
                <Badge key={interest} variant="secondary" className="text-xs rounded-full">
                  {interest}
                </Badge>
              ))}
            </div>
          )}

          {/* Action buttons — Instagram style */}
          <div className="flex gap-2 mt-3">
            <button
              onClick={() => setShowEditProfile(true)}
              className="flex-1 bg-secondary hover:bg-secondary/70 text-foreground text-sm font-semibold py-1.5 rounded-lg transition-colors flex items-center justify-center gap-1.5"
            >
              <Edit2 className="w-3.5 h-3.5" />
              Editar perfil
            </button>
            <button
              onClick={() => setShowAddPhoto(true)}
              className="flex-1 bg-secondary hover:bg-secondary/70 text-foreground text-sm font-semibold py-1.5 rounded-lg transition-colors flex items-center justify-center gap-1.5"
            >
              <Plus className="w-3.5 h-3.5" />
              Añadir foto
            </button>
          </div>
          {onOpenShop && (
            <button
              onClick={onOpenShop}
              className="w-full mt-2 bg-primary/10 hover:bg-primary/20 text-primary text-sm font-semibold py-1.5 rounded-lg transition-colors flex items-center justify-center gap-1.5"
            >
              🛒 Ir a la tienda
            </button>
          )}
        </div>

        {/* ── Story highlights row ── */}
        {photos.length > 0 && (
          <div className="border-t border-border/50 pt-4 pb-3">
            <div className="flex gap-4 px-4 overflow-x-auto no-scrollbar">
              {photos.slice(0, 6).map((photo, i) => (
                <div key={photo.id} className="flex flex-col items-center gap-1.5 flex-shrink-0">
                  <div className="w-16 h-16 rounded-full overflow-hidden ring-2 ring-primary ring-offset-2 ring-offset-background">
                    <img
                      src={photo.url}
                      alt={`Foto ${i + 1}`}
                      className="w-full h-full object-cover rounded-full"
                      referrerPolicy="no-referrer"
                    />
                  </div>
                  <span className="text-[10px] text-muted-foreground">Foto {i + 1}</span>
                </div>
              ))}
              <div className="flex flex-col items-center gap-1.5 flex-shrink-0">
                <button
                  onClick={() => setShowAddPhoto(true)}
                  className="w-16 h-16 rounded-full border-2 border-dashed border-border flex items-center justify-center hover:border-primary transition-colors"
                >
                  <Plus className="w-5 h-5 text-muted-foreground" />
                </button>
                <span className="text-[10px] text-muted-foreground">Nueva</span>
              </div>
            </div>
          </div>
        )}

        {/* ── Photo grid divider ── */}
        <div className="border-t border-border flex items-center justify-center py-2">
          <Grid3X3 className="w-4 h-4 text-muted-foreground" />
        </div>

        {/* ── Photo grid — tight 3-col Instagram style ── */}
        {photos.length === 0 ? (
          <div className="text-center py-16 px-8">
            <div className="w-16 h-16 rounded-full border-2 border-border flex items-center justify-center mx-auto mb-3">
              <Grid3X3 className="w-7 h-7 text-muted-foreground" />
            </div>
            <p className="font-bold text-foreground mb-1">Sin fotos aún</p>
            <p className="text-sm text-muted-foreground">Comparte tu primera foto</p>
            <button
              onClick={() => setShowAddPhoto(true)}
              className="mt-4 bg-primary text-primary-foreground text-sm font-semibold px-5 py-2 rounded-lg hover:bg-primary/90 transition-colors"
            >
              Añadir foto
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-3 gap-px bg-border">
            {photos.map((photo) => (
              <div key={photo.id} className="relative aspect-square group bg-background">
                <img
                  src={photo.url}
                  alt="Foto"
                  className="w-full h-full object-cover"
                  referrerPolicy="no-referrer"
                />
                <button
                  onClick={() => handleDeletePhoto(photo.id)}
                  disabled={deletingId === photo.id}
                  className="absolute top-1.5 right-1.5 p-1 bg-black/50 backdrop-blur-sm rounded-full text-white opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  {deletingId === photo.id ? (
                    <div className="w-3 h-3 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  ) : (
                    <Trash2 className="w-3 h-3" />
                  )}
                </button>
              </div>
            ))}
            {/* Add tile at end of grid */}
            <button
              onClick={() => setShowAddPhoto(true)}
              className="aspect-square bg-muted flex items-center justify-center hover:bg-muted/70 transition-colors"
            >
              <Plus className="w-6 h-6 text-muted-foreground" />
            </button>
          </div>
        )}
      </div>

      {showEditProfile && (
        <EditProfileModal onClose={() => setShowEditProfile(false)} onSaved={refreshUser} />
      )}
      {showAddPhoto && (
        <AddPhotoModal onClose={() => setShowAddPhoto(false)} onAdd={handleAddPhoto} />
      )}
    </div>
  )
}
