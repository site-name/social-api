CREATE TABLE IF NOT EXISTS page_types (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  slug character varying(255),
  metadata jsonb,
  privatemetadata jsonb
);

ALTER TABLE ONLY page_types
    ADD CONSTRAINT page_types_slug_key UNIQUE (slug);

CREATE INDEX idx_page_types_name ON page_types USING btree (name);

CREATE INDEX idx_page_types_name_lower_textpattern ON page_types USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_page_types_slug ON page_types USING btree (slug);
