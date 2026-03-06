package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
)

type matchService struct {
	matchRepo   repository.IMatchRepository
	profileRepo repository.IProfileRepository
}

func NewMatchService(matchRepo repository.IMatchRepository, profileRepo repository.IProfileRepository) IMatchService {
	return &matchService{matchRepo: matchRepo, profileRepo: profileRepo}
}

func (s *matchService) GetCandidates(ctx context.Context, userID string, page, limit int) ([]domain.Candidate, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit

	prefs, err := s.profileRepo.GetPreferences(ctx, userID)
	if err != nil {
		// Use defaults if no preferences set
		prefs = &domain.UserPreferences{
			MinAge:        18,
			MaxAge:        100,
			MaxDistanceKm: 50,
			InterestedIn:  "both",
		}
	}

	return s.matchRepo.GetCandidates(ctx, userID, prefs, limit, offset)
}

func (s *matchService) Swipe(ctx context.Context, userID, targetID string, direction domain.SwipeDirection) (*SwipeResponse, error) {
	if userID == targetID {
		return nil, domain.ErrSelfAction
	}

	swipe := &domain.Swipe{
		ID:        uuid.New().String(),
		SwiperID:  userID,
		SwipedID:  targetID,
		Direction: direction,
		CreatedAt: time.Now(),
	}

	if err := s.matchRepo.CreateSwipe(ctx, swipe); err != nil {
		return nil, fmt.Errorf("creating swipe: %w", err)
	}

	// Check for mutual like (match)
	if direction == domain.SwipeRight || direction == domain.SwipeSuper {
		reverseSwipe, err := s.matchRepo.GetSwipe(ctx, targetID, userID)
		if err == nil && (reverseSwipe.Direction == domain.SwipeRight || reverseSwipe.Direction == domain.SwipeSuper) {
			match := &domain.Match{
				ID:        uuid.New().String(),
				User1ID:   userID,
				User2ID:   targetID,
				CreatedAt: time.Now(),
			}
			if err := s.matchRepo.CreateMatch(ctx, match); err != nil {
				return nil, fmt.Errorf("creating match: %w", err)
			}
			return &SwipeResponse{IsMatch: true, MatchID: match.ID}, nil
		}
	}

	return &SwipeResponse{IsMatch: false}, nil
}

func (s *matchService) GetMatches(ctx context.Context, userID string) ([]domain.MatchWithProfile, error) {
	return s.matchRepo.GetMatchesByUserID(ctx, userID)
}

func (s *matchService) GetMatch(ctx context.Context, userID, matchID string) (*domain.MatchWithProfile, error) {
	match, err := s.matchRepo.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Ensure the requesting user is part of this match
	if match.User1ID != userID && match.User2ID != userID {
		return nil, domain.ErrForbidden
	}

	// Get the other user's profile
	otherUserID := match.User2ID
	if match.User2ID == userID {
		otherUserID = match.User1ID
	}

	profile, err := s.profileRepo.GetByUserID(ctx, otherUserID)
	if err != nil {
		return nil, fmt.Errorf("getting profile: %w", err)
	}

	return &domain.MatchWithProfile{Match: *match, Profile: *profile}, nil
}
