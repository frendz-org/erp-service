-- Members table: saving-domain operator/admin registration data.
-- Extends user_tenant_registrations (1-to-1) for MEMBER type registrations.
-- Mirrors the participants table pattern.
--
-- Cross-domain boundary: tenant_id, product_id, user_id are denormalized UUIDs (no FK)
-- following the same pattern as participants.tenant_id / participants.product_id.
-- user_tenant_registration_id has a FK because members is a dependent extension of UTR.

CREATE TABLE IF NOT EXISTS members (
    -- Primary Key
    id                          UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Registration Link (FK to UTR - lifecycle coupling)
    user_tenant_registration_id UUID NOT NULL,

    -- Tenant/Product Scoping (denormalized, no FK - cross-domain boundary)
    tenant_id                   UUID NOT NULL,
    product_id                  UUID NOT NULL,

    -- User Link (denormalized from UTR, no FK - cross-domain boundary)
    user_id                     UUID NOT NULL,

    -- Member-specific data
    participant_number          VARCHAR(50)  NOT NULL,
    identity_number             VARCHAR(16)  NOT NULL,
    organization_code           VARCHAR(50)  NOT NULL,

    -- Profile snapshot (captured at registration time, immutable)
    full_name                   VARCHAR(255) NOT NULL,
    gender                      VARCHAR(20),
    date_of_birth               DATE,

    -- Audit Fields
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ,

    -- Optimistic Locking
    version                     INTEGER NOT NULL DEFAULT 1,

    -- FK to UTR (lifecycle coupling, RESTRICT prevents orphan deletion)
    CONSTRAINT fk_members_user_tenant_registration
        FOREIGN KEY (user_tenant_registration_id) REFERENCES user_tenant_registrations(id)
        ON DELETE RESTRICT
);

-- 1-to-1 with UTR: each UTR has at most one members row
CREATE UNIQUE INDEX uk_members_user_tenant_registration_id
    ON members(user_tenant_registration_id)
    WHERE deleted_at IS NULL;

-- Prevent duplicate participant_number per tenant+product
CREATE UNIQUE INDEX uk_members_tenant_product_participant_number
    ON members(tenant_id, product_id, participant_number)
    WHERE deleted_at IS NULL;

-- One member per user per tenant+product
CREATE UNIQUE INDEX uk_members_user_tenant_product
    ON members(user_id, tenant_id, product_id)
    WHERE deleted_at IS NULL;

-- Tenant+product scoped list queries
CREATE INDEX idx_members_tenant_product
    ON members(tenant_id, product_id)
    WHERE deleted_at IS NULL;

-- Soft-delete filter
CREATE INDEX idx_members_deleted_at
    ON members(deleted_at);

-- Auto-update updated_at
CREATE TRIGGER trg_members_updated_at
    BEFORE UPDATE ON members
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE members IS 'Saving-domain operator/admin registration data. Extends user_tenant_registrations (1-to-1) for MEMBER type. Mirrors participants table pattern.';
COMMENT ON COLUMN members.user_tenant_registration_id IS 'FK to user_tenant_registrations.id — the UTR this member extends';
COMMENT ON COLUMN members.tenant_id IS 'Denormalized from UTR. No FK — cross-domain boundary.';
COMMENT ON COLUMN members.product_id IS 'Denormalized from UTR. No FK — cross-domain boundary.';
COMMENT ON COLUMN members.user_id IS 'Denormalized from UTR. No FK — cross-domain boundary.';
COMMENT ON COLUMN members.participant_number IS 'Employee number validated against employee_data.emp_no. No FK — external table.';
COMMENT ON COLUMN members.identity_number IS '16-digit Indonesian NIK (Nomor Induk Kependudukan).';
COMMENT ON COLUMN members.organization_code IS 'Masterdata TENANT code. Validated at registration time, no FK.';
COMMENT ON COLUMN members.full_name IS 'Profile snapshot at registration time.';
COMMENT ON COLUMN members.gender IS 'Profile snapshot: gender masterdata code at registration time.';
COMMENT ON COLUMN members.date_of_birth IS 'Profile snapshot: date of birth at registration time.';
