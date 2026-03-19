package repository

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type newsRepository struct{ db *database.DB }

func NewNewsRepository(db *database.DB) INewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) Create(ctx context.Context, article *domain.NewsArticle) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO news_articles (id, author_id, title, summary, content, image_url, category, published, published_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		article.ID, article.AuthorID, article.Title, article.Summary, article.Content,
		article.ImageURL, article.Category, article.Published, article.PublishedAt,
	)
	return err
}

func (r *newsRepository) Update(ctx context.Context, article *domain.NewsArticle) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`UPDATE news_articles
		 SET title=?, summary=?, content=?, image_url=?, category=?, published=?, published_at=?
		 WHERE id=? AND deleted_at IS NULL`,
		article.Title, article.Summary, article.Content, article.ImageURL,
		article.Category, article.Published, article.PublishedAt, article.ID,
	)
	return err
}

func (r *newsRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`UPDATE news_articles SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`, id,
	)
	return err
}

func (r *newsRepository) GetByID(ctx context.Context, id string) (*domain.NewsArticle, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var a domain.NewsArticle
	err := r.db.QueryRowContext(ctx, `
		SELECT na.id, na.author_id, na.title, na.summary, na.content, na.image_url,
		       na.category, na.published, na.published_at, na.created_at, na.updated_at,
		       up.name AS author_name
		FROM news_articles na
		JOIN user_profiles up ON up.user_id = na.author_id
		WHERE na.id = ? AND na.deleted_at IS NULL`, id,
	).Scan(
		&a.ID, &a.AuthorID, &a.Title, &a.Summary, &a.Content, &a.ImageURL,
		&a.Category, &a.Published, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt, &a.AuthorName,
	)
	if err != nil {
		return nil, fmt.Errorf("getting article: %w", err)
	}
	return &a, nil
}

func (r *newsRepository) List(ctx context.Context, category string, publishedOnly bool, limit, offset int) ([]domain.NewsArticle, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	query := `
		SELECT na.id, na.author_id, na.title, na.summary, na.content, na.image_url,
		       na.category, na.published, na.published_at, na.created_at, na.updated_at,
		       up.name AS author_name
		FROM news_articles na
		JOIN user_profiles up ON up.user_id = na.author_id
		WHERE na.deleted_at IS NULL`

	args := []interface{}{}
	if publishedOnly {
		query += " AND na.published = TRUE"
	}
	if category != "" && category != "general" {
		query += " AND na.category = ?"
		args = append(args, category)
	}
	query += " ORDER BY na.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing articles: %w", err)
	}
	defer rows.Close()

	var articles []domain.NewsArticle
	for rows.Next() {
		var a domain.NewsArticle
		if err := rows.Scan(
			&a.ID, &a.AuthorID, &a.Title, &a.Summary, &a.Content, &a.ImageURL,
			&a.Category, &a.Published, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt, &a.AuthorName,
		); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, rows.Err()
}
