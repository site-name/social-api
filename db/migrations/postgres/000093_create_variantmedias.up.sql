CREATE TABLE IF NOT EXISTS variantmedias (
  id character varying(36) NOT NULL PRIMARY KEY,
  variantid character varying(36),
  mediaid character varying(36)
);

ALTER TABLE ONLY variantmedias
    ADD CONSTRAINT variantmedias_variantid_mediaid_key UNIQUE (variantid, mediaid);
ALTER TABLE ONLY variantmedias
    ADD CONSTRAINT fk_variantmedias_productmedias FOREIGN KEY (mediaid) REFERENCES productmedias(id) ON DELETE CASCADE;
ALTER TABLE ONLY variantmedias
    ADD CONSTRAINT fk_variantmedias_productvariants FOREIGN KEY (variantid) REFERENCES productvariants(id) ON DELETE CASCADE;
