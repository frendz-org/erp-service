CREATE INDEX idx_participants_tenant_id ON participants(tenant_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participants_tenant_application ON participants(tenant_id, application_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participants_tenant_status ON participants(tenant_id, status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participants_ktp_number ON participants(ktp_number)
    WHERE ktp_number IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_participants_employee_number ON participants(employee_number)
    WHERE employee_number IS NOT NULL AND deleted_at IS NULL;
-- User linking: unique participant per user within a tenant+application (deferred Phase 2)
CREATE UNIQUE INDEX uq_participants_user_tenant_app ON participants(user_id, tenant_id, application_id)
    WHERE user_id IS NOT NULL AND deleted_at IS NULL;

-- Soft-delete filter
CREATE INDEX idx_participants_deleted_at ON participants(deleted_at);
CREATE INDEX idx_participant_identities_participant_id ON participant_identities(participant_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participant_identities_deleted_at ON participant_identities(deleted_at);
CREATE INDEX idx_participant_addresses_participant_id ON participant_addresses(participant_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participant_addresses_deleted_at ON participant_addresses(deleted_at);
CREATE INDEX idx_participant_bank_accounts_participant_id ON participant_bank_accounts(participant_id)
    WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX uq_participant_bank_accounts_primary ON participant_bank_accounts(participant_id)
    WHERE is_primary = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_participant_bank_accounts_deleted_at ON participant_bank_accounts(deleted_at);
CREATE INDEX idx_participant_family_members_participant_id ON participant_family_members(participant_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participant_family_members_deleted_at ON participant_family_members(deleted_at);
CREATE INDEX idx_participant_employments_participant_id ON participant_employments(participant_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participant_employments_deleted_at ON participant_employments(deleted_at);
CREATE INDEX idx_participant_beneficiaries_participant_id ON participant_beneficiaries(participant_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participant_beneficiaries_family_member_id ON participant_beneficiaries(family_member_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_participant_beneficiaries_deleted_at ON participant_beneficiaries(deleted_at);
CREATE INDEX idx_participant_status_history_participant_id ON participant_status_history(participant_id, changed_at DESC);
CREATE INDEX idx_participant_status_history_changed_by ON participant_status_history(changed_by);
