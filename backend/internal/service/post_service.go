package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
	"github.com/google/uuid"
)

type postService struct {
	postRepo repository.IPostRepository
}

func NewPostService(postRepo repository.IPostRepository) IPostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) GetFeed(ctx context.Context, viewerID string, page, limit int) ([]PostResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	posts, err := s.postRepo.GetFeed(ctx, viewerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("getting feed: %w", err)
	}
	out := make([]PostResponse, len(posts))
	for i, p := range posts {
		out[i] = toPostResponse(p)
	}
	return out, nil
}

func (s *postService) CreatePost(ctx context.Context, userID, content string, imageURL *string) (*PostResponse, error) {
	post := &domain.Post{
		ID:        uuid.New().String(),
		UserID:    userID,
		Content:   content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("creating post: %w", err)
	}
	// Re-fetch to get author name
	created, err := s.postRepo.GetByID(ctx, post.ID, userID)
	if err != nil {
		resp := toPostResponse(*post)
		return &resp, nil
	}
	resp := toPostResponse(*created)
	return &resp, nil
}

func (s *postService) DeletePost(ctx context.Context, userID, postID string) error {
	return s.postRepo.Delete(ctx, postID, userID)
}

func (s *postService) LikePost(ctx context.Context, userID, postID string) error {
	like := &domain.PostLike{
		ID:        uuid.New().String(),
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	return s.postRepo.LikePost(ctx, like)
}

func (s *postService) UnlikePost(ctx context.Context, userID, postID string) error {
	return s.postRepo.UnlikePost(ctx, postID, userID)
}

func (s *postService) GetComments(ctx context.Context, postID string) ([]CommentResponse, error) {
	comments, err := s.postRepo.GetComments(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("getting comments: %w", err)
	}
	out := make([]CommentResponse, len(comments))
	for i, c := range comments {
		out[i] = CommentResponse{
			ID:           c.ID,
			PostID:       c.PostID,
			UserID:       c.UserID,
			Content:      c.Content,
			AuthorName:   c.AuthorName,
			AuthorAvatar: c.AuthorAvatar,
			CreatedAt:    c.CreatedAt,
		}
	}
	return out, nil
}

func (s *postService) AddComment(ctx context.Context, userID, postID, content string) (*CommentResponse, error) {
	comment := &domain.PostComment{
		ID:        uuid.New().String(),
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	created, err := s.postRepo.AddComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("adding comment: %w", err)
	}
	resp := &CommentResponse{
		ID:           created.ID,
		PostID:       created.PostID,
		UserID:       created.UserID,
		Content:      created.Content,
		AuthorName:   created.AuthorName,
		AuthorAvatar: created.AuthorAvatar,
		CreatedAt:    created.CreatedAt,
	}
	return resp, nil
}

func toPostResponse(p domain.Post) PostResponse {
	return PostResponse{
		ID:           p.ID,
		UserID:       p.UserID,
		Content:      p.Content,
		ImageURL:     p.ImageURL,
		LikesCount:   p.LikesCount,
		AuthorName:   p.AuthorName,
		AuthorAvatar: p.AuthorAvatar,
		IsLikedByMe:  p.IsLikedByMe,
		CreatedAt:    p.CreatedAt,
	}
}
