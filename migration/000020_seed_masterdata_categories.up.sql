-- ============================================================================
-- SEED: Masterdata Categories
-- Uses dynamic UUID generation via uuidv7()
-- ============================================================================

-- System Categories (standardized, not tenant-extensible)
INSERT INTO masterdata_categories (code, name, description, parent_category_id, is_system, is_tenant_extensible, sort_order) VALUES
    ('GENDER',          'Gender',          'Administrative gender (FHIR R4 AdministrativeGender)',                    NULL, TRUE, FALSE, 1),
    ('MARITAL_STATUS',  'Marital Status',  'Marital status (FHIR R4 / HL7 v3-MaritalStatus)',                          NULL, TRUE, FALSE, 2),
    ('RELIGION',        'Religion',        'Religious affiliation (HL7 v3-ReligiousAffiliation)',                      NULL, TRUE, FALSE, 3),
    ('EDUCATION_LEVEL', 'Education Level', 'Education attainment (UNESCO ISCED 2011)',                                 NULL, TRUE, FALSE, 4),
    ('BLOOD_TYPE',      'Blood Type',      'ABO and Rh blood group system',                                            NULL, TRUE, FALSE, 5),
    ('IDENTITY_TYPE',   'Identity Type',   'Identity document type (ISO/ICAO)',                                        NULL, TRUE, FALSE, 6);

-- Tenant-Extensible Categories (tenants can add their own items)
INSERT INTO masterdata_categories (code, name, description, parent_category_id, is_system, is_tenant_extensible, sort_order) VALUES
    ('DEPARTMENT',    'Department',    'Organizational department',  NULL, FALSE, TRUE, 20),
    ('JOB_LEVEL',     'Job Level',     'Employee job level/grade',   NULL, FALSE, TRUE, 21),
    ('EMPLOYEE_TYPE', 'Employee Type', 'Type of employment',         NULL, FALSE, TRUE, 22);
