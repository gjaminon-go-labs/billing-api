-- Drop trigger first
DROP TRIGGER IF EXISTS update_storage_records_updated_at ON billing.storage_records;

-- Drop indexes
DROP INDEX IF EXISTS billing.idx_storage_records_created_at;

-- Drop table
DROP TABLE IF EXISTS billing.storage_records;