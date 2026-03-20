package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/SX110903/match_app/backend/internal/auth"
	"github.com/SX110903/match_app/backend/internal/config"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/service"
)

// ---- mock IUserRepository ----

type mockUserRepo struct {
	users        map[string]*domain.User // key: email
	byID         map[string]*domain.User
	failedLogins map[string]int
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:        make(map[string]*domain.User),
		byID:         make(map[string]*domain.User),
		failedLogins: make(map[string]int),
	}
}

func (m *mockUserRepo) Create(_ context.Context, u *domain.User) error {
	if _, exists := m.users[u.Email]; exists {
		return domain.ErrConflict
	}
	m.users[u.Email] = u
	m.byID[u.ID] = u
	return nil
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id string) (*domain.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) Update(_ context.Context, u *domain.User) error {
	m.users[u.Email] = u
	m.byID[u.ID] = u
	return nil
}

func (m *mockUserRepo) SoftDelete(_ context.Context, _ string) error               { return nil }
func (m *mockUserRepo) SetEmailVerified(_ context.Context, _ string) error          { return nil }
func (m *mockUserRepo) UpdateLastLogin(_ context.Context, _ string) error           { return nil }
func (m *mockUserRepo) UpdatePassword(_ context.Context, _, _ string) error         { return nil }
func (m *mockUserRepo) UpdateBadge(_ context.Context, _, _ string) error            { return nil }
func (m *mockUserRepo) UpdateCredits(_ context.Context, _ string, _ int) error      { return nil }
func (m *mockUserRepo) EnableTOTP(_ context.Context, _ string) error                { return nil }
func (m *mockUserRepo) DisableTOTP(_ context.Context, _ string) error               { return nil }
func (m *mockUserRepo) SetTOTPSecret(_ context.Context, _, _, _ string) error       { return nil }
func (m *mockUserRepo) ResetFailedLogins(_ context.Context, id string) error {
	m.failedLogins[id] = 0
	if u, ok := m.byID[id]; ok {
		u.FailedLoginAttempts = 0
	}
	return nil
}

func (m *mockUserRepo) IncrementFailedLogins(_ context.Context, id string) error {
	m.failedLogins[id]++
	if u, ok := m.byID[id]; ok {
		u.FailedLoginAttempts = m.failedLogins[id]
	}
	return nil
}

func (m *mockUserRepo) LockUntil(_ context.Context, id string, until time.Time) error {
	if u, ok := m.byID[id]; ok {
		u.LockedUntil = &until
	}
	return nil
}

// ---- mock IProfileRepository ----

type mockProfileRepo struct{}

func (m *mockProfileRepo) Create(_ context.Context, _ *domain.UserProfile) error { return nil }
func (m *mockProfileRepo) Update(_ context.Context, _ *domain.UserProfile) error { return nil }
func (m *mockProfileRepo) GetByUserID(_ context.Context, _ string) (*domain.UserProfile, error) {
	return &domain.UserProfile{}, nil
}
func (m *mockProfileRepo) GetPreferences(_ context.Context, _ string) (*domain.UserPreferences, error) {
	return &domain.UserPreferences{}, nil
}
func (m *mockProfileRepo) UpsertPreferences(_ context.Context, _ *domain.UserPreferences) error {
	return nil
}
func (m *mockProfileRepo) AddPhoto(_ context.Context, _ *domain.UserPhoto) error { return nil }
func (m *mockProfileRepo) DeletePhoto(_ context.Context, _, _ string) error      { return nil }
func (m *mockProfileRepo) GetPhotoCount(_ context.Context, _ string) (int, error) {
	return 0, nil
}
func (m *mockProfileRepo) ReplaceInterests(_ context.Context, _ string, _ []string) error {
	return nil
}

// ---- mock ITokenRepository ----

type mockTokenRepo struct{}

func (m *mockTokenRepo) CreateEmailToken(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}
func (m *mockTokenRepo) GetEmailToken(_ context.Context, _ string) (string, time.Time, bool, error) {
	return "", time.Time{}, false, errors.New("not found")
}
func (m *mockTokenRepo) MarkEmailTokenUsed(_ context.Context, _ string) error { return nil }
func (m *mockTokenRepo) CreatePasswordResetToken(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}
func (m *mockTokenRepo) GetPasswordResetToken(_ context.Context, _ string) (string, time.Time, bool, error) {
	return "", time.Time{}, false, errors.New("not found")
}
func (m *mockTokenRepo) MarkPasswordResetTokenUsed(_ context.Context, _ string) error  { return nil }
func (m *mockTokenRepo) InvalidateAllPasswordResetTokens(_ context.Context, _ string) error {
	return nil
}
func (m *mockTokenRepo) CreateRefreshToken(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}
func (m *mockTokenRepo) GetRefreshToken(_ context.Context, _ string) (string, time.Time, error) {
	return "", time.Time{}, errors.New("not found")
}
func (m *mockTokenRepo) DeleteRefreshToken(_ context.Context, _ string) error      { return nil }
func (m *mockTokenRepo) DeleteAllRefreshTokens(_ context.Context, _ string) error  { return nil }

// ---- mock IEmailService ----

type mockEmailSvc struct{ sendErr error }

