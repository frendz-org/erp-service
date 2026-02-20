UPDATE masterdata_items
SET status = 'INACTIVE'
WHERE code = 'TENANT_004'
  AND deleted_at IS NULL;
