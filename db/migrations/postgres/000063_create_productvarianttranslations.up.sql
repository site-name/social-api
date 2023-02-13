CREATE TABLE IF NOT EXISTS productvarianttranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  productvariantid character varying(36),
  name character varying(255)
);

ALTER TABLE ONLY productvarianttranslations
    ADD CONSTRAINT productvarianttranslations_languagecode_productvariantid_key UNIQUE (languagecode, productvariantid);
ALTER TABLE ONLY productvarianttranslations
    ADD CONSTRAINT fk_productvarianttranslations_productvariants FOREIGN KEY (productvariantid) REFERENCES productvariants(id) ON DELETE CASCADE;
