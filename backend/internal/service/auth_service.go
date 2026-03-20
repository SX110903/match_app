package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/png"
	"time"

	"github.com/google/uuid"
	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/config"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/email"
	"github.com/SX110903/match_app/backend/internal/repository"
	"github.com/SX110903/match_app/backend/pkg/logger"
)

type authService struct {
	userRepo    repository.IUserRepository
	profileRepo repository.IProfileRepository
	tokenRepo   repository.ITokenRepository
	emailSvc    email.IEmailService
	jwtSvc      auth.IJWTService
	totpSvc     auth.ITOTPService
	blacklist   auth.ITokenBlacklist
	cfg         *config.Config
}

func NewAuthService(
	userRepo repository.IUserRepository,
	profileRepo repository.IProfileRepository,
	tokenRepo repository.ITokenRepository,
	emailSvc email.IEmailService,
	jwtSvc auth.IJWTService,
	totpSvc auth.ITOTPService,
	blacklist auth.ITokenBlacklist,
	cfg *config.Config,
) IAuthService {
	return &authService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		tokenRepo:   tokenRepo,
		emailSvc:    emailSvc,
		jwtSvc:      jwtSvc,
		totpSvc:     totpSvc,
		blacklist:   blacklist,
		cfg:         cfg,
	}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) error {
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return domain.ErrConflict
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		// MySQL 1062: Duplicate entry — concurrent registration race
		if repository.IsDuplicateEntry(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("creating user: %w", err)
	}

	profile := &domain.UserProfile{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Name:      req.Name,
		Age:       req.Age,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.profileRepo.Create(ctx, profile); err != nil {
		return fmt.Errorf("creating profile: %w", err)
	}

	// Send verification email
	token, tokenHash := generateSecureToken()
	expiry := time.Now().Add(24 * time.Hour)
	if err := s.tokenRepo.CreateEmailToken(ctx, user.ID, tokenHash, expiry); err != nil {
		return fmt.Errorf("creating email token: %w", err)
	}

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.cfg.Security.FrontendURL, token)
	userID := user.ID
	userName := req.Name
	userEmail := req.Email
	go func() {
		if err := s.emailSvc.SendVerificationEmail(context.Background(), userEmail, userName, verifyURL); err != nil {
			logger.Error().Err(err).Str("user_id", userID).Msg("verification email failed")
		}
	}()

	return nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user.IsDeleted() {
		return nil, domain.ErrInvalidCredentials
	}

	if user.IsLocked() {
		return nil, domain.ErrAccountLocked
	}

	ok, err := auth.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !ok {
		_ = s.userRepo.IncrementFailedLogins(ctx, user.ID)
		// Lock after 5 failed attempts
		if user.FailedLoginAttempts+1 >= 5 {
			_ = s.userRepo.LockUntil(ctx, user.ID, time.Now().Add(15*time.Minute))
		}
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsEmailVerified() {
		return nil, domain.ErrEmailNotVerified
	}

	_ = s.userRepo.ResetFailedLogins(ctx, user.ID)
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// If 2FA is enabled, return temp token
	if user.TOTPEnabled {
		tempToken, err := s.jwtSvc.GenerateTempToken(user.ID)
		if err != nil {
			return nil, fmt.Errorf("generating temp token: %w", err)
		}
		return &LoginResponse{Requires2FA: true, TempToken: tempToken}, nil
	}

	return s.issueTokens(ctx, user.ID)
}

func (s *authService) LoginWith2FA(ctx context.Context, tempToken, code string) (*LoginResponse, error) {
	claims, err := s.jwtSvc.ValidateTempToken(tempToken)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	user, err := s.userRepo.GetByID(ctx, claims.Subject)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if !user.TOTPEnabled || user.TOTPSecret == nil {
		return nil, domain.ErrTwoFANotEnabled
	}

	secret, err := auth.Decrypt(*user.TOTPSecret, s.cfg.Security.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("decrypting TOTP secret: %w", err)
	}

	if !s.totpSvc.Validate(code, secret) {
		return nil, domain.ErrTwoFAInvalid
	}

	return s.issueTokens(ctx, user.ID)
}

func (s *authService) issueTokens(ctx context.Context, userID string) (*LoginResponse, error) {
	accessToken, _, err := s.jwtSvc.GenerateAccessToken(userID, []string{"user"}, true)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := s.jwtSvc.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	tokenHash := hashToken(refreshToken)
	expiry := time.Now().Add(s.cfg.JWT.RefreshTokenExpiry)
	if err := s.tokenRepo.CreateRefreshToken(ctx, userID, tokenHash, expiry); err != nil {
		return nil, fmt.Errorf("storing refresh token: %w", err)
	}

	return &LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *authService) Logout(ctx context.Context, accessJTI, refreshToken string, accessExpiry int64) error {
	remaining := time.Until(time.Unix(accessExpiry, 0))
	if remaining > 0 {
		_ = s.blacklist.Add(ctx, accessJTI, remaining)
	}
	if refreshToken != "" {
		_ = s.tokenRepo.DeleteRefreshToken(ctx, hashToken(refreshToken))
	}
	return nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	tokenHash := hashToken(refreshToken)
	userID, expiry, err := s.tokenRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}
	if time.Now().After(expiry) {
		_ = s.tokenRepo.DeleteRefreshToken(ctx, tokenHash)
		return nil, domain.ErrTokenExpired
	}

	// Rotate: delete old, issue new
	_ = s.tokenRepo.DeleteRefreshToken(ctx, tokenHash)
	return s.issueTokens(ctx, userID)
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	tokenHash := hashToken(token)
	userID, expiry, used, err := s.tokenRepo.GetEmailToken(ctx, tokenHash)
	if err != nil {
		return domain.ErrTokenInvalid
	}
	if used {
		return domain.ErrTokenInvalid
	}
	if time.Now().After(expiry) {
		return domain.ErrTokenExpired
	}

	if err := s.userRepo.SetEmailVerified(ctx, userID); err != nil {
		return fmt.Errorf("setting email verified: %w", err)
	}
	return s.tokenRepo.MarkEmailTokenUsed(ctx, tokenHash)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Always respond the same - prevent user enumeration
		return nil
	}

	token, tokenHash := generateSecureToken()
	expiry := time.Now().Add(1 * time.Hour)
	if err := s.tokenRepo.CreatePasswordResetToken(ctx, user.ID, tokenHash, expiry); err != nil {
		return nil // Still silent
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.Security.FrontendURL, token)
	_ = s.emailSvc.SendPasswordResetEmail(ctx, email, resetURL)
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	tokenHash := hashToken(token)
	userID, expiry, used, err := s.tokenRepo.GetPasswordResetToken(ctx, tokenHash)
	if err != nil || used || time.Now().After(expiry) {
		return domain.ErrTokenInvalid
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, hash); err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	_ = s.tokenRepo.MarkPasswordResetTokenUsed(ctx, tokenHash)
	_ = s.tokenRepo.DeleteAllRefreshTokens(ctx, userID)

	user, _ := s.userRepo.GetByID(ctx, userID)
	if user != nil {
		_ = s.emailSvc.SendPasswordChangedEmail(ctx, user.Email)
	}
	return nil
}

