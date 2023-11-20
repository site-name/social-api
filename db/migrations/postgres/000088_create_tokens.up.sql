CREATE TABLE IF NOT EXISTS tokens (
  token varchar(64) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  type varchar(64) NOT NULL,
  extra varchar(2048)
);