import { API_URL } from './constants'

let _token: string | null = null

export function setAuthToken(token: string | null): void {
  _token = token
}

export function getAuthToken(): string | null {
  return _token
}

interface RequestOptions extends Omit<RequestInit, 'body'> {
  body?: object
  skipRefresh?: boolean
}

async function doRequest(path: string, opts: RequestOptions = {}): Promise<Response> {
  const { body, skipRefresh, ...rest } = opts
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(rest.headers as Record<string, string> | undefined),
  }
  if (_token) headers['Authorization'] = `Bearer ${_token}`

  const res = await fetch(`${API_URL}${path}`, {
    ...rest,
    credentials: 'include',
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })

  if (res.status === 401 && !skipRefresh) {
    const refreshed = await tryRefresh()
    if (refreshed) return doRequest(path, { ...opts, skipRefresh: true })
    setAuthToken(null)
    throw new APIError(401, 'UNAUTHORIZED')
  }

  return res
}

async function tryRefresh(): Promise<boolean> {
  try {
    const res = await fetch(`${API_URL}/api/v1/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
    })
    if (!res.ok) return false
    const data = await res.json() as { data?: { access_token?: string } }
    const token = data?.data?.access_token
    if (token) { setAuthToken(token); return true }
    return false
  } catch {
    return false
  }
}

export async function apiClient<T = unknown>(path: string, opts: RequestOptions = {}): Promise<T> {
  const res = await doRequest(path, opts)
  if (!res.ok) {
    const body = await res.json().catch(() => ({})) as Record<string, string>
    throw new APIError(res.status, body?.error ?? 'Request failed')
  }
  if (res.status === 204) return null as T
  const data = await res.json() as { data: T }
  return data.data
}

export class APIError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'APIError'
  }
}
