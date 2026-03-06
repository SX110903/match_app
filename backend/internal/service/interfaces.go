package service

import (
	"context"

	"github.com/SX110903/match_app/backend/internal/domain"
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

// --- Request / Response DTOs ---

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
	Requires2FA  bool   `json:"requires_2fa,omitempty"`
	TempToken    string `json:"temp_token,omitempty"`
}

type Setup2FAResponse struct {
	QRCodeBase64 string   `json:"qr_code"`
	BackupCodes  []string `json:"backup_codes"`
}

type UserProfileResponse struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	Age        int      `json:"age"`
	Bio        *string  `json:"bio"`
	Occupation *string  `json:"occupation"`
	Location   *string  `json:"location"`
	Photos     []string `json:"photos"`
	Interests  []string `json:"interests"`
	TOTPEnabled bool    `json:"totp_enabled"`
}

type UpdateProfileRequest struct {
	Name       *string  `json:"name"       validate:"omitempty,min=2,max=50"`
	Bio        *string  `json:"bio"        validate:"omitempty,max=500"`
	Occupation *string  `json:"occupation" validate:"omitempty,max=100"`
	Location   *string  `json:"location"   validate:"omitempty,max=100"`
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
