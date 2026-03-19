export type BadgeType = "none" | "influencer" | "verified" | "verified_gov"

export interface Profile {
  id: string
  name: string
  age: number
  bio: string
  images: string[]
  distance: number
  occupation: string
  interests: string[]
  badge?: BadgeType
}

export interface Match {
  id: string
  profile: Profile
  matchedAt: Date
  lastMessage?: string
  unread?: boolean
}

export interface Message {
  id: string
  senderId: string
  text: string
  timestamp: Date
  read: boolean
}

export interface Conversation {
  id: string
  match: Match
  messages: Message[]
}
