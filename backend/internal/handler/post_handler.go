package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/service"
	ws "github.com/SX110903/match_app/backend/internal/websocket"
	"github.com/SX110903/match_app/backend/pkg/response"
	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postSvc service.IPostService
	hub     *ws.Hub
}

func NewPostHandler(postSvc service.IPostService, hub *ws.Hub) *PostHandler {
	return &PostHandler{postSvc: postSvc, hub: hub}
}

func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 20
	}
	posts, err := h.postSvc.GetFeed(r.Context(), claims.Subject, page, limit)
	if err != nil {
		response.InternalError(w)
		return
	}
	if posts == nil {
		posts = []service.PostResponse{}
	}
	response.OK(w, posts)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	var req service.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Content == "" {
		response.BadRequest(w, "content is required")
		return
	}
	if len(req.Content) > 2000 {
		response.BadRequest(w, "content too long (max 2000 chars)")
		return
	}
	post, err := h.postSvc.CreatePost(r.Context(), claims.Subject, req.Content, req.ImageURL)
	if err != nil {
		response.InternalError(w)
		return
	}
	go h.hub.BroadcastAll(map[string]interface{}{
		"type": "new_post",
		"post": map[string]interface{}{
			"id":             post.ID,
			"user_id":        post.UserID,
			"content":        post.Content,
			"image_url":      post.ImageURL,
			"likes_count":    post.LikesCount,
			"author_name":    post.AuthorName,
			"author_avatar":  post.AuthorAvatar,
			"created_at":     post.CreatedAt,
			"is_liked_by_me": false,
		},
	}, claims.Subject)
	response.Created(w, post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	postID := chi.URLParam(r, "postId")
	if err := h.postSvc.DeletePost(r.Context(), claims.Subject, postID); err != nil {
		response.NotFound(w, "post not found")
		return
	}
	response.NoContent(w)
}

func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	postID := chi.URLParam(r, "postId")
	if err := h.postSvc.LikePost(r.Context(), claims.Subject, postID); err != nil {
		response.InternalError(w)
		return
	}
	response.NoContent(w)
}

func (h *PostHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	postID := chi.URLParam(r, "postId")
	if err := h.postSvc.UnlikePost(r.Context(), claims.Subject, postID); err != nil {
		response.InternalError(w)
		return
	}
	response.NoContent(w)
}

func (h *PostHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postId")
	comments, err := h.postSvc.GetComments(r.Context(), postID)
	if err != nil {
		response.InternalError(w)
		return
	}
	if comments == nil {
		comments = []service.CommentResponse{}
	}
	response.OK(w, comments)
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "not authenticated")
		return
	}
	postID := chi.URLParam(r, "postId")
	var req service.AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Content == "" {
		response.BadRequest(w, "content is required")
		return
	}
	comment, err := h.postSvc.AddComment(r.Context(), claims.Subject, postID, req.Content)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Created(w, comment)
}
