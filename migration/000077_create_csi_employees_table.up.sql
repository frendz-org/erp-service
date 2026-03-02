-- CSI legacy employee/member master data imported from co.txt (6,250 records).
-- Standalone table — no FKs to existing domain tables.
-- Natural key csi_employee_id preserves source identity.

CREATE TABLE csi_employees (
    csi_employee_id     INTEGER PRIMARY KEY,
    employee_no         VARCHAR(50),
    employee_name       VARCHAR(255),
    gender              VARCHAR(5),
    birth_place         VARCHAR(100),
    birth_date          TIMESTAMPTZ,
    retirement_date     TIMESTAMPTZ,
    marital_status      VARCHAR(5),
    mobile_phone_no     VARCHAR(50),
    opu_no              VARCHAR(20),
    group_name          VARCHAR(50),
    status_name         VARCHAR(50),
    cost_center_no      VARCHAR(50),
    join_date           TIMESTAMPTZ,
    start_date          TIMESTAMPTZ,
    end_date            TIMESTAMPTZ,
    pension_category_no VARCHAR(10),
    amount_balance      NUMERIC(18,4),
    account_no          VARCHAR(100),
    photo               VARCHAR(255),
    is_active           BOOLEAN NOT NULL DEFAULT FALSE,
    last_updated        TIMESTAMPTZ,
    last_updater        INTEGER
);

CREATE INDEX idx_csi_employees_employee_no ON csi_employees(employee_no);
CREATE INDEX idx_csi_employees_opu_no ON csi_employees(opu_no);
CREATE INDEX idx_csi_employees_status_name ON csi_employees(status_name);
CREATE INDEX idx_csi_employees_is_active ON csi_employees(is_active);

COMMENT ON TABLE csi_employees IS 'Legacy CSI employee/member master data. Imported from DAPEN pension fund system.';
