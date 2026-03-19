package repository

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type adminRepository struct{ db *database.DB }

func NewAdminRepository(db *database.DB) IAdminRepository {
	return &adminRepository{db: db}
}

type AdminUserSummary struct {
	ID       string `db:"id"       json:"id"`
	Email    string `db:"email"    json:"email"`
	Name     string `db:"name"     json:"name"`
	IsAdmin  bool   `db:"is_admin" json:"is_admin"`
	IsFrozen bool   `db:"is_frozen" json:"is_frozen"`
	VIPLevel int    `db:"vip_level" json:"vip_level"`
	Credits  int    `db:"credits"  json:"credits"`
}

func (r *adminRepository) ListUsers(ctx context.Context, limit, offset int) ([]AdminUserSummary, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT u.id, u.email, COALESCE(up.name, u.email) AS name,
		       u.is_admin, u.is_frozen, u.vip_level, u.credits
		FROM users u
		LEFT JOIN user_profiles up ON up.user_id = u.id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	defer rows.Close()

	var users []AdminUserSummary
	for rows.Next() {
		var u AdminUserSummary
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.IsAdmin, &u.IsFrozen, &u.VIPLevel, &u.Credits); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *adminRepository) SetFrozen(ctx context.Context, userID string, frozen bool) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `UPDATE users SET is_frozen = ? WHERE id = ?`, frozen, userID)
	return err
}

func (r *adminRepository) SetVIPLevel(ctx context.Context, userID string, level int) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `UPDATE users SET vip_level = ? WHERE id = ?`, level, userID)
	return err
}

func (r *adminRepository) AddCredits(ctx context.Context, userID string, delta int) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET credits = GREATEST(credits + ?, 0) WHERE id = ?`, delta, userID)
	return err
}

func (r *adminRepository) SetAdmin(ctx context.Context, userID string, admin bool) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx, `UPDATE users SET is_admin = ? WHERE id = ?`, admin, userID)
	return err
}

func (r *adminRepository) LogAction(ctx context.Context, log *domain.AdminLog) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO admin_logs (id, admin_id, target_id, action, details) VALUES (?, ?, ?, ?, ?)`,
		log.ID, log.AdminID, log.TargetID, log.Action, log.Details,
	)
	return err
}

func (r *adminRepository) GetNotificationSettings(ctx context.Context, userID string) (*domain.NotificationSettings, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var s domain.NotificationSettings
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, new_matches, new_messages, news_updates, marketing
		 FROM user_notification_settings WHERE user_id = ?`, userID,
	).Scan(&s.UserID, &s.NewMatches, &s.NewMessages, &s.NewsUpdates, &s.Marketing)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *adminRepository) UpsertNotificationSettings(ctx context.Context, s *domain.NotificationSettings) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_notification_settings (user_id, new_matches, new_messages, news_updates, marketing)
		 VALUES (?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE new_matches=?, new_messages=?, news_updates=?, marketing=?`,
		s.UserID, s.NewMatches, s.NewMessages, s.NewsUpdates, s.Marketing,
		s.NewMatches, s.NewMessages, s.NewsUpdates, s.Marketing,
	)
	return err
}

func (r *adminRepository) GetPrivacySettings(ctx context.Context, userID string) (*domain.PrivacySettings, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var s domain.PrivacySettings
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, show_online_status, show_last_seen, show_distance, incognito_mode
		 FROM user_privacy_settings WHERE user_id = ?`, userID,
	).Scan(&s.UserID, &s.ShowOnlineStatus, &s.ShowLastSeen, &s.ShowDistance, &s.IncognitoMode)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *adminRepository) UpsertPrivacySettings(ctx context.Context, s *domain.PrivacySettings) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_privacy_settings (user_id, show_online_status, show_last_seen, show_distance, incognito_mode)
		 VALUES (?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE show_online_status=?, show_last_seen=?, show_distance=?, incognito_mode=?`,
		s.UserID, s.ShowOnlineStatus, s.ShowLastSeen, s.ShowDistance, s.IncognitoMode,
		s.ShowOnlineStatus, s.ShowLastSeen, s.ShowDistance, s.IncognitoMode,
	)
	return err
}
