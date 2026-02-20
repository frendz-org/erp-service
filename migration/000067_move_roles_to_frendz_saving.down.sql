DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_iam_product_id UUID;
    v_frendz_product_id UUID;
BEGIN
    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';
    IF v_platform_tenant_id IS NULL THEN RAISE NOTICE 'Platform tenant not found, skipping'; RETURN; END IF;

    SELECT id INTO v_iam_product_id FROM products WHERE tenant_id = v_platform_tenant_id AND code = 'iam-admin';
    SELECT id INTO v_frendz_product_id FROM products WHERE tenant_id = v_platform_tenant_id AND code = 'frendz-saving';

    IF v_iam_product_id IS NULL OR v_frendz_product_id IS NULL THEN
        RAISE NOTICE 'iam-admin or frendz-saving product not found, skipping'; RETURN;
    END IF;

    UPDATE permissions SET product_id = v_iam_product_id, updated_at = NOW()
    WHERE product_id = v_frendz_product_id AND (code LIKE 'participant:%' OR code LIKE 'member:%') AND deleted_at IS NULL;

    UPDATE roles SET product_id = v_iam_product_id, updated_at = NOW()
    WHERE product_id = v_frendz_product_id AND code IN ('PARTICIPANT_CREATOR', 'PARTICIPANT_APPROVER', 'TENANT_PRODUCT_ADMIN') AND deleted_at IS NULL;

    RAISE NOTICE 'Moved participant/member roles and permissions back to iam-admin';
END $$;
