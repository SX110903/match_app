package domain

import "time"

type Post struct {
	ID         string     `db:"id"          json:"id"`
	UserID     string     `db:"user_id"     json:"user_id"`
	Content    string     `db:"content"     json:"content"`
	ImageURL   *string    `db:"image_url"   json:"image_url,omitempty"`
	LikesCount int        `db:"likes_count" json:"likes_count"`
	CreatedAt  time.Time  `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"  json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"  json:"-"`

	// Populated from JOINs
	AuthorName   string `db:"-" json:"author_name"`
	AuthorAvatar string `db:"-" json:"author_avatar"`
	IsLikedByMe  bool   `db:"-" json:"is_liked_by_me"`
}

type PostComment struct {
	ID        string     `db:"id"         json:"id"`
	PostID    string     `db:"post_id"    json:"post_id"`
	UserID    string     `db:"user_id"    json:"user_id"`
	Content   string     `db:"content"    json:"content"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`

	// Populated from JOINs
	AuthorName   string `db:"-" json:"author_name"`
	AuthorAvatar string `db:"-" json:"author_avatar"`
}

type PostLike struct {
	ID        string    `db:"id"         json:"id"`
	PostID    string    `db:"post_id"    json:"post_id"`
	UserID    string    `db:"user_id"    json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
