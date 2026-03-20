"use client"

import { motion } from "framer-motion"
import { Home, Newspaper, Flame, MessageCircle, User, Heart } from "lucide-react"
import { TabType } from "@/components/match-hub/bottom-nav"

interface SidebarNavProps {
  activeTab: TabType
  onTabChange: (tab: TabType) => void
  unreadMessages?: number
  newMatches?: number
  userPhoto?: string
  userName?: string
}

export function SidebarNav({
  activeTab,
  onTabChange,
  unreadMessages = 0,
  newMatches = 0,
  userPhoto,
  userName,
}: SidebarNavProps) {
  const tabs = [
    { id: "home" as const, icon: Home, label: "Inicio" },
    { id: "news" as const, icon: Newspaper, label: "Noticias" },
    { id: "discover" as const, icon: Flame, label: "Citas" },
    { id: "messages" as const, icon: MessageCircle, label: "Mensajes", badge: unreadMessages + newMatches },
    { id: "profile" as const, icon: User, label: "Perfil" },
  ]

  return (
    <aside className="hidden md:flex flex-col w-56 h-screen bg-card border-r border-border shrink-0">
      {/* Logo */}
      <div className="flex items-center gap-3 px-5 py-5 border-b border-border">
        <div className="w-9 h-9 bg-primary rounded-xl flex items-center justify-center shrink-0">
          <Heart className="w-5 h-5 text-primary-foreground fill-current" />
        </div>
        <span className="text-lg font-bold text-card-foreground">Match Hub</span>
      </div>

      {/* Nav items */}
      <nav className="flex flex-col gap-1 p-3 flex-1 overflow-y-auto">
        {tabs.map((tab) => {
          const Icon = tab.icon
          const isActive = activeTab === tab.id
          const isProfile = tab.id === "profile"

          return (
            <button
              key={tab.id}
              onClick={() => onTabChange(tab.id)}
              className={`relative flex items-center gap-3 px-4 py-3 rounded-xl transition-colors text-left ${
                isActive
                  ? "bg-primary/10 text-primary"
                  : "text-muted-foreground hover:bg-muted/50 hover:text-card-foreground"
              }`}
              aria-label={tab.label}
              aria-current={isActive ? "page" : undefined}
            >
              {isActive && (
                <motion.div
                  layoutId="sidebarActiveIndicator"
                  className="absolute left-0 top-1/2 -translate-y-1/2 w-[3px] h-6 bg-primary rounded-r-full"
                />
              )}
              <div className="relative shrink-0">
                <motion.div
                  animate={{ scale: isActive ? 1.1 : 1 }}
                  transition={{ type: "spring", stiffness: 400, damping: 20 }}
                >
                  {isProfile && userPhoto ? (
                    <div
                      className={`w-5 h-5 rounded-full overflow-hidden ring-1 ${
                        isActive ? "ring-primary" : "ring-muted-foreground/40"
                      }`}
                    >
                      <img
                        src={userPhoto}
                        alt={userName ?? "Perfil"}
                        className="w-full h-full object-cover"
                        referrerPolicy="no-referrer"
                      />
                    </div>
                  ) : isProfile && userName ? (
                    <div
                      className={`w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold ring-1 ${
                        isActive
                          ? "bg-primary text-primary-foreground ring-primary"
                          : "bg-muted text-muted-foreground ring-muted-foreground/40"
                      }`}
                    >
                      {userName.charAt(0).toUpperCase()}
                    </div>
                  ) : (
                    <Icon
                      className={`w-5 h-5 ${
                        isActive && tab.id === "discover" ? "fill-primary" : ""
                      }`}
                    />
                  )}
                </motion.div>
                {tab.badge && tab.badge > 0 ? (
                  <span className="absolute -top-1 -right-1 w-4 h-4 bg-primary text-primary-foreground text-[10px] font-bold rounded-full flex items-center justify-center">
                    {tab.badge > 9 ? "9+" : tab.badge}
                  </span>
                ) : null}
              </div>
              <span className={`text-sm font-medium ${isActive ? "text-primary" : ""}`}>
                {tab.label}
              </span>
            </button>
          )
        })}
      </nav>
    </aside>
  )
}
