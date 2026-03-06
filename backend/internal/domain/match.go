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
	ID        string    `db:"id"`
	User1ID   string    `db:"user1_id"`
	User2ID   string    `db:"user2_id"`
	CreatedAt time.Time `db:"created_at"`
}

type MatchWithProfile struct {
	Match
	Profile     UserProfile
	LastMessage *string
	UnreadCount int
}

type Candidate struct {
	Profile  UserProfile
	Distance float64
}
