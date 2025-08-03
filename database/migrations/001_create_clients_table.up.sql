-- Create billing schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS billing;

-- Create clients table
CREATE TABLE billing.clients (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL CHECK (LENGTH(name) >= 2),
    email VARCHAR(254) NOT NULL UNIQUE,
    phone VARCHAR(20),
    address VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_clients_email ON billing.clients(email);
CREATE INDEX idx_clients_created_at ON billing.clients(created_at);
CREATE INDEX idx_clients_name ON billing.clients(name);

-- Add comments for documentation
COMMENT ON TABLE billing.clients IS 'Stores client information for billing purposes';
COMMENT ON COLUMN billing.clients.id IS 'Unique identifier for the client (UUID)';
COMMENT ON COLUMN billing.clients.name IS 'Client full name (2-100 characters)';
COMMENT ON COLUMN billing.clients.email IS 'Client email address (unique)';
COMMENT ON COLUMN billing.clients.phone IS 'Client phone number (optional)';
COMMENT ON COLUMN billing.clients.address IS 'Client address (optional, up to 500 characters)';
COMMENT ON COLUMN billing.clients.created_at IS 'Timestamp when the client was created';
COMMENT ON COLUMN billing.clients.updated_at IS 'Timestamp when the client was last updated';

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION billing.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_clients_updated_at 
    BEFORE UPDATE ON billing.clients 
    FOR EACH ROW 
    EXECUTE FUNCTION billing.update_updated_at_column();