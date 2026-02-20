-- ============================================================================
-- ROLLBACK: Restore original tenant and organization codes
-- ============================================================================

DO $$
DECLARE
    v_org_category_id UUID;
    v_tenant_category_id UUID;
BEGIN
    -- ========================================================================
    -- Step 1: Restore ORGANIZATION item codes
    -- ========================================================================
    SELECT id INTO v_org_category_id
    FROM masterdata_categories
    WHERE code = 'ORGANIZATION' AND deleted_at IS NULL;

    IF v_org_category_id IS NOT NULL THEN
        UPDATE masterdata_items SET code = 'ISM-BOGASARI'
        WHERE code = 'ORGANIZATION_001' AND category_id = v_org_category_id AND deleted_at IS NULL;

        UPDATE masterdata_items SET code = 'INTI-ABADI-KEMASINDO'
        WHERE code = 'ORGANIZATION_002' AND category_id = v_org_category_id AND deleted_at IS NULL;

        UPDATE masterdata_items SET code = 'DPIP-BOGASARI'
        WHERE code = 'ORGANIZATION_003' AND category_id = v_org_category_id AND deleted_at IS NULL;

        UPDATE masterdata_items SET code = 'DPMP-BOGASARI'
        WHERE code = 'ORGANIZATION_004' AND category_id = v_org_category_id AND deleted_at IS NULL;
    END IF;

    -- ========================================================================
    -- Step 2: Restore tenant codes
    -- ========================================================================
    UPDATE tenants SET code = 'ISM-BOGASARI'
    WHERE code = 'TENANT_001' AND deleted_at IS NULL;

    UPDATE tenants SET code = 'INTI-ABADI-KEMASINDO'
    WHERE code = 'TENANT_002' AND deleted_at IS NULL;

    UPDATE tenants SET code = 'DPIP-BOGASARI'
    WHERE code = 'TENANT_003' AND deleted_at IS NULL;

    UPDATE tenants SET code = 'DPMP-BOGASARI'
    WHERE code = 'TENANT_004' AND deleted_at IS NULL;

    -- ========================================================================
    -- Step 3: Delete TENANT masterdata items and category (children first)
    -- ========================================================================
    SELECT id INTO v_tenant_category_id
    FROM masterdata_categories
    WHERE code = 'TENANT' AND deleted_at IS NULL;

    IF v_tenant_category_id IS NOT NULL THEN
        DELETE FROM masterdata_items
        WHERE category_id = v_tenant_category_id;

        DELETE FROM masterdata_categories
        WHERE id = v_tenant_category_id;
    END IF;

    RAISE NOTICE 'Rolled back tenant and organization code restructuring';
END $$;
