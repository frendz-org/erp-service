-- Enforce uniqueness of employee_number within a tenant+product scope.
-- A participant must have a unique employee number per tenant/product combination.
CREATE UNIQUE INDEX IF NOT EXISTS idx_participants_tenant_product_employee_number
    ON participants (tenant_id, product_id, employee_number)
    WHERE employee_number IS NOT NULL
      AND deleted_at IS NULL;
