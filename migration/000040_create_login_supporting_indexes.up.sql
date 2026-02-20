-- Login-supporting indexes for JWT population queries.
-- After OTP verification, the system fetches all active tenant registrations
-- and role assignments for the user to populate JWT claims.

-- JWT population: fetch all active tenant registrations for a user
CREATE INDEX IF NOT EXISTS idx_utr_user_active
    ON user_tenant_registrations(user_id)
    WHERE status = 'ACTIVE' AND deleted_at IS NULL;

-- JWT population: fetch all active role assignments for a user
CREATE INDEX IF NOT EXISTS idx_ura_user_active_lean
    ON user_role_assignments(user_id)
    WHERE status = 'ACTIVE' AND deleted_at IS NULL;
