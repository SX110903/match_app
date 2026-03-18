package websocket

import (
	"encoding/json"
	"sync"

	"github.com/SX110903/match_app/backend/internal/service"
)

// OutboundMessage is what the server sends to WebSocket clients.
type OutboundMessage struct {
	Type      string                  `json:"type"`
	MatchID   string                  `json:"match_id"`
	Message   *service.MessageResponse `json:"message,omitempty"`
}

type userBroadcast struct {
	userID  string
	payload []byte
}

// Hub manages all active WebSocket client connections, keyed by userID.
type Hub struct {
	mu         sync.RWMutex
	users      map[string]map[string]*Client // userID → clientID → *Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan userBroadcast
}

func NewHub() *Hub {
	return &Hub{
		users:      make(map[string]map[string]*Client),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		broadcast:  make(chan userBroadcast, 256),
	}
}

// Run processes hub events. Must be called in a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if h.users[c.userID] == nil {
				h.users[c.userID] = make(map[string]*Client)
			}
			h.users[c.userID][c.id] = c
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.users[c.userID]; ok {
				delete(conns, c.id)
				if len(conns) == 0 {
					delete(h.users, c.userID)
				}
			}
			h.mu.Unlock()
			close(c.send)

		case msg := <-h.broadcast:
			h.mu.RLock()
			conns := h.users[msg.userID]
			h.mu.RUnlock()
			for _, c := range conns {
				select {
				case c.send <- msg.payload:
				default:
					// slow client — drop message
				}
			}
		}
	}
}

// BroadcastMessage sends a new message event to both participants of a match.
// otherUserID may be empty — in that case only the senderID's connections receive it.
func (h *Hub) BroadcastMessage(matchID, senderID, otherUserID string, msg service.MessageResponse) {
	env := OutboundMessage{
		Type:    "message",
		MatchID: matchID,
		Message: &msg,
	}
	payload, err := json.Marshal(env)
	if err != nil {
		return
	}

	h.broadcast <- userBroadcast{userID: senderID, payload: payload}
	if otherUserID != "" && otherUserID != senderID {
		h.broadcast <- userBroadcast{userID: otherUserID, payload: payload}
	}
}

// BroadcastToUser sends a raw JSON payload to all connections of a specific user.
func (h *Hub) BroadcastToUser(userID string, payload []byte) {
	h.broadcast <- userBroadcast{userID: userID, payload: payload}
}

// RegisterClient adds a client to the hub. Called from ws_handler.
func (h *Hub) RegisterClient(c *Client) {
	h.register <- c
}

func (h *Hub) unregisterClient(c *Client) {
	h.unregister <- c
}
