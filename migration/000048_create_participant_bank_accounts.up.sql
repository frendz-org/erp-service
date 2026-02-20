CREATE TABLE IF NOT EXISTS participant_bank_accounts (
    -- Primary Key
    id                      UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Key (intra-domain)
    participant_id          UUID NOT NULL,

    -- Bank Account Data
    bank_code               VARCHAR(50) NOT NULL,
    account_number          VARCHAR(50) NOT NULL,
    account_holder_name     VARCHAR(255) NOT NULL,
    account_type            VARCHAR(50) NULL,
    currency_code           VARCHAR(50) NOT NULL DEFAULT 'IDR',

    -- Primary Flag
    is_primary              BOOLEAN NOT NULL DEFAULT FALSE,

    -- Validity Dates
    issue_date              DATE NULL,
    expiry_date             DATE NULL,

    -- Audit Fields
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ,

    -- Optimistic Locking
    version                 INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_participant_bank_accounts_participant FOREIGN KEY (participant_id)
        REFERENCES participants(id) ON DELETE CASCADE
);

CREATE TRIGGER trg_participant_bank_accounts_updated_at
    BEFORE UPDATE ON participant_bank_accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE participant_bank_accounts IS 'Bank accounts for a participant. One primary per participant enforced by partial unique index.';
COMMENT ON COLUMN participant_bank_accounts.bank_code IS 'Masterdata string code from BANK category.';
COMMENT ON COLUMN participant_bank_accounts.account_type IS 'Masterdata string code from ACCOUNT_TYPE category.';
COMMENT ON COLUMN participant_bank_accounts.currency_code IS 'Masterdata string code from CURRENCY category. Defaults to IDR.';
COMMENT ON COLUMN participant_bank_accounts.is_primary IS 'Whether this is the primary bank account. Uniqueness enforced by partial index.';
