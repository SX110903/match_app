-- Migration 010: Backfill missing user_profile rows
-- Garantiza que todo usuario tenga al menos una fila en user_profiles.
-- Usa INSERT IGNORE para no sobrescribir perfiles existentes.
-- Necesario para que los JOINs en matches/candidates no descarten usuarios.

INSERT IGNORE INTO user_profiles (id, user_id, name, age)
SELECT UUID(), u.id, COALESCE(u.name, u.email, 'Usuario'), 18
FROM users u
WHERE u.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM user_profiles up WHERE up.user_id = u.id
  );
