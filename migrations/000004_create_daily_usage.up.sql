CREATE TABLE IF NOT EXISTS daily_usage (
  client_id UUID NOT NULL REFERENCES clients(id),
  day DATE NOT NULL,
  total BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (client_id, day)
);

CREATE INDEX IF NOT EXISTS idx_daily_usage_day 
  ON daily_usage (day DESC);
