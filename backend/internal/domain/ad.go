package domain

import "time"

type Ad struct {
	ID          string    `db:"id"           json:"id"`
	Title       string    `db:"title"        json:"title"`
	Description *string   `db:"description"  json:"description,omitempty"`
	ImageURL    *string   `db:"image_url"    json:"image_url,omitempty"`
	CTAText     string    `db:"cta_text"     json:"cta_text"`
	CTAURL      string    `db:"cta_url"      json:"cta_url"`
	TargetBadge string    `db:"target_badge" json:"target_badge"`
	Active      bool      `db:"active"       json:"active"`
	Impressions int       `db:"impressions"  json:"impressions"`
	Clicks      int       `db:"clicks"       json:"clicks"`
	CreatedBy   string    `db:"created_by"   json:"created_by"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}
