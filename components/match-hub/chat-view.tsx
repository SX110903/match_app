"use client"

import { useState, useRef, useEffect } from "react"
import { motion, AnimatePresence } from "framer-motion"
import { ArrowLeft, Send, MoreVertical, Phone, Video } from "lucide-react"
import { Match, Message } from "@/lib/types"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

interface ChatViewProps {
  match: Match
  messages: Message[]
  onBack: () => void
  onSendMessage: (text: string) => void
}

export function ChatView({ match, messages, onBack, onSendMessage }: ChatViewProps) {
  const [newMessage, setNewMessage] = useState("")
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  const handleSend = () => {
    if (newMessage.trim()) {
      onSendMessage(newMessage.trim())
      setNewMessage("")
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString("es-ES", { hour: "2-digit", minute: "2-digit" })
  }

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Header */}
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
            <p className="text-xs text-muted-foreground">En línea</p>
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

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {/* Match info */}
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
          <p className="text-xs text-muted-foreground mt-2">
            Hicieron match. ¡Empieza la conversación!
          </p>
        </div>

        {/* Messages list */}
        <AnimatePresence>
          {messages.map((message, index) => {
            const isMe = message.senderId === "me"
            return (
              <motion.div
                key={message.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.05 }}
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
                  <p
                    className={`text-[10px] mt-1 ${
                      isMe ? "text-primary-foreground/70" : "text-muted-foreground"
                    }`}
                  >
                    {formatTime(message.timestamp)}
                  </p>
                </div>
              </motion.div>
            )
          })}
        </AnimatePresence>
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <div className="p-4 bg-card border-t border-border">
        <div className="flex items-center gap-2">
          <Input
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            onKeyDown={handleKeyPress}
            placeholder="Escribe un mensaje..."
            className="flex-1 bg-muted border-0 rounded-full px-4"
          />
          <Button
            onClick={handleSend}
            disabled={!newMessage.trim()}
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
