CREATE TABLE IF NOT EXISTS product_variant_translations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  productvariantid character varying(36),
  name character varying(255)
);

ALTER TABLE ONLY product_variant_translations
    ADD CONSTRAINT product_variant_translations_languagecode_productvariantid_key UNIQUE (languagecode, productvariantid);
