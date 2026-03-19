'use client'

import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'
import { apiClient, setAuthToken } from './api-client'
import { connectWS, disconnectWS } from './websocket-client'

import { API_URL } from './constants'

export interface AuthUser {
  id: string
  email: string
  name: string
  age: number
  bio?: string
  occupation?: string
  location?: string
  photos: Array<{ id: string; url: string }>
  interests: string[]
  totp_enabled: boolean
  is_admin: boolean
  is_frozen: boolean
  vip_level: number
  credits: number
  badge?: string
  follower_count?: number
}

interface AuthContextValue {
  user: AuthUser | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string, age: number) => Promise<void>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const fetchMe = useCallback(async () => {
    const data = await apiClient<AuthUser>('/api/v1/users/me')
    setUser(data)
  }, [])

  // On mount, try to restore session via refresh token cookie
  useEffect(() => {
    async function init() {
      try {
        const res = await fetch(`${API_URL}/api/v1/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
        })
        if (res.ok) {
          const data = await res.json() as { data: { access_token: string } }
          const token = data?.data?.access_token
          if (token) {
            setAuthToken(token)
            await fetchMe()
            connectWS(token)
          }
        }
      } catch {
        // no session — fine
      } finally {
        setIsLoading(false)
      }
    }
    init()
  }, [fetchMe])

  const login = useCallback(async (email: string, password: string) => {
    const data = await apiClient<{ access_token: string }>('/api/v1/auth/login', {
      method: 'POST',
      body: { email, password },
    })
    setAuthToken(data.access_token)
    await fetchMe()
    connectWS(data.access_token)
  }, [fetchMe])

  const register = useCallback(async (email: string, password: string, name: string, age: number) => {
    await apiClient('/api/v1/auth/register', {
      method: 'POST',
      body: { email, password, name, age },
    })
  }, [])

  const logout = useCallback(async () => {
    try {
      await apiClient('/api/v1/auth/logout', { method: 'POST' })
    } finally {
      setAuthToken(null)
      setUser(null)
      disconnectWS()
    }
  }, [])

  return (
    <AuthContext.Provider value={{
      user,
      isLoading,
      isAuthenticated: !!user,
      login,
      register,
      logout,
      refreshUser: fetchMe,
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
