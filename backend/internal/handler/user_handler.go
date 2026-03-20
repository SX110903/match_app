package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/internal/validator"
	"github.com/SX110903/match_app/backend/pkg/logger"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type UserHandler struct {
	userSvc service.IUserService
}

func NewUserHandler(userSvc service.IUserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	profile, err := h.userSvc.GetMe(r.Context(), claims.Subject)
	if err != nil {
		logger.Error().Err(err).Msg("get me failed")
		response.InternalError(w)
		return
	}

	response.OK(w, profile)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	var req service.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	if err := h.userSvc.UpdateMe(r.Context(), claims.Subject, req); err != nil {
		logger.Error().Err(err).Msg("update me failed")
		response.InternalError(w)
		return
	}

	response.OK(w, map[string]string{"message": "profile updated successfully"})
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	if err := h.userSvc.DeleteMe(r.Context(), claims.Subject); err != nil {
		logger.Error().Err(err).Msg("delete me failed")
		response.InternalError(w)
		return
	}

	response.OK(w, map[string]string{"message": "account deleted"})
}

func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	var req service.UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	if err := h.userSvc.UpdatePreferences(r.Context(), claims.Subject, req); err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(w, "user not found")
		default:
			logger.Error().Err(err).Msg("update preferences failed")
			response.InternalError(w)
		}
		return
	}

	response.OK(w, map[string]string{"message": "preferences updated"})
}

func (h *UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	targetID := chi.URLParam(r, "id")
	profile, err := h.userSvc.GetPublicProfile(r.Context(), claims.Subject, targetID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(w, "profile not found")
		default:
			logger.Error().Err(err).Msg("get public profile failed")
			response.InternalError(w)
		}
		return
	}
	response.OK(w, profile)
}
