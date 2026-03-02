DO $$
DECLARE
    v_tenant_001_id UUID;
    v_tenant_003_id UUID;
    v_product_001_id UUID;
    v_product_003_id UUID;
BEGIN

    SELECT id INTO v_tenant_001_id FROM tenants WHERE code = 'TENANT_001' AND deleted_at IS NULL;
    SELECT id INTO v_tenant_003_id FROM tenants WHERE code = 'TENANT_003' AND deleted_at IS NULL;

    IF v_tenant_001_id IS NULL THEN
        RAISE NOTICE 'Tenant TENANT_001 not found, skipping';
        RETURN;
    END IF;

    IF v_tenant_003_id IS NULL THEN
        RAISE NOTICE 'Tenant TENANT_003 not found, skipping';
        RETURN;
    END IF;

    -- 1. Seed frendz-saving product for TENANT_001
    INSERT INTO products (tenant_id, code, name, description, status)
    VALUES (
        v_tenant_001_id,
        'frendz-saving',
        'Frendz Saving',
        'Frendz Saving program for participant management',
        'ACTIVE'
    )
    ON CONFLICT (tenant_id, code) DO NOTHING;

    SELECT id INTO v_product_001_id
    FROM products WHERE tenant_id = v_tenant_001_id AND code = 'frendz-saving';

    RAISE NOTICE 'Ensured frendz-saving product for TENANT_001 with ID: %', v_product_001_id;

    -- 2. Seed frendz-saving product for TENANT_003
    INSERT INTO products (tenant_id, code, name, description, status)
    VALUES (
        v_tenant_003_id,
        'frendz-saving',
        'Frendz Saving',
        'Frendz Saving program for participant management',
        'ACTIVE'
    )
    ON CONFLICT (tenant_id, code) DO NOTHING;

    SELECT id INTO v_product_003_id
    FROM products WHERE tenant_id = v_tenant_003_id AND code = 'frendz-saving';

    RAISE NOTICE 'Ensured frendz-saving product for TENANT_003 with ID: %', v_product_003_id;

    -- 3. Seed PARTICIPANT registration config for TENANT_001 frendz-saving
    INSERT INTO product_registration_configs (product_id, registration_type, requires_approval, is_active)
    VALUES (v_product_001_id, 'PARTICIPANT', TRUE, TRUE)
    ON CONFLICT (product_id, registration_type) DO NOTHING;

    RAISE NOTICE 'Ensured PARTICIPANT registration config for TENANT_001 frendz-saving';

    -- 4. Seed MEMBER registration config for TENANT_001 frendz-saving
    INSERT INTO product_registration_configs (product_id, registration_type, requires_approval, is_active)
    VALUES (v_product_001_id, 'MEMBER', TRUE, TRUE)
    ON CONFLICT (product_id, registration_type) DO NOTHING;

    RAISE NOTICE 'Ensured MEMBER registration config for TENANT_001 frendz-saving';

    -- 5. Seed PARTICIPANT registration config for TENANT_003 frendz-saving
    INSERT INTO product_registration_configs (product_id, registration_type, requires_approval, is_active)
    VALUES (v_product_003_id, 'PARTICIPANT', TRUE, TRUE)
    ON CONFLICT (product_id, registration_type) DO NOTHING;

    RAISE NOTICE 'Ensured PARTICIPANT registration config for TENANT_003 frendz-saving';

    -- 6. Seed MEMBER registration config for TENANT_003 frendz-saving
    INSERT INTO product_registration_configs (product_id, registration_type, requires_approval, is_active)
    VALUES (v_product_003_id, 'MEMBER', TRUE, TRUE)
    ON CONFLICT (product_id, registration_type) DO NOTHING;

    RAISE NOTICE 'Ensured MEMBER registration config for TENANT_003 frendz-saving';

    RAISE NOTICE '=== Tenant frendz-saving products and registration configs seeded successfully ===';
END $$;
