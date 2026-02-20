DO $$
DECLARE
    v_category_id UUID;
BEGIN
    -- 1. Create TENANT_TYPE category
    INSERT INTO masterdata_categories (
        code, name, description,
        parent_category_id, is_system, is_tenant_extensible,
        sort_order, status, metadata
    ) VALUES (
        'TENANT_TYPE',
        'Tenant Type',
        'Classification of tenant types',
        NULL,           -- flat category (no hierarchy)
        TRUE,           -- system category (cannot be modified by tenant admins)
        FALSE,          -- NOT tenant-extensible (platform-level only)
        100,            -- sort order
        'ACTIVE',
        '{}'::jsonb
    )
    RETURNING id INTO v_category_id;

    -- 2. Seed items
    INSERT INTO masterdata_items (
        category_id, tenant_id, parent_item_id,
        code, name, alt_name, description,
        sort_order, is_system, is_default, status, metadata
    ) VALUES
        -- Mitra Pendiri
        (v_category_id, NULL, NULL,
         'MITRA_PENDIRI', 'Mitra Pendiri', NULL, 'Founding partner tenant type',
         1, TRUE, FALSE, 'ACTIVE', '{}'::jsonb),
        -- Pension Fund
        (v_category_id, NULL, NULL,
         'PENSION_FUND', 'Pension Fund', 'Dana Pensiun', 'Pension fund tenant type',
         2, TRUE, FALSE, 'ACTIVE', '{}'::jsonb);

    RAISE NOTICE 'Seeded TENANT_TYPE category with 2 items';
END $$;
