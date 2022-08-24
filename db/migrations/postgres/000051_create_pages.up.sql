CREATE TABLE IF NOT EXISTS pages (
  id character varying(36) NOT NULL PRIMARY KEY,
  title character varying(250),
  slug character varying(255),
  pagetypeid character varying(36),
  content text,
  createat bigint,
  metadata text,
  privatemetadata text,
  publicationdate timestamp with time zone,
  ispublished boolean,
  seotitle character varying(70),
  seodescription character varying(300)
);

ALTER TABLE ONLY pages
    ADD CONSTRAINT pages_slug_key UNIQUE (slug);

CREATE INDEX idx_pages_metadata ON pages USING btree (metadata);

CREATE INDEX idx_pages_private_metadata ON pages USING btree (privatemetadata);

CREATE INDEX idx_pages_slug ON pages USING btree (slug);

CREATE INDEX idx_pages_title ON pages USING btree (title);

CREATE INDEX idx_pages_title_lower_textpattern ON pages USING btree (lower((title)::text) text_pattern_ops);

ALTER TABLE ONLY pages
    ADD CONSTRAINT fk_pages_pagetypes FOREIGN KEY (pagetypeid) REFERENCES pagetypes(id) ON DELETE CASCADE;

