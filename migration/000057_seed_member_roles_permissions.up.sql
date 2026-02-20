DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_iam_app_id UUID;
    v_product_admin_role_id UUID;
BEGIN

    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';

    IF v_platform_tenant_id IS NULL THEN
        RAISE NOTICE 'Platform tenant not found, skipping member seed';
        RETURN;
    END IF;

    SELECT id INTO v_iam_app_id
    FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'iam-admin';

    IF v_iam_app_id IS NULL THEN
        RAISE NOTICE 'IAM admin application not found, skipping member seed';
        RETURN;
    END IF;

    -- 1. Seed member permissions (6)
    INSERT INTO permissions (application_id, code, name, resource_type, action, status) VALUES
        (v_iam_app_id, 'member:read',       'View Member',       'member', 'read',       'ACTIVE'),
        (v_iam_app_id, 'member:approve',     'Approve Member',    'member', 'approve',    'ACTIVE'),
        (v_iam_app_id, 'member:reject',      'Reject Member',     'member', 'reject',     'ACTIVE'),
        (v_iam_app_id, 'member:update',      'Update Member',     'member', 'update',     'ACTIVE'),
        (v_iam_app_id, 'member:deactivate',  'Deactivate Member', 'member', 'deactivate', 'ACTIVE'),
        (v_iam_app_id, 'member:register',    'Register Member',   'member', 'register',   'ACTIVE')
    ON CONFLICT DO NOTHING;

    RAISE NOTICE 'Ensured 6 member permissions';

    -- 2. Seed TENANT_PRODUCT_ADMIN role
    INSERT INTO roles (application_id, code, name, description, is_system, status)
    VALUES (
        v_iam_app_id,
        'TENANT_PRODUCT_ADMIN',
        'Product Administrator',
        'Manages product membership lifecycle: approve/reject members, assign roles, deactivate members',
        TRUE,
        'ACTIVE'
    )
    ON CONFLICT (application_id, code) DO NOTHING;

    SELECT id INTO v_product_admin_role_id
    FROM roles
    WHERE application_id = v_iam_app_id AND code = 'TENANT_PRODUCT_ADMIN';

    RAISE NOTICE 'Ensured TENANT_PRODUCT_ADMIN role with ID: %', v_product_admin_role_id;

    -- 3. Assign all member:* permissions to TENANT_PRODUCT_ADMIN
    IF v_product_admin_role_id IS NOT NULL THEN
        INSERT INTO role_permissions (role_id, permission_id)
        SELECT v_product_admin_role_id, p.id
        FROM permissions p
        WHERE p.application_id = v_iam_app_id
          AND p.code LIKE 'member:%'
        ON CONFLICT DO NOTHING;

        RAISE NOTICE 'Assigned member permissions to TENANT_PRODUCT_ADMIN';
    END IF;

    -- 4. Assign all participant:* permissions to TENANT_PRODUCT_ADMIN
    IF v_product_admin_role_id IS NOT NULL THEN
        INSERT INTO role_permissions (role_id, permission_id)
        SELECT v_product_admin_role_id, p.id
        FROM permissions p
        WHERE p.application_id = v_iam_app_id
          AND p.code LIKE 'participant:%'
        ON CONFLICT DO NOTHING;

        RAISE NOTICE 'Assigned participant permissions to TENANT_PRODUCT_ADMIN';
    END IF;

    -- 5. Grant all member:* permissions to PLATFORM_ADMIN
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT r.id, p.id
    FROM roles r, permissions p
    WHERE r.application_id = v_iam_app_id
      AND r.code = 'PLATFORM_ADMIN'
      AND p.application_id = v_iam_app_id
      AND p.code LIKE 'member:%'
    ON CONFLICT DO NOTHING;

    RAISE NOTICE 'Assigned member permissions to PLATFORM_ADMIN role';
    RAISE NOTICE '=== Member roles and permissions seeded successfully ===';
END $$;
