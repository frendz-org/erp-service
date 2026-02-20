CREATE TABLE branches (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,

    -- Business Identifiers
    code                VARCHAR(50) NOT NULL,
    name                VARCHAR(255) NOT NULL,

    -- Optional Metadata
    address             TEXT,
    metadata            JSONB NOT NULL DEFAULT '{}',

    -- Status Management
    status              VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_branches_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(id) ON DELETE RESTRICT,
    CONSTRAINT uq_branches_tenant_code UNIQUE (tenant_id, code),
    CONSTRAINT chk_branches_status CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE TRIGGER trg_branches_updated_at
    BEFORE UPDATE ON branches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE branches IS 'Organizational units within a tenant for geographic or departmental scoping';
COMMENT ON COLUMN branches.code IS 'Unique identifier within tenant (e.g., HQ, BRANCH-001). Immutable once created.';
COMMENT ON COLUMN branches.metadata IS 'Tenant-specific branch data (e.g., swift_code, region). IAM stores but does not interpret.';
