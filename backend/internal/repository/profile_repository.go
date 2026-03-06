package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type profileRepository struct {
	db *database.DB
}

func NewProfileRepository(db *database.DB) IProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) Create(ctx context.Context, profile *domain.UserProfile) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_profiles (id, user_id, name, age, bio, occupation, location, latitude, longitude, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		profile.ID, profile.UserID, profile.Name, profile.Age,
		profile.Bio, profile.Occupation, profile.Location,
		profile.Latitude, profile.Longitude,
		profile.CreatedAt, profile.UpdatedAt,
	)
	return err
}

func (r *profileRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserProfile, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var profile domain.UserProfile
	err := r.db.GetContext(ctx, &profile,
		`SELECT * FROM user_profiles WHERE user_id = ?`, userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting profile: %w", err)
	}

	// Load photos
	photos := []struct{ URL string `db:"url"` }{}
	if err := r.db.SelectContext(ctx, &photos,
		`SELECT url FROM user_photos WHERE user_id = ? ORDER BY sort_order`, userID,
	); err == nil {
		for _, p := range photos {
			profile.Photos = append(profile.Photos, p.URL)
		}
	}

	// Load interests
	interests := []struct{ Interest string `db:"interest"` }{}
	if err := r.db.SelectContext(ctx, &interests,
		`SELECT interest FROM user_interests WHERE user_id = ?`, userID,
	); err == nil {
		for _, i := range interests {
			profile.Interests = append(profile.Interests, i.Interest)
		}
	}

	return &profile, nil
}

func (r *profileRepository) Update(ctx context.Context, profile *domain.UserProfile) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE user_profiles SET name=?, age=?, bio=?, occupation=?, location=?, updated_at=? WHERE user_id=?`,
		profile.Name, profile.Age, profile.Bio, profile.Occupation, profile.Location,
		time.Now(), profile.UserID,
	)
	return err
}

func (r *profileRepository) GetPreferences(ctx context.Context, userID string) (*domain.UserPreferences, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var prefs domain.UserPreferences
	err := r.db.GetContext(ctx, &prefs,
		`SELECT * FROM user_preferences WHERE user_id = ?`, userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting preferences: %w", err)
	}
	return &prefs, nil
}

func (r *profileRepository) UpsertPreferences(ctx context.Context, prefs *domain.UserPreferences) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	if prefs.ID == "" {
		prefs.ID = uuid.New().String()
	}
	now := time.Now()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_preferences (id, user_id, min_age, max_age, max_distance_km, interested_in, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE min_age=?, max_age=?, max_distance_km=?, interested_in=?, updated_at=?`,
		prefs.ID, prefs.UserID, prefs.MinAge, prefs.MaxAge, prefs.MaxDistanceKm, prefs.InterestedIn, now, now,
		prefs.MinAge, prefs.MaxAge, prefs.MaxDistanceKm, prefs.InterestedIn, now,
	)
	return err
}
