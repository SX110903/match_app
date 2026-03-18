"use client"

import { useState } from "react"
import { motion } from "framer-motion"
import { Settings, Edit2, MapPin, Briefcase, Camera, Trash2, Plus, Loader2, LogOut } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { useAuth } from "@/lib/auth-context"
import { apiClient, APIError } from "@/lib/api-client"

export function ProfileView() {
  const { user, logout, refreshUser } = useAuth()
  const [showAddPhoto, setShowAddPhoto] = useState(false)
  const [photoUrl, setPhotoUrl] = useState("")
  const [photoError, setPhotoError] = useState("")
  const [addingPhoto, setAddingPhoto] = useState(false)
  const [deletingId, setDeletingId] = useState<string | null>(null)

  if (!user) return null

  const photos = user.photos ?? []

  const handleAddPhoto = async () => {
    if (!photoUrl.trim()) return
    setPhotoError("")
    setAddingPhoto(true)
    try {
      await apiClient('/api/v1/users/me/photos', {
        method: 'POST',
        body: { url: photoUrl.trim() },
      })
      setPhotoUrl("")
      setShowAddPhoto(false)
      await refreshUser()
    } catch (err) {
      if (err instanceof APIError) {
        if (err.status === 400) setPhotoError("URL inválida. Usa un enlace de Imgur (https://i.imgur.com/...)")
        else if (err.status === 409) setPhotoError("Máximo 6 fotos permitidas")
        else setPhotoError("Error al agregar la foto")
      }
    } finally {
      setAddingPhoto(false)
    }
  }

  const handleDeletePhoto = async (photoId: string) => {
    setDeletingId(photoId)
    try {
      await apiClient(`/api/v1/users/me/photos/${photoId}`, { method: 'DELETE' })
      await refreshUser()
    } catch { /* silent */ }
    finally { setDeletingId(null) }
  }

  const mainPhoto = photos[0]?.url ?? `https://i.pravatar.cc/400?u=${user.id}`

  return (
    <div className="flex-1 overflow-y-auto pb-24">
      <header className="flex items-center justify-between p-4">
        <h1 className="text-xl font-bold text-card-foreground">Mi Perfil</h1>
        <Button variant="ghost" size="icon" onClick={logout} aria-label="Cerrar sesión">
          <LogOut className="w-5 h-5" />
        </Button>
      </header>

      {/* Main photo */}
      <div className="flex flex-col items-center px-4 pb-6">
        <motion.div initial={{ scale: 0.9, opacity: 0 }} animate={{ scale: 1, opacity: 1 }} className="relative">
          <div className="w-32 h-32 rounded-full overflow-hidden ring-4 ring-primary ring-offset-4 ring-offset-background">
            <img src={mainPhoto} alt={user.name} className="w-full h-full object-cover" crossOrigin="anonymous" />
          </div>
          <button
            onClick={() => setShowAddPhoto(!showAddPhoto)}
            className="absolute bottom-0 right-0 p-2 bg-primary rounded-full text-primary-foreground shadow-lg hover:bg-primary/90 transition-colors"
            aria-label="Agregar foto"
          >
            <Camera className="w-5 h-5" />
          </button>
        </motion.div>

        {/* Add photo form */}
        {showAddPhoto && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            className="mt-4 w-full max-w-xs space-y-2"
          >
            <p className="text-xs text-muted-foreground text-center">Pega un enlace de Imgur</p>
            <div className="flex gap-2">
              <Input
                placeholder="https://i.imgur.com/abc.jpg"
                value={photoUrl}
                onChange={(e) => setPhotoUrl(e.target.value)}
                className="text-sm"
              />
              <Button size="sm" onClick={handleAddPhoto} disabled={addingPhoto || !photoUrl.trim()}>
                {addingPhoto ? <Loader2 className="w-4 h-4 animate-spin" /> : <Plus className="w-4 h-4" />}
              </Button>
            </div>
            {photoError && <p className="text-xs text-destructive text-center">{photoError}</p>}
          </motion.div>
        )}

        <h2 className="text-2xl font-bold text-card-foreground mt-4">{user.name}, {user.age}</h2>

        {user.occupation && (
          <div className="flex items-center gap-2 mt-1 text-muted-foreground">
            <Briefcase className="w-4 h-4" />
            <span>{user.occupation}</span>
          </div>
        )}

        {user.location && (
          <div className="flex items-center gap-2 mt-1 text-muted-foreground">
            <MapPin className="w-4 h-4" />
            <span>{user.location}</span>
          </div>
        )}

        <Button variant="outline" className="mt-4 gap-2">
          <Edit2 className="w-4 h-4" />
          Editar perfil
        </Button>
      </div>

      {/* Photo gallery */}
      {photos.length > 0 && (
        <div className="px-4 pb-4">
          <h3 className="font-semibold text-card-foreground mb-3">Mis fotos</h3>
          <div className="grid grid-cols-3 gap-2">
            {photos.map((photo, i) => (
              <div key={photo.id} className="relative aspect-square rounded-xl overflow-hidden group">
                <img src={photo.url} alt={`Foto ${i + 1}`} className="w-full h-full object-cover" crossOrigin="anonymous" />
                <button
                  onClick={() => handleDeletePhoto(photo.id)}
                  disabled={deletingId === photo.id}
                  className="absolute top-1 right-1 p-1 bg-black/60 rounded-full text-white opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  {deletingId === photo.id ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Bio */}
      {user.bio && (
        <div className="p-4">
          <h3 className="font-semibold text-card-foreground mb-2">Sobre mí</h3>
          <p className="text-muted-foreground">{user.bio}</p>
        </div>
      )}

      {/* Interests */}
      {user.interests?.length > 0 && (
        <div className="p-4">
          <h3 className="font-semibold text-card-foreground mb-3">Mis intereses</h3>
          <div className="flex flex-wrap gap-2">
            {user.interests.map((interest, index) => (
              <motion.div
                key={interest}
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: index * 0.05 }}
              >
                <Badge variant="secondary" className="bg-primary/10 text-primary px-3 py-1">
                  {interest}
                </Badge>
              </motion.div>
            ))}
          </div>
        </div>
      )}

      {/* Settings */}
      <div className="p-4 space-y-2">
        <h3 className="font-semibold text-card-foreground mb-3">Configuración</h3>
        {[
          { label: "Preferencias de búsqueda", desc: "Edad, distancia, género" },
          { label: "Notificaciones", desc: "Matches, mensajes, likes" },
          { label: "Privacidad", desc: "Visibilidad, bloqueos" },
        ].map((item, index) => (
          <motion.button
            key={item.label}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.3 + index * 0.05 }}
            className="w-full flex items-center justify-between p-4 bg-card rounded-xl hover:bg-muted/50 transition-colors border border-border"
          >
            <div className="text-left">
              <p className="font-medium text-card-foreground">{item.label}</p>
              <p className="text-sm text-muted-foreground">{item.desc}</p>
            </div>
            <span className="text-muted-foreground">→</span>
          </motion.button>
        ))}
      </div>
    </div>
  )
}
