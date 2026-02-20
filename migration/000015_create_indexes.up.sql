-- Branches
CREATE INDEX idx_branches_tenant_id ON branches(tenant_id) WHERE deleted_at IS NULL;

-- Users
CREATE INDEX idx_users_tenant_id ON users(tenant_id) WHERE deleted_at IS NULL;

-- User Auth Methods
CREATE INDEX idx_user_auth_methods_user_id ON user_auth_methods(user_id);

-- User Branches
CREATE INDEX idx_user_branches_user_id ON user_branches(user_id);
CREATE INDEX idx_user_branches_branch_id ON user_branches(branch_id);

-- Password History
CREATE INDEX idx_password_history_user_id ON password_history(user_id);

-- Applications
CREATE INDEX idx_applications_tenant_id ON applications(tenant_id) WHERE deleted_at IS NULL;

-- Roles
CREATE INDEX idx_roles_application_id ON roles(application_id) WHERE deleted_at IS NULL;

-- Permissions
CREATE INDEX idx_permissions_application_id ON permissions(application_id) WHERE deleted_at IS NULL;

-- Role Permissions
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- User Role Assignments
CREATE INDEX idx_ura_user_id ON user_role_assignments(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ura_role_id ON user_role_assignments(role_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ura_branch_id ON user_role_assignments(branch_id)
    WHERE deleted_at IS NULL AND branch_id IS NOT NULL;

-- User lookup by email (login flow) - CRITICAL
CREATE INDEX idx_users_tenant_email_active ON users(tenant_id, LOWER(email))
    WHERE deleted_at IS NULL AND status IN ('ACTIVE', 'PENDING_APPROVAL');

-- User auth method lookup (login flow: find PASSWORD method for user)
CREATE INDEX idx_user_auth_methods_user_type ON user_auth_methods(user_id, method_type)
    WHERE is_active = TRUE;

-- Google OAuth lookup (find user by Google ID)
CREATE INDEX idx_user_auth_methods_google ON user_auth_methods(
    (credential_data->>'google_id')
) WHERE method_type = 'GOOGLE' AND is_active = TRUE;

-- Users pending approval (admin dashboard)
CREATE INDEX idx_users_pending_approval ON users(tenant_id, created_at)
    WHERE status = 'PENDING_APPROVAL' AND deleted_at IS NULL;

-- Application lookup by code
CREATE INDEX idx_applications_tenant_code_active ON applications(tenant_id, code)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Role lookup by code within application
CREATE INDEX idx_roles_app_code_active ON roles(application_id, code)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Permission lookup by code within application
CREATE INDEX idx_permissions_app_code_active ON permissions(application_id, code)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Get all active roles for a user (permission check flow)
CREATE INDEX idx_ura_user_active ON user_role_assignments(user_id, status)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Get all users with a specific role (admin view)
CREATE INDEX idx_ura_role_users ON user_role_assignments(role_id, user_id)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Password history lookup (last N passwords)
CREATE INDEX idx_password_history_user_recent ON password_history(user_id, created_at DESC);

-- Tenant code lookup (login flow)
CREATE INDEX idx_tenants_code_active ON tenants(code)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Branch code lookup within tenant
CREATE INDEX idx_branches_tenant_code_active ON branches(tenant_id, code)
    WHERE deleted_at IS NULL AND status = 'ACTIVE';

-- Security state: locked users (for admin monitoring)
CREATE INDEX idx_user_security_locked ON user_security_states(locked_until)
    WHERE locked_until IS NOT NULL;

-- Security state: users needing forced password change (admin view)
CREATE INDEX idx_user_security_force_pw ON user_security_states(user_id)
    WHERE force_password_change = TRUE;
