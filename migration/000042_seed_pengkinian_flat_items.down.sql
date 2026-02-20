DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'PROVINCE')
  AND code LIKE 'PROVINCE_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'NATIONALITY')
  AND code LIKE 'NATIONALITY_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'TAX_STATUS')
  AND code LIKE 'TAX_STATUS_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'TERMINATION_REASON')
  AND code LIKE 'TERMINATION_REASON_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'JOB_LEVEL')
  AND code LIKE 'JOB_LEVEL_%';
DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'EMPLOYEE_TYPE')
  AND code LIKE 'EMPLOYEE_TYPE_%';
