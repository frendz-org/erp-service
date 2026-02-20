-- Append-only table: tokens are created and optionally revoked.
-- Exception to database-conventions Rule 2: NO updated_at, deleted_at, or version columns.
-- Rationale: refresh tokens are never updated (only revoked_at is SET once), never soft-deleted,
-- and do not require optimistic locking.
-- NO tenant_id: refresh tokens are platform-level. JWT carries all tenant associations.

CREATE TABLE IF NOT EXISTS refresh_tokens (
    -- Primary Key
    id                   UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Owner
    user_id              UUID NOT NULL,

    -- Token Data
    token_hash           VARCHAR(64) NOT NULL,
    token_family         UUID NOT NULL DEFAULT uuidv7(),

    -- Lifecycle
    expires_at           TIMESTAMPTZ NOT NULL,
    revoked_at           TIMESTAMPTZ,
    revoked_reason       VARCHAR(100),
    replaced_by_token_id UUID,

    -- Client Context
    ip_address           INET,
    user_agent           TEXT,

    -- Audit
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign Keys
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_refresh_tokens_replaced_by FOREIGN KEY (replaced_by_token_id)
        REFERENCES refresh_tokens(id) ON DELETE SET NULL
);

-- Token lookup by hash (login/refresh flow, CRITICAL path)
CREATE UNIQUE INDEX IF NOT EXISTS uq_refresh_tokens_token_hash
    ON refresh_tokens(token_hash);

-- FK index: find all tokens for a user (revoke-all-on-logout, admin view)
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id
    ON refresh_tokens(user_id);

-- Token family index: revoke entire family on rotation theft detection
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_family
    ON refresh_tokens(token_family);

-- Active tokens for a user: used in "list active sessions" and "revoke all active"
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_active
    ON refresh_tokens(user_id)
    WHERE revoked_at IS NULL;

COMMENT ON TABLE refresh_tokens IS 'Append-only refresh token store. Platform-level (no tenant_id). Revocation sets revoked_at once.';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hex digest of the raw refresh token. Never stored in plaintext.';
COMMENT ON COLUMN refresh_tokens.token_family IS 'UUIDv7 grouping rotated tokens. All tokens in a family are revoked on theft detection.';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'Set once when token is revoked. NULL = active token.';
COMMENT ON COLUMN refresh_tokens.revoked_reason IS 'Human-readable reason: logout, rotation, theft_detected, admin_revoke.';
COMMENT ON COLUMN refresh_tokens.replaced_by_token_id IS 'Self-referential FK: points to the new token that replaced this one during rotation.';
