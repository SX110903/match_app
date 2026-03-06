"use client"

import { motion } from "framer-motion"
import { Flame, Bell, Sliders } from "lucide-react"
import { Button } from "@/components/ui/button"

interface HeaderProps {
  title?: string
  showFilters?: boolean
  showNotifications?: boolean
  notificationCount?: number
}

export function Header({
  title = "Match Hub",
  showFilters = true,
  showNotifications = true,
  notificationCount = 0,
}: HeaderProps) {
  return (
    <header className="flex items-center justify-between px-4 py-3 bg-card border-b border-border">
      {showFilters ? (
        <Button variant="ghost" size="icon" aria-label="Filtros">
          <Sliders className="w-6 h-6 text-muted-foreground" />
        </Button>
      ) : (
        <div className="w-10" />
      )}

      <motion.div
        initial={{ y: -10, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        className="flex items-center gap-2"
      >
        <Flame className="w-8 h-8 text-primary fill-primary" />
        <h1 className="text-xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
          {title}
        </h1>
      </motion.div>

      {showNotifications ? (
        <Button variant="ghost" size="icon" className="relative" aria-label="Notificaciones">
          <Bell className="w-6 h-6 text-muted-foreground" />
          {notificationCount > 0 && (
            <span className="absolute top-1 right-1 w-4 h-4 bg-primary text-primary-foreground text-[10px] font-bold rounded-full flex items-center justify-center">
              {notificationCount > 9 ? "9+" : notificationCount}
            </span>
          )}
        </Button>
      ) : (
        <div className="w-10" />
      )}
    </header>
  )
}
