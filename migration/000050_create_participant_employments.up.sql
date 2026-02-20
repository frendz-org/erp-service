CREATE TABLE IF NOT EXISTS participant_employments (
    -- Primary Key
    id                      UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Key (intra-domain)
    participant_id          UUID NOT NULL,

    -- Employment Data
    personnel_number        VARCHAR(50) NULL,
    date_of_hire            DATE NULL,

    -- Organization Hierarchy (string codes and names - denormalized for display)
    corporate_group_name    VARCHAR(255) NULL,
    legal_entity_code       VARCHAR(50) NULL,
    legal_entity_name       VARCHAR(255) NULL,
    business_unit_code      VARCHAR(50) NULL,
    business_unit_name      VARCHAR(255) NULL,
    tenant_name             VARCHAR(255) NULL,

    -- Job Details
    employment_status       VARCHAR(50) NULL,
    position_name           VARCHAR(255) NULL,
    job_level               VARCHAR(50) NULL,

    -- Location
    location_code           VARCHAR(50) NULL,
    location_name           VARCHAR(255) NULL,
    sub_location_name       VARCHAR(255) NULL,

    -- Retirement
    retirement_date         DATE NULL,
    retirement_type_code    VARCHAR(50) NULL,

    -- Audit Fields
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ,

    -- Optimistic Locking
    version                 INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_participant_employments_participant FOREIGN KEY (participant_id)
        REFERENCES participants(id) ON DELETE CASCADE
);

CREATE TRIGGER trg_participant_employments_updated_at
    BEFORE UPDATE ON participant_employments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participant_employments IS 'Employment/job details for a participant. Denormalized organization hierarchy for display.';
COMMENT ON COLUMN participant_employments.personnel_number IS 'HR system personnel/employee number.';
COMMENT ON COLUMN participant_employments.employment_status IS 'Masterdata string code from EMPLOYMENT_STATUS category.';
COMMENT ON COLUMN participant_employments.job_level IS 'Masterdata string code from JOB_LEVEL category.';
COMMENT ON COLUMN participant_employments.retirement_type_code IS 'Masterdata string code from RETIREMENT_TYPE category.';
COMMENT ON COLUMN participant_employments.corporate_group_name IS 'Denormalized corporate group name for display. Not a FK.';
COMMENT ON COLUMN participant_employments.legal_entity_code IS 'Denormalized legal entity code for display.';
COMMENT ON COLUMN participant_employments.tenant_name IS 'Denormalized tenant name at time of employment record creation.';
