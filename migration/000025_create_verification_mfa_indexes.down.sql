-- Verification Challenges Indexes
DROP INDEX IF EXISTS idx_vc_rate_limit;
DROP INDEX IF EXISTS idx_vc_cleanup;
DROP INDEX IF EXISTS idx_vc_token_hash;
DROP INDEX IF EXISTS idx_vc_expires_at;
DROP INDEX IF EXISTS idx_vc_user_purpose;
DROP INDEX IF EXISTS idx_vc_tenant_identifier;

-- MFA Enrollments Indexes
DROP INDEX IF EXISTS idx_mfa_user_active;
