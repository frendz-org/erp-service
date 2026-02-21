CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL,
    product_id UUID NOT NULL,
    uploaded_by UUID NOT NULL,
    bucket VARCHAR(255) NOT NULL,
    storage_key TEXT NOT NULL,
    original_name VARCHAR(500) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    failed_delete_attempts INT NOT NULL DEFAULT 0,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_files_tenant_product ON files (tenant_id, product_id);
CREATE INDEX idx_files_expires_at ON files (expires_at) WHERE expires_at IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_files_uploaded_by ON files (uploaded_by);
CREATE INDEX idx_files_deleted_at ON files (deleted_at) WHERE deleted_at IS NOT NULL;
