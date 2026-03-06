"use client"

import { motion } from "framer-motion"
import { X, Heart, Star, RotateCcw } from "lucide-react"
import { Button } from "@/components/ui/button"

interface ActionButtonsProps {
  onLike: () => void
  onDislike: () => void
  onSuperLike: () => void
  onUndo: () => void
  canUndo: boolean
}

export function ActionButtons({
  onLike,
  onDislike,
  onSuperLike,
  onUndo,
  canUndo,
}: ActionButtonsProps) {
  return (
    <div className="flex items-center justify-center gap-4">
      <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
        <Button
          onClick={onUndo}
          disabled={!canUndo}
          size="lg"
          variant="outline"
          className="w-12 h-12 rounded-full p-0 border-2 border-muted-foreground/30 hover:border-accent hover:bg-accent/10 disabled:opacity-30"
          aria-label="Deshacer"
        >
          <RotateCcw className="w-5 h-5 text-accent" />
        </Button>
      </motion.div>

      <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
        <Button
          onClick={onDislike}
          size="lg"
          variant="outline"
          className="w-16 h-16 rounded-full p-0 border-2 border-destructive/50 hover:border-destructive hover:bg-destructive/10"
          aria-label="No me gusta"
        >
          <X className="w-8 h-8 text-destructive" />
        </Button>
      </motion.div>

      <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
        <Button
          onClick={onSuperLike}
          size="lg"
          variant="outline"
          className="w-12 h-12 rounded-full p-0 border-2 border-blue-500/50 hover:border-blue-500 hover:bg-blue-500/10"
          aria-label="Super Like"
        >
          <Star className="w-5 h-5 text-blue-500 fill-blue-500" />
        </Button>
      </motion.div>

      <motion.div whileHover={{ scale: 1.1 }} whileTap={{ scale: 0.9 }}>
        <Button
          onClick={onLike}
          size="lg"
          variant="outline"
          className="w-16 h-16 rounded-full p-0 border-2 border-green-500/50 hover:border-green-500 hover:bg-green-500/10"
          aria-label="Me gusta"
        >
          <Heart className="w-8 h-8 text-green-500 fill-green-500" />
        </Button>
      </motion.div>
    </div>
  )
}
