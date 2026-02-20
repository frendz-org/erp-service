ALTER TABLE users ADD COLUMN tenant_id UUID;
UPDATE users SET tenant_id = (SELECT id FROM tenants WHERE code = 'platform' LIMIT 1);
ALTER TABLE users ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE users DROP CONSTRAINT IF EXISTS uq_users_email;
ALTER TABLE users ADD CONSTRAINT fk_users_tenant
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;
ALTER TABLE users ADD CONSTRAINT uq_users_tenant_email UNIQUE (tenant_id, email);
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_status;
ALTER TABLE users ADD CONSTRAINT chk_users_status CHECK (status IN (
    'PENDING_APPROVAL',
    'ACTIVE',
    'INACTIVE',
    'SUSPENDED',
    'LOCKED'
));
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'PENDING_APPROVAL';
UPDATE users SET status = 'PENDING_APPROVAL' WHERE status = 'PENDING_VERIFICATION';
COMMENT ON TABLE users IS 'Core user identity - minimal fields, authentication/profile/security in related tables';
COMMENT ON COLUMN users.email IS 'Login identifier, unique within tenant, stored lowercase';
COMMENT ON COLUMN users.status IS 'User lifecycle: PENDING_APPROVAL → ACTIVE. Email verification handled in Redis.';
