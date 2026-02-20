CREATE TABLE masterdata_items (
    -- Primary Key
    id                      UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    category_id             UUID NOT NULL,

    -- Tenant Scoping
    tenant_id               UUID,

    -- Hierarchy (self-referencing)
    parent_item_id          UUID,

    -- Business Identifiers
    code                    VARCHAR(100) NOT NULL,
    name                    VARCHAR(255) NOT NULL,
    alt_name                VARCHAR(255),
    description             TEXT,

    -- Display & Behavior
    sort_order              INTEGER NOT NULL DEFAULT 0,
    is_system               BOOLEAN NOT NULL DEFAULT FALSE,
    is_default              BOOLEAN NOT NULL DEFAULT FALSE,

    -- Status Management
    status                  VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Temporal Validity (optional)
    effective_from          DATE,
    effective_until         DATE,

    -- Flexible Metadata
    metadata                JSONB NOT NULL DEFAULT '{}',

    -- Audit Fields
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ,

    -- Optimistic Locking
    version                 INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_masterdata_items_category FOREIGN KEY (category_id)
        REFERENCES masterdata_categories(id) ON DELETE RESTRICT,
    CONSTRAINT fk_masterdata_items_parent FOREIGN KEY (parent_item_id)
        REFERENCES masterdata_items(id) ON DELETE RESTRICT,
    CONSTRAINT chk_masterdata_items_status CHECK (status IN ('ACTIVE', 'INACTIVE')),
    CONSTRAINT chk_masterdata_items_no_self_ref CHECK (id != parent_item_id),
    CONSTRAINT chk_masterdata_items_effective_dates CHECK (
        effective_from IS NULL OR effective_until IS NULL OR effective_from <= effective_until
    )
);

CREATE TRIGGER trg_masterdata_items_updated_at
    BEFORE UPDATE ON masterdata_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE masterdata_items IS 'Reference data values. Supports hierarchy (city->province->country) and tenant-scoping.';
COMMENT ON COLUMN masterdata_items.category_id IS 'Which category this item belongs to (GENDER, COUNTRY, PROVINCE, CITY, etc.)';
COMMENT ON COLUMN masterdata_items.tenant_id IS 'NULL = global item (all tenants). Set = visible only to this tenant.';
COMMENT ON COLUMN masterdata_items.parent_item_id IS 'Hierarchical parent. Application validates parent belongs to parent category.';
COMMENT ON COLUMN masterdata_items.code IS 'Unique within category+tenant. This is what consuming services store (e.g., IAM stores gender=MALE).';
COMMENT ON COLUMN masterdata_items.alt_name IS 'Alternative name for display (e.g., local language, abbreviation). Jakarta Selatan alt_name=Jaksel.';
COMMENT ON COLUMN masterdata_items.is_default IS 'Suggested default for UI dropdowns. Only one item per category+tenant should be default.';
COMMENT ON COLUMN masterdata_items.effective_from IS 'For time-bound validity. NULL means valid since the beginning of time.';
COMMENT ON COLUMN masterdata_items.effective_until IS 'For time-bound validity. NULL means no expiration.';
COMMENT ON COLUMN masterdata_items.metadata IS 'Category-specific extra data (e.g., ISO codes for countries).';
