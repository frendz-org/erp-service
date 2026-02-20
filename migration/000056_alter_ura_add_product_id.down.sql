-- Reverse: remove product_id from user_role_assignments

DROP INDEX IF EXISTS idx_ura_product_id;
DROP INDEX IF EXISTS idx_ura_user_product_active;
DROP INDEX IF EXISTS idx_ura_user_role_branch_product_unique;

ALTER TABLE user_role_assignments
    DROP CONSTRAINT IF EXISTS fk_ura_product;

ALTER TABLE user_role_assignments
    DROP COLUMN IF EXISTS product_id;

-- Restore original unique index
CREATE UNIQUE INDEX idx_ura_user_role_branch_unique
    ON user_role_assignments (user_id, role_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid))
    WHERE deleted_at IS NULL;
