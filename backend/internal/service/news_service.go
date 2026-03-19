package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
	"github.com/google/uuid"
)

type newsService struct {
	newsRepo repository.INewsRepository
}

func NewNewsService(newsRepo repository.INewsRepository) INewsService {
	return &newsService{newsRepo: newsRepo}
}

func (s *newsService) List(ctx context.Context, category string, adminView bool, page, limit int) ([]NewsArticleResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	// Non-admins see published only
	publishedOnly := !adminView
	articles, err := s.newsRepo.List(ctx, category, publishedOnly, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("listing news: %w", err)
	}
	out := make([]NewsArticleResponse, len(articles))
	for i, a := range articles {
		out[i] = toNewsResponse(a)
	}
	return out, nil
}

func (s *newsService) GetByID(ctx context.Context, id string) (*NewsArticleResponse, error) {
	a, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	resp := toNewsResponse(*a)
	return &resp, nil
}

func (s *newsService) Create(ctx context.Context, authorID string, req CreateNewsRequest) (*NewsArticleResponse, error) {
	var publishedAt *time.Time
	if req.Published {
		now := time.Now()
		publishedAt = &now
	}
	article := &domain.NewsArticle{
		ID:          uuid.New().String(),
		AuthorID:    authorID,
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		ImageURL:    req.ImageURL,
		Category:    domain.NewsCategory(req.Category),
		Published:   req.Published,
		PublishedAt: publishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.newsRepo.Create(ctx, article); err != nil {
		return nil, fmt.Errorf("creating article: %w", err)
	}
	created, err := s.newsRepo.GetByID(ctx, article.ID)
	if err != nil {
		resp := toNewsResponse(*article)
		return &resp, nil
	}
	resp := toNewsResponse(*created)
	return &resp, nil
}

func (s *newsService) Update(ctx context.Context, id string, req UpdateNewsRequest) (*NewsArticleResponse, error) {
	existing, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Summary != nil {
		existing.Summary = *req.Summary
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.ImageURL != nil {
		existing.ImageURL = req.ImageURL
	}
	if req.Category != nil {
		existing.Category = domain.NewsCategory(*req.Category)
	}
	if req.Published != nil {
		existing.Published = *req.Published
		if *req.Published && existing.PublishedAt == nil {
			now := time.Now()
			existing.PublishedAt = &now
		}
	}
	if err := s.newsRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("updating article: %w", err)
	}
	updated, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		resp := toNewsResponse(*existing)
		return &resp, nil
	}
	resp := toNewsResponse(*updated)
	return &resp, nil
}

func (s *newsService) Delete(ctx context.Context, id string) error {
	return s.newsRepo.Delete(ctx, id)
}

func toNewsResponse(a domain.NewsArticle) NewsArticleResponse {
	return NewsArticleResponse{
		ID:          a.ID,
		AuthorID:    a.AuthorID,
		AuthorName:  a.AuthorName,
		Title:       a.Title,
		Summary:     a.Summary,
		Content:     a.Content,
		ImageURL:    a.ImageURL,
		Category:    string(a.Category),
		Published:   a.Published,
		PublishedAt: a.PublishedAt,
		CreatedAt:   a.CreatedAt,
	}
}
