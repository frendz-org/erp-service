CREATE TABLE recovery_codes (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- User Association
    user_id             UUID NOT NULL,

    -- Code (hashed)
    code_hash           VARCHAR(255) NOT NULL,

    -- Usage Tracking
    is_used             BOOLEAN NOT NULL DEFAULT FALSE,
    used_at             TIMESTAMPTZ,
    used_ip             INET,
    used_user_agent     TEXT,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_recovery_codes_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE
);

-- Index for user lookup (unused codes)
CREATE INDEX idx_recovery_codes_user_unused
    ON recovery_codes(user_id)
    WHERE is_used = FALSE;

-- Comments
COMMENT ON TABLE recovery_codes IS 'One-time use backup codes. Typically 10 codes generated at MFA setup.';
COMMENT ON COLUMN recovery_codes.code_hash IS 'Bcrypt hash. Code format: XXXX-XXXX (8 alphanumeric, shown to user once)';
COMMENT ON COLUMN recovery_codes.is_used IS 'TRUE after code has been used for recovery. Cannot be reused.';
COMMENT ON COLUMN recovery_codes.used_at IS 'Timestamp when the code was used.';
COMMENT ON COLUMN recovery_codes.used_ip IS 'IP address from which the code was used.';
COMMENT ON COLUMN recovery_codes.used_user_agent IS 'User agent from which the code was used.';
