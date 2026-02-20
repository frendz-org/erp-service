CREATE TABLE applications (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,

    -- Business Identifiers
    code                VARCHAR(50) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    description         TEXT,

    -- Application Settings (overrides tenant settings — see TRD Section 3.1.6)
    settings            JSONB NOT NULL DEFAULT '{}',

    -- Status Management
    status              VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_applications_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(id) ON DELETE RESTRICT,
    CONSTRAINT uq_applications_tenant_code UNIQUE (tenant_id, code),
    CONSTRAINT chk_applications_status CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE TRIGGER trg_applications_updated_at
    BEFORE UPDATE ON applications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE applications IS 'Applications that use IAM - each defines its own roles and permissions';
COMMENT ON COLUMN applications.code IS 'Unique identifier within tenant (e.g., pension-fund, backoffice). Immutable.';
COMMENT ON COLUMN applications.settings IS 'Application-level overrides for tenant settings (highest priority in hierarchy). See TRD Section 3.1.6.';
