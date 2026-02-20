-- ============================================================================
-- ROLLBACK: Drop participant_pensions table
-- ============================================================================

DROP TRIGGER IF EXISTS trg_participant_pensions_updated_at ON participant_pensions;
DROP TABLE IF EXISTS participant_pensions;
