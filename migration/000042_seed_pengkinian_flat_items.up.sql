INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('PROVINCE_001', 'Aceh',                          1,  '{"bps_code": "11", "legacy_sap_code": "06", "iso_3166_2": "ID-AC"}'),
    ('PROVINCE_002', 'Sumatera Utara',                2,  '{"bps_code": "12", "legacy_sap_code": "07", "iso_3166_2": "ID-SU"}'),
    ('PROVINCE_003', 'Sumatera Barat',                3,  '{"bps_code": "13", "legacy_sap_code": "08", "iso_3166_2": "ID-SB"}'),
    ('PROVINCE_004', 'Riau',                          4,  '{"bps_code": "14", "legacy_sap_code": "09", "iso_3166_2": "ID-RI"}'),
    ('PROVINCE_005', 'Jambi',                         5,  '{"bps_code": "15", "legacy_sap_code": "10", "iso_3166_2": "ID-JA"}'),
    ('PROVINCE_006', 'Sumatera Selatan',              6,  '{"bps_code": "16", "legacy_sap_code": "11", "iso_3166_2": "ID-SS"}'),
    ('PROVINCE_007', 'Bengkulu',                      7,  '{"bps_code": "17", "legacy_sap_code": "12", "iso_3166_2": "ID-BE"}'),
    ('PROVINCE_008', 'Lampung',                       8,  '{"bps_code": "18", "legacy_sap_code": "13", "iso_3166_2": "ID-LA"}'),
    ('PROVINCE_009', 'Kepulauan Bangka Belitung',     9,  '{"bps_code": "19", "legacy_sap_code": "29", "iso_3166_2": "ID-BB"}'),
    ('PROVINCE_010', 'Kepulauan Riau',                10, '{"bps_code": "21", "legacy_sap_code": null, "iso_3166_2": "ID-KR"}'),
    ('PROVINCE_011', 'DKI Jakarta',                   11, '{"bps_code": "31", "legacy_sap_code": "01", "iso_3166_2": "ID-JK"}'),
    ('PROVINCE_012', 'Jawa Barat',                    12, '{"bps_code": "32", "legacy_sap_code": "02", "iso_3166_2": "ID-JB"}'),
    ('PROVINCE_013', 'Jawa Tengah',                   13, '{"bps_code": "33", "legacy_sap_code": "03", "iso_3166_2": "ID-JT"}'),
    ('PROVINCE_014', 'DI Yogyakarta',                 14, '{"bps_code": "34", "legacy_sap_code": "05", "iso_3166_2": "ID-YO"}'),
    ('PROVINCE_015', 'Jawa Timur',                    15, '{"bps_code": "35", "legacy_sap_code": "04", "iso_3166_2": "ID-JI"}'),
    ('PROVINCE_016', 'Banten',                        16, '{"bps_code": "36", "legacy_sap_code": "27", "iso_3166_2": "ID-BT"}'),
    ('PROVINCE_017', 'Bali',                          17, '{"bps_code": "51", "legacy_sap_code": "22", "iso_3166_2": "ID-BA"}'),
    ('PROVINCE_018', 'Nusa Tenggara Barat',           18, '{"bps_code": "52", "legacy_sap_code": "23", "iso_3166_2": "ID-NB"}'),
    ('PROVINCE_019', 'Nusa Tenggara Timur',           19, '{"bps_code": "53", "legacy_sap_code": "24", "iso_3166_2": "ID-NT"}'),
    ('PROVINCE_020', 'Kalimantan Barat',              20, '{"bps_code": "61", "legacy_sap_code": "15", "iso_3166_2": "ID-KB"}'),
    ('PROVINCE_021', 'Kalimantan Tengah',             21, '{"bps_code": "62", "legacy_sap_code": null, "iso_3166_2": "ID-KT"}'),
    ('PROVINCE_022', 'Kalimantan Selatan',            22, '{"bps_code": "63", "legacy_sap_code": "14", "iso_3166_2": "ID-KS"}'),
    ('PROVINCE_023', 'Kalimantan Timur',              23, '{"bps_code": "64", "legacy_sap_code": "17", "iso_3166_2": "ID-KI"}'),
    ('PROVINCE_024', 'Kalimantan Utara',              24, '{"bps_code": "65", "legacy_sap_code": null, "iso_3166_2": "ID-KU"}'),
    ('PROVINCE_025', 'Sulawesi Utara',                25, '{"bps_code": "71", "legacy_sap_code": "21", "iso_3166_2": "ID-SA"}'),
    ('PROVINCE_026', 'Sulawesi Tengah',               26, '{"bps_code": "72", "legacy_sap_code": "20", "iso_3166_2": "ID-ST"}'),
    ('PROVINCE_027', 'Sulawesi Selatan',              27, '{"bps_code": "73", "legacy_sap_code": "18", "iso_3166_2": "ID-SN"}'),
    ('PROVINCE_028', 'Sulawesi Tenggara',             28, '{"bps_code": "74", "legacy_sap_code": "19", "iso_3166_2": "ID-SG"}'),
    ('PROVINCE_029', 'Gorontalo',                     29, '{"bps_code": "75", "legacy_sap_code": null, "iso_3166_2": "ID-GO"}'),
    ('PROVINCE_030', 'Sulawesi Barat',                30, '{"bps_code": "76", "legacy_sap_code": null, "iso_3166_2": "ID-SR"}'),
    ('PROVINCE_031', 'Maluku',                        31, '{"bps_code": "81", "legacy_sap_code": "25", "iso_3166_2": "ID-ML"}'),
    ('PROVINCE_032', 'Maluku Utara',                  32, '{"bps_code": "82", "legacy_sap_code": null, "iso_3166_2": "ID-MU"}'),
    ('PROVINCE_033', 'Papua',                         33, '{"bps_code": "91", "legacy_sap_code": "33", "iso_3166_2": "ID-PA"}'),
    ('PROVINCE_034', 'Papua Barat',                   34, '{"bps_code": "92", "legacy_sap_code": null, "iso_3166_2": "ID-PB"}'),
    ('PROVINCE_035', 'Papua Selatan',                 35, '{"bps_code": "93", "legacy_sap_code": null, "iso_3166_2": "ID-PS"}'),
    ('PROVINCE_036', 'Papua Tengah',                  36, '{"bps_code": "94", "legacy_sap_code": null, "iso_3166_2": "ID-PT"}'),
    ('PROVINCE_037', 'Papua Pegunungan',              37, '{"bps_code": "95", "legacy_sap_code": null, "iso_3166_2": "ID-PP"}'),
    ('PROVINCE_038', 'Papua Barat Daya',              38, '{"bps_code": "96", "legacy_sap_code": null, "iso_3166_2": "ID-PD"}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'PROVINCE'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('NATIONALITY_001', 'Indonesian', 1, '{"iso_alpha2": "ID", "iso_alpha3": "IDN", "iso_numeric": "360", "legacy_sap_code": "ID"}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'NATIONALITY'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('TAX_STATUS_001', 'TK/0 - Tidak Kawin Tanpa Tanggungan',  1, '{"legacy_sap_code": "TK-0", "ptkp_code": "TK/0", "marital": "unmarried", "dependents": 0}'),
    ('TAX_STATUS_002', 'TK/1 - Tidak Kawin 1 Tanggungan',      2, '{"legacy_sap_code": "TK-1", "ptkp_code": "TK/1", "marital": "unmarried", "dependents": 1}'),
    ('TAX_STATUS_003', 'TK/2 - Tidak Kawin 2 Tanggungan',      3, '{"legacy_sap_code": "TK-2", "ptkp_code": "TK/2", "marital": "unmarried", "dependents": 2}'),
    ('TAX_STATUS_004', 'TK/3 - Tidak Kawin 3 Tanggungan',      4, '{"legacy_sap_code": "TK-3", "ptkp_code": "TK/3", "marital": "unmarried", "dependents": 3}'),
    ('TAX_STATUS_005', 'K/0 - Kawin Tanpa Tanggungan',          5, '{"legacy_sap_code": "K-0",  "ptkp_code": "K/0",  "marital": "married",   "dependents": 0}'),
    ('TAX_STATUS_006', 'K/1 - Kawin 1 Tanggungan',              6, '{"legacy_sap_code": "K-1",  "ptkp_code": "K/1",  "marital": "married",   "dependents": 1}'),
    ('TAX_STATUS_007', 'K/2 - Kawin 2 Tanggungan',              7, '{"legacy_sap_code": "K-2",  "ptkp_code": "K/2",  "marital": "married",   "dependents": 2}'),
    ('TAX_STATUS_008', 'K/3 - Kawin 3 Tanggungan',              8, '{"legacy_sap_code": "K-3",  "ptkp_code": "K/3",  "marital": "married",   "dependents": 3}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'TAX_STATUS'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, description, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.description, v.sort_order, TRUE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('TERMINATION_REASON_001', 'Mengundurkan Diri (Resign)',           'Permohonan sukarela, surat resign min. 30 hari',                          1,  '{"legacy_sap_code": "RESIGN",              "initiative": "employee",  "hris_code": "Resign"}'),
    ('TERMINATION_REASON_002', 'PKWT Habis Kontrak',                   'Masa kontrak berakhir sesuai tanggal PKWT',                                2,  '{"legacy_sap_code": "CONTRACT_END",        "initiative": "automatic", "hris_code": "Contract End"}'),
    ('TERMINATION_REASON_003', 'PHK Efisiensi',                        'Efisiensi / pengurangan biaya, dasar bisnis terdokumentasi',               3,  '{"legacy_sap_code": "EFFICIENCY",          "initiative": "company",   "hris_code": "Efficiency PHK"}'),
    ('TERMINATION_REASON_004', 'PHK Restrukturisasi',                  'Perubahan organisasi, dokumen restrukturisasi',                             4,  '{"legacy_sap_code": "RESTRUCTURING",       "initiative": "company",   "hris_code": "Restructuring"}'),
    ('TERMINATION_REASON_005', 'PHK Kinerja',                          'Performa tidak memenuhi, ada evaluasi & SP1-SP3',                          5,  '{"legacy_sap_code": "PERFORMANCE",         "initiative": "company",   "hris_code": "Performance PHK"}'),
    ('TERMINATION_REASON_006', 'PHK Pelanggaran Disiplin',             'Pelanggaran berat, ada investigasi & bukti',                               6,  '{"legacy_sap_code": "DISCIPLINE",          "initiative": "company",   "hris_code": "Discipline PHK"}'),
    ('TERMINATION_REASON_007', 'PHK Mangkir',                          'Mangkir berturut-turut, pemanggilan resmi',                                7,  '{"legacy_sap_code": "ABSENTEEISM",         "initiative": "company",   "hris_code": "Absenteeism PHK"}'),
    ('TERMINATION_REASON_008', 'Kesepakatan Bersama',                  'Mutual agreement, perjanjian tertulis',                                    8,  '{"legacy_sap_code": "MUTUAL_AGREEMENT",    "initiative": "both",      "hris_code": "Mutual Agreement"}'),
    ('TERMINATION_REASON_009', 'Sakit Berkepanjangan',                 'Tidak mampu kerja, bukti medis',                                           9,  '{"legacy_sap_code": "PROLONGED_ILLNESS",   "initiative": "company",   "hris_code": "Prolonged Illness"}'),
    ('TERMINATION_REASON_010', 'Disabilitas Permanen',                 'Tidak dapat lanjut kerja, bukti medis',                                    10, '{"legacy_sap_code": "PERMANENT_DISABILITY","initiative": "company",   "hris_code": "Permanent Disability"}'),
    ('TERMINATION_REASON_011', 'Pensiun',                              'Usia pensiun sesuai aturan perusahaan',                                    11, '{"legacy_sap_code": "RETIREMENT",          "initiative": "automatic", "hris_code": "Retirement"}'),
    ('TERMINATION_REASON_012', 'Meninggal Dunia',                      'Wafat, dokumen kematian, hak dibayar ke ahli waris',                       12, '{"legacy_sap_code": "DECEASED",            "initiative": "natural",   "hris_code": "Deceased"}'),
    ('TERMINATION_REASON_013', 'Perusahaan Tutup',                     'Penutupan usaha, legal closure',                                           13, '{"legacy_sap_code": "COMPANY_CLOSURE",     "initiative": "company",   "hris_code": "Company Closure"}'),
    ('TERMINATION_REASON_014', 'Pailit',                               'Putusan pengadilan, status pailit resmi',                                  14, '{"legacy_sap_code": "BANKRUPTCY",          "initiative": "legal",     "hris_code": "Bankruptcy Termination"}'),
    ('TERMINATION_REASON_015', 'Force Majeure',                        'Keadaan kahar, bukti kondisi force majeure',                               15, '{"legacy_sap_code": "FORCE_MAJEURE",       "initiative": "company",   "hris_code": "Force Majeure"}'),
    ('TERMINATION_REASON_016', 'Akhir Masa Percobaan Tidak Lulus',     'Tidak lolos probation, evaluasi probation',                                16, '{"legacy_sap_code": "PROBATION_FAIL",      "initiative": "company",   "hris_code": "Probation Fail"}'),
    ('TERMINATION_REASON_017', 'Pemutusan PKWT Sebelum Waktu',        'Kontrak diputus dini, alasan sah, kompensasi sisa kontrak',                17, '{"legacy_sap_code": "EARLY_CONTRACT_END",  "initiative": "company",   "hris_code": "Early Contract End"}')
) AS v(code, name, description, sort_order, metadata)
WHERE c.code = 'TERMINATION_REASON'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.sort_order, FALSE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('JOB_LEVEL_001', 'Operator / Teknisi',    1, '{"legacy_sap_code": "Operator/Teknisi"}'),
    ('JOB_LEVEL_002', 'Foreman / Supervisor',   2, '{"legacy_sap_code": "Foreman/Supervisor"}'),
    ('JOB_LEVEL_003', 'Section Head',           3, '{"legacy_sap_code": "Section Head"}'),
    ('JOB_LEVEL_004', 'Assistant Manager',      4, '{"legacy_sap_code": "Ass.Manager"}'),
    ('JOB_LEVEL_005', 'Manager',                5, '{"legacy_sap_code": "Manager"}'),
    ('JOB_LEVEL_006', 'Vice President',         6, '{"legacy_sap_code": "Vice President"}'),
    ('JOB_LEVEL_007', 'Senior Vice President',  7, '{"legacy_sap_code": "Senior VP"}')
) AS v(code, name, sort_order, metadata)
WHERE c.code = 'JOB_LEVEL'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;

