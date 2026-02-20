DO $$
DECLARE
    v_platform_tenant_id UUID;
    v_iam_app_id UUID;
    v_frendz_app_id UUID;
    v_creator_role_id UUID;
    v_approver_role_id UUID;
BEGIN

    SELECT id INTO v_platform_tenant_id FROM tenants WHERE code = 'platform';

    IF v_platform_tenant_id IS NULL THEN
        RAISE NOTICE 'Platform tenant not found, skipping participant seed';
        RETURN;
    END IF;


    SELECT id INTO v_iam_app_id
    FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'iam-admin';

    IF v_iam_app_id IS NULL THEN
        RAISE NOTICE 'IAM admin application not found, skipping participant seed';
        RETURN;
    END IF;


    INSERT INTO applications (tenant_id, code, name, description, status)
    VALUES (
        v_platform_tenant_id,
        'frendz-saving',
        'Frendz Saving',
        'Participant management application for Frendz Saving program',
        'ACTIVE'
    )
    ON CONFLICT (tenant_id, code) DO NOTHING;

    SELECT id INTO v_frendz_app_id
    FROM applications
    WHERE tenant_id = v_platform_tenant_id AND code = 'frendz-saving';

    RAISE NOTICE 'Ensured frendz-saving application with ID: %', v_frendz_app_id;


    INSERT INTO permissions (application_id, code, name, resource_type, action, status) VALUES
        (v_iam_app_id, 'participant:create',  'Create Participant',  'participant', 'create',  'ACTIVE'),
        (v_iam_app_id, 'participant:read',    'View Participant',    'participant', 'read',    'ACTIVE'),
        (v_iam_app_id, 'participant:update',  'Update Participant',  'participant', 'update',  'ACTIVE'),
        (v_iam_app_id, 'participant:delete',  'Delete Participant',  'participant', 'delete',  'ACTIVE'),
        (v_iam_app_id, 'participant:submit',  'Submit Participant',  'participant', 'submit',  'ACTIVE'),
        (v_iam_app_id, 'participant:approve', 'Approve Participant', 'participant', 'approve', 'ACTIVE'),
        (v_iam_app_id, 'participant:reject',  'Reject Participant',  'participant', 'reject',  'ACTIVE')
    ON CONFLICT DO NOTHING;

    RAISE NOTICE 'Ensured 7 participant permissions';


    INSERT INTO roles (application_id, code, name, description, is_system, status)
    VALUES (
        v_iam_app_id,
        'PARTICIPANT_CREATOR',
        'Participant Creator',
        'Can create, edit, and submit participants for approval',
        TRUE,
        'ACTIVE'
    )
    ON CONFLICT (application_id, code) DO NOTHING;

    SELECT id INTO v_creator_role_id
    FROM roles
    WHERE application_id = v_iam_app_id AND code = 'PARTICIPANT_CREATOR';

    RAISE NOTICE 'Ensured PARTICIPANT_CREATOR role with ID: %', v_creator_role_id;


    INSERT INTO roles (application_id, code, name, description, is_system, status)
    VALUES (
        v_iam_app_id,
        'PARTICIPANT_APPROVER',
        'Participant Approver',
        'Can approve or reject submitted participants',
        TRUE,
        'ACTIVE'
    )
    ON CONFLICT (application_id, code) DO NOTHING;

    SELECT id INTO v_approver_role_id
    FROM roles
    WHERE application_id = v_iam_app_id AND code = 'PARTICIPANT_APPROVER';

    RAISE NOTICE 'Ensured PARTICIPANT_APPROVER role with ID: %', v_approver_role_id;



    IF v_creator_role_id IS NOT NULL THEN
        INSERT INTO role_permissions (role_id, permission_id)
        SELECT v_creator_role_id, p.id
        FROM permissions p
        WHERE p.application_id = v_iam_app_id
          AND p.code IN ('participant:create', 'participant:read', 'participant:update', 'participant:delete', 'participant:submit')
        ON CONFLICT DO NOTHING;

        RAISE NOTICE 'Assigned 5 permissions to PARTICIPANT_CREATOR role';
    END IF;



    IF v_approver_role_id IS NOT NULL THEN
        INSERT INTO role_permissions (role_id, permission_id)
        SELECT v_approver_role_id, p.id
        FROM permissions p
        WHERE p.application_id = v_iam_app_id
          AND p.code IN ('participant:read', 'participant:approve', 'participant:reject')
        ON CONFLICT DO NOTHING;

        RAISE NOTICE 'Assigned 3 permissions to PARTICIPANT_APPROVER role';
    END IF;


    INSERT INTO role_permissions (role_id, permission_id)
    SELECT r.id, p.id
    FROM roles r, permissions p
    WHERE r.application_id = v_iam_app_id
      AND r.code = 'PLATFORM_ADMIN'
      AND p.application_id = v_iam_app_id
      AND p.code LIKE 'participant:%'
    ON CONFLICT DO NOTHING;

    RAISE NOTICE 'Assigned participant permissions to PLATFORM_ADMIN role';
    RAISE NOTICE '=== Participant application, roles, and permissions seeded successfully ===';
END $$;
