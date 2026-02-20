DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_frendz_app_id UUID;
BEGIN
    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';
    IF v_platform_tenant_id IS NULL THEN
        RETURN;
    END IF;

    SELECT id INTO v_frendz_app_id FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'frendz-saving';
    IF v_frendz_app_id IS NULL THEN
        RETURN;
    END IF;

    DELETE FROM product_registration_configs
    WHERE application_id = v_frendz_app_id AND registration_type = 'PARTICIPANT';
END $$;
