DO $$
DECLARE
    v_tenant_id UUID;
    v_application_id UUID;
    v_role_id UUID;
    v_user_id UUID;
BEGIN
    -- 1. Create Platform Tenant
    INSERT INTO tenants (code, name, status, settings)
    VALUES (
        'platform',
        'Platform Administration',
        'ACTIVE',
        '{
            "approval_required": false,
            "is_platform_tenant": true,
            "password_policy": {
                "min_length": 12,
                "require_uppercase": true,
                "require_lowercase": true,
                "require_number": true,
                "require_special": true,
                "history_count": 5,
                "max_age_days": 90
            },
            "session": {
                "access_token_ttl_minutes": 30,
                "refresh_token_ttl_days": 7,
                "max_concurrent_sessions": 3
            }
        }'::jsonb
    )
    RETURNING id INTO v_tenant_id;

    RAISE NOTICE 'Created platform tenant with ID: %', v_tenant_id;

    -- 2. Create IAM Admin Application
    INSERT INTO applications (tenant_id, code, name, description, status)
    VALUES (
        v_tenant_id,
        'iam-admin',
        'IAM Administration',
        'Internal application for managing the IAM system',
        'ACTIVE'
    )
    RETURNING id INTO v_application_id;

    RAISE NOTICE 'Created IAM admin application with ID: %', v_application_id;

    -- 3. Create Platform Admin Role
    INSERT INTO roles (application_id, code, name, description, is_system, status)
    VALUES (
        v_application_id,
        'PLATFORM_ADMIN',
        'Platform Administrator',
        'Full access to all IAM operations across all tenants',
        TRUE,
        'ACTIVE'
    )
    RETURNING id INTO v_role_id;

    RAISE NOTICE 'Created platform admin role with ID: %', v_role_id;

    -- 4. Create Core IAM Permissions
    INSERT INTO permissions (application_id, code, name, resource_type, action, status) VALUES
        (v_application_id, 'tenant:create', 'Create Tenant', 'tenant', 'create', 'ACTIVE'),
        (v_application_id, 'tenant:read', 'View Tenant', 'tenant', 'read', 'ACTIVE'),
        (v_application_id, 'tenant:update', 'Update Tenant', 'tenant', 'update', 'ACTIVE'),
        (v_application_id, 'tenant:delete', 'Delete Tenant', 'tenant', 'delete', 'ACTIVE'),
        (v_application_id, 'user:create', 'Create User', 'user', 'create', 'ACTIVE'),
        (v_application_id, 'user:read', 'View User', 'user', 'read', 'ACTIVE'),
        (v_application_id, 'user:update', 'Update User', 'user', 'update', 'ACTIVE'),
        (v_application_id, 'user:delete', 'Delete User', 'user', 'delete', 'ACTIVE'),
        (v_application_id, 'user:approve', 'Approve User', 'user', 'approve', 'ACTIVE'),
        (v_application_id, 'application:create', 'Create Application', 'application', 'create', 'ACTIVE'),
        (v_application_id, 'application:read', 'View Application', 'application', 'read', 'ACTIVE'),
        (v_application_id, 'application:update', 'Update Application', 'application', 'update', 'ACTIVE'),
        (v_application_id, 'application:delete', 'Delete Application', 'application', 'delete', 'ACTIVE'),
        (v_application_id, 'role:create', 'Create Role', 'role', 'create', 'ACTIVE'),
        (v_application_id, 'role:read', 'View Role', 'role', 'read', 'ACTIVE'),
        (v_application_id, 'role:update', 'Update Role', 'role', 'update', 'ACTIVE'),
        (v_application_id, 'role:delete', 'Delete Role', 'role', 'delete', 'ACTIVE'),
        (v_application_id, 'permission:create', 'Create Permission', 'permission', 'create', 'ACTIVE'),
        (v_application_id, 'permission:read', 'View Permission', 'permission', 'read', 'ACTIVE'),
        (v_application_id, 'permission:update', 'Update Permission', 'permission', 'update', 'ACTIVE'),
        (v_application_id, 'permission:delete', 'Delete Permission', 'permission', 'delete', 'ACTIVE'),
        (v_application_id, 'audit:read', 'View Audit Logs', 'audit', 'read', 'ACTIVE'),
        (v_application_id, 'audit:export', 'Export Audit Logs', 'audit', 'export', 'ACTIVE');

    RAISE NOTICE 'Created 23 IAM permissions';

    -- 5. Assign all permissions to Platform Admin role
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT v_role_id, p.id
    FROM permissions p
    WHERE p.application_id = v_application_id;

    RAISE NOTICE 'Assigned all permissions to platform admin role';

    -- 6. Create initial platform admin user (core identity)
    INSERT INTO users (tenant_id, email, status, registration_source)
    VALUES (
        v_tenant_id,
        'admin@platform.local',
        'ACTIVE',
        'ADMIN'
    )
    RETURNING id INTO v_user_id;

    RAISE NOTICE 'Created platform admin user with ID: %', v_user_id;

    -- 7. Create admin user auth methods (PASSWORD + PIN)
    -- Password: ChangeMe123!
    INSERT INTO user_auth_methods (user_id, method_type, credential_data, is_active)
    VALUES (
        v_user_id,
        'PASSWORD',
        '{"hash": "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4/tPBxqK1/XW.zSm", "last_changed_at": null}'::jsonb,
        TRUE
    );

    -- PIN: 123456
    INSERT INTO user_auth_methods (user_id, method_type, credential_data, is_active)
    VALUES (
        v_user_id,
        'PIN',
        '{"hash": "$2a$12$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", "pin_set": true, "last_changed_at": null}'::jsonb,
        TRUE
    );

    RAISE NOTICE 'Created admin auth methods (PASSWORD + PIN)';

    -- 8. Create admin user profile
    INSERT INTO user_profiles (user_id, first_name, last_name)
    VALUES (
        v_user_id,
        'Platform',
        'Administrator'
    );

    RAISE NOTICE 'Created admin user profile';

    -- 9. Create admin user security state
    INSERT INTO user_security_states (user_id, email_verified, email_verified_at, pin_verified, force_password_change)
    VALUES (
        v_user_id,
        TRUE,
        NOW(),
        TRUE,
        TRUE  -- Force password change on first login
    );

    RAISE NOTICE 'Created admin user security state';

    -- 10. Assign Platform Admin role to admin user
    INSERT INTO user_role_assignments (user_id, role_id, assigned_by, status)
    VALUES (
        v_user_id,
        v_role_id,
        v_user_id,
        'ACTIVE'
    );

    RAISE NOTICE 'Assigned platform admin role to admin user';
    RAISE NOTICE '=== Platform seed data created successfully ===';
END $$;
