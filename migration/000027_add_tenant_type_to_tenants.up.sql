-- Add tenant_type column to tenants table
-- Stores masterdata code string (e.g., 'MITRA_PENDIRI', 'PENSION_FUND')
-- Validated against Masterdata service TENANT_TYPE category at write time

ALTER TABLE tenants
ADD COLUMN tenant_type VARCHAR(50);

COMMENT ON COLUMN tenants.tenant_type IS
    'Masterdata code string from TENANT_TYPE category. Validated against Masterdata service at write time.';
