package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

func splitNonEmpty(s string, sep rune) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == sep })
	return parts
}

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
		`INSERT IGNORE INTO swipes (id, swiper_id, swiped_id, direction, created_at) VALUES (?, ?, ?, ?, ?)`,
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
		`INSERT IGNORE INTO matches (id, user1_id, user2_id, created_at) VALUES (?, ?, ?, ?)`,
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

// GetMatchesByUserID devuelve todos los matches de un usuario con el perfil del otro participante.
// TODO(perf-v2): las subqueries last_message y unread_count son correlacionadas (N+1).
// Para escalar, desnormalizar en columnas last_message_text/last_message_at/unread_u1/unread_u2
// en la tabla matches y actualizar tras cada INSERT en messages.
func (r *matchRepository) GetMatchesByUserID(ctx context.Context, userID string) ([]domain.MatchWithProfile, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	query := `
		SELECT
			m.id, m.user1_id, m.user2_id, m.created_at,
			COALESCE(up.id, '')          as profile_id,
			COALESCE(up.user_id, '')     as profile_user_id,
			COALESCE(up.name, '')        as name,
			COALESCE(up.age, 0)          as age,
			up.bio, up.occupation, up.location,
			(SELECT text FROM messages WHERE match_id = m.id ORDER BY created_at DESC LIMIT 1) as last_message,
			(SELECT COUNT(*) FROM messages WHERE match_id = m.id AND sender_id != ? AND read_at IS NULL) as unread_count,
			(SELECT GROUP_CONCAT(url ORDER BY sort_order SEPARATOR '|') FROM user_photos WHERE user_id = up.user_id) as photos_str
		FROM matches m
		LEFT JOIN user_profiles up ON up.user_id = CASE
			WHEN m.user1_id = ? THEN m.user2_id
			ELSE m.user1_id
		END
		WHERE (m.user1_id = ? OR m.user2_id = ?)
		  AND m.deleted_at IS NULL
		ORDER BY m.created_at DESC`

	rows, err := r.db.QueryxContext(ctx, query, userID, userID, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("getting matches: %w", err)
	}
	defer rows.Close()

	var results []domain.MatchWithProfile
	for rows.Next() {
		var mwp domain.MatchWithProfile
		var photosStr *string
		if err := rows.Scan(
			&mwp.Match.ID, &mwp.Match.User1ID, &mwp.Match.User2ID, &mwp.Match.CreatedAt,
			&mwp.Profile.ID, &mwp.Profile.UserID, &mwp.Profile.Name, &mwp.Profile.Age,
			&mwp.Profile.Bio, &mwp.Profile.Occupation, &mwp.Profile.Location,
			&mwp.LastMessage, &mwp.UnreadCount, &photosStr,
		); err != nil {
			return nil, fmt.Errorf("scanning match: %w", err)
		}
		if photosStr != nil && *photosStr != "" {
			for _, url := range splitNonEmpty(*photosStr, '|') {
				mwp.Profile.Photos = append(mwp.Profile.Photos, url)
			}
		}
		results = append(results, mwp)
	}
	return results, rows.Err()
}

func (r *matchRepository) DeleteMatch(ctx context.Context, matchID, userID string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		`UPDATE matches SET deleted_at = NOW()
		 WHERE id = ? AND (user1_id = ? OR user2_id = ?) AND deleted_at IS NULL`,
		matchID, userID, userID,
	)
	if err != nil {
		return fmt.Errorf("deleting match: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *matchRepository) GetCandidates(ctx context.Context, userID string, prefs *domain.UserPreferences, limit, offset int) ([]domain.Candidate, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	// When either the current user or a candidate lacks coordinates,
	// ST_Distance_Sphere returns NULL and HAVING filters out the row.
	// Use COALESCE so that missing coordinates result in 0 distance (include everyone).
	query := `
		SELECT
			up.id, up.user_id, up.name, up.age, up.bio, up.occupation, up.location,
			COALESCE(
				ST_Distance_Sphere(
					POINT(up.longitude, up.latitude),
					(SELECT POINT(longitude, latitude) FROM user_profiles WHERE user_id = ?)
				) / 1000,
				0
			) as distance,
			(SELECT GROUP_CONCAT(url ORDER BY sort_order SEPARATOR '|')
			 FROM user_photos WHERE user_id = up.user_id) as photos_str,
			(SELECT GROUP_CONCAT(interest SEPARATOR ',')
			 FROM user_interests WHERE user_id = up.user_id) as interests_str
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
		var photosStr, interestsStr *string
		if err := rows.Scan(
			&c.Profile.ID, &c.Profile.UserID, &c.Profile.Name, &c.Profile.Age,
			&c.Profile.Bio, &c.Profile.Occupation, &c.Profile.Location, &c.Distance,
			&photosStr, &interestsStr,
		); err != nil {
			return nil, fmt.Errorf("scanning candidate: %w", err)
		}
		if photosStr != nil && *photosStr != "" {
			for _, url := range splitNonEmpty(*photosStr, '|') {
				c.Profile.Photos = append(c.Profile.Photos, url)
			}
		}
		if interestsStr != nil && *interestsStr != "" {
			c.Profile.Interests = splitNonEmpty(*interestsStr, ',')
		}
		candidates = append(candidates, c)
	}
	return candidates, rows.Err()
}
