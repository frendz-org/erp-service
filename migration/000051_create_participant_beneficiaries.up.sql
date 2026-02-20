CREATE TABLE IF NOT EXISTS participant_beneficiaries (
    -- Primary Key
    id                              UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys (intra-domain)
    participant_id                  UUID NOT NULL,
    family_member_id                UUID NOT NULL,

    -- File References (MinIO object keys)
    identity_photo_file_path        VARCHAR(1024) NULL,
    family_card_photo_file_path     VARCHAR(1024) NULL,
    bank_book_photo_file_path       VARCHAR(1024) NULL,

    -- Bank Account
    account_number                  VARCHAR(50) NULL,

    -- Audit Fields
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                      TIMESTAMPTZ,

    -- Optimistic Locking
    version                         INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_participant_beneficiaries_participant FOREIGN KEY (participant_id)
        REFERENCES participants(id) ON DELETE CASCADE,
    CONSTRAINT fk_participant_beneficiaries_family_member FOREIGN KEY (family_member_id)
        REFERENCES participant_family_members(id) ON DELETE RESTRICT
);

CREATE TRIGGER trg_participant_beneficiaries_updated_at
    BEFORE UPDATE ON participant_beneficiaries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participant_beneficiaries IS 'Beneficiary designations linking to family members. A beneficiary must be an existing family member.';
COMMENT ON COLUMN participant_beneficiaries.family_member_id IS 'FK to participant_family_members. ON DELETE RESTRICT - cannot delete a family member who is a beneficiary.';
COMMENT ON COLUMN participant_beneficiaries.identity_photo_file_path IS 'MinIO object key for beneficiary identity photo.';
COMMENT ON COLUMN participant_beneficiaries.family_card_photo_file_path IS 'MinIO object key for family card (kartu keluarga) photo.';
COMMENT ON COLUMN participant_beneficiaries.bank_book_photo_file_path IS 'MinIO object key for bank book photo.';
COMMENT ON COLUMN participant_beneficiaries.account_number IS 'Beneficiary bank account number for benefit disbursement.';
