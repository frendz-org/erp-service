CREATE TABLE roles (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    application_id      UUID NOT NULL,

    -- Business Identifiers
    code                VARCHAR(50) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    description         TEXT,

    -- Role Type
    is_system           BOOLEAN NOT NULL DEFAULT FALSE,

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
    CONSTRAINT fk_roles_application FOREIGN KEY (application_id)
        REFERENCES applications(id) ON DELETE CASCADE,
    CONSTRAINT uq_roles_application_code UNIQUE (application_id, code),
    CONSTRAINT chk_roles_status CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE TRIGGER trg_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE roles IS 'Application-specific roles that bundle permissions together';
COMMENT ON COLUMN roles.code IS 'Unique within application (e.g., ADMIN, HR_STAFF, VIEWER). Immutable.';
COMMENT ON COLUMN roles.is_system IS 'System roles (e.g., PLATFORM_ADMIN) cannot be modified or deleted';
