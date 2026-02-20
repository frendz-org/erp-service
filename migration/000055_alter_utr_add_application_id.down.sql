-- Reverse: remove application_id from user_tenant_registrations

DROP INDEX IF EXISTS idx_utr_application_id;
DROP INDEX IF EXISTS idx_utr_tenant_app_type;
DROP INDEX IF EXISTS uq_utr_user_tenant_type_app;

ALTER TABLE user_tenant_registrations
    DROP CONSTRAINT IF EXISTS fk_utr_application;

ALTER TABLE user_tenant_registrations
    DROP COLUMN IF EXISTS application_id;

-- Restore original unique constraint
ALTER TABLE user_tenant_registrations
    ADD CONSTRAINT uq_utr_user_tenant_type UNIQUE (user_id, tenant_id, registration_type);
