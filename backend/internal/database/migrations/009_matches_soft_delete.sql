-- Migration 009: Add soft-delete support to matches
-- Permite que los usuarios "eliminen" su match sin afectar al otro participante.
-- El backend filtra matches donde deleted_at IS NULL.

ALTER TABLE matches ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_matches_deleted_at ON matches(deleted_at);
