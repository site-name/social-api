CREATE TABLE IF NOT EXISTS page_types (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(250) NOT NULL,
  slug character varying(255)NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY page_types
    ADD CONSTRAINT page_types_slug_key UNIQUE (slug);

CREATE INDEX idx_page_types_name ON page_types USING btree (name);

CREATE INDEX idx_page_types_name_lower_textpattern ON page_types USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_page_types_slug ON page_types USING btree (slug);