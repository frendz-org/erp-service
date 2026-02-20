DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_admin_user_id UUID;
    v_iam_app_id UUID;
    v_platform_admin_role_id UUID;
    v_role_assign_perm_id UUID;
BEGIN
    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';
    SELECT id INTO v_admin_user_id FROM users WHERE email = 'admin@platform.local';
    SELECT id INTO v_iam_app_id
    FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'iam-admin';
    SELECT id INTO v_platform_admin_role_id
    FROM roles
    WHERE application_id = v_iam_app_id AND code = 'PLATFORM_ADMIN';

    IF v_platform_tenant_id IS NULL THEN
        RAISE NOTICE 'Platform tenant not found, skipping seed data update';
        RETURN;
    END IF;

    IF v_admin_user_id IS NOT NULL AND v_platform_tenant_id IS NOT NULL THEN
        INSERT INTO user_tenant_registrations (
            user_id, tenant_id, registration_type, status, metadata
        )
        VALUES (
            v_admin_user_id,
            v_platform_tenant_id,
            'MEMBER',
            'ACTIVE',
            '{"migrated_from": "seed_data_v1", "note": "Platform admin auto-linked during v2 migration"}'::jsonb
        )
        ON CONFLICT (user_id, tenant_id, registration_type) DO NOTHING;

        RAISE NOTICE 'Linked platform admin user to platform tenant via user_tenant_registrations';
    END IF;

    IF v_iam_app_id IS NOT NULL THEN
        INSERT INTO permissions (application_id, code, name, resource_type, action, status)
        VALUES (
            v_iam_app_id,
            'role:assign',
            'Assign Role',
            'role',
            'assign',
            'ACTIVE'
        )
        ON CONFLICT DO NOTHING;


        SELECT id INTO v_role_assign_perm_id
        FROM permissions
        WHERE application_id = v_iam_app_id AND code = 'role:assign';

        RAISE NOTICE 'Ensured role:assign permission exists with ID: %', v_role_assign_perm_id;
    END IF;

    IF v_platform_admin_role_id IS NOT NULL AND v_role_assign_perm_id IS NOT NULL THEN
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES (v_platform_admin_role_id, v_role_assign_perm_id)
        ON CONFLICT DO NOTHING;

        RAISE NOTICE 'Assigned role:assign permission to PLATFORM_ADMIN role';
    END IF;

    RAISE NOTICE '=== V2 seed data migration complete ===';
END $$;
