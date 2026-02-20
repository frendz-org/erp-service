DROP INDEX IF EXISTS idx_prc_application_id;
DROP INDEX IF EXISTS idx_utr_pending_approval;
DROP INDEX IF EXISTS idx_utr_tenant_status;
DROP INDEX IF EXISTS idx_utr_user_id;
DROP INDEX IF EXISTS idx_users_pending_verification;
DROP INDEX IF EXISTS idx_users_email_active;

CREATE INDEX idx_users_tenant_id ON users(tenant_id) WHERE deleted_at IS NULL;

CREATE INDEX idx_users_tenant_email_active ON users(tenant_id, LOWER(email))
    WHERE deleted_at IS NULL AND status IN ('ACTIVE', 'PENDING_APPROVAL');

CREATE INDEX idx_users_pending_approval ON users(tenant_id, created_at)
    WHERE status = 'PENDING_APPROVAL' AND deleted_at IS NULL;
