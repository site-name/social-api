CREATE TABLE IF NOT EXISTS collections (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(250),
  slug character varying(255),
  background_image character varying(200),
  background_image_alt character varying(128),
  description text,
  metadata jsonb,
  private_metadata jsonb,
  seo_title character varying(70),
  seo_description character varying(300)
);

CREATE INDEX idx_collections_name ON collections USING btree (name);

CREATE INDEX idx_collections_name_lower_textpattern ON collections USING btree (lower((name)::text) text_pattern_ops);