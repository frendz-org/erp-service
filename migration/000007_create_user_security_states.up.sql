CREATE TABLE user_security_states (
    -- Primary Key (same as user_id for 1:1 relationship)
    user_id                 UUID PRIMARY KEY,

    -- Failed Attempt Tracking
    failed_login_attempts   INTEGER NOT NULL DEFAULT 0,
    failed_pin_attempts     INTEGER NOT NULL DEFAULT 0,
    locked_until            TIMESTAMPTZ,

    -- Login History
    last_login_at           TIMESTAMPTZ,
    last_login_ip           INET,

    -- Email Verification
    email_verified          BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified_at       TIMESTAMPTZ,

    -- PIN Verification
    pin_verified            BOOLEAN NOT NULL DEFAULT FALSE,

    -- Password Policy Enforcement
    force_password_change   BOOLEAN NOT NULL DEFAULT FALSE,

    -- Audit Fields
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_security_states_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_user_security_failed_login CHECK (failed_login_attempts >= 0),
    CONSTRAINT chk_user_security_failed_pin CHECK (failed_pin_attempts >= 0)
);

CREATE TRIGGER trg_user_security_states_updated_at
    BEFORE UPDATE ON user_security_states
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE user_security_states IS 'Security state that changes frequently during authentication. No version column — uses atomic operations.';
COMMENT ON COLUMN user_security_states.failed_login_attempts IS 'Counter for failed password attempts, resets on success. Incremented atomically.';
COMMENT ON COLUMN user_security_states.failed_pin_attempts IS 'Counter for failed PIN attempts, resets on success. Incremented atomically.';
COMMENT ON COLUMN user_security_states.locked_until IS 'Account is locked until this timestamp (NULL = not locked)';
COMMENT ON COLUMN user_security_states.force_password_change IS 'When TRUE, user must change password on next login. Set by admin when password policy changes or on security events.';
