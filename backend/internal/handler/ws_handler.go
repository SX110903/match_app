package handler

import (
	"context"
	"encoding/json"
	"net/http"

	gws "github.com/gorilla/websocket"
	"github.com/SX110903/match_app/backend/internal/auth"
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
	jwtSvc   auth.IJWTService
	blklist  auth.ITokenBlacklist
	msgSvc   service.IMessageService
	matchSvc service.IMatchService
}

func NewWSHandler(
	hub *ws.Hub,
	jwtSvc auth.IJWTService,
	blklist auth.ITokenBlacklist,
	msgSvc service.IMessageService,
	matchSvc service.IMatchService,
) *WSHandler {
	return &WSHandler{
		hub:      hub,
		jwtSvc:   jwtSvc,
		blklist:  blklist,
		msgSvc:   msgSvc,
		matchSvc: matchSvc,
	}
}

// ServeWS upgrades the HTTP connection to WebSocket.
// The JWT is passed as a query parameter: GET /ws?token=<access_token>
func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		response.Unauthorized(w, "missing token")
		return
	}

	claims, err := h.jwtSvc.ValidateAccessToken(tokenStr)
	if err != nil {
		response.Unauthorized(w, "invalid token")
		return
	}

	blacklisted, err := h.blklist.IsBlacklisted(r.Context(), claims.ID)
	if err != nil || blacklisted {
		response.Unauthorized(w, "token revoked")
		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Err(err).Msg("websocket upgrade failed")
		return
	}

	client := ws.NewClient(h.hub, conn, claims.Subject)
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
	}
}

func (h *WSHandler) handleChatMessage(userID, matchID, text string) {
	ctx := context.Background()

	created, err := h.msgSvc.SendMessage(ctx, userID, matchID, text)
	if err != nil {
		return
	}

	// Determine the other participant so we broadcast to both sides.
	match, err := h.matchSvc.GetMatch(ctx, userID, matchID)
	if err != nil {
		// Still broadcast to sender at minimum.
		h.hub.BroadcastMessage(matchID, userID, "", *created)
		return
	}

	otherUserID := match.User1ID
	if match.User1ID == userID {
		otherUserID = match.User2ID
	}

	h.hub.BroadcastMessage(matchID, userID, otherUserID, *created)
}
