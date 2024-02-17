CREATE TABLE IF NOT EXISTS collections (
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(250) NOT NULL,
  slug varchar(255) NOT NULL,
  background_image varchar(200),
  background_image_alt varchar(128) NOT NULL,
  description text,
  metadata jsonb,
  private_metadata jsonb,
  seo_title varchar(70) NOT NULL,
  seo_description varchar(300) NOT NULL
);

CREATE INDEX idx_collections_name ON collections USING btree (name);

CREATE INDEX idx_collections_name_lower_textpattern ON collections USING btree (lower((name)::text) text_pattern_ops);