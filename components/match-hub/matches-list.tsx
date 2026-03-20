"use client"

import { useState } from "react"
import { motion } from "framer-motion"
import { Trash2 } from "lucide-react"
import { Match } from "@/lib/types"
import { BadgeIcon } from "@/components/match-hub/badge"
import { apiClient, APIError } from "@/lib/api-client"

interface MatchesListProps {
  matches: Match[]
  onSelectMatch: (match: Match) => void
  onMatchDeleted?: (matchId: string) => void
}

export function MatchesList({ matches, onSelectMatch, onMatchDeleted }: MatchesListProps) {
  const [confirmDelete, setConfirmDelete] = useState<Match | null>(null)
  const [deleting, setDeleting] = useState(false)

  const handleDelete = async () => {
    if (!confirmDelete) return
    setDeleting(true)
    try {
      await apiClient(`/api/v1/matches/${confirmDelete.id}`, { method: "DELETE" })
      onMatchDeleted?.(confirmDelete.id)
      setConfirmDelete(null)
    } catch (e) {
      console.error("[MatchesList] delete failed", e instanceof APIError ? e.message : e)
    } finally {
      setDeleting(false)
    }
  }

  const formatTime = (date: Date) => {
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (minutes < 60) return `hace ${minutes}m`
    if (hours < 24) return `hace ${hours}h`
    return `hace ${days}d`
  }

  const newMatches = matches.filter((m) => !m.lastMessage)
  const conversations = matches.filter((m) => m.lastMessage)

  return (
    <div className="flex-1 overflow-y-auto">
      {/* Confirm delete modal */}
      {confirmDelete && (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
          <div className="bg-card rounded-2xl p-6 max-w-sm w-full shadow-xl">
            <p className="text-sm text-card-foreground mb-1 font-semibold">¿Eliminar match?</p>
            <p className="text-xs text-muted-foreground mb-6">
              El match con <span className="font-medium">{confirmDelete.profile.name}</span> se eliminará solo para ti.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setConfirmDelete(null)}
                className="flex-1 py-2 rounded-xl border border-border text-sm font-medium text-card-foreground hover:bg-muted/50 transition-colors"
              >
                Cancelar
              </button>
              <button
                onClick={handleDelete}
                disabled={deleting}
                className="flex-1 py-2 rounded-xl bg-destructive text-destructive-foreground text-sm font-medium disabled:opacity-50 transition-colors"
              >
                {deleting ? "Eliminando..." : "Eliminar"}
              </button>
            </div>
          </div>
        </div>
      )}
      {/* Nuevos matches */}
      {newMatches.length > 0 && (
        <div className="p-4">
          <h3 className="text-sm font-semibold text-muted-foreground mb-3">Nuevos Matches</h3>
          <div className="flex gap-4 overflow-x-auto pb-2">
            {newMatches.map((match, index) => (
              <motion.button
                key={match.id}
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ delay: index * 0.1 }}
                onClick={() => onSelectMatch(match)}
                className="flex flex-col items-center gap-2 flex-shrink-0"
              >
                <div className="relative">
                  <div className="w-20 h-20 rounded-full overflow-hidden ring-2 ring-primary ring-offset-2 ring-offset-background">
                    <img
                      src={match.profile.images[0]}
                      alt={match.profile.name}
                      className="w-full h-full object-cover"
                      referrerPolicy="no-referrer"
                    />
                  </div>
                  <div className="absolute -bottom-1 -right-1 w-6 h-6 bg-primary rounded-full flex items-center justify-center">
                    <span className="text-primary-foreground text-xs">✨</span>
                  </div>
                </div>
                <span className="text-sm font-medium text-card-foreground">{match.profile.name}</span>
              </motion.button>
            ))}
          </div>
        </div>
      )}

      {/* Conversaciones */}
      <div className="p-4">
        <h3 className="text-sm font-semibold text-muted-foreground mb-3">Mensajes</h3>
        <div className="space-y-2">
          {conversations.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">
              Aún no tienes conversaciones. ¡Empieza a chatear con tus matches!
            </p>
          ) : (
            conversations.map((match, index) => (
              <motion.div
                key={match.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: index * 0.05 }}
                className="flex items-center gap-1"
              >
              <button
                onClick={() => onSelectMatch(match)}
                className="flex-1 flex items-center gap-3 p-3 rounded-xl hover:bg-muted/50 transition-colors min-w-0"
              >
                <div className="relative">
                  <div className="w-16 h-16 rounded-full overflow-hidden">
                    <img
                      src={match.profile.images[0]}
                      alt={match.profile.name}
                      className="w-full h-full object-cover"
                      referrerPolicy="no-referrer"
                    />
                  </div>
                  {match.unread && (
                    <div className="absolute bottom-0 right-0 w-4 h-4 bg-primary rounded-full border-2 border-card" />
                  )}
                </div>
                <div className="flex-1 text-left">
                  <div className="flex items-center justify-between">
                    <h4 className="font-semibold text-card-foreground flex items-center gap-1">
                      {match.profile.name}
                      <BadgeIcon badge={match.profile.badge} />
                    </h4>
                    <span className="text-xs text-muted-foreground">{formatTime(match.matchedAt)}</span>
                  </div>
                  <p
                    className={`text-sm truncate ${
                      match.unread ? "text-card-foreground font-medium" : "text-muted-foreground"
                    }`}
                  >
                    {match.lastMessage}
                  </p>
                </div>
              </button>
              <button
                onClick={() => setConfirmDelete(match)}
                className="flex-shrink-0 p-2 text-muted-foreground hover:text-destructive transition-colors rounded-lg"
                aria-label="Eliminar match"
              >
                <Trash2 className="w-4 h-4" />
              </button>
              </motion.div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