func (s *authService) Setup2FA(ctx context.Context, userID string) (*Setup2FAResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if user.TOTPEnabled {
		return nil, domain.ErrTwoFAAlreadyEnabled
	}

	key, err := s.totpSvc.GenerateSecret(user.Email)
	if err != nil {
		return nil, fmt.Errorf("generating TOTP key: %w", err)
	}

	// Generate QR code as base64 PNG
	img, err := key.Image(256, 256)
	if err != nil {
		return nil, fmt.Errorf("generating QR image: %w", err)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encoding QR PNG: %w", err)
	}
	qrBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	// Generate backup codes
	backupCodes, err := s.totpSvc.GenerateBackupCodes(10)
	if err != nil {
		return nil, fmt.Errorf("generating backup codes: %w", err)
	}

	// Encrypt and store secret (not yet enabled - user must verify first)
	encSecret, err := auth.Encrypt(key.Secret(), s.cfg.Security.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("encrypting TOTP secret: %w", err)
	}

	backupJSON, _ := json.Marshal(backupCodes)
	encBackup, err := auth.Encrypt(string(backupJSON), s.cfg.Security.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("encrypting backup codes: %w", err)
	}

	if err := s.userRepo.SetTOTPSecret(ctx, userID, encSecret, encBackup); err != nil {
		return nil, fmt.Errorf("storing TOTP secret: %w", err)
	}

	return &Setup2FAResponse{QRCodeBase64: qrBase64, BackupCodes: backupCodes}, nil
}

func (s *authService) Verify2FA(ctx context.Context, userID, code string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrNotFound
	}
	if user.TOTPSecret == nil {
		return domain.ErrTwoFANotEnabled
	}

	secret, err := auth.Decrypt(*user.TOTPSecret, s.cfg.Security.EncryptionKey)
	if err != nil {
		return fmt.Errorf("decrypting TOTP secret: %w", err)
	}

	if !s.totpSvc.Validate(code, secret) {
		return domain.ErrTwoFAInvalid
	}

	return s.userRepo.EnableTOTP(ctx, userID)
}

func (s *authService) Disable2FA(ctx context.Context, userID, password string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrNotFound
	}

	ok, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil || !ok {
		return domain.ErrInvalidCredentials
	}

	if !user.TOTPEnabled {
		return domain.ErrTwoFANotEnabled
	}

	return s.userRepo.DisableTOTP(ctx, userID)
}

// generateSecureToken creates a cryptographically secure random token.
// Returns the raw token (to send to user) and its SHA256 hash (to store in DB).
func generateSecureToken() (token, hash string) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	token = hex.EncodeToString(b)
	hash = hashToken(token)
	return token, hash
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
