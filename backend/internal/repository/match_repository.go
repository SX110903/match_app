package repository

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type matchRepository struct {
	db *database.DB
}

func NewMatchRepository(db *database.DB) IMatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) CreateSwipe(ctx context.Context, swipe *domain.Swipe) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO swipes (id, swiper_id, swiped_id, direction, created_at) VALUES (?, ?, ?, ?, ?)`,
		swipe.ID, swipe.SwiperID, swipe.SwipedID, swipe.Direction, swipe.CreatedAt,
	)
	return err
}

func (r *matchRepository) GetSwipe(ctx context.Context, swiperID, swipedID string) (*domain.Swipe, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var swipe domain.Swipe
	err := r.db.GetContext(ctx, &swipe,
		`SELECT * FROM swipes WHERE swiper_id = ? AND swiped_id = ?`,
		swiperID, swipedID,
	)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &swipe, nil
}

func (r *matchRepository) CreateMatch(ctx context.Context, match *domain.Match) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO matches (id, user1_id, user2_id, created_at) VALUES (?, ?, ?, ?)`,
		match.ID, match.User1ID, match.User2ID, match.CreatedAt,
	)
	return err
}

func (r *matchRepository) GetMatchByID(ctx context.Context, id string) (*domain.Match, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var match domain.Match
	if err := r.db.GetContext(ctx, &match, `SELECT * FROM matches WHERE id = ?`, id); err != nil {
		return nil, domain.ErrNotFound
	}
	return &match, nil
}

func (r *matchRepository) GetMatchesByUserID(ctx context.Context, userID string) ([]domain.MatchWithProfile, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	query := `
		SELECT
			m.id, m.user1_id, m.user2_id, m.created_at,
			up.id as profile_id, up.name, up.age, up.bio, up.occupation, up.location,
			(SELECT text FROM messages WHERE match_id = m.id ORDER BY created_at DESC LIMIT 1) as last_message,
			(SELECT COUNT(*) FROM messages WHERE match_id = m.id AND sender_id != ? AND read_at IS NULL) as unread_count
		FROM matches m
		JOIN user_profiles up ON up.user_id = CASE
			WHEN m.user1_id = ? THEN m.user2_id
			ELSE m.user1_id
		END
		WHERE (m.user1_id = ? OR m.user2_id = ?)
		ORDER BY m.created_at DESC`

	rows, err := r.db.QueryxContext(ctx, query, userID, userID, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("getting matches: %w", err)
	}
	defer rows.Close()

	var results []domain.MatchWithProfile
	for rows.Next() {
		var mwp domain.MatchWithProfile
		if err := rows.StructScan(&mwp); err != nil {
			return nil, fmt.Errorf("scanning match: %w", err)
		}
		results = append(results, mwp)
	}
	return results, rows.Err()
}

func (r *matchRepository) GetCandidates(ctx context.Context, userID string, prefs *domain.UserPreferences, limit, offset int) ([]domain.Candidate, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	query := `
		SELECT
			up.id, up.user_id, up.name, up.age, up.bio, up.occupation, up.location,
			ST_Distance_Sphere(
				POINT(up.longitude, up.latitude),
				(SELECT POINT(longitude, latitude) FROM user_profiles WHERE user_id = ?)
			) / 1000 as distance
		FROM user_profiles up
		JOIN users u ON u.id = up.user_id
		WHERE up.user_id != ?
			AND u.deleted_at IS NULL
			AND up.age BETWEEN ? AND ?
			AND up.user_id NOT IN (SELECT swiped_id FROM swipes WHERE swiper_id = ?)
			AND up.user_id NOT IN (SELECT user2_id FROM matches WHERE user1_id = ? UNION SELECT user1_id FROM matches WHERE user2_id = ?)
		HAVING distance <= ?
		ORDER BY distance ASC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryxContext(ctx, query,
		userID, userID,
		prefs.MinAge, prefs.MaxAge,
		userID, userID, userID,
		prefs.MaxDistanceKm,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("getting candidates: %w", err)
	}
	defer rows.Close()

	var candidates []domain.Candidate
	for rows.Next() {
		var c domain.Candidate
		if err := rows.Scan(
			&c.Profile.ID, &c.Profile.UserID, &c.Profile.Name, &c.Profile.Age,
			&c.Profile.Bio, &c.Profile.Occupation, &c.Profile.Location, &c.Distance,
		); err != nil {
			return nil, fmt.Errorf("scanning candidate: %w", err)
		}
		candidates = append(candidates, c)
	}
	return candidates, rows.Err()
}
