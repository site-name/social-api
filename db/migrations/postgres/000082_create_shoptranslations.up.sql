CREATE TABLE IF NOT EXISTS shoptranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  shopid character varying(36),
  languagecode character varying(5),
  name character varying(110),
  description character varying(110),
  createat bigint,
  updateat bigint
);

ALTER TABLE ONLY shoptranslations
    ADD CONSTRAINT shoptranslations_languagecode_shopid_key UNIQUE (languagecode, shopid);

ALTER TABLE ONLY shoptranslations
    ADD CONSTRAINT fk_shoptranslations_shops FOREIGN KEY (shopid) REFERENCES shops(id) ON DELETE CASCADE;
