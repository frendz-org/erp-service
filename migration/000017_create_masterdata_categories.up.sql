CREATE TABLE masterdata_categories (
    -- Primary Key
    id                      UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Business Identifiers
    code                    VARCHAR(50) NOT NULL,
    name                    VARCHAR(255) NOT NULL,
    description             TEXT,

    -- Hierarchy Definition (self-referencing)
    parent_category_id      UUID,

    -- Category Behavior Flags
    is_system               BOOLEAN NOT NULL DEFAULT FALSE,
    is_tenant_extensible    BOOLEAN NOT NULL DEFAULT FALSE,

    -- Display
    sort_order              INTEGER NOT NULL DEFAULT 0,

    -- Status Management
    status                  VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    -- Flexible Configuration
    metadata                JSONB NOT NULL DEFAULT '{}',

    -- Audit Fields
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ,

    -- Optimistic Locking
    version                 INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT uq_masterdata_categories_code UNIQUE (code),
    CONSTRAINT fk_masterdata_categories_parent FOREIGN KEY (parent_category_id)
        REFERENCES masterdata_categories(id) ON DELETE RESTRICT,
    CONSTRAINT chk_masterdata_categories_status CHECK (status IN ('ACTIVE', 'INACTIVE')),
    CONSTRAINT chk_masterdata_categories_no_self_ref CHECK (id != parent_category_id)
);

CREATE TRIGGER trg_masterdata_categories_updated_at
    BEFORE UPDATE ON masterdata_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE masterdata_categories IS 'Defines types of reference data. Self-referencing parent creates hierarchy rules.';
COMMENT ON COLUMN masterdata_categories.code IS 'Unique machine-readable identifier: GENDER, COUNTRY, PROVINCE, CITY';
COMMENT ON COLUMN masterdata_categories.parent_category_id IS 'NULL = flat/root category. Set = child category (CITY parent is PROVINCE).';
COMMENT ON COLUMN masterdata_categories.is_system IS 'System categories (GENDER, COUNTRY, etc.) cannot be modified or deleted by tenant admins.';
COMMENT ON COLUMN masterdata_categories.is_tenant_extensible IS 'If TRUE, tenant admins can add their own items to this category (e.g., DEPARTMENT). If FALSE, only platform admins can manage items (e.g., COUNTRY).';
COMMENT ON COLUMN masterdata_categories.metadata IS 'Category-level config. Example: {"max_code_length": 10, "code_pattern": "^[A-Z_]+$"}';
