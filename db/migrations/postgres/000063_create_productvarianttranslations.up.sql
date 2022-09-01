CREATE TABLE IF NOT EXISTS productvarianttranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  productvariantid character varying(36),
  name character varying(255)
);

ALTER TABLE ONLY productvarianttranslations
    ADD CONSTRAINT productvarianttranslations_languagecode_productvariantid_key UNIQUE (languagecode, productvariantid);
