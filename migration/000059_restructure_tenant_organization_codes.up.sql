-- ============================================================================
-- RESTRUCTURE: Tenant codes → TENANT_NNN, Organization codes → ORGANIZATION_NNN
-- Description: Standardizes tenant and organization masterdata codes.
--              Platform tenant (code='platform') is NOT modified.
--              ORGANIZATION category update is skipped if category doesn't exist.
-- ============================================================================

DO $$
DECLARE
    v_tenant_type_category_id UUID;
    v_tenant_category_id UUID;
    v_org_category_id UUID;
    v_mitra_pendiri_item_id UUID;
    v_pension_fund_item_id UUID;
BEGIN
    -- ========================================================================
    -- Step 0: Look up existing category/item IDs
    -- ========================================================================
    SELECT id INTO v_tenant_type_category_id
    FROM masterdata_categories
    WHERE code = 'TENANT_TYPE' AND deleted_at IS NULL;

    IF v_tenant_type_category_id IS NULL THEN
        RAISE EXCEPTION 'TENANT_TYPE category not found.';
    END IF;

    -- ORGANIZATION category is optional — may not exist in all environments
    SELECT id INTO v_org_category_id
    FROM masterdata_categories
    WHERE code = 'ORGANIZATION' AND deleted_at IS NULL;

    -- Look up TENANT_TYPE items by actual codes (TENANT_TYPE_001 / TENANT_TYPE_002)
    SELECT id INTO v_mitra_pendiri_item_id
    FROM masterdata_items
    WHERE category_id = v_tenant_type_category_id
      AND code = 'TENANT_TYPE_001'
      AND deleted_at IS NULL;

    SELECT id INTO v_pension_fund_item_id
    FROM masterdata_items
    WHERE category_id = v_tenant_type_category_id
      AND code = 'TENANT_TYPE_002'
      AND deleted_at IS NULL;

    IF v_mitra_pendiri_item_id IS NULL OR v_pension_fund_item_id IS NULL THEN
        RAISE EXCEPTION 'TENANT_TYPE items (TENANT_TYPE_001 / TENANT_TYPE_002) not found.';
    END IF;

    -- ========================================================================
    -- Step 1: Create TENANT masterdata category (idempotent)
    -- ========================================================================
    SELECT id INTO v_tenant_category_id
    FROM masterdata_categories
    WHERE code = 'TENANT' AND deleted_at IS NULL;

    IF v_tenant_category_id IS NULL THEN
        INSERT INTO masterdata_categories (
            code, name, description,
            parent_category_id,
            is_system, is_tenant_extensible,
            sort_order, status, metadata
        ) VALUES (
            'TENANT',
            'Tenant',
            'Registered tenants in the system',
            NULL,
            TRUE,
            FALSE,
            102,
            'ACTIVE',
            '{}'::jsonb
        )
        RETURNING id INTO v_tenant_category_id;
    END IF;

    -- ========================================================================
    -- Step 2: Seed TENANT masterdata items (idempotent, non-platform only)
    -- ========================================================================
    INSERT INTO masterdata_items (
        category_id, tenant_id, parent_item_id,
        code, name, alt_name, description,
        sort_order, is_system, is_default, status, metadata
    ) VALUES
        (v_tenant_category_id, NULL, v_mitra_pendiri_item_id,
         'TENANT_001', 'PT ISM Bogasari Flour Mills', 'ISM Bogasari',
         'Tenant for PT ISM Bogasari Flour Mills',
         1, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb),

        (v_tenant_category_id, NULL, v_mitra_pendiri_item_id,
         'TENANT_002', 'PT Inti Abadi Kemasindo', 'IAK',
         'Tenant for PT Inti Abadi Kemasindo',
         2, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb),

        (v_tenant_category_id, NULL, v_pension_fund_item_id,
         'TENANT_003', 'DPIP Bogasari', 'DPIP',
         'Tenant for DPIP Bogasari',
         3, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb),

        (v_tenant_category_id, NULL, v_pension_fund_item_id,
         'TENANT_004', 'DPMP Bogasari', 'DPMP',
         'Tenant for DPMP Bogasari',
         4, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb)
    ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
        WHERE deleted_at IS NULL
    DO NOTHING;

    -- ========================================================================
    -- Step 3: Update non-platform tenant codes to TENANT_NNN
    -- ========================================================================
    UPDATE tenants SET code = 'TENANT_001'
    WHERE code = 'ISM-BOGASARI' AND deleted_at IS NULL;

    UPDATE tenants SET code = 'TENANT_002'
    WHERE code = 'INTI-ABADI-KEMASINDO' AND deleted_at IS NULL;

    UPDATE tenants SET code = 'TENANT_003'
    WHERE code = 'DPIP-BOGASARI' AND deleted_at IS NULL;

    UPDATE tenants SET code = 'TENANT_004'
    WHERE code = 'DPMP-BOGASARI' AND deleted_at IS NULL;

    -- ========================================================================
    -- Step 4: Update ORGANIZATION masterdata item codes (if category exists)
    -- ========================================================================
    IF v_org_category_id IS NOT NULL THEN
        UPDATE masterdata_items SET code = 'ORGANIZATION_001'
        WHERE code = 'ISM-BOGASARI' AND category_id = v_org_category_id AND deleted_at IS NULL;

        UPDATE masterdata_items SET code = 'ORGANIZATION_002'
        WHERE code = 'INTI-ABADI-KEMASINDO' AND category_id = v_org_category_id AND deleted_at IS NULL;

        UPDATE masterdata_items SET code = 'ORGANIZATION_003'
        WHERE code = 'DPIP-BOGASARI' AND category_id = v_org_category_id AND deleted_at IS NULL;

        UPDATE masterdata_items SET code = 'ORGANIZATION_004'
        WHERE code = 'DPMP-BOGASARI' AND category_id = v_org_category_id AND deleted_at IS NULL;

        RAISE NOTICE 'Updated ORGANIZATION codes to ORGANIZATION_NNN';
    ELSE
        RAISE NOTICE 'ORGANIZATION category not found — skipping organization code update';
    END IF;

    RAISE NOTICE 'Restructured tenant codes to TENANT_NNN';
END $$;
