DROP TRIGGER IF EXISTS trg_mfa_enrollments_updated_at ON mfa_enrollments;
DROP INDEX IF EXISTS idx_mfa_enrollments_user_method;
DROP INDEX IF EXISTS idx_mfa_enrollments_primary;
DROP TABLE IF EXISTS mfa_enrollments;
