INSERT INTO masterdata_items (category_id, code, name, description, sort_order, is_system, status, metadata)
SELECT id, v.code, v.name, v.description, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('male',    'Male',    'Male',                                        1, '{"fhir_system": "http://hl7.org/fhir/administrative-gender", "fhir_version": "R4"}'),
    ('female',  'Female',  'Female',                                      2, '{"fhir_system": "http://hl7.org/fhir/administrative-gender", "fhir_version": "R4"}'),
    ('other',   'Other',   'Other',                                       3, '{"fhir_system": "http://hl7.org/fhir/administrative-gender", "fhir_version": "R4"}'),
    ('unknown', 'Unknown', 'A proper value is applicable, but not known', 4, '{"fhir_system": "http://hl7.org/fhir/administrative-gender", "fhir_version": "R4"}')
) AS v(code, name, description, sort_order, metadata)
WHERE c.code = 'GENDER';

INSERT INTO masterdata_items (category_id, code, name, description, sort_order, is_system, status, metadata)
SELECT id, v.code, v.name, v.description, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('A',   'Annulled',          'Marriage contract has been declared null and to not have existed',                      1,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('D',   'Divorced',          'Marriage contract has been declared dissolved and inactive',                            2,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('I',   'Interlocutory',     'Subject to an Interlocutory Decree',                                                    3,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('L',   'Legally Separated', 'Legally Separated',                                                                     4,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('M',   'Married',           'A current marriage contract is active',                                                 5,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('C',   'Common Law',        'A marriage recognized in some jurisdictions and based on the parties living together', 6,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('P',   'Polygamous',        'More than 1 current spouse',                                                            7,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('T',   'Domestic Partner',  'Person declares that a domestic partner relationship exists',                          8,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('U',   'Unmarried',         'Currently not in a marriage contract',                                                  9,  '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('S',   'Never Married',     'No marriage contract has ever been entered',                                            10, '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('W',   'Widowed',           'The spouse has died',                                                                   11, '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus", "fhir_version": "R4"}'),
    ('UNK', 'Unknown',           'A proper value is applicable, but not known',                                           12, '{"fhir_system": "http://terminology.hl7.org/CodeSystem/v3-NullFlavor", "fhir_version": "R4"}')
) AS v(code, name, description, sort_order, metadata)
WHERE c.code = 'MARITAL_STATUS';

INSERT INTO masterdata_items (category_id, code, name, description, sort_order, is_system, status, metadata)
SELECT id, v.code, v.name, v.description, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('1026', 'Islam',             'Islam',                                  1,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1013', 'Christianity',      'Christian (non-Catholic, non-specific)', 2,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1041', 'Roman Catholicism', 'Roman Catholic Church',                  3,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1008', 'Buddhism',          'Buddhism',                               4,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1020', 'Hinduism',          'Hinduism',                               5,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1027', 'Judaism',           'Judaism',                                6,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1004', 'Agnosticism',       'Agnosticism',                            7,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1007', 'Atheism',           'Atheism',                                8,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1062', 'Confucianism',      'Confucianism',                           9,  '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1046', 'Shinto',            'Shinto',                                 10, '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1047', 'Sikhism',           'Sikhism',                                11, '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1001', 'Adventist',         'Adventist',                              12, '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1038', 'Protestantism',     'Protestant',                             13, '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1036', 'Orthodox',          'Orthodox',                               14, '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1059', 'Other',             'Other',                                  98, '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}'),
    ('1061', 'Unknown',           'Unknown',                                99, '{"hl7_system": "http://terminology.hl7.org/CodeSystem/v3-ReligiousAffiliation"}')
) AS v(code, name, description, sort_order, metadata)
WHERE c.code = 'RELIGION';

INSERT INTO masterdata_items (category_id, code, name, description, sort_order, is_system, status, metadata)
SELECT id, v.code, v.name, v.description, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('0', 'Early childhood education',      'Education designed to support early development (age 0-2)',     0, '{"isced_level": 0, "standard": "ISCED 2011"}'),
    ('1', 'Primary education',              'Basic education, typically starting at age 5-7 (6 years)',      1, '{"isced_level": 1, "standard": "ISCED 2011"}'),
    ('2', 'Lower secondary education',      'First stage of secondary education (3 years)',                  2, '{"isced_level": 2, "standard": "ISCED 2011"}'),
    ('3', 'Upper secondary education',      'Second/final stage of secondary education (3 years)',           3, '{"isced_level": 3, "standard": "ISCED 2011"}'),
    ('4', 'Post-secondary non-tertiary',    'Programs bridging secondary and tertiary education',            4, '{"isced_level": 4, "standard": "ISCED 2011"}'),
    ('5', 'Short-cycle tertiary education', 'Vocational/professional higher education (2-3 years)',          5, '{"isced_level": 5, "standard": "ISCED 2011"}'),
    ('6', 'Bachelor or equivalent',         'First university degree or equivalent (3-4 years)',             6, '{"isced_level": 6, "standard": "ISCED 2011"}'),
    ('7', 'Master or equivalent',           'Second university degree or equivalent (1-2 years)',            7, '{"isced_level": 7, "standard": "ISCED 2011"}'),
    ('8', 'Doctoral or equivalent',         'Third university degree or equivalent (3+ years)',              8, '{"isced_level": 8, "standard": "ISCED 2011"}'),
    ('9', 'Not elsewhere classified',       'Education level not elsewhere classified',                      9, '{"isced_level": 9, "standard": "ISCED 2011"}')
) AS v(code, name, description, sort_order, metadata)
WHERE c.code = 'EDUCATION_LEVEL';

INSERT INTO masterdata_items (category_id, code, name, description, sort_order, is_system, status, metadata)
SELECT id, v.code, v.name, v.description, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('A+',  'A Positive',  'Blood type A, Rh positive',  1, '{"abo": "A", "rh": "positive", "can_donate_to": ["A+", "AB+"], "can_receive_from": ["A+", "A-", "O+", "O-"]}'),
    ('A-',  'A Negative',  'Blood type A, Rh negative',  2, '{"abo": "A", "rh": "negative", "can_donate_to": ["A+", "A-", "AB+", "AB-"], "can_receive_from": ["A-", "O-"]}'),
    ('B+',  'B Positive',  'Blood type B, Rh positive',  3, '{"abo": "B", "rh": "positive", "can_donate_to": ["B+", "AB+"], "can_receive_from": ["B+", "B-", "O+", "O-"]}'),
    ('B-',  'B Negative',  'Blood type B, Rh negative',  4, '{"abo": "B", "rh": "negative", "can_donate_to": ["B+", "B-", "AB+", "AB-"], "can_receive_from": ["B-", "O-"]}'),
    ('AB+', 'AB Positive', 'Blood type AB, Rh positive', 5, '{"abo": "AB", "rh": "positive", "can_donate_to": ["AB+"], "can_receive_from": ["A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"]}'),
    ('AB-', 'AB Negative', 'Blood type AB, Rh negative', 6, '{"abo": "AB", "rh": "negative", "can_donate_to": ["AB+", "AB-"], "can_receive_from": ["A-", "B-", "AB-", "O-"]}'),
    ('O+',  'O Positive',  'Blood type O, Rh positive',  7, '{"abo": "O", "rh": "positive", "can_donate_to": ["A+", "B+", "AB+", "O+"], "can_receive_from": ["O+", "O-"]}'),
    ('O-',  'O Negative',  'Blood type O, Rh negative',  8, '{"abo": "O", "rh": "negative", "can_donate_to": ["A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"], "can_receive_from": ["O-"]}')
) AS v(code, name, description, sort_order, metadata)
WHERE c.code = 'BLOOD_TYPE';

INSERT INTO masterdata_items (category_id, code, name, description, sort_order, is_system, status, metadata)
SELECT id, v.code, v.name, v.description, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('P',  'Passport',              'Machine Readable Passport (ICAO Doc 9303)',                1, '{"icao_type": "P", "standard": "ICAO Doc 9303"}'),
    ('I',  'National ID Card',      'Machine Readable National ID Card (ICAO Doc 9303 Part 5)', 2, '{"icao_type": "I", "standard": "ICAO Doc 9303"}'),
    ('D',  'Driving License',       'Driving License (ISO/IEC 18013)',                          3, '{"standard": "ISO/IEC 18013"}'),
    ('V',  'Visa',                  'Machine Readable Visa (ICAO Doc 9303 Part 7)',             4, '{"icao_type": "V", "standard": "ICAO Doc 9303"}'),
    ('T',  'Travel Document',       'Machine Readable Travel Document (ICAO Doc 9303)',         5, '{"icao_type": "T", "standard": "ICAO Doc 9303"}'),
    ('RP', 'Residence Permit',      'Residence Permit or Temporary Stay Permit',                6, '{"standard": "national"}'),
    ('WP', 'Work Permit',           'Work Permit or Employment Authorization',                  7, '{"standard": "national"}'),
    ('BC', 'Birth Certificate',     'Official Birth Certificate',                               8, '{"standard": "national"}'),
    ('MC', 'Marriage Certificate',  'Official Marriage Certificate',                            9, '{"standard": "national"}'),
    ('O',  'Other',                 'Other identity document',                                  99,'{"standard": "other"}')
) AS v(code, name, description, sort_order, metadata)
WHERE c.code = 'IDENTITY_TYPE';
