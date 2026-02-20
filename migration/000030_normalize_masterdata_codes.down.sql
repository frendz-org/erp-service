DO $$
BEGIN
    UPDATE tenants t
    SET tenant_type = mi.metadata->>'legacy_code'
    FROM masterdata_items mi
    JOIN masterdata_categories mc ON mi.category_id = mc.id
    WHERE mc.code = 'TENANT_TYPE'
      AND mi.code = t.tenant_type
      AND mi.metadata->>'legacy_code' IS NOT NULL
      AND mi.deleted_at IS NULL;

    RAISE NOTICE 'Restored tenants.tenant_type to original codes';

    UPDATE masterdata_items
    SET code = metadata->>'legacy_code'
    WHERE metadata->>'legacy_code' IS NOT NULL
      AND deleted_at IS NULL;

    UPDATE masterdata_items
    SET metadata = metadata - 'legacy_code'
    WHERE metadata ? 'legacy_code'
      AND deleted_at IS NULL;

    RAISE NOTICE 'Restored original masterdata item codes';

    UPDATE masterdata_categories
    SET
        code = metadata->>'legacy_code',
        name = 'Organization',
        description = 'Organizations/companies within each tenant type'
    WHERE code = 'TENANT'
      AND metadata->>'legacy_code' = 'ORGANIZATION'
      AND deleted_at IS NULL;

    UPDATE masterdata_categories
    SET metadata = metadata - 'legacy_code'
    WHERE metadata ? 'legacy_code'
      AND deleted_at IS NULL;

    RAISE NOTICE 'Restored TENANT category back to ORGANIZATION';

END $$;
