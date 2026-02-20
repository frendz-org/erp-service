CREATE TABLE user_branches (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys (Composite relationship)
    user_id             UUID NOT NULL,
    branch_id           UUID NOT NULL,

    -- Assignment Details
    is_primary          BOOLEAN NOT NULL DEFAULT FALSE,
    assigned_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by         UUID,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_branches_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_branches_branch FOREIGN KEY (branch_id)
        REFERENCES branches(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_branches_user_branch UNIQUE (user_id, branch_id)
);

-- Ensure only one primary branch per user
CREATE UNIQUE INDEX idx_user_branches_primary
    ON user_branches(user_id)
    WHERE is_primary = TRUE;

COMMENT ON TABLE user_branches IS 'Junction table linking users to their assigned branches';
COMMENT ON COLUMN user_branches.is_primary IS 'Only one branch can be primary per user (enforced by partial unique index)';
