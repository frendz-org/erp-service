UPDATE verification_challenges
SET tenant_id = (SELECT id FROM tenants WHERE code = 'platform' LIMIT 1)
WHERE tenant_id IS NULL;

ALTER TABLE verification_challenges ALTER COLUMN tenant_id SET NOT NULL;

COMMENT ON COLUMN verification_challenges.tenant_id IS 'Tenant scoping for the verification challenge';
