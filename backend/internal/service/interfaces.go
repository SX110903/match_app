package service

import (
	"context"
	"time"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
)

type IAuthService interface {
	Register(ctx context.Context, req RegisterRequest) error
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	LoginWith2FA(ctx context.Context, tempToken, code string) (*LoginResponse, error)
	Logout(ctx context.Context, accessJTI, refreshToken string, accessExpiry int64) error
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
	VerifyEmail(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	Setup2FA(ctx context.Context, userID string) (*Setup2FAResponse, error)
	Verify2FA(ctx context.Context, userID, code string) error
	Disable2FA(ctx context.Context, userID, password string) error
}

type IUserService interface {
	GetMe(ctx context.Context, userID string) (*UserProfileResponse, error)
	UpdateMe(ctx context.Context, userID string, req UpdateProfileRequest) error
	DeleteMe(ctx context.Context, userID string) error
	UpdatePreferences(ctx context.Context, userID string, req UpdatePreferencesRequest) error
}

type IMatchService interface {
	GetCandidates(ctx context.Context, userID string, page, limit int) ([]domain.Candidate, error)
	Swipe(ctx context.Context, userID, targetID string, direction domain.SwipeDirection) (*SwipeResponse, error)
	GetMatches(ctx context.Context, userID string) ([]domain.MatchWithProfile, error)
	GetMatch(ctx context.Context, userID, matchID string) (*domain.MatchWithProfile, error)
}

type IMessageService interface {
	GetMessages(ctx context.Context, userID, matchID string, page, limit int) ([]MessageResponse, error)
	SendMessage(ctx context.Context, userID, matchID, text string) (*MessageResponse, string, error)
	MarkRead(ctx context.Context, userID, matchID string) error
}

type IPhotoService interface {
	AddPhoto(ctx context.Context, userID, url string) (*PhotoResponse, error)
	DeletePhoto(ctx context.Context, userID, photoID string) error
}

type IPostService interface {
	GetFeed(ctx context.Context, viewerID string, page, limit int) ([]PostResponse, error)
	CreatePost(ctx context.Context, userID, content string, imageURL *string) (*PostResponse, error)
	DeletePost(ctx context.Context, userID, postID string) error
	LikePost(ctx context.Context, userID, postID string) error
	UnlikePost(ctx context.Context, userID, postID string) error
	GetComments(ctx context.Context, postID string) ([]CommentResponse, error)
	AddComment(ctx context.Context, userID, postID, content string) (*CommentResponse, error)
}

type INewsService interface {
	List(ctx context.Context, category string, adminView bool, page, limit int) ([]NewsArticleResponse, error)
	GetByID(ctx context.Context, id string) (*NewsArticleResponse, error)
	Create(ctx context.Context, authorID string, req CreateNewsRequest) (*NewsArticleResponse, error)
	Update(ctx context.Context, id string, req UpdateNewsRequest) (*NewsArticleResponse, error)
	Delete(ctx context.Context, id string) error
}

type IAdminService interface {
	AssertAdmin(ctx context.Context, userID string) error
	ListUsers(ctx context.Context, page, limit int) ([]repository.AdminUserSummary, error)
	FreezeUser(ctx context.Context, adminID, targetID string) error
	UnfreezeUser(ctx context.Context, adminID, targetID string) error
	SetVIPLevel(ctx context.Context, adminID, targetID string, level int) error
	AdjustCredits(ctx context.Context, adminID, targetID string, delta int) error
	SetAdmin(ctx context.Context, adminID, targetID string, isAdmin bool) error
	DeleteUser(ctx context.Context, adminID, targetID string) error
	GetAuditLog(ctx context.Context, adminID string, page, limit int) ([]domain.AdminLog, error)
	GetNotificationSettings(ctx context.Context, userID string) (*domain.NotificationSettings, error)
	SaveNotificationSettings(ctx context.Context, settings *domain.NotificationSettings) error
	GetPrivacySettings(ctx context.Context, userID string) (*domain.PrivacySettings, error)
	SavePrivacySettings(ctx context.Context, settings *domain.PrivacySettings) error
}

// --- Auth DTOs ---

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Name     string `json:"name"     validate:"required,min=2,max=50"`
	Age      int    `json:"age"      validate:"required,min=18,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"` // set as httpOnly cookie by handler, never in body
	Requires2FA  bool   `json:"requires_2fa,omitempty"`
	TempToken    string `json:"temp_token,omitempty"`
}

type Setup2FAResponse struct {
	QRCodeBase64 string   `json:"qr_code"`
	BackupCodes  []string `json:"backup_codes"`
}

// --- User DTOs ---

