package domain

import "time"

type SwipeDirection string

const (
	SwipeLeft  SwipeDirection = "left"
	SwipeRight SwipeDirection = "right"
	SwipeSuper SwipeDirection = "super"
)

type Swipe struct {
	ID        string         `db:"id"`
	SwiperID  string         `db:"swiper_id"`
	SwipedID  string         `db:"swiped_id"`
	Direction SwipeDirection `db:"direction"`
	CreatedAt time.Time      `db:"created_at"`
}

type Match struct {
	ID        string     `db:"id"         json:"id"`
	User1ID   string     `db:"user1_id"   json:"user1_id"`
	User2ID   string     `db:"user2_id"   json:"user2_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

type MatchWithProfile struct {
	Match
	Profile     UserProfile `json:"profile"`
	LastMessage *string     `json:"last_message"`
	UnreadCount int         `json:"unread_count"`
}

type Candidate struct {
	Profile  UserProfile `json:"profile"`
	Distance float64     `json:"distance"`
}
