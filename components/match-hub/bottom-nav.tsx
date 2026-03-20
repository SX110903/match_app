"use client"

import { motion } from "framer-motion"
import { Home, Newspaper, Flame, MessageCircle, User } from "lucide-react"

export type TabType = "home" | "news" | "discover" | "messages" | "profile"

interface BottomNavProps {
  activeTab: TabType
  onTabChange: (tab: TabType) => void
  unreadMessages?: number
  newMatches?: number
  userPhoto?: string
  userName?: string
}

export function BottomNav({ activeTab, onTabChange, unreadMessages = 0, newMatches = 0, userPhoto, userName }: BottomNavProps) {
  const tabs = [
    { id: "home" as const, icon: Home, label: "Inicio" },
    { id: "news" as const, icon: Newspaper, label: "Noticias" },
    { id: "discover" as const, icon: Flame, label: "Citas" },
    { id: "messages" as const, icon: MessageCircle, label: "Chat", badge: unreadMessages + newMatches },
    { id: "profile" as const, icon: User, label: "Perfil" },
  ]

  return (
    <nav className="bg-card border-t border-border safe-area-inset-bottom z-50">
      <div className="flex items-center justify-around py-2 px-2">
        {tabs.map((tab) => {
          const Icon = tab.icon
          const isActive = activeTab === tab.id
          const isProfile = tab.id === "profile"

          return (
            <button
              key={tab.id}
              onClick={() => onTabChange(tab.id)}
              className="relative flex flex-col items-center gap-1 py-2 px-3 transition-colors flex-1"
              aria-label={tab.label}
              aria-current={isActive ? "page" : undefined}
            >
              <div className="relative">
                <motion.div
                  animate={{ scale: isActive ? 1.1 : 1 }}
                  transition={{ type: "spring", stiffness: 400, damping: 20 }}
                >
                  {isProfile && userPhoto ? (
                    <div className={`w-5 h-5 rounded-full overflow-hidden ring-1 ${isActive ? "ring-primary" : "ring-muted-foreground/40"}`}>
                      <img
                        src={userPhoto}
                        alt={userName ?? "Perfil"}
                        className="w-full h-full object-cover"
                        referrerPolicy="no-referrer"
                      />
                    </div>
                  ) : isProfile && userName ? (
                    <div className={`w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold ring-1 ${isActive ? "bg-primary text-primary-foreground ring-primary" : "bg-muted text-muted-foreground ring-muted-foreground/40"}`}>
                      {userName.charAt(0).toUpperCase()}
                    </div>
                  ) : (
                    <Icon
                      className={`w-5 h-5 transition-colors ${
                        isActive ? "text-primary" : "text-muted-foreground"
                      } ${isActive && tab.id === "discover" ? "fill-primary" : ""}`}
                    />
                  )}
                </motion.div>
                {tab.badge && tab.badge > 0 && (
                  <span className="absolute -top-1 -right-1 w-4 h-4 bg-primary text-primary-foreground text-[10px] font-bold rounded-full flex items-center justify-center">
                    {tab.badge > 9 ? "9+" : tab.badge}
                  </span>
                )}
              </div>
              <span
                className={`text-[10px] transition-colors ${
                  isActive ? "text-primary font-medium" : "text-muted-foreground"
                }`}
              >
                {tab.label}
              </span>
              {isActive && (
                <motion.div
                  layoutId="activeTab"
                  className="absolute -bottom-2 w-1 h-1 bg-primary rounded-full"
                />
              )}
            </button>
          )
        })}
      </div>
    </nav>
  )
}
