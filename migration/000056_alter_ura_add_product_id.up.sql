-- Add product_id to user_role_assignments for product-scoped role assignments
-- The entity already references product_id but the column was missing from the table

ALTER TABLE user_role_assignments
    ADD COLUMN product_id UUID;

-- FK to applications table
ALTER TABLE user_role_assignments
    ADD CONSTRAINT fk_ura_product FOREIGN KEY (product_id)
        REFERENCES applications(id) ON DELETE CASCADE;

-- Drop old unique index and create new one including product_id
DROP INDEX IF EXISTS idx_ura_user_role_branch_unique;

CREATE UNIQUE INDEX idx_ura_user_role_branch_product_unique
    ON user_role_assignments (
        user_id,
        role_id,
        COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid),
        COALESCE(product_id, '00000000-0000-0000-0000-000000000000'::uuid)
    )
    WHERE deleted_at IS NULL;

-- Index for querying active role assignments by user and product
CREATE INDEX idx_ura_user_product_active
    ON user_role_assignments (user_id, product_id)
    WHERE deleted_at IS NULL;

-- Standalone FK index for ON DELETE CASCADE enforcement on applications
CREATE INDEX idx_ura_product_id
    ON user_role_assignments (product_id)
    WHERE deleted_at IS NULL;

COMMENT ON COLUMN user_role_assignments.product_id IS 'Product this role assignment is scoped to. NULL means tenant-wide.';
