BEGIN;

-- 9. participants: reverse rename column + indexes
DROP INDEX IF EXISTS uq_participants_user_tenant_product;
DROP INDEX IF EXISTS idx_participants_tenant_product;
ALTER TABLE participants RENAME COLUMN product_id TO application_id;
CREATE INDEX idx_participants_tenant_application ON participants(tenant_id, application_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX uq_participants_user_tenant_app ON participants(user_id, tenant_id, application_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL;

-- Drop FK constraints on dependent tables that reference products before renaming
ALTER TABLE user_tenant_registrations DROP CONSTRAINT IF EXISTS fk_utr_product;
ALTER TABLE product_registration_configs DROP CONSTRAINT IF EXISTS fk_prc_product;
ALTER TABLE permissions DROP CONSTRAINT IF EXISTS fk_permissions_product;
ALTER TABLE roles DROP CONSTRAINT IF EXISTS fk_roles_product;

-- Reverse indexes on products table (before rename)
DROP INDEX IF EXISTS idx_products_tenant_code_active;
DROP INDEX IF EXISTS idx_products_tenant_id;

-- 4. Reverse FK and unique constraint renames on products
ALTER TABLE products RENAME CONSTRAINT uq_products_tenant_code TO uq_applications_tenant_code;
ALTER TABLE products RENAME CONSTRAINT fk_products_tenant TO fk_applications_tenant;

-- 3. Reverse CHECK constraint rename
ALTER TABLE products RENAME CONSTRAINT chk_products_status TO chk_applications_status;

-- 2. Reverse trigger rename
ALTER TRIGGER trg_products_updated_at ON products RENAME TO trg_applications_updated_at;

-- 1. Reverse table rename (must happen before re-adding FKs that reference applications)
ALTER TABLE products RENAME TO applications;

-- Repair user_role_assignments FK after table rename
ALTER TABLE user_role_assignments DROP CONSTRAINT IF EXISTS fk_ura_product;
ALTER TABLE user_role_assignments
    ADD CONSTRAINT fk_ura_product FOREIGN KEY (product_id)
        REFERENCES applications(id) ON DELETE CASCADE;

-- Recreate indexes on the now-renamed applications table
CREATE INDEX idx_applications_tenant_id ON applications(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_applications_tenant_code_active ON applications(tenant_id, code) WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- 8. user_tenant_registrations: reverse rename column + constraints + indexes
DROP INDEX IF EXISTS idx_utr_product_id;
DROP INDEX IF EXISTS idx_utr_tenant_product_type;
DROP INDEX IF EXISTS uq_utr_user_tenant_type_product;
ALTER TABLE user_tenant_registrations RENAME COLUMN product_id TO application_id;
ALTER TABLE user_tenant_registrations ADD CONSTRAINT fk_utr_application FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE RESTRICT;
CREATE UNIQUE INDEX uq_utr_user_tenant_type_app ON user_tenant_registrations(user_id, tenant_id, registration_type, COALESCE(application_id, '00000000-0000-0000-0000-000000000000'::uuid)) WHERE deleted_at IS NULL;
CREATE INDEX idx_utr_tenant_app_type ON user_tenant_registrations(tenant_id, application_id, registration_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_utr_application_id ON user_tenant_registrations(application_id) WHERE deleted_at IS NULL;

-- 7. product_registration_configs: reverse rename column + constraints + indexes
DROP INDEX IF EXISTS idx_prc_product_id;
ALTER TABLE product_registration_configs DROP CONSTRAINT IF EXISTS uq_prc_product_reg_type;
ALTER TABLE product_registration_configs RENAME COLUMN product_id TO application_id;
ALTER TABLE product_registration_configs ADD CONSTRAINT fk_prc_application FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE;
ALTER TABLE product_registration_configs ADD CONSTRAINT uq_prc_app_reg_type UNIQUE (application_id, registration_type);
CREATE INDEX idx_prc_application_id ON product_registration_configs(application_id);

-- 6. permissions: reverse rename column + constraints + indexes
DROP INDEX IF EXISTS idx_permissions_product_code_active;
DROP INDEX IF EXISTS idx_permissions_product_id;
ALTER TABLE permissions DROP CONSTRAINT IF EXISTS uq_permissions_product_code;
ALTER TABLE permissions RENAME COLUMN product_id TO application_id;
ALTER TABLE permissions ADD CONSTRAINT fk_permissions_application FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE;
ALTER TABLE permissions ADD CONSTRAINT uq_permissions_application_code UNIQUE (application_id, code);
CREATE INDEX idx_permissions_application_id ON permissions(application_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_permissions_app_code_active ON permissions(application_id, code) WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- 5. roles: reverse rename column + constraints + indexes
DROP INDEX IF EXISTS idx_roles_product_code_active;
DROP INDEX IF EXISTS idx_roles_product_id;
ALTER TABLE roles DROP CONSTRAINT IF EXISTS uq_roles_product_code;
ALTER TABLE roles RENAME COLUMN product_id TO application_id;
ALTER TABLE roles ADD CONSTRAINT fk_roles_application FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE;
ALTER TABLE roles ADD CONSTRAINT uq_roles_application_code UNIQUE (application_id, code);
CREATE INDEX idx_roles_application_id ON roles(application_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_app_code_active ON roles(application_id, code) WHERE deleted_at IS NULL AND status = 'ACTIVE';

COMMIT;
