DO $$
BEGIN
    UPDATE masterdata_items mi
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = new_code
    FROM (
        SELECT id,
               'GENDER_' || LPAD(ROW_NUMBER() OVER (ORDER BY sort_order)::text, 3, '0') as new_code
        FROM masterdata_items
        WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'GENDER')
          AND deleted_at IS NULL
    ) AS mapping
    WHERE mi.id = mapping.id;

    RAISE NOTICE 'Updated GENDER codes';

    UPDATE masterdata_items mi
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = new_code
    FROM (
        SELECT id,
               'MARITAL_STATUS_' || LPAD(ROW_NUMBER() OVER (ORDER BY sort_order)::text, 3, '0') as new_code
        FROM masterdata_items
        WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'MARITAL_STATUS')
          AND deleted_at IS NULL
    ) AS mapping
    WHERE mi.id = mapping.id;

    RAISE NOTICE 'Updated MARITAL_STATUS codes';

    UPDATE masterdata_items mi
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = new_code
    FROM (
        SELECT id,
               'RELIGION_' || LPAD(ROW_NUMBER() OVER (ORDER BY sort_order)::text, 3, '0') as new_code
        FROM masterdata_items
        WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'RELIGION')
          AND deleted_at IS NULL
    ) AS mapping
    WHERE mi.id = mapping.id;

    RAISE NOTICE 'Updated RELIGION codes';

    UPDATE masterdata_items mi
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = new_code
    FROM (
        SELECT id,
               'EDUCATION_LEVEL_' || LPAD(ROW_NUMBER() OVER (ORDER BY sort_order)::text, 3, '0') as new_code
        FROM masterdata_items
        WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'EDUCATION_LEVEL')
          AND deleted_at IS NULL
    ) AS mapping
    WHERE mi.id = mapping.id;

    RAISE NOTICE 'Updated EDUCATION_LEVEL codes';

    UPDATE masterdata_items mi
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = new_code
    FROM (
        SELECT id,
               'BLOOD_TYPE_' || LPAD(ROW_NUMBER() OVER (ORDER BY sort_order)::text, 3, '0') as new_code
        FROM masterdata_items
        WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'BLOOD_TYPE')
          AND deleted_at IS NULL
    ) AS mapping
    WHERE mi.id = mapping.id;

    RAISE NOTICE 'Updated BLOOD_TYPE codes';

    UPDATE masterdata_items mi
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = new_code
    FROM (
        SELECT id,
               'IDENTITY_TYPE_' || LPAD(ROW_NUMBER() OVER (ORDER BY sort_order)::text, 3, '0') as new_code
        FROM masterdata_items
        WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'IDENTITY_TYPE')
          AND deleted_at IS NULL
    ) AS mapping
    WHERE mi.id = mapping.id;

    RAISE NOTICE 'Updated IDENTITY_TYPE codes';

    UPDATE masterdata_items mi
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = new_code
    FROM (
        SELECT id,
               'TENANT_TYPE_' || LPAD(ROW_NUMBER() OVER (ORDER BY sort_order)::text, 3, '0') as new_code
        FROM masterdata_items
        WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'TENANT_TYPE')
          AND deleted_at IS NULL
    ) AS mapping
    WHERE mi.id = mapping.id;

    RAISE NOTICE 'Updated TENANT_TYPE codes';

    -- ========================================================================
    -- ORGANIZATION → TENANT (rename category and items)
    -- ========================================================================

    -- First, rename the category ORGANIZATION → TENANT
    UPDATE masterdata_categories
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = 'TENANT',
        name = 'Tenant',
        description = 'Tenants/organizations registered in the system'
    WHERE code = 'ORGANIZATION'
      AND deleted_at IS NULL;

    RAISE NOTICE 'Renamed ORGANIZATION category to TENANT';

    -- Then, update the items under the renamed TENANT category
    UPDATE masterdata_items mi
    SET
        metadata = jsonb_set(
            COALESCE(metadata, '{}'::jsonb),
            '{legacy_code}',
            to_jsonb(code)
        ),
        code = new_code
    FROM (
        SELECT id,
               'TENANT_' || LPAD(ROW_NUMBER() OVER (ORDER BY sort_order)::text, 3, '0') as new_code
        FROM masterdata_items
        WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'TENANT')
          AND deleted_at IS NULL
    ) AS mapping
    WHERE mi.id = mapping.id;

    RAISE NOTICE 'Updated TENANT item codes';


    UPDATE tenants t
    SET tenant_type = mi.code
    FROM masterdata_items mi
    JOIN masterdata_categories mc ON mi.category_id = mc.id
    WHERE mc.code = 'TENANT_TYPE'
      AND mi.metadata->>'legacy_code' = t.tenant_type
      AND mi.deleted_at IS NULL;

    RAISE NOTICE 'Updated tenants.tenant_type references';

    RAISE NOTICE '';
    RAISE NOTICE '=== MASTERDATA CODE NORMALIZATION COMPLETE ===';
    RAISE NOTICE 'All codes now follow format: {CATEGORY_CODE}_{SEQUENCE_NUMBER}';
    RAISE NOTICE 'Original codes preserved in metadata.legacy_code';
    RAISE NOTICE 'Referencing tables (tenants) updated to new codes';
    RAISE NOTICE '';

END $$;
