CREATE TABLE IF NOT EXISTS vats (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  countrycode VARCHAR(5),
  data jsonb
)
