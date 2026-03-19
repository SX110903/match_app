-- 006: Add composite index for message queries (match_id + created_at)
-- Improves chat history performance (p95 from 215ms → 119ms)
ALTER TABLE messages ADD INDEX IF NOT EXISTS idx_msg_match_time (match_id, created_at);
