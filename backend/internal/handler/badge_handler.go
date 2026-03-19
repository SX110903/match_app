package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type BadgeHandler struct {
	badgeSvc service.IBadgeService
}

func NewBadgeHandler(badgeSvc service.IBadgeService) *BadgeHandler {
	return &BadgeHandler{badgeSvc: badgeSvc}
}

func (h *BadgeHandler) Follow(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	targetID := chi.URLParam(r, "targetId")
	if targetID == "" {
		response.BadRequest(w, "missing targetId")
		return
	}
	if err := h.badgeSvc.Follow(r.Context(), claims.Subject, targetID); err != nil {
		switch err {
		case domain.ErrSelfAction:
			response.BadRequest(w, "cannot follow yourself")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"status": "followed"})
}

func (h *BadgeHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	targetID := chi.URLParam(r, "targetId")
	if targetID == "" {
		response.BadRequest(w, "missing targetId")
		return
	}
	if err := h.badgeSvc.Unfollow(r.Context(), claims.Subject, targetID); err != nil {
		switch err {
		case domain.ErrSelfAction:
			response.BadRequest(w, "cannot unfollow yourself")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"status": "unfollowed"})
}

func (h *BadgeHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	ids, err := h.badgeSvc.GetFollowers(r.Context(), userID, page, 50)
	if err != nil {
		response.InternalError(w)
		return
	}
	if ids == nil {
		ids = []string{}
	}
	response.OK(w, ids)
}

func (h *BadgeHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	ids, err := h.badgeSvc.GetFollowing(r.Context(), userID, page, 50)
	if err != nil {
		response.InternalError(w)
		return
	}
	if ids == nil {
		ids = []string{}
	}
	response.OK(w, ids)
}

func (h *BadgeHandler) RequestVerify(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	if err := h.badgeSvc.RequestVerify(r.Context(), claims.Subject); err != nil {
		switch err {
		case domain.ErrConflict:
			response.Conflict(w, "already has verified_gov badge")
		case domain.ErrInvalidInput:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "insufficient credits"})
		case domain.ErrNotFound:
			response.NotFound(w, "user not found")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"status": "verified"})
}

func (h *BadgeHandler) AdminSetBadge(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.SetBadgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.badgeSvc.AdminSetBadge(r.Context(), claims.Subject, req.UserID, req.Badge); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		case domain.ErrInvalidInput:
			response.BadRequest(w, "invalid badge value")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"badge": req.Badge})
}
