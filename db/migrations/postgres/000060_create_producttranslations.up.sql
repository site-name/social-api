CREATE TABLE IF NOT EXISTS producttranslations (
id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  productid character varying(36),
  name character varying(250),
  description text,
  seotitle character varying(70),
  seodescription character varying(300)
);

ALTER TABLE ONLY producttranslations
    ADD CONSTRAINT producttranslations_languagecode_productid_key UNIQUE (languagecode, productid);

ALTER TABLE ONLY producttranslations
    ADD CONSTRAINT producttranslations_name_key UNIQUE (name);
ALTER TABLE ONLY producttranslations
    ADD CONSTRAINT fk_producttranslations_products FOREIGN KEY (productid) REFERENCES products(id) ON DELETE CASCADE;
