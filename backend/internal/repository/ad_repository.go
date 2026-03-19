package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type adRepository struct {
	db    *database.DB
	redis *redis.Client
}

func NewAdRepository(db *database.DB, rdb *redis.Client) IAdRepository {
	return &adRepository{db: db, redis: rdb}
}

func (r *adRepository) GetActive(ctx context.Context, userBadge string) (*domain.Ad, error) {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var ad domain.Ad
	err := r.db.QueryRowContext(dbCtx,
		`SELECT id, title, description, image_url, cta_text, cta_url, target_badge, active, impressions, clicks, created_by, created_at
		 FROM ads
		 WHERE active = true AND (target_badge = 'all' OR target_badge = ?)
		 ORDER BY RAND() LIMIT 1`,
		userBadge,
	).Scan(
		&ad.ID, &ad.Title, &ad.Description, &ad.ImageURL,
		&ad.CTAText, &ad.CTAURL, &ad.TargetBadge, &ad.Active,
		&ad.Impressions, &ad.Clicks, &ad.CreatedBy, &ad.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get active ad: %w", err)
	}
	return &ad, nil
}

func (r *adRepository) RegisterClick(ctx context.Context, adID, userID string) error {
	// Rate limit: one click per user per ad per 24h
	key := fmt.Sprintf("ad_click:%s:%s", adID, userID)
	set, err := r.redis.SetNX(ctx, key, 1, 24*time.Hour).Result()
	if err != nil {
		return fmt.Errorf("redis setnx: %w", err)
	}
	if !set {
		// Already clicked — idempotent, not an error
		return nil
	}

	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err = r.db.ExecContext(dbCtx,
		`UPDATE ads SET clicks = clicks + 1 WHERE id = ?`, adID,
	)
	return err
}

func (r *adRepository) Create(ctx context.Context, ad *domain.Ad) error {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx,
		`INSERT INTO ads (id, title, description, image_url, cta_text, cta_url, target_badge, active, created_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ad.ID, ad.Title, ad.Description, ad.ImageURL,
		ad.CTAText, ad.CTAURL, ad.TargetBadge, ad.Active, ad.CreatedBy,
	)
	return err
}

func (r *adRepository) Update(ctx context.Context, ad *domain.Ad) error {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx,
		`UPDATE ads SET title=?, description=?, image_url=?, cta_text=?, cta_url=?, target_badge=?, active=?
		 WHERE id=?`,
		ad.Title, ad.Description, ad.ImageURL, ad.CTAText, ad.CTAURL, ad.TargetBadge, ad.Active, ad.ID,
	)
	return err
}

func (r *adRepository) Delete(ctx context.Context, id string) error {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, `DELETE FROM ads WHERE id = ?`, id)
	return err
}

func (r *adRepository) Toggle(ctx context.Context, id string) error {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx,
		`UPDATE ads SET active = NOT active WHERE id = ?`, id,
	)
	return err
}

func (r *adRepository) ListAll(ctx context.Context) ([]domain.Ad, error) {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(dbCtx,
		`SELECT id, title, description, image_url, cta_text, cta_url, target_badge, active, impressions, clicks, created_by, created_at
		 FROM ads ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list ads: %w", err)
	}
	defer rows.Close()

	var ads []domain.Ad
	for rows.Next() {
		var ad domain.Ad
		if err := rows.Scan(
			&ad.ID, &ad.Title, &ad.Description, &ad.ImageURL,
			&ad.CTAText, &ad.CTAURL, &ad.TargetBadge, &ad.Active,
			&ad.Impressions, &ad.Clicks, &ad.CreatedBy, &ad.CreatedAt,
		); err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}
	return ads, rows.Err()
}

func (r *adRepository) GetByID(ctx context.Context, id string) (*domain.Ad, error) {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var ad domain.Ad
	err := r.db.QueryRowContext(dbCtx,
		`SELECT id, title, description, image_url, cta_text, cta_url, target_badge, active, impressions, clicks, created_by, created_at
		 FROM ads WHERE id = ?`, id,
	).Scan(
		&ad.ID, &ad.Title, &ad.Description, &ad.ImageURL,
		&ad.CTAText, &ad.CTAURL, &ad.TargetBadge, &ad.Active,
		&ad.Impressions, &ad.Clicks, &ad.CreatedBy, &ad.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get ad by id: %w", err)
	}
	return &ad, nil
}

func (r *adRepository) IncrementImpressions(ctx context.Context, id string) error {
	dbCtx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx,
		`UPDATE ads SET impressions = impressions + 1 WHERE id = ?`, id,
	)
	return err
}
