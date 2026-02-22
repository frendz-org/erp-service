DROP INDEX IF EXISTS idx_files_claimed_at;
ALTER TABLE files DROP COLUMN IF EXISTS claimed_at;
