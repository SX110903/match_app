package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	gws "github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/SX110903/match_app/backend/internal/service"
	ws "github.com/SX110903/match_app/backend/internal/websocket"
	"github.com/SX110903/match_app/backend/pkg/logger"
	"github.com/SX110903/match_app/backend/pkg/response"
)

var wsUpgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Restrict in production via config; for now accept all origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

// inboundWSMessage is what clients send over WebSocket.
type inboundWSMessage struct {
	Type    string `json:"type"`
	MatchID string `json:"match_id"`
	Text    string `json:"text"`
}

type WSHandler struct {
	hub      *ws.Hub
	redis    *redis.Client
	msgSvc   service.IMessageService
	matchSvc service.IMatchService
}

func NewWSHandler(
	hub *ws.Hub,
	redisClient *redis.Client,
	msgSvc service.IMessageService,
	matchSvc service.IMatchService,
) *WSHandler {
	return &WSHandler{
		hub:      hub,
		redis:    redisClient,
		msgSvc:   msgSvc,
		matchSvc: matchSvc,
	}
}

// ServeWS upgrades the HTTP connection to WebSocket.
// Autenticación mediante ticket de un solo uso (30s TTL) — el JWT nunca viaja en la URL.
func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	ticket := r.URL.Query().Get("ticket")
	if ticket == "" {
		response.Unauthorized(w, "missing ticket")
		return
	}

	// GetDel: atómico — lee y borra en una sola operación, evitando reuso del ticket
	key := fmt.Sprintf("ws:ticket:%s", ticket)
	userID, err := h.redis.GetDel(r.Context(), key).Result()
	if err != nil || userID == "" {
		response.Unauthorized(w, "invalid or expired ticket")
		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Err(err).Msg("websocket upgrade failed")
		return
	}

	client := ws.NewClient(h.hub, conn, userID)
	h.hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump(h.handleInbound)
}

func (h *WSHandler) handleInbound(userID string, raw []byte) {
	var msg inboundWSMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}

	switch msg.Type {
	case "message":
		if msg.MatchID == "" || msg.Text == "" {
			return
		}
		h.handleChatMessage(userID, msg.MatchID, msg.Text)
	case "ping":
		h.hub.BroadcastToUser(userID, []byte(`{"type":"pong"}`))
	}
}

func (h *WSHandler) handleChatMessage(userID, matchID, text string) {
	ctx := context.Background()

	created, otherUserID, err := h.msgSvc.SendMessage(ctx, userID, matchID, text)
	if err != nil {
		return
	}

	h.hub.BroadcastMessage(matchID, userID, otherUserID, *created)
}
