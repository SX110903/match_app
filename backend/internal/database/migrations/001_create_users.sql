-- Migration 001: Create users table
CREATE TABLE IF NOT EXISTS users (
    id                      VARCHAR(36)     NOT NULL PRIMARY KEY,
    email                   VARCHAR(255)    NOT NULL UNIQUE,
    password_hash           VARCHAR(255)    NOT NULL,
    email_verified_at       TIMESTAMP       NULL DEFAULT NULL,
    totp_secret             TEXT            NULL COMMENT 'AES-256-GCM encrypted TOTP secret',
    totp_enabled            BOOLEAN         NOT NULL DEFAULT FALSE,
    backup_codes            TEXT            NULL COMMENT 'AES-256-GCM encrypted JSON array of backup codes',
    last_login_at           TIMESTAMP       NULL DEFAULT NULL,
    failed_login_attempts   INT             NOT NULL DEFAULT 0,
    locked_until            TIMESTAMP       NULL DEFAULT NULL,
    deleted_at              TIMESTAMP       NULL DEFAULT NULL COMMENT 'Soft delete',
    created_at              TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_users_email       (email),
    INDEX idx_users_deleted_at  (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
