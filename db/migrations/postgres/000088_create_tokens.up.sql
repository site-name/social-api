CREATE TABLE IF NOT EXISTS tokens (
  token character varying(64) NOT NULL PRIMARY KEY,
  createat bigint,
  type character varying(64),
  extra character varying(2048)
);