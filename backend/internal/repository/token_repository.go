package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type tokenRepository struct {
	db *database.DB
}

func NewTokenRepository(db *database.DB) ITokenRepository {
	return &tokenRepository{db: db}
}

// Email verification tokens

func (r *tokenRepository) CreateEmailToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO email_verification_tokens (user_id, token_hash, expires_at, created_at) VALUES (?, ?, ?, ?)`,
		userID, tokenHash, expiresAt, time.Now(),
	)
	return err
}

func (r *tokenRepository) GetEmailToken(ctx context.Context, tokenHash string) (string, time.Time, bool, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var row struct {
		UserID    string    `db:"user_id"`
		ExpiresAt time.Time `db:"expires_at"`
		UsedAt    *time.Time `db:"used_at"`
	}
	err := r.db.GetContext(ctx, &row,
		`SELECT user_id, expires_at, used_at FROM email_verification_tokens WHERE token_hash = ?`,
		tokenHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, false, domain.ErrNotFound
		}
		return "", time.Time{}, false, fmt.Errorf("getting email token: %w", err)
	}
	return row.UserID, row.ExpiresAt, row.UsedAt != nil, nil
}

func (r *tokenRepository) MarkEmailTokenUsed(ctx context.Context, tokenHash string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE email_verification_tokens SET used_at = ? WHERE token_hash = ?`,
		time.Now(), tokenHash,
	)
	return err
}

// Password reset tokens

func (r *tokenRepository) CreatePasswordResetToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO password_reset_tokens (user_id, token_hash, expires_at, created_at) VALUES (?, ?, ?, ?)`,
		userID, tokenHash, expiresAt, time.Now(),
	)
	return err
}

func (r *tokenRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (string, time.Time, bool, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var row struct {
		UserID    string     `db:"user_id"`
		ExpiresAt time.Time  `db:"expires_at"`
		UsedAt    *time.Time `db:"used_at"`
	}
	err := r.db.GetContext(ctx, &row,
		`SELECT user_id, expires_at, used_at FROM password_reset_tokens WHERE token_hash = ?`,
		tokenHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, false, domain.ErrNotFound
		}
		return "", time.Time{}, false, fmt.Errorf("getting password reset token: %w", err)
	}
	return row.UserID, row.ExpiresAt, row.UsedAt != nil, nil
}

func (r *tokenRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE password_reset_tokens SET used_at = ? WHERE token_hash = ?`,
		time.Now(), tokenHash,
	)
	return err
}

func (r *tokenRepository) InvalidateAllPasswordResetTokens(ctx context.Context, userID string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE password_reset_tokens SET used_at = ? WHERE user_id = ? AND used_at IS NULL`,
		time.Now(), userID,
	)
	return err
}

// Refresh tokens

func (r *tokenRepository) CreateRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) VALUES (?, ?, ?, ?)`,
		userID, tokenHash, expiresAt, time.Now(),
	)
	return err
}

func (r *tokenRepository) GetRefreshToken(ctx context.Context, tokenHash string) (string, time.Time, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var row struct {
		UserID    string    `db:"user_id"`
		ExpiresAt time.Time `db:"expires_at"`
	}
	err := r.db.GetContext(ctx, &row,
		`SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = ? AND revoked_at IS NULL`,
		tokenHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, domain.ErrNotFound
		}
		return "", time.Time{}, fmt.Errorf("getting refresh token: %w", err)
	}
	return row.UserID, row.ExpiresAt, nil
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = ? WHERE token_hash = ?`,
		time.Now(), tokenHash,
	)
	return err
}

func (r *tokenRepository) DeleteAllRefreshTokens(ctx context.Context, userID string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`,
		time.Now(), userID,
	)
	return err
}
