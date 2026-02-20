BEGIN;

-- 1. Rename table
ALTER TABLE applications RENAME TO products;

-- 2. Rename trigger (must happen before FK renames since trigger is on the table)
ALTER TRIGGER trg_applications_updated_at ON products RENAME TO trg_products_updated_at;

-- 3. Rename CHECK constraint
ALTER TABLE products RENAME CONSTRAINT chk_applications_status TO chk_products_status;

-- 4. Rename FK and unique constraint on products itself
ALTER TABLE products RENAME CONSTRAINT fk_applications_tenant TO fk_products_tenant;
ALTER TABLE products RENAME CONSTRAINT uq_applications_tenant_code TO uq_products_tenant_code;

-- Rename indexes on products
DROP INDEX IF EXISTS idx_applications_tenant_id;
DROP INDEX IF EXISTS idx_applications_tenant_code_active;
CREATE INDEX idx_products_tenant_id ON products(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_tenant_code_active ON products(tenant_id, code) WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- 5. roles: rename column + constraints + indexes
ALTER TABLE roles DROP CONSTRAINT IF EXISTS fk_roles_application;
ALTER TABLE roles DROP CONSTRAINT IF EXISTS uq_roles_application_code;
DROP INDEX IF EXISTS idx_roles_application_id;
DROP INDEX IF EXISTS idx_roles_app_code_active;
ALTER TABLE roles RENAME COLUMN application_id TO product_id;
ALTER TABLE roles ADD CONSTRAINT fk_roles_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE roles ADD CONSTRAINT uq_roles_product_code UNIQUE (product_id, code);
CREATE INDEX idx_roles_product_id ON roles(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_product_code_active ON roles(product_id, code) WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- 6. permissions: rename column + constraints + indexes
ALTER TABLE permissions DROP CONSTRAINT IF EXISTS fk_permissions_application;
ALTER TABLE permissions DROP CONSTRAINT IF EXISTS uq_permissions_application_code;
DROP INDEX IF EXISTS idx_permissions_application_id;
DROP INDEX IF EXISTS idx_permissions_app_code_active;
ALTER TABLE permissions RENAME COLUMN application_id TO product_id;
ALTER TABLE permissions ADD CONSTRAINT fk_permissions_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE permissions ADD CONSTRAINT uq_permissions_product_code UNIQUE (product_id, code);
CREATE INDEX idx_permissions_product_id ON permissions(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_permissions_product_code_active ON permissions(product_id, code) WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- 7. product_registration_configs: rename column + constraints + indexes
ALTER TABLE product_registration_configs DROP CONSTRAINT IF EXISTS fk_prc_application;
ALTER TABLE product_registration_configs DROP CONSTRAINT IF EXISTS uq_prc_app_reg_type;
DROP INDEX IF EXISTS idx_prc_application_id;
ALTER TABLE product_registration_configs RENAME COLUMN application_id TO product_id;
ALTER TABLE product_registration_configs ADD CONSTRAINT fk_prc_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE product_registration_configs ADD CONSTRAINT uq_prc_product_reg_type UNIQUE (product_id, registration_type);
CREATE INDEX idx_prc_product_id ON product_registration_configs(product_id);

-- 8. user_tenant_registrations: rename column + constraints + indexes
ALTER TABLE user_tenant_registrations DROP CONSTRAINT IF EXISTS fk_utr_application;
DROP INDEX IF EXISTS uq_utr_user_tenant_type_app;
DROP INDEX IF EXISTS idx_utr_tenant_app_type;
DROP INDEX IF EXISTS idx_utr_application_id;
ALTER TABLE user_tenant_registrations RENAME COLUMN application_id TO product_id;
ALTER TABLE user_tenant_registrations ADD CONSTRAINT fk_utr_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT;
CREATE UNIQUE INDEX uq_utr_user_tenant_type_product ON user_tenant_registrations(user_id, tenant_id, registration_type, COALESCE(product_id, '00000000-0000-0000-0000-000000000000'::uuid)) WHERE deleted_at IS NULL;
CREATE INDEX idx_utr_tenant_product_type ON user_tenant_registrations(tenant_id, product_id, registration_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_utr_product_id ON user_tenant_registrations(product_id) WHERE deleted_at IS NULL;

-- 9. participants: rename column + indexes (no FK — cross-domain boundary)
DROP INDEX IF EXISTS idx_participants_tenant_application;
DROP INDEX IF EXISTS uq_participants_user_tenant_app;
ALTER TABLE participants RENAME COLUMN application_id TO product_id;
CREATE INDEX idx_participants_tenant_product ON participants(tenant_id, product_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX uq_participants_user_tenant_product ON participants(user_id, tenant_id, product_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;

-- Update comments
COMMENT ON TABLE products IS 'Products (formerly applications) that use IAM';
COMMENT ON COLUMN participants.product_id IS 'UUID of owning product. No FK - cross-domain boundary to Organization domain.';

COMMIT;
