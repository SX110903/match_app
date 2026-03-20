"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Heart } from "lucide-react"
import { apiClient, APIError } from "@/lib/api-client"
import { BadgeIcon, BadgeType } from "@/components/match-hub/badge"
import { UserAvatar } from "@/components/match-hub/user-avatar"

interface ExploreUser {
  id: string
  name: string
  age: number
  avatar: string
  badge: string
  vip_level: number
  follower_count: number
  is_following: boolean
}

interface ExplorePost {
  id: string
  user_id: string
  content: string
  image_url?: string
  likes_count: number
  author_name: string
  author_avatar: string
  is_liked_by_me: boolean
  created_at: string
}

export function ExploreView() {
  const router = useRouter()
  const [users, setUsers] = useState<ExploreUser[]>([])
  const [posts, setPosts] = useState<ExplorePost[]>([])
  const [loading, setLoading] = useState(true)
  const [cursor] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    async function load() {
      setLoading(true)
      try {
        const [fetchedUsers, fetchedPosts] = await Promise.all([
          apiClient<ExploreUser[]>("/api/v1/explore/users?limit=20"),
          apiClient<ExplorePost[]>("/api/v1/explore/posts?limit=20"),
        ])
        if (!cancelled) {
          setUsers(fetchedUsers ?? [])
          setPosts(fetchedPosts ?? [])
        }
      } catch (err) {
        if (err instanceof APIError) {
          console.error("Explore fetch failed:", err.message)
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    load()
    return () => { cancelled = true }
  }, [])

  async function handleFollowToggle(user: ExploreUser) {
    const method = user.is_following ? "DELETE" : "POST"
    try {
      await apiClient(`/api/v1/users/me/follow/${user.id}`, { method })
      setUsers(prev =>
        prev.map(u =>
          u.id === user.id
            ? {
                ...u,
                is_following: !u.is_following,
                follower_count: u.follower_count + (u.is_following ? -1 : 1),
              }
            : u
        )
      )
    } catch (err) {
      if (err instanceof APIError) {
        console.error("Follow toggle failed:", err.message)
      }
    }
  }

  if (loading) {
    return (
      <div className="p-4 space-y-6">
        {/* Users skeleton */}
        <div>
          <div className="h-6 w-32 bg-gray-200 rounded mb-3 animate-pulse" />
          <div className="flex gap-3 overflow-x-auto pb-2">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="flex-shrink-0 w-24 flex flex-col items-center gap-2">
                <div className="w-16 h-16 rounded-full bg-gray-200 animate-pulse" />
                <div className="h-3 w-16 bg-gray-200 rounded animate-pulse" />
                <div className="h-6 w-16 bg-gray-200 rounded animate-pulse" />
              </div>
            ))}
          </div>
        </div>
        {/* Posts skeleton */}
        <div>
          <div className="h-6 w-40 bg-gray-200 rounded mb-3 animate-pulse" />
          <div className="grid grid-cols-2 gap-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="rounded-xl bg-gray-100 p-3 space-y-2 animate-pulse">
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 rounded-full bg-gray-200" />
                  <div className="h-3 w-20 bg-gray-200 rounded" />
                </div>
                <div className="h-3 w-full bg-gray-200 rounded" />
                <div className="h-3 w-3/4 bg-gray-200 rounded" />
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="p-4 space-y-6">
      {/* Section A: Users */}
      <section>
        <h2 className="text-base font-semibold text-gray-800 mb-3">Descubrir personas</h2>
        {users.length === 0 ? (
          <p className="text-sm text-gray-500">No hay usuarios para explorar.</p>
        ) : (
          <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide">
            {users.map(user => (
              <div
                key={user.id}
                className="flex-shrink-0 w-24 flex flex-col items-center gap-1"
              >
                <button
                  className="focus:outline-none"
                  onClick={() => router.push(`/profile/${user.id}`)}
                  aria-label={`Ver perfil de ${user.name}`}
                >
                  <UserAvatar
                    src={user.avatar || undefined}
                    alt={user.name}
                    fallbackId={user.id}
                    className="w-16 h-16 rounded-full object-cover border-2 border-white shadow"
                  />
                </button>
                <div className="flex items-center gap-1 max-w-full">
                  <span className="text-xs font-medium text-gray-700 truncate max-w-[72px]">
                    {user.name}
                  </span>
                  <BadgeIcon badge={user.badge as BadgeType} />
                </div>
                <button
                  onClick={() => handleFollowToggle(user)}
                  className={`text-xs px-2 py-1 rounded-full font-medium transition-colors ${
                    user.is_following
                      ? "bg-gray-200 text-gray-700 hover:bg-gray-300"
                      : "bg-pink-500 text-white hover:bg-pink-600"
                  }`}
                >
                  {user.is_following ? "Siguiendo" : "Seguir"}
                </button>
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Section B: Trending Posts */}
      <section>
        <h2 className="text-base font-semibold text-gray-800 mb-3">Publicaciones en tendencia</h2>
        {posts.length === 0 ? (
          <p className="text-sm text-gray-500">No hay publicaciones en tendencia.</p>
        ) : (
          <div className="grid grid-cols-2 gap-3">
            {posts.map(post => (
              <div
                key={post.id}
                className="rounded-xl bg-white border border-gray-100 shadow-sm p-3 flex flex-col gap-2"
              >
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => router.push(`/profile/${post.user_id}`)}
                    aria-label={`Ver perfil de ${post.author_name}`}
                    className="focus:outline-none"
                  >
                    <UserAvatar
                      src={post.author_avatar || undefined}
                      alt={post.author_name}
                      fallbackId={post.user_id}
                      className="w-8 h-8 rounded-full object-cover"
                    />
                  </button>
                  <span className="text-xs font-semibold text-gray-700 truncate">
                    {post.author_name}
                  </span>
                </div>
                <p className="text-xs text-gray-600 line-clamp-2 leading-relaxed">
                  {post.content}
                </p>
                <div className="flex items-center gap-1 mt-auto">
                  <Heart
                    className={`w-3.5 h-3.5 ${post.is_liked_by_me ? "fill-pink-500 text-pink-500" : "text-gray-400"}`}
                  />
                  <span className="text-xs text-gray-500">{post.likes_count}</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  )
}