type UserProfileResponse struct {
	ID          string          `json:"id"`
	Email       string          `json:"email"`
	Name        string          `json:"name"`
	Age         int             `json:"age"`
	Bio         *string         `json:"bio"`
	Occupation  *string         `json:"occupation"`
	Location    *string         `json:"location"`
	Photos      []PhotoResponse `json:"photos"`
	Interests   []string        `json:"interests"`
	TOTPEnabled bool            `json:"totp_enabled"`
	IsAdmin     bool            `json:"is_admin"`
	IsFrozen    bool            `json:"is_frozen"`
	VIPLevel    int             `json:"vip_level"`
	Credits     int             `json:"credits"`
}

type UpdateProfileRequest struct {
	Name       *string  `json:"name"       validate:"omitempty,min=2,max=50"`
	Bio        *string  `json:"bio"        validate:"omitempty,max=500"`
	Occupation *string  `json:"occupation" validate:"omitempty,max=100"`
	Location   *string  `json:"location"   validate:"omitempty,max=100"`
	Interests  []string `json:"interests"  validate:"omitempty,max=10,dive,max=50"`
}

type UpdatePreferencesRequest struct {
	MinAge        int    `json:"min_age"         validate:"required,min=18,max=100"`
	MaxAge        int    `json:"max_age"         validate:"required,min=18,max=100"`
	MaxDistanceKm int    `json:"max_distance_km" validate:"required,min=1,max=500"`
	InterestedIn  string `json:"interested_in"   validate:"required,oneof=male female both"`
}

type SwipeResponse struct {
	IsMatch bool   `json:"is_match"`
	MatchID string `json:"match_id,omitempty"`
}

// --- Message DTOs ---

type MessageResponse struct {
	ID        string     `json:"id"`
	MatchID   string     `json:"match_id"`
	SenderID  string     `json:"sender_id"`
	Text      string     `json:"text"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// --- Photo DTOs ---

type AddPhotoRequest struct {
	URL string `json:"url" validate:"required,max=500"`
}

type PhotoResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// --- Post DTOs ---

type PostResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"content"`
	ImageURL     *string   `json:"image_url,omitempty"`
	LikesCount   int       `json:"likes_count"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	IsLikedByMe  bool      `json:"is_liked_by_me"`
	CreatedAt    time.Time `json:"created_at"`
}

type CommentResponse struct {
	ID           string    `json:"id"`
	PostID       string    `json:"post_id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"content"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreatePostRequest struct {
	Content  string  `json:"content"   validate:"required,min=1,max=2000"`
	ImageURL *string `json:"image_url" validate:"omitempty,max=500"`
}

type AddCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}

// --- News DTOs ---

type NewsArticleResponse struct {
	ID          string     `json:"id"`
	AuthorID    string     `json:"author_id"`
	AuthorName  string     `json:"author_name"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content"`
	ImageURL    *string    `json:"image_url,omitempty"`
	Category    string     `json:"category"`
	Published   bool       `json:"published"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateNewsRequest struct {
	Title     string  `json:"title"     validate:"required,min=3,max=255"`
	Summary   string  `json:"summary"   validate:"required,min=10,max=500"`
	Content   string  `json:"content"   validate:"required,min=10"`
	ImageURL  *string `json:"image_url" validate:"omitempty,max=500"`
	Category  string  `json:"category"  validate:"required,oneof=tendencias tech seguridad negocios general"`
	Published bool    `json:"published"`
}

type UpdateNewsRequest struct {
	Title     *string `json:"title"     validate:"omitempty,min=3,max=255"`
	Summary   *string `json:"summary"   validate:"omitempty,min=10,max=500"`
	Content   *string `json:"content"   validate:"omitempty,min=10"`
	ImageURL  *string `json:"image_url" validate:"omitempty,max=500"`
	Category  *string `json:"category"  validate:"omitempty,oneof=tendencias tech seguridad negocios general"`
	Published *bool   `json:"published"`
}

// --- Admin DTOs ---

type UserAdminActionRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type SetVIPRequest struct {
	UserID   string `json:"user_id"   validate:"required"`
	VIPLevel int    `json:"vip_level" validate:"min=0,max=5"`
}

type AdjustCreditsRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Delta  int    `json:"delta"`
}

// --- Settings DTOs ---

type NotificationSettingsRequest struct {
	NewMatches  bool `json:"new_matches"`
	NewMessages bool `json:"new_messages"`
	NewsUpdates bool `json:"news_updates"`
	Marketing   bool `json:"marketing"`
}

type PrivacySettingsRequest struct {
	ShowOnlineStatus bool `json:"show_online_status"`
	ShowLastSeen     bool `json:"show_last_seen"`
	ShowDistance     bool `json:"show_distance"`
	IncognitoMode    bool `json:"incognito_mode"`
}
