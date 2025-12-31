CREATE TABLE IF NOT EXISTS api_hits (
  id BIGSERIAL PRIMARY KEY,
  client_id UUID NOT NULL REFERENCES clients(id),
  ip INET NOT NULL,
  endpoint TEXT NOT NULL,
  ts TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_api_hits_client_ts 
  ON api_hits (client_id, ts DESC);

CREATE INDEX IF NOT EXISTS idx_api_hits_ts 
  ON api_hits (ts DESC);
