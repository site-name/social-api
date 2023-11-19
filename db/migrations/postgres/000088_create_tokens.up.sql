CREATE TABLE IF NOT EXISTS tokens (
  token character varying(64) NOT NULL PRIMARY KEY,
  created_at bigint NOT NULL,
  type character varying(64) NOT NULL,
  extra character varying(2048)
);