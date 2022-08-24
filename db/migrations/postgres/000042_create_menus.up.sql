CREATE TABLE IF NOT EXISTS menus (
  id character varying(36) NOT NULL PRIMARY KEY,
  name character varying(250),
  slug character varying(255),
  createat bigint,
  metadata text,
  privatemetadata text
);

ALTER TABLE ONLY menus
    ADD CONSTRAINT menus_name_key UNIQUE (name);

ALTER TABLE ONLY menus
    ADD CONSTRAINT menus_slug_key UNIQUE (slug);

CREATE INDEX idx_menus_name ON menus USING btree (name);

CREATE INDEX idx_menus_name_lower_textpattern ON menus USING btree (lower((name)::text) text_pattern_ops);

CREATE INDEX idx_menus_slug ON menus USING btree (slug);
