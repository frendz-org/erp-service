-- CSI legacy monthly ledger transactions imported from ledger0126v2.txt (2,454,551 records).
-- References csi_employees for referential integrity within CSI domain.
-- amount_trans is NUMERIC(18,4) because some item types carry decimal precision.

CREATE TABLE csi_ledgers (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    csi_employee_id     INTEGER NOT NULL REFERENCES csi_employees(csi_employee_id),
    year_period         INTEGER NOT NULL,
    month_period        INTEGER NOT NULL,
    csi_item_id         INTEGER NOT NULL,
    amount_trans        NUMERIC(18,4) NOT NULL,
    last_updated        TIMESTAMPTZ,
    last_updater        INTEGER
);

CREATE INDEX idx_csi_ledgers_employee_id ON csi_ledgers(csi_employee_id);
CREATE INDEX idx_csi_ledgers_period ON csi_ledgers(year_period, month_period);
CREATE INDEX idx_csi_ledgers_item ON csi_ledgers(csi_item_id);
CREATE UNIQUE INDEX uk_csi_ledgers_employee_period_item
    ON csi_ledgers(csi_employee_id, year_period, month_period, csi_item_id);

COMMENT ON TABLE csi_ledgers IS 'Legacy CSI monthly ledger transactions. Imported from DAPEN pension fund system.';
