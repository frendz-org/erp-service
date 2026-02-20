DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_iam_app_id UUID;
BEGIN
    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';

    IF v_platform_tenant_id IS NULL THEN
        RAISE NOTICE 'Platform tenant not found, nothing to delete';
        RETURN;
    END IF;

    SELECT id INTO v_iam_app_id
    FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'iam-admin';

    -- 1. Remove member permissions from PLATFORM_ADMIN role
    IF v_iam_app_id IS NOT NULL THEN
        DELETE FROM role_permissions
        WHERE permission_id IN (
            SELECT id FROM permissions
            WHERE application_id = v_iam_app_id AND code LIKE 'member:%'
        );

        RAISE NOTICE 'Removed member permission assignments';
    END IF;

    -- 2. Remove TENANT_PRODUCT_ADMIN role (role_permissions cascade-deleted by FK)
    IF v_iam_app_id IS NOT NULL THEN
        DELETE FROM role_permissions
        WHERE role_id IN (
            SELECT id FROM roles
            WHERE application_id = v_iam_app_id AND code = 'TENANT_PRODUCT_ADMIN'
        );

        DELETE FROM roles
        WHERE application_id = v_iam_app_id
          AND code = 'TENANT_PRODUCT_ADMIN';

        RAISE NOTICE 'Removed TENANT_PRODUCT_ADMIN role';
    END IF;

    -- 3. Remove member permissions
    IF v_iam_app_id IS NOT NULL THEN
        DELETE FROM permissions
        WHERE application_id = v_iam_app_id
          AND code LIKE 'member:%';

        RAISE NOTICE 'Removed member permissions';
    END IF;

    RAISE NOTICE '=== Member roles and permissions removed ===';
END $$;
