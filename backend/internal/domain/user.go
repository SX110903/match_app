package domain

import "time"

type User struct {
	ID                  string     `db:"id"`
	Email               string     `db:"email"`
	PasswordHash        string     `db:"password_hash"`
	EmailVerifiedAt     *time.Time `db:"email_verified_at"`
	TOTPSecret          *string    `db:"totp_secret"`
	TOTPEnabled         bool       `db:"totp_enabled"`
	BackupCodes         *string    `db:"backup_codes"`
	LastLoginAt         *time.Time `db:"last_login_at"`
	FailedLoginAttempts int        `db:"failed_login_attempts"`
	LockedUntil         *time.Time `db:"locked_until"`
	DeletedAt           *time.Time `db:"deleted_at"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

type UserProfile struct {
	ID         string     `db:"id"`
	UserID     string     `db:"user_id"`
	Name       string     `db:"name"`
	Age        int        `db:"age"`
	Bio        *string    `db:"bio"`
	Occupation *string    `db:"occupation"`
	Location   *string    `db:"location"`
	Photos     []string   `db:"-"`
	Interests  []string   `db:"-"`
	Latitude   *float64   `db:"latitude"`
	Longitude  *float64   `db:"longitude"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

type UserPreferences struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	MinAge        int       `db:"min_age"`
	MaxAge        int       `db:"max_age"`
	MaxDistanceKm int       `db:"max_distance_km"`
	InterestedIn  string    `db:"interested_in"` // male, female, both
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
