"use client"

import { motion } from "framer-motion"
import { Flame, Heart, MessageCircle, User } from "lucide-react"

export type TabType = "discover" | "matches" | "messages" | "profile"

interface BottomNavProps {
  activeTab: TabType
  onTabChange: (tab: TabType) => void
  unreadMessages?: number
  newMatches?: number
}

export function BottomNav({ activeTab, onTabChange, unreadMessages = 0, newMatches = 0 }: BottomNavProps) {
  const tabs = [
    { id: "discover" as const, icon: Flame, label: "Descubrir" },
    { id: "matches" as const, icon: Heart, label: "Matches", badge: newMatches },
    { id: "messages" as const, icon: MessageCircle, label: "Chat", badge: unreadMessages },
    { id: "profile" as const, icon: User, label: "Perfil" },
  ]

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-card border-t border-border safe-area-inset-bottom">
      <div className="flex items-center justify-around py-2 px-4 max-w-md mx-auto">
        {tabs.map((tab) => {
          const Icon = tab.icon
          const isActive = activeTab === tab.id

          return (
            <button
              key={tab.id}
              onClick={() => onTabChange(tab.id)}
              className="relative flex flex-col items-center gap-1 py-2 px-4 transition-colors"
              aria-label={tab.label}
              aria-current={isActive ? "page" : undefined}
            >
              <div className="relative">
                <motion.div
                  animate={{ scale: isActive ? 1.1 : 1 }}
                  transition={{ type: "spring", stiffness: 400, damping: 20 }}
                >
                  <Icon
                    className={`w-6 h-6 transition-colors ${
                      isActive ? "text-primary" : "text-muted-foreground"
                    } ${isActive && tab.id === "discover" ? "fill-primary" : ""}`}
                  />
                </motion.div>
                {tab.badge && tab.badge > 0 && (
                  <span className="absolute -top-1 -right-1 w-4 h-4 bg-primary text-primary-foreground text-[10px] font-bold rounded-full flex items-center justify-center">
                    {tab.badge > 9 ? "9+" : tab.badge}
                  </span>
                )}
              </div>
              <span
                className={`text-xs transition-colors ${
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
