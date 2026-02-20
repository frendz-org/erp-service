-- ============================================================================
-- UNIQUE INDEXES: Business Rules (using expressions)
-- ============================================================================

-- Unique code per category+tenant (treats NULL tenant as global)
CREATE UNIQUE INDEX uq_masterdata_items_category_tenant_code
    ON masterdata_items(
        category_id,
        COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid),
        code
    )
    WHERE deleted_at IS NULL;

-- ============================================================================
-- INDEXES: Foreign Key Columns
-- ============================================================================

-- Categories: parent lookup
CREATE INDEX idx_masterdata_categories_parent
    ON masterdata_categories(parent_category_id)
    WHERE deleted_at IS NULL AND parent_category_id IS NOT NULL;

-- Items: category lookup (most common query pattern)
CREATE INDEX idx_masterdata_items_category
    ON masterdata_items(category_id)
    WHERE deleted_at IS NULL;

-- Items: parent item lookup (for hierarchy traversal)
CREATE INDEX idx_masterdata_items_parent
    ON masterdata_items(parent_item_id)
    WHERE deleted_at IS NULL AND parent_item_id IS NOT NULL;

-- Items: tenant scoping
CREATE INDEX idx_masterdata_items_tenant
    ON masterdata_items(tenant_id)
    WHERE deleted_at IS NULL AND tenant_id IS NOT NULL;

-- ============================================================================
-- INDEXES: Query Optimization
-- ============================================================================

-- Primary query: Get all active items in a category (for dropdowns)
CREATE INDEX idx_masterdata_items_category_active
    ON masterdata_items(category_id, status, tenant_id)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Category lookup by code (API endpoint: GET /categories/{code})
CREATE INDEX idx_masterdata_categories_code_active
    ON masterdata_categories(code)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Item lookup by code within category (validation: "is 'MALE' a valid GENDER?")
CREATE INDEX idx_masterdata_items_category_code
    ON masterdata_items(category_id, code)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Hierarchical query: Get children of a parent item (e.g., all cities in DKI Jakarta)
CREATE INDEX idx_masterdata_items_parent_active
    ON masterdata_items(parent_item_id, category_id, status)
    WHERE deleted_at IS NULL AND status = 'ACTIVE' AND parent_item_id IS NOT NULL;

-- Search by name (for autocomplete/search features)
-- Note: Requires pg_trgm extension for fuzzy matching
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_masterdata_items_name_trgm
    ON masterdata_items USING gin (LOWER(name) gin_trgm_ops)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Temporal validity: find currently effective items
CREATE INDEX idx_masterdata_items_effective
    ON masterdata_items(effective_from, effective_until)
    WHERE deleted_at IS NULL AND status = 'ACTIVE'
    AND (effective_from IS NOT NULL OR effective_until IS NOT NULL);

-- Default item per category (for UI pre-selection)
CREATE INDEX idx_masterdata_items_default
    ON masterdata_items(category_id, tenant_id)
    WHERE is_default = TRUE AND deleted_at IS NULL AND status = 'ACTIVE';
