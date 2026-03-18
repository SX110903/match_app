package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
	"github.com/google/uuid"
)

const maxPhotosPerUser = 6

type photoService struct {
	profileRepo repository.IProfileRepository
}

func NewPhotoService(profileRepo repository.IProfileRepository) IPhotoService {
	return &photoService{profileRepo: profileRepo}
}

func (s *photoService) AddPhoto(ctx context.Context, userID, url string) (*PhotoResponse, error) {
	if !isValidImgurURL(url) {
		return nil, domain.ErrInvalidInput
	}

	count, err := s.profileRepo.GetPhotoCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("checking photo count: %w", err)
	}
	if count >= maxPhotosPerUser {
		return nil, domain.ErrConflict
	}

	photo := &domain.UserPhoto{
		ID:        uuid.New().String(),
		UserID:    userID,
		URL:       url,
		SortOrder: count,
		CreatedAt: time.Now(),
	}

	if err := s.profileRepo.AddPhoto(ctx, photo); err != nil {
		return nil, fmt.Errorf("adding photo: %w", err)
	}

	return &PhotoResponse{ID: photo.ID, URL: photo.URL}, nil
}

func (s *photoService) DeletePhoto(ctx context.Context, userID, photoID string) error {
	if err := s.profileRepo.DeletePhoto(ctx, userID, photoID); err != nil {
		return err
	}
	return nil
}

func isValidImgurURL(url string) bool {
	return strings.HasPrefix(url, "https://i.imgur.com/") ||
		strings.HasPrefix(url, "https://imgur.com/")
}
