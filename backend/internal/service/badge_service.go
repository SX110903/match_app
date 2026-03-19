package service

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
)

type badgeService struct {
	followRepo repository.IFollowRepository
	userRepo   repository.IUserRepository
	adminSvc   IAdminService
}

func NewBadgeService(
	followRepo repository.IFollowRepository,
	userRepo repository.IUserRepository,
	adminSvc IAdminService,
) IBadgeService {
	return &badgeService{
		followRepo: followRepo,
		userRepo:   userRepo,
		adminSvc:   adminSvc,
	}
}

func (s *badgeService) Follow(ctx context.Context, followerID, targetID string) error {
	if followerID == targetID {
		return domain.ErrSelfAction
	}
	if err := s.followRepo.Follow(ctx, followerID, targetID); err != nil {
		return fmt.Errorf("follow: %w", err)
	}
	return s.maybeUpgradeBadge(ctx, targetID)
}

func (s *badgeService) Unfollow(ctx context.Context, followerID, targetID string) error {
	if followerID == targetID {
		return domain.ErrSelfAction
	}
	if err := s.followRepo.Unfollow(ctx, followerID, targetID); err != nil {
		return fmt.Errorf("unfollow: %w", err)
	}
	return s.maybeDowngradeBadge(ctx, targetID)
}

func (s *badgeService) maybeUpgradeBadge(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil // best-effort
	}
	if user.FollowerCount >= domain.InfluencerThreshold &&
		(user.Badge == domain.BadgeNone || user.Badge == domain.BadgeInfluencer) {
		return s.userRepo.UpdateBadge(ctx, userID, domain.BadgeInfluencer)
	}
	return nil
}

func (s *badgeService) maybeDowngradeBadge(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil // best-effort
	}
	if user.FollowerCount < domain.InfluencerThreshold && user.Badge == domain.BadgeInfluencer {
		return s.userRepo.UpdateBadge(ctx, userID, domain.BadgeNone)
	}
	return nil
}

func (s *badgeService) GetFollowers(ctx context.Context, userID string, page, limit int) ([]string, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return s.followRepo.GetFollowers(ctx, userID, limit, offset)
}

func (s *badgeService) GetFollowing(ctx context.Context, userID string, page, limit int) ([]string, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}
	return s.followRepo.GetFollowing(ctx, userID, limit, offset)
}

func (s *badgeService) RequestVerify(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrNotFound
	}
	if user.Badge == domain.BadgeVerified || user.Badge == domain.BadgeVerifiedGov {
		return domain.ErrConflict
	}
	if user.Credits < domain.VerifyCost {
		return domain.ErrInvalidInput
	}
	if err := s.userRepo.UpdateCredits(ctx, userID, -domain.VerifyCost); err != nil {
		return fmt.Errorf("deducting credits: %w", err)
	}
	return s.userRepo.UpdateBadge(ctx, userID, domain.BadgeVerified)
}

func (s *badgeService) AdminSetBadge(ctx context.Context, adminID, targetID, badge string) error {
	if err := s.adminSvc.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	if !domain.ValidBadges[badge] {
		return domain.ErrInvalidInput
	}
	return s.userRepo.UpdateBadge(ctx, targetID, badge)
}
