-- Create storage_records table for PostgreSQL storage abstraction
-- This table is used by the PostgreSQL storage implementation to store key-value pairs

CREATE TABLE billing.storage_records (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_storage_records_created_at ON billing.storage_records(created_at);

-- Add comments for documentation
COMMENT ON TABLE billing.storage_records IS 'Key-value storage for PostgreSQL storage abstraction';
COMMENT ON COLUMN billing.storage_records.key IS 'Unique storage key (up to 255 characters)';
COMMENT ON COLUMN billing.storage_records.value IS 'JSON-serialized storage value';
COMMENT ON COLUMN billing.storage_records.created_at IS 'Timestamp when the record was created';
COMMENT ON COLUMN billing.storage_records.updated_at IS 'Timestamp when the record was last updated';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_storage_records_updated_at 
    BEFORE UPDATE ON billing.storage_records 
    FOR EACH ROW 
    EXECUTE FUNCTION billing.update_updated_at_column();