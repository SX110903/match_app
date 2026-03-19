package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/database"
)

type followRepository struct {
	db *database.DB
}

func NewFollowRepository(db *database.DB) IFollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Follow(ctx context.Context, followerID, followedID string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		`INSERT IGNORE INTO follows (follower_id, followed_id) VALUES (?, ?)`,
		followerID, followedID,
	)
	if err != nil {
		return fmt.Errorf("follow insert: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		// Increment follower count only when a new row was inserted
		_, err = r.db.ExecContext(ctx,
			`UPDATE users SET follower_count = follower_count + 1 WHERE id = ?`,
			followedID,
		)
		if err != nil {
			return fmt.Errorf("follow increment: %w", err)
		}
	}
	return nil
}

func (r *followRepository) Unfollow(ctx context.Context, followerID, followedID string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	res, err := r.db.ExecContext(ctx,
		`DELETE FROM follows WHERE follower_id = ? AND followed_id = ?`,
		followerID, followedID,
	)
	if err != nil {
		return fmt.Errorf("unfollow delete: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		_, err = r.db.ExecContext(ctx,
			`UPDATE users SET follower_count = GREATEST(follower_count - 1, 0) WHERE id = ?`,
			followedID,
		)
		if err != nil {
			return fmt.Errorf("unfollow decrement: %w", err)
		}
	}
	return nil
}

func (r *followRepository) GetFollowerCount(ctx context.Context, userID string) (int, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT follower_count FROM users WHERE id = ?`, userID,
	).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("get follower count: %w", err)
	}
	return count, nil
}

func (r *followRepository) GetFollowingCount(ctx context.Context, userID string) (int, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM follows WHERE follower_id = ?`, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get following count: %w", err)
	}
	return count, nil
}

func (r *followRepository) GetFollowers(ctx context.Context, userID string, limit, offset int) ([]string, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(ctx,
		`SELECT follower_id FROM follows WHERE followed_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get followers: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *followRepository) GetFollowing(ctx context.Context, userID string, limit, offset int) ([]string, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(ctx,
		`SELECT followed_id FROM follows WHERE follower_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get following: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *followRepository) IsFollowing(ctx context.Context, followerID, followedID string) (bool, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM follows WHERE follower_id = ? AND followed_id = ?`,
		followerID, followedID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("is following: %w", err)
	}
	return count > 0, nil
}
