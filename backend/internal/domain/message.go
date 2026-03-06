package domain

import "time"

type Message struct {
	ID        string    `db:"id"`
	MatchID   string    `db:"match_id"`
	SenderID  string    `db:"sender_id"`
	Text      string    `db:"text"`
	ReadAt    *time.Time `db:"read_at"`
	CreatedAt time.Time  `db:"created_at"`
}

func (m *Message) IsRead() bool {
	return m.ReadAt != nil
}
