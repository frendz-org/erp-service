ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_tenant;
ALTER TABLE users DROP CONSTRAINT IF EXISTS uq_users_tenant_email;
ALTER TABLE users DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE users ADD CONSTRAINT uq_users_email UNIQUE (email);
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_status;
ALTER TABLE users ADD CONSTRAINT chk_users_status CHECK (status IN (
    'PENDING_VERIFICATION',
    'ACTIVE',
    'INACTIVE',
    'SUSPENDED',
    'LOCKED'
));
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'PENDING_VERIFICATION';
UPDATE users SET status = 'PENDING_VERIFICATION' WHERE status = 'PENDING_APPROVAL';
COMMENT ON TABLE users IS 'Platform-level user identity - tenant association via user_tenant_registrations';
COMMENT ON COLUMN users.email IS 'Login identifier, globally unique, stored lowercase';
COMMENT ON COLUMN users.status IS 'User lifecycle: PENDING_VERIFICATION → ACTIVE. Email verification via verification_challenges.';
