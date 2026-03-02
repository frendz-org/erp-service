ALTER TABLE user_sessions DROP CONSTRAINT IF EXISTS fk_user_sessions_parent_session;
DROP INDEX IF EXISTS idx_user_sessions_parent_session_id;
ALTER TABLE user_sessions DROP COLUMN IF EXISTS parent_session_id;
