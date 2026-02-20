-- Delete tenant-extensible categories
DELETE FROM masterdata_categories WHERE code IN ('DEPARTMENT', 'JOB_LEVEL', 'EMPLOYEE_TYPE');

-- Delete system categories
DELETE FROM masterdata_categories WHERE code IN ('GENDER', 'MARITAL_STATUS', 'RELIGION', 'EDUCATION_LEVEL', 'BLOOD_TYPE', 'IDENTITY_TYPE');
