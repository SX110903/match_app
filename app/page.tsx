"use client"

import { useState, useCallback } from "react"
import { AnimatePresence } from "framer-motion"
import { Header } from "@/components/match-hub/header"
import { BottomNav, TabType } from "@/components/match-hub/bottom-nav"
import { SwipeCard } from "@/components/match-hub/swipe-card"
import { ActionButtons } from "@/components/match-hub/action-buttons"
import { ProfileModal } from "@/components/match-hub/profile-modal"
import { MatchModal } from "@/components/match-hub/match-modal"
import { MatchesList } from "@/components/match-hub/matches-list"
import { ChatView } from "@/components/match-hub/chat-view"
import { ProfileView } from "@/components/match-hub/profile-view"
import { mockProfiles, mockMatches, mockMessages } from "@/lib/mock-data"
import { Profile, Match, Message } from "@/lib/types"

export default function MatchHub() {
  const [activeTab, setActiveTab] = useState<TabType>("discover")
  const [profiles, setProfiles] = useState<Profile[]>(mockProfiles)
  const [currentIndex, setCurrentIndex] = useState(0)
  const [swipeHistory, setSwipeHistory] = useState<{ profile: Profile; direction: "left" | "right" }[]>([])
  const [showProfileModal, setShowProfileModal] = useState(false)
  const [showMatchModal, setShowMatchModal] = useState(false)
  const [matchedProfile, setMatchedProfile] = useState<Profile | null>(null)
  const [matches, setMatches] = useState<Match[]>(mockMatches)
  const [selectedMatch, setSelectedMatch] = useState<Match | null>(null)
  const [messages, setMessages] = useState<Message[]>(mockMessages)

  const currentProfile = profiles[currentIndex]

  const handleSwipe = useCallback((direction: "left" | "right") => {
    if (!currentProfile) return

    setSwipeHistory((prev) => [...prev, { profile: currentProfile, direction }])

    // Simular match aleatorio al dar like
    if (direction === "right" && Math.random() > 0.5) {
      setMatchedProfile(currentProfile)
      setShowMatchModal(true)
      
      // Agregar a matches
      const newMatch: Match = {
        id: `m-${Date.now()}`,
        profile: currentProfile,
        matchedAt: new Date(),
        unread: false,
      }
      setMatches((prev) => [newMatch, ...prev])
    }

    setCurrentIndex((prev) => prev + 1)
  }, [currentProfile])

  const handleLike = () => handleSwipe("right")
  const handleDislike = () => handleSwipe("left")
  
  const handleSuperLike = () => {
    if (!currentProfile) return
    setSwipeHistory((prev) => [...prev, { profile: currentProfile, direction: "right" }])
    setMatchedProfile(currentProfile)
    setShowMatchModal(true)
    
    const newMatch: Match = {
      id: `m-${Date.now()}`,
      profile: currentProfile,
      matchedAt: new Date(),
      unread: false,
    }
    setMatches((prev) => [newMatch, ...prev])
    setCurrentIndex((prev) => prev + 1)
  }

  const handleUndo = () => {
    if (swipeHistory.length === 0) return
    
    const lastSwipe = swipeHistory[swipeHistory.length - 1]
    setSwipeHistory((prev) => prev.slice(0, -1))
    setCurrentIndex((prev) => prev - 1)
  }

  const handleSelectMatch = (match: Match) => {
    setSelectedMatch(match)
  }

  const handleSendMessage = (text: string) => {
    const newMessage: Message = {
      id: `msg-${Date.now()}`,
      senderId: "me",
      text,
      timestamp: new Date(),
      read: true,
    }
    setMessages((prev) => [...prev, newMessage])
  }

  const handleSendMessageFromMatch = () => {
    setShowMatchModal(false)
    if (matchedProfile) {
      const match = matches.find((m) => m.profile.id === matchedProfile.id)
      if (match) {
        setSelectedMatch(match)
        setActiveTab("messages")
      }
    }
  }

  const unreadMessages = matches.filter((m) => m.unread).length
  const newMatches = matches.filter((m) => !m.lastMessage).length

  // Chat view
  if (selectedMatch) {
    return (
      <ChatView
        match={selectedMatch}
        messages={messages}
        onBack={() => setSelectedMatch(null)}
        onSendMessage={handleSendMessage}
      />
    )
  }

  return (
    <div className="flex flex-col h-screen bg-background">
      {/* Header */}
      {activeTab !== "profile" && (
        <Header
          title={activeTab === "discover" ? "Match Hub" : activeTab === "matches" ? "Matches" : "Mensajes"}
          showFilters={activeTab === "discover"}
          showNotifications={true}
          notificationCount={unreadMessages + newMatches}
        />
      )}

      {/* Content */}
      <main className="flex-1 overflow-hidden">
        <AnimatePresence mode="wait">
          {activeTab === "discover" && (
            <div className="flex flex-col h-full">
              {/* Swipe area */}
              <div className="flex-1 relative p-4 overflow-hidden">
                {currentProfile ? (
                  <div className="relative w-full h-full max-w-md mx-auto">
                    <SwipeCard
                      key={currentProfile.id}
                      profile={currentProfile}
                      onSwipe={handleSwipe}
                      onInfoClick={() => setShowProfileModal(true)}
                    />
                  </div>
                ) : (
                  <div className="flex flex-col items-center justify-center h-full text-center p-8">
                    <div className="w-24 h-24 bg-muted rounded-full flex items-center justify-center mb-4">
                      <span className="text-4xl">🔍</span>
                    </div>
                    <h2 className="text-xl font-semibold text-card-foreground mb-2">
                      No hay más perfiles
                    </h2>
                    <p className="text-muted-foreground">
                      Vuelve más tarde para ver nuevas personas cerca de ti
                    </p>
                  </div>
                )}
              </div>

              {/* Action buttons */}
              {currentProfile && (
                <div className="pb-24 pt-4">
                  <ActionButtons
                    onLike={handleLike}
                    onDislike={handleDislike}
                    onSuperLike={handleSuperLike}
                    onUndo={handleUndo}
                    canUndo={swipeHistory.length > 0}
                  />
                </div>
              )}
            </div>
          )}

          {(activeTab === "matches" || activeTab === "messages") && (
            <MatchesList matches={matches} onSelectMatch={handleSelectMatch} />
          )}

          {activeTab === "profile" && <ProfileView />}
        </AnimatePresence>
      </main>

      {/* Bottom navigation */}
      <BottomNav
        activeTab={activeTab}
        onTabChange={setActiveTab}
        unreadMessages={unreadMessages}
        newMatches={newMatches}
      />

      {/* Modals */}
      <ProfileModal
        profile={currentProfile}
        isOpen={showProfileModal}
        onClose={() => setShowProfileModal(false)}
      />

      <MatchModal
        profile={matchedProfile}
        isOpen={showMatchModal}
        onClose={() => setShowMatchModal(false)}
        onSendMessage={handleSendMessageFromMatch}
      />
    </div>
  )
}
