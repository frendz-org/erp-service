CREATE TABLE tenants (
    -- Primary Key (UUIDv7 for time-ordered, sortable IDs)
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Business Identifiers
    code                VARCHAR(50) NOT NULL,
    name                VARCHAR(255) NOT NULL,

    -- Configuration (see TRD Section 3.1.6 for hierarchy)
    settings            JSONB NOT NULL DEFAULT '{}',

    -- Status Management
    status              VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking (application-managed)
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT uq_tenants_code UNIQUE (code),
    CONSTRAINT chk_tenants_status CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED'))
);

CREATE TRIGGER trg_tenants_updated_at
    BEFORE UPDATE ON tenants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE tenants IS 'Root table for multi-tenant organization management';
COMMENT ON COLUMN tenants.code IS 'Unique identifier code for the tenant (URL-safe, immutable)';
COMMENT ON COLUMN tenants.settings IS 'Per-tenant configuration: password_policy, pin_policy, session, branding, approval_required';
COMMENT ON COLUMN tenants.status IS 'Tenant lifecycle status: ACTIVE, INACTIVE, SUSPENDED';
COMMENT ON COLUMN tenants.deleted_at IS 'Soft delete timestamp - NULL means active';
COMMENT ON COLUMN tenants.version IS 'Optimistic lock version - managed by application';
