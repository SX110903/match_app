"use client"

import { apiClient } from "@/lib/api-client"

interface AdData {
  id: string
  title: string
  description?: string
  image_url?: string
  cta_text: string
  cta_url: string
}

export function AdBanner({ ad }: { ad: AdData | null }) {
  if (!ad) return null

  const handleClick = async () => {
    try {
      await apiClient(`/api/v1/ads/${ad.id}/click`, { method: "POST" })
    } catch { /* best-effort */ }
    window.open(ad.cta_url, "_blank", "noopener,noreferrer")
  }

  return (
    <div className="mx-4 my-3 rounded-xl border border-border bg-card overflow-hidden">
      <div className="px-3 py-1 bg-muted/50 border-b border-border/50">
        <span className="text-[10px] text-muted-foreground font-medium uppercase tracking-wide">Publicidad</span>
      </div>
      <div className="p-3 flex items-center gap-3">
        {ad.image_url && (
          <img
            src={ad.image_url}
            alt={ad.title}
            className="w-14 h-14 rounded-lg object-cover flex-shrink-0"
            referrerPolicy="no-referrer"
          />
        )}
        <div className="flex-1 min-w-0">
          <p className="text-sm font-semibold text-card-foreground truncate">{ad.title}</p>
          {ad.description && (
            <p className="text-xs text-muted-foreground line-clamp-2">{ad.description}</p>
          )}
        </div>
        <button
          onClick={handleClick}
          className="flex-shrink-0 bg-primary text-primary-foreground text-xs font-semibold px-3 py-1.5 rounded-lg hover:bg-primary/90 transition-colors"
        >
          {ad.cta_text}
        </button>
      </div>
    </div>
  )
}
