CREATE TABLE IF NOT EXISTS hourly_usage (
  client_id UUID NOT NULL REFERENCES clients(id),
  hour TIMESTAMPTZ NOT NULL,
  total BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (client_id, hour)
);

CREATE INDEX IF NOT EXISTS idx_hourly_usage_hour 
  ON hourly_usage (hour DESC);
