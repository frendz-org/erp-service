ALTER TABLE participant_identities ADD COLUMN IF NOT EXISTS photo_file_id UUID;
ALTER TABLE participant_family_members ADD COLUMN IF NOT EXISTS supporting_doc_file_id UUID;
ALTER TABLE participant_beneficiaries ADD COLUMN IF NOT EXISTS identity_photo_file_id UUID;
ALTER TABLE participant_beneficiaries ADD COLUMN IF NOT EXISTS family_card_photo_file_id UUID;
ALTER TABLE participant_beneficiaries ADD COLUMN IF NOT EXISTS bank_book_photo_file_id UUID;

-- Indexes for FK lookups on the new file-ID columns
CREATE INDEX IF NOT EXISTS idx_participant_identities_photo_file_id
    ON participant_identities (photo_file_id)
    WHERE photo_file_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_participant_family_members_supporting_doc_file_id
    ON participant_family_members (supporting_doc_file_id)
    WHERE supporting_doc_file_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_participant_beneficiaries_identity_photo_file_id
    ON participant_beneficiaries (identity_photo_file_id)
    WHERE identity_photo_file_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_participant_beneficiaries_family_card_photo_file_id
    ON participant_beneficiaries (family_card_photo_file_id)
    WHERE family_card_photo_file_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_participant_beneficiaries_bank_book_photo_file_id
    ON participant_beneficiaries (bank_book_photo_file_id)
    WHERE bank_book_photo_file_id IS NOT NULL;
