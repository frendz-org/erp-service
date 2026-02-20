-- ============================================================================
-- ROLLBACK: Restore ORGANIZATION → DIVISION
-- ============================================================================

DO $$
DECLARE
    v_org_category_id UUID;
BEGIN
    SELECT id INTO v_org_category_id
    FROM masterdata_categories
    WHERE code = 'ORGANIZATION' AND deleted_at IS NULL;

    IF v_org_category_id IS NULL THEN
        RAISE EXCEPTION 'ORGANIZATION category not found.';
    END IF;

    -- Step 1: Restore item codes
    UPDATE masterdata_items
    SET code = REPLACE(code, 'ORGANIZATION_', 'DIVISION_')
    WHERE category_id = v_org_category_id
      AND code LIKE 'ORGANIZATION_%'
      AND deleted_at IS NULL;

    -- Step 2: Restore category
    UPDATE masterdata_categories
    SET code = 'DIVISION', name = 'Division',
        description = 'Divisions within each legal entity'
    WHERE id = v_org_category_id;

    RAISE NOTICE 'Restored ORGANIZATION category and items back to DIVISION';
END $$;
