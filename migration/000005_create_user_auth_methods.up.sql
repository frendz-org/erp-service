CREATE TABLE user_auth_methods (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    user_id             UUID NOT NULL,

    -- Method Identification
    method_type         VARCHAR(30) NOT NULL,

    -- Method-Specific Credential Data (encrypted at application level for sensitive data)
    credential_data     JSONB NOT NULL,

    -- Status
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_auth_methods_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_auth_methods UNIQUE (user_id, method_type),
    CONSTRAINT chk_user_auth_methods_type CHECK (method_type IN (
        'PASSWORD',
        'PIN',
        'GOOGLE',
        'APPLE',
        'MICROSOFT',
        'WEBAUTHN'
    ))
);

CREATE TRIGGER trg_user_auth_methods_updated_at
    BEFORE UPDATE ON user_auth_methods
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE user_auth_methods IS 'Normalized auth methods. One row per method per user. Extensible without schema changes.';
COMMENT ON COLUMN user_auth_methods.method_type IS 'Authentication method identifier. Adding new methods only requires adding to CHECK constraint.';
COMMENT ON COLUMN user_auth_methods.credential_data IS 'Method-specific data as JSONB. Sensitive values (hashes) are application-level encrypted.';
COMMENT ON COLUMN user_auth_methods.is_active IS 'FALSE to temporarily disable a method without deleting it.';
