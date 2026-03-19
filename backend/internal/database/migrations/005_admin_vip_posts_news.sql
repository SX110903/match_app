-- Migration 005: Admin, VIP, Credits, Posts, News, Settings

-- Add admin/VIP/credits/frozen fields to users
ALTER TABLE users
  ADD COLUMN is_admin   BOOLEAN     NOT NULL DEFAULT FALSE AFTER deleted_at,
  ADD COLUMN is_frozen  BOOLEAN     NOT NULL DEFAULT FALSE AFTER is_admin,
  ADD COLUMN vip_level  TINYINT     NOT NULL DEFAULT 0     AFTER is_frozen,
  ADD COLUMN credits    INT         NOT NULL DEFAULT 0     AFTER vip_level;

-- Posts table (social feed)
CREATE TABLE IF NOT EXISTS posts (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY,
    user_id     VARCHAR(36)     NOT NULL,
    content     TEXT            NOT NULL,
    image_url   VARCHAR(500)    NULL,
    likes_count INT             NOT NULL DEFAULT 0,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP       NULL DEFAULT NULL,

    CONSTRAINT fk_posts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_posts_user_id    (user_id),
    INDEX idx_posts_created_at (created_at),
    INDEX idx_posts_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Post likes table
CREATE TABLE IF NOT EXISTS post_likes (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY,
    post_id     VARCHAR(36)     NOT NULL,
    user_id     VARCHAR(36)     NOT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_postlikes_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT fk_postlikes_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uq_post_user_like (post_id, user_id),
    INDEX idx_postlikes_post_id  (post_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Post comments table
CREATE TABLE IF NOT EXISTS post_comments (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY,
    post_id     VARCHAR(36)     NOT NULL,
    user_id     VARCHAR(36)     NOT NULL,
    content     TEXT            NOT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP       NULL DEFAULT NULL,

    CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_comments_post_id (post_id),
    INDEX idx_comments_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- News articles table (admin-managed)
CREATE TABLE IF NOT EXISTS news_articles (
    id           VARCHAR(36)     NOT NULL PRIMARY KEY,
    author_id    VARCHAR(36)     NOT NULL,
    title        VARCHAR(255)    NOT NULL,
    summary      TEXT            NOT NULL,
    content      LONGTEXT        NOT NULL,
    image_url    VARCHAR(500)    NULL,
    category     ENUM('tendencias','tech','seguridad','negocios','general') NOT NULL DEFAULT 'general',
    published    BOOLEAN         NOT NULL DEFAULT FALSE,
    published_at TIMESTAMP       NULL DEFAULT NULL,
    created_at   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP       NULL DEFAULT NULL,

    CONSTRAINT fk_news_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_news_category    (category),
    INDEX idx_news_published   (published),
    INDEX idx_news_deleted_at  (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Admin action logs
CREATE TABLE IF NOT EXISTS admin_logs (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY,
    admin_id    VARCHAR(36)     NOT NULL,
    target_id   VARCHAR(36)     NULL,
    action      VARCHAR(100)    NOT NULL,
    details     JSON            NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_adminlogs_admin   (admin_id),
    INDEX idx_adminlogs_target  (target_id),
    INDEX idx_adminlogs_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- User notification settings
CREATE TABLE IF NOT EXISTS user_notification_settings (
    user_id          VARCHAR(36) NOT NULL PRIMARY KEY,
    new_matches      BOOLEAN     NOT NULL DEFAULT TRUE,
    new_messages     BOOLEAN     NOT NULL DEFAULT TRUE,
    news_updates     BOOLEAN     NOT NULL DEFAULT FALSE,
    marketing        BOOLEAN     NOT NULL DEFAULT FALSE,
    updated_at       TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_notif_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- User privacy settings
CREATE TABLE IF NOT EXISTS user_privacy_settings (
    user_id             VARCHAR(36) NOT NULL PRIMARY KEY,
    show_online_status  BOOLEAN     NOT NULL DEFAULT TRUE,
    show_last_seen      BOOLEAN     NOT NULL DEFAULT TRUE,
    show_distance       BOOLEAN     NOT NULL DEFAULT TRUE,
    incognito_mode      BOOLEAN     NOT NULL DEFAULT FALSE,
    updated_at          TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_privacy_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
