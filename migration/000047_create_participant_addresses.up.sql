CREATE TABLE IF NOT EXISTS participant_addresses (
    -- Primary Key
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Key (intra-domain)
    participant_id      UUID NOT NULL,

    -- Address Data
    address_type        VARCHAR(50) NOT NULL,
    country_code        VARCHAR(50) NULL,
    province_code       VARCHAR(50) NULL,
    city_code           VARCHAR(50) NULL,
    district_code       VARCHAR(50) NULL,
    subdistrict_code    VARCHAR(50) NULL,
    postal_code         VARCHAR(20) NULL,
    rt                  VARCHAR(10) NULL,
    rw                  VARCHAR(10) NULL,
    address_line        TEXT NULL,

    -- Primary Flag
    is_primary          BOOLEAN NOT NULL DEFAULT FALSE,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_participant_addresses_participant FOREIGN KEY (participant_id)
        REFERENCES participants(id) ON DELETE CASCADE
);

CREATE TRIGGER trg_participant_addresses_updated_at
    BEFORE UPDATE ON participant_addresses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participant_addresses IS 'Multiple addresses per participant (home, mailing, etc.).';
COMMENT ON COLUMN participant_addresses.address_type IS 'Masterdata string code from ADDRESS_TYPE category.';
COMMENT ON COLUMN participant_addresses.country_code IS 'Masterdata string code for country.';
COMMENT ON COLUMN participant_addresses.province_code IS 'Masterdata string code for province.';
COMMENT ON COLUMN participant_addresses.city_code IS 'Masterdata string code for city.';
COMMENT ON COLUMN participant_addresses.district_code IS 'Masterdata string code for district (kecamatan).';
COMMENT ON COLUMN participant_addresses.subdistrict_code IS 'Masterdata string code for subdistrict (kelurahan).';
COMMENT ON COLUMN participant_addresses.rt IS 'RT (Rukun Tetangga) - Indonesian neighborhood unit.';
COMMENT ON COLUMN participant_addresses.rw IS 'RW (Rukun Warga) - Indonesian community unit.';
COMMENT ON COLUMN participant_addresses.is_primary IS 'Whether this is the primary address for the participant.';
