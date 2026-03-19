package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/internal/validator"
	"github.com/SX110903/match_app/backend/pkg/logger"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type MessageHandler struct {
	msgSvc service.IMessageService
	hub    MessageBroadcaster
}

// MessageBroadcaster allows broadcasting a new message to WebSocket clients.
// Implemented by websocket.Hub.
type MessageBroadcaster interface {
	BroadcastMessage(matchID, senderID, otherUserID string, msg service.MessageResponse)
}

func NewMessageHandler(msgSvc service.IMessageService, hub MessageBroadcaster) *MessageHandler {
	return &MessageHandler{msgSvc: msgSvc, hub: hub}
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	matchID := chi.URLParam(r, "matchId")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	messages, err := h.msgSvc.GetMessages(r.Context(), claims.Subject, matchID, page, limit)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(w, "match not found")
		case domain.ErrForbidden:
			response.Forbidden(w, "not a participant of this match")
		default:
			logger.Error().Err(err).Msg("get messages failed")
			response.InternalError(w)
		}
		return
	}

	response.OK(w, messages)
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	matchID := chi.URLParam(r, "matchId")

	var req struct {
		Text string `json:"text" validate:"required,min=1,max=2000"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	msg, otherUserID, err := h.msgSvc.SendMessage(r.Context(), claims.Subject, matchID, req.Text)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(w, "match not found")
		case domain.ErrForbidden:
			response.Forbidden(w, "not a participant of this match")
		default:
			logger.Error().Err(err).Msg("send message failed")
			response.InternalError(w)
		}
		return
	}

	if h.hub != nil {
		h.hub.BroadcastMessage(matchID, claims.Subject, otherUserID, *msg)
	}

	response.Created(w, msg)
}

func (h *MessageHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	matchID := chi.URLParam(r, "matchId")

	if err := h.msgSvc.MarkRead(r.Context(), claims.Subject, matchID); err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(w, "match not found")
		case domain.ErrForbidden:
			response.Forbidden(w, "not a participant of this match")
		default:
			response.InternalError(w)
		}
		return
	}

	response.NoContent(w)
}
