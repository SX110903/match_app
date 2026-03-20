package handler

import (
	"net/http"
	"strconv"

	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/pkg/response"
)

type ExploreHandler struct {
	exploreSvc service.IExploreService
}

func NewExploreHandler(exploreSvc service.IExploreService) *ExploreHandler {
	return &ExploreHandler{exploreSvc: exploreSvc}
}

func (h *ExploreHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	cursor := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 20
	}
	users, err := h.exploreSvc.GetUsers(r.Context(), claims.Subject, cursor, limit)
	if err != nil {
		response.InternalError(w)
		return
	}
	if users == nil {
		users = []service.ExploreUserResponse{}
	}
	response.OK(w, users)
}

func (h *ExploreHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	cursor := r.URL.Query().Get("cursor")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 20
	}
	posts, err := h.exploreSvc.GetPosts(r.Context(), claims.Subject, cursor, limit)
	if err != nil {
		response.InternalError(w)
		return
	}
	if posts == nil {
		posts = []service.PostResponse{}
	}
	response.OK(w, posts)
}
