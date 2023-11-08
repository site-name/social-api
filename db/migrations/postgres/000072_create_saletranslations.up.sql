CREATE TABLE IF NOT EXISTS sale_translations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(10),
  name character varying(255),
  saleid character varying(36)
);

ALTER TABLE ONLY sale_translations
    ADD CONSTRAINT sale_translations_languagecode_saleid_key UNIQUE (languagecode, saleid);

CREATE INDEX idx_sale_translations_language_code ON sale_translations USING btree (languagecode);

CREATE INDEX idx_sale_translations_name ON sale_translations USING btree (name);
