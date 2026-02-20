CREATE TABLE mfa_enrollments (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- User Association
    user_id             UUID NOT NULL,

    -- MFA Method
    method_type         VARCHAR(30) NOT NULL,

    -- Method-Specific Credential Data (encrypted at application level)
    credential_data     JSONB NOT NULL,

    -- Configuration
    is_primary          BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified         BOOLEAN NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    -- Usage Tracking
    last_used_at        TIMESTAMPTZ,
    use_count           INTEGER NOT NULL DEFAULT 0,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_mfa_enrollments_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_mfa_method_type CHECK (method_type IN (
        'TOTP',
        'SMS',
        'EMAIL',
        'PUSH',
        'WEBAUTHN',
        'BACKUP_CODES'
    ))
);

-- Ensure only one primary MFA per user
CREATE UNIQUE INDEX idx_mfa_enrollments_primary
    ON mfa_enrollments(user_id)
    WHERE is_primary = TRUE AND is_active = TRUE;

-- Prevent duplicate method types per user
CREATE UNIQUE INDEX idx_mfa_enrollments_user_method
    ON mfa_enrollments(user_id, method_type);

-- Trigger for updated_at
CREATE TRIGGER trg_mfa_enrollments_updated_at
    BEFORE UPDATE ON mfa_enrollments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE mfa_enrollments IS 'User MFA method configurations. One row per MFA method per user.';
COMMENT ON COLUMN mfa_enrollments.method_type IS 'MFA method: TOTP, SMS, EMAIL, PUSH, WEBAUTHN, BACKUP_CODES';
COMMENT ON COLUMN mfa_enrollments.credential_data IS 'Method-specific data: TOTP secret, phone number, device tokens, etc. (encrypted at app level)';
COMMENT ON COLUMN mfa_enrollments.is_primary IS 'Default MFA method presented to user. Only one active primary per user.';
COMMENT ON COLUMN mfa_enrollments.is_verified IS 'FALSE during setup, TRUE after user proves they can use the method.';
COMMENT ON COLUMN mfa_enrollments.is_active IS 'Can be temporarily disabled without removing enrollment.';
COMMENT ON COLUMN mfa_enrollments.use_count IS 'Number of times this MFA method has been used successfully.';
