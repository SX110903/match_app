package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/pkg/logger"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type PhotoHandler struct {
	photoSvc service.IPhotoService
}

func NewPhotoHandler(photoSvc service.IPhotoService) *PhotoHandler {
	return &PhotoHandler{photoSvc: photoSvc}
}

func (h *PhotoHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	var req service.AddPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	photo, err := h.photoSvc.AddPhoto(r.Context(), claims.Subject, req.URL)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			response.BadRequest(w, "URL must be a valid Imgur link (https://i.imgur.com/ or https://imgur.com/)")
		case domain.ErrConflict:
			response.Conflict(w, "maximum number of photos reached (6)")
		default:
			logger.Error().Err(err).Msg("add photo failed")
			response.InternalError(w)
		}
		return
	}

	response.Created(w, photo)
}

func (h *PhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}

	photoID := chi.URLParam(r, "id")

	if err := h.photoSvc.DeletePhoto(r.Context(), claims.Subject, photoID); err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(w, "photo not found")
		default:
			logger.Error().Err(err).Msg("delete photo failed")
			response.InternalError(w)
		}
		return
	}

	response.NoContent(w)
}
