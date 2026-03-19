export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'
export const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? 'ws://localhost:8080'
export const AVATAR_BASE = 'https://i.pravatar.cc/300'
export const BADGE_LABELS = {
  none: null,
  influencer: 'Influencer',
  verified: 'Verificado',
  verified_gov: 'Verificado Gubernamental',
} as const
export const VIP_LEVELS = [0, 1, 2, 3, 4, 5] as const
export const NEWS_CATEGORIES = ['Tendencias', 'Tech', 'Seguridad', 'Negocios'] as const
export const MAX_CREDITS_DELTA = 10000
export const MAX_POST_LENGTH = 2000
export const MAX_PHOTO_SIZE_MB = 5
export const VERIFY_BADGE_COST = 50000
