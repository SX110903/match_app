"use client"

import { useState, useEffect, useCallback } from "react"
import { apiClient, APIError } from "@/lib/api-client"
import { useAuth } from "@/lib/auth-context"
import { UserAvatar } from "@/components/match-hub/user-avatar"
import { BadgeIcon } from "@/components/match-hub/badge"
import type { BadgeType } from "@/components/match-hub/badge"

interface PublicProfile {
  id: string
  name: string
  age: number
  bio?: string
  occupation?: string
  location?: string
  photos: { id: string; url: string }[]
  interests: string[]
  badge: string
  follower_count: number
}

interface Props {
  userId: string
  onClose?: () => void
}

export function PublicProfileView({ userId, onClose }: Props) {
  const { user } = useAuth()
  const [profile, setProfile] = useState<PublicProfile | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [currentImageIndex, setCurrentImageIndex] = useState(0)
  const [isFollowing, setIsFollowing] = useState(false)
  const [followLoading, setFollowLoading] = useState(false)

  const fetchProfile = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await apiClient<PublicProfile>(`/api/v1/users/${userId}/profile`)
      setProfile(data)
    } catch (err) {
      if (err instanceof APIError && err.status === 404) {
        setError("Profile not found.")
      } else {
        setError("Failed to load profile.")
      }
    } finally {
      setLoading(false)
    }
  }, [userId])

  const fetchFollowStatus = useCallback(async () => {
    if (!user) return
    try {
      const following = await apiClient<string[]>(`/api/v1/users/${user.id}/following?limit=1000&page=1`)
      setIsFollowing(Array.isArray(following) && following.includes(userId))
    } catch {
      // silently ignore — follow status is best-effort
    }
  }, [user, userId])

  useEffect(() => {
    fetchProfile()
  }, [fetchProfile])

  useEffect(() => {
    fetchFollowStatus()
  }, [fetchFollowStatus])

  const handleFollow = async () => {
    if (followLoading) return
    setFollowLoading(true)
    try {
      if (isFollowing) {
        await apiClient(`/api/v1/users/me/follow/${userId}`, { method: "DELETE" })
        setIsFollowing(false)
        if (profile) setProfile({ ...profile, follower_count: Math.max(0, profile.follower_count - 1) })
      } else {
        await apiClient(`/api/v1/users/me/follow/${userId}`, { method: "POST" })
        setIsFollowing(true)
        if (profile) setProfile({ ...profile, follower_count: profile.follower_count + 1 })
      }
    } catch {
      // silently ignore follow errors
    } finally {
      setFollowLoading(false)
    }
  }

  const prevImage = () => {
    if (!profile || profile.photos.length === 0) return
    setCurrentImageIndex((i) => (i - 1 + profile.photos.length) % profile.photos.length)
  }

  const nextImage = () => {
    if (!profile || profile.photos.length === 0) return
    setCurrentImageIndex((i) => (i + 1) % profile.photos.length)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  if (error || !profile) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-4 p-6">
        {onClose && (
          <button onClick={onClose} className="self-start text-sm text-muted-foreground">
            ← Back
          </button>
        )}
        <p className="text-muted-foreground">{error ?? "Profile not found."}</p>
      </div>
    )
  }

  const mainPhoto = profile.photos[currentImageIndex]?.url

  return (
    <div className="flex flex-col h-full overflow-y-auto">
      {/* Header */}
      <div className="flex items-center justify-between p-4">
        {onClose ? (
          <button onClick={onClose} className="text-sm text-muted-foreground hover:text-foreground transition-colors">
            ← Back
          </button>
        ) : (
          <div />
        )}
        {user && user.id !== userId && (
          <button
            onClick={handleFollow}
            disabled={followLoading}
            className={`px-4 py-1.5 rounded-full text-sm font-medium transition-colors ${
              isFollowing
                ? "bg-muted text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
                : "bg-primary text-primary-foreground hover:bg-primary/90"
            }`}
          >
            {followLoading ? "..." : isFollowing ? "Unfollow" : "Follow"}
          </button>
        )}
      </div>

      {/* Photo carousel */}
      <div className="relative w-full aspect-square bg-muted flex-shrink-0">
        {mainPhoto ? (
          <UserAvatar
            src={mainPhoto}
            alt={profile.name}
            className="w-full h-full object-cover"
            fallbackId={profile.id}
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-muted-foreground">
            No photos
          </div>
        )}
        {profile.photos.length > 1 && (
          <>
            <button
              onClick={prevImage}
              className="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/40 text-white flex items-center justify-center hover:bg-black/60 transition-colors"
              aria-label="Previous photo"
            >
              ‹
            </button>
            <button
              onClick={nextImage}
              className="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/40 text-white flex items-center justify-center hover:bg-black/60 transition-colors"
              aria-label="Next photo"
            >
              ›
            </button>
            <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1">
              {profile.photos.map((_, idx) => (
                <span
                  key={idx}
                  className={`w-1.5 h-1.5 rounded-full transition-colors ${
                    idx === currentImageIndex ? "bg-white" : "bg-white/50"
                  }`}
                />
              ))}
            </div>
          </>
        )}
      </div>

      {/* Profile info */}
      <div className="p-4 flex flex-col gap-3">
        {/* Name, age, badge */}
        <div className="flex items-center gap-2">
          <h1 className="text-xl font-bold">{profile.name}, {profile.age}</h1>
          <BadgeIcon badge={profile.badge as BadgeType} />
        </div>

        {/* Follower count */}
        <p className="text-sm text-muted-foreground">{profile.follower_count} followers</p>

        {/* Occupation */}
        {profile.occupation && (
          <p className="text-sm text-muted-foreground">{profile.occupation}</p>
        )}

        {/* Location */}
        {profile.location && (
          <p className="text-sm text-muted-foreground">📍 {profile.location}</p>
        )}

        {/* Bio */}
        {profile.bio && (
          <p className="text-sm leading-relaxed">{profile.bio}</p>
        )}

        {/* Interests */}
        {profile.interests && profile.interests.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-1">
            {profile.interests.map((interest) => (
              <span
                key={interest}
                className="px-3 py-1 rounded-full bg-muted text-xs font-medium"
              >
                {interest}
              </span>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
