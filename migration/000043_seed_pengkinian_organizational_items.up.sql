INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('LEGAL_ENTITY_001', 'Indofood Sukses Makmur',  1, '{"legacy_sap_code": "1100", "company_type": "PT"}'),
    ('LEGAL_ENTITY_002', 'Inti Abadi Kemasindo',     2, '{"legacy_sap_code": "3300", "company_type": "PT"}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'LEGAL_ENTITY'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('BUSINESS_UNIT_001', 'Bogasari - Corporate',       1, '{"legacy_sap_code": "1251", "legal_entity_code": "1100"}'),
    ('BUSINESS_UNIT_002', 'Bogasari - Flour Jakarta',    2, '{"legacy_sap_code": "1252", "legal_entity_code": "1100"}'),
    ('BUSINESS_UNIT_003', 'Bogasari - Flour Surabaya',   3, '{"legacy_sap_code": "1253", "legal_entity_code": "1100"}'),
    ('BUSINESS_UNIT_004', 'Bogasari - Pasta',            4, '{"legacy_sap_code": "1254", "legal_entity_code": "1100"}'),
    ('BUSINESS_UNIT_005', 'Bogasari - Flour Cibitung',   5, '{"legacy_sap_code": "1256", "legal_entity_code": "1100"}'),
    ('BUSINESS_UNIT_006', 'IAK Packaging',               6, '{"legacy_sap_code": "3302", "legal_entity_code": "3300"}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'BUSINESS_UNIT'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('WORK_LOCATION_001', 'BS-Jakarta (Corporate)',      1,  '{"legacy_sap_code": "BS-Jakarta",        "location_code": "K001", "business_unit_code": "1251"}'),

    ('WORK_LOCATION_002', 'BS-Aceh',                     2,  '{"legacy_sap_code": "BS-Aceh",           "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_003', 'BS-Bandung',                  3,  '{"legacy_sap_code": "BS-Bandung",        "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_004', 'BS-Bogor',                    4,  '{"legacy_sap_code": "BS-Bogor",          "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_005', 'BS-Cirebon',                  5,  '{"legacy_sap_code": "BS-Cirebon",        "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_006', 'BS-Garut',                    6,  '{"legacy_sap_code": "BS-Garut",          "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_007', 'BS-Jakarta (Flour)',           7,  '{"legacy_sap_code": "BS-Jakarta",        "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_008', 'BS-Jambi UM',                 8,  '{"legacy_sap_code": "BS-Jambi UM",       "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_009', 'BS-Jatabek',                  9,  '{"legacy_sap_code": "BS-Jatabek",        "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_010', 'BS-Kalimantan Barat',         10, '{"legacy_sap_code": "BS-Kalimantan Barat","location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_011', 'BS-Karawang',                 11, '{"legacy_sap_code": "BS-Karawang",       "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_012', 'BS-Lampung',                  12, '{"legacy_sap_code": "BS-Lampung",        "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_013', 'BS-Serang',                   13, '{"legacy_sap_code": "BS-Serang",         "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_014', 'BS-Sukabumi',                 14, '{"legacy_sap_code": "BS-Sukabumi",       "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_015', 'BS-Sumatera Barat',           15, '{"legacy_sap_code": "BS-Sumatera Barat", "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_016', 'BS-Sumatera Selatan',         16, '{"legacy_sap_code": "BS-Sumatera Selatan","location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_017', 'BS-Sumatera Utara',           17, '{"legacy_sap_code": "BS-Sumatera Utara", "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_018', 'BS-Tasikmalaya',              18, '{"legacy_sap_code": "BS-Tasikmalaya",    "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_019', 'BS-Tegal',                    19, '{"legacy_sap_code": "BS-Tegal",          "location_code": "K002", "business_unit_code": "1252"}'),
    ('WORK_LOCATION_020', 'Batam',                       20, '{"legacy_sap_code": "Batam",             "location_code": "K002", "business_unit_code": "1252"}'),

    ('WORK_LOCATION_021', 'BS-Banjarmasin',              21, '{"legacy_sap_code": "BS-Banjarmasin",    "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_022', 'BS-Blitar',                   22, '{"legacy_sap_code": "BS-Blitar",         "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_023', 'BS-Denpasar',                 23, '{"legacy_sap_code": "BS-Denpasar",       "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_024', 'BS-Jember',                   24, '{"legacy_sap_code": "BS-Jember",         "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_025', 'BS-Kediri',                   25, '{"legacy_sap_code": "BS-Kediri",         "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_026', 'BS-Lamongan',                 26, '{"legacy_sap_code": "BS-Lamongan",       "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_027', 'BS-Madiun',                   27, '{"legacy_sap_code": "BS-Madiun",         "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_028', 'BS-Malang',                   28, '{"legacy_sap_code": "BS-Malang",         "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_029', 'BS-Pati',                     29, '{"legacy_sap_code": "BS-Pati",           "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_030', 'BS-Probolinggo',              30, '{"legacy_sap_code": "BS-Probolinggo",    "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_031', 'BS-Semarang',                 31, '{"legacy_sap_code": "BS-Semarang",       "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_032', 'BS-Situbondo',                32, '{"legacy_sap_code": "BS-Situbondo",      "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_033', 'BS-Solo',                     33, '{"legacy_sap_code": "BS-Solo",           "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_034', 'BS-Surabaya',                 34, '{"legacy_sap_code": "BS-Surabaya",       "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_035', 'BS-Yogyakarta',               35, '{"legacy_sap_code": "BS-Yogyakarta",     "location_code": "K003", "business_unit_code": "1253"}'),
    ('WORK_LOCATION_036', 'Lumajang',                    36, '{"legacy_sap_code": "Lumajang",          "location_code": "K003", "business_unit_code": "1253"}'),

    ('WORK_LOCATION_037', 'BS-Jakarta (Pasta)',           37, '{"legacy_sap_code": "BS-Jakarta",        "location_code": "K004", "business_unit_code": "1254"}'),

    ('WORK_LOCATION_038', 'BS-Cibitung',                 38, '{"legacy_sap_code": "BS-Cibitung",       "location_code": "K005", "business_unit_code": "1256"}'),

    ('WORK_LOCATION_039', 'BS-Jakarta (K006)',            39, '{"legacy_sap_code": "BS-Jakarta",        "location_code": "K006", "business_unit_code": null}'),

    ('WORK_LOCATION_040', 'IAK-Citeureup',               40, '{"legacy_sap_code": "IAK-Citeureup",     "location_code": "L001", "business_unit_code": "3302"}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'WORK_LOCATION'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, FALSE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('DIVISION_001', 'Bogasari',               1, '{"legacy_sap_code": "BGS"}'),
    ('DIVISION_002', 'Inti Abadi Kemasindo',    2, '{"legacy_sap_code": "IAK"}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'DIVISION'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, FALSE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('DEPARTMENT_FUNCTION_001', 'Administration - Finance',         1, '{"legacy_sap_code": "ADM Finance"}'),
    ('DEPARTMENT_FUNCTION_002', 'Administration - Human Resources', 2, '{"legacy_sap_code": "ADM HR"}'),
    ('DEPARTMENT_FUNCTION_003', 'Administration - Management',      3, '{"legacy_sap_code": "ADM Man.Office"}'),
    ('DEPARTMENT_FUNCTION_004', 'Manufacturing',                    4, '{"legacy_sap_code": "MFG Manuf."}'),
    ('DEPARTMENT_FUNCTION_005', 'Marketing',                        5, '{"legacy_sap_code": "MKT Marketing"}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'DEPARTMENT_FUNCTION'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;
