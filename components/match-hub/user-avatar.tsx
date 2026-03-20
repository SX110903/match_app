'use client'
import { useState } from 'react'
import { AVATAR_BASE } from '@/lib/constants'

interface Props {
  src?: string
  alt: string
  className?: string
  fallbackId?: string
}

export function UserAvatar({ src, alt, className, fallbackId }: Props) {
  const [err, setErr] = useState(false)
  const fallback = `${AVATAR_BASE}?u=${fallbackId ?? 'default'}`
  const imgSrc = err ? fallback : (src || fallback)

  return (
    <img
      src={imgSrc}
      alt={alt}
      className={className}
      referrerPolicy="no-referrer"
      onError={() => setErr(true)}
    />
  )
}
