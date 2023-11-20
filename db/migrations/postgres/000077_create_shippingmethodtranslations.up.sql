CREATE TABLE IF NOT EXISTS shipping_method_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  shipping_method_id uuid NOT NULL,
  language_code varchar(10) NOT NULL,
  name varchar(100) NOT NULL,
  description text NOT NULL
);

ALTER TABLE ONLY shipping_method_translations
    ADD CONSTRAINT shipping_method_translations_language_code_shipping_method_id_key UNIQUE (language_code, shipping_method_id);

CREATE INDEX idx_shipping_method_translations_name ON shipping_method_translations USING btree (name);

CREATE INDEX idx_shipping_method_translations_name_lower_text_pattern ON shipping_method_translations USING btree (lower((name)::text) text_pattern_ops);