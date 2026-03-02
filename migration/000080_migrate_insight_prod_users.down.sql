-- ============================================================================
-- Rollback: Delete all users migrated from insight_prod_users
-- Identified by user_profiles.metadata->>'migration_source' = 'insight_prod_users'
-- ============================================================================

-- Step 1: Collect migrated user IDs
CREATE TEMP TABLE _rollback_user_ids AS
SELECT up.user_id
FROM user_profiles up
WHERE up.metadata->>'migration_source' = 'insight_prod_users';

-- Step 2: Delete from user_security_states
DELETE FROM user_security_states
WHERE user_id IN (SELECT user_id FROM _rollback_user_ids);

-- Step 3: Delete from user_auth_methods
DELETE FROM user_auth_methods
WHERE user_id IN (SELECT user_id FROM _rollback_user_ids);

-- Step 4: Delete from user_profiles
DELETE FROM user_profiles
WHERE user_id IN (SELECT user_id FROM _rollback_user_ids);

-- Step 5: Delete from users
DELETE FROM users
WHERE id IN (SELECT user_id FROM _rollback_user_ids);

-- Step 6: Revert registration_source CHECK constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_registration_source;
ALTER TABLE users ADD CONSTRAINT chk_users_registration_source CHECK (
    registration_source IN ('SELF', 'ADMIN', 'IMPORT', 'GOOGLE')
);

-- Step 7: Cleanup
DROP TABLE IF EXISTS _rollback_user_ids;
