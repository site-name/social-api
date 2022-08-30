CREATE TABLE IF NOT EXISTS attributetranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  attributeid character varying(36),
  languagecode character varying(5),
  name character varying(100)
);

ALTER TABLE ONLY attributetranslations
    ADD CONSTRAINT attributetranslations_languagecode_attributeid_key UNIQUE (languagecode, attributeid);

CREATE INDEX IF NOT EXISTS idx_attributetranslations_name ON attributetranslations USING btree (name);

CREATE INDEX IF NOT EXISTS idx_attributetranslations_name_lower_textpattern ON attributetranslations USING btree (lower((name)::text) text_pattern_ops);