func (m *mockEmailSvc) SendVerificationEmail(_ context.Context, _, _, _ string) error {
	return m.sendErr
}
func (m *mockEmailSvc) SendPasswordResetEmail(_ context.Context, _, _ string) error {
	return m.sendErr
}
func (m *mockEmailSvc) SendPasswordChangedEmail(_ context.Context, _ string) error { return nil }
func (m *mockEmailSvc) SendWelcomeEmail(_ context.Context, _, _ string) error      { return nil }

// ---- mock IJWTService ----

type mockJWTSvc struct{}

func (m *mockJWTSvc) GenerateAccessToken(_ string, _ []string, _ bool) (string, string, error) {
	return "access_tok", "jti-123", nil
}
func (m *mockJWTSvc) GenerateRefreshToken(_ string) (string, error) { return "refresh_tok", nil }
func (m *mockJWTSvc) GenerateTempToken(_ string) (string, error)    { return "temp_tok", nil }
func (m *mockJWTSvc) ValidateAccessToken(_ string) (*auth.Claims, error) {
	return nil, errors.New("invalid")
}
func (m *mockJWTSvc) ValidateTempToken(_ string) (*auth.Claims, error) {
	return nil, errors.New("invalid")
}

// ---- mock ITOTPService ----

type mockTOTPSvc struct{}

func (m *mockTOTPSvc) GenerateSecret(_ string) (*otp.Key, error)   { return nil, nil }
func (m *mockTOTPSvc) Validate(_, _ string) bool                   { return true }
func (m *mockTOTPSvc) GenerateBackupCodes(_ int) ([]string, error) { return nil, nil }

// ---- mock ITokenBlacklist ----

type mockBlacklist struct{}

func (m *mockBlacklist) Add(_ context.Context, _ string, _ time.Duration) error { return nil }
func (m *mockBlacklist) IsBlacklisted(_ context.Context, _ string) (bool, error) {
	return false, nil
}

// ---- helper ----

func newAuthSvc(userRepo *mockUserRepo, emailErr error) service.IAuthService {
	cfg := &config.Config{}
	cfg.Security.FrontendURL = "http://localhost:3000"
	return service.NewAuthService(
		userRepo,
		&mockProfileRepo{},
		&mockTokenRepo{},
		&mockEmailSvc{sendErr: emailErr},
		&mockJWTSvc{},
		&mockTOTPSvc{},
		&mockBlacklist{},
		cfg,
	)
}

func makeVerifiedUser(email, password string) *domain.User {
	hash, _ := auth.HashPassword(password)
	now := time.Now()
	emailVerifiedAt := now
	return &domain.User{
		ID:                  "user-" + email,
		Email:               email,
		PasswordHash:        hash,
		EmailVerifiedAt:     &emailVerifiedAt,
		FailedLoginAttempts: 0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// ---- Tests ----

func TestRegister_DuplicateEmail_ReturnsConflict(t *testing.T) {
	repo := newMockUserRepo()
	svc := newAuthSvc(repo, nil)
	ctx := context.Background()

	req := service.RegisterRequest{Email: "dup@test.com", Password: "Test1234!", Name: "Dup", Age: 25}
	if err := svc.Register(ctx, req); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	err := svc.Register(ctx, req)
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected ErrConflict on duplicate, got %v", err)
	}
}

func TestLogin_Locked_ReturnsErrAccountLocked(t *testing.T) {
	repo := newMockUserRepo()
	u := makeVerifiedUser("locked@test.com", "Test1234!")
	lockTime := time.Now().Add(15 * time.Minute)
	u.LockedUntil = &lockTime
	repo.users[u.Email] = u
	repo.byID[u.ID] = u

	svc := newAuthSvc(repo, nil)
	_, err := svc.Login(context.Background(), service.LoginRequest{
		Email:    "locked@test.com",
		Password: "Test1234!",
	})
	if !errors.Is(err, domain.ErrAccountLocked) {
		t.Fatalf("expected ErrAccountLocked, got %v", err)
	}
}

func TestLogin_WrongPassword_IncrementsFailedCounter(t *testing.T) {
	repo := newMockUserRepo()
	u := makeVerifiedUser("ok@test.com", "Test1234!")
	repo.users[u.Email] = u
	repo.byID[u.ID] = u

	svc := newAuthSvc(repo, nil)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := svc.Login(ctx, service.LoginRequest{Email: "ok@test.com", Password: "wrong"})
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Fatalf("attempt %d: expected ErrInvalidCredentials, got %v", i+1, err)
		}
	}
	if repo.failedLogins[u.ID] < 3 {
		t.Fatalf("expected >=3 failed login increments, got %d", repo.failedLogins[u.ID])
	}
}

func TestLogin_5Failures_LocksAccount(t *testing.T) {
	repo := newMockUserRepo()
	u := makeVerifiedUser("willlock@test.com", "Test1234!")
	repo.users[u.Email] = u
	repo.byID[u.ID] = u

	svc := newAuthSvc(repo, nil)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		svc.Login(ctx, service.LoginRequest{Email: "willlock@test.com", Password: "wrong"}) //nolint:errcheck
	}

	updated := repo.users["willlock@test.com"]
	if updated.LockedUntil == nil || updated.LockedUntil.Before(time.Now()) {
		t.Fatal("expected account to be locked after 5 failed logins")
	}
}
