package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
)

var adSanitizer = bluemonday.StrictPolicy()

type adService struct {
	adRepo   repository.IAdRepository
	adminSvc IAdminService
}

func NewAdService(adRepo repository.IAdRepository, adminSvc IAdminService) IAdService {
	return &adService{adRepo: adRepo, adminSvc: adminSvc}
}

func (s *adService) GetActive(ctx context.Context, userID, userBadge string) (*domain.Ad, error) {
	ad, err := s.adRepo.GetActive(ctx, userBadge)
	if err != nil {
		return nil, fmt.Errorf("get active ad: %w", err)
	}
	if ad != nil {
		// best-effort increment
		_ = s.adRepo.IncrementImpressions(ctx, ad.ID)
	}
	return ad, nil
}

func (s *adService) RegisterClick(ctx context.Context, adID, userID string) error {
	return s.adRepo.RegisterClick(ctx, adID, userID)
}

func (s *adService) AdminCreate(ctx context.Context, adminID string, req AdCreateRequest) (*domain.Ad, error) {
	if err := s.adminSvc.AssertAdmin(ctx, adminID); err != nil {
		return nil, err
	}

	title := adSanitizer.Sanitize(req.Title)
	targetBadge := req.TargetBadge
	if targetBadge == "" {
		targetBadge = "all"
	}
	ctaText := req.CTAText
	if ctaText == "" {
		ctaText = "Ver más"
	}

	var desc *string
	if req.Description != nil {
		clean := adSanitizer.Sanitize(*req.Description)
		desc = &clean
	}

	ad := &domain.Ad{
		ID:          uuid.New().String(),
		Title:       title,
		Description: desc,
		ImageURL:    req.ImageURL,
		CTAText:     ctaText,
		CTAURL:      req.CTAURL,
		TargetBadge: targetBadge,
		Active:      req.Active,
		CreatedBy:   adminID,
	}

	if err := s.adRepo.Create(ctx, ad); err != nil {
		return nil, fmt.Errorf("create ad: %w", err)
	}
	return ad, nil
}

func (s *adService) AdminUpdate(ctx context.Context, adminID, adID string, req AdUpdateRequest) (*domain.Ad, error) {
	if err := s.adminSvc.AssertAdmin(ctx, adminID); err != nil {
		return nil, err
	}
	ad, err := s.adRepo.GetByID(ctx, adID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		ad.Title = adSanitizer.Sanitize(*req.Title)
	}
	if req.Description != nil {
		clean := adSanitizer.Sanitize(*req.Description)
		ad.Description = &clean
	}
	if req.ImageURL != nil {
		ad.ImageURL = req.ImageURL
	}
	if req.CTAText != nil {
		ad.CTAText = *req.CTAText
	}
	if req.CTAURL != nil {
		ad.CTAURL = *req.CTAURL
	}
	if req.TargetBadge != nil {
		ad.TargetBadge = *req.TargetBadge
	}
	if req.Active != nil {
		ad.Active = *req.Active
	}

	if err := s.adRepo.Update(ctx, ad); err != nil {
		return nil, fmt.Errorf("update ad: %w", err)
	}
	return ad, nil
}

func (s *adService) AdminDelete(ctx context.Context, adminID, adID string) error {
	if err := s.adminSvc.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	return s.adRepo.Delete(ctx, adID)
}

func (s *adService) AdminToggle(ctx context.Context, adminID, adID string) error {
	if err := s.adminSvc.AssertAdmin(ctx, adminID); err != nil {
		return err
	}
	return s.adRepo.Toggle(ctx, adID)
}

func (s *adService) AdminList(ctx context.Context, adminID string) ([]domain.Ad, error) {
	if err := s.adminSvc.AssertAdmin(ctx, adminID); err != nil {
		return nil, err
	}
	return s.adRepo.ListAll(ctx)
}
