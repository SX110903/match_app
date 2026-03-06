-- Migration 004: Create matches and messages tables
CREATE TABLE IF NOT EXISTS swipes (
    id          VARCHAR(36)                 NOT NULL PRIMARY KEY,
    swiper_id   VARCHAR(36)                 NOT NULL,
    swiped_id   VARCHAR(36)                 NOT NULL,
    direction   ENUM('left','right','super') NOT NULL,
    created_at  TIMESTAMP                   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_swipe_swiper FOREIGN KEY (swiper_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_swipe_swiped FOREIGN KEY (swiped_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uq_swipe (swiper_id, swiped_id),
    INDEX idx_swipe_swiper_id (swiper_id),
    INDEX idx_swipe_swiped_id (swiped_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS matches (
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    user1_id    VARCHAR(36) NOT NULL,
    user2_id    VARCHAR(36) NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_match_user1 FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_match_user2 FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE KEY uq_match_users (user1_id, user2_id),
    INDEX idx_match_user1_id (user1_id),
    INDEX idx_match_user2_id (user2_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS messages (
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    match_id    VARCHAR(36) NOT NULL,
    sender_id   VARCHAR(36) NOT NULL,
    text        TEXT        NOT NULL,
    read_at     TIMESTAMP   NULL DEFAULT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_message_match  FOREIGN KEY (match_id)  REFERENCES matches(id) ON DELETE CASCADE,
    CONSTRAINT fk_message_sender FOREIGN KEY (sender_id) REFERENCES users(id)   ON DELETE CASCADE,
    INDEX idx_message_match_id  (match_id),
    INDEX idx_message_sender_id (sender_id),
    INDEX idx_message_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
