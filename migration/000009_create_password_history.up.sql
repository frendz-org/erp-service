CREATE TABLE password_history (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    user_id             UUID NOT NULL,

    -- Password Hash
    password_hash       VARCHAR(255) NOT NULL,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_password_history_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE
);

COMMENT ON TABLE password_history IS 'Stores last N password hashes to prevent reuse (N configured per tenant)';
COMMENT ON COLUMN password_history.password_hash IS 'Bcrypt hash - compare new password against these to prevent reuse';
