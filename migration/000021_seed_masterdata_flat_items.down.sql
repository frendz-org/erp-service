-- Delete items by category code (reverse order of insertion)

-- Delete Identity Type items
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'IDENTITY_TYPE');

-- Delete Blood Type items
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'BLOOD_TYPE');

-- Delete Education Level items
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'EDUCATION_LEVEL');

-- Delete Religion items
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'RELIGION');

-- Delete Marital Status items
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'MARITAL_STATUS');

-- Delete Gender items
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'GENDER');
