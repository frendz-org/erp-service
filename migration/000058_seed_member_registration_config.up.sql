DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_frendz_app_id UUID;
BEGIN
    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';

    IF v_platform_tenant_id IS NULL THEN
        RAISE NOTICE 'Platform tenant not found, skipping member registration config seed';
        RETURN;
    END IF;

    SELECT id INTO v_frendz_app_id
    FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'frendz-saving';

    IF v_frendz_app_id IS NULL THEN
        RAISE NOTICE 'frendz-saving application not found, skipping member registration config seed';
        RETURN;
    END IF;

    -- Seed MEMBER registration config for frendz-saving (requires approval)
    INSERT INTO product_registration_configs (application_id, registration_type, requires_approval, is_active)
    VALUES (v_frendz_app_id, 'MEMBER', TRUE, TRUE)
    ON CONFLICT (application_id, registration_type) DO NOTHING;

    RAISE NOTICE 'Ensured MEMBER registration config for frendz-saving';
END $$;
