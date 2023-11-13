CREATE TABLE IF NOT EXISTS pages (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  title character varying(250),
  slug character varying(255),
  page_type_id uuid,
  content text,
  created_at bigint,
  metadata jsonb,
  private_metadata jsonb,
  publication_date timestamp with time zone,
  is_published boolean,
  seo_title character varying(70),
  seo_description character varying(300)
);

ALTER TABLE ONLY pages
    ADD CONSTRAINT pages_slug_key UNIQUE (slug);

CREATE INDEX idx_pages_metadata ON pages USING btree (metadata);

CREATE INDEX idx_pages_private_metadata ON pages USING btree (private_metadata);

CREATE INDEX idx_pages_slug ON pages USING btree (slug);

CREATE INDEX idx_pages_title ON pages USING btree (title);

CREATE INDEX idx_pages_title_lower_textpattern ON pages USING btree (lower((title)::text) text_pattern_ops);