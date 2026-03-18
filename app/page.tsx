"use client"

import { useState, useCallback, useEffect } from "react"
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
import { Profile, Match } from "@/lib/types"
import { apiClient } from "@/lib/api-client"
import { useAuth } from "@/lib/auth-context"

// Backend response shapes
interface APICandidate {
  profile: {
    id: string
    user_id: string
    name: string
    age: number
    bio?: string
    occupation?: string
    location?: string
    photos?: string[]
    interests?: string[]
  }
  distance: number
}

interface APIMatch {
  id: string
  user1_id: string
  user2_id: string
  created_at: string
  profile: APICandidate['profile']
  last_message?: string
  unread_count: number
}

function candidateToProfile(c: APICandidate): Profile {
  return {
    id: c.profile.user_id,
    name: c.profile.name,
    age: c.profile.age,
    bio: c.profile.bio ?? '',
    images: c.profile.photos?.length ? c.profile.photos : [`https://i.pravatar.cc/800?u=${c.profile.user_id}`],
    distance: Math.round(c.distance),
    occupation: c.profile.occupation ?? '',
    interests: c.profile.interests ?? [],
  }
}

function apiMatchToMatch(m: APIMatch): Match {
  return {
    id: m.id,
    profile: {
      id: m.profile.user_id,
      name: m.profile.name,
      age: m.profile.age,
      bio: m.profile.bio ?? '',
      images: m.profile.photos?.length ? m.profile.photos : [`https://i.pravatar.cc/800?u=${m.profile.user_id}`],
      distance: 0,
      occupation: m.profile.occupation ?? '',
      interests: m.profile.interests ?? [],
    },
    matchedAt: new Date(m.created_at),
    lastMessage: m.last_message,
    unread: m.unread_count > 0,
  }
}

export default function MatchHub() {
  const { user, isLoading } = useAuth()
  const [activeTab, setActiveTab] = useState<TabType>("discover")
  const [profiles, setProfiles] = useState<Profile[]>([])
  const [currentIndex, setCurrentIndex] = useState(0)
  const [swipeHistory, setSwipeHistory] = useState<{ profile: Profile; direction: "left" | "right" }[]>([])
  const [showProfileModal, setShowProfileModal] = useState(false)
  const [showMatchModal, setShowMatchModal] = useState(false)
  const [matchedProfile, setMatchedProfile] = useState<Profile | null>(null)
  const [matches, setMatches] = useState<Match[]>([])
  const [selectedMatch, setSelectedMatch] = useState<Match | null>(null)
  const [loadingCandidates, setLoadingCandidates] = useState(true)

  const currentProfile = profiles[currentIndex]

  // Load candidates and matches on mount
  useEffect(() => {
    if (!user) return
    async function load() {
      setLoadingCandidates(true)
      try {
        const [candidates, apiMatches] = await Promise.all([
          apiClient<APICandidate[]>('/api/v1/matches/candidates?limit=20'),
          apiClient<APIMatch[]>('/api/v1/matches/'),
        ])
        setProfiles((candidates ?? []).map(candidateToProfile))
        setMatches((apiMatches ?? []).map(apiMatchToMatch))
      } catch {
        // errors handled silently — user will see empty state
      } finally {
        setLoadingCandidates(false)
      }
    }
    load()
  }, [user])

  const handleSwipe = useCallback(async (direction: "left" | "right") => {
    if (!currentProfile) return
    setSwipeHistory((prev) => [...prev, { profile: currentProfile, direction }])
    setCurrentIndex((prev) => prev + 1)

    if (direction === "right") {
      try {
        const result = await apiClient<{ is_match: boolean; match_id?: string }>(
          '/api/v1/matches/swipe',
          { method: 'POST', body: { user_id: currentProfile.id, direction: 'right' } }
        )
        if (result?.is_match && result.match_id) {
          setMatchedProfile(currentProfile)
          setShowMatchModal(true)
          const newMatch: Match = {
            id: result.match_id,
            profile: currentProfile,
            matchedAt: new Date(),
            unread: false,
          }
          setMatches((prev) => [newMatch, ...prev])
        }
      } catch {
        // swipe failed silently
      }
    } else {
      try {
        await apiClient('/api/v1/matches/swipe', {
          method: 'POST',
          body: { user_id: currentProfile.id, direction: 'left' },
        })
      } catch { /* silent */ }
    }
  }, [currentProfile])

  const handleLike = () => handleSwipe("right")
  const handleDislike = () => handleSwipe("left")

  const handleSuperLike = useCallback(async () => {
    if (!currentProfile) return
    setSwipeHistory((prev) => [...prev, { profile: currentProfile, direction: "right" }])
    setCurrentIndex((prev) => prev + 1)
    try {
      const result = await apiClient<{ is_match: boolean; match_id?: string }>(
        '/api/v1/matches/swipe',
        { method: 'POST', body: { user_id: currentProfile.id, direction: 'super' } }
      )
      if (result?.is_match && result.match_id) {
        setMatchedProfile(currentProfile)
        setShowMatchModal(true)
        const newMatch: Match = {
          id: result.match_id,
          profile: currentProfile,
          matchedAt: new Date(),
          unread: false,
        }
        setMatches((prev) => [newMatch, ...prev])
      }
    } catch { /* silent */ }
  }, [currentProfile])

  const handleUndo = () => {
    if (swipeHistory.length === 0) return
    setSwipeHistory((prev) => prev.slice(0, -1))
    setCurrentIndex((prev) => prev - 1)
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
  const newMatchCount = matches.filter((m) => !m.lastMessage).length

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-background">
        <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  // Chat view
  if (selectedMatch) {
    return (
      <ChatView
        match={selectedMatch}
        currentUserId={user?.id ?? ''}
        onBack={() => setSelectedMatch(null)}
      />
    )
  }

  return (
    <div className="flex flex-col h-screen bg-background">
      {activeTab !== "profile" && (
        <Header
          title={activeTab === "discover" ? "Match Hub" : activeTab === "matches" ? "Matches" : "Mensajes"}
          showFilters={activeTab === "discover"}
          showNotifications={true}
          notificationCount={unreadMessages + newMatchCount}
        />
      )}

      <main className="flex-1 overflow-hidden">
        <AnimatePresence mode="wait">
          {activeTab === "discover" && (
            <div className="flex flex-col h-full">
              <div className="flex-1 relative p-4 overflow-hidden">
                {loadingCandidates ? (
                  <div className="flex items-center justify-center h-full">
                    <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
                  </div>
                ) : currentProfile ? (
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
                    <h2 className="text-xl font-semibold text-card-foreground mb-2">No hay más perfiles</h2>
                    <p className="text-muted-foreground">Vuelve más tarde para ver nuevas personas cerca de ti</p>
                  </div>
                )}
              </div>

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
            <MatchesList matches={matches} onSelectMatch={setSelectedMatch} />
          )}

          {activeTab === "profile" && <ProfileView />}
        </AnimatePresence>
      </main>

      <BottomNav
        activeTab={activeTab}
        onTabChange={setActiveTab}
        unreadMessages={unreadMessages}
        newMatches={newMatchCount}
      />

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
