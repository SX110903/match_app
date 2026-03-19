package domain

import "time"

type Message struct {
	ID        string     `db:"id"         json:"id"`
	MatchID   string     `db:"match_id"   json:"match_id"`
	SenderID  string     `db:"sender_id"  json:"sender_id"`
	Text      string     `db:"text"       json:"text"`
	ReadAt    *time.Time `db:"read_at"    json:"read_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

func (m *Message) IsRead() bool {
	return m.ReadAt != nil
}

type UserPhoto struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	URL       string    `db:"url"`
	SortOrder int       `db:"sort_order"`
	CreatedAt time.Time `db:"created_at"`
}
