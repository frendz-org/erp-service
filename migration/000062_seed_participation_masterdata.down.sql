-- ============================================================================
-- ROLLBACK: Remove PARTICIPANT_PENSION_CATEGORY and PARTICIPANT_PENSION_STATUS
-- Delete children (status items) first, then category items, then categories.
-- ============================================================================

-- Delete status items first (child category items)
DELETE FROM masterdata_items
WHERE category_id IN (
    SELECT id FROM masterdata_categories
    WHERE code = 'PARTICIPANT_PENSION_STATUS' AND deleted_at IS NULL
);

-- Delete category items (parent category items)
DELETE FROM masterdata_items
WHERE category_id IN (
    SELECT id FROM masterdata_categories
    WHERE code = 'PARTICIPANT_PENSION_CATEGORY' AND deleted_at IS NULL
);

-- Delete child category first (FK to parent)
DELETE FROM masterdata_categories WHERE code = 'PARTICIPANT_PENSION_STATUS';

-- Delete parent category
DELETE FROM masterdata_categories WHERE code = 'PARTICIPANT_PENSION_CATEGORY';
