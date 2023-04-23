CREATE TABLE IF NOT EXISTS shoptranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  name character varying(110),
  description character varying(110),
  createat bigint,
  updateat bigint
);
