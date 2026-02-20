INSERT INTO masterdata_categories (code, name, description, parent_category_id, is_system, is_tenant_extensible, sort_order, status, metadata) VALUES
    ('PROVINCE',
     'Province',
     'Indonesian province / administrative region (BPS code reference)',
     NULL, TRUE, FALSE, 7,
     'ACTIVE',
     '{"standard": "BPS", "country": "ID"}'::jsonb),

    ('NATIONALITY',
     'Nationality',
     'Country nationality / citizenship (ISO 3166-1 alpha-2)',
     NULL, TRUE, FALSE, 8,
     'ACTIVE',
     '{"standard": "ISO 3166-1"}'::jsonb),

    ('TAX_STATUS',
     'Tax Status',
     'Indonesian personal tax filing status (PTKP classification)',
     NULL, TRUE, FALSE, 9,
     'ACTIVE',
     '{"standard": "DJP_PTKP", "country": "ID"}'::jsonb),

    ('TERMINATION_REASON',
     'Termination Reason',
     'Employee separation / termination category per Indonesian labor law (UU Ketenagakerjaan)',
     NULL, TRUE, FALSE, 10,
     'ACTIVE',
     '{"standard": "UU_13_2003", "country": "ID"}'::jsonb)
ON CONFLICT (code) DO NOTHING;

INSERT INTO masterdata_categories (code, name, description, parent_category_id, is_system, is_tenant_extensible, sort_order, status, metadata) VALUES
    ('LEGAL_ENTITY',
     'Legal Entity',
     'Company legal entity classification',
     NULL, TRUE, FALSE, 102,
     'ACTIVE',
     '{}'::jsonb),

    ('BUSINESS_UNIT',
     'Business Unit',
     'Business unit / plant division within a legal entity',
     NULL, TRUE, FALSE, 103,
     'ACTIVE',
     '{}'::jsonb),

    ('WORK_LOCATION',
     'Work Location',
     'Employee work location / sales area',
     NULL, TRUE, TRUE, 104,
     'ACTIVE',
     '{}'::jsonb),

    ('DIVISION',
     'Division',
     'Organizational division / business group',
     NULL, FALSE, TRUE, 23,
     'ACTIVE',
     '{}'::jsonb),

    ('DEPARTMENT_FUNCTION',
     'Department Function',
     'Functional classification of departments (e.g., Finance, HR, Manufacturing)',
     NULL, FALSE, TRUE, 24,
     'ACTIVE',
     '{}'::jsonb),

    ('POSITION',
     'Position',
     'Employee job position / role title',
     NULL, FALSE, TRUE, 25,
     'ACTIVE',
     '{}'::jsonb)
ON CONFLICT (code) DO NOTHING;
