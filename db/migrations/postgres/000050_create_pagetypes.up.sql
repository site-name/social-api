CREATE TABLE IF NOT EXISTS pagetypes (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  slug character varying(255),
  metadata text,
  privatemetadata text
);

ALTER TABLE ONLY pagetypes
    ADD CONSTRAINT pagetypes_slug_key UNIQUE (slug);

CREATE INDEX idx_page_types_name ON pagetypes USING btree (name);

CREATE INDEX idx_page_types_name_lower_textpattern ON pagetypes USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_page_types_slug ON pagetypes USING btree (slug);
