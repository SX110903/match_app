"use client"

import { useState, useEffect, useCallback } from "react"
import { Heart, MessageCircle, Repeat2, Share2, Image, X, Send } from "lucide-react"
import { apiClient } from "@/lib/api-client"
import { useAuth } from "@/lib/auth-context"
import { AVATAR_BASE } from "@/lib/constants"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { AdBanner } from "@/components/match-hub/ad-banner"

interface Post {
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

interface Comment {
  id: string
  user_id: string
  content: string
  author_name: string
  author_avatar: string
  created_at: string
}

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return "ahora"
  if (mins < 60) return `${mins}m`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h`
  return `${Math.floor(hrs / 24)}d`
}

function TweetCard({ post, onUpdated }: {
  post: Post
  onUpdated: (updated: Post) => void
}) {
  const [liked, setLiked] = useState(post.is_liked_by_me)
  const [likesCount, setLikesCount] = useState(post.likes_count)
  const [showReplies, setShowReplies] = useState(false)
  const [comments, setComments] = useState<Comment[]>([])
  const [commentText, setCommentText] = useState("")
  const [loadingComments, setLoadingComments] = useState(false)

  const toggleLike = async () => {
    if (liked) {
      setLiked(false)
      setLikesCount((n) => n - 1)
      await apiClient(`/api/v1/posts/${post.id}/like`, { method: "DELETE" }).catch(() => {})
    } else {
      setLiked(true)
      setLikesCount((n) => n + 1)
      await apiClient(`/api/v1/posts/${post.id}/like`, { method: "POST" }).catch(() => {})
    }
  }

  const loadComments = async () => {
    setLoadingComments(true)
    try {
      const data = await apiClient<Comment[]>(`/api/v1/posts/${post.id}/comments`)
      setComments(data ?? [])
    } catch {
      //
    } finally {
      setLoadingComments(false)
    }
  }

  const handleToggleReplies = () => {
    if (!showReplies && comments.length === 0) loadComments()
    setShowReplies((v) => !v)
  }

  const addComment = async () => {
    if (!commentText.trim()) return
    try {
      const comment = await apiClient<Comment>(`/api/v1/posts/${post.id}/comments`, {
        method: "POST",
        body: { content: commentText.trim() },
      })
      setComments((prev) => [...prev, comment])
      setCommentText("")
    } catch {
      //
    }
  }

  const avatar = post.author_avatar || `${AVATAR_BASE}?u=${post.user_id}`

  return (
    <article className="flex gap-3 px-4 py-3 border-b border-border hover:bg-muted/20 transition-colors">
      {/* Avatar */}
      <div className="flex-shrink-0">
        <Avatar className="w-10 h-10">
          <AvatarImage src={avatar} referrerPolicy="no-referrer" />
          <AvatarFallback className="bg-primary text-primary-foreground text-sm font-bold">
            {post.author_name?.[0]?.toUpperCase() ?? "?"}
          </AvatarFallback>
        </Avatar>
      </div>

      {/* Right column */}
      <div className="flex-1 min-w-0">
        {/* Header row */}
        <div className="flex items-center gap-1.5 mb-0.5">
          <span className="font-bold text-sm text-foreground truncate">{post.author_name}</span>
          <span className="text-muted-foreground text-sm">·</span>
          <span className="text-muted-foreground text-sm flex-shrink-0">{timeAgo(post.created_at)}</span>
        </div>

        {/* Tweet text */}
        <p className="text-sm text-foreground leading-relaxed mb-2 whitespace-pre-wrap">{post.content}</p>

        {/* Image */}
        {post.image_url && (
          <div className="rounded-2xl overflow-hidden mb-3 border border-border/50">
            <img
              src={post.image_url}
              alt="Imagen"
              className="w-full max-h-72 object-cover"
              referrerPolicy="no-referrer"
            />
          </div>
        )}

        {/* Action bar — Twitter style */}
        <div className="flex items-center gap-0 -ml-2 mt-1">
          <button
            onClick={handleToggleReplies}
            className="flex items-center gap-1.5 text-muted-foreground hover:text-primary transition-colors px-2 py-1.5 rounded-full hover:bg-primary/10 text-sm"
          >
            <MessageCircle className="w-4 h-4" />
            {comments.length > 0 && <span className="text-xs">{comments.length}</span>}
          </button>
          <button className="flex items-center gap-1.5 text-muted-foreground hover:text-green-500 transition-colors px-2 py-1.5 rounded-full hover:bg-green-500/10 text-sm">
            <Repeat2 className="w-4 h-4" />
          </button>
          <button
            onClick={toggleLike}
            className={`flex items-center gap-1.5 transition-colors px-2 py-1.5 rounded-full text-sm ${
              liked
                ? "text-primary hover:bg-primary/10"
                : "text-muted-foreground hover:text-primary hover:bg-primary/10"
            }`}
          >
            <Heart className={`w-4 h-4 ${liked ? "fill-primary" : ""}`} />
            {likesCount > 0 && <span className="text-xs">{likesCount}</span>}
          </button>
          <button className="flex items-center gap-1.5 text-muted-foreground hover:text-primary transition-colors px-2 py-1.5 rounded-full hover:bg-primary/10 text-sm ml-auto">
            <Share2 className="w-4 h-4" />
          </button>
        </div>

        {/* Replies */}
        {showReplies && (
          <div className="mt-3 space-y-2 border-l-2 border-border/50 pl-3">
            {loadingComments ? (
              <p className="text-xs text-muted-foreground py-1">Cargando...</p>
            ) : comments.length === 0 ? (
              <p className="text-xs text-muted-foreground py-1">Sin respuestas aún</p>
            ) : (
              comments.map((c) => (
                <div key={c.id} className="flex gap-2">
                  <Avatar className="w-7 h-7 flex-shrink-0">
                    <AvatarImage
                      src={c.author_avatar || `${AVATAR_BASE}?u=${c.user_id}`}
                      referrerPolicy="no-referrer"
                    />
                    <AvatarFallback className="text-[10px] bg-secondary">
                      {c.author_name?.[0] ?? "?"}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 bg-secondary/60 rounded-2xl px-3 py-1.5">
                    <span className="font-semibold text-xs text-foreground">{c.author_name} </span>
                    <span className="text-xs text-foreground/80">{c.content}</span>
                  </div>
                </div>
              ))
            )}
            {/* Reply input */}
            <div className="flex gap-2 pt-1">
              <input
                value={commentText}
                onChange={(e) => setCommentText(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && addComment()}
                placeholder="Responder..."
                className="flex-1 bg-secondary/50 border border-border rounded-full px-3 py-1.5 text-xs focus:outline-none focus:ring-1 focus:ring-primary"
              />
              <button
                onClick={addComment}
                disabled={!commentText.trim()}
                className="text-primary disabled:opacity-40"
              >
                <Send className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}
      </div>
    </article>
  )
}

interface AdData {
  id: string
  title: string
  description?: string
  image_url?: string
  cta_text: string
  cta_url: string
}

export function InicioView() {
  const { user } = useAuth()
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [draftText, setDraftText] = useState("")
  const [draftImage, setDraftImage] = useState("")
  const [showImageInput, setShowImageInput] = useState(false)
  const [posting, setPosting] = useState(false)
  const [composing, setComposing] = useState(false)
  const [activeAd, setActiveAd] = useState<AdData | null>(null)

  const loadFeed = useCallback(async () => {
    setLoading(true)
    try {
      const [data, ad] = await Promise.all([
        apiClient<Post[]>("/api/v1/posts?limit=30"),
        apiClient<AdData | null>(`/api/v1/ads/active?badge=${user?.badge ?? "none"}`).catch(() => null),
      ])
      setPosts(data ?? [])
      setActiveAd(ad ?? null)
    } catch {
      //
    } finally {
      setLoading(false)
    }
  }, [user?.badge])

  useEffect(() => {
    loadFeed()
  }, [loadFeed])

  const submitPost = async () => {
    if (!draftText.trim()) return
    setPosting(true)
    try {
      const post = await apiClient<Post>("/api/v1/posts", {
        method: "POST",
        body: { content: draftText.trim(), image_url: draftImage.trim() || undefined },
      })
      setPosts((prev) => [post, ...prev])
      setDraftText("")
      setDraftImage("")
      setShowImageInput(false)
      setComposing(false)
    } catch {
      //
    } finally {
      setPosting(false)
    }
  }

  if (!user) return null

  const userAvatar = user.photos?.[0]?.url || `${AVATAR_BASE}?u=${user.id}`

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header — Twitter style */}
      <div className="sticky top-0 z-10 bg-background/90 backdrop-blur-sm border-b border-border">
        <div className="flex items-center justify-center px-4 py-3">
          <h1 className="text-lg font-bold text-foreground">Match Hub</h1>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto pb-24">
        {/* Compose box */}
        <div className="flex gap-3 px-4 py-3 border-b border-border">
          <Avatar className="w-10 h-10 flex-shrink-0 mt-0.5">
            <AvatarImage src={userAvatar} referrerPolicy="no-referrer" />
            <AvatarFallback className="bg-primary text-primary-foreground font-bold">
              {user.name?.[0]?.toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div className="flex-1 min-w-0">
            <textarea
              placeholder="¿Qué está pasando?"
              value={draftText}
              onChange={(e) => setDraftText(e.target.value)}
              onFocus={() => setComposing(true)}
              rows={composing ? 3 : 1}
              className="w-full resize-none bg-transparent text-foreground placeholder:text-muted-foreground focus:outline-none text-base py-1 leading-relaxed"
            />
            {showImageInput && (
              <div className="flex items-center gap-2 mt-1 mb-2">
                <input
                  placeholder="URL de imagen (opcional)..."
                  value={draftImage}
                  onChange={(e) => setDraftImage(e.target.value)}
                  className="flex-1 text-sm bg-transparent border-b border-border focus:outline-none focus:border-primary py-0.5 text-foreground placeholder:text-muted-foreground"
                />
                <button
                  onClick={() => { setShowImageInput(false); setDraftImage("") }}
                  className="text-muted-foreground hover:text-foreground"
                >
                  <X className="w-3.5 h-3.5" />
                </button>
              </div>
            )}
            <div className="flex items-center justify-between pt-2 border-t border-border/40">
              <button
                onClick={() => setShowImageInput(!showImageInput)}
                className="text-primary hover:text-primary/70 p-1 rounded-full hover:bg-primary/10 transition-colors"
                title="Añadir imagen"
              >
                <Image className="w-5 h-5" />
              </button>
              <button
                onClick={submitPost}
                disabled={posting || !draftText.trim()}
                className="bg-primary hover:bg-primary/90 text-primary-foreground text-sm font-bold px-5 py-1.5 rounded-full disabled:opacity-50 transition-colors"
              >
                {posting ? "..." : "Publicar"}
              </button>
            </div>
          </div>
        </div>

        {/* Feed */}
        {loading ? (
          <div className="flex items-center justify-center py-16">
            <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
          </div>
        ) : posts.length === 0 ? (
          <div className="text-center py-20">
            <p className="text-muted-foreground font-medium">No hay publicaciones aún</p>
            <p className="text-sm text-muted-foreground/60 mt-1">¡Sé el primero en publicar!</p>
          </div>
        ) : (
          posts.map((post, idx) => (
            <div key={post.id}>
              <TweetCard
                post={post}
                onUpdated={(updated) =>
                  setPosts((prev) => prev.map((p) => (p.id === updated.id ? updated : p)))
                }
              />
              {activeAd && (idx + 1) % 5 === 0 && <AdBanner ad={activeAd} />}
            </div>
          ))
        )}
      </div>
    </div>
  )
}
