CREATE TABLE IF NOT EXISTS vats (
  id varchar(36) NOT NULL PRIMARY KEY,
  country_code country_code NOT NULL,
  data jsonb
);
