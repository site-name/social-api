CREATE TABLE IF NOT EXISTS collections (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(250) NOT NULL,
  slug character varying(255) NOT NULL,
  background_image character varying(200),
  background_image_alt character varying(128) NOT NULL,
  description text,
  metadata jsonb,
  private_metadata jsonb,
  seo_title character varying(70) NOT NULL,
  seo_description character varying(300) NOT NULL
);

CREATE INDEX idx_collections_name ON collections USING btree (name);

CREATE INDEX idx_collections_name_lower_textpattern ON collections USING btree (lower((name)::text) text_pattern_ops);