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
