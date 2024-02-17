CREATE TABLE IF NOT EXISTS pages (
  id varchar(36) NOT NULL PRIMARY KEY,
  title varchar(250) NOT NULL,
  slug varchar(255) NOT NULL,
  page_type_id varchar(36) NOT NULL,
  content jsonb NOT NULL,
  created_at bigint NOT NULL,
  metadata jsonb,
  private_metadata jsonb,
  publication_date timestamp with time zone,
  is_published boolean NOT NULL,
  seo_title varchar(70) NOT NULL,
  seo_description varchar(300) NOT NULL
);

ALTER TABLE ONLY pages
    ADD CONSTRAINT pages_slug_key UNIQUE (slug);

CREATE INDEX idx_pages_metadata ON pages USING btree (metadata);

CREATE INDEX idx_pages_private_metadata ON pages USING btree (private_metadata);

CREATE INDEX idx_pages_slug ON pages USING btree (slug);

CREATE INDEX idx_pages_title ON pages USING btree (title);

CREATE INDEX idx_pages_title_lower_textpattern ON pages USING btree (lower((title)::text) text_pattern_ops);