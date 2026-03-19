package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type AdHandler struct {
	adSvc service.IAdService
}

func NewAdHandler(adSvc service.IAdService) *AdHandler {
	return &AdHandler{adSvc: adSvc}
}

func (h *AdHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	// Get user badge from query param or default to "none"
	userBadge := r.URL.Query().Get("badge")
	if userBadge == "" {
		userBadge = domain.BadgeNone
	}
	ad, err := h.adSvc.GetActive(r.Context(), claims.Subject, userBadge)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, ad)
}

func (h *AdHandler) RegisterClick(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	adID := chi.URLParam(r, "id")
	if adID == "" {
		response.BadRequest(w, "missing ad id")
		return
	}
	if err := h.adSvc.RegisterClick(r.Context(), adID, claims.Subject); err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, map[string]string{"status": "ok"})
}

func (h *AdHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	ads, err := h.adSvc.AdminList(r.Context(), claims.Subject)
	if err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		default:
			response.InternalError(w)
		}
		return
	}
	if ads == nil {
		ads = []domain.Ad{}
	}
	response.OK(w, ads)
}

func (h *AdHandler) AdminCreate(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.AdCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	ad, err := h.adSvc.AdminCreate(r.Context(), claims.Subject, req)
	if err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		default:
			response.InternalError(w)
		}
		return
	}
	response.Created(w, ad)
}

func (h *AdHandler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	adID := chi.URLParam(r, "id")
	var req service.AdUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	ad, err := h.adSvc.AdminUpdate(r.Context(), claims.Subject, adID, req)
	if err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		case domain.ErrNotFound:
			response.NotFound(w, "ad not found")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, ad)
}

func (h *AdHandler) AdminDelete(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	adID := chi.URLParam(r, "id")
	if err := h.adSvc.AdminDelete(r.Context(), claims.Subject, adID); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		case domain.ErrNotFound:
			response.NotFound(w, "ad not found")
		default:
			response.InternalError(w)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdHandler) AdminToggle(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	adID := chi.URLParam(r, "id")
	if err := h.adSvc.AdminToggle(r.Context(), claims.Subject, adID); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		case domain.ErrNotFound:
			response.NotFound(w, "ad not found")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"status": "toggled"})
}
