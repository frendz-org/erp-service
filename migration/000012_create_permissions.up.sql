CREATE TABLE permissions (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    application_id      UUID NOT NULL,

    -- Business Identifiers
    code                VARCHAR(100) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    description         TEXT,

    -- Permission Classification (for UI grouping)
    resource_type       VARCHAR(50),
    action              VARCHAR(50),

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
    CONSTRAINT fk_permissions_application FOREIGN KEY (application_id)
        REFERENCES applications(id) ON DELETE CASCADE,
    CONSTRAINT uq_permissions_application_code UNIQUE (application_id, code),
    CONSTRAINT chk_permissions_status CHECK (status IN ('ACTIVE', 'INACTIVE')),
    CONSTRAINT chk_permissions_code_format CHECK (code ~ '^[a-z_]+:[a-z_]+$')
);

CREATE TRIGGER trg_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE permissions IS 'Fine-grained permissions in format resource:action';
COMMENT ON COLUMN permissions.code IS 'Permission code in format resource:action (e.g., user:create, loan:approve). Immutable.';
