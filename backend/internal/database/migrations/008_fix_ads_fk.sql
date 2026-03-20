-- Migration 008: Fix ads FK to use ON DELETE SET NULL
-- Evita que borrar un usuario admin deje huérfanas las filas de ads

ALTER TABLE ads MODIFY COLUMN created_by CHAR(36) NULL;
ALTER TABLE ads DROP FOREIGN KEY IF EXISTS ads_ibfk_1;
ALTER TABLE ads ADD CONSTRAINT ads_created_by_fk
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
