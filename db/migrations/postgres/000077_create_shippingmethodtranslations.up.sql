CREATE TABLE IF NOT EXISTS shipping_method_translations (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  shipping_method_id uuid,
  language_code character varying(5),
  name character varying(100),
  description text
);

ALTER TABLE ONLY shipping_method_translations
    ADD CONSTRAINT shipping_method_translations_language_code_shipping_method_id_key UNIQUE (language_code, shipping_method_id);

CREATE INDEX idx_shipping_method_translations_name ON shipping_method_translations USING btree (name);

CREATE INDEX idx_shipping_method_translations_name_lower_text_pattern ON shipping_method_translations USING btree (lower((name)::text) text_pattern_ops);