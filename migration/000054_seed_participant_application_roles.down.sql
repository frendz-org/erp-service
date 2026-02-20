DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_iam_app_id UUID;
    v_frendz_app_id UUID;
BEGIN
    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';

    IF v_platform_tenant_id IS NULL THEN
        RAISE NOTICE 'Platform tenant not found, nothing to delete';
        RETURN;
    END IF;

    SELECT id INTO v_iam_app_id
    FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'iam-admin';

    SELECT id INTO v_frendz_app_id
    FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'frendz-saving';

    -- 1. Remove participant permissions from PLATFORM_ADMIN role
    IF v_iam_app_id IS NOT NULL THEN
        DELETE FROM role_permissions
        WHERE permission_id IN (
            SELECT id FROM permissions
            WHERE application_id = v_iam_app_id AND code LIKE 'participant:%'
        );

        RAISE NOTICE 'Removed participant permission assignments';
    END IF;

    -- 2. Remove PARTICIPANT_CREATOR and PARTICIPANT_APPROVER roles
    --    (role_permissions cascade-deleted by FK)
    IF v_iam_app_id IS NOT NULL THEN
        DELETE FROM roles
        WHERE application_id = v_iam_app_id
          AND code IN ('PARTICIPANT_CREATOR', 'PARTICIPANT_APPROVER');

        RAISE NOTICE 'Removed participant roles';
    END IF;

    -- 3. Remove participant permissions
    IF v_iam_app_id IS NOT NULL THEN
        DELETE FROM permissions
        WHERE application_id = v_iam_app_id
          AND code LIKE 'participant:%';

        RAISE NOTICE 'Removed participant permissions';
    END IF;

    -- 4. Remove frendz-saving application
    IF v_frendz_app_id IS NOT NULL THEN
        DELETE FROM applications WHERE id = v_frendz_app_id;
        RAISE NOTICE 'Removed frendz-saving application';
    END IF;

    RAISE NOTICE '=== Participant application, roles, and permissions removed ===';
END $$;
