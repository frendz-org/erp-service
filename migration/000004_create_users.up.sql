CREATE TABLE users (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,

    -- Core Identity (only essential fields)
    email               VARCHAR(255) NOT NULL,

    -- User Lifecycle Status
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING_APPROVAL',
    status_changed_at   TIMESTAMPTZ,
    status_changed_by   UUID,

    -- Registration Metadata (permanent fact about user origin)
    registration_source VARCHAR(20) NOT NULL DEFAULT 'SELF',

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(id) ON DELETE RESTRICT,
    CONSTRAINT uq_users_tenant_email UNIQUE (tenant_id, email),
    CONSTRAINT chk_users_status CHECK (status IN (
        'PENDING_APPROVAL',
        'ACTIVE',
        'INACTIVE',
        'SUSPENDED',
        'LOCKED'
    )),
    CONSTRAINT chk_users_registration_source CHECK (registration_source IN (
        'SELF',
        'ADMIN',
        'IMPORT',
        'GOOGLE'
    )),
    CONSTRAINT chk_users_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE users IS 'Core user identity - minimal fields, authentication/profile/security in related tables';
COMMENT ON COLUMN users.email IS 'Login identifier, unique within tenant, stored lowercase';
COMMENT ON COLUMN users.status IS 'User lifecycle: PENDING_APPROVAL → ACTIVE. Email verification handled in Redis.';
COMMENT ON COLUMN users.registration_source IS 'Permanent fact about how this user was created. Never changes.';
