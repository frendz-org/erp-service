CREATE TABLE IF NOT EXISTS participants (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Tenant Scoping (UUID column, NOT FK to tenants - cross-domain boundary)
    tenant_id           UUID NOT NULL,

    -- Application Scoping (UUID column, NOT FK to applications - cross-domain boundary)
    application_id      UUID NOT NULL,

    -- User Linking (UUID column, NOT FK to users - cross-domain boundary)
    -- NULL until deferred user-linking matches this participant to a user.
    -- Unique per (user_id, tenant_id, application_id) - enforced by partial unique index.
    user_id             UUID NULL,

    -- Personal Data
    full_name           VARCHAR(255) NOT NULL,
    gender              VARCHAR(50) NULL,
    place_of_birth      VARCHAR(255) NULL,
    date_of_birth       DATE NULL,
    marital_status      VARCHAR(50) NULL,
    citizenship         VARCHAR(50) NULL,
    religion            VARCHAR(50) NULL,

    -- Identification (for deferred user-linking matching)
    ktp_number          VARCHAR(50) NULL,
    employee_number     VARCHAR(50) NULL,
    phone_number        VARCHAR(50) NULL,

    -- Approval Workflow
    status              VARCHAR(20) NOT NULL DEFAULT 'DRAFT',

    -- Workflow Actors (UUID columns, NOT FK to users - cross-domain boundary)
    created_by          UUID NOT NULL,
    submitted_by        UUID NULL,
    submitted_at        TIMESTAMPTZ NULL,
    approved_by         UUID NULL,
    approved_at         TIMESTAMPTZ NULL,
    rejected_by         UUID NULL,
    rejected_at         TIMESTAMPTZ NULL,
    rejection_reason    TEXT NULL,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT chk_participants_status CHECK (status IN (
        'DRAFT', 'PENDING_APPROVAL', 'APPROVED', 'REJECTED'
    ))
);

CREATE TRIGGER trg_participants_updated_at
    BEFORE UPDATE ON participants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participants IS 'Core participant record with personal data and approval workflow. Belongs to Participant domain. tenant_id and application_id are UUID references (no FK) respecting cross-domain boundary.';
COMMENT ON COLUMN participants.tenant_id IS 'UUID of owning tenant. No FK - cross-domain boundary to Organization domain.';
COMMENT ON COLUMN participants.application_id IS 'UUID of owning application. No FK - cross-domain boundary to Authorization domain.';
COMMENT ON COLUMN participants.user_id IS 'UUID of linked user. No FK - cross-domain boundary to Identity domain. NULL until deferred user-linking Phase 2. 1:1 per tenant+application enforced by partial unique index.';
COMMENT ON COLUMN participants.gender IS 'Masterdata string code from GENDER category. Validated at application level, not DB FK.';
COMMENT ON COLUMN participants.marital_status IS 'Masterdata string code from MARITAL_STATUS category.';
COMMENT ON COLUMN participants.citizenship IS 'Masterdata string code from CITIZENSHIP category.';
COMMENT ON COLUMN participants.religion IS 'Masterdata string code from RELIGION category.';
COMMENT ON COLUMN participants.ktp_number IS 'Indonesian national ID number. Used for deferred user-linking matching.';
COMMENT ON COLUMN participants.employee_number IS 'Employee/personnel number. Used for deferred user-linking matching.';
COMMENT ON COLUMN participants.status IS 'Approval workflow: DRAFT -> PENDING_APPROVAL -> APPROVED/REJECTED.';
COMMENT ON COLUMN participants.created_by IS 'UUID of user who created this record. No FK - cross-domain boundary to Identity domain.';
COMMENT ON COLUMN participants.submitted_by IS 'UUID of user who submitted for approval.';
COMMENT ON COLUMN participants.approved_by IS 'UUID of user who approved this record.';
COMMENT ON COLUMN participants.rejected_by IS 'UUID of user who rejected this record.';
COMMENT ON COLUMN participants.rejection_reason IS 'Free-text reason when status is REJECTED.';
