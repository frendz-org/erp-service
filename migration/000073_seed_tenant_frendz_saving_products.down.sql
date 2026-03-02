DO $$
DECLARE
    v_tenant_001_id UUID;
    v_tenant_003_id UUID;
    v_product_001_id UUID;
    v_product_003_id UUID;
BEGIN

    SELECT id INTO v_tenant_001_id FROM tenants WHERE code = 'TENANT_001' AND deleted_at IS NULL;
    SELECT id INTO v_tenant_003_id FROM tenants WHERE code = 'TENANT_003' AND deleted_at IS NULL;

    -- 1. Remove registration configs and product for TENANT_001
    IF v_tenant_001_id IS NOT NULL THEN
        SELECT id INTO v_product_001_id
        FROM products WHERE tenant_id = v_tenant_001_id AND code = 'frendz-saving';

        IF v_product_001_id IS NOT NULL THEN
            DELETE FROM product_registration_configs
            WHERE product_id = v_product_001_id;

            RAISE NOTICE 'Removed registration configs for TENANT_001 frendz-saving';

            DELETE FROM products WHERE id = v_product_001_id;

            RAISE NOTICE 'Removed frendz-saving product for TENANT_001';
        END IF;
    END IF;

    -- 2. Remove registration configs and product for TENANT_003
    IF v_tenant_003_id IS NOT NULL THEN
        SELECT id INTO v_product_003_id
        FROM products WHERE tenant_id = v_tenant_003_id AND code = 'frendz-saving';

        IF v_product_003_id IS NOT NULL THEN
            DELETE FROM product_registration_configs
            WHERE product_id = v_product_003_id;

            RAISE NOTICE 'Removed registration configs for TENANT_003 frendz-saving';

            DELETE FROM products WHERE id = v_product_003_id;

            RAISE NOTICE 'Removed frendz-saving product for TENANT_003';
        END IF;
    END IF;

    RAISE NOTICE '=== Tenant frendz-saving products and registration configs removed ===';
END $$;
