

INSERT INTO tenants (code, name, tenant_type, status, settings)
VALUES
    ('ISM-BOGASARI', 'PT ISM Bogasari Flour Mills', 'MITRA_PENDIRI', 'ACTIVE',
     '{"approval_required": true, "password_policy": {"min_length": 8, "require_special": true}}'::jsonb),

    ('INTI-ABADI-KEMASINDO', 'PT Inti Abadi Kemasindo', 'MITRA_PENDIRI', 'ACTIVE',
     '{"approval_required": true, "password_policy": {"min_length": 8, "require_special": true}}'::jsonb),

    ('DPIP-BOGASARI', 'DPIP Bogasari', 'PENSION_FUND', 'ACTIVE',
     '{"approval_required": true, "password_policy": {"min_length": 8, "require_special": true}}'::jsonb),

    ('DPMP-BOGASARI', 'DPMP Bogasari', 'PENSION_FUND', 'ACTIVE',
     '{"approval_required": true, "password_policy": {"min_length": 8, "require_special": true}}'::jsonb)

ON CONFLICT (code) DO NOTHING;
COMMENT ON TABLE tenants IS 'Root table for multi-tenant organization management. Codes should match masterdata ORGANIZATION items.';
