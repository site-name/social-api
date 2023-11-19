CREATE TABLE IF NOT EXISTS vats (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  country_code VARCHAR(5) NOT NULL,
  data jsonb
)
