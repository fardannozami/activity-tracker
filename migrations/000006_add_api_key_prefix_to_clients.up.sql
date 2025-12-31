ALTER TABLE clients
ADD COLUMN IF NOT EXISTS api_key_prefix TEXT;

CREATE INDEX IF NOT EXISTS idx_clients_api_key_prefix
ON clients(api_key_prefix);