CREATE TABLE IF NOT EXISTS menus (
  id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name character varying(250) NOT NULL,
  slug character varying(255) NOT NULL,
  created_at bigint NOT NULL,
  metadata jsonb,
  private_metadata jsonb
);

ALTER TABLE ONLY menus
    ADD CONSTRAINT menus_name_key UNIQUE (name);

ALTER TABLE ONLY menus
    ADD CONSTRAINT menus_slug_key UNIQUE (slug);

CREATE INDEX idx_menus_name ON menus USING btree (name);

CREATE INDEX idx_menus_name_lower_textpattern ON menus USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_menus_slug ON menus USING btree (slug);