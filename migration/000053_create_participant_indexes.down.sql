-- participant_status_history
DROP INDEX IF EXISTS idx_participant_status_history_changed_by;
DROP INDEX IF EXISTS idx_participant_status_history_participant_id;

-- participant_beneficiaries
DROP INDEX IF EXISTS idx_participant_beneficiaries_deleted_at;
DROP INDEX IF EXISTS idx_participant_beneficiaries_family_member_id;
DROP INDEX IF EXISTS idx_participant_beneficiaries_participant_id;

-- participant_employments
DROP INDEX IF EXISTS idx_participant_employments_deleted_at;
DROP INDEX IF EXISTS idx_participant_employments_participant_id;

-- participant_family_members
DROP INDEX IF EXISTS idx_participant_family_members_deleted_at;
DROP INDEX IF EXISTS idx_participant_family_members_participant_id;

-- participant_bank_accounts
DROP INDEX IF EXISTS idx_participant_bank_accounts_deleted_at;
DROP INDEX IF EXISTS uq_participant_bank_accounts_primary;
DROP INDEX IF EXISTS idx_participant_bank_accounts_participant_id;

-- participant_addresses
DROP INDEX IF EXISTS idx_participant_addresses_deleted_at;
DROP INDEX IF EXISTS idx_participant_addresses_participant_id;

-- participant_identities
DROP INDEX IF EXISTS idx_participant_identities_deleted_at;
DROP INDEX IF EXISTS idx_participant_identities_participant_id;

-- participants
DROP INDEX IF EXISTS idx_participants_deleted_at;
DROP INDEX IF EXISTS uq_participants_user_tenant_app;
DROP INDEX IF EXISTS idx_participants_employee_number;
DROP INDEX IF EXISTS idx_participants_ktp_number;
DROP INDEX IF EXISTS idx_participants_tenant_status;
DROP INDEX IF EXISTS idx_participants_tenant_application;
DROP INDEX IF EXISTS idx_participants_tenant_id;
