package service

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
	"github.com/google/uuid"
)

type adminService struct {
	adminRepo repository.IAdminRepository
	userRepo  repository.IUserRepository
}

func NewAdminService(adminRepo repository.IAdminRepository, userRepo repository.IUserRepository) IAdminService {
	return &adminService{adminRepo: adminRepo, userRepo: userRepo}
}

func (s *adminService) AssertAdmin(ctx context.Context, userID string) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrNotFound
	}
	if !u.IsAdmin {
		return domain.ErrForbidden
	}
	return nil
}

func (s *adminService) ListUsers(ctx context.Context, page, limit int) ([]repository.AdminUserSummary, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return s.adminRepo.ListUsers(ctx, limit, offset)
}

func (s *adminService) FreezeUser(ctx context.Context, adminID, targetID string) error {
	if err := s.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	if err := s.adminRepo.SetFrozen(ctx, targetID, true); err != nil {
		return fmt.Errorf("freezing user: %w", err)
	}
	return s.logAction(ctx, adminID, &targetID, "freeze_user", nil)
}

func (s *adminService) UnfreezeUser(ctx context.Context, adminID, targetID string) error {
	if err := s.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	if err := s.adminRepo.SetFrozen(ctx, targetID, false); err != nil {
		return fmt.Errorf("unfreezing user: %w", err)
	}
	return s.logAction(ctx, adminID, &targetID, "unfreeze_user", nil)
}

func (s *adminService) SetVIPLevel(ctx context.Context, adminID, targetID string, level int) error {
	if err := s.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	if err := s.adminRepo.SetVIPLevel(ctx, targetID, level); err != nil {
		return fmt.Errorf("setting VIP: %w", err)
	}
	details := fmt.Sprintf(`{"level":%d}`, level)
	return s.logAction(ctx, adminID, &targetID, "set_vip", &details)
}

func (s *adminService) AdjustCredits(ctx context.Context, adminID, targetID string, delta int) error {
	if err := s.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	if delta > 10000 || delta < -10000 {
		return domain.ErrInvalidInput
	}
	if delta < 0 {
		target, err := s.userRepo.GetByID(ctx, targetID)
		if err != nil {
			return domain.ErrNotFound
		}
		if target.Credits+delta < 0 {
			return domain.ErrInvalidInput
		}
	}
	if err := s.adminRepo.AddCredits(ctx, targetID, delta); err != nil {
		return fmt.Errorf("adjusting credits: %w", err)
	}
	details := fmt.Sprintf(`{"delta":%d}`, delta)
	return s.logAction(ctx, adminID, &targetID, "adjust_credits", &details)
}

func (s *adminService) SetAdmin(ctx context.Context, adminID, targetID string, isAdmin bool) error {
	if err := s.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	if adminID == targetID {
		return domain.ErrSelfAction
	}
	if err := s.adminRepo.SetAdmin(ctx, targetID, isAdmin); err != nil {
		return fmt.Errorf("setting admin: %w", err)
	}
	action := "grant_admin"
	if !isAdmin {
		action = "revoke_admin"
	}
	return s.logAction(ctx, adminID, &targetID, action, nil)
}

func (s *adminService) DeleteUser(ctx context.Context, adminID, targetID string) error {
	if err := s.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	if adminID == targetID {
		return domain.ErrSelfAction
	}
	if err := s.userRepo.SoftDelete(ctx, targetID); err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return s.logAction(ctx, adminID, &targetID, "delete_user", nil)
}

func (s *adminService) GetAuditLog(ctx context.Context, adminID string, page, limit int) ([]domain.AdminLog, error) {
	if err := s.AssertAdmin(ctx, adminID); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return s.adminRepo.GetAuditLog(ctx, limit, offset)
}

func (s *adminService) GetNotificationSettings(ctx context.Context, userID string) (*domain.NotificationSettings, error) {
	settings, err := s.adminRepo.GetNotificationSettings(ctx, userID)
	if err != nil {
		// Return defaults
		return &domain.NotificationSettings{
			UserID:      userID,
			NewMatches:  true,
			NewMessages: true,
			NewsUpdates: false,
			Marketing:   false,
		}, nil
	}
	return settings, nil
}

func (s *adminService) SaveNotificationSettings(ctx context.Context, settings *domain.NotificationSettings) error {
	return s.adminRepo.UpsertNotificationSettings(ctx, settings)
}

func (s *adminService) GetPrivacySettings(ctx context.Context, userID string) (*domain.PrivacySettings, error) {
	settings, err := s.adminRepo.GetPrivacySettings(ctx, userID)
	if err != nil {
		return &domain.PrivacySettings{
			UserID:           userID,
			ShowOnlineStatus: true,
			ShowLastSeen:     true,
			ShowDistance:     true,
			IncognitoMode:    false,
		}, nil
	}
	return settings, nil
}

func (s *adminService) SavePrivacySettings(ctx context.Context, settings *domain.PrivacySettings) error {
	return s.adminRepo.UpsertPrivacySettings(ctx, settings)
}

func (s *adminService) logAction(ctx context.Context, adminID string, targetID *string, action string, details *string) error {
	log := &domain.AdminLog{
		ID:       uuid.New().String(),
		AdminID:  adminID,
		TargetID: targetID,
		Action:   action,
		Details:  details,
	}
	return s.adminRepo.LogAction(ctx, log)
}
