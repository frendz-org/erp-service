-- ============================================================================
-- RENAME: DIVISION category/items → ORGANIZATION
-- Description: Renames the DIVISION masterdata category to ORGANIZATION and
--              updates item codes from DIVISION_NNN to ORGANIZATION_NNN.
-- ============================================================================

DO $$
DECLARE
    v_division_category_id UUID;
BEGIN
    SELECT id INTO v_division_category_id
    FROM masterdata_categories
    WHERE code = 'DIVISION' AND deleted_at IS NULL;

    IF v_division_category_id IS NULL THEN
        RAISE EXCEPTION 'DIVISION category not found.';
    END IF;

    -- Step 1: Rename category
    UPDATE masterdata_categories
    SET code = 'ORGANIZATION', name = 'Organization',
        description = 'Organizations/companies within the system'
    WHERE id = v_division_category_id;

    -- Step 2: Rename item codes
    UPDATE masterdata_items
    SET code = REPLACE(code, 'DIVISION_', 'ORGANIZATION_')
    WHERE category_id = v_division_category_id
      AND code LIKE 'DIVISION_%'
      AND deleted_at IS NULL;

    RAISE NOTICE 'Renamed DIVISION category and items to ORGANIZATION';
END $$;
