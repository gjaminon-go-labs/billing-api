-- Drop trigger first
DROP TRIGGER IF EXISTS update_clients_updated_at ON billing.clients;

-- Drop function
DROP FUNCTION IF EXISTS billing.update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS billing.idx_clients_name;
DROP INDEX IF EXISTS billing.idx_clients_created_at;
DROP INDEX IF EXISTS billing.idx_clients_email;

-- Drop table
DROP TABLE IF EXISTS billing.clients;

-- Note: We don't drop the schema as it might be used by other tables
-- DROP SCHEMA IF EXISTS billing CASCADE;