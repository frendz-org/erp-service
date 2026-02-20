DROP TRIGGER IF EXISTS trg_user_sessions_updated_at ON user_sessions;
DROP INDEX IF EXISTS idx_user_sessions_refresh_token_id;
DROP INDEX IF EXISTS idx_user_sessions_user_active;
DROP TABLE IF EXISTS user_sessions;
