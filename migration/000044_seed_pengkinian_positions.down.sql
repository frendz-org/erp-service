DELETE FROM masterdata_items
WHERE category_id = (SELECT id FROM masterdata_categories WHERE code = 'POSITION')
  AND code LIKE 'POSITION_%';
