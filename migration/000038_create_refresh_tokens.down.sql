DROP INDEX IF EXISTS idx_refresh_tokens_user_active;
DROP INDEX IF EXISTS idx_refresh_tokens_token_family;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS uq_refresh_tokens_token_hash;
DROP TABLE IF EXISTS refresh_tokens;
