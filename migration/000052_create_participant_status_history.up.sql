-- Append-only audit trail table for participant status transitions.
-- EXCEPTION: No deleted_at column (append-only, records are never soft-deleted).
-- EXCEPTION: No version column (append-only, records are never updated).
CREATE TABLE IF NOT EXISTS participant_status_history (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Key (intra-domain)
    participant_id      UUID NOT NULL,

    -- Status Transition
    from_status         VARCHAR(20) NULL,
    to_status           VARCHAR(20) NOT NULL,

    -- Actor (UUID column, NOT FK to users - cross-domain boundary)
    changed_by          UUID NOT NULL,

    -- Context
    reason              TEXT NULL,
    changed_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Audit Fields (created_at and updated_at only - append-only table)
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_participant_status_history_participant FOREIGN KEY (participant_id)
        REFERENCES participants(id) ON DELETE CASCADE,
    CONSTRAINT chk_participant_status_history_from_status CHECK (
        from_status IS NULL OR from_status IN ('DRAFT', 'PENDING_APPROVAL', 'APPROVED', 'REJECTED')
    ),
    CONSTRAINT chk_participant_status_history_to_status CHECK (
        to_status IN ('DRAFT', 'PENDING_APPROVAL', 'APPROVED', 'REJECTED')
    )
);

CREATE TRIGGER trg_participant_status_history_updated_at
    BEFORE UPDATE ON participant_status_history
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participant_status_history IS 'Append-only audit trail of participant status transitions. No deleted_at or version - records are immutable.';
COMMENT ON COLUMN participant_status_history.from_status IS 'Previous status. NULL for initial creation (DRAFT).';
COMMENT ON COLUMN participant_status_history.to_status IS 'New status after transition.';
COMMENT ON COLUMN participant_status_history.changed_by IS 'UUID of user who triggered the transition. No FK - cross-domain boundary.';
COMMENT ON COLUMN participant_status_history.reason IS 'Free-text reason for transition. Required for REJECTED, optional otherwise.';
COMMENT ON COLUMN participant_status_history.changed_at IS 'Timestamp of the status change event.';
