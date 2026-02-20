-- Prevents TOCTOU race on concurrent self-registrations with the same KTP.
-- Partial index excludes soft-deleted rows.
CREATE UNIQUE INDEX uk_participants_ktp_tenant_app
    ON participants (ktp_number, tenant_id, application_id)
    WHERE deleted_at IS NULL;
