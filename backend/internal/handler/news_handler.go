package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
	"github.com/SX110903/match_app/backend/pkg/response"
	"github.com/go-chi/chi/v5"
)

type NewsHandler struct {
	newsSvc  service.INewsService
	adminSvc service.IAdminService
}

func NewNewsHandler(newsSvc service.INewsService, adminSvc service.IAdminService) *NewsHandler {
	return &NewsHandler{newsSvc: newsSvc, adminSvc: adminSvc}
}

func (h *NewsHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	adminView := false
	if claims, ok := auth.ClaimsFromContext(r.Context()); ok {
		adminView = h.adminSvc.AssertAdmin(r.Context(), claims.Subject) == nil
	}

	articles, err := h.newsSvc.List(r.Context(), category, adminView, page, limit)
	if err != nil {
		response.InternalError(w)
		return
	}
	if articles == nil {
		articles = []service.NewsArticleResponse{}
	}
	response.OK(w, articles)
}

func (h *NewsHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	article, err := h.newsSvc.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "article not found")
		return
	}
	response.OK(w, article)
}

var validNewsCategories = map[string]bool{
	"Tendencias": true, "Tech": true, "Seguridad": true, "Negocios": true,
}

func (h *NewsHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	if err := h.adminSvc.AssertAdmin(r.Context(), claims.Subject); err != nil {
		response.Forbidden(w, "admin required")
		return
	}
	var req service.CreateNewsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Category != "" && !validNewsCategories[req.Category] {
		response.BadRequest(w, "invalid category: must be Tendencias, Tech, Seguridad or Negocios")
		return
	}
	article, err := h.newsSvc.Create(r.Context(), claims.Subject, req)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Created(w, article)
}

func (h *NewsHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	if err := h.adminSvc.AssertAdmin(r.Context(), claims.Subject); err != nil {
		response.Forbidden(w, "admin required")
		return
	}
	id := chi.URLParam(r, "id")
	var req service.UpdateNewsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Category != nil && *req.Category != "" && !validNewsCategories[*req.Category] {
		response.BadRequest(w, "invalid category: must be Tendencias, Tech, Seguridad or Negocios")
		return
	}
	article, err := h.newsSvc.Update(r.Context(), id, req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(w, "article not found")
		default:
			response.InternalError(w)
		}
		return
	}
	response.OK(w, article)
}

func (h *NewsHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	if err := h.adminSvc.AssertAdmin(r.Context(), claims.Subject); err != nil {
		response.Forbidden(w, "admin required")
		return
	}
	id := chi.URLParam(r, "id")
	if err := h.newsSvc.Delete(r.Context(), id); err != nil {
		response.InternalError(w)
		return
	}
	response.NoContent(w)
}
