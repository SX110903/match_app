-- Migration 003: Create token tables
CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)     NOT NULL,
    token_hash  VARCHAR(64)     NOT NULL UNIQUE COMMENT 'SHA-256 of the actual token',
    expires_at  TIMESTAMP       NOT NULL,
    used_at     TIMESTAMP       NULL DEFAULT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_evtoken_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_evtoken_hash    (token_hash),
    INDEX idx_evtoken_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)     NOT NULL,
    token_hash  VARCHAR(64)     NOT NULL UNIQUE COMMENT 'SHA-256 of the actual token',
    expires_at  TIMESTAMP       NOT NULL,
    used_at     TIMESTAMP       NULL DEFAULT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_prtoken_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_prtoken_hash    (token_hash),
    INDEX idx_prtoken_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          VARCHAR(36)     NOT NULL PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)     NOT NULL,
    token_hash  VARCHAR(64)     NOT NULL UNIQUE COMMENT 'SHA-256 of the refresh token',
    expires_at  TIMESTAMP       NOT NULL,
    revoked_at  TIMESTAMP       NULL DEFAULT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_rttoken_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_rttoken_hash    (token_hash),
    INDEX idx_rttoken_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
