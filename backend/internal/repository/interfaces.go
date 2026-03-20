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
	UpdateBadge(ctx context.Context, userID string, badge string) error
	UpdateCredits(ctx context.Context, userID string, delta int) error
}

type IFollowRepository interface {
	Follow(ctx context.Context, followerID, followedID string) error
	Unfollow(ctx context.Context, followerID, followedID string) error
	GetFollowerCount(ctx context.Context, userID string) (int, error)
	GetFollowingCount(ctx context.Context, userID string) (int, error)
	GetFollowers(ctx context.Context, userID string, limit, offset int) ([]string, error)
	GetFollowing(ctx context.Context, userID string, limit, offset int) ([]string, error)
	IsFollowing(ctx context.Context, followerID, followedID string) (bool, error)
}

type IShopRepository interface {
	CreateTransaction(ctx context.Context, tx *domain.ShopTransaction) error
	GetTransactionsByUser(ctx context.Context, userID string, limit, offset int) ([]domain.ShopTransaction, error)
	PurchaseVIP(ctx context.Context, userID string, itemValue, cost int) error
}

type IAdRepository interface {
	GetActive(ctx context.Context, userBadge string) (*domain.Ad, error)
	RegisterClick(ctx context.Context, adID, userID string) error
	Create(ctx context.Context, ad *domain.Ad) error
	Update(ctx context.Context, ad *domain.Ad) error
	Delete(ctx context.Context, id string) error
	Toggle(ctx context.Context, id string) error
	ListAll(ctx context.Context) ([]domain.Ad, error)
	GetByID(ctx context.Context, id string) (*domain.Ad, error)
	IncrementImpressions(ctx context.Context, id string) error
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
	ReplaceInterests(ctx context.Context, userID string, interests []string) error
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
	DeleteMatch(ctx context.Context, matchID, userID string) error
}

type IPostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	GetFeed(ctx context.Context, viewerID string, limit, offset int) ([]domain.Post, error)
	GetByID(ctx context.Context, postID, viewerID string) (*domain.Post, error)
	Delete(ctx context.Context, postID, userID string) error
	LikePost(ctx context.Context, like *domain.PostLike) error
	UnlikePost(ctx context.Context, postID, userID string) error
	GetComments(ctx context.Context, postID string) ([]domain.PostComment, error)
	AddComment(ctx context.Context, comment *domain.PostComment) (*domain.PostComment, error)
}

type INewsRepository interface {
	Create(ctx context.Context, article *domain.NewsArticle) error
	Update(ctx context.Context, article *domain.NewsArticle) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.NewsArticle, error)
	List(ctx context.Context, category string, publishedOnly bool, limit, offset int) ([]domain.NewsArticle, error)
}

type IAdminRepository interface {
	ListUsers(ctx context.Context, limit, offset int) ([]AdminUserSummary, error)
	SetFrozen(ctx context.Context, userID string, frozen bool) error
	SetVIPLevel(ctx context.Context, userID string, level int) error
	AddCredits(ctx context.Context, userID string, delta int) error
	SetAdmin(ctx context.Context, userID string, admin bool) error
	LogAction(ctx context.Context, log *domain.AdminLog) error
	GetAuditLog(ctx context.Context, limit, offset int) ([]domain.AdminLog, error)
	GetNotificationSettings(ctx context.Context, userID string) (*domain.NotificationSettings, error)
	UpsertNotificationSettings(ctx context.Context, s *domain.NotificationSettings) error
	GetPrivacySettings(ctx context.Context, userID string) (*domain.PrivacySettings, error)
	UpsertPrivacySettings(ctx context.Context, s *domain.PrivacySettings) error
}
