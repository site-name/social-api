CREATE TABLE IF NOT EXISTS attributevaluetranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  attributevalueid character varying(36),
  name character varying(100),
  richtext text
);

ALTER TABLE ONLY attributevaluetranslations
    ADD CONSTRAINT attributevaluetranslations_languagecode_attributevalueid_key UNIQUE (languagecode, attributevalueid);

CREATE INDEX IF NOT EXISTS idx_attribute_value_translations_name ON attributevaluetranslations USING btree (name);

CREATE INDEX IF NOT EXISTS idx_attribute_value_translations_name_lower_textpattern ON attributevaluetranslations USING btree (lower((name)::text) text_pattern_ops);
