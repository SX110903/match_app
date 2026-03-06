"use client"

import { motion } from "framer-motion"
import { Settings, Edit2, MapPin, Briefcase, Camera } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"

export function ProfileView() {
  const myProfile = {
    name: "Carlos",
    age: 28,
    bio: "Desarrollador apasionado por la tecnología y los viajes. Siempre en busca de nuevas aventuras y conexiones genuinas.",
    images: [
      "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&q=80",
    ],
    occupation: "Desarrollador Full Stack",
    interests: ["Viajes", "Tecnología", "Música", "Fotografía", "Cocina"],
    location: "Ciudad de México",
  }

  return (
    <div className="flex-1 overflow-y-auto pb-24">
      {/* Header */}
      <header className="flex items-center justify-between p-4">
        <h1 className="text-xl font-bold text-card-foreground">Mi Perfil</h1>
        <Button variant="ghost" size="icon" aria-label="Configuración">
          <Settings className="w-6 h-6" />
        </Button>
      </header>

      {/* Profile photo */}
      <div className="flex flex-col items-center px-4 pb-6">
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          className="relative"
        >
          <div className="w-32 h-32 rounded-full overflow-hidden ring-4 ring-primary ring-offset-4 ring-offset-background">
            <img
              src={myProfile.images[0]}
              alt={myProfile.name}
              className="w-full h-full object-cover"
              crossOrigin="anonymous"
            />
          </div>
          <button
            className="absolute bottom-0 right-0 p-2 bg-primary rounded-full text-primary-foreground shadow-lg hover:bg-primary/90 transition-colors"
            aria-label="Cambiar foto"
          >
            <Camera className="w-5 h-5" />
          </button>
        </motion.div>

        <h2 className="text-2xl font-bold text-card-foreground mt-4">
          {myProfile.name}, {myProfile.age}
        </h2>
        
        <div className="flex items-center gap-2 mt-1 text-muted-foreground">
          <Briefcase className="w-4 h-4" />
          <span>{myProfile.occupation}</span>
        </div>
        
        <div className="flex items-center gap-2 mt-1 text-muted-foreground">
          <MapPin className="w-4 h-4" />
          <span>{myProfile.location}</span>
        </div>

        <Button variant="outline" className="mt-4 gap-2">
          <Edit2 className="w-4 h-4" />
          Editar perfil
        </Button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-4 px-4 py-6 border-y border-border">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="text-center"
        >
          <p className="text-2xl font-bold text-primary">24</p>
          <p className="text-sm text-muted-foreground">Matches</p>
        </motion.div>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="text-center"
        >
          <p className="text-2xl font-bold text-primary">156</p>
          <p className="text-sm text-muted-foreground">Likes</p>
        </motion.div>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="text-center"
        >
          <p className="text-2xl font-bold text-primary">89%</p>
          <p className="text-sm text-muted-foreground">Perfil</p>
        </motion.div>
      </div>

      {/* About */}
      <div className="p-4">
        <h3 className="font-semibold text-card-foreground mb-2">Sobre mí</h3>
        <p className="text-muted-foreground">{myProfile.bio}</p>
      </div>

      {/* Interests */}
      <div className="p-4">
        <h3 className="font-semibold text-card-foreground mb-3">Mis intereses</h3>
        <div className="flex flex-wrap gap-2">
          {myProfile.interests.map((interest, index) => (
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

      {/* Settings quick links */}
      <div className="p-4 space-y-2">
        <h3 className="font-semibold text-card-foreground mb-3">Configuración</h3>
        {[
          { label: "Preferencias de búsqueda", desc: "Edad, distancia, género" },
          { label: "Notificaciones", desc: "Matches, mensajes, likes" },
          { label: "Privacidad", desc: "Visibilidad, bloqueos" },
          { label: "Verificar perfil", desc: "Obtén la insignia verificada" },
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
