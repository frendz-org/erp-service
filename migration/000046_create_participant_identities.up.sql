CREATE TABLE IF NOT EXISTS participant_identities (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Key (intra-domain)
    participant_id      UUID NOT NULL,

    -- Identity Document Data
    identity_type       VARCHAR(50) NOT NULL,
    identity_number     VARCHAR(100) NOT NULL,
    identity_authority  VARCHAR(255) NULL,
    issue_date          DATE NULL,
    expiry_date         DATE NULL,

    -- File Reference (MinIO object key)
    photo_file_path     VARCHAR(1024) NULL,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_participant_identities_participant FOREIGN KEY (participant_id)
        REFERENCES participants(id) ON DELETE CASCADE
);

CREATE TRIGGER trg_participant_identities_updated_at
    BEFORE UPDATE ON participant_identities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participant_identities IS 'Identity documents (KTP, passport, etc.) for a participant. Multiple documents allowed.';
COMMENT ON COLUMN participant_identities.identity_type IS 'Masterdata string code from IDENTITY_TYPE category.';
COMMENT ON COLUMN participant_identities.identity_number IS 'Document number (e.g., KTP number, passport number).';
COMMENT ON COLUMN participant_identities.identity_authority IS 'Issuing authority name.';
COMMENT ON COLUMN participant_identities.photo_file_path IS 'MinIO object key for scanned identity document image.';
