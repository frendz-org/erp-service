-- ============================================================================
-- Migration: Migrate insight_prod_users into IAM platform users
-- Source: insight_prod_users (active users only)
-- Target: users, user_profiles, user_auth_methods, user_security_states
--
-- Prerequisites:
--   1. insight_prod_users table must exist in the database
--
-- Filters:
--   - Only status = 'active' (skip pending, inactive, suspended, deleted)
--   - Only is_deleted = false
--   - Deduplicate by email (keep lowest id per email)
--   - Skip emails already present in IAM users table
-- ============================================================================

-- Step 1: Expand registration_source CHECK constraint to include DATA_MIGRATION
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_registration_source;
ALTER TABLE users ADD CONSTRAINT chk_users_registration_source CHECK (
    registration_source IN ('SELF', 'ADMIN', 'IMPORT', 'GOOGLE', 'DATA_MIGRATION')
);

-- Step 2: Create temporary mapping table with UUIDv7 for each eligible user
CREATE TEMP TABLE _migration_user_map (
    old_id        INTEGER PRIMARY KEY,
    new_user_id   UUID NOT NULL DEFAULT uuidv7(),
    email         TEXT NOT NULL,
    full_name     TEXT,
    has_password  BOOLEAN NOT NULL DEFAULT FALSE,
    has_google_id BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO _migration_user_map (old_id, email, full_name, has_password, has_google_id)
SELECT DISTINCT ON (LOWER(TRIM(src.email)))
    src.id,
    LOWER(TRIM(src.email)),
    TRIM(src.full_name),
    (src.password IS NOT NULL AND src.password != ''),
    (src.google_id IS NOT NULL AND src.google_id != '')
FROM insight_prod_users src
WHERE src.is_deleted = FALSE
  AND src.status = 'active'
  AND NOT EXISTS (
      SELECT 1 FROM users u
      WHERE LOWER(u.email) = LOWER(TRIM(src.email))
        AND u.deleted_at IS NULL
  )
ORDER BY LOWER(TRIM(src.email)), src.id ASC;

-- Step 3: Insert into users
INSERT INTO users (id, email, status, status_changed_at, registration_source, version, created_at, updated_at)
SELECT
    m.new_user_id,
    m.email,
    'ACTIVE',
    NOW(),
    'DATA_MIGRATION',
    1,
    COALESCE(src.created_at, NOW()),
    COALESCE(src.updated_at, NOW())
FROM _migration_user_map m
JOIN insight_prod_users src ON src.id = m.old_id;

-- Step 4: Insert into user_profiles
-- Split full_name on last space into first_name + last_name
-- Store unmapped fields in metadata JSONB
INSERT INTO user_profiles (user_id, first_name, last_name, phone_number, date_of_birth, id_number, profile_picture_url, metadata, updated_at)
SELECT
    m.new_user_id,
    -- first_name: everything before the last space, or entire name if no space (max 100 chars)
    LEFT(CASE
        WHEN m.full_name IS NULL OR TRIM(m.full_name) = '' THEN 'User'
        WHEN POSITION(' ' IN TRIM(m.full_name)) = 0 THEN TRIM(m.full_name)
        ELSE TRIM(LEFT(TRIM(m.full_name), LENGTH(TRIM(m.full_name)) - LENGTH(SUBSTRING(TRIM(m.full_name) FROM '[^ ]+$')) - 1))
    END, 100),
    -- last_name: last word after the last space, or empty string (max 100 chars)
    LEFT(CASE
        WHEN m.full_name IS NULL OR TRIM(m.full_name) = '' THEN ''
        WHEN POSITION(' ' IN TRIM(m.full_name)) = 0 THEN ''
        ELSE TRIM(SUBSTRING(TRIM(m.full_name) FROM '[^ ]+$'))
    END, 100),
    LEFT(src.phone_no, 20),
    src.birth_date::DATE,
    LEFT(src.no_ktp, 50),
    LEFT(src.profile_pic, 500),
    jsonb_strip_nulls(jsonb_build_object(
        'migration_source',    'insight_prod_users',
        'migration_source_id', src.id,
        'migrated_at',         NOW()::TEXT,
        'nik_emp',             src.nik_emp,
        'opu_code',            src.opu_code,
        'division',            src.division,
        'job_class',           src.job_class,
        'heirs',               src.heirs,
        'status_peserta',      src.status_peserta,
        'old_emp_no',          src.old_emp_no,
        'type',                src.type,
        'kepatuhan',           src.kepatuhan,
        'join_date',           src.join_date::TEXT,
        'd_join_date',         src.d_join_date::TEXT,
        'no_ktp',              src.no_ktp,
        'old_role',            src.role
    )),
    NOW()
FROM _migration_user_map m
JOIN insight_prod_users src ON src.id = m.old_id;

-- Step 5a: Insert PASSWORD auth methods (only for users with passwords)
INSERT INTO user_auth_methods (id, user_id, method_type, credential_data, is_active, created_at, updated_at)
SELECT
    uuidv7(),
    m.new_user_id,
    'PASSWORD',
    jsonb_build_object(
        'password_hash', src.password,
        'password_history', '[]'::jsonb
    ),
    TRUE,
    COALESCE(src.created_at, NOW()),
    COALESCE(src.updated_at, NOW())
FROM _migration_user_map m
JOIN insight_prod_users src ON src.id = m.old_id
WHERE m.has_password = TRUE;

-- Step 5b: Insert GOOGLE auth methods (only for users with google_id)
INSERT INTO user_auth_methods (id, user_id, method_type, credential_data, is_active, created_at, updated_at)
SELECT
    uuidv7(),
    m.new_user_id,
    'GOOGLE',
    jsonb_build_object(
        'google_id',      src.google_id,
        'email',          m.email,
        'email_verified', TRUE,
        'name',           COALESCE(TRIM(src.full_name), ''),
        'picture',        COALESCE(src.profile_pic, '')
    ),
    TRUE,
    COALESCE(src.created_at, NOW()),
    COALESCE(src.updated_at, NOW())
FROM _migration_user_map m
JOIN insight_prod_users src ON src.id = m.old_id
WHERE m.has_google_id = TRUE;

-- Step 6: Insert user_security_states
INSERT INTO user_security_states (
    user_id, failed_login_attempts, failed_pin_attempts,
    email_verified, email_verified_at, pin_verified,
    force_password_change, updated_at
)
SELECT
    m.new_user_id,
    0,
    0,
    TRUE,
    NOW(),
    FALSE,
    FALSE,
    NOW()
FROM _migration_user_map m;

-- Step 7: Cleanup
DROP TABLE IF EXISTS _migration_user_map;
