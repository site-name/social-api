CREATE TABLE IF NOT EXISTS sale_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code character varying(10),
  name character varying(255),
  sale_id character varying(36)
);

ALTER TABLE ONLY sale_translations
    ADD CONSTRAINT sale_translations_language_code_sale_id_key UNIQUE (language_code, sale_id);

CREATE INDEX idx_sale_translations_language_code ON sale_translations USING btree (language_code);

CREATE INDEX idx_sale_translations_name ON sale_translations USING btree (name);