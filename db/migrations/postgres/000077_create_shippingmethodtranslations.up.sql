CREATE TABLE IF NOT EXISTS shipping_method_translations (
  id varchar(36) NOT NULL PRIMARY KEY,
  shipping_method_id varchar(36) NOT NULL,
  language_code language_code NOT NULL,
  name varchar(100) NOT NULL,
  description text NOT NULL
);

ALTER TABLE ONLY shipping_method_translations
    ADD CONSTRAINT shipping_method_translations_language_code_shipping_method_id_key UNIQUE (language_code, shipping_method_id);

CREATE INDEX idx_shipping_method_translations_name ON shipping_method_translations USING btree (name);

CREATE INDEX idx_shipping_method_translations_name_lower_text_pattern ON shipping_method_translations USING btree (lower((name)::text) text_pattern_ops);