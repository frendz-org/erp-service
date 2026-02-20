CREATE TABLE user_profiles (
    -- Primary Key (same as user_id for 1:1 relationship)
    user_id             UUID PRIMARY KEY,

    -- Required Profile Fields
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,

    -- Optional Profile Fields
    phone_number        VARCHAR(20),
    date_of_birth       DATE,
    gender              VARCHAR(20),
    marital_status      VARCHAR(20),
    address             TEXT,
    id_number           VARCHAR(50),
    profile_picture_url VARCHAR(500),

    -- Flexible Metadata (tenant-specific fields IAM stores but does not interpret)
    metadata            JSONB NOT NULL DEFAULT '{}',

    -- Audit Fields
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER trg_user_profiles_updated_at
    BEFORE UPDATE ON user_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE user_profiles IS 'Personal profile information separated from user identity';
COMMENT ON COLUMN user_profiles.gender IS 'Masterdata code string (e.g., MALE). Validated against Masterdata service at write time.';
COMMENT ON COLUMN user_profiles.marital_status IS 'Masterdata code string (e.g., SINGLE). Validated against Masterdata service at write time.';
COMMENT ON COLUMN user_profiles.metadata IS 'Tenant-specific extensible fields (e.g., employee_id, department). IAM stores, consuming apps interpret.';
COMMENT ON COLUMN user_profiles.id_number IS 'National ID - consider application-level encryption';
