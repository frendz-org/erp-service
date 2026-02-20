UPDATE masterdata_items
SET status = 'ACTIVE'
WHERE code = 'TENANT_004'
  AND deleted_at IS NULL;
