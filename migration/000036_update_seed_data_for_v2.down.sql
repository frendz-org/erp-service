DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_iam_app_id UUID;
    v_role_assign_perm_id UUID;
BEGIN
    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';

    IF v_platform_tenant_id IS NOT NULL THEN
        SELECT id INTO v_iam_app_id
        FROM applications
        WHERE tenant_id = v_platform_tenant_id AND code = 'iam-admin';
    END IF;

    DELETE FROM user_tenant_registrations;
    RAISE NOTICE 'Removed all user_tenant_registrations';

    IF v_iam_app_id IS NOT NULL THEN
        SELECT id INTO v_role_assign_perm_id
        FROM permissions
        WHERE application_id = v_iam_app_id AND code = 'role:assign';

        IF v_role_assign_perm_id IS NOT NULL THEN
            DELETE FROM role_permissions WHERE permission_id = v_role_assign_perm_id;
            RAISE NOTICE 'Removed role:assign from role_permissions';
        END IF;


        DELETE FROM permissions
        WHERE application_id = v_iam_app_id AND code = 'role:assign';
        RAISE NOTICE 'Removed role:assign permission';
    END IF;

    RAISE NOTICE '=== V2 seed data rollback complete ===';
END $$;
