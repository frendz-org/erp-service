CREATE TABLE user_tenant_registrations (
    id                    UUID PRIMARY KEY DEFAULT uuidv7(),

    user_id               UUID NOT NULL,
    tenant_id             UUID NOT NULL,

    registration_type     VARCHAR(20) NOT NULL,         
    identification_number VARCHAR(100),                 

    status                VARCHAR(20) NOT NULL DEFAULT 'PENDING_APPROVAL',

    approved_by           UUID,
    approved_at           TIMESTAMPTZ,

    metadata              JSONB NOT NULL DEFAULT '{}',

    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,

    version               INTEGER NOT NULL DEFAULT 1,

    CONSTRAINT fk_utr_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT fk_utr_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(id) ON DELETE RESTRICT,
    CONSTRAINT fk_utr_approved_by FOREIGN KEY (approved_by)
        REFERENCES users(id),
    CONSTRAINT uq_utr_user_tenant_type UNIQUE (user_id, tenant_id, registration_type),
    CONSTRAINT chk_utr_registration_type CHECK (registration_type IN (
        'PARTICIPANT',
        'MEMBER'
    )),
    CONSTRAINT chk_utr_status CHECK (status IN (
        'PENDING_APPROVAL',
        'ACTIVE',
        'REJECTED',
        'INACTIVE'
    ))
);

CREATE TRIGGER trg_user_tenant_registrations_updated_at
    BEFORE UPDATE ON user_tenant_registrations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
COMMENT ON TABLE user_tenant_registrations IS 'Links platform-level users to tenants via participant or member registration';
COMMENT ON COLUMN user_tenant_registrations.user_id IS 'Reference to the platform-level user';
COMMENT ON COLUMN user_tenant_registrations.tenant_id IS 'Reference to the tenant the user is registering with';
COMMENT ON COLUMN user_tenant_registrations.registration_type IS 'PARTICIPANT = auto-approved with auto-granted roles; MEMBER = requires product-level approval';
COMMENT ON COLUMN user_tenant_registrations.identification_number IS 'External identifier: ParticipantNumber for PARTICIPANT, EmployeeNumber for MEMBER';
COMMENT ON COLUMN user_tenant_registrations.status IS 'PENDING_APPROVAL (member only) → ACTIVE → INACTIVE. PARTICIPANT auto-set to ACTIVE.';
COMMENT ON COLUMN user_tenant_registrations.approved_by IS 'User ID of the TENANT_PRODUCT_ADMIN who approved a MEMBER registration';
COMMENT ON COLUMN user_tenant_registrations.approved_at IS 'Timestamp when a MEMBER registration was approved';
COMMENT ON COLUMN user_tenant_registrations.metadata IS 'Flexible context: {branch_id, department, notes, etc.}';
