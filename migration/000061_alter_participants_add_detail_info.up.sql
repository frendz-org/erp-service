-- ============================================================================
-- CREATE TABLE: participant_pensions
-- Description: Pension detail information for a participant.
--              Separate table from participants because this data relates to
--              pension program enrollment, not participant identity.
--              Balance-related fields are intentionally excluded — they belong
--              in a dedicated ledger table (future: participant_balances).
-- Relationship: 1:1 owned child of participants (ON DELETE CASCADE).
-- ============================================================================

CREATE TABLE IF NOT EXISTS participant_pensions (
    -- Primary Key
    id                          UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Key (intra-domain, same as other participant child tables)
    participant_id              UUID NOT NULL,

    -- Pension Detail Information
    participant_number          VARCHAR(50) NULL,
    pension_category            VARCHAR(50) NULL,
    pension_status              VARCHAR(50) NULL,
    effective_date              DATE NULL,
    end_date                    DATE NULL,
    projected_retirement_date   DATE NULL,

    -- Audit Fields
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ,

    -- Optimistic Locking
    version                     INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_participant_pensions_participant FOREIGN KEY (participant_id)
        REFERENCES participants(id) ON DELETE CASCADE,

    -- Date range: end_date >= effective_date when both are set
    CONSTRAINT chk_participant_pensions_dates CHECK (
        effective_date IS NULL
        OR end_date IS NULL
        OR end_date >= effective_date
    )
);

-- One pension record per participant (1:1 relationship)
CREATE UNIQUE INDEX IF NOT EXISTS uq_participant_pensions_participant
    ON participant_pensions(participant_id)
    WHERE deleted_at IS NULL;

-- Search by participant_number (critical lookup path)
CREATE INDEX IF NOT EXISTS idx_participant_pensions_number
    ON participant_pensions(participant_number)
    WHERE participant_number IS NOT NULL AND deleted_at IS NULL;

-- Filter by pension_status (list page filter)
CREATE INDEX IF NOT EXISTS idx_participant_pensions_status
    ON participant_pensions(pension_status)
    WHERE deleted_at IS NULL;

-- Upcoming retirements (reporting/batch jobs)
CREATE INDEX IF NOT EXISTS idx_participant_pensions_retirement_date
    ON participant_pensions(projected_retirement_date)
    WHERE projected_retirement_date IS NOT NULL AND deleted_at IS NULL;

CREATE TRIGGER trg_participant_pensions_updated_at
    BEFORE UPDATE ON participant_pensions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participant_pensions IS 'Pension program enrollment details for a participant. 1:1 with participants. Balance tracking is separate.';
COMMENT ON COLUMN participant_pensions.participant_number IS 'Unique pension participation number. Assigned on enrollment or import.';
COMMENT ON COLUMN participant_pensions.pension_category IS 'Masterdata string code from PARTICIPANT_PENSION_CATEGORY (INACTIVE/ACTIVE/PASSIVE). Validated at application level, no DB FK.';
COMMENT ON COLUMN participant_pensions.pension_status IS 'Masterdata string code from PARTICIPANT_PENSION_STATUS (child of pension_category). Validated at application level, no DB FK.';
COMMENT ON COLUMN participant_pensions.effective_date IS 'Date pension participation becomes effective.';
COMMENT ON COLUMN participant_pensions.end_date IS 'Date pension participation ends. Must be >= effective_date.';
COMMENT ON COLUMN participant_pensions.projected_retirement_date IS 'Projected retirement date for the participant.';
