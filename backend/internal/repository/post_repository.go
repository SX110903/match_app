package repository

import (
	"context"
	"fmt"

	"github.com/SX110903/match_app/backend/internal/database"
	"github.com/SX110903/match_app/backend/internal/domain"
)

type postRepository struct{ db *database.DB }

func NewPostRepository(db *database.DB) IPostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *domain.Post) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO posts (id, user_id, content, image_url) VALUES (?, ?, ?, ?)`,
		post.ID, post.UserID, post.Content, post.ImageURL,
	)
	return err
}

// GetFeed devuelve posts ordenados por created_at DESC con paginación OFFSET.
// TODO(perf-v2 BP-6): con inserciones concurrentes, OFFSET provoca duplicados/saltos entre páginas.
// Migrar a cursor-based: añadir parámetro beforeID y usar WHERE p.id < beforeID ORDER BY p.id DESC LIMIT ?.
// El frontend debe guardar el ID del último post recibido y pasarlo como cursor en el siguiente fetch.
func (r *postRepository) GetFeed(ctx context.Context, viewerID string, limit, offset int) ([]domain.Post, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			p.id, p.user_id, p.content, p.image_url, p.likes_count, p.created_at, p.updated_at,
			up.name   AS author_name,
			COALESCE((SELECT ph.url FROM user_photos ph WHERE ph.user_id = p.user_id ORDER BY ph.sort_order LIMIT 1), '') AS author_avatar,
			EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = ?) AS is_liked_by_me
		FROM posts p
		JOIN user_profiles up ON up.user_id = p.user_id
		WHERE p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`,
		viewerID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying feed: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var p domain.Post
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Content, &p.ImageURL, &p.LikesCount, &p.CreatedAt, &p.UpdatedAt,
			&p.AuthorName, &p.AuthorAvatar, &p.IsLikedByMe,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *postRepository) GetByID(ctx context.Context, postID, viewerID string) (*domain.Post, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	var p domain.Post
	err := r.db.QueryRowContext(ctx, `
		SELECT
			p.id, p.user_id, p.content, p.image_url, p.likes_count, p.created_at, p.updated_at,
			up.name AS author_name,
			COALESCE((SELECT ph.url FROM user_photos ph WHERE ph.user_id = p.user_id ORDER BY ph.sort_order LIMIT 1), '') AS author_avatar,
			EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = ?) AS is_liked_by_me
		FROM posts p
		JOIN user_profiles up ON up.user_id = p.user_id
		WHERE p.id = ? AND p.deleted_at IS NULL`,
		viewerID, postID,
	).Scan(
		&p.ID, &p.UserID, &p.Content, &p.ImageURL, &p.LikesCount, &p.CreatedAt, &p.UpdatedAt,
		&p.AuthorName, &p.AuthorAvatar, &p.IsLikedByMe,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *postRepository) Delete(ctx context.Context, postID, userID string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	res, err := r.db.ExecContext(ctx,
		`UPDATE posts SET deleted_at = NOW() WHERE id = ? AND user_id = ? AND deleted_at IS NULL`,
		postID, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *postRepository) LikePost(ctx context.Context, like *domain.PostLike) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint
	res, err := tx.ExecContext(ctx,
		`INSERT IGNORE INTO post_likes (id, post_id, user_id) VALUES (?, ?, ?)`,
		like.ID, like.PostID, like.UserID,
	)
	if err != nil {
		return err
	}
	// Only increment if actually inserted (INSERT IGNORE is no-op on duplicate)
	if n, _ := res.RowsAffected(); n > 0 {
		if _, err := tx.ExecContext(ctx,
			`UPDATE posts SET likes_count = likes_count + 1 WHERE id = ?`,
			like.PostID,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *postRepository) UnlikePost(ctx context.Context, postID, userID string) error {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint
	res, err := tx.ExecContext(ctx,
		`DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`,
		postID, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n > 0 {
		if _, err := tx.ExecContext(ctx,
			`UPDATE posts SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = ?`,
			postID,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *postRepository) GetComments(ctx context.Context, postID string) ([]domain.PostComment, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			c.id, c.post_id, c.user_id, c.content, c.created_at,
			up.name AS author_name,
			COALESCE((SELECT ph.url FROM user_photos ph WHERE ph.user_id = c.user_id ORDER BY ph.sort_order LIMIT 1), '') AS author_avatar
		FROM post_comments c
		JOIN user_profiles up ON up.user_id = c.user_id
		WHERE c.post_id = ? AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.PostComment
	for rows.Next() {
		var c domain.PostComment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.AuthorName, &c.AuthorAvatar); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *postRepository) AddComment(ctx context.Context, comment *domain.PostComment) (*domain.PostComment, error) {
	ctx, cancel := r.db.WithTimeout(ctx)
	defer cancel()
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO post_comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`,
		comment.ID, comment.PostID, comment.UserID, comment.Content,
	); err != nil {
		return nil, err
	}
	// Re-fetch with author JOIN so response includes name/avatar
	var c domain.PostComment
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at,
			up.name AS author_name,
			COALESCE((SELECT ph.url FROM user_photos ph WHERE ph.user_id = c.user_id ORDER BY ph.sort_order LIMIT 1), '') AS author_avatar
		FROM post_comments c
		JOIN user_profiles up ON up.user_id = c.user_id
		WHERE c.id = ?`, comment.ID,
	).Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.AuthorName, &c.AuthorAvatar)
	if err != nil {
		return comment, nil // fallback: return without author info
	}
	return &c, nil
}
