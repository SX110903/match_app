const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? 'ws://localhost:8080'

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
let currentToken: string | null = null
const listeners = new Set<WSListener>()

export function connectWS(token: string): void {
  currentToken = token
  shouldReconnect = true
  reconnectDelay = 1000
  _connect()
}

export function disconnectWS(): void {
  shouldReconnect = false
  currentToken = null
  if (reconnectTimer) { clearTimeout(reconnectTimer); reconnectTimer = null }
  if (ws) { ws.close(); ws = null }
}

export function addWSListener(fn: WSListener): void { listeners.add(fn) }
export function removeWSListener(fn: WSListener): void { listeners.delete(fn) }

export function sendWS(data: object): void {
  if (ws?.readyState === WebSocket.OPEN) ws.send(JSON.stringify(data))
}

function _connect(): void {
  if (!currentToken) return
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) return

  ws = new WebSocket(`${WS_URL}/ws?token=${encodeURIComponent(currentToken)}`)

  ws.onopen = () => { reconnectDelay = 1000 }

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data as string) as WSMessage
      listeners.forEach((fn) => fn(data))
    } catch { /* ignore */ }
  }

  ws.onclose = () => {
    ws = null
    if (shouldReconnect && currentToken) {
      reconnectTimer = setTimeout(() => {
        reconnectDelay = Math.min(reconnectDelay * 2, 30_000)
        _connect()
      }, reconnectDelay)
    }
  }

  ws.onerror = () => { ws?.close() }
}
