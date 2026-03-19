package domain

import "time"

type NewsCategory string

const (
	NewsCategoryTendencias NewsCategory = "tendencias"
	NewsCategoryTech       NewsCategory = "tech"
	NewsCategorySeguridad  NewsCategory = "seguridad"
	NewsCategoryNegocios   NewsCategory = "negocios"
	NewsCategoryGeneral    NewsCategory = "general"
)

type NewsArticle struct {
	ID          string       `db:"id"           json:"id"`
	AuthorID    string       `db:"author_id"    json:"author_id"`
	Title       string       `db:"title"        json:"title"`
	Summary     string       `db:"summary"      json:"summary"`
	Content     string       `db:"content"      json:"content"`
	ImageURL    *string      `db:"image_url"    json:"image_url,omitempty"`
	Category    NewsCategory `db:"category"     json:"category"`
	Published   bool         `db:"published"    json:"published"`
	PublishedAt *time.Time   `db:"published_at" json:"published_at,omitempty"`
	CreatedAt   time.Time    `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"   json:"updated_at"`
	DeletedAt   *time.Time   `db:"deleted_at"   json:"-"`

	// Populated from JOINs
	AuthorName string `db:"-" json:"author_name"`
}
