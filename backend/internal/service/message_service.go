package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/SX110903/match_app/backend/internal/repository"
	"github.com/google/uuid"
)

type messageService struct {
	msgRepo   repository.IMessageRepository
	matchRepo repository.IMatchRepository
}

func NewMessageService(msgRepo repository.IMessageRepository, matchRepo repository.IMatchRepository) IMessageService {
	return &messageService{
		msgRepo:   msgRepo,
		matchRepo: matchRepo,
	}
}

func (s *messageService) GetMessages(ctx context.Context, userID, matchID string, page, limit int) ([]MessageResponse, error) {
	if err := s.assertMatchParticipant(ctx, userID, matchID); err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	msgs, err := s.msgRepo.GetByMatchID(ctx, matchID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("getting messages: %w", err)
	}

	// Mark messages as read for this user
	_ = s.msgRepo.MarkAllRead(ctx, matchID, userID)

	responses := make([]MessageResponse, len(msgs))
	for i, m := range msgs {
		responses[i] = toMessageResponse(m)
	}
	return responses, nil
}

func (s *messageService) SendMessage(ctx context.Context, userID, matchID, text string) (*MessageResponse, error) {
	if err := s.assertMatchParticipant(ctx, userID, matchID); err != nil {
		return nil, err
	}

	msg := &domain.Message{
		ID:        uuid.New().String(),
		MatchID:   matchID,
		SenderID:  userID,
		Text:      text,
		CreatedAt: time.Now(),
	}

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("creating message: %w", err)
	}

	resp := toMessageResponse(*msg)
	return &resp, nil
}

func (s *messageService) MarkRead(ctx context.Context, userID, matchID string) error {
	if err := s.assertMatchParticipant(ctx, userID, matchID); err != nil {
		return err
	}
	return s.msgRepo.MarkAllRead(ctx, matchID, userID)
}

func (s *messageService) assertMatchParticipant(ctx context.Context, userID, matchID string) error {
	match, err := s.matchRepo.GetMatchByID(ctx, matchID)
	if err != nil {
		return domain.ErrNotFound
	}
	if match.User1ID != userID && match.User2ID != userID {
		return domain.ErrForbidden
	}
	return nil
}

func toMessageResponse(m domain.Message) MessageResponse {
	return MessageResponse{
		ID:        m.ID,
		MatchID:   m.MatchID,
		SenderID:  m.SenderID,
		Text:      m.Text,
		ReadAt:    m.ReadAt,
		CreatedAt: m.CreatedAt,
	}
}
