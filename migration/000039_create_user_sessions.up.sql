-- Basic session/device tracking for logged-in users.
-- Exception to database-conventions Rule 2: NO deleted_at or version columns.
-- Rationale: sessions transition ACTIVE -> REVOKED/EXPIRED, never soft-deleted.

CREATE TABLE IF NOT EXISTS user_sessions (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Owner
    user_id             UUID NOT NULL,

    -- Linked Refresh Token
    refresh_token_id    UUID,

    -- Client Context
    ip_address          INET NOT NULL,
    user_agent          TEXT,
    device_fingerprint  VARCHAR(255),

    -- Authentication
    login_method        VARCHAR(20) NOT NULL,

    -- Status
    status              VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Activity Tracking
    last_active_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ NOT NULL,
    revoked_at          TIMESTAMPTZ,

    -- Audit
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign Keys
    CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_sessions_refresh_token FOREIGN KEY (refresh_token_id)
        REFERENCES refresh_tokens(id) ON DELETE SET NULL,

    -- Check Constraints
    CONSTRAINT chk_user_sessions_status CHECK (status IN (
        'ACTIVE', 'REVOKED', 'EXPIRED'
    )),
    CONSTRAINT chk_user_sessions_login_method CHECK (login_method IN (
        'EMAIL_OTP'
    ))
);

CREATE TRIGGER trg_user_sessions_updated_at
    BEFORE UPDATE ON user_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Active sessions for a user
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_active
    ON user_sessions(user_id)
    WHERE status = 'ACTIVE';

-- FK index: lookup session by refresh token
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token_id
    ON user_sessions(refresh_token_id);

COMMENT ON TABLE user_sessions IS 'Tracks active login sessions with device context. Sessions transition ACTIVE -> REVOKED/EXPIRED.';
COMMENT ON COLUMN user_sessions.refresh_token_id IS 'Linked refresh token. SET NULL on token deletion.';
COMMENT ON COLUMN user_sessions.device_fingerprint IS 'Reserved for future device fingerprinting.';
COMMENT ON COLUMN user_sessions.login_method IS 'Authentication method: EMAIL_OTP. Extensible via CHECK update.';
COMMENT ON COLUMN user_sessions.last_active_at IS 'Updated on each token refresh. Used for idle timeout.';
