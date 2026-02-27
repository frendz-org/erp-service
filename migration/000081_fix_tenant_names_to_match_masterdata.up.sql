UPDATE tenants SET name = 'DPIP Bogasari', tenant_type = 'TENANT_TYPE_002', updated_at = NOW()
WHERE code = 'TENANT_001' AND deleted_at IS NULL;

UPDATE tenants SET name = 'PT ISM Bogasari Flour Mills', tenant_type = 'TENANT_TYPE_001', updated_at = NOW()
WHERE code = 'TENANT_002' AND deleted_at IS NULL;

UPDATE tenants SET name = 'DPMP Bogasari', tenant_type = 'TENANT_TYPE_002', updated_at = NOW()
WHERE code = 'TENANT_003' AND deleted_at IS NULL;

UPDATE tenants SET name = 'PT Inti Abadi Kemasindo', tenant_type = 'TENANT_TYPE_001', updated_at = NOW()
WHERE code = 'TENANT_004' AND deleted_at IS NULL;
