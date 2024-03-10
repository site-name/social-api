CREATE TABLE IF NOT EXISTS attribute_value_translations (
  id varchar(36) NOT NULL PRIMARY KEY,
  language_code language_code NOT NULL,
  attribute_value_id varchar(36) NOT NULL,
  name varchar(250) NOT NULL,
  rich_text text,
  plain_text text
);

ALTER TABLE ONLY attribute_value_translations
    ADD CONSTRAINT attribute_value_translations_language_code_attribute_value_id_key UNIQUE (language_code, attribute_value_id);

CREATE INDEX IF NOT EXISTS idx_attribute_value_translations_name ON attribute_value_translations USING btree (name);

CREATE INDEX IF NOT EXISTS idx_attribute_value_translations_name_lower_textpattern ON attribute_value_translations USING btree (lower((name)::text) text_pattern_ops);
