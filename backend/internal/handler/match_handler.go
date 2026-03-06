package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/pkg/logger"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type MatchHandler struct {
	matchSvc service.IMatchService
}

func NewMatchHandler(matchSvc service.IMatchService) *MatchHandler {
	return &MatchHandler{matchSvc: matchSvc}
}

func (h *MatchHandler) GetCandidates(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	candidates, err := h.matchSvc.GetCandidates(r.Context(), claims.Subject, page, limit)
	if err != nil {
		logger.Error().Err(err).Msg("get candidates failed")
		response.InternalError(w)
		return
	}

	response.OK(w, candidates)
}

func (h *MatchHandler) Swipe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	var req struct {
		UserID    string `json:"user_id"   validate:"required"`
		Direction string `json:"direction" validate:"required,oneof=left right super"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	direction := domain.SwipeDirection(req.Direction)
	result, err := h.matchSvc.Swipe(r.Context(), claims.Subject, req.UserID, direction)
	if err != nil {
		switch err {
		case domain.ErrSelfAction:
			response.BadRequest(w, "cannot swipe on yourself")
		case domain.ErrNotFound:
			response.NotFound(w, "user not found")
		default:
			logger.Error().Err(err).Msg("swipe failed")
			response.InternalError(w)
		}
		return
	}

	response.OK(w, result)
}

func (h *MatchHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	matches, err := h.matchSvc.GetMatches(r.Context(), claims.Subject)
	if err != nil {
		logger.Error().Err(err).Msg("get matches failed")
		response.InternalError(w)
		return
	}

	response.OK(w, matches)
}

func (h *MatchHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	matchID := chi.URLParam(r, "id")
	match, err := h.matchSvc.GetMatch(r.Context(), claims.Subject, matchID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(w, "match not found")
		case domain.ErrForbidden:
			response.Forbidden(w, "access denied")
		default:
			response.InternalError(w)
		}
		return
	}

	response.OK(w, match)
}
