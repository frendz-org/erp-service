-- Add index on employee_data.emp_no for lookup during self-registration.
-- The employee_data table is externally managed (dump-imported), but we
-- add this index to optimize the emp_no lookup performed on every
-- self-register request.
CREATE INDEX IF NOT EXISTS idx_employee_data_emp_no ON employee_data (emp_no);
