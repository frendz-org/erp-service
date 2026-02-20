-- ============================================================================
-- ROLLBACK: Remove ORGANIZATION Category and Items
-- ============================================================================

-- Delete items first (foreign key constraint)
DELETE FROM masterdata_items
WHERE category_id = (
    SELECT id FROM masterdata_categories WHERE code = 'ORGANIZATION'
);

-- Delete category
DELETE FROM masterdata_categories WHERE code = 'ORGANIZATION';
