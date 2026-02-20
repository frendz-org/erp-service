DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'LEGAL_ENTITY')
  AND code LIKE 'LEGAL_ENTITY_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'BUSINESS_UNIT')
  AND code LIKE 'BUSINESS_UNIT_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'WORK_LOCATION')
  AND code LIKE 'WORK_LOCATION_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'DIVISION')
  AND code LIKE 'DIVISION_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'DEPARTMENT_FUNCTION')
  AND code LIKE 'DEPARTMENT_FUNCTION_%';
