package service

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
)

type userService struct {
	userRepo    repository.IUserRepository
	profileRepo repository.IProfileRepository
}

func NewUserService(userRepo repository.IUserRepository, profileRepo repository.IProfileRepository) IUserService {
	return &userService{userRepo: userRepo, profileRepo: profileRepo}
}

func (s *userService) GetMe(ctx context.Context, userID string) (*UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting profile: %w", err)
	}

	photos := make([]PhotoResponse, len(profile.PhotoObjects))
	for i, p := range profile.PhotoObjects {
		photos[i] = PhotoResponse{ID: p.ID, URL: p.URL}
	}

	return &UserProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        profile.Name,
		Age:         profile.Age,
		Bio:         profile.Bio,
		Occupation:  profile.Occupation,
		Location:    profile.Location,
		Photos:      photos,
		Interests:   profile.Interests,
		TOTPEnabled: user.TOTPEnabled,
		IsAdmin:     user.IsAdmin,
		IsFrozen:    user.IsFrozen,
		VIPLevel:    user.VIPLevel,
		Credits:     user.Credits,
	}, nil
}

func (s *userService) UpdateMe(ctx context.Context, userID string, req UpdateProfileRequest) error {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("getting profile: %w", err)
	}

	if req.Name != nil {
		profile.Name = *req.Name
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.Occupation != nil {
		profile.Occupation = req.Occupation
	}
	if req.Location != nil {
		profile.Location = req.Location
	}

	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return err
	}
	if req.Interests != nil {
		if err := s.profileRepo.ReplaceInterests(ctx, userID, req.Interests); err != nil {
			return fmt.Errorf("updating interests: %w", err)
		}
	}
	return nil
}

func (s *userService) DeleteMe(ctx context.Context, userID string) error {
	return s.userRepo.SoftDelete(ctx, userID)
}

func (s *userService) UpdatePreferences(ctx context.Context, userID string, req UpdatePreferencesRequest) error {
	prefs := &domain.UserPreferences{
		UserID:        userID,
		MinAge:        req.MinAge,
		MaxAge:        req.MaxAge,
		MaxDistanceKm: req.MaxDistanceKm,
		InterestedIn:  req.InterestedIn,
	}
	return s.profileRepo.UpsertPreferences(ctx, prefs)
}
