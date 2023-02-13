CREATE TABLE IF NOT EXISTS saletranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(10),
  name character varying(255),
  saleid character varying(36)
);

ALTER TABLE ONLY saletranslations
    ADD CONSTRAINT saletranslations_languagecode_saleid_key UNIQUE (languagecode, saleid);

CREATE INDEX idx_sale_translations_language_code ON saletranslations USING btree (languagecode);

CREATE INDEX idx_sale_translations_name ON saletranslations USING btree (name);
ALTER TABLE ONLY saletranslations
    ADD CONSTRAINT fk_saletranslations_sales FOREIGN KEY (saleid) REFERENCES sales(id) ON DELETE CASCADE;
