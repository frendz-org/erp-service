DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_users_tenant_email_active;
DROP INDEX IF EXISTS idx_users_pending_approval;
CREATE UNIQUE INDEX idx_users_email_active ON users(LOWER(email))
    WHERE deleted_at IS NULL;
CREATE INDEX idx_users_pending_verification ON users(created_at)
    WHERE status = 'PENDING_VERIFICATION' AND deleted_at IS NULL;
CREATE INDEX idx_utr_user_id ON user_tenant_registrations(user_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_utr_tenant_status ON user_tenant_registrations(tenant_id, status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_utr_pending_approval ON user_tenant_registrations(tenant_id, created_at)
    WHERE status = 'PENDING_APPROVAL' AND deleted_at IS NULL;
CREATE INDEX idx_prc_application_id ON product_registration_configs(application_id);
