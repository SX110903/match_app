"use client"

import { BADGE_LABELS } from "@/lib/constants"

export type BadgeType = "none" | "influencer" | "verified" | "verified_gov"

export function BadgeIcon({ badge }: { badge?: BadgeType }) {
  if (!badge || badge === "none") return null

  if (badge === "influencer") {
    return (
      <span
        title={BADGE_LABELS.influencer ?? ""}
        className="inline-flex items-center text-amber-500 text-[13px] leading-none"
        aria-label={BADGE_LABELS.influencer ?? ""}
      >
        ⭐
      </span>
    )
  }

  if (badge === "verified") {
    return (
      <span
        title={BADGE_LABELS.verified ?? ""}
        className="inline-flex items-center justify-center w-4 h-4 rounded-full bg-blue-500 text-white text-[9px] font-bold leading-none flex-shrink-0"
        aria-label={BADGE_LABELS.verified ?? ""}
      >
        ✓
      </span>
    )
  }

  if (badge === "verified_gov") {
    return (
      <span
        title={BADGE_LABELS.verified_gov ?? ""}
        className="inline-flex items-center text-purple-600 text-[13px] leading-none"
        aria-label={BADGE_LABELS.verified_gov ?? ""}
      >
        🛡
      </span>
    )
  }

  return null
}
