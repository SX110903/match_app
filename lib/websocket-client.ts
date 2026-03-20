import { API_URL, WS_URL } from './constants'

export interface WSMessage {
  type: string
  match_id: string
  message?: {
    id: string
    match_id: string
    sender_id: string
    text: string
    read_at: string | null
    created_at: string
  }
}

type WSListener = (msg: WSMessage) => void

let ws: WebSocket | null = null
let reconnectDelay = 1000
let shouldReconnect = false
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let pingTimer: ReturnType<typeof setInterval> | null = null
let currentToken: string | null = null
const listeners = new Set<WSListener>()

export function connectWS(token: string): void {
  currentToken = token
  shouldReconnect = true
  reconnectDelay = 1000
  void _connectWithTicket()
}

export function disconnectWS(): void {
  shouldReconnect = false
  currentToken = null
  if (reconnectTimer) { clearTimeout(reconnectTimer); reconnectTimer = null }
  if (pingTimer) { clearInterval(pingTimer); pingTimer = null }
  if (ws) { ws.close(); ws = null }
}

export function addWSListener(fn: WSListener): void { listeners.add(fn) }
export function removeWSListener(fn: WSListener): void { listeners.delete(fn) }

export function sendWS(data: object): void {
  if (ws?.readyState === WebSocket.OPEN) ws.send(JSON.stringify(data))
}

// Solicita un ticket de un solo uso al backend y abre el WebSocket con él.
// El JWT nunca aparece en la URL.
async function _connectWithTicket(): Promise<void> {
  if (!currentToken) return
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) return

  try {
    const res = await fetch(`${API_URL}/api/v1/auth/ws-ticket`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${currentToken}`,
      },
    })
    if (!res.ok) { _scheduleReconnect(); return }
    const body = await res.json() as { data: { ticket: string } }
    const ticket = body?.data?.ticket
    if (!ticket) { _scheduleReconnect(); return }

    ws = new WebSocket(`${WS_URL}/ws?ticket=${encodeURIComponent(ticket)}`)
    _attachHandlers()
  } catch {
    _scheduleReconnect()
  }
}

function _scheduleReconnect(): void {
  if (shouldReconnect && currentToken) {
    reconnectTimer = setTimeout(() => {
      reconnectDelay = Math.min(reconnectDelay * 2, 30_000)
      void _connectWithTicket()
    }, reconnectDelay)
  }
}

function _attachHandlers(): void {
  if (!ws) return

  ws.onopen = () => {
    reconnectDelay = 1000
    if (pingTimer) clearInterval(pingTimer)
    pingTimer = setInterval(() => {
      if (ws?.readyState === WebSocket.OPEN) ws.send(JSON.stringify({ type: 'ping' }))
    }, 30_000)
  }

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data as string) as WSMessage
      listeners.forEach((fn) => fn(data))
    } catch (err) {
      if (process.env.NODE_ENV === 'development')
        console.warn('[WS] parse error:', err, typeof event.data === 'string' ? event.data.slice(0, 200) : '(binary)')
    }
  }

  ws.onclose = () => {
    ws = null
    if (pingTimer) { clearInterval(pingTimer); pingTimer = null }
    _scheduleReconnect()
  }

  ws.onerror = () => { ws?.close() }
}
