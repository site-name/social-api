CREATE TABLE IF NOT EXISTS categorytranslations (
  id character varying(36) NOT NULL PRIMARY KEY,
  languagecode character varying(5),
  categoryid character varying(36),
  name character varying(250),
  description text,
  seotitle character varying(70),
  seodescription character varying(300)
);

CREATE INDEX idx_category_translations_name ON categorytranslations USING btree (name);

CREATE INDEX idx_category_translations_name_lower_textpattern ON categorytranslations USING btree (lower((name)::text) text_pattern_ops);
