-- Primary lookup: find pending challenges by tenant and identifier
CREATE INDEX idx_vc_tenant_identifier ON verification_challenges(tenant_id, identifier, identifier_type)
    WHERE status = 'PENDING';

-- User's active challenges by purpose
CREATE INDEX idx_vc_user_purpose ON verification_challenges(user_id, purpose)
    WHERE status IN ('PENDING', 'VERIFIED') AND user_id IS NOT NULL;

-- Find expired challenges (for cleanup job)
CREATE INDEX idx_vc_expires_at ON verification_challenges(expires_at)
    WHERE status = 'PENDING';

-- Magic link token lookup
CREATE INDEX idx_vc_token_hash ON verification_challenges(token_hash)
    WHERE token_hash IS NOT NULL AND status = 'PENDING';

-- Cleanup job: find challenges to expire/archive
CREATE INDEX idx_vc_cleanup ON verification_challenges(status, expires_at)
    WHERE status = 'PENDING';

-- Rate limiting: count recent challenges per identifier
CREATE INDEX idx_vc_rate_limit ON verification_challenges(tenant_id, identifier, created_at)
    WHERE identifier_type = 'EMAIL';

-- ============================================================================
-- MFA Enrollments Indexes
-- ============================================================================

-- User's active and verified MFA methods (used during login)
CREATE INDEX idx_mfa_user_active ON mfa_enrollments(user_id)
    WHERE is_active = TRUE AND is_verified = TRUE;

COMMENT ON INDEX idx_vc_tenant_identifier IS 'Primary lookup for pending verification challenges by email/phone';
COMMENT ON INDEX idx_vc_user_purpose IS 'Find user active challenges for specific purpose (e.g., password reset)';
COMMENT ON INDEX idx_vc_expires_at IS 'Cleanup job: efficiently find expired pending challenges';
COMMENT ON INDEX idx_vc_token_hash IS 'Magic link token verification lookup';
COMMENT ON INDEX idx_vc_cleanup IS 'Background job to mark expired challenges';
COMMENT ON INDEX idx_vc_rate_limit IS 'Rate limiting: count recent registration/reset attempts per email';
COMMENT ON INDEX idx_mfa_user_active IS 'MFA login: get user verified and active MFA methods';
