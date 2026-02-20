CREATE TABLE user_role_assignments (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Core Assignment
    user_id             UUID NOT NULL,
    role_id             UUID NOT NULL,

    -- Optional Scoping
    branch_id           UUID,

    -- Assignment Metadata
    assigned_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by         UUID,
    expires_at          TIMESTAMPTZ,

    -- Status Management
    status              VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT fk_ura_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_ura_role FOREIGN KEY (role_id)
        REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_ura_branch FOREIGN KEY (branch_id)
        REFERENCES branches(id) ON DELETE CASCADE,
    CONSTRAINT chk_ura_status CHECK (status IN ('ACTIVE', 'INACTIVE', 'EXPIRED'))
);

-- Unique index using COALESCE to handle NULL branch_id
-- This prevents duplicate (user_id, role_id, branch_id) combinations
-- where NULL branch_id is treated as a specific value
CREATE UNIQUE INDEX idx_ura_user_role_branch_unique
    ON user_role_assignments (user_id, role_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid))
    WHERE deleted_at IS NULL;

CREATE TRIGGER trg_ura_updated_at
    BEFORE UPDATE ON user_role_assignments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE user_role_assignments IS 'Assigns roles to users, with optional branch-level scoping';
COMMENT ON COLUMN user_role_assignments.branch_id IS 'NULL means tenant-wide; specific ID means branch-scoped';
COMMENT ON COLUMN user_role_assignments.expires_at IS 'Optional expiration for time-bound roles (e.g., temporary authority)';
