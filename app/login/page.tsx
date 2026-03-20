'use client'

import { useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { motion } from 'framer-motion'
import { Heart, Eye, EyeOff, Loader2, ShieldCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuth } from '@/lib/auth-context'
import { apiClient, APIError } from '@/lib/api-client'

export default function LoginPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { login } = useAuth()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  // 2FA second step
  const [requires2FA, setRequires2FA] = useState(false)
  const [tempToken, setTempToken] = useState('')
  const [totpCode, setTotpCode] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const result = await login(email, password)
      // login() may resolve with a requires_2fa payload if auth-context exposes it
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const r = result as any
      if (r?.requires_2fa && r?.temp_token) {
        setTempToken(r.temp_token)
        setRequires2FA(true)
        return
      }
      const returnTo = searchParams.get('returnTo') ?? '/'
      router.push(returnTo)
    } catch (err) {
      if (err instanceof APIError) {
        // Backend returns requires_2fa in the error body for some implementations
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const body = (err as any).body
        if (body?.requires_2fa && body?.temp_token) {
          setTempToken(body.temp_token)
          setRequires2FA(true)
          return
        }
        if (err.status === 401) setError('Email o contraseña incorrectos')
        else if (err.status === 403) setError('Cuenta bloqueada temporalmente')
        else if (err.status === 429) setError('Demasiados intentos. Intenta más tarde.')
        else setError(err.message)
      } else {
        setError('Error de conexión. Intenta de nuevo.')
      }
    } finally {
      setLoading(false)
    }
  }

  const handle2FASubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await apiClient('/api/v1/auth/login/2fa', { method: 'POST', body: { temp_token: tempToken, code: totpCode } })
      const returnTo = searchParams.get('returnTo') ?? '/'
      router.push(returnTo)
    } catch (err) {
      if (err instanceof APIError) {
        if (err.status === 401) setError('Código incorrecto o expirado')
        else setError(err.message)
      } else {
        setError('Error de conexión. Intenta de nuevo.')
      }
    } finally {
      setLoading(false)
    }
  }

  // 2FA step UI
  if (requires2FA) {
    return (
      <div className="min-h-screen bg-background flex flex-col items-center justify-center p-4">
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="w-full max-w-sm space-y-8">
          <div className="text-center">
            <div className="flex justify-center mb-4">
              <div className="w-16 h-16 bg-primary rounded-2xl flex items-center justify-center shadow-lg">
                <ShieldCheck className="w-9 h-9 text-primary-foreground" />
              </div>
            </div>
            <h1 className="text-2xl font-bold text-card-foreground">Verificación 2FA</h1>
            <p className="text-muted-foreground mt-1 text-sm">Introduce el código de tu app de autenticación</p>
          </div>
          <form onSubmit={handle2FASubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="totp">Código de 6 dígitos</Label>
              <Input
                id="totp"
                type="text"
                inputMode="numeric"
                maxLength={6}
                placeholder="123456"
                value={totpCode}
                onChange={(e) => setTotpCode(e.target.value.replace(/\D/g, ''))}
                required
                autoComplete="one-time-code"
                className="text-center tracking-widest text-lg"
              />
            </div>
            {error && <motion.p initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="text-sm text-destructive text-center">{error}</motion.p>}
            <Button type="submit" className="w-full" disabled={loading || totpCode.length < 6}>
              {loading ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : null}
              Verificar
            </Button>
            <button type="button" onClick={() => { setRequires2FA(false); setTempToken(''); setTotpCode(''); setError('') }} className="w-full text-sm text-muted-foreground hover:text-foreground text-center">
              ← Volver al inicio de sesión
            </button>
          </form>
        </motion.div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-sm space-y-8"
      >
        {/* Logo */}
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <div className="w-16 h-16 bg-primary rounded-2xl flex items-center justify-center shadow-lg">
              <Heart className="w-9 h-9 text-primary-foreground fill-current" />
            </div>
          </div>
          <h1 className="text-3xl font-bold text-card-foreground">Match Hub</h1>
          <p className="text-muted-foreground mt-1">Inicia sesión en tu cuenta</p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="tu@email.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoComplete="email"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Contraseña</Label>
            <div className="relative">
              <Input
                id="password"
                type={showPassword ? 'text' : 'password'}
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                autoComplete="current-password"
                className="pr-10"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
              </button>
            </div>
          </div>

          {error && (
            <motion.p
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="text-sm text-destructive text-center"
            >
              {error}
            </motion.p>
          )}

          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : null}
            Iniciar sesión
          </Button>
        </form>

        <p className="text-center text-sm text-muted-foreground">
          ¿No tienes cuenta?{' '}
          <Link href="/register" className="text-primary font-medium hover:underline">
            Regístrate aquí
          </Link>
        </p>
      </motion.div>
    </div>
  )
}
