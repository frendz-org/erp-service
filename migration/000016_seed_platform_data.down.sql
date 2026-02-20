DO $$
DECLARE
    v_tenant_id UUID;
    v_application_id UUID;
    v_user_id UUID;
BEGIN
    -- Get the platform tenant ID
    SELECT id INTO v_tenant_id FROM tenants WHERE code = 'platform';

    IF v_tenant_id IS NULL THEN
        RAISE NOTICE 'Platform tenant not found, nothing to delete';
        RETURN;
    END IF;

    -- Get the IAM admin application ID
    SELECT id INTO v_application_id
    FROM applications
    WHERE tenant_id = v_tenant_id AND code = 'iam-admin';

    -- Get the platform admin user ID
    SELECT id INTO v_user_id
    FROM users
    WHERE tenant_id = v_tenant_id AND email = 'admin@platform.local';

    -- Delete in reverse dependency order

    -- 1. Remove user role assignments
    IF v_user_id IS NOT NULL THEN
        DELETE FROM user_role_assignments WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted user role assignments';
    END IF;

    -- 2. Remove admin user security state
    IF v_user_id IS NOT NULL THEN
        DELETE FROM user_security_states WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted admin user security state';
    END IF;

    -- 3. Remove admin user profile
    IF v_user_id IS NOT NULL THEN
        DELETE FROM user_profiles WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted admin user profile';
    END IF;

    -- 4. Remove admin user auth methods
    IF v_user_id IS NOT NULL THEN
        DELETE FROM user_auth_methods WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted admin user auth methods';
    END IF;

    -- 5. Remove admin user
    IF v_user_id IS NOT NULL THEN
        DELETE FROM users WHERE id = v_user_id;
        RAISE NOTICE 'Deleted admin user';
    END IF;

    -- 6. Remove role permissions (for roles in this application)
    IF v_application_id IS NOT NULL THEN
        DELETE FROM role_permissions
        WHERE role_id IN (SELECT id FROM roles WHERE application_id = v_application_id);
        RAISE NOTICE 'Deleted role permissions';
    END IF;

    -- 7. Remove permissions
    IF v_application_id IS NOT NULL THEN
        DELETE FROM permissions WHERE application_id = v_application_id;
        RAISE NOTICE 'Deleted permissions';
    END IF;

    -- 8. Remove roles (including PLATFORM_ADMIN)
    IF v_application_id IS NOT NULL THEN
        DELETE FROM roles WHERE application_id = v_application_id;
        RAISE NOTICE 'Deleted roles';
    END IF;

    -- 9. Remove IAM admin application
    IF v_application_id IS NOT NULL THEN
        DELETE FROM applications WHERE id = v_application_id;
        RAISE NOTICE 'Deleted IAM admin application';
    END IF;

    -- 10. Remove platform tenant
    DELETE FROM tenants WHERE id = v_tenant_id;
    RAISE NOTICE 'Deleted platform tenant';

    RAISE NOTICE '=== Platform seed data removed successfully ===';
END $$;
