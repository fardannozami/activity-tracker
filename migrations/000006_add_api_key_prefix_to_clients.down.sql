DROP INDEX IF EXISTS idx_clients_api_key_prefix;

ALTER TABLE clients
DROP COLUMN IF EXISTS api_key_prefix;