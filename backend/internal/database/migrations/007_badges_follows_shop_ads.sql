-- Migration 007: Badges, Follows, Shop transactions, Ads

-- Add badge and follower_count to users
ALTER TABLE users
  ADD COLUMN badge          ENUM('none','influencer','verified','verified_gov') NOT NULL DEFAULT 'none' AFTER credits,
  ADD COLUMN follower_count INT NOT NULL DEFAULT 0 AFTER badge;

-- Follows table
CREATE TABLE IF NOT EXISTS follows (
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    follower_id VARCHAR(36) NOT NULL,
    following_id VARCHAR(36) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uq_follow (follower_id, following_id),
    CONSTRAINT fk_follows_follower  FOREIGN KEY (follower_id)  REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_follows_following FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_follows_following (following_id),
    INDEX idx_follows_follower  (follower_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Shop transactions table
CREATE TABLE IF NOT EXISTS shop_transactions (
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    user_id     VARCHAR(36) NOT NULL,
    item_type   VARCHAR(50) NOT NULL,
    item_value  INT         NOT NULL,
    cost        INT         NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_shop_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_shop_user_id    (user_id),
    INDEX idx_shop_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Ads table
CREATE TABLE IF NOT EXISTS ads (
    id           VARCHAR(36)  NOT NULL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    description  TEXT         NULL,
    image_url    VARCHAR(500) NULL,
    cta_text     VARCHAR(100) NOT NULL DEFAULT 'Ver más',
    cta_url      VARCHAR(500) NOT NULL,
    target_badge ENUM('none','influencer','verified','verified_gov','all') NOT NULL DEFAULT 'all',
    active       BOOLEAN      NOT NULL DEFAULT TRUE,
    impressions  INT          NOT NULL DEFAULT 0,
    clicks       INT          NOT NULL DEFAULT 0,
    created_by   VARCHAR(36)  NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_ads_creator FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_ads_active       (active),
    INDEX idx_ads_target_badge (target_badge)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
