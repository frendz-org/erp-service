-- Add application_id to user_tenant_registrations for per-product membership
-- A user can have separate registrations (PARTICIPANT, MEMBER) per product

ALTER TABLE user_tenant_registrations
    ADD COLUMN application_id UUID;

-- FK to applications table
ALTER TABLE user_tenant_registrations
    ADD CONSTRAINT fk_utr_application FOREIGN KEY (application_id)
        REFERENCES applications(id) ON DELETE RESTRICT;

-- Drop old unique constraint that only considers (user_id, tenant_id, registration_type)
ALTER TABLE user_tenant_registrations
    DROP CONSTRAINT IF EXISTS uq_utr_user_tenant_type;

-- New unique index includes application_id (COALESCE for NULL safety)
-- Allows one registration per (user, tenant, type, application) combination
CREATE UNIQUE INDEX uq_utr_user_tenant_type_app
    ON user_tenant_registrations (
        user_id,
        tenant_id,
        registration_type,
        COALESCE(application_id, '00000000-0000-0000-0000-000000000000'::uuid)
    )
    WHERE deleted_at IS NULL;

-- Index for member list queries: filter by tenant + application + type
CREATE INDEX idx_utr_tenant_app_type
    ON user_tenant_registrations (tenant_id, application_id, registration_type)
    WHERE deleted_at IS NULL;

-- Standalone FK index for ON DELETE RESTRICT enforcement on applications
CREATE INDEX idx_utr_application_id
    ON user_tenant_registrations (application_id)
    WHERE deleted_at IS NULL;

COMMENT ON COLUMN user_tenant_registrations.application_id IS 'Product this registration belongs to. NULL for legacy registrations.';
