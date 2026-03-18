package repository

import (
	"context"
	"time"

	"github.com/SX110903/match_app/backend/internal/domain"
)

type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	SoftDelete(ctx context.Context, id string) error
	IncrementFailedLogins(ctx context.Context, id string) error
	ResetFailedLogins(ctx context.Context, id string) error
	LockUntil(ctx context.Context, id string, until time.Time) error
	SetEmailVerified(ctx context.Context, id string) error
	SetTOTPSecret(ctx context.Context, id string, encryptedSecret string, backupCodes string) error
	EnableTOTP(ctx context.Context, id string) error
	DisableTOTP(ctx context.Context, id string) error
	UpdatePassword(ctx context.Context, id string, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id string) error
}

type IProfileRepository interface {
	Create(ctx context.Context, profile *domain.UserProfile) error
	GetByUserID(ctx context.Context, userID string) (*domain.UserProfile, error)
	Update(ctx context.Context, profile *domain.UserProfile) error
	GetPreferences(ctx context.Context, userID string) (*domain.UserPreferences, error)
	UpsertPreferences(ctx context.Context, prefs *domain.UserPreferences) error
	AddPhoto(ctx context.Context, photo *domain.UserPhoto) error
	DeletePhoto(ctx context.Context, userID, photoID string) error
	GetPhotoCount(ctx context.Context, userID string) (int, error)
}

type ITokenRepository interface {
	// Email verification
	CreateEmailToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	GetEmailToken(ctx context.Context, tokenHash string) (userID string, expiresAt time.Time, used bool, err error)
	MarkEmailTokenUsed(ctx context.Context, tokenHash string) error

	// Password reset
	CreatePasswordResetToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (userID string, expiresAt time.Time, used bool, err error)
	MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error
	InvalidateAllPasswordResetTokens(ctx context.Context, userID string) error

	// Refresh tokens
	CreateRefreshToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenHash string) (userID string, expiresAt time.Time, err error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	DeleteAllRefreshTokens(ctx context.Context, userID string) error
}

type IMessageRepository interface {
	Create(ctx context.Context, msg *domain.Message) error
	GetByMatchID(ctx context.Context, matchID string, limit, offset int) ([]domain.Message, error)
	MarkAllRead(ctx context.Context, matchID, recipientID string) error
	GetUnreadCount(ctx context.Context, matchID, recipientID string) (int, error)
	GetLastMessage(ctx context.Context, matchID string) (*domain.Message, error)
}

type IMatchRepository interface {
	CreateSwipe(ctx context.Context, swipe *domain.Swipe) error
	GetSwipe(ctx context.Context, swiperID, swipedID string) (*domain.Swipe, error)
	CreateMatch(ctx context.Context, match *domain.Match) error
	GetMatchByID(ctx context.Context, id string) (*domain.Match, error)
	GetMatchesByUserID(ctx context.Context, userID string) ([]domain.MatchWithProfile, error)
	GetCandidates(ctx context.Context, userID string, prefs *domain.UserPreferences, limit, offset int) ([]domain.Candidate, error)
}
