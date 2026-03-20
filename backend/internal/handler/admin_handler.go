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

type AdminHandler struct {
	adminSvc service.IAdminService
}

func NewAdminHandler(adminSvc service.IAdminService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	if err := h.adminSvc.AssertAdmin(r.Context(), claims.Subject); err != nil {
		response.Forbidden(w, "admin required")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	users, err := h.adminSvc.ListUsers(r.Context(), page, 50)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, users)
}

func (h *AdminHandler) FreezeUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.UserAdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.adminSvc.FreezeUser(r.Context(), claims.Subject, req.UserID); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"status": "frozen"})
}

func (h *AdminHandler) UnfreezeUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.UserAdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.adminSvc.UnfreezeUser(r.Context(), claims.Subject, req.UserID); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"status": "unfrozen"})
}

func (h *AdminHandler) SetVIP(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.SetVIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.VIPLevel < 0 || req.VIPLevel > 5 {
		response.BadRequest(w, "vip_level must be 0-5")
		return
	}
	if err := h.adminSvc.SetVIPLevel(r.Context(), claims.Subject, req.UserID, req.VIPLevel); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		case domain.ErrNotFound:
			response.NotFound(w, "user not found")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]int{"vip_level": req.VIPLevel})
}

func (h *AdminHandler) AdjustCredits(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.AdjustCreditsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	const maxCreditsDelta = 10000
	if req.Delta > maxCreditsDelta || req.Delta < -maxCreditsDelta {
		response.BadRequest(w, "delta must be between -10000 and 10000")
		return
	}
	if err := h.adminSvc.AdjustCredits(r.Context(), claims.Subject, req.UserID, req.Delta); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		case domain.ErrInvalidInput:
			response.BadRequest(w, "delta exceeds limit or would result in negative balance")
		case domain.ErrNotFound:
			response.NotFound(w, "user not found")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]string{"status": "ok"})
}

func (h *AdminHandler) SetAdminRole(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req struct {
		UserID  string `json:"user_id"`
		IsAdmin bool   `json:"is_admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.adminSvc.SetAdmin(r.Context(), claims.Subject, req.UserID, req.IsAdmin); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		case domain.ErrSelfAction:
			response.Forbidden(w, "cannot modify your own admin role")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, map[string]bool{"is_admin": req.IsAdmin})
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.BadRequest(w, "missing user id")
		return
	}
	if err := h.adminSvc.DeleteUser(r.Context(), claims.Subject, userID); err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		case domain.ErrSelfAction:
			response.Forbidden(w, "cannot delete your own account")
		case domain.ErrNotFound:
			response.NotFound(w, "user not found")
		default:
			response.InternalError(w)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	logs, err := h.adminSvc.GetAuditLog(r.Context(), claims.Subject, page, 50)
	if err != nil {
		switch err {
		case domain.ErrForbidden:
			response.Forbidden(w, "admin required")
		default:
			response.InternalError(w)
		}
		return
	}
	if logs == nil {
		logs = []domain.AdminLog{}
	}
	response.OK(w, logs)
}

func (h *AdminHandler) GetNotificationSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	settings, err := h.adminSvc.GetNotificationSettings(r.Context(), claims.Subject)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, settings)
}

func (h *AdminHandler) SaveNotificationSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.NotificationSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	settings := &domain.NotificationSettings{
		UserID:      claims.Subject,
		NewMatches:  req.NewMatches,
		NewMessages: req.NewMessages,
		NewsUpdates: req.NewsUpdates,
		Marketing:   req.Marketing,
	}
	if err := h.adminSvc.SaveNotificationSettings(r.Context(), settings); err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, settings)
}

func (h *AdminHandler) GetPrivacySettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	settings, err := h.adminSvc.GetPrivacySettings(r.Context(), claims.Subject)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, settings)
}

func (h *AdminHandler) SavePrivacySettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.PrivacySettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	settings := &domain.PrivacySettings{
		UserID:           claims.Subject,
		ShowOnlineStatus: req.ShowOnlineStatus,
		ShowLastSeen:     req.ShowLastSeen,
		ShowDistance:     req.ShowDistance,
		IncognitoMode:    req.IncognitoMode,
	}
	if err := h.adminSvc.SavePrivacySettings(r.Context(), settings); err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, settings)
}
