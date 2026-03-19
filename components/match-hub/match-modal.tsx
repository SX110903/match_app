"use client"

import { motion, AnimatePresence } from "framer-motion"
import { Heart, MessageCircle, X } from "lucide-react"
import { Profile } from "@/lib/types"
import { AVATAR_BASE } from "@/lib/constants"
import { Button } from "@/components/ui/button"

interface MatchModalProps {
  profile: Profile | null
  isOpen: boolean
  onClose: () => void
  onSendMessage: () => void
  currentUserPhoto?: string
}

export function MatchModal({ profile, isOpen, onClose, onSendMessage, currentUserPhoto }: MatchModalProps) {
  if (!profile) return null

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 z-50 bg-gradient-to-b from-primary/90 to-primary flex items-center justify-center p-6"
        >
          <motion.div
            initial={{ scale: 0.5, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.5, opacity: 0 }}
            transition={{ type: "spring", stiffness: 300, damping: 25 }}
            className="text-center max-w-sm"
          >
            {/* Botón cerrar */}
            <button
              onClick={onClose}
              className="absolute top-6 right-6 p-2 text-primary-foreground/80 hover:text-primary-foreground transition-colors"
              aria-label="Cerrar"
            >
              <X className="w-8 h-8" />
            </button>

            {/* Título */}
            <motion.div
              initial={{ y: -20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.2 }}
            >
              <h1 className="text-4xl font-bold text-primary-foreground mb-2">
                {"¡Es un Match!"}
              </h1>
              <p className="text-primary-foreground/80">
                Tú y {profile.name} se han gustado mutuamente
              </p>
            </motion.div>

            {/* Fotos */}
            <motion.div
              initial={{ scale: 0.8, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              transition={{ delay: 0.3 }}
              className="flex justify-center items-center gap-4 my-8"
            >
              <div className="relative">
                <div className="w-32 h-32 rounded-full overflow-hidden border-4 border-primary-foreground shadow-xl">
                  <img
                    src={currentUserPhoto ?? `${AVATAR_BASE}?u=me`}
                    alt="Tu perfil"
                    className="w-full h-full object-cover"
                    referrerPolicy="no-referrer"
                  />
                </div>
              </div>
              
              <motion.div
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ repeat: Infinity, duration: 1.5 }}
                className="p-3 bg-primary-foreground rounded-full"
              >
                <Heart className="w-8 h-8 text-primary fill-primary" />
              </motion.div>

              <div className="relative">
                <div className="w-32 h-32 rounded-full overflow-hidden border-4 border-primary-foreground shadow-xl">
                  <img
                    src={profile.images[0]}
                    alt={profile.name}
                    className="w-full h-full object-cover"
                    referrerPolicy="no-referrer"
                  />
                </div>
              </div>
            </motion.div>

            {/* Botones */}
            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.5 }}
              className="space-y-3"
            >
              <Button
                onClick={onSendMessage}
                className="w-full bg-primary-foreground text-primary hover:bg-primary-foreground/90 font-semibold"
                size="lg"
              >
                <MessageCircle className="w-5 h-5 mr-2" />
                Enviar mensaje
              </Button>
              <Button
                onClick={onClose}
                variant="ghost"
                className="w-full text-primary-foreground hover:bg-primary-foreground/10"
                size="lg"
              >
                Seguir deslizando
              </Button>
            </motion.div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}
