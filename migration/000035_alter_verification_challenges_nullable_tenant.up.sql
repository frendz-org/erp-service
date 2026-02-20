ALTER TABLE verification_challenges ALTER COLUMN tenant_id DROP NOT NULL;

COMMENT ON COLUMN verification_challenges.tenant_id IS 'Tenant context for the challenge. NULL for pre-registration flows where user has no tenant association yet.';
