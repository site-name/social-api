CREATE TABLE IF NOT EXISTS attribute_translations (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  languagecode character varying(5),
  name character varying(100)
);

ALTER TABLE ONLY attribute_translations
    ADD CONSTRAINT attribute_translations_languagecode_attributeid_key UNIQUE (languagecode, attributeid);

CREATE INDEX IF NOT EXISTS idx_attribute_translations_name ON attribute_translations USING btree (name);

CREATE INDEX IF NOT EXISTS idx_attribute_translations_name_lower_textpattern ON attribute_translations USING btree (lower((name)::text) text_pattern_ops);

