package service

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type exploreService struct {
	db *database.DB
}

func NewExploreService(db *database.DB) IExploreService {
	return &exploreService{db: db}
}

func (s *exploreService) GetUsers(ctx context.Context, callerID, cursor string, limit int) ([]ExploreUserResponse, error) {
	ctx, cancel := s.db.WithTimeout(ctx)
	defer cancel()

	var args []interface{}
	cursorClause := ""
	if cursor != "" {
		cursorClause = "AND u.id > ?"
		args = append(args, cursor)
	}

	args = append(args, callerID, callerID, callerID, limit)

	query := fmt.Sprintf(`
		SELECT u.id, COALESCE(up.name,'') as name, COALESCE(up.age,0) as age,
		       COALESCE((SELECT ph.url FROM user_photos ph WHERE ph.user_id = u.id ORDER BY ph.sort_order LIMIT 1),'') as avatar,
		       u.badge, u.vip_level, u.follower_count,
		       EXISTS(SELECT 1 FROM follows f WHERE f.follower_id = ? AND f.followed_id = u.id) as is_following
		FROM users u
		LEFT JOIN user_profiles up ON up.user_id = u.id
		WHERE u.deleted_at IS NULL
		  AND u.is_frozen = 0
		  AND u.id != ?
		  AND u.id NOT IN (SELECT followed_id FROM follows WHERE follower_id = ?)
		  %s
		ORDER BY
		  CASE u.badge WHEN 'verified_gov' THEN 0 WHEN 'verified' THEN 1 WHEN 'influencer' THEN 2 ELSE 3 END ASC,
		  u.vip_level DESC,
		  u.follower_count DESC
		LIMIT ?`, cursorClause)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("explore users: %w", err)
	}
	defer rows.Close()

	var out []ExploreUserResponse
	for rows.Next() {
		var u ExploreUserResponse
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Avatar, &u.Badge, &u.VIPLevel, &u.FollowerCount, &u.IsFollowing); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (s *exploreService) GetPosts(ctx context.Context, callerID, cursor string, limit int) ([]PostResponse, error) {
	ctx, cancel := s.db.WithTimeout(ctx)
	defer cancel()

	cursorClause := ""
	var args []interface{}
	args = append(args, callerID)
	if cursor != "" {
		cursorClause = "AND p.id < ?"
		args = append(args, cursor)
	}
	args = append(args, limit)

	query := fmt.Sprintf(`
		SELECT p.id, p.user_id, p.content, p.image_url, p.likes_count, p.created_at, p.updated_at,
		       up.name AS author_name,
		       COALESCE((SELECT ph.url FROM user_photos ph WHERE ph.user_id = p.user_id ORDER BY ph.sort_order LIMIT 1),'') AS author_avatar,
		       EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = ?) AS is_liked_by_me
		FROM posts p
		JOIN user_profiles up ON up.user_id = p.user_id
		WHERE p.deleted_at IS NULL
		  %s
		  AND p.created_at >= NOW() - INTERVAL 24 HOUR
		ORDER BY p.likes_count DESC, p.created_at DESC
		LIMIT ?`, cursorClause)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("explore posts: %w", err)
	}
	defer rows.Close()

	var out []PostResponse
	for rows.Next() {
		var p domain.Post
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Content, &p.ImageURL, &p.LikesCount, &p.CreatedAt, &p.UpdatedAt,
			&p.AuthorName, &p.AuthorAvatar, &p.IsLikedByMe,
		); err != nil {
			return nil, err
		}
		out = append(out, toPostResponse(p))
	}
	return out, rows.Err()
}
