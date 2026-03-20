-- Migration 011: Performance indexes for production load
-- Note: run each statement separately; ignore "Duplicate key name" errors on re-run.

ALTER TABLE posts ADD INDEX idx_posts_user_created (user_id, created_at DESC);
ALTER TABLE posts ADD INDEX idx_posts_likes_created (likes_count DESC, created_at DESC);
ALTER TABLE users ADD INDEX idx_users_badge_vip (badge, vip_level DESC, follower_count DESC);
ALTER TABLE ads ADD INDEX idx_ads_active_badge (active, target_badge);
ALTER TABLE users ADD INDEX idx_users_deleted_at (deleted_at);
