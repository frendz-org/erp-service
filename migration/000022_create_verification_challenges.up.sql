CREATE TABLE verification_challenges (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Tenant Scoping
    tenant_id           UUID NOT NULL,

    -- User Association (NULL for pre-registration flows)
    user_id             UUID,

    -- Identifier (what we're verifying)
    identifier          VARCHAR(255) NOT NULL,
    identifier_type     VARCHAR(20) NOT NULL,

    -- Challenge Type & Purpose
    challenge_type      VARCHAR(30) NOT NULL,
    purpose             VARCHAR(50) NOT NULL,

    -- Challenge Data (only one will be populated based on challenge_type)
    otp_hash            VARCHAR(255),
    token_hash          VARCHAR(255),

    -- Flexible Metadata
    metadata            JSONB NOT NULL DEFAULT '{}',

    -- Status Tracking
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    -- Attempt Tracking
    attempts            INTEGER NOT NULL DEFAULT 0,
    max_attempts        INTEGER NOT NULL DEFAULT 5,

    -- Resend Tracking
    resend_count        INTEGER NOT NULL DEFAULT 0,
    max_resends         INTEGER NOT NULL DEFAULT 3,
    last_resent_at      TIMESTAMPTZ,
    resend_cooldown_seconds INTEGER NOT NULL DEFAULT 60,

    -- Lifecycle Timestamps
    expires_at          TIMESTAMPTZ NOT NULL,
    verified_at         TIMESTAMPTZ,
    consumed_at         TIMESTAMPTZ,

    -- Security Context
    ip_address          INET,
    user_agent          TEXT,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_verification_challenges_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_verification_challenges_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_vc_identifier_type CHECK (identifier_type IN (
        'EMAIL',
        'PHONE',
        'USER_ID'
    )),
    CONSTRAINT chk_vc_challenge_type CHECK (challenge_type IN (
        'OTP',
        'MAGIC_LINK',
        'TOTP',
        'PUSH',
        'WEBAUTHN'
    )),
    CONSTRAINT chk_vc_purpose CHECK (purpose IN (
        'REGISTRATION',
        'LOGIN',
        'PASSWORD_RESET',
        'PASSWORD_CHANGE',
        'EMAIL_CHANGE',
        'PHONE_CHANGE',
        'MFA_SETUP',
        'MFA_LOGIN',
        'MFA_SENSITIVE_OP',
        'ACCOUNT_RECOVERY',
        'ADMIN_ACTION'
    )),
    CONSTRAINT chk_vc_status CHECK (status IN (
        'PENDING',
        'VERIFIED',
        'CONSUMED',
        'EXPIRED',
        'FAILED',
        'CANCELLED'
    )),
    CONSTRAINT chk_vc_attempts CHECK (attempts >= 0),
    CONSTRAINT chk_vc_resend_count CHECK (resend_count >= 0),
    CONSTRAINT chk_vc_has_challenge CHECK (
        otp_hash IS NOT NULL OR token_hash IS NOT NULL OR challenge_type IN ('PUSH', 'TOTP')
    )
);

-- Comments
COMMENT ON TABLE verification_challenges IS 'Generic verification system supporting OTP, magic links, TOTP for various purposes';
COMMENT ON COLUMN verification_challenges.user_id IS 'NULL for pre-registration flows where user does not exist yet';
COMMENT ON COLUMN verification_challenges.identifier IS 'The value being verified: email address, phone number, or user ID';
COMMENT ON COLUMN verification_challenges.identifier_type IS 'Type of identifier: EMAIL, PHONE, USER_ID';
COMMENT ON COLUMN verification_challenges.challenge_type IS 'Verification method: OTP, MAGIC_LINK, TOTP, PUSH, WEBAUTHN';
COMMENT ON COLUMN verification_challenges.purpose IS 'Business purpose determines flow behavior and permissions granted on success';
COMMENT ON COLUMN verification_challenges.otp_hash IS 'Bcrypt hash of OTP code';
COMMENT ON COLUMN verification_challenges.token_hash IS 'SHA256 hash of magic link token';
COMMENT ON COLUMN verification_challenges.metadata IS 'Flexible context: {redirect_url, device_info, session_id, etc.}';
COMMENT ON COLUMN verification_challenges.status IS 'Challenge state: PENDING, VERIFIED, CONSUMED, EXPIRED, FAILED, CANCELLED';
COMMENT ON COLUMN verification_challenges.verified_at IS 'When the challenge was successfully verified (e.g., correct OTP entered)';
COMMENT ON COLUMN verification_challenges.consumed_at IS 'When the verification was used for its purpose (e.g., password actually changed)';
