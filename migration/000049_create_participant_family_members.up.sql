CREATE TABLE IF NOT EXISTS participant_family_members (
    -- Primary Key
    id                          UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Key (intra-domain)
    participant_id              UUID NOT NULL,

    -- Family Member Data
    full_name                   VARCHAR(255) NOT NULL,
    relationship_type           VARCHAR(50) NOT NULL,
    is_dependent                BOOLEAN NOT NULL DEFAULT FALSE,

    -- File Reference (MinIO object key)
    supporting_doc_file_path    VARCHAR(1024) NULL,

    -- Audit Fields
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ,

    -- Optimistic Locking
    version                     INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_participant_family_members_participant FOREIGN KEY (participant_id)
        REFERENCES participants(id) ON DELETE CASCADE
);

CREATE TRIGGER trg_participant_family_members_updated_at
    BEFORE UPDATE ON participant_family_members
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participant_family_members IS 'Family members and dependents of a participant.';
COMMENT ON COLUMN participant_family_members.relationship_type IS 'Masterdata string code from FAMILY_RELATIONSHIP category.';
COMMENT ON COLUMN participant_family_members.is_dependent IS 'Whether this family member is financially dependent on the participant.';
COMMENT ON COLUMN participant_family_members.supporting_doc_file_path IS 'MinIO object key for supporting document (e.g., family card scan).';
