CREATE TABLE product_registration_configs (
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    application_id      UUID NOT NULL,

    registration_type   VARCHAR(20) NOT NULL,           
    auto_grant_role_id  UUID,                           
    requires_approval   BOOLEAN NOT NULL DEFAULT TRUE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_prc_application FOREIGN KEY (application_id)
        REFERENCES applications(id) ON DELETE CASCADE,
    CONSTRAINT fk_prc_role FOREIGN KEY (auto_grant_role_id)
        REFERENCES roles(id) ON DELETE SET NULL,
    CONSTRAINT uq_prc_app_reg_type UNIQUE (application_id, registration_type),
    CONSTRAINT chk_prc_registration_type CHECK (registration_type IN (
        'PARTICIPANT',
        'MEMBER'
    ))
);

CREATE TRIGGER trg_product_registration_configs_updated_at
    BEFORE UPDATE ON product_registration_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
COMMENT ON TABLE product_registration_configs IS 'Per-product registration configuration: auto-grant roles and approval settings';
COMMENT ON COLUMN product_registration_configs.application_id IS 'Reference to the product (applications table)';
COMMENT ON COLUMN product_registration_configs.registration_type IS 'PARTICIPANT = typically auto-approved; MEMBER = typically requires approval';
COMMENT ON COLUMN product_registration_configs.auto_grant_role_id IS 'Role automatically assigned when registration is approved/auto-approved';
COMMENT ON COLUMN product_registration_configs.requires_approval IS 'FALSE = auto-approved (typical for PARTICIPANT); TRUE = needs TENANT_PRODUCT_ADMIN approval';
COMMENT ON COLUMN product_registration_configs.is_active IS 'Whether this registration path is currently accepting new registrations';
