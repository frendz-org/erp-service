DROP INDEX IF EXISTS idx_ura_user_role_branch_unique;
DROP TRIGGER IF EXISTS trg_ura_updated_at ON user_role_assignments;
DROP TABLE IF EXISTS user_role_assignments;