INSERT INTO masterdata_items (category_id, code, name, alt_name, sort_order, is_system, status, metadata)
SELECT c.id, v.code, v.name, v.alt_name, v.sort_order, FALSE, 'ACTIVE', v.metadata::jsonb
FROM masterdata_categories c
CROSS JOIN (VALUES
    ('EMPLOYEE_TYPE_001', 'Permanent',     'Karyawan Tetap (PKWTT)',      1,  '{"legacy_sap_code": "Permanent"}'),
    ('EMPLOYEE_TYPE_002', 'Contract',      'Karyawan Kontrak (PKWT)',     2,  '{"legacy_sap_code": "Contract"}'),
    ('EMPLOYEE_TYPE_003', 'Probation',     'Masa Percobaan',              3,  '{"legacy_sap_code": "Probation"}'),
    ('EMPLOYEE_TYPE_004', 'Internship',    'Magang',                      4,  '{"legacy_sap_code": "Internship"}'),
    ('EMPLOYEE_TYPE_005', 'Temporary',     'Sementara',                   5,  '{"legacy_sap_code": "Temporary"}'),
    ('EMPLOYEE_TYPE_006', 'Outsource',     'Tenaga Alih Daya',            6,  '{"legacy_sap_code": "Outsource"}'),
    ('EMPLOYEE_TYPE_007', 'Part-Time',     'Paruh Waktu',                 7,  '{"legacy_sap_code": "Part-Time"}'),
    ('EMPLOYEE_TYPE_008', 'Full-Time',     'Penuh Waktu',                 8,  '{"legacy_sap_code": "Full-Time"}'),
    ('EMPLOYEE_TYPE_009', 'Freelance',     'Lepas',                       9,  '{"legacy_sap_code": "Freelance"}'),
    ('EMPLOYEE_TYPE_010', 'Consultant',    'Konsultan',                   10, '{"legacy_sap_code": "Consultant"}'),
    ('EMPLOYEE_TYPE_011', 'Daily Worker',  'Harian Lepas',                11, '{"legacy_sap_code": "Daily Worker"}')
) AS v(code, name, alt_name, sort_order, metadata)
WHERE c.code = 'EMPLOYEE_TYPE'
ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
    WHERE deleted_at IS NULL DO NOTHING;
