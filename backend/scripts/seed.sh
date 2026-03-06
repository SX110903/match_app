#!/bin/bash
# Seed the database with test data
set -euo pipefail

echo "Seeding database with test data..."

mysql -h "${DB_HOST:-localhost}" -P "${DB_PORT:-3306}" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" << 'SQL'
-- Test users (passwords are 'Test1234!' hashed with Argon2id — generate real hashes in production)
INSERT IGNORE INTO users (id, email, password_hash, email_verified_at, created_at, updated_at) VALUES
  ('00000000-0000-0000-0000-000000000001', 'test@matchhub.com', '$argon2id$placeholder', NOW(), NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000002', 'demo@matchhub.com', '$argon2id$placeholder', NOW(), NOW(), NOW());

INSERT IGNORE INTO user_profiles (id, user_id, name, age, bio, occupation, location, created_at, updated_at) VALUES
  ('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'Test User', 25, 'Test bio', 'Developer', 'Madrid', NOW(), NOW()),
  ('10000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000002', 'Demo User', 28, 'Demo bio', 'Designer', 'Barcelona', NOW(), NOW());

SELECT 'Seed completed successfully' as status;
SQL

echo "Done."
