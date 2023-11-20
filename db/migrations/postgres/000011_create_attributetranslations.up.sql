CREATE TABLE IF NOT EXISTS attribute_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  attribute_id uuid NOT NULL,
  language_code varchar(10) NOT NULL,
  name varchar(100) NOT NULL
);

ALTER TABLE ONLY attribute_translations
    ADD CONSTRAINT attribute_translations_language_code_attribute_id_key UNIQUE (language_code, attribute_id);

CREATE INDEX IF NOT EXISTS idx_attribute_translations_name ON attribute_translations USING btree (name);

CREATE INDEX IF NOT EXISTS idx_attribute_translations_name_lower_textpattern ON attribute_translations USING btree (lower((name)::text) text_pattern_ops);
