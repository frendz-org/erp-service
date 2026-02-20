DROP INDEX IF EXISTS uq_masterdata_items_category_tenant_code;
DROP INDEX IF EXISTS idx_masterdata_items_default;
DROP INDEX IF EXISTS idx_masterdata_items_effective;
DROP INDEX IF EXISTS idx_masterdata_items_name_trgm;
DROP INDEX IF EXISTS idx_masterdata_items_parent_active;
DROP INDEX IF EXISTS idx_masterdata_items_category_code;
DROP INDEX IF EXISTS idx_masterdata_categories_code_active;
DROP INDEX IF EXISTS idx_masterdata_items_category_active;
DROP INDEX IF EXISTS idx_masterdata_items_tenant;
DROP INDEX IF EXISTS idx_masterdata_items_parent;
DROP INDEX IF EXISTS idx_masterdata_items_category;
DROP INDEX IF EXISTS idx_masterdata_categories_parent;
-- Note: pg_trgm extension is not dropped as it may be used elsewhere
