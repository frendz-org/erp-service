-- Delete items first (FK constraint)
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'TENANT_TYPE');

-- Delete category
DELETE FROM masterdata_categories WHERE code = 'TENANT_TYPE';
