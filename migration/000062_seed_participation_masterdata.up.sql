-- ============================================================================
-- SEED: PARTICIPANT_PENSION_CATEGORY dan PARTICIPANT_PENSION_STATUS masterdata
-- Deskripsi: Membuat masterdata hierarkis mengikuti pola cascading dropdown:
--            PARTICIPANT_PENSION_CATEGORY (dropdown 1) → PARTICIPANT_PENSION_STATUS (dropdown 2)
--
-- Hierarki:
--   Tidak Aktif (kategori) → Tidak Aktif (status)
--   Aktif       (kategori) → Aktif, Life Cycle (status)
--   Pasif       (kategori) → Pensiun Dipercepat, Pensiun Normal, Life Cycle (status)
--
-- Catatan: "Life Cycle" muncul di bawah grup Aktif dan Pasif.
--          Kode item unik (PARTICIPANT_STATUS_NNN), nama tampil sama: "Life Cycle".
-- ============================================================================

DO $$
DECLARE
    v_category_id UUID;
    v_status_category_id UUID;
    v_tidak_aktif_item_id UUID;
    v_aktif_item_id UUID;
    v_pasif_item_id UUID;
BEGIN
    -- ========================================================================
    -- 1. Buat kategori PARTICIPANT_PENSION_CATEGORY (parent, tanpa parent)
    -- ========================================================================
    INSERT INTO masterdata_categories (
        code, name, description,
        parent_category_id,
        is_system, is_tenant_extensible,
        sort_order, status, metadata
    ) VALUES (
        'PARTICIPANT_PENSION_CATEGORY',
        'Kategori Pensiun Peserta',
        'Pengelompokan tingkat tinggi kepesertaan pensiun (Tidak Aktif/Aktif/Pasif)',
        NULL,
        TRUE,
        FALSE,
        32,
        'ACTIVE',
        '{"domain": "participant", "usage": "cascading_dropdown_parent"}'::jsonb
    )
    ON CONFLICT (code) DO NOTHING
    RETURNING id INTO v_category_id;

    IF v_category_id IS NULL THEN
        SELECT id INTO v_category_id
        FROM masterdata_categories
        WHERE code = 'PARTICIPANT_PENSION_CATEGORY' AND deleted_at IS NULL;
    END IF;

    IF v_category_id IS NULL THEN
        RAISE EXCEPTION 'Gagal membuat atau menemukan PARTICIPANT_PENSION_CATEGORY';
    END IF;

    -- ========================================================================
    -- 2. Seed item PARTICIPANT_PENSION_CATEGORY (3 grup)
    -- ========================================================================
    INSERT INTO masterdata_items (
        category_id, tenant_id, parent_item_id,
        code, name, description,
        sort_order, is_system, is_default, status, metadata
    ) VALUES
        (v_category_id, NULL, NULL,
         'PARTICIPANT_CATEGORY_001', 'Tidak Aktif',
         'Peserta tidak aktif dalam program pensiun',
         1, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb),

        (v_category_id, NULL, NULL,
         'PARTICIPANT_CATEGORY_002', 'Aktif',
         'Peserta aktif terdaftar dan berkontribusi dalam program pensiun',
         2, TRUE, TRUE, 'ACTIVE',
         '{}'::jsonb),

        (v_category_id, NULL, NULL,
         'PARTICIPANT_CATEGORY_003', 'Pasif',
         'Peserta dalam fase pasif/pensiun dari program pensiun',
         3, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb)
    ON CONFLICT (category_id, COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code)
        WHERE deleted_at IS NULL DO NOTHING;

    -- Ambil ID item grup untuk referensi parent_item_id
    SELECT id INTO v_tidak_aktif_item_id
    FROM masterdata_items
    WHERE category_id = v_category_id AND code = 'PARTICIPANT_CATEGORY_001' AND deleted_at IS NULL;

    SELECT id INTO v_aktif_item_id
    FROM masterdata_items
    WHERE category_id = v_category_id AND code = 'PARTICIPANT_CATEGORY_002' AND deleted_at IS NULL;

    SELECT id INTO v_pasif_item_id
    FROM masterdata_items
    WHERE category_id = v_category_id AND code = 'PARTICIPANT_CATEGORY_003' AND deleted_at IS NULL;

    IF v_tidak_aktif_item_id IS NULL OR v_aktif_item_id IS NULL OR v_pasif_item_id IS NULL THEN
        RAISE EXCEPTION 'Item PARTICIPANT_PENSION_CATEGORY tidak ditemukan setelah insert';
    END IF;

    -- ========================================================================
    -- 3. Buat kategori PARTICIPANT_PENSION_STATUS (child dari PARTICIPANT_PENSION_CATEGORY)
    -- ========================================================================
    INSERT INTO masterdata_categories (
        code, name, description,
        parent_category_id,
        is_system, is_tenant_extensible,
        sort_order, status, metadata
    ) VALUES (
        'PARTICIPANT_PENSION_STATUS',
        'Status Pensiun Peserta',
        'Status kepesertaan pensiun detail, cascading dari kategori pensiun',
        v_category_id,
        TRUE,
        FALSE,
        33,
        'ACTIVE',
        '{"domain": "participant", "usage": "cascading_dropdown_child"}'::jsonb
    )
    ON CONFLICT (code) DO NOTHING
    RETURNING id INTO v_status_category_id;

    IF v_status_category_id IS NULL THEN
        SELECT id INTO v_status_category_id
        FROM masterdata_categories
        WHERE code = 'PARTICIPANT_PENSION_STATUS' AND deleted_at IS NULL;
    END IF;

    IF v_status_category_id IS NULL THEN
        RAISE EXCEPTION 'Gagal membuat atau menemukan PARTICIPANT_PENSION_STATUS';
    END IF;

    -- ========================================================================
    -- 4. Seed item PARTICIPANT_PENSION_STATUS (6 status, dipetakan ke grup)
    -- ========================================================================

    -- Grup Tidak Aktif → 1 status
    INSERT INTO masterdata_items (
        category_id, tenant_id, parent_item_id,
        code, name, description,
        sort_order, is_system, is_default, status, metadata
    ) VALUES
        (v_status_category_id, NULL, v_tidak_aktif_item_id,
         'PARTICIPANT_STATUS_001', 'Tidak Aktif',
         'Peserta tidak aktif dalam program pensiun',
         1, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb);

    -- Grup Aktif → 2 status
    INSERT INTO masterdata_items (
        category_id, tenant_id, parent_item_id,
        code, name, description,
        sort_order, is_system, is_default, status, metadata
    ) VALUES
        (v_status_category_id, NULL, v_aktif_item_id,
         'PARTICIPANT_STATUS_002', 'Aktif',
         'Peserta aktif berkontribusi',
         2, TRUE, TRUE, 'ACTIVE',
         '{}'::jsonb),

        (v_status_category_id, NULL, v_aktif_item_id,
         'PARTICIPANT_STATUS_003', 'Life Cycle',
         'Peserta dalam fase life cycle di bawah kepesertaan aktif',
         3, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb);

    -- Grup Pasif → 3 status
    INSERT INTO masterdata_items (
        category_id, tenant_id, parent_item_id,
        code, name, description,
        sort_order, is_system, is_default, status, metadata
    ) VALUES
        (v_status_category_id, NULL, v_pasif_item_id,
         'PARTICIPANT_STATUS_004', 'Pensiun Dipercepat',
         'Peserta pensiun lebih awal (pensiun dipercepat)',
         4, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb),

        (v_status_category_id, NULL, v_pasif_item_id,
         'PARTICIPANT_STATUS_005', 'Pensiun Normal',
         'Peserta pensiun pada usia pensiun normal',
         5, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb),

        (v_status_category_id, NULL, v_pasif_item_id,
         'PARTICIPANT_STATUS_006', 'Life Cycle',
         'Peserta dalam fase life cycle di bawah status pasif/pensiun',
         6, TRUE, FALSE, 'ACTIVE',
         '{}'::jsonb);

    RAISE NOTICE 'Seed PARTICIPANT_PENSION_CATEGORY (3 grup) dan PARTICIPANT_PENSION_STATUS (6 status) berhasil';
END $$;
