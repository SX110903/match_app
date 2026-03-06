"use client"

import { useState } from "react"
import { motion, useMotionValue, useTransform, PanInfo } from "framer-motion"
import { MapPin, Briefcase, Info, ChevronLeft, ChevronRight } from "lucide-react"
import { Profile } from "@/lib/types"
import { Badge } from "@/components/ui/badge"

interface SwipeCardProps {
  profile: Profile
  onSwipe: (direction: "left" | "right") => void
  onInfoClick: () => void
}

export function SwipeCard({ profile, onSwipe, onInfoClick }: SwipeCardProps) {
  const [currentImageIndex, setCurrentImageIndex] = useState(0)
  const [exitX, setExitX] = useState(0)

  const x = useMotionValue(0)
  const rotate = useTransform(x, [-200, 200], [-25, 25])
  const opacity = useTransform(x, [-200, -100, 0, 100, 200], [0.5, 1, 1, 1, 0.5])

  const likeOpacity = useTransform(x, [0, 100], [0, 1])
  const nopeOpacity = useTransform(x, [-100, 0], [1, 0])

  const handleDragEnd = (_: MouseEvent | TouchEvent | PointerEvent, info: PanInfo) => {
    if (info.offset.x > 100) {
      setExitX(300)
      onSwipe("right")
    } else if (info.offset.x < -100) {
      setExitX(-300)
      onSwipe("left")
    }
  }

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
    <motion.div
      className="absolute w-full h-full cursor-grab active:cursor-grabbing"
      drag="x"
      dragConstraints={{ left: 0, right: 0 }}
      onDragEnd={handleDragEnd}
      animate={{ x: exitX }}
      transition={{ type: "spring", stiffness: 300, damping: 30 }}
      style={{ x, rotate, opacity }}
    >
      <div className="relative w-full h-full rounded-3xl overflow-hidden shadow-2xl bg-card">
        {/* Imagen de fondo */}
        <div className="absolute inset-0">
          <img
            src={profile.images[currentImageIndex]}
            alt={profile.name}
            className="w-full h-full object-cover"
            crossOrigin="anonymous"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent" />
        </div>

        {/* Indicadores de imagen */}
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

        {/* Navegación de imágenes */}
        <div className="absolute inset-0 flex">
          <button
            onClick={prevImage}
            className="w-1/2 h-full flex items-center justify-start pl-2 opacity-0 hover:opacity-100 transition-opacity"
            aria-label="Imagen anterior"
          >
            <ChevronLeft className="w-8 h-8 text-primary-foreground drop-shadow-lg" />
          </button>
          <button
            onClick={nextImage}
            className="w-1/2 h-full flex items-center justify-end pr-2 opacity-0 hover:opacity-100 transition-opacity"
            aria-label="Siguiente imagen"
          >
            <ChevronRight className="w-8 h-8 text-primary-foreground drop-shadow-lg" />
          </button>
        </div>

        {/* Indicadores LIKE/NOPE */}
        <motion.div
          className="absolute top-20 left-6 border-4 border-green-500 text-green-500 px-4 py-2 rounded-lg font-bold text-3xl -rotate-12"
          style={{ opacity: likeOpacity }}
        >
          LIKE
        </motion.div>
        <motion.div
          className="absolute top-20 right-6 border-4 border-red-500 text-red-500 px-4 py-2 rounded-lg font-bold text-3xl rotate-12"
          style={{ opacity: nopeOpacity }}
        >
          NOPE
        </motion.div>

        {/* Información del perfil */}
        <div className="absolute bottom-0 left-0 right-0 p-6 text-primary-foreground">
          <div className="flex items-end justify-between">
            <div className="flex-1">
              <h2 className="text-3xl font-bold">
                {profile.name}, {profile.age}
              </h2>
              <div className="flex items-center gap-2 mt-1 text-primary-foreground/80">
                <Briefcase className="w-4 h-4" />
                <span className="text-sm">{profile.occupation}</span>
              </div>
              <div className="flex items-center gap-2 mt-1 text-primary-foreground/80">
                <MapPin className="w-4 h-4" />
                <span className="text-sm">a {profile.distance} km</span>
              </div>
              <div className="flex flex-wrap gap-2 mt-3">
                {profile.interests.slice(0, 3).map((interest) => (
                  <Badge
                    key={interest}
                    variant="secondary"
                    className="bg-primary-foreground/20 text-primary-foreground border-0"
                  >
                    {interest}
                  </Badge>
                ))}
              </div>
            </div>
            <button
              onClick={onInfoClick}
              className="p-3 bg-primary-foreground/20 rounded-full hover:bg-primary-foreground/30 transition-colors"
              aria-label="Ver más información"
            >
              <Info className="w-6 h-6 text-primary-foreground" />
            </button>
          </div>
        </div>
      </div>
    </motion.div>
  )
}
