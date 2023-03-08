CREATE TABLE IF NOT EXISTS variantmedias (
  id character varying(36) NOT NULL PRIMARY KEY,
  variantid character varying(36),
  mediaid character varying(36)
);

ALTER TABLE ONLY variantmedias
    ADD CONSTRAINT variantmedias_variantid_mediaid_key UNIQUE (variantid, mediaid);
