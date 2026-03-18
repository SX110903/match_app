package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
	"github.com/google/uuid"
)

type messageRepository struct {
	db *database.DB
}

func NewMessageRepository(db *database.DB) IMessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, msg *domain.Message) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO messages (id, match_id, sender_id, text, created_at) VALUES (?, ?, ?, ?, ?)`,
		msg.ID, msg.MatchID, msg.SenderID, msg.Text, msg.CreatedAt,
	)
	return err
}

func (r *messageRepository) GetByMatchID(ctx context.Context, matchID string, limit, offset int) ([]domain.Message, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var messages []domain.Message
	err := r.db.SelectContext(ctx, &messages,
		`SELECT id, match_id, sender_id, text, read_at, created_at
		 FROM messages
		 WHERE match_id = ?
		 ORDER BY created_at ASC
		 LIMIT ? OFFSET ?`,
		matchID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("getting messages: %w", err)
	}
	return messages, nil
}

func (r *messageRepository) MarkAllRead(ctx context.Context, matchID, recipientID string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE messages SET read_at = ?
		 WHERE match_id = ? AND sender_id != ? AND read_at IS NULL`,
		now, matchID, recipientID,
	)
	return err
}

func (r *messageRepository) GetUnreadCount(ctx context.Context, matchID, recipientID string) (int, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM messages
		 WHERE match_id = ? AND sender_id != ? AND read_at IS NULL`,
		matchID, recipientID,
	)
	return count, err
}

func (r *messageRepository) GetLastMessage(ctx context.Context, matchID string) (*domain.Message, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var msg domain.Message
	err := r.db.GetContext(ctx, &msg,
		`SELECT id, match_id, sender_id, text, read_at, created_at
		 FROM messages WHERE match_id = ?
		 ORDER BY created_at DESC LIMIT 1`,
		matchID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting last message: %w", err)
	}
	return &msg, nil
}
