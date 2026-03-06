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

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var user domain.User
	query := `SELECT * FROM users WHERE id = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var user domain.User
	query := `SELECT * FROM users WHERE email = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	query := `
		UPDATE users SET
			email = ?,
			password_hash = ?,
			email_verified_at = ?,
			totp_secret = ?,
			totp_enabled = ?,
			backup_codes = ?,
			last_login_at = ?,
			failed_login_attempts = ?,
			locked_until = ?,
			updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.PasswordHash, user.EmailVerifiedAt,
		user.TOTPSecret, user.TOTPEnabled, user.BackupCodes,
		user.LastLoginAt, user.FailedLoginAttempts, user.LockedUntil,
		time.Now(), user.ID,
	)
	return err
}

func (r *userRepository) SoftDelete(ctx context.Context, id string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET deleted_at = ?, updated_at = ? WHERE id = ?`,
		time.Now(), time.Now(), id,
	)
	return err
}

func (r *userRepository) IncrementFailedLogins(ctx context.Context, id string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET failed_login_attempts = failed_login_attempts + 1, updated_at = ? WHERE id = ?`,
		time.Now(), id,
	)
	return err
}

func (r *userRepository) ResetFailedLogins(ctx context.Context, id string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET failed_login_attempts = 0, locked_until = NULL, updated_at = ? WHERE id = ?`,
		time.Now(), id,
	)
	return err
}

func (r *userRepository) LockUntil(ctx context.Context, id string, until time.Time) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET locked_until = ?, updated_at = ? WHERE id = ?`,
		until, time.Now(), id,
	)
	return err
}

func (r *userRepository) SetEmailVerified(ctx context.Context, id string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET email_verified_at = ?, updated_at = ? WHERE id = ?`,
		now, now, id,
	)
	return err
}

func (r *userRepository) SetTOTPSecret(ctx context.Context, id, encryptedSecret, backupCodes string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET totp_secret = ?, backup_codes = ?, updated_at = ? WHERE id = ?`,
		encryptedSecret, backupCodes, time.Now(), id,
	)
	return err
}

func (r *userRepository) EnableTOTP(ctx context.Context, id string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET totp_enabled = TRUE, updated_at = ? WHERE id = ?`,
		time.Now(), id,
	)
	return err
}

func (r *userRepository) DisableTOTP(ctx context.Context, id string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET totp_enabled = FALSE, totp_secret = NULL, backup_codes = NULL, updated_at = ? WHERE id = ?`,
		time.Now(), id,
	)
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`,
		passwordHash, time.Now(), id,
	)
	return err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET last_login_at = ?, updated_at = ? WHERE id = ?`,
		now, now, id,
	)
	return err
}
