ALTER TABLE files ADD COLUMN claimed_at TIMESTAMPTZ;
CREATE INDEX idx_files_claimed_at ON files (claimed_at) WHERE claimed_at IS NOT NULL AND deleted_at IS NULL;
