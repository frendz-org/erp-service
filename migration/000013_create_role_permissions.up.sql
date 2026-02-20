CREATE TABLE role_permissions (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    role_id             UUID NOT NULL,
    permission_id       UUID NOT NULL,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID,

    -- Constraints
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id)
        REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id)
        REFERENCES permissions(id) ON DELETE CASCADE,
    CONSTRAINT uq_role_permissions UNIQUE (role_id, permission_id)
);

COMMENT ON TABLE role_permissions IS 'Junction table: which permissions belong to which roles';
