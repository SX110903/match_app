"use client"

import { useState } from "react"
import { motion, AnimatePresence } from "framer-motion"
import { X, MapPin, Briefcase, ChevronLeft, ChevronRight } from "lucide-react"
import { Profile } from "@/lib/types"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"

interface ProfileModalProps {
  profile: Profile | null
  isOpen: boolean
  onClose: () => void
}

export function ProfileModal({ profile, isOpen, onClose }: ProfileModalProps) {
  const [currentImageIndex, setCurrentImageIndex] = useState(0)

  if (!profile) return null

  const nextImage = () => {
    if (currentImageIndex < profile.images.length - 1) {
      setCurrentImageIndex(currentImageIndex + 1)
    }
  }

  const prevImage = () => {
    if (currentImageIndex > 0) {
      setCurrentImageIndex(currentImageIndex - 1)
    }
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4"
          onClick={onClose}
        >
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.9, opacity: 0 }}
            onClick={(e) => e.stopPropagation()}
            className="bg-card rounded-3xl overflow-hidden max-w-md w-full max-h-[90vh] overflow-y-auto shadow-2xl"
          >
            {/* Galería de imágenes */}
            <div className="relative aspect-[3/4]">
              <img
                src={profile.images[currentImageIndex]}
                alt={profile.name}
                className="w-full h-full object-cover"
                crossOrigin="anonymous"
              />
              
              {/* Indicadores */}
              <div className="absolute top-4 left-4 right-4 flex gap-1">
                {profile.images.map((_, index) => (
                  <div
                    key={index}
                    className={`h-1 flex-1 rounded-full transition-colors ${
                      index === currentImageIndex ? "bg-primary-foreground" : "bg-primary-foreground/40"
                    }`}
                  />
                ))}
              </div>

              {/* Navegación */}
              {profile.images.length > 1 && (
                <>
                  <button
                    onClick={prevImage}
                    className="absolute left-2 top-1/2 -translate-y-1/2 p-2 bg-black/30 rounded-full hover:bg-black/50 transition-colors"
                    aria-label="Imagen anterior"
                  >
                    <ChevronLeft className="w-6 h-6 text-primary-foreground" />
                  </button>
                  <button
                    onClick={nextImage}
                    className="absolute right-2 top-1/2 -translate-y-1/2 p-2 bg-black/30 rounded-full hover:bg-black/50 transition-colors"
                    aria-label="Siguiente imagen"
                  >
                    <ChevronRight className="w-6 h-6 text-primary-foreground" />
                  </button>
                </>
              )}

              {/* Cerrar */}
              <button
                onClick={onClose}
                className="absolute top-4 right-4 p-2 bg-black/30 rounded-full hover:bg-black/50 transition-colors"
                aria-label="Cerrar"
              >
                <X className="w-6 h-6 text-primary-foreground" />
              </button>

              {/* Gradiente */}
              <div className="absolute inset-x-0 bottom-0 h-32 bg-gradient-to-t from-card to-transparent" />
            </div>

            {/* Información */}
            <div className="p-6 -mt-16 relative">
              <h2 className="text-2xl font-bold text-card-foreground">
                {profile.name}, {profile.age}
              </h2>
              
              <div className="flex items-center gap-2 mt-2 text-muted-foreground">
                <Briefcase className="w-4 h-4" />
                <span>{profile.occupation}</span>
              </div>
              
              <div className="flex items-center gap-2 mt-1 text-muted-foreground">
                <MapPin className="w-4 h-4" />
                <span>a {profile.distance} km</span>
              </div>

              <div className="mt-4">
                <h3 className="font-semibold text-card-foreground mb-2">Sobre mí</h3>
                <p className="text-muted-foreground">{profile.bio}</p>
              </div>

              <div className="mt-4">
                <h3 className="font-semibold text-card-foreground mb-2">Intereses</h3>
                <div className="flex flex-wrap gap-2">
                  {profile.interests.map((interest) => (
                    <Badge key={interest} variant="secondary" className="bg-primary/10 text-primary">
                      {interest}
                    </Badge>
                  ))}
                </div>
              </div>

              <Button onClick={onClose} className="w-full mt-6 bg-primary hover:bg-primary/90 text-primary-foreground">
                Cerrar
              </Button>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}
