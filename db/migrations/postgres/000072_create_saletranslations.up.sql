CREATE TABLE IF NOT EXISTS sale_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  language_code varchar(10) NOT NULL,
  name varchar(255) NOT NULL,
  sale_id uuid NOT NULL
);

ALTER TABLE ONLY sale_translations
    ADD CONSTRAINT sale_translations_language_code_sale_id_key UNIQUE (language_code, sale_id);

CREATE INDEX idx_sale_translations_language_code ON sale_translations USING btree (language_code);

CREATE INDEX idx_sale_translations_name ON sale_translations USING btree (name);