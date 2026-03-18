"use client"

import { useState, useRef, useEffect, useCallback } from "react"
import { motion, AnimatePresence } from "framer-motion"
import { ArrowLeft, Send, MoreVertical, Phone, Video } from "lucide-react"
import { Match } from "@/lib/types"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { apiClient } from "@/lib/api-client"
import { addWSListener, removeWSListener, type WSMessage } from "@/lib/websocket-client"

interface BackendMessage {
  id: string
  match_id: string
  sender_id: string
  text: string
  read_at: string | null
  created_at: string
}

interface DisplayMessage {
  id: string
  senderId: string
  text: string
  timestamp: Date
  read: boolean
}

interface ChatViewProps {
  match: Match
  currentUserId: string
  onBack: () => void
}

function toDisplay(m: BackendMessage): DisplayMessage {
  return {
    id: m.id,
    senderId: m.sender_id,
    text: m.text,
    timestamp: new Date(m.created_at),
    read: m.read_at !== null,
  }
}

export function ChatView({ match, currentUserId, onBack }: ChatViewProps) {
  const [messages, setMessages] = useState<DisplayMessage[]>([])
  const [newMessage, setNewMessage] = useState("")
  const [sending, setSending] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  // Fetch message history on mount
  useEffect(() => {
    apiClient<BackendMessage[]>(`/api/v1/matches/${match.id}/messages?limit=50`)
      .then((data) => setMessages((data ?? []).map(toDisplay)))
      .catch(() => {/* silent */})

    // Mark as read
    apiClient(`/api/v1/matches/${match.id}/messages/read`, { method: 'PUT' }).catch(() => {})
  }, [match.id])

  // WebSocket listener for incoming messages
  useEffect(() => {
    const handler = (wsMsg: WSMessage) => {
      if (wsMsg.type === 'message' && wsMsg.match_id === match.id && wsMsg.message) {
        const incoming = toDisplay(wsMsg.message)
        setMessages((prev) => {
          // Deduplicate by id
          if (prev.some((m) => m.id === incoming.id)) return prev
          return [...prev, incoming]
        })
      }
    }
    addWSListener(handler)
    return () => removeWSListener(handler)
  }, [match.id])

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  const handleSend = useCallback(async () => {
    const text = newMessage.trim()
    if (!text || sending) return
    setNewMessage("")
    setSending(true)

    // Optimistic update
    const optimistic: DisplayMessage = {
      id: `opt-${Date.now()}`,
      senderId: currentUserId,
      text,
      timestamp: new Date(),
      read: true,
    }
    setMessages((prev) => [...prev, optimistic])

    try {
      const created = await apiClient<BackendMessage>(`/api/v1/matches/${match.id}/messages`, {
        method: 'POST',
        body: { text },
      })
      // Replace optimistic with real message
      setMessages((prev) => prev.map((m) => m.id === optimistic.id ? toDisplay(created) : m))
    } catch {
      // Remove optimistic on failure
      setMessages((prev) => prev.filter((m) => m.id !== optimistic.id))
      setNewMessage(text) // restore input
    } finally {
      setSending(false)
    }
  }, [newMessage, sending, match.id, currentUserId])

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const formatTime = (date: Date) =>
    date.toLocaleTimeString("es-ES", { hour: "2-digit", minute: "2-digit" })

  return (
    <div className="flex flex-col h-full bg-background">
      <header className="flex items-center gap-3 p-4 bg-card border-b border-border">
        <Button variant="ghost" size="icon" onClick={onBack} aria-label="Volver">
          <ArrowLeft className="w-6 h-6" />
        </Button>
        <div className="flex items-center gap-3 flex-1">
          <div className="w-10 h-10 rounded-full overflow-hidden">
            <img
              src={match.profile.images[0]}
              alt={match.profile.name}
              className="w-full h-full object-cover"
              crossOrigin="anonymous"
            />
          </div>
          <div>
            <h2 className="font-semibold text-card-foreground">{match.profile.name}</h2>
            <p className="text-xs text-muted-foreground">Match Hub</p>
          </div>
        </div>
        <div className="flex items-center gap-1">
          <Button variant="ghost" size="icon" aria-label="Llamar">
            <Phone className="w-5 h-5 text-muted-foreground" />
          </Button>
          <Button variant="ghost" size="icon" aria-label="Videollamada">
            <Video className="w-5 h-5 text-muted-foreground" />
          </Button>
          <Button variant="ghost" size="icon" aria-label="Más opciones">
            <MoreVertical className="w-5 h-5 text-muted-foreground" />
          </Button>
        </div>
      </header>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        <div className="text-center py-4">
          <div className="w-20 h-20 rounded-full overflow-hidden mx-auto mb-3 ring-2 ring-primary ring-offset-2 ring-offset-background">
            <img
              src={match.profile.images[0]}
              alt={match.profile.name}
              className="w-full h-full object-cover"
              crossOrigin="anonymous"
            />
          </div>
          <h3 className="font-semibold text-card-foreground">{match.profile.name}</h3>
          <p className="text-sm text-muted-foreground">{match.profile.occupation}</p>
          <p className="text-xs text-muted-foreground mt-2">¡Hicieron match! Empieza la conversación.</p>
        </div>

        <AnimatePresence>
          {messages.map((message, index) => {
            const isMe = message.senderId === currentUserId
            return (
              <motion.div
                key={message.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: Math.min(index * 0.03, 0.3) }}
                className={`flex ${isMe ? "justify-end" : "justify-start"}`}
              >
                <div
                  className={`max-w-[75%] px-4 py-2 rounded-2xl ${
                    isMe
                      ? "bg-primary text-primary-foreground rounded-br-md"
                      : "bg-muted text-card-foreground rounded-bl-md"
                  }`}
                >
                  <p className="text-sm">{message.text}</p>
                  <p className={`text-[10px] mt-1 ${isMe ? "text-primary-foreground/70" : "text-muted-foreground"}`}>
                    {formatTime(message.timestamp)}
                  </p>
                </div>
              </motion.div>
            )
          })}
        </AnimatePresence>
        <div ref={messagesEndRef} />
      </div>

      <div className="p-4 bg-card border-t border-border">
        <div className="flex items-center gap-2">
          <Input
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Escribe un mensaje..."
            className="flex-1 bg-muted border-0 rounded-full px-4"
            disabled={sending}
          />
          <Button
            onClick={handleSend}
            disabled={!newMessage.trim() || sending}
            size="icon"
            className="rounded-full bg-primary hover:bg-primary/90"
            aria-label="Enviar mensaje"
          >
            <Send className="w-5 h-5" />
          </Button>
        </div>
      </div>
    </div>
  )
}
