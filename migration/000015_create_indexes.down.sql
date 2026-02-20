-- Query Optimization Indexes
DROP INDEX IF EXISTS idx_user_security_force_pw;
DROP INDEX IF EXISTS idx_user_security_locked;
DROP INDEX IF EXISTS idx_branches_tenant_code_active;
DROP INDEX IF EXISTS idx_tenants_code_active;
DROP INDEX IF EXISTS idx_password_history_user_recent;
DROP INDEX IF EXISTS idx_ura_role_users;
DROP INDEX IF EXISTS idx_ura_user_active;
DROP INDEX IF EXISTS idx_permissions_app_code_active;
DROP INDEX IF EXISTS idx_roles_app_code_active;
DROP INDEX IF EXISTS idx_applications_tenant_code_active;
DROP INDEX IF EXISTS idx_users_pending_approval;
DROP INDEX IF EXISTS idx_user_auth_methods_google;
DROP INDEX IF EXISTS idx_user_auth_methods_user_type;
DROP INDEX IF EXISTS idx_users_tenant_email_active;

-- Foreign Key Indexes
DROP INDEX IF EXISTS idx_ura_branch_id;
DROP INDEX IF EXISTS idx_ura_role_id;
DROP INDEX IF EXISTS idx_ura_user_id;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP INDEX IF EXISTS idx_permissions_application_id;
DROP INDEX IF EXISTS idx_roles_application_id;
DROP INDEX IF EXISTS idx_applications_tenant_id;
DROP INDEX IF EXISTS idx_password_history_user_id;
DROP INDEX IF EXISTS idx_user_branches_branch_id;
DROP INDEX IF EXISTS idx_user_branches_user_id;
DROP INDEX IF EXISTS idx_user_auth_methods_user_id;
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_branches_tenant_id;
