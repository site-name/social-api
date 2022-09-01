CREATE TABLE IF NOT EXISTS vouchertranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(10),
  name character varying(255),
  voucherid character varying(36),
  createat bigint
);

ALTER TABLE ONLY vouchertranslations
    ADD CONSTRAINT vouchertranslations_languagecode_voucherid_key UNIQUE (languagecode, voucherid);
